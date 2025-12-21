package filestore

import (
	"context"
	"io"
)

// FileStore defines the behavior for file operations.
// The Application Logic (Service) depends on THIS, not on MinIO.
type FileStore interface {
	// Upload streams data to the storage
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error

	// Delete removes a file (useful for rollbacks)
	Delete(ctx context.Context, objectName string) error

	// (optional) PresignedURL generates a temporary public link
	// PresignedURL(ctx context.Context, objectName string) (string, error)
}
