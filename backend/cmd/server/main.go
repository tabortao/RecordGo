package main

import (
    "log"
    "recordgo/internal/config"
    "recordgo/internal/logger"
    "recordgo/internal/router"
    "recordgo/internal/storage"
    "recordgo/internal/db"
)

// 中文注释：后端入口，初始化配置、日志、存储目录与路由，启动 HTTP 服务
func main() {
    // 初始化配置（环境变量优先）
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }

    // 初始化日志
    lg := logger.New(cfg.Env)

    // 初始化存储目录（确保 storage 路径存在）
    if err := storage.EnsurePaths(cfg.StorageRoot); err != nil {
        lg.Sugar().Fatalf("初始化存储目录失败: %v", err)
    }

    // 初始化数据库（SQLite 纯 Go 驱动）
    if _, err := db.Init(cfg, lg); err != nil {
        lg.Sugar().Fatalf("数据库初始化失败: %v", err)
    }

    // 初始化路由与服务
    r := router.New(cfg, lg)
    lg.Sugar().Infof("RecordGo API 启动，端口: %s", cfg.Port)
    if err := r.Run(":" + cfg.Port); err != nil {
        lg.Sugar().Fatalf("服务启动失败: %v", err)
    }
}
