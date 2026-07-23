package main

import (
	"context"
	"embed"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/kardianos/service"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.org/x/sys/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func isElevated() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY, 2,
		windows.SECURITY_BUILTIN_DOMAIN_RID, windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0, &sid)
	if err != nil {
		return false
	}
	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

func elevate() {
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verb := syscall.StringToUTF16Ptr("runas")
	exePtr := syscall.StringToUTF16Ptr(exe)
	cwdPtr := syscall.StringToUTF16Ptr(cwd)
	argPtr := syscall.StringToUTF16Ptr(args)

	err := windows.ShellExecute(0, verb, exePtr, argPtr, cwdPtr, windows.SW_NORMAL)
	if err != nil {
		log.Fatalf("Không thể mở ứng dụng quyền Admin / Cannot run as Admin: %v", err)
	}
	os.Exit(0)
}

// FIX: Hàm kiểm tra xem Windows Service của app có đang hoạt động không
func IsServiceRunning() bool {
	svcConfig := &service.Config{Name: "SteamRouterService"}
	s, err := service.New(&program{}, svcConfig)
	if err != nil {
		return false
	}
	status, err := s.Status()
	if err != nil {
		return false
	}
	return status == service.StatusRunning
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go func() {
		updateHosts(true)
		runProxy()
	}()
	return nil
}

func (p *program) Stop(s service.Service) error {
	stopProxy()
	updateHosts(false)
	return nil
}

func main() {
	if !isElevated() {
		elevate()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "-service" {
		svcConfig := &service.Config{
			Name: "SteamRouterService",
		}
		prg := &program{}
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}
		s.Run()
		return
	}

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Steam Router Proxy v0.3",
		Width:  480, // TĂNG KÍCH THƯỚC: Mở rộng chiều ngang
		Height: 360, // TĂNG TỪ 320 LÊN 360: Giúp hiển thị thoải mái không lo bị cấn
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 40, B: 56, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			// FIX: Chỉ khôi phục cài đặt nếu Service KHÔNG chạy
			// Nếu Service đang chạy thì im lặng đóng app, nhường service lo liệu
			if !IsServiceRunning() {
				stopProxy()
				updateHosts(false)
			}
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
