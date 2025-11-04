package models

import "time"

// 中文注释：用户模型（简化版），后续可按 PRD 完善
type User struct {
    ID          uint      `gorm:"primaryKey"`
    Username    string    `gorm:"uniqueIndex;size:64"`
    PasswordSha string    `gorm:"size:128"`
    Role        string    `gorm:"size:16"` // admin/user
    Permissions string    `gorm:"type:text"` // JSON
    Coins       int64
    Tomatoes    int64
    Nickname    string    `gorm:"size:64"`
    AvatarPath  string    `gorm:"size:256"`
    Phone       string    `gorm:"size:32"`
    ParentID    *uint
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
