package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" env-default:"9000"`

	PostgresURL     string `env:"PG_URL" env-required:"true"`
	PostgresMaxCons int    `env:"PG_MAX_CONS" env-default:"15"`

	OccupiesExpireTime int `env:"OCCUPIES_EXPIRE_TIME" env-default:"5"`

	ServeSwagger bool   `env:"SERVE_SWAGGER" env-default:"true"`
	LogLevel     string `env:"LOG_LEVEL"     env-default:"info"`

	initialized bool
}

var config Config

func NewConfig() Config {
	if config.initialized {
		return config
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatal(err)
	}
	return config
}
