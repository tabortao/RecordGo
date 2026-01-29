package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"
	"strconv"
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
	var today int64
	day := time.Now().Format("2006-01-02")
	if err := db.DB().Model(&models.UserDailyLogin{}).Where("day = ?", day).Count(&today).Error; err != nil {
		common.Error(c, 50002, "统计失败")
		return
	}
	common.Ok(c, gin.H{"total_users": total, "today_logins": today})
}

func AdminLoginStats(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	rangeDays := 7
	if v := strings.TrimSpace(c.Query("range")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			rangeDays = n
		}
	}
	if rangeDays != 7 && rangeDays != 30 && rangeDays != 365 {
		rangeDays = 7
	}

	now := time.Now()
	endDay := now.Format("2006-01-02")
	startTime := now.AddDate(0, 0, -(rangeDays - 1))
	startDay := startTime.Format("2006-01-02")

	type row struct {
		Day   string `gorm:"column:day"`
		Count int64  `gorm:"column:count"`
	}
	var rows []row
	if err := db.DB().Model(&models.UserDailyLogin{}).
		Select("day, COUNT(*) as count").
		Where("day >= ? AND day <= ?", startDay, endDay).
		Group("day").
		Order("day ASC").
		Scan(&rows).Error; err != nil {
		common.Error(c, 50003, "统计失败")
		return
	}
	countMap := make(map[string]int64, len(rows))
	for _, r := range rows {
		countMap[r.Day] = r.Count
	}

	items := make([]gin.H, 0, rangeDays)
	for i := 0; i < rangeDays; i++ {
		day := startTime.AddDate(0, 0, i).Format("2006-01-02")
		items = append(items, gin.H{"day": day, "count": countMap[day]})
	}
	common.Ok(c, gin.H{
		"range_days": rangeDays,
		"start_day":  startDay,
		"end_day":    endDay,
		"items":      items,
	})
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

	ids := make([]uint, 0, len(users))
	for _, u := range users {
		if u.ID != 0 {
			ids = append(ids, u.ID)
		}
	}
	todayDay := time.Now().Format("2006-01-02")
	start30Day := time.Now().AddDate(0, 0, -29).Format("2006-01-02")
	todayMap := map[uint]bool{}
	active30Map := map[uint]bool{}
	if len(ids) > 0 {
		var todayIDs []uint
		_ = db.DB().Model(&models.UserDailyLogin{}).Where("day = ? AND user_id IN ?", todayDay, ids).Pluck("user_id", &todayIDs).Error
		for _, id := range todayIDs {
			todayMap[id] = true
		}
		var active30IDs []uint
		_ = db.DB().Model(&models.UserDailyLogin{}).Distinct("user_id").Where("day >= ? AND user_id IN ?", start30Day, ids).Pluck("user_id", &active30IDs).Error
		for _, id := range active30IDs {
			active30Map[id] = true
		}
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
			"login_today":   todayMap[u.ID],
			"inactive_30d":  !active30Map[u.ID],
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
		zap.L().Error("AdminUpdateVIP: bind error", zap.Error(err))
		common.Error(c, 40000, "参数错误: "+err.Error())
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		common.Error(c, 40001, "用户ID不能为空")
		return
	}

	// 中文注释：显式检查并同步 VIP 相关字段，防止 schema 不一致导致保存失败
	m := db.DB().Migrator()
	if m != nil {
		// 检查并添加 IsVIP
		if !m.HasColumn(&models.User{}, "IsVIP") {
			zap.L().Warn("AdminUpdateVIP: adding missing column IsVIP")
			if err := m.AddColumn(&models.User{}, "IsVIP"); err != nil {
				zap.L().Error("AdminUpdateVIP: failed to add IsVIP", zap.Error(err))
			}
		}
		// 检查并添加 VIPExpireTime
		if !m.HasColumn(&models.User{}, "VIPExpireTime") {
			zap.L().Warn("AdminUpdateVIP: adding missing column VIPExpireTime")
			if err := m.AddColumn(&models.User{}, "VIPExpireTime"); err != nil {
				zap.L().Error("AdminUpdateVIP: failed to add VIPExpireTime", zap.Error(err))
			}
		}
		// 检查并添加 IsLifetimeVIP
		if !m.HasColumn(&models.User{}, "IsLifetimeVIP") {
			zap.L().Warn("AdminUpdateVIP: adding missing column IsLifetimeVIP")
			if err := m.AddColumn(&models.User{}, "IsLifetimeVIP"); err != nil {
				zap.L().Error("AdminUpdateVIP: failed to add IsLifetimeVIP", zap.Error(err))
			}
		}
	}

	// 执行更新 (改用 First + Save 模式以规避 Updates 在某些 GORM/SQLite 版本下的兼容性问题)
	var user models.User
	if err := db.DB().First(&user, id).Error; err != nil {
		common.Error(c, 40400, "用户不存在")
		return
	}

	if payload.IsVIP != nil {
		user.IsVIP = *payload.IsVIP
	}
	if payload.IsLifetimeVIP != nil {
		user.IsLifetimeVIP = *payload.IsLifetimeVIP
	}
	if payload.VIPExpireTime != nil {
		s := strings.TrimSpace(*payload.VIPExpireTime)
		if s == "" {
			user.VIPExpireTime = nil
		} else {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				common.Error(c, 40002, "VIP到期时间格式错误")
				return
			}
			user.VIPExpireTime = &t
		}
	}

	if err := db.DB().Save(&user).Error; err != nil {
		zap.L().Error("AdminUpdateVIP save error", zap.String("user_id", id), zap.Error(err))
		common.Error(c, 50020, "更新失败："+err.Error())
		return
	}
	zap.L().Info("管理员更新VIP", zap.String("user_id", id), zap.Bool("is_vip", user.IsVIP))
	common.Ok(c, gin.H{"updated": true})
}

