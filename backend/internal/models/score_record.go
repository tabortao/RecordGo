package models

import "time"

type ScoreRecord struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Subject   string     `gorm:"size:64" json:"subject"`
	ExamName  string     `gorm:"size:128" json:"exam_name"`
	ExamDate  time.Time  `gorm:"index" json:"exam_date"`
	Score     float64    `json:"score"`
	FullScore float64    `json:"full_score"`
	Rank      *int       `json:"rank"`
	ClassAvg  *float64   `json:"class_avg"`
	PhotoURL  string     `gorm:"size:512" json:"photo_url"`
	Remark    string     `gorm:"size:512" json:"remark"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
