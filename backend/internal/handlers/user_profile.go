package handlers

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "mime/multipart"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    jwt "github.com/golang-jwt/jwt/v5"
    "go.uber.org/zap"
    "recordgo/internal/common"
    "recordgo/internal/config"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

// 中文注释：从 Authorization Bearer Token 中提取用户ID；失败返回 0
func extractUserIDFromToken(c *gin.Context) uint {
    auth := c.Request.Header.Get("Authorization")
    if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
        tokenStr := strings.TrimSpace(auth[7:])
        cfg, err := config.Load()
        if err != nil || cfg.SecretKey == "" {
            return 0
        }
        claims := &Claims{}
        t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(cfg.SecretKey), nil
        })
        if err == nil && t.Valid {
            return claims.UserID
        }
    }
    return 0
}

// UpdateProfile 更新用户资料（昵称/电话/邮箱），字段可选
type UpdateProfileReq struct {
    Nickname *string `json:"nickname"`
    Phone    *string `json:"phone"`
    Email    *string `json:"email"`
}

func UpdateProfile(c *gin.Context) {
    var req UpdateProfileReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    uid := extractUserIDFromToken(c)
    if uid == 0 {
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    var u models.User
    if err := db.DB().First(&u, uid).Error; err != nil {
        common.Error(c, 40401, "用户不存在")
        return
    }
    if req.Nickname != nil {
        u.Nickname = strings.TrimSpace(*req.Nickname)
    }
    if req.Phone != nil {
        u.Phone = strings.TrimSpace(*req.Phone)
    }
    if req.Email != nil {
        e := strings.TrimSpace(*req.Email)
        // 中文注释：基础邮箱格式校验（后端兜底），允许为空
        if e != "" {
            if !strings.Contains(e, "@") || !strings.Contains(e, ".") {
                common.Error(c, 40002, "邮箱格式不正确")
                return
            }
        }
        u.Email = e
    }
    if err := db.DB().Save(&u).Error; err != nil {
        common.Error(c, 50010, "更新昵称失败")
        return
    }
    common.Ok(c, gin.H{"nickname": u.Nickname, "phone": u.Phone, "email": u.Email})
}

// ChangePassword 修改密码
type ChangePasswordReq struct {
    OldPassword string `json:"old_password"`
    NewPassword string `json:"new_password"`
}

func ChangePassword(c *gin.Context) {
    var req ChangePasswordReq
    if err := c.ShouldBindJSON(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
        common.Error(c, 40001, "参数错误：请填写原密码与新密码")
        return
    }
    uid := extractUserIDFromToken(c)
    if uid == 0 {
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    var u models.User
    if err := db.DB().First(&u, uid).Error; err != nil {
        common.Error(c, 40401, "用户不存在")
        return
    }
    // 校验原密码
    oh := sha256.Sum256([]byte(req.OldPassword))
    if u.PasswordSha != hex.EncodeToString(oh[:]) {
        common.Error(c, 40003, "原密码错误")
        return
    }
    nh := sha256.Sum256([]byte(req.NewPassword))
    u.PasswordSha = hex.EncodeToString(nh[:])
    if err := db.DB().Save(&u).Error; err != nil {
        common.Error(c, 50011, "修改密码失败")
        return
    }
    common.Ok(c, gin.H{"message": "ok"})
}

type AvatarObjectReq struct {
    ObjectKey string `json:"object_key"`
}

func UploadAvatarObject(c *gin.Context) {
    var req AvatarObjectReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.ObjectKey) == "" {
        common.Error(c, 40001, "参数错误")
        return
    }
    uid := extractUserIDFromToken(c)
    if uid == 0 {
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    var u models.User
    if err := db.DB().First(&u, uid).Error; err != nil {
        common.Error(c, 40401, "用户不存在")
        return
    }
    u.AvatarPath = strings.TrimSpace(req.ObjectKey)
    if err := db.DB().Save(&u).Error; err != nil {
        common.Error(c, 50018, "保存头像失败")
        return
    }
    common.Ok(c, gin.H{"path": u.AvatarPath})
}

// UploadAvatar 上传用户头像（前端需先转换为 webp），并更新用户 AvatarPath
func UploadAvatar(c *gin.Context) {
    userIDStr := strings.TrimSpace(c.PostForm("user_id"))
    var uid uint
    if userIDStr != "" {
        if v, err := strconv.Atoi(userIDStr); err == nil && v > 0 { uid = uint(v) }
    }
    // 中文注释：权限校验——允许本人上传；家长可为其子账号上传（需有 manage_children 权限）
    cl := extractClaims(c)
    if cl == nil {
        // 若无有效登录，则不允许上传
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    if uid == 0 { uid = cl.UserID }
    // 本人上传：放行
    allowed := (uid == cl.UserID)
    if !allowed {
        // 家长上传子账号头像：需权限，且目标用户的 ParentID 为当前家长
        if cl.ParentID == nil {
            if hasPermission(c, "account", "manage_children") {
                var target models.User
                if err := db.DB().First(&target, uid).Error; err == nil && target.ParentID != nil && *target.ParentID == cl.UserID {
                    allowed = true
                }
            }
        }
    }
    if !allowed {
        deny(c, "无权限上传该用户头像")
        return
    }
    // 解析文件（兼容 image 与 file）
    if ct := c.Request.Header.Get("Content-Type"); !strings.HasPrefix(strings.ToLower(ct), "multipart/form-data") {
        zap.L().Warn("UploadAvatar: invalid content-type", zap.String("content_type", ct))
    }
    _ = c.Request.ParseMultipartForm(10 << 20)
    file, err := c.FormFile("image")
    if err != nil || file == nil {
        file, err = c.FormFile("file")
    }
    if err != nil || file == nil {
        common.Error(c, 40002, "缺少文件")
        return
    }
    // 尺寸与类型检查
    const maxSize = int64(5 * 1024 * 1024)
    if file.Size > maxSize {
        common.Error(c, 40005, "文件过大，需小于5MB")
        return
    }
    contentType, ext, err := sniffImageContentType(file)
    if err != nil || !strings.HasPrefix(contentType, "image/") {
        common.Error(c, 40004, "不支持的文件类型")
        return
    }
    // 保存到 storage/uploads/images/avatars/{用户id}
    path, err := saveAvatarFile(file, fmt.Sprintf("%d", uid), ext)
    if err != nil {
        common.Error(c, 50018, "保存头像失败")
        return
    }
    // 更新用户头像路径
    var u models.User
    if err := db.DB().First(&u, uid).Error; err == nil {
        u.AvatarPath = path
        _ = db.DB().Save(&u).Error
    }
    common.Ok(c, gin.H{"path": path})
}

// saveAvatarFile 保存头像到约定目录并命名为 avatar_用户ID_时间戳_uuid.ext
func saveAvatarFile(file *multipart.FileHeader, userID string, ext string) (string, error) {
    root := os.Getenv("STORAGE_ROOT")
    if strings.TrimSpace(root) == "" { root = "storage" }
    dir := filepath.Join(root, "uploads", "images", userID, "avatars")
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return "", fmt.Errorf("创建目录失败: %w", err)
    }
    filename := fmt.Sprintf("avatar_%s_%d_%d.%s", userID, time.Now().Unix(), time.Now().UnixNano(), ext)
    full := filepath.Join(dir, filename)
    if err := saveUploadedFileGeneric(file, full); err != nil {
        return "", err
    }
    rel := filepath.ToSlash(filepath.Join("uploads", "images", userID, "avatars", filename))
    return rel, nil
}
