package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NissesSenap/sbom-api/internal/config"
	"github.com/NissesSenap/sbom-api/internal/db"
	"github.com/NissesSenap/sbom-api/internal/sbom"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	db.CreateTables(ctx, dbpool)

	bom, err := sbom.ParseSBOM("../../go-bom.json")
	err = sbom.StoreSBOM(ctx, dbpool, bom)
	if err != nil {
		log.Fatalf("Failed to parse and store SBOM: %v\n", err)
	}

	fmt.Println("Database connected, tables created, and SBOM data stored successfully!")
}
