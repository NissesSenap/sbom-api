package sbom

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SBOMParser interface {
	Parse(filePath string) (interface{}, error)
	Store(ctx context.Context, dbpool *pgxpool.Pool, bom interface{}, sbomURL string) error
}
