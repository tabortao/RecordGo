package handlers

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"recordgo/internal/common"
	"recordgo/internal/config"
	"recordgo/internal/db"
	"recordgo/internal/models"
)

// 中文注释：从 Authorization 中解析 JWT Claims；失败返回 nil
func extractClaims(c *gin.Context) *Claims {
	auth := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		tokenStr := strings.TrimSpace(auth[7:])
		cfg, err := config.Load()
		if err != nil || cfg.SecretKey == "" {
			return nil
		}
		claims := &Claims{}
		t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.SecretKey), nil
		})
		if err == nil && t.Valid {
			// 中文注释：子账号令牌刷新校验——若 JWT 中含有登录令牌字段，则需与数据库中当前 LoginToken 一致；不一致则判定令牌失效
			if claims.ParentID != nil && strings.TrimSpace(claims.LoginToken) != "" {
				var u models.User
				if dberr := db.DB().First(&u, claims.UserID).Error; dberr != nil {
					return nil
				}
				if strings.TrimSpace(u.LoginToken) != strings.TrimSpace(claims.LoginToken) {
					return nil
				}
			}
			return claims
		}
	}
	return nil
}

// 中文注释：检查指定分类/动作的权限；父账号（无 ParentID）默认放行；view_only=true 全部拒绝
func hasPermission(c *gin.Context, category string, action string) bool {
	cl := extractClaims(c)
	if cl == nil {
		return false
	}
	// 父账号：无 ParentID → 默认有全部权限
	if cl.ParentID == nil {
		return true
	}
	// 子账号：视图只读直接拒绝
	var perm map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(cl.Permissions)), &perm); err != nil {
		// JSON 异常时按只读处理
		return false
	}
	if v, ok := perm["view_only"].(bool); ok && v {
		return false
	}
	catRaw, ok := perm[category]
	if !ok {
		return false
	}
	// 结构如 { create: true, edit: true, ... }
	if m, ok := catRaw.(map[string]any); ok {
		if v, ok := m[action].(bool); ok {
			return v
		}
	}
	return false
}

// 中文注释：统一的权限拒绝响应
func deny(c *gin.Context, msg string) {
	if strings.TrimSpace(msg) == "" {
		msg = "无权限操作"
	}
	common.Error(c, 40301, msg)
}

func canAccessUser(c *gin.Context, uid uint) bool {
	cl := extractClaims(c)
	if cl == nil {
		return false
	}
	if uid == cl.UserID {
		return true
	}
	if cl.ParentID == nil {
		var u models.User
		if err := db.DB().First(&u, uid).Error; err == nil && u.ParentID != nil && *u.ParentID == cl.UserID {
			return true
		}
	}
	if cl.ParentID != nil {
		if uid == *cl.ParentID {
			return true
		}
	}
	return false
}

// 中文注释：获取“读取范围”可访问的用户ID列表
// - 父账号：包含自己与所有子账号
// - 子账号：包含自己与主账号（不包含兄弟子账号）
func readableUserIDs(cl *Claims) ([]uint, error) {
	if cl == nil {
		return nil, nil
	}

	var me models.User
	if err := db.DB().Select("id", "parent_id").First(&me, cl.UserID).Error; err == nil {
		if me.ParentID != nil && *me.ParentID != me.ID {
			ids := []uint{me.ID}
			if *me.ParentID != 0 && *me.ParentID != me.ID {
				ids = append(ids, *me.ParentID)
			}
			return ids, nil
		}

		ids := []uint{me.ID}
		var childIDs []uint
		if err := db.DB().Model(&models.User{}).Where("parent_id = ?", me.ID).Pluck("id", &childIDs).Error; err != nil {
			return ids, err
		}
		seen := map[uint]bool{me.ID: true}
		for _, id := range childIDs {
			if id == 0 || seen[id] {
				continue
			}
			ids = append(ids, id)
			seen[id] = true
		}
		return ids, nil
	}

	if cl.ParentID != nil {
		if *cl.ParentID == cl.UserID {
			return []uint{cl.UserID}, nil
		}
		return []uint{cl.UserID, *cl.ParentID}, nil
	}
	return []uint{cl.UserID}, nil
}

// requireAdmin 仅允许用户ID为1的管理员访问（简化版本）
func requireAdmin(c *gin.Context) bool {
	cl := extractClaims(c)
	if cl == nil {
		deny(c, "未登录或令牌无效")
		return false
	}
	if cl.UserID != 1 {
		deny(c, "仅管理员可访问")
		return false
	}
	return true
}
