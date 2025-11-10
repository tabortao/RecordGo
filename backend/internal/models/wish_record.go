package models

import "time"

// 中文注释：心愿兑换记录模型，记录用户在何时使用多少金币兑换了什么心愿
type WishRecord struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    UserID     uint      `gorm:"index" json:"user_id"`
    WishID     uint      `gorm:"index" json:"wish_id"`
    WishName   string    `gorm:"size:128" json:"wish_name"`
    CoinsUsed  int       `json:"coins_used"`
    Amount     int       `json:"amount"`      // 中文注释：一次兑换对应的数量（如：10）
    Unit       string    `gorm:"size:16" json:"unit"` // 中文注释：单位（个/次/分钟等）
    Status     string    `gorm:"size:16" json:"status"` // 中文注释：状态（成功/失败），前端展示用
    Remark     string    `gorm:"size:256" json:"remark"` // 中文注释：兑换时填写的备注信息，便于记录具体情况
    CreatedAt  time.Time `json:"created_at"`
}

