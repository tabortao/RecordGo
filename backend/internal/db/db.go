package db

import (
    "recordgo/internal/config"
    "recordgo/internal/models"
    "go.uber.org/zap"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/glebarez/sqlite"
)

var globalDB *gorm.DB

// 中文注释：初始化 GORM 连接（纯 Go SQLite 驱动）并自动迁移基础模型
func Init(cfg *config.Config, lg *zap.Logger) (*gorm.DB, error) {
    gormLogger := logger.Default
    if cfg.Env == "debug" {
        gormLogger = logger.Default.LogMode(logger.Info)
    }

    db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{Logger: gormLogger})
    if err != nil {
        return nil, err
    }

    // 自动迁移基础模型（后续可扩展）
    if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.Wish{}, &models.UserSettings{}, &models.TaskHistory{}, &models.WishRecord{}, &models.TaskCategory{}, &models.TaskOccurrence{}); err != nil {
        return nil, err
    }

    globalDB = db
    lg.Sugar().Infof("SQLite 初始化成功: %s", cfg.DBPath)
    return db, nil
}

// DB 返回全局数据库对象
func DB() *gorm.DB {
    return globalDB
}
