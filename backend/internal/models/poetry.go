package models

import "gorm.io/gorm"

// 中文注释：古诗词数据模型，存储用户的古诗学习状态、自定义诗词等所有数据（JSON格式）
type PoetryData struct {
	gorm.Model
	UserID   uint   `gorm:"uniqueIndex"`
	DataJSON string `gorm:"type:text"` // 包含 studyRecords, customPoems, collections, challengeRecords 等
}
