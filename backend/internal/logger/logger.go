package logger

import (
    "go.uber.org/zap"
)

// 中文注释：根据环境初始化 Zap 日志
func New(env string) *zap.Logger {
    if env == "debug" {
        l, _ := zap.NewDevelopment()
        // 中文注释：替换全局日志，便于 handlers 使用 zap.L()/zap.S()
        zap.ReplaceGlobals(l)
        return l
    }
    l, _ := zap.NewProduction()
    // 中文注释：替换全局日志，便于 handlers 使用 zap.L()/zap.S()
    zap.ReplaceGlobals(l)
    return l
}
