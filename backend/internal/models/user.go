package models

import "time"

// 中文注释：用户模型（简化版），后续可按 PRD 完善
type User struct {
    ID          uint      `gorm:"primaryKey"`
    Username    string    `gorm:"uniqueIndex;size:64"`
    PasswordSha string    `gorm:"size:128"`
    Role        string    `gorm:"size:16"` // admin/user
    Permissions string    `gorm:"type:text"` // JSON
    // 中文注释：子账号登录令牌（随机生成，用于免密码登录子账号）
    LoginToken  string    `gorm:"size:128;index"`
    Coins       int64
    Tomatoes    int64
    Nickname    string    `gorm:"size:64"`
    AvatarPath  string    `gorm:"size:256"`
    Phone       string    `gorm:"size:32"`
    Email       string    `gorm:"size:128"`
    ParentID    *uint
    LoginTokenExpireAt *time.Time
    // 新增：最后登录时间
    LastLoginTime *time.Time
    // 新增：VIP 状态与过期时间、终身 VIP 标记
    IsVIP          bool
    VIPExpireTime  *time.Time
    IsLifetimeVIP  bool
    IsDisabled     bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
