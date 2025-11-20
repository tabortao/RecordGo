package handlers

import (
    "time"
    "github.com/gin-gonic/gin"
    "recordgo/internal/common"
)

// 中文注释：健康检查接口，返回 code=0
func Health(c *gin.Context) {
    now := time.Now()
    common.Ok(c, gin.H{
        "status": "healthy",
        "timezone": now.Location().String(),
        "now_local": now.Format(time.RFC3339),
        "now_utc": now.UTC().Format(time.RFC3339),
    })
}
