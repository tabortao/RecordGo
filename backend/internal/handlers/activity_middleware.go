package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"recordgo/internal/db"
	"recordgo/internal/models"
)

func AuthContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := extractClaims(c)
		if cl != nil && cl.UserID != 0 {
			if err := recordDailyLogin(cl.UserID, time.Now()); err != nil {
				zap.L().Warn("record daily login failed", zap.Uint("user_id", cl.UserID), zap.Error(err))
			}
		}
		c.Next()
	}
}

func recordDailyLogin(userID uint, now time.Time) error {
	day := now.Format("2006-01-02")
	item := models.UserDailyLogin{UserID: userID, Day: day, LastSeenAt: now}
	if err := db.DB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "day"}},
		DoUpdates: clause.Assignments(map[string]any{"last_seen_at": now, "updated_at": now}),
	}).Create(&item).Error; err != nil {
		return err
	}

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return db.DB().Model(&models.User{}).
		Where("id = ? AND (last_login_time IS NULL OR last_login_time < ?)", userID, start).
		Update("last_login_time", now).Error
}

