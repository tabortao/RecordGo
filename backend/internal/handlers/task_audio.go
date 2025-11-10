package handlers

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "recordgo/internal/common"
)

// 中文注释：上传任务音频（wav/mp3），保存至与任务图片相同的目录：storage/uploads/images/task_images/{用户id}/{任务id}
// 请求：multipart/form-data，字段 user_id, task_id, audio 或 file
// 响应：{ code: 0, message: "ok", data: { path: "uploads/images/task_images/{user}/{task}/xxx.wav" } }
func UploadTaskAudio(c *gin.Context) {
    zap.L().Info("UploadTaskAudio: incoming",
        zap.String("content_type", c.Request.Header.Get("Content-Type")),
        zap.Int64("content_length", c.Request.ContentLength),
    )

    // 解析 multipart
    if ct := c.Request.Header.Get("Content-Type"); !strings.HasPrefix(strings.ToLower(ct), "multipart/form-data") {
        zap.L().Warn("UploadTaskAudio: invalid content-type", zap.String("content_type", ct))
        _ = c.Request.ParseMultipartForm(10 << 20)
    } else {
        _ = c.Request.ParseMultipartForm(50 << 20)
    }

    userID := strings.TrimSpace(c.PostForm("user_id"))
    taskID := strings.TrimSpace(c.PostForm("task_id"))
    if userID == "" || taskID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 user_id 或 task_id", "data": nil})
        return
    }

    // 取文件（优先 audio，其次 file）
    file, err := c.FormFile("audio")
    if err != nil || file == nil {
        file, err = c.FormFile("file")
    }
    if err != nil || file == nil {
        zap.L().Warn("UploadTaskAudio: missing file", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少音频文件", "data": nil})
        return
    }

    // 简单类型校验：允许 audio/wav, audio/mpeg
    ct := file.Header.Get("Content-Type")
    lower := strings.ToLower(ct)
    var ext string
    if strings.Contains(lower, "audio/wav") || strings.Contains(lower, "audio/x-wav") || strings.HasSuffix(strings.ToLower(file.Filename), ".wav") {
        ext = ".wav"
    } else if strings.Contains(lower, "audio/mpeg") || strings.HasSuffix(strings.ToLower(file.Filename), ".mp3") {
        ext = ".mp3"
    } else {
        // 回退默认 wav
        ext = filepath.Ext(file.Filename)
        if ext == "" { ext = ".wav" }
    }

    // 构造保存路径（与图片相同目录层级）
    root := os.Getenv("STORAGE_ROOT")
    if strings.TrimSpace(root) == "" {
        root = "storage"
    }
    filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)) + ext
    full := filepath.Join(root, "uploads", "images", "task_images", userID, taskID, filename)
    // 确保目录存在
    if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
        zap.L().Error("UploadTaskAudio: mkdir failed", zap.Error(err))
        common.Error(c, 50017, fmt.Sprintf("创建目录失败: %v", err))
        return
    }
    if err := saveUploadedFileGeneric(file, full); err != nil {
        zap.L().Error("UploadTaskAudio: save failed", zap.Error(err))
        common.Error(c, 50018, fmt.Sprintf("保存失败: %v", err))
        return
    }
    rel := filepath.ToSlash(filepath.Join("uploads", "images", "task_images", userID, taskID, filename))
    common.Ok(c, gin.H{"path": rel})
}