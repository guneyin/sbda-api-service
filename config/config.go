package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const defaultHttpPort = 8080

type Config struct {
	ApiName          string `env:"API_NAME"`
	DiscoverySvcAddr string `env:"DISCOVERY_SVC_ADDR"`
	ApiVersion       string `env:"API_VERSION"`
	HttpPort         int    `env:"HTTP_PORT"`
	NetworkAlias     string `env:"NETWORK_ALIAS"`
}

func GetConfig() *Config {
	var config Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return nil
	}

	if config.HttpPort == 0 {
		config.HttpPort = defaultHttpPort
	}

	return &config
}
