package router

import (
    "path/filepath"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "recordgo/internal/config"
    "recordgo/internal/handlers"
)

// 中文注释：初始化 Gin 路由，注册健康检查与未来的模块路由
func New(cfg *config.Config, lg *zap.Logger) *gin.Engine {
    if cfg.Env == "debug" {
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()
    r.Use(gin.Recovery())

    // 中文注释：静态文件服务（上传图片），映射到 /api/uploads
    r.Static("/api/uploads", filepath.Join(cfg.StorageRoot, "uploads"))

    api := r.Group("/api")
    {
        // 中文注释：认证路由（登录/注册）
        api.POST("/auth/register", handlers.Register)
        api.POST("/auth/login", handlers.Login)

        api.GET("/health", handlers.Health)
        // 中文注释：默认心愿占位接口，可用于前端初始化
        api.GET("/wishes/default", handlers.DefaultWishes)

        // 中文注释：心愿管理与兑换记录相关接口
        api.GET("/wishes", handlers.ListWishes)
        api.POST("/wishes", handlers.CreateWish)
        api.PUT("/wishes/:id", handlers.UpdateWish)
        api.DELETE("/wishes/:id", handlers.DeleteWish)
        api.POST("/wishes/:id/exchange", handlers.ExchangeWish)
        api.GET("/wishes/records", handlers.ListWishRecords)
        api.POST("/upload/wish-icon", handlers.UploadWishIcon)
        // 中文注释：任务图片上传接口
        api.POST("/upload/task-image", handlers.UploadTaskImage)

        // 中文注释：任务管理 RESTful 路由
        api.POST("/tasks", handlers.CreateTask)
        api.GET("/tasks", handlers.ListTasks)
        api.GET("/tasks/:id", handlers.GetTask)
        api.PUT("/tasks/:id", handlers.UpdateTask)
        api.PATCH("/tasks/:id/status", handlers.UpdateStatus)
        api.DELETE("/tasks/:id", handlers.DeleteTask)
        api.DELETE("/tasks", handlers.BatchDelete) // ?ids=1,2,3
        api.GET("/tasks/recycle-bin", handlers.ListRecycleBin)
        api.POST("/tasks/recycle-bin/restore", handlers.RestoreTasks) // ?ids=1,2
        api.POST("/tasks/:id/tomato/complete", handlers.CompleteTomato)
    }

    return r
}
