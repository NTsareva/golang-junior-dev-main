package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Log    LogConfig    `toml:"log"`
	Server ServerConfig `toml:"servers"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

type ServerConfig struct {
	Port     string `toml:"port"`
	Hostname string `toml:"hostname"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return Config{}, fmt.Errorf("can't parse log toml file: %v", err)
	}

	return config, nil
}
