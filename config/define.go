package config

import (
	"github.com/BurntSushi/toml"
	"github.com/lxn/win"
	"os"
	"os/signal"
	"syscall"
)

const (
	AppName     = "os_manage"
	AppVersion  = "0.1.1"
	DefaultAddr = ":7799"
	AppIconPath = "icon128.ico"
	RegeditKey  = "MyOSManage"
)

var (
	ProcessWorkDir = "./"
	GlobalQuit     chan os.Signal
	MainPanelHWND  win.HWND // 主窗口句柄
	GlobalConfig   Config
)

func init() {
	pwd, err := os.Getwd()
	if err == nil {
		ProcessWorkDir = pwd
	}

	GlobalQuit = make(chan os.Signal, 1)
	signal.Notify(GlobalQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	LoadConfig()
}

func LoadConfig() {
	configBytes, err := os.ReadFile("config.toml")
	if err == nil {
		err = toml.Unmarshal(configBytes, &GlobalConfig)
		if err != nil {
			panic(err)
		}
	} else {
		// 配置文件不存在则创建
		GlobalConfig = Config{
			App: AppConfig{
				Name:      AppName,
				Version:   AppVersion,
				ServeAddr: DefaultAddr,
			},
			Log: LogConfig{
				MemoryMaxLogLine: 128,
				LogLevel:         1,
				LogPath:          "log/",
			},
			Bili: BiliConfig{
				Cookie:    "",
				SavePath:  "./",
				SaveCover: false,
			},
		}
		err = SaveConfigToTomlFile()
		if err != nil {
			panic("generate config.toml failed: " + err.Error())
		}
	}
}
