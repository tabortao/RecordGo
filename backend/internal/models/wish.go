package models

import "time"

// 中文注释：心愿模型（简化版）
type Wish struct {
    ID            uint      `gorm:"primaryKey"`
    UserID        uint      `gorm:"index"`
    Name          string    `gorm:"size:128"`
    Content       string    `gorm:"type:text"`
    Icon          string    `gorm:"size:256"`
    NeedCoins     int       `gorm:"default:1"`
    ExchangeAmount int      `gorm:"default:1"`
    Unit          string    `gorm:"size:16"`
    Exchanged     int       `gorm:"default:0"`
    BuiltIn       bool      `gorm:"default:false"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
