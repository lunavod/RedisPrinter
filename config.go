package main

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type RedisConfig struct {
	IP       string
	Port     int
	Password string
	Channel  string
	Database int
}

type MainConfig struct {
	UploadsDir string
}

type Config struct {
	Redis RedisConfig
	Main  MainConfig
}

func GetConfig() Config {
	dat, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic("Config file not found")
	}

	config := Config{}
	_ = toml.Unmarshal(dat, &config)

	if config.Redis.IP == "" {
		panic("Redis IP not specified")
	}

	if config.Redis.Port == 0 {
		panic("Redis port not specified")
	}

	if config.Redis.Channel == "" {
		panic("Redis channel not specified")
	}

	return config
}
