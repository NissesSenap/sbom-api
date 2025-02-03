package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTables creates the necessary tables in the database using the schema.sql file.
func CreateTables(ctx context.Context, pool *pgxpool.Pool) error {
	schema, err := os.ReadFile("internal/db/schema/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %e", err)
	}

	_, err = pool.Exec(ctx, string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema.sql: %e", err)
	}
	return nil
}
