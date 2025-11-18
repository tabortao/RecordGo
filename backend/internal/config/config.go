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
    AllowedOrigins  []string
    StorageDriver   string
    S3Provider      string
    S3Endpoint      string
    S3Region        string
    S3Bucket        string
    S3AccessKeyID   string
    S3SecretAccessKey string
    S3ForcePathStyle bool
    S3PresignExpire  int64
    S3PublicBaseURL  string
    MaxUploadSizeMB  int
    AllowMIME        string
}

// Load 读取环境变量并提供默认值
func Load() (*Config, error) {
    v := viper.New()
    v.AutomaticEnv()
    // 中文注释：为满足“子账号与主账号金币池共享”的默认规则，开启父子金币同步默认值
    v.SetDefault("PARENT_COINS_SYNC", true)
    v.SetDefault("ALLOW_ORIGINS", "https://recordgo.netlify.app")
    v.SetDefault("STORAGE_DRIVER", "local")
    v.SetDefault("S3_FORCE_PATH_STYLE", false)
    v.SetDefault("S3_PRESIGN_EXPIRE", 300)
    v.SetDefault("MAX_UPLOAD_SIZE_MB", 10)
    v.SetDefault("ALLOW_MIME", "image/webp,audio/mpeg,audio/wav")

    cfg := &Config{
        SecretKey:   v.GetString("SECRET_KEY"),
        DBPath:      defaultString(v.GetString("DB_PATH"), "storage/database/recordgo.db"),
        StorageRoot: defaultString(v.GetString("STORAGE_ROOT"), "storage"),
        Port:        defaultString(v.GetString("PORT"), "8080"),
        Env:         defaultString(v.GetString("GIN_MODE"), "release"),
        // 中文注释：PARENT_COINS_SYNC 读取为布尔值，默认 true（可通过环境变量覆盖为 false）
        ParentCoinsSync: v.GetBool("PARENT_COINS_SYNC"),
        StorageDriver:   defaultString(v.GetString("STORAGE_DRIVER"), "local"),
        S3Provider:      v.GetString("S3_PROVIDER"),
        S3Endpoint:      v.GetString("S3_ENDPOINT"),
        S3Region:        v.GetString("S3_REGION"),
        S3Bucket:        v.GetString("S3_BUCKET"),
        S3AccessKeyID:   v.GetString("S3_ACCESS_KEY_ID"),
        S3SecretAccessKey: v.GetString("S3_SECRET_ACCESS_KEY"),
        S3ForcePathStyle: v.GetBool("S3_FORCE_PATH_STYLE"),
        S3PresignExpire:  v.GetInt64("S3_PRESIGN_EXPIRE"),
        S3PublicBaseURL:  v.GetString("S3_PUBLIC_BASE_URL"),
        MaxUploadSizeMB:  v.GetInt("MAX_UPLOAD_SIZE_MB"),
        AllowMIME:        v.GetString("ALLOW_MIME"),
    }

	originsRaw := v.GetString("ALLOW_ORIGINS")
	cfg.AllowedOrigins = splitAndTrim(originsRaw)

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

func splitAndTrim(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			// trim spaces
			l := 0
			r := len(part)
			for l < r && (part[l] == ' ' || part[l] == '\t' || part[l] == '\n' || part[l] == '\r') {
				l++
			}
			for r > l && (part[r-1] == ' ' || part[r-1] == '\t' || part[r-1] == '\n' || part[r-1] == '\r') {
				r--
			}
			if r > l {
				out = append(out, part[l:r])
			}
			start = i + 1
		}
	}
	if len(out) == 0 {
		return []string{"https://recordgo.netlify.app"}
	}
	return out
}
