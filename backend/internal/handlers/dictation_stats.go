package handlers

import (
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
)

type DictationStats struct {
	TotalDuration   int                       `json:"total_duration"`
	TotalWords      int                       `json:"total_words"`
	RecentHistory   []models.DictationHistory `json:"recent_history"`
	TotalDictations int                       `json:"total_dictations"`
}

func GetDictationStats(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var stats DictationStats

	// Calculate total duration
	var totalDuration int64
	db.DB().Model(&models.DictationHistory{}).Where("user_id = ?", userID).Select("sum(duration_seconds)").Scan(&totalDuration)
	stats.TotalDuration = int(totalDuration)

	// Calculate total dictations
	var totalCount int64
	db.DB().Model(&models.DictationHistory{}).Where("user_id = ?", userID).Count(&totalCount)
	stats.TotalDictations = int(totalCount)

	// We don't track words count explicitly per history item yet (only content snapshot),
	// so we might skip TotalWords or estimate it.
	// Let's just return 0 or implement a way to count.
	stats.TotalWords = 0 // Placeholder

	// Recent history
	if err := db.DB().Where("user_id = ?", userID).Order("created_at desc").Limit(5).Find(&stats.RecentHistory).Error; err != nil {
		common.Error(c, 500, "查询失败")
		return
	}

	common.Ok(c, stats)
}
