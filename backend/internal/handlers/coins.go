package handlers

import (
    "github.com/gin-gonic/gin"
    "recordgo/internal/common"
    "recordgo/internal/config"
    "recordgo/internal/db"
    "recordgo/internal/models"
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