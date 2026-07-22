package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kardianos/service"
)

var (
	steamDomains = []string{
		"steamcommunity.com", "store.steampowered.com", "checkout.steampowered.com",
		"login.steampowered.com", "help.steampowered.com", "community.cloudflare.steamstatic.com",
		"steamcommunity-a.akamaihd.net",
	}
	hostsPath    = `C:\Windows\System32\drivers\etc\hosts`
	coreListener net.Listener
)

type App struct {
	ctx     context.Context
	logging bool
}

func NewApp() *App {
	return &App{}
}

// Khi GUI mở lên -> Tự động khởi chạy Proxy ngay lập tức
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.SetOutput(io.Discard) // Mặc định không ghi log
	if err := updateHosts(true); err == nil {
		go runProxy()
	}
}

// Bật/Tắt Ghi Log từ Frontend
func (a *App) ToggleLogging(enable bool) string {
	a.logging = enable
	if enable {
		exePath, _ := os.Executable()
		logDir := filepath.Join(filepath.Dir(exePath), "logs")
		rotator := &LogRotator{LogDir: logDir, MaxSize: 8 * 1024 * 1024} // 8MB
		log.SetOutput(rotator)
		log.Println("[INIT] --- Bắt đầu ghi log phiên làm việc ---")
		return "Đã bật ghi Log."
	} else {
		log.SetOutput(io.Discard)
		return "Đã tắt ghi Log."
	}
}

// Quản lý Windows Service chạy ngầm
func (a *App) ManageService(action string) string {
	exePath, _ := os.Executable()
	svcConfig := &service.Config{
		Name:        "SteamRouterService",
		DisplayName: "Steam Network Router Service",
		Description: "Dịch vụ chạy ngầm tự động định tuyến và tăng tốc kết nối Steam.",
		Executable:  exePath,
		Arguments:   []string{"-service"},
		Option: service.KeyValue{
			"StartType": "automatic", // Tự động chạy cùng Windows
		},
	}
	s, err := service.New(&program{}, svcConfig)
	if err != nil {
		return "Lỗi khởi tạo Service: " + err.Error()
	}

	switch action {
	case "install":
		_ = service.Control(s, "stop")
		_ = service.Control(s, "uninstall")
		if err := service.Control(s, "install"); err != nil {
			return "Lỗi cài đặt Service: " + err.Error()
		}
		if err := service.Control(s, "start"); err != nil {
			return "Đã cài Service nhưng chưa bật được: " + err.Error()
		}
		return "Đã cài đặt & kích hoạt Service tự chạy cùng Windows!"
	case "uninstall":
		_ = service.Control(s, "stop")
		if err := service.Control(s, "uninstall"); err != nil {
			return "Lỗi gỡ cài đặt: " + err.Error()
		}
		return "Đã gỡ bỏ Service khỏi hệ thống."
	}
	return "Lệnh không hợp lệ."
}

// Mở thư mục chứa File Log
func (a *App) OpenLogFolder() {
	exePath, _ := os.Executable()
	logDir := filepath.Join(filepath.Dir(exePath), "logs")
	os.MkdirAll(logDir, 0755)
	exec.Command("explorer", logDir).Start()
}

// -----------------------------------------
// CORE PROXY ENGINE
// -----------------------------------------

func runProxy() {
	var err error
	coreListener, err = net.Listen("tcp", "127.0.0.1:443")
	if err != nil {
		log.Println("[ERROR] Port 443 đã bị chiếm dụng hoặc không đủ quyền:", err)
		return
	}
	log.Println("[+] Steam Router Proxy đang hoạt động tại 127.0.0.1:443")

	for {
		conn, err := coreListener.Accept()
		if err != nil {
			break
		}
		go handleConnection(conn)
	}
}

func stopProxy() {
	if coreListener != nil {
		coreListener.Close()
		coreListener = nil
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 4096)
	n, err := clientConn.Read(buf)
	if err != nil || n == 0 {
		return
	}
	clientConn.SetReadDeadline(time.Time{})

	host := parseSNI(buf[:n])
	if host == "" {
		host = "steamcommunity.com"
	}

	log.Printf("[REQ] Yêu cầu kết nối Domain: %s", host)

	targetIP, err := resolveDoH(host)
	if err != nil || targetIP == "" {
		log.Printf("[ERR] DoH Phân giải thất bại: %s (%v)", host, err)
		return
	}

	log.Printf("[DOH] Phân giải thành công: %s => %s", host, targetIP)

	targetConn, err := net.DialTimeout("tcp", net.JoinHostPort(targetIP, "443"), 5*time.Second)
	if err != nil {
		log.Printf("[ERR] Không thể kết nối IP mục tiêu: %s (%s)", host, targetIP)
		return
	}
	defer targetConn.Close()

	log.Printf("[OK] Đang chuyển tiếp dữ liệu: %s", host)

	targetConn.Write(buf[:1])
	time.Sleep(3 * time.Millisecond)
	targetConn.Write(buf[1:n])

	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}

func resolveDoH(domain string) (string, error) {
	apiURL := "https://cloudflare-dns.com/dns-query?name=" + url.QueryEscape(domain) + "&type=A"
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Accept", "application/dns-json")
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Answer []struct {
			Data string `json:"data"`
		}
	}
	json.NewDecoder(resp.Body).Decode(&result)

	for _, ans := range result.Answer {
		if !strings.HasPrefix(ans.Data, "127.") && !strings.Contains(ans.Data, ":") {
			return ans.Data, nil
		}
	}
	return "", os.ErrNotExist
}

func parseSNI(data []byte) string {
	if len(data) < 5 || data[0] != 0x16 {
		return ""
	}
	for i := 0; i < len(data)-4; i++ {
		if data[i] == 0x00 && data[i+1] == 0x00 {
			length := int(data[i+2])<<8 | int(data[i+3])
			if i+4+length <= len(data) {
				sniData := data[i+4 : i+4+length]
				if len(sniData) > 5 && sniData[2] == 0 {
					nameLen := int(sniData[3])<<8 | int(sniData[4])
					if 5+nameLen <= len(sniData) {
						return string(sniData[5 : 5+nameLen])
					}
				}
			}
		}
	}
	return ""
}

func updateHosts(add bool) error {
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	strContent := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(strContent, "\n")

	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, "# SteamRouter") && trimmed != "" {
			newLines = append(newLines, line)
		}
	}

	if add {
		for _, domain := range steamDomains {
			newLines = append(newLines, "127.0.0.1 "+domain+" # SteamRouter")
		}
	}

	output := strings.Join(newLines, "\r\n") + "\r\n"
	err = os.WriteFile(hostsPath, []byte(output), 0644)
	exec.Command("ipconfig", "/flushdns").Run()
	return err
}

// -----------------------------------------
// CUSTOM LOG ROTATOR (File proxy-log-timestamp.log)
// -----------------------------------------

type LogRotator struct {
	LogDir  string
	MaxSize int64
	mu      sync.Mutex
	file    *os.File
	size    int64
}

func (r *LogRotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.file == nil {
		if err := r.openNew(); err != nil {
			return 0, err
		}
	}

	if r.size+int64(len(p)) > r.MaxSize {
		r.file.Close()
		if err := r.openNew(); err != nil {
			return 0, err
		}
	}

	n, err = r.file.Write(p)
	r.size += int64(n)
	return n, err
}

func (r *LogRotator) openNew() error {
	os.MkdirAll(r.LogDir, 0755)
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("proxy-log-%s.log", timestamp)
	filePath := filepath.Join(r.LogDir, fileName)

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	r.file = f
	r.size = 0
	return nil
}
