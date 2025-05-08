package config

import "time"

type Config struct {
	Env                string `yaml:"env" env-default:"development"`
	HTTPServer         `yaml:"http_server"`
	DbConnectionString string `yaml:"postgres_connection_string" env-required:"true"`
	IndexFilePath      string `yaml:"index_file_path" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
