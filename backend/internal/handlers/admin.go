package handlers

import (
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminUsersOverview 返回用户总数与今日登录数
func AdminUsersOverview(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	var total int64
	if err := db.DB().Model(&models.User{}).Count(&total).Error; err != nil {
		common.Error(c, 50001, "统计失败")
		return
	}
	start := time.Now().Truncate(24 * time.Hour)
	var today int64
	if err := db.DB().Model(&models.User{}).Where("last_login_time >= ?", start).Count(&today).Error; err != nil {
		common.Error(c, 50002, "统计失败")
		return
	}
	common.Ok(c, gin.H{"total_users": total, "today_logins": today})
}

// AdminListUsers 用户列表，支持搜索
func AdminListUsers(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	q := strings.TrimSpace(c.Query("q"))
	var users []models.User
	tx := db.DB().Model(&models.User{})
	if q != "" {
		like := "%" + q + "%"
		tx = tx.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", like, like, like)
	}
	if err := tx.Order("id ASC").Find(&users).Error; err != nil {
		common.Error(c, 50010, "查询失败")
		return
	}
	// 精简返回字段
	list := make([]gin.H, 0, len(users))
	for _, u := range users {
		list = append(list, gin.H{
			"id":       u.ID,
			"username": u.Username,
			"nickname": u.Nickname,
			"email":    u.Email,
			"phone":    u.Phone,
			"is_vip":   u.IsVIP,
			"vip_expire_time": func() any {
				if u.VIPExpireTime != nil {
					return u.VIPExpireTime.Format(time.RFC3339)
				}
				return nil
			}(),
			"is_lifetime_vip": u.IsLifetimeVIP,
			"last_login_time": func() any {
				if u.LastLoginTime != nil {
					return u.LastLoginTime.Format(time.RFC3339)
				}
				return nil
			}(),
			"is_disabled": u.IsDisabled,
		})
	}
	common.Ok(c, gin.H{"items": list})
}

// AdminUpdateVIP 设置/取消 VIP 状态
func AdminUpdateVIP(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	var payload struct {
		IsVIP         *bool   `json:"is_vip"`
		VIPExpireTime *string `json:"vip_expire_time"`
		IsLifetimeVIP *bool   `json:"is_lifetime_vip"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.Error(c, 40000, "参数错误")
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	updates := map[string]any{}
	if payload.IsVIP != nil {
		updates["is_vip"] = *payload.IsVIP
	}
	if payload.IsLifetimeVIP != nil {
		updates["is_lifetime_vip"] = *payload.IsLifetimeVIP
	}
	if payload.VIPExpireTime != nil {
		s := strings.TrimSpace(*payload.VIPExpireTime)
		if s == "" {
			updates["vip_expire_time"] = nil
		} else {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				common.Error(c, 40002, "VIP到期时间格式错误")
				return
			}
			updates["vip_expire_time"] = t
		}
	}
	if len(updates) == 0 {
		common.Error(c, 40001, "无更新字段")
		return
	}
	if err := db.DB().Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		common.Error(c, 50020, "更新失败")
		return
	}
	zap.L().Info("管理员更新VIP", zap.String("user_id", id), zap.Any("updates", updates))
	common.Ok(c, gin.H{"updated": true})
}

// AdminSetDisabled 禁用/启用用户
func AdminSetDisabled(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	var payload struct {
		Disabled bool `json:"disabled"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.Error(c, 40000, "参数错误")
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if err := db.DB().Model(&models.User{}).Where("id = ?", id).Update("is_disabled", payload.Disabled).Error; err != nil {
		common.Error(c, 50030, "更新失败")
		return
	}
	zap.L().Warn("管理员更改禁用状态", zap.String("user_id", id), zap.Bool("disabled", payload.Disabled))
	common.Ok(c, gin.H{"updated": true})
}

// AdminDeleteUser 删除用户
func AdminDeleteUser(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if err := db.DB().Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		common.Error(c, 50040, "删除失败")
		return
	}
	zap.L().Warn("管理员删除用户", zap.String("user_id", id))
	common.Ok(c, gin.H{"deleted": true})
}
