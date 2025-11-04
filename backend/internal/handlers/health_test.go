package handlers

import (
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)

// 中文注释：简单的健康检查单元测试，验证返回结构
func TestHealth(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.GET("/api/health", Health)

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/health", nil)
    r.ServeHTTP(w, req)

    if w.Code != 200 {
        t.Fatalf("期待 200，实际 %d", w.Code)
    }
    body := w.Body.String()
    if !(len(body) > 0 && body[0] == '{') {
        t.Fatalf("响应不是 JSON: %s", body)
    }
}
