package models

import (
	"time"

	"gorm.io/gorm"
)

// GrowthRecord 小成长记录
type GrowthRecord struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Date      time.Time      `gorm:"index" json:"date"`
	Content   string         `gorm:"type:text" json:"content"`
	Images    string         `gorm:"type:text" json:"images"` // JSON array of paths
	Audio     string         `gorm:"size:256" json:"audio"`   // Path to audio file
	Tags      string         `gorm:"type:text" json:"tags"`   // JSON array of Tag Names or IDs
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// GrowthTag 小成长标签（可选，如果需要管理标签树）
type GrowthTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Name      string    `gorm:"size:64;index" json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
