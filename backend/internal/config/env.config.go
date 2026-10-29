package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	HOST            string
	PORT            string
	POSTGRES        *Postgres
	ALLOWED_ORIGINS []string
}

func LoadEnv() *Env {
	godotenv.Load()
	origins := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.Trim(o, " \t\n\r\"'")
		if trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}

	return &Env{
		HOST: os.Getenv("HOST_SERVER"),
		PORT: os.Getenv("PORT_SERVER"),
		POSTGRES: &Postgres{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Db_name:  os.Getenv("POSTGRES_DB"),
		},
		ALLOWED_ORIGINS: allowedOrigins,
	}
}
