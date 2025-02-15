package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Host     string `envconfig:"DATABASE_HOST" required:"true"`
	Port     string `envconfig:"DATABASE_PORT" required:"true"`
	User     string `envconfig:"DATABASE_USER" required:"true"`
	Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	Name     string `envconfig:"DATABASE_NAME" required:"true"`
}

func NewDBConfig() *DBConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Warning: No .env file found, using system env variables")
	} else {
		log.Println(".env file loaded successfully")
	}
	var config DBConfig

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &config
}
