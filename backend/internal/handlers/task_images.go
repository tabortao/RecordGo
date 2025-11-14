package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "mime/multipart"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "recordgo/internal/common"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

// UploadTaskImage 上传任务图片（前端需先压缩并转换为 webp）
// 中文注释：接收表单字段 user_id 与文件 file，保存到 storage/uploads/images/task_images/{用户id}
func UploadTaskImage(c *gin.Context) {
    // 中文注释：按需求必须提供用户ID与任务ID，图片已在前端转换为 webp（失败回退原图）
    userIDStr := strings.TrimSpace(c.PostForm("user_id"))
    taskIDStr := strings.TrimSpace(c.PostForm("task_id"))
    // 中文注释：记录请求基础信息，便于排查
    zap.L().Info("UploadTaskImage: received form",
        zap.String("path", c.FullPath()),
        zap.String("method", c.Request.Method),
        zap.String("user_id", userIDStr),
        zap.String("task_id", taskIDStr),
    )
    // 中文注释：记录请求头的 Content-Type 与内容长度，辅助定位 multipart 解析问题
    zap.L().Info("UploadTaskImage: request meta",
        zap.String("content_type", c.Request.Header.Get("Content-Type")),
        zap.Int64("content_length", c.Request.ContentLength),
    )
    if userIDStr == "" {
        zap.L().Warn("UploadTaskImage: missing user_id")
        common.Error(c, 40001, "缺少 user_id")
        return
    }
    if taskIDStr == "" {
        zap.L().Warn("UploadTaskImage: missing task_id")
        common.Error(c, 40003, "缺少 task_id")
        return
    }
    // 中文注释：将 user_id 与 task_id 转为数字，防止目录遍历与非法字符
    uid, parseErr := strconv.Atoi(userIDStr)
    if parseErr != nil || uid <= 0 {
        zap.L().Warn("UploadTaskImage: invalid user_id", zap.String("user_id", userIDStr), zap.Error(parseErr))
        common.Error(c, 40001, "非法 user_id")
        return
    }
    tid, parseErr2 := strconv.Atoi(taskIDStr)
    if parseErr2 != nil || tid <= 0 {
        zap.L().Warn("UploadTaskImage: invalid task_id", zap.String("task_id", taskIDStr), zap.Error(parseErr2))
        common.Error(c, 40003, "非法 task_id")
        return
    }

    // 中文注释：显式解析 multipart 表单，避免默认行为在部分代理/环境下解析失败
    if ct := c.Request.Header.Get("Content-Type"); !strings.HasPrefix(strings.ToLower(ct), "multipart/form-data") {
        zap.L().Warn("UploadTaskImage: invalid content-type", zap.String("content_type", ct))
    }
    if perr := c.Request.ParseMultipartForm(10 << 20); perr != nil {
        zap.L().Warn("UploadTaskImage: ParseMultipartForm error", zap.Error(perr))
    }
    if !canAccessUser(c, uint(uid)) {
        deny(c, "无权限上传该用户图片")
        return
    }
    var t models.Task
    if dbErr := db.DB().First(&t, tid).Error; dbErr != nil {
        common.Error(c, 40401, "任务不存在")
        return
    }
    if t.UserID != uint(uid) {
        deny(c, "任务与用户不匹配")
        return
    }

    // 中文注释：优先兼容字段名 image 与 file（前端统一使用 image）；同时兼容常见数组键 image[]
    file, formErr := c.FormFile("image")
    if formErr != nil || file == nil {
        file, formErr = c.FormFile("file")
    }
    if (formErr != nil || file == nil) && c.Request.MultipartForm != nil {
        if v := c.Request.MultipartForm.File["image[]"]; len(v) > 0 {
            file = v[0]
            formErr = nil
            zap.L().Info("UploadTaskImage: fallback via image[]", zap.String("filename", file.Filename))
        }
    }
    // 中文注释：尝试解析 multipart，记录可用的文件键与数量，必要时进行回退选取
    var parsedKeys []string
    if mf, perr := c.MultipartForm(); perr == nil && mf != nil {
        for k, v := range mf.File {
            parsedKeys = append(parsedKeys, fmt.Sprintf("%s(%d)", k, len(v)))
        }
        zap.L().Info("UploadTaskImage: multipart parsed", zap.Strings("file_keys", parsedKeys))
        // 中文注释：如果常规获取失败但表单中存在文件，回退使用第一个可用文件
        if (formErr != nil || file == nil) && len(mf.File) > 0 {
            for k, v := range mf.File {
                if len(v) > 0 {
                    file = v[0]
                    formErr = nil
                    zap.L().Info("UploadTaskImage: fallback file via multipart", zap.String("key", k), zap.String("filename", file.Filename))
                    break
                }
            }
        }
    } else {
        if perr != nil {
            zap.L().Warn("UploadTaskImage: parse multipart failed", zap.Error(perr))
        }
    }
    if formErr != nil || file == nil {
        zap.L().Warn("UploadTaskImage: missing file", zap.Error(formErr), zap.Strings("available_keys", parsedKeys))
        common.Error(c, 40002, "缺少文件")
        return
    }
    zap.L().Info("UploadTaskImage: file meta",
        zap.String("filename", file.Filename),
        zap.Int64("size", file.Size),
    )

    // 中文注释：安全校验——限制文件大小 < 5MB，校验 MIME 类型为图片
    const maxSize = int64(5 * 1024 * 1024)
    if file.Size > maxSize {
        zap.L().Warn("UploadTaskImage: file too large", zap.Int64("size", file.Size))
        common.Error(c, 40005, "文件过大，需小于5MB")
        return
    }
    contentType, ext, sniffErr := sniffImageContentType(file)
    if sniffErr != nil {
        zap.L().Warn("UploadTaskImage: sniff type failed", zap.Error(sniffErr))
        common.Error(c, 40004, "无法识别文件类型")
        return
    }
    if !strings.HasPrefix(contentType, "image/") {
        zap.L().Warn("UploadTaskImage: unsupported type", zap.String("content_type", contentType))
        common.Error(c, 40004, "不支持的文件类型")
        return
    }
    zap.L().Info("UploadTaskImage: type detected", zap.String("content_type", contentType), zap.String("ext", ext))

    // 保存到约定目录：storage/uploads/images/task_images/{用户id}
    // 中文注释：根据最新规范，图片按用户维度进行归档，不再按任务ID分目录；任务ID保留在文件名中。
    path, saveErr := saveTaskImage(file, fmt.Sprintf("%d", uid), fmt.Sprintf("%d", tid), ext)
    if saveErr != nil {
        zap.L().Error("UploadTaskImage: save failed",
            zap.Error(saveErr),
            zap.Int("user_id", uid),
            zap.Int("task_id", tid),
        )
        common.Error(c, 50018, "保存任务图片失败")
        return
    }
    // 返回相对路径，前端记录到 image_json（JSON 数组）
    zap.L().Info("UploadTaskImage: saved",
        zap.Int("user_id", uid),
        zap.Int("task_id", tid),
        zap.String("path", path),
    )
    common.Ok(c, gin.H{"path": path})
}

