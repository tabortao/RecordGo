package models

import "time"

type UserDailyLogin struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null;index:idx_user_day,unique"`
	Day        string    `gorm:"size:10;not null;index:idx_user_day,unique"`
	LastSeenAt time.Time `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

