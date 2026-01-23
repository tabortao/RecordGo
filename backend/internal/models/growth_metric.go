package models

import "time"

type GrowthMetricRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	MetricType string    `gorm:"size:32;index" json:"metric_type"`
	RecordDate time.Time `gorm:"index" json:"record_date"`
	Value      float64   `json:"value"`
	LeftValue  float64   `json:"left_value"`
	RightValue float64   `json:"right_value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
