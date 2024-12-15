package config

import (
	"github.com/lxn/win"
	"os"
	"os/signal"
	"syscall"
)

const (
	AppIconPath = "icon128.ico"
	RegeditKey  = "MyOSManage"
)

var (
	ProcessWorkDir = "./"
	GlobalQuit     chan os.Signal
	MainPanelHWND  win.HWND // 主窗口句柄
)

func init() {
	pwd, err := os.Getwd()
	if err == nil {
		ProcessWorkDir = pwd
	}

	GlobalQuit = make(chan os.Signal, 1)
	signal.Notify(GlobalQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}
