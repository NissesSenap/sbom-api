package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NissesSenap/sbom-api/internal/db"
	"github.com/NissesSenap/sbom-api/internal/sbom"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://sbom:sbom@localhost:5432/sbom"
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	db.CreateTables(ctx, dbpool)

	err = sbom.ParseAndStoreSBOM(ctx, dbpool, "go-bom.json")
	if err != nil {
		log.Fatalf("Failed to parse and store SBOM: %v\n", err)
	}

	fmt.Println("Database connected, tables created, and SBOM data stored successfully!")
}
