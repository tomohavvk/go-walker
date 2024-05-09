package configs

import (
	"github.com/caarlos0/env/v9"
	"log"
)

type DBConfig struct {
	Host           string `env:"DB_HOST" envDefault:"localhost"`
	Port           uint16 `env:"DB_PORT" envDefault:"5432"`
	Name           string `env:"DB_NAME" envDefault:"walker_service"`
	User           string `env:"DB_USER" envDefault:"root"`
	Password       string `env:"DB_PASSWORD" envDefault:"password"`
	MaxConnections int    `env:"DB_MAX_CONNECTIONS" envDefault:"10"`
}

type HttpServer struct {
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
	Port uint16 `env:"SERVER_PORT" envDefault:"9000"`
}

type AppConfig struct {
	DB         DBConfig
	HttpServer HttpServer
	/*
	   LevelDebug Level = -4
	   LevelInfo  Level = 0
	   LevelWarn  Level = 4
	   LevelError Level = 8
	*/
	LogLevel int `env:"LOG_LEVEL" envDefault:"0"`
}

func LoadConfig() AppConfig {
	config := AppConfig{}

	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}

	return config
}

func LoadDBConfig() DBConfig {
	config := DBConfig{}

	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}

	return config
}
