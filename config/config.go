package config

import (
	"time"
)

type Config struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Banner  string `mapstructure:"banner"`
	Port    int    `mapstructure:"port"`
	Memory  Memory `mapstructure:"memory"`
}

type Memory struct {
	ExpirationTime  time.Duration `mapstructure:"expiration-time"`  // 有效期
	CleanupInterval time.Duration `mapstructure:"cleanup-Interval"` // 清理间隔
}

const (
	banner_bulbhead = `
	 __  __    __    ___  ____  ____  ____ 
	(  \/  )  /__\  / __)(  _ \(_  _)( ___)
	 )    (  /(__)\( (_-. )___/ _)(_  )__) 
	(_/\/\_)(__)(__)\___/(__)  (____)(____)
	
	`
)

var defaultConf = Config{
	Name:    "Magpie",
	Version: "1.0.0",
	Banner:  banner_bulbhead,
	Port:    8843,
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
