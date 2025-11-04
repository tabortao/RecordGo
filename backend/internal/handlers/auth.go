package handlers

import (
    "crypto/sha256"
    "encoding/hex"
    "time"
    
    "github.com/gin-gonic/gin"
    jwt "github.com/golang-jwt/jwt/v5"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/config"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

// 中文注释：登录与注册请求结构体（字段与前端一致）
type RegisterReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Nickname string `json:"nickname"`
}

type LoginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// 中文注释：JWT Claims 结构，包含必要的用户信息
type Claims struct {
    UserID     uint   `json:"user_id"`
    Username   string `json:"username"`
    Role       string `json:"role"`
    Permissions string `json:"permissions"`
    ParentID   *uint  `json:"parent_id"`
    jwt.RegisteredClaims
}

// 中文注释：Register 用户注册；首次注册的用户为 admin，后续为 user
func Register(c *gin.Context) {
    var req RegisterReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
        common.Error(c, 40001, "参数错误：请填写用户名与密码")
        return
    }

    // 用户名唯一性校验
    var exists models.User
    if err := db.DB().Where("username = ?", req.Username).First(&exists).Error; err == nil {
        common.Error(c, 40002, "用户名已存在")
        return
    }

    // 判断是否首个管理员
    var cnt int64
    _ = db.DB().Model(&models.User{}).Where("role = ?", "admin").Count(&cnt).Error
    role := "user"
    if cnt == 0 {
        role = "admin"
    }

    // 密码 SHA256 加密存储（与 PRD 一致）
    h := sha256.Sum256([]byte(req.Password))
    u := models.User{
        Username:    req.Username,
        PasswordSha: hex.EncodeToString(h[:]),
        Role:        role,
        Permissions: `{"view_only": false}`,
        Nickname:    ifEmpty(req.Nickname, req.Username),
        AvatarPath:  "storage/images/avatars/default.png",
        Coins:       0,
        Tomatoes:    0,
    }
    if err := db.DB().Create(&u).Error; err != nil {
        common.Error(c, 50001, "注册失败")
        return
    }
    common.Ok(c, gin.H{"id": u.ID, "username": u.Username})
}

// 中文注释：Login 用户登录，成功返回 JWT 与用户信息
func Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
        common.Error(c, 40001, "参数错误：请填写用户名与密码")
        return
    }

    var u models.User
    if err := db.DB().Where("username = ?", req.Username).First(&u).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            common.Error(c, 40401, "用户不存在")
            return
        }
        common.Error(c, 50002, "查询用户失败")
        return
    }
    // 校验密码
    h := sha256.Sum256([]byte(req.Password))
    if u.PasswordSha != hex.EncodeToString(h[:]) {
        common.Error(c, 40003, "密码错误")
        return
    }

    // 生成 JWT
    cfg, _ := config.Load()
    claims := Claims{
        UserID:     u.ID,
        Username:   u.Username,
        Role:       u.Role,
        Permissions: ifEmpty(u.Permissions, `{"view_only": false}`),
        ParentID:   u.ParentID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(cfg.SecretKey))
    if err != nil {
        common.Error(c, 50003, "签发令牌失败")
        return
    }
    // 返回与前端一致的用户字段
    common.Ok(c, gin.H{
        "token": signed,
        "user": gin.H{
            "id": u.ID,
            "username": u.Username,
            "nickname": u.Nickname,
            "role": u.Role,
            "permissions": u.Permissions,
            "parent_id": u.ParentID,
            "coins": u.Coins,
            "tomatoes": u.Tomatoes,
            "avatar_path": u.AvatarPath,
        },
    })
}

// 中文注释：空字符串兼容处理
func ifEmpty(s string, def string) string {
    if s == "" { return def }
    return s
}

