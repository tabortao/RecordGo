package logger

import (
    "go.uber.org/zap"
)

// 中文注释：根据环境初始化 Zap 日志
func New(env string) *zap.Logger {
    if env == "debug" {
        l, _ := zap.NewDevelopment()
        return l
    }
    l, _ := zap.NewProduction()
    return l
}
