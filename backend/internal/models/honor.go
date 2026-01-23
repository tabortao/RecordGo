package models

import "time"

// 中文注释：荣誉记录模型，记录荣誉照片与关键信息
type Honor struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Title      string    `gorm:"size:128" json:"title"`
	Issuer     string    `gorm:"size:128" json:"issuer"`
	AwardedAt  time.Time `gorm:"index" json:"awarded_at"`
	PhotoURL   string    `gorm:"size:512" json:"photo_url"`
	Remark     string    `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
