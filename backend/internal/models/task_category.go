package models

import "time"

type TaskCategory struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `gorm:"size:64;uniqueIndex" json:"name"`
    Color string `gorm:"size:16" json:"color"`
    Order int    `gorm:"default:1" json:"order"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}