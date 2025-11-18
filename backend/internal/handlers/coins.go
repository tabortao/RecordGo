package handlers

import (
    "github.com/gin-gonic/gin"
    "recordgo/internal/common"
    "recordgo/internal/config"
    "recordgo/internal/db"
    "recordgo/internal/models"
    "strings"
)

// GetCoins 返回当前登录用户的“有效金币”
// 中文注释：若启用父子金币同步且为子账号，则返回父账号金币；否则返回当前账号金币
func GetCoins(c *gin.Context) {
    cl := extractClaims(c)
    if cl == nil {
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    var u models.User
    if err := db.DB().First(&u, cl.UserID).Error; err != nil {
        common.Error(c, 40401, "用户不存在")
        return
    }
    coins := u.Coins
    // 中文注释：父子金币池共享（默认开启）。子账号始终返回“父账号金币”，父账号返回自己的金币。
    if cfg, _ := config.Load(); cfg != nil && cfg.ParentCoinsSync {
        if u.ParentID != nil {
            var parent models.User
            if err := db.DB().First(&parent, *u.ParentID).Error; err == nil {
                coins = parent.Coins
            } else {
                // 中文注释：父账号不存在时回退到当前账号金币，避免 0 显示
                coins = u.Coins
            }
        } else {
            coins = u.Coins
        }
    }
    common.Ok(c, gin.H{"coins": coins})
}

type SetCoinsReq struct {
    Coins  int64  `json:"coins"`
    Reason string `json:"reason"`
}

// SetCoins 设置当前登录用户的总金币（考虑父子同步）
func SetCoins(c *gin.Context) {
    cl := extractClaims(c)
    if cl == nil { common.Error(c, 40100, "未登录或令牌无效"); return }
    var req SetCoinsReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Coins < 0 || strings.TrimSpace(req.Reason) == "" {
        common.Error(c, 40001, "参数错误")
        return
    }
    var u models.User
    if err := db.DB().First(&u, cl.UserID).Error; err != nil { common.Error(c, 40401, "用户不存在"); return }
    target := &u
    if cfg, _ := config.Load(); cfg != nil && cfg.ParentCoinsSync && u.ParentID != nil {
        var parent models.User
        if err := db.DB().First(&parent, *u.ParentID).Error; err == nil { target = &parent }
    }
    target.Coins = req.Coins
    if err := db.DB().Save(target).Error; err != nil { common.Error(c, 50020, "保存金币失败"); return }
    common.Ok(c, gin.H{"coins": target.Coins})
}