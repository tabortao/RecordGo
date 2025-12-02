package storage

import (
    "context"
    "io"
    "os"
    "path/filepath"
)

type LocalStorage struct {
    Root string
}

func NewLocal(root string) *LocalStorage { return &LocalStorage{Root: root} }

func (l *LocalStorage) DriverName() string { return "local" }

func (l *LocalStorage) PresignPut(ctx context.Context, key string, contentType string, expireSeconds int64) (string, map[string]string, error) {
    url := "/api/storage/put?key=" + key
    headers := map[string]string{"Content-Type": contentType}
    return url, headers, nil
}

func (l *LocalStorage) PresignGet(ctx context.Context, key string, expireSeconds int64) (string, error) {
    return "/api/uploads/" + key, nil
}

func (l *LocalStorage) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
    dir := filepath.Join(l.Root, "uploads")
    full := filepath.Join(dir, filepath.FromSlash(key))
    if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil { return err }
    f, err := os.Create(full)
    if err != nil { return err }
    defer f.Close()
    if _, err := io.Copy(f, r); err != nil { return err }
    return nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
    dir := filepath.Join(l.Root, "uploads")
    full := filepath.Join(dir, filepath.FromSlash(key))
    return os.Remove(full)
}

func (l *LocalStorage) PublicURL(key string) (string, bool) {
    return "/api/uploads/" + key, true
}