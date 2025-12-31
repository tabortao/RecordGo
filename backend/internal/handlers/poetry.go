package handlers

import (
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
)

// GetPoetryData 获取古诗词数据
func GetPoetryData(c *gin.Context) {
	userID := c.GetUint("userID")
	var data models.PoetryData
	if err := db.DB().Where("user_id = ?", userID).First(&data).Error; err != nil {
		// 不存在则返回空
		common.Ok(c, gin.H{"data": "{}"})
		return
	}
	common.Ok(c, gin.H{"data": data.DataJSON})
}

// SavePoetryData 保存古诗词数据
func SavePoetryData(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}

	var data models.PoetryData
	db.DB().Where("user_id = ?", userID).FirstOrInit(&data)
	data.UserID = userID
	data.DataJSON = req.Data
	if err := db.DB().Save(&data).Error; err != nil {
		common.Error(c, 500, "保存失败")
		return
	}
	common.Ok(c, nil)
}
