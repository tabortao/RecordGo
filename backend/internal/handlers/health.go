package handlers

import (
    "github.com/gin-gonic/gin"
    "recordgo/internal/common"
)

// 中文注释：健康检查接口，返回 code=0
func Health(c *gin.Context) {
    common.Ok(c, gin.H{"status": "healthy"})
}
