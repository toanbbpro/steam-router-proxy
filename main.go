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
		log.Fatalf("Không thể mở ứng dụng với quyền Admin: %v", err)
	}
	os.Exit(0)
}

// Logic thực thi khi ứng dụng chạy dưới dạng Windows Service ngầm
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

	// Trường hợp khởi chạy bởi Windows Service Manager
	if len(os.Args) > 1 && os.Args[1] == "-service" {
		svcConfig := &service.Config{
			Name:        "SteamRouterService",
			DisplayName: "Steam Network Router Service",
			Description: "Dịch vụ chạy ngầm tự động định tuyến và tăng tốc kết nối Steam.",
		}
		prg := &program{}
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}
		s.Run()
		return
	}

	// Trường hợp người dùng mở giao diện Wails GUI
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Steam Network Router",
		Width:  420,
		Height: 280, // Tối ưu kích thước giao diện gọn hơn
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 40, B: 56, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			// Tự động dọn dẹp khi đóng cửa sổ GUI
			stopProxy()
			updateHosts(false)
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
