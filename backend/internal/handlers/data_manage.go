package handlers

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/config"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

type clearDataReq struct { Password string `json:"password"` }

func ClearData(c *gin.Context) {
    cl := extractClaims(c)
    if cl == nil { common.Error(c, 40100, "未登录或令牌无效"); return }
    var req clearDataReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Password) == "" { common.Error(c, 40001, "参数错误"); return }
    var u models.User
    if err := db.DB().First(&u, cl.UserID).Error; err != nil { common.Error(c, 40401, "用户不存在"); return }
    h := sha256.Sum256([]byte(req.Password))
    if u.PasswordSha != hex.EncodeToString(h[:]) { common.Error(c, 40003, "密码错误"); return }
    var tasks []models.Task
    if err := db.DB().Where("user_id = ?", u.ID).Find(&tasks).Error; err != nil { common.Error(c, 50021, "查询任务失败"); return }
    ids := make([]uint, 0, len(tasks))
    for _, t := range tasks { ids = append(ids, t.ID) }
    err := db.DB().Transaction(func(tx *gorm.DB) error {
        for _, t := range tasks {
            if strings.TrimSpace(t.ImageJSON) != "" {
                var paths []string
                _ = jsonUnmarshal(t.ImageJSON, &paths)
                if len(paths) > 0 {
                    root := os.Getenv("STORAGE_ROOT")
                    if strings.TrimSpace(root) == "" { root = "storage" }
                    for _, rel := range paths {
                        full := filepath.Join(root, filepath.FromSlash(rel))
                        _ = os.Remove(full)
                    }
                }
            }
            root := os.Getenv("STORAGE_ROOT")
            if strings.TrimSpace(root) == "" { root = "storage" }
            dir := filepath.Join(root, "uploads", "images", "task_images", fmt.Sprintf("%d", t.UserID))
            entries, readErr := os.ReadDir(dir)
            if readErr == nil {
                prefix := fmt.Sprintf("task_%d_", t.ID)
                for _, ent := range entries {
                    if !ent.IsDir() && strings.HasPrefix(ent.Name(), prefix) {
                        _ = os.Remove(filepath.Join(dir, ent.Name()))
                    }
                }
            }
        }
        if len(ids) > 0 {
            if err := tx.Where("task_id IN ?", ids).Delete(&models.TaskOccurrence{}).Error; err != nil { return err }
        }
        if err := tx.Where("user_id = ?", u.ID).Delete(&models.Task{}).Error; err != nil { return err }
        if err := tx.Where("user_id = ?", u.ID).Delete(&models.WishRecord{}).Error; err != nil { return err }
        if err := tx.Where("user_id = ?", u.ID).Delete(&models.Wish{}).Error; err != nil { return err }
        target := &u
        if cfg, _ := config.Load(); cfg != nil && cfg.ParentCoinsSync && u.ParentID != nil {
            var parent models.User
            if err := tx.First(&parent, *u.ParentID).Error; err == nil { target = &parent }
        }
        target.Coins = 0
        if err := tx.Save(target).Error; err != nil { return err }
        return nil
    })
    if err != nil { common.Error(c, 50022, "清除数据失败"); return }
    targetCoins := u.Coins
    if cfg, _ := config.Load(); cfg != nil && u.ParentID != nil && cfg.ParentCoinsSync {
        var parent models.User
        if err := db.DB().First(&parent, *u.ParentID).Error; err == nil { targetCoins = parent.Coins }
    }
    common.Ok(c, gin.H{"user_coins": targetCoins})
}

func jsonUnmarshal(s string, v interface{}) error { return json.Unmarshal([]byte(strings.TrimSpace(s)), v) }