// AdminSetDisabled 禁用/启用用户
func AdminSetDisabled(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	_ = db.DB().AutoMigrate(&models.User{})
	if m := db.DB().Migrator(); m != nil {
		if !m.HasColumn(&models.User{}, "IsDisabled") && !m.HasColumn(&models.User{}, "is_disabled") {
			_ = m.AddColumn(&models.User{}, "IsDisabled")
		}
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
		msg := "更新失败"
		if strings.Contains(strings.ToLower(err.Error()), "no such column") {
			// 自动迁移并补列后重试一次
			_ = db.DB().AutoMigrate(&models.User{})
			if m := db.DB().Migrator(); m != nil {
				if !m.HasColumn(&models.User{}, "IsDisabled") && !m.HasColumn(&models.User{}, "is_disabled") {
					_ = m.AddColumn(&models.User{}, "IsDisabled")
				}
			}
			if e2 := db.DB().Model(&models.User{}).Where("id = ?", id).Update("is_disabled", payload.Disabled).Error; e2 != nil {
				msg = "更新失败：数据库缺少字段，请重启后端或升级版本"
				common.Error(c, 50030, msg)
				return
			}
			zap.L().Warn("管理员更改禁用状态(重试成功)", zap.String("user_id", id), zap.Bool("disabled", payload.Disabled))
			common.Ok(c, gin.H{"updated": true})
			return
		}
		common.Error(c, 50030, msg)
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

func AdminResetPassword(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		common.Error(c, 40001, "用户ID不能为空")
		return
	}
	var u models.User
	if err := db.DB().First(&u, id).Error; err != nil {
		common.Error(c, 40400, "用户不存在")
		return
	}
	tmp, err := generateTempPassword()
	if err != nil {
		common.Error(c, 50020, "生成临时密码失败")
		return
	}
	h := sha256.Sum256([]byte(tmp))
	if err := db.DB().Model(&u).Updates(map[string]any{
		"password_sha":         hex.EncodeToString(h[:]),
		"must_change_password": true,
	}).Error; err != nil {
		common.Error(c, 50021, "重置密码失败")
		return
	}
	common.Ok(c, gin.H{"temp_password": tmp})
}
