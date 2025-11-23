package handlers

import (
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    jwt "github.com/golang-jwt/jwt/v5"
    "recordgo/internal/common"
    "recordgo/internal/config"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

type TokenLoginReq struct {
    Token string `json:"token"`
}

// TokenLogin 使用子账号令牌登录：返回与普通登录一致的结构
func TokenLogin(c *gin.Context) {
    var req TokenLoginReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Token) == "" {
        common.Error(c, 40001, "参数错误：缺少 token")
        return
    }
    var u models.User
    if err := db.DB().Where("login_token = ?", strings.TrimSpace(req.Token)).First(&u).Error; err != nil {
        common.Error(c, 40405, "令牌无效或子账号不存在")
        return
    }
    // 仅允许子账号使用令牌登录（必须有 ParentID）
    if u.ParentID == nil {
        common.Error(c, 40302, "仅子账号可使用令牌登录")
        return
    }
    cfg, err := config.Load()
    if err != nil || cfg.SecretKey == "" {
        common.Error(c, 50020, "服务配置错误")
        return
    }
    claims := &Claims{
        UserID:     u.ID,
        Username:   u.Username,
        Role:       u.Role,
        Permissions: u.Permissions,
        ParentID:   u.ParentID,
        LoginToken: strings.TrimSpace(u.LoginToken),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTExpireDays) * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString([]byte(cfg.SecretKey))
    if err != nil {
        common.Error(c, 50021, "令牌签发失败")
        return
    }
    if u.LoginTokenExpireAt != nil && time.Now().After(*u.LoginTokenExpireAt) {
        common.Error(c, 40303, "令牌已过期")
        return
    }
    // 中文注释：计算返回的金币值——若启用了父子金币同步，则使用父账号金币
    coinsToReturn := u.Coins
    cfg2, _ := config.Load()
    if cfg2 != nil && cfg2.ParentCoinsSync && u.ParentID != nil {
        var parent models.User
        if err := db.DB().First(&parent, *u.ParentID).Error; err == nil {
            coinsToReturn = parent.Coins
        }
    }
    common.Ok(c, gin.H{
        "token": tokenStr,
        "user": gin.H{
            "id": u.ID,
            "username": u.Username,
            "role": u.Role,
            "permissions": u.Permissions,
            "parent_id": u.ParentID,
            "nickname": u.Nickname,
            "avatar_path": u.AvatarPath,
            "coins": coinsToReturn,
        },
    })
}