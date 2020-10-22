package config

import (
	"time"
)

type Config struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Memory  Memory `mapstructure:"memory"`
}

type Memory struct {
	ExpirationTime  time.Duration `mapstructure:"expiration-time"`  // 有效期
	CleanupInterval time.Duration `mapstructure:"cleanup-Interval"` // 清理间隔
}

var defaultConf = Config{
	Name:    "Magpie",
	Version: "1.0.0",
	Memory: Memory{
		ExpirationTime:  5 * time.Minute,
		CleanupInterval: 10 * time.Minute,
	},
}

var globalConf = defaultConf

func init() {

}

func GetConfig() *Config {
	return &globalConf
}
