package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL  string
	SBOMformat   string
	S3Bucket     string
	AWSEndpoint  string
	AWSAccessKey string
	AWSSecretKey string
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

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		fmt.Println("S3_BUCKET environment variable is not set, using defaults")
		s3Bucket = "sbom"
	}

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		fmt.Println("AWS_ENDPOINT environment variable is not set, using defaults")
		awsEndpoint = "http://localhost:9000"
	}

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	if awsAccessKey == "" {
		fmt.Println("AWS_ACCESS_KEY environment variable is not set, using defaults")
		awsAccessKey = "minioadmin"
	}

	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	if awsSecretKey == "" {
		fmt.Println("AWS_SECRET_KEY environment variable is not set, using defaults")
		awsSecretKey = "minioadmin"
	}

	return &Config{
		DatabaseURL:  databaseURL,
		SBOMformat:   sbomFormat,
		S3Bucket:     s3Bucket,
		AWSEndpoint:  awsEndpoint,
		AWSAccessKey: awsAccessKey,
		AWSSecretKey: awsSecretKey,
	}, nil
}
