package config

import (
	"time"
)

type Config struct {
	Name      string  `mapstructure:"name"`
	Version   string  `mapstructure:"version"`
	Banner    string  `mapstructure:"banner"`
	Port      int     `mapstructure:"port"`
	Directory string  `mapstructure:"directory"`
	Memory    Memory  `mapstructure:"memory"`
	Etcd      Etcd    `mapstructure:"etcd"`
	Cluster   Cluster `mapstructure:"cluster"`
}

type Memory struct {
	ExpirationTime  time.Duration `mapstructure:"expiration-time"`  // 有效期
	CleanupInterval time.Duration `mapstructure:"cleanup-Interval"` // 清理间隔
}

type Etcd struct {
	Endpoints      []string `mapstructure:"endpoints"`
	DialTimeout    int      `mapstructure:"dial-timeout"`
	RequestTimeout int      `mapstructure:"request-timeout"`
}

type Cluster struct {
	HeartbeatInterval        int `mapstructure:"heartbeat-interval"`
	HeartbeatRecheckInterval int `mapstructure:"heartbeat-recheck-interval"`
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
	Name:      "Magpie",
	Version:   "1.0.0",
	Banner:    banner_bulbhead,
	Port:      8843,
	Directory: "@magpie",
	Memory: Memory{
		ExpirationTime:  5 * time.Minute,
		CleanupInterval: 10 * time.Minute,
	},
	Etcd: Etcd{
		Endpoints:      []string{"127.0.0.1:2379"},
		DialTimeout:    5,
		RequestTimeout: 5,
	},
	Cluster: Cluster{
		HeartbeatInterval:        9,
		HeartbeatRecheckInterval: 5,
	},
}

var globalConf = defaultConf

func init() {

}

func GetConfig() *Config {
	return &globalConf
}
