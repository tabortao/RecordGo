package models

import "time"

// CourseDict 课程字典表
type CourseDict struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"` // 为空表示系统预设，不为空表示用户自定义
	Name      string    `gorm:"size:64;not null" json:"name"`
	Color     string    `gorm:"size:32" json:"color"` // 颜色代码，如 #FFFFFF
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Timetable 课表数据表
type Timetable struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"index;not null" json:"user_id"`
	Grade     string      `gorm:"size:32;index" json:"grade"`    // 年级
	Semester  string      `gorm:"size:32;index" json:"semester"` // 学期：上学期/下学期
	DayOfWeek int         `gorm:"not null" json:"day_of_week"`   // 1-7 (周一到周日)
	Period    int         `gorm:"not null" json:"period"`        // 第几节课
	CourseID  uint        `gorm:"not null" json:"course_id"`
	Course    *CourseDict `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// TimetableConfig 课表配置表
type TimetableConfig struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	ShowSaturday       bool      `gorm:"default:false" json:"show_saturday"`
	ShowSunday         bool      `gorm:"default:false" json:"show_sunday"`
	CurrentGrade       string    `gorm:"size:32" json:"current_grade"`          // 当前默认展示的年级
	CurrentSemester    string    `gorm:"size:32" json:"current_semester"`       // 当前默认展示的学期
	PeriodSettingsJSON string    `gorm:"type:text" json:"period_settings_json"` // 课程时间设置 JSON: [{"period":1,"start":"08:00","end":"08:40"}]
	BackgroundEmojis   string    `gorm:"size:255" json:"background_emojis"`     // 课表背景 Emoji，逗号分隔
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PeriodSetting struct {
	Period    int    `json:"period"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
