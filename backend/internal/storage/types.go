package storage

import (
    "context"
    "io"
)

type Storage interface {
    DriverName() string
    PresignPut(ctx context.Context, key string, contentType string, expireSeconds int64) (string, map[string]string, error)
    PresignGet(ctx context.Context, key string, expireSeconds int64) (string, error)
    Put(ctx context.Context, key string, r io.Reader, contentType string) error
    Delete(ctx context.Context, key string) error
    PublicURL(key string) (string, bool)
}