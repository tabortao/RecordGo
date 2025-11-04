package common

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// 中文注释：统一响应结构体，满足 {code, message, data}
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// Ok 成功响应
func Ok(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// Error 失败响应
func Error(c *gin.Context, code int, msg string) {
    c.JSON(http.StatusOK, Response{Code: code, Message: msg, Data: nil})
}
