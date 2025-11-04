package models

import "time"

// 中文注释：任务编辑历史模型，记录每次修改的快照与类型
type TaskHistory struct {
    ID        uint      `gorm:"primaryKey"`
    TaskID    uint      `gorm:"index"`
    ChangeType string   `gorm:"size:32"` // create/update/status/delete/restore
    SnapshotJSON string `gorm:"type:text"` // 中文注释：修改前或修改后的任务 JSON 快照
    CreatedAt time.Time
}

