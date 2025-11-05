package handlers

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "recordgo/internal/common"
)

// UploadTaskImage 上传任务图片（前端需先压缩并转换为 webp）
// 中文注释：接收表单字段 user_id 与文件 file，保存到 storage/uploads/images/task_images/{用户id}
func UploadTaskImage(c *gin.Context) {
    // 中文注释：按需求必须提供用户ID与任务ID，图片已在前端转换为 webp
    userID := c.PostForm("user_id")
    taskID := c.PostForm("task_id")
    if strings.TrimSpace(userID) == "" {
        common.Error(c, 40001, "缺少 user_id")
        return
    }
    if strings.TrimSpace(taskID) == "" {
        common.Error(c, 40003, "缺少 task_id")
        return
    }
    file, err := c.FormFile("file")
    if err != nil {
        common.Error(c, 40002, "缺少文件")
        return
    }
    // 保存到约定目录：storage/uploads/images/task_images/{用户id}/{taskid}
    path, err := saveTaskImage(file, userID, taskID)
    if err != nil {
        common.Error(c, 50018, "保存任务图片失败")
        return
    }
    // 返回相对路径，前端记录到 image_json（JSON 数组）
    common.Ok(c, gin.H{"path": path})
}

// saveTaskImage 将上传文件保存到 storage/uploads/images/task_images/{用户id}/{taskid}
// 中文注释：文件命名为 task_任务ID_时间戳_uuid.webp（uuid 用纳秒替代）
func saveTaskImage(file *multipart.FileHeader, userID string, taskID string) (string, error) {
    root := os.Getenv("STORAGE_ROOT")
    if strings.TrimSpace(root) == "" {
        root = "storage"
    }
    dir := filepath.Join(root, "uploads", "images", "task_images", userID, taskID)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return "", fmt.Errorf("创建目录失败: %w", err)
    }
    filename := fmt.Sprintf("task_%s_%d_%d.webp", taskID, time.Now().Unix(), time.Now().UnixNano())
    full := filepath.Join(dir, filename)
    if err := saveUploadedFileGeneric(file, full); err != nil {
        return "", err
    }
    rel := filepath.ToSlash(filepath.Join("uploads", "images", "task_images", userID, taskID, filename))
    return rel, nil
}

// saveUploadedFileGeneric 兼容性的保存封装（避免跨平台路径问题）
func saveUploadedFileGeneric(file *multipart.FileHeader, dst string) error {
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()
    if _, err := io.Copy(out, src); err != nil {
        return err
    }
    return nil
}