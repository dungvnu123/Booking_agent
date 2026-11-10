package storage

import (
	"context"
	"io"
)

// Storage is the attachment boundary. Implement disk/object storage after the
// patient/examination ownership rules and allowed file types have been decided.
type Storage interface {
	Put(ctx context.Context, key string, body io.Reader) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
