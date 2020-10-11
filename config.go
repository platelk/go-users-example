package main

import (
	"github.com/ilyakaznacheev/cleanenv"

	"go-users-example/infra/logger"
	"go-users-example/transport/http"
)

// Config will hold all the based the configuration for the app
type Config struct {
	HTTP   http.Config   `env:"HTTP"`
	Logger logger.Config `env:"LOG"`
}

// Load will retrieve the configuration from different sources by order of priority `flag > ENV > file`
func Load() *Config {
	c := &Config{}

	_ = cleanenv.ReadEnv(c)

	return c
}
