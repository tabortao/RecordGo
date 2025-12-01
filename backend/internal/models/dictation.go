package models

import (
	"time"
	"gorm.io/gorm"
)

// DictationWordBank 听写词库
type DictationWordBank struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"index" json:"user_id"` // 所属用户ID
	Title          string         `json:"title"`                // 词库标题
	EducationStage string         `json:"education_stage"`      // 教育阶段（小学、初中、高中）
	Subject        string         `json:"subject"`              // 科目（语文、英语等）
	Version        string         `json:"version"`              // 教材版本
	Grade          string         `json:"grade"`                // 年级
	Content        string         `gorm:"type:text" json:"content"` // 词库内容（JSON string array or simple text）
	IsPublic       bool           `json:"is_public"`            // 是否公开（预留）
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// DictationSettings 听写设置
type DictationSettings struct {
	UserID          uint    `gorm:"primaryKey" json:"user_id"`
	SplitRule       string  `json:"split_rule"`       // 分词规则：space, newline
	PlayMode        string  `json:"play_mode"`        // 模式：read, dictate
	OrderMode       string  `json:"order_mode"`       // 顺序：sequential, random
	RepeatCount     int     `json:"repeat_count"`     // 重复次数
	IntervalSeconds int     `json:"interval_seconds"` // 间隔时间
	VoiceType       string  `json:"voice_type"`       // 音色
	Speed           float64 `json:"speed"`            // 语速
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DictationHistory 听写历史
type DictationHistory struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"index" json:"user_id"`
	WordBankID      *uint          `json:"word_bank_id"` // 关联词库ID（可选）
	ContentSnapshot string         `gorm:"type:text" json:"content_snapshot"` // 听写内容快照
	DurationSeconds int            `json:"duration_seconds"` // 耗时
	MistakeCount    int            `json:"mistake_count"`    // 错题数
	CreatedAt       time.Time      `json:"created_at"`
}

// MistakeNotebook 错题本
type MistakeNotebook struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Word      string         `json:"word"`    // 错词
	Context   string         `json:"context"` // 上下文/来源
	CreatedAt time.Time      `json:"created_at"`
}
