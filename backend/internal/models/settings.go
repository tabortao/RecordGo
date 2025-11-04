package models

import "time"

// 中文注释：用户设置模型（简化版）
type UserSettings struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"uniqueIndex"`
    FixedTomatoPage    bool
    TaskAutoSort       bool
    TomatoSettingsJSON string  `gorm:"type:text"`
    TaskMigrationJSON  string  `gorm:"type:text"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
