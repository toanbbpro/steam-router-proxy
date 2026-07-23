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

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.SetOutput(io.Discard)

	// FIX: Nếu Service đang chạy, App GUI không bật Proxy cục bộ để tránh xung đột Port
	if !IsServiceRunning() {
		if err := updateHosts(true); err == nil {
			go runProxy()
		}
	}
}

func (a *App) ToggleLogging(enable bool) string {
	a.logging = enable
	if enable {
		exePath, _ := os.Executable()
		logDir := filepath.Join(filepath.Dir(exePath), "logs")
		rotator := &LogRotator{LogDir: logDir, MaxSize: 8 * 1024 * 1024}
		log.SetOutput(rotator)
		log.Println("[INIT] --- Bắt đầu ghi log / Start logging ---")
		return "Đã bật ghi Log / Logging enabled."
	} else {
		log.SetOutput(io.Discard)
		return "Đã tắt ghi Log / Logging disabled."
	}
}

func (a *App) ManageService(action string) string {
	exePath, _ := os.Executable()
	svcConfig := &service.Config{
		Name:        "SteamRouterService",
		DisplayName: "Steam Network Router Service",
		Description: "Dịch vụ chạy ngầm tự động định tuyến và tăng tốc kết nối Steam.",
		Executable:  exePath,
		Arguments:   []string{"-service"},
		Option: service.KeyValue{
			"StartType": "automatic",
		},
	}
	s, err := service.New(&program{}, svcConfig)
	if err != nil {
		return "Lỗi khởi tạo / Init error: " + err.Error()
	}

	switch action {
	case "install":
		// FIX: Tắt Proxy của GUI trước khi nhường Port 443 cho Service
		stopProxy()

		_ = service.Control(s, "stop")
		_ = service.Control(s, "uninstall")
		if err := service.Control(s, "install"); err != nil {
			go runProxy() // Bật lại GUI proxy nếu lỗi
			return "Lỗi cài đặt / Install error: " + err.Error()
		}
		if err := service.Control(s, "start"); err != nil {
			return "Lỗi khởi chạy / Start error: " + err.Error()
		}
		return "Thành công: Đã cài Service ngầm! / Success: Service installed!"
	case "uninstall":
		_ = service.Control(s, "stop")
		if err := service.Control(s, "uninstall"); err != nil {
			return "Lỗi gỡ cài đặt / Uninstall error: " + err.Error()
		}
		// FIX: Lấy lại quyền chạy Proxy về cho GUI sau khi gỡ Service
		updateHosts(true)
		go runProxy()
		return "Thành công: Đã gỡ Service! / Success: Service uninstalled!"
	}
	return "Lệnh không hợp lệ / Invalid action."
}

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
		log.Println("[ERROR] Port 443 đã bị chiếm / Port in use:", err)
		return
	}
	log.Println("[+] Proxy đang hoạt động / Proxy running at 127.0.0.1:443")

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

	targetIP, err := resolveDoH(host)
	if err != nil || targetIP == "" {
		return
	}

	targetConn, err := net.DialTimeout("tcp", net.JoinHostPort(targetIP, "443"), 5*time.Second)
	if err != nil {
		return
	}
	defer targetConn.Close()

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

// FIX: Hàm updateHosts chống trùng lặp dữ liệu
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
		// Lọc bỏ 100% các dòng cũ do app tạo ra để tránh trùng lặp
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

// Mở trang chủ GitHub của dự án
func (a *App) OpenHomePage() {
	exec.Command("cmd", "/c", "start", "https://github.com/toanbbpro/steam-router-proxy").Start()
}
