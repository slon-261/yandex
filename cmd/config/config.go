package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	DefaultRunAddr string
	DefaultBaseUrl string
	EnvRunAddr     string `env:"SERVER_ADDRESS"`
	EnvBaseUrl     string `env:"BASE_URL"`
}

func NewConfig() *Config {

	var cfg Config
	env.Parse(&cfg)

	cfg.DefaultRunAddr = ":8080"
	cfg.DefaultBaseUrl = "http://localhost:8080"

	return &cfg
}