// saveTaskImage 将上传文件保存到 storage/uploads/images/task_images/{用户id}
// 中文注释：文件命名为 task_任务ID_时间戳_uuid.webp（uuid 用纳秒替代），目录不再包含任务ID。
func saveTaskImage(file *multipart.FileHeader, userID string, taskID string, ext string) (string, error) {
    root := os.Getenv("STORAGE_ROOT")
    if strings.TrimSpace(root) == "" {
        root = "storage"
    }
    // 中文注释：按最新规范仅按用户ID分目录
    dir := filepath.Join(root, "uploads", "images", "task_images", userID)
    if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil {
        return "", fmt.Errorf("创建目录失败: %w", mkdirErr)
    }
    // 中文注释：文件命名为 task_任务ID_时间戳_uuid.后缀（uuid 用纳秒替代）；后缀来源于 MIME 检测
    filename := fmt.Sprintf("task_%s_%d_%d.%s", taskID, time.Now().Unix(), time.Now().UnixNano(), ext)
    full := filepath.Join(dir, filename)
    if writeErr := saveUploadedFileGeneric(file, full); writeErr != nil {
        return "", writeErr
    }
    // 中文注释：记录保存路径信息
    zap.L().Debug("saveTaskImage: file saved",
        zap.String("dir", dir),
        zap.String("filename", filename),
        zap.String("full", full),
    )
    // 中文注释：返回相对路径 uploads/images/task_images/{用户id}/{文件名}
    rel := filepath.ToSlash(filepath.Join("uploads", "images", "task_images", userID, filename))
    return rel, nil
}

