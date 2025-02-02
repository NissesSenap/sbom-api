package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL environment variable is not set, using defaults")
		databaseURL = "postgres://sbom:sbom@localhost:5432/sbom"
	}

	return &Config{
		DatabaseURL: databaseURL,
	}, nil
}
