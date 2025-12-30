package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"recordgo/internal/common"
	"time"

	"github.com/gin-gonic/gin"
)

// OCRProxyReq 代理请求结构体
type OCRProxyReq struct {
	TargetURL string                 `json:"target_url" binding:"required"`
	Token     string                 `json:"token" binding:"required"`
	Payload   map[string]interface{} `json:"payload" binding:"required"`
}

// HandleOCRProxy 代理 OCR 请求以解决跨域问题
func HandleOCRProxy(c *gin.Context) {
	var req OCRProxyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误：需要 target_url, token 和 payload")
		return
	}

	// 序列化 payload
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		common.Error(c, 40002, "Payload 序列化失败")
		return
	}

	// 创建代理请求
	proxyReq, err := http.NewRequest("POST", req.TargetURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		common.Error(c, 50001, "创建请求失败: "+err.Error())
		return
	}

	// 设置 Header (模拟前端直接请求时的 Header)
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("Authorization", "token "+req.Token)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		common.Error(c, 50002, "请求目标服务失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		common.Error(c, 50003, "读取响应失败: "+err.Error())
		return
	}

	// Log the raw response from OCR service for debugging
	// fmt.Printf("[OCR Proxy] Raw Response: %s\n", string(bodyBytes))

	// 解析目标响应为 JSON (为了放入 common.Ok 的 data 中)
	var targetResp interface{}
	if err := json.Unmarshal(bodyBytes, &targetResp); err != nil {
		// 如果不是 JSON，直接返回字符串
		targetResp = string(bodyBytes)
	}

	if resp.StatusCode != 200 {
		// 即使目标服务报错，我们也返回给前端，让前端知道发生了什么
		// 这里使用 50004 状态码表示上游错误
		c.JSON(http.StatusOK, common.Response{
			Code:    50004,
			Message: "OCR 服务返回错误: " + resp.Status,
			Data:    targetResp,
		})
		return
	}

	common.Ok(c, targetResp)
}