// saveUploadedFileGeneric 兼容性的保存封装（避免跨平台路径问题）
func saveUploadedFileGeneric(file *multipart.FileHeader, dst string) error {
    src, openErr := file.Open()
    if openErr != nil {
        return openErr
    }
    defer src.Close()
    out, createErr := os.Create(dst)
    if createErr != nil {
        return createErr
    }
    defer out.Close()
    if _, copyErr := io.Copy(out, src); copyErr != nil {
        return copyErr
    }
    return nil
}

// sniffImageContentType 读取少量内容检测 MIME，并返回建议扩展名（webp/jpg/png/gif）
// 中文注释：防止用户伪造扩展名，后端以内容为准；失败时返回错误
func sniffImageContentType(file *multipart.FileHeader) (string, string, error) {
    f, openErr2 := file.Open()
    if openErr2 != nil {
        return "", "", openErr2
    }
    defer f.Close()
    buf := make([]byte, 1024)
    n, _ := f.Read(buf)
    if n == 0 {
        return "", "", fmt.Errorf("空文件")
    }
    ct := http.DetectContentType(buf[:n])
    // 映射扩展名
    var ext string
    switch ct {
    case "image/webp":
        ext = "webp"
    case "image/jpeg":
        ext = "jpg"
    case "image/png":
        ext = "png"
    case "image/gif":
        ext = "gif"
    default:
        if strings.HasPrefix(ct, "image/") {
            // 兜底：未知图片类型统一使用原后缀或 webp
            ext = "webp"
        } else {
            return "", "", fmt.Errorf("不支持的类型: %s", ct)
        }
    }
    return ct, ext, nil
}

// DeleteTaskImage 删除单个任务图片文件并更新任务的 image_json
// 中文注释：请求路径包含任务ID（/tasks/:id/image），请求体或查询参数提供相对路径 path
func DeleteTaskImage(c *gin.Context) {
    idStr := strings.TrimSpace(c.Param("id"))
    path := strings.TrimSpace(c.Query("path"))
    if path == "" {
        var body struct{ Path string `json:"path"` }
        if bindErr := c.ShouldBindJSON(&body); bindErr == nil {
            path = strings.TrimSpace(body.Path)
        }
    }
    if idStr == "" || path == "" {
        zap.L().Warn("DeleteTaskImage: missing params", zap.String("id", idStr), zap.String("path", path))
        common.Error(c, 40001, "缺少必要参数")
        return
    }
    // 中文注释：查询任务并解析现有图片列表
    var t models.Task
    if dbErr := db.DB().First(&t, idStr).Error; dbErr != nil {
        zap.L().Warn("DeleteTaskImage: task not found", zap.Error(dbErr), zap.String("id", idStr))
        common.Error(c, 40401, "任务不存在")
        return
    }
    if !canAccessUser(c, t.UserID) {
        deny(c, "无权限删除该图片")
        return
    }
    // 中文注释：标准化路径，确保斜杠一致
    norm := filepath.ToSlash(path)
    var list []string
    if strings.TrimSpace(t.ImageJSON) != "" {
        _ = json.Unmarshal([]byte(t.ImageJSON), &list)
    }
    // 中文注释：过滤掉待删除的路径（严格匹配相对路径）
    kept := make([]string, 0, len(list))
    removed := false
    for _, p := range list {
        if filepath.ToSlash(p) == norm {
            removed = true
            continue
        }
        kept = append(kept, p)
    }
    if !removed {
        // 中文注释：若任务中不存在该路径，直接返回当前列表
        zap.L().Info("DeleteTaskImage: path not in task images", zap.String("path", norm), zap.Uint("task_id", t.ID))
        common.Ok(c, gin.H{"images": kept})
        return
    }
    // 中文注释：尝试删除物理文件
    root := os.Getenv("STORAGE_ROOT")
    if strings.TrimSpace(root) == "" { root = "storage" }
    full := filepath.Join(root, filepath.FromSlash(norm))
    if rmErr := os.Remove(full); rmErr != nil {
        zap.L().Warn("DeleteTaskImage: remove file failed", zap.String("full", full), zap.Error(rmErr))
    } else {
        zap.L().Info("DeleteTaskImage: removed file", zap.String("full", full))
    }
    // 中文注释：更新任务的 image_json
    b, _ := json.Marshal(kept)
    t.ImageJSON = string(b)
    if saveErr3 := db.DB().Save(&t).Error; saveErr3 != nil {
        zap.L().Error("DeleteTaskImage: db save failed", zap.Error(saveErr3), zap.Uint("task_id", t.ID))
        common.Error(c, 50019, "更新任务图片列表失败")
        return
    }
    common.Ok(c, gin.H{"images": kept})
}