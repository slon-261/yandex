package config

type Config struct {
	RunAddr  string
	BaseAddr string
}

func NewConfig() *Config {

	return &Config{
		RunAddr:  ":8080",
		BaseAddr: "http://localhost:8080",
	}
}
