package models

import (
	"gorm.io/gorm"
	"time"
)

// 中文注释：任务模型（简化版）
type Task struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"index" json:"user_id"`
	SeriesID         *uint          `json:"series_id,omitempty"`
	Name             string         `gorm:"size:128" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	Icon             string         `gorm:"size:128" json:"icon"`
	Category         string         `gorm:"size:64" json:"category"`
	PlanMinutes      int            `gorm:"default:10" json:"plan_minutes"`
	ActualMinutes    int            `gorm:"default:0" json:"actual_minutes"`
	Score            int            `gorm:"default:1" json:"score"` // 可为负数
	ScoreMode        string         `gorm:"size:16;default:fixed" json:"score_mode"`
	CustomScoreMax   int            `gorm:"default:5" json:"custom_score_max"`
	CompletedScore   int            `gorm:"default:0" json:"completed_score"`
	DailyMaxCheckins int            `gorm:"default:1" json:"daily_max_checkins"`
	Repeat           string         `gorm:"size:32" json:"repeat"`             // 无/每天/每周/每月/自定义星期
	RepeatDaysJSON   string         `gorm:"type:text" json:"repeat_days_json"` // 周期附加数据（如 weekly 的星期集合）
	WeeklyDays       []int          `gorm:"-" json:"weekly_days,omitempty"`    // 运行时解包，数据库不存该列
	StartDate        time.Time      `json:"start_date"`
	EndDate          *time.Time     `json:"end_date,omitempty"`
	Status           int            `gorm:"default:0" json:"status"` // 0待完成 1进行中 2已完成
	Remark           string         `gorm:"type:text" json:"remark"`
	ImageJSON        string         `gorm:"type:text" json:"image_json"`
	TomatoCount      int            `gorm:"default:0" json:"tomato_count"`     // 中文注释：任务关联的番茄钟完成次数
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 中文注释：GORM 软删除支持
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
