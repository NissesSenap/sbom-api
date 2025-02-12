package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NissesSenap/sbom-api/internal/config"
	"github.com/NissesSenap/sbom-api/internal/db"
	"github.com/NissesSenap/sbom-api/internal/sbom"
	"github.com/NissesSenap/sbom-api/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful termination
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Received termination signal, shutting down gracefully...")
		cancel()
	}()

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	defer dbpool.Close()

	if err := db.CreateTables(ctx, dbpool); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	var parser sbom.SBOMParser
	switch cfg.SBOMformat {
	case "cyclonedx":
		parser = &sbom.CycloneDXParser{}
	default:
		return fmt.Errorf("unsupported SBOM format: %s", cfg.SBOMformat)
	}

	storageService, err := storage.NewS3Storage(cfg.AWSEndpoint, cfg.AWSAccessKey, cfg.AWSSecretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize storage service: %w", err)
	}

	const bomFilename = "go-bom.json"
	const appname = "sbom-api"
	bom, err := parser.Parse("go-bom.json")
	if err != nil {
		return fmt.Errorf("failed to parse SBOM: %w", err)
	}

	sbomURL := fmt.Sprintf("%v/%v", appname, bomFilename)
	if err := parser.Store(ctx, dbpool, bom, sbomURL); err != nil {
		return fmt.Errorf("failed to store SBOM: %w", err)
	}

	if err := storageService.Upload(ctx, cfg.S3Bucket, sbomURL, bomFilename); err != nil {
		return fmt.Errorf("failed to upload SBOM file to storage: %w", err)
	}

	fmt.Println("Database connected, tables created, and SBOM data stored successfully!")
	return nil
}
