package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	DefaultRunAddr string
	DefaultBaseURL string
	EnvRunAddr     string `env:"SERVER_ADDRESS"`
	EnvBaseURL     string `env:"BASE_URL"`
}

func NewConfig() *Config {

	var cfg Config
	env.Parse(&cfg)

	cfg.DefaultRunAddr = ":8080"
	cfg.DefaultBaseURL = "http://localhost:8080"

	return &cfg
}
