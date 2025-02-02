package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	SBOMformat  string
}

func LoadConfig() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL environment variable is not set, using defaults")
		databaseURL = "postgres://sbom:sbom@localhost:5432/sbom"
	}
	sbomFormat := os.Getenv("SBOM_FORMAT")
	if sbomFormat == "" {
		fmt.Println("SBOM_FORMAT environment variable is not set, using defaults")
		sbomFormat = "cyclonedx"
	}

	return &Config{
		DatabaseURL: databaseURL,
		SBOMformat:  sbomFormat,
	}, nil
}
