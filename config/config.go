package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Port            string `env:"PORT" env-default:"9000"`
	PostgresURL     string `env:"PG_URL" env-required:"true"`
	PostgresMaxCons int    `env:"PG_MAX_CONS" env-default:"15"`
	LogLevel        string `env:"LOG_LEVEL" env-default:"info"`
}

func NewConfig() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
