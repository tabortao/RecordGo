package storage

import (
    "fmt"
    "os"
    "path/filepath"
)

// 中文注释：确保项目需要的存储目录存在
func EnsurePaths(root string) error {
    paths := []string{
        filepath.Join(root, "database"),
        filepath.Join(root, "uploads", "images", "avatars"),
        filepath.Join(root, "uploads", "images", "task_images"),
        filepath.Join(root, "uploads", "images", "wish"),
        filepath.Join(root, "images", "honors"),
    }
    for _, p := range paths {
        if err := os.MkdirAll(p, 0o755); err != nil {
            return fmt.Errorf("创建目录失败 %s: %w", p, err)
        }
    }
    return nil
}
