package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	App  AppConfig  `toml:"app"`
	Log  LogConfig  `toml:"log"`
	Bili BiliConfig `toml:"bili"`
}

// AppConfig 应用程序基本配置
type AppConfig struct {
	Name      string `toml:"name"`       // 应用名称
	Version   string `toml:"version"`    // 应用版本
	ServeAddr string `toml:"serve_addr"` // 服务端口
}

// LogConfig 日志配置
type LogConfig struct {
	MemoryMaxLogLine int    `toml:"memory_max_log_line"` // 最大内存日志行数
	LogLevel         int    `toml:"log_level"`           // 日志级别 (例如: 0=info, 1=error)
	LogPath          string `toml:"log_path"`            // 日志存储路径
}

// BiliConfig B站下载相关配置
type BiliConfig struct {
	Cookie    string `toml:"cookie"`     // B站账号的Cookie值
	SavePath  string `toml:"save_path"`  // 下载保存路径
	SaveCover bool   `toml:"save_cover"` // 下载视频的同时是否保存封面
}

func SaveConfigToTomlFile() error {
	globalCfgBytes, _ := toml.Marshal(GlobalConfig)
	return os.WriteFile("config.toml", globalCfgBytes, 0644)
}
