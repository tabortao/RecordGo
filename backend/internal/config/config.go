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
}

// Load 读取环境变量并提供默认值
func Load() (*Config, error) {
    v := viper.New()
    v.AutomaticEnv()

    cfg := &Config{
        SecretKey:   v.GetString("SECRET_KEY"),
        DBPath:      defaultString(v.GetString("DB_PATH"), "storage/database/recordgo.db"),
        StorageRoot: defaultString(v.GetString("STORAGE_ROOT"), "storage"),
        Port:        defaultString(v.GetString("PORT"), "8080"),
        Env:         defaultString(v.GetString("GIN_MODE"), "release"),
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
