package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTables creates the necessary tables in the database using the schema.sql file.
func CreateTables(ctx context.Context, pool *pgxpool.Pool) {
	schema, err := os.ReadFile("internal/db/schema/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema.sql: %v\n", err)
	}

	_, err = pool.Exec(ctx, string(schema))
	if err != nil {
		log.Fatalf("Failed to execute schema.sql: %v\n", err)
	}
}
