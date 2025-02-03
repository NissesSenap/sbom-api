package storage

import "context"

type StorageService interface {
	Upload(ctx context.Context, bucket, key, filePath string) error
}
