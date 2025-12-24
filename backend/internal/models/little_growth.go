package models

import (
	"time"

	"gorm.io/gorm"
)

// GrowthRecord 小成长记录
type GrowthRecord struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	UserID     uint            `gorm:"index" json:"user_id"`
	User       User            `gorm:"foreignKey:UserID" json:"user"` // New: Author info
	Date       time.Time       `gorm:"index" json:"date"`
	Content    string          `gorm:"type:text" json:"content"`
	Images     string          `gorm:"type:text" json:"images"` // JSON array of paths
	Audio      string          `gorm:"size:256" json:"audio"`   // Path to audio file
	Tags       string          `gorm:"type:text" json:"tags"`   // JSON array of Tag Names or IDs
	IsPinned   bool            `gorm:"default:false" json:"is_pinned"`
	IsFavorite bool            `gorm:"default:false" json:"is_favorite"`    // New: Favorite status
	Visibility int             `gorm:"default:0" json:"visibility"`         // 0: Family, 1: Private
	Comments   []GrowthComment `gorm:"foreignKey:RecordID" json:"comments"` // New: Comments relation
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// GrowthComment 小成长评论
type GrowthComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RecordID  uint      `gorm:"index" json:"record_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` // Preload user info if needed
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// GrowthTag 小成长标签（可选，如果需要管理标签树）
type GrowthTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Name      string    `gorm:"size:64;index" json:"name"`
	Color     string    `gorm:"size:16" json:"color"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
