package initialize

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	HTTPHost        string `env:"HTTP_HOST" envDefault:"localhost"`
	HTTPPort        string `env:"HTTP_PORT" envDefault:"3000"`
	PGMaxAttemption int    `env:"PG_MAX_ATTEMPTION" envDefault:"5"`
	PGHost          string `env:"PG_HOST" envDefault:"localhost"`
	PGPort          string `env:"PG_PORT" envDefault:"5432"`
	PGUser          string `env:"PG_USER" envDefault:"postgres"`
	PGPassword      string `env:"PG_PASSWORD" envDefault:"22578"`
	PGDatabase      string `env:"PG_DATABASE" envDefault:"urlshortener"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
