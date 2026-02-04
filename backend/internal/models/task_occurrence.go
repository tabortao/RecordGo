package models

import "time"

type TaskOccurrence struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TaskID        uint      `gorm:"index" json:"task_id"`
	Date          time.Time `gorm:"index" json:"date"`
	Status        int       `gorm:"default:0" json:"status"` // 0待完成 2已完成
	CheckinsCount int       `gorm:"default:0" json:"checkins_count"`
	Minutes       int       `gorm:"default:0" json:"minutes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
