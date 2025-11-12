package config

import (
    "fmt"
    "github.com/spf13/viper"
)

// 中文注释：配置结构体，统一管理环境变量与默认值
type Config struct {
    SecretKey   string
    DBPath      string
    StorageRoot string
    Port        string
    Env         string // release 或 debug
    // 中文注释：是否启用父子账号金币池同步（子账号显示与操作均使用父账号金币）
    ParentCoinsSync bool
}

// Load 读取环境变量并提供默认值
func Load() (*Config, error) {
    v := viper.New()
    v.AutomaticEnv()
    // 中文注释：为满足“子账号与主账号金币池共享”的默认规则，开启父子金币同步默认值
    v.SetDefault("PARENT_COINS_SYNC", true)

    cfg := &Config{
        SecretKey:   v.GetString("SECRET_KEY"),
        DBPath:      defaultString(v.GetString("DB_PATH"), "storage/database/recordgo.db"),
        StorageRoot: defaultString(v.GetString("STORAGE_ROOT"), "storage"),
        Port:        defaultString(v.GetString("PORT"), "8080"),
        Env:         defaultString(v.GetString("GIN_MODE"), "release"),
        // 中文注释：PARENT_COINS_SYNC 读取为布尔值，默认 true（可通过环境变量覆盖为 false）
        ParentCoinsSync: v.GetBool("PARENT_COINS_SYNC"),
    }

    if cfg.SecretKey == "" {
        return nil, fmt.Errorf("SECRET_KEY 未设置")
    }
    return cfg, nil
}

func defaultString(val, def string) string {
    if val == "" {
        return def
    }
    return val
}
