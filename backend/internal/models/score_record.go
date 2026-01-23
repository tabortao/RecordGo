package models

import "time"

type ScoreRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Subject      string    `gorm:"size:64" json:"subject"`
	ExamName     string    `gorm:"size:128" json:"exam_name"`
	ExamType     string    `gorm:"size:64" json:"exam_type"`
	Grade        string    `gorm:"size:32" json:"grade"`
	Term         string    `gorm:"size:32" json:"term"`
	Topic        string    `gorm:"size:128" json:"topic"`
	Difficulty   string    `gorm:"size:32" json:"difficulty"`
	ExamDate     time.Time `gorm:"index" json:"exam_date"`
	Score        float64   `json:"score"`
	FullScore    float64   `json:"full_score"`
	ClassRank    *int      `json:"class_rank"`
	RankType     string    `gorm:"size:32" json:"rank_type"`
	GradeRank    *int      `json:"grade_rank"`
	ClassAvg     *float64  `json:"class_avg"`
	ClassHighest *float64  `json:"class_highest"`
	PhotoURL     string    `gorm:"size:512" json:"photo_url"`
	Remark       string    `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
