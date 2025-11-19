package handlers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

type TaskCategoriesPayload struct {
    Items []models.TaskCategory `json:"items"`
}

func ListTaskCategories(c *gin.Context) {
    var list []models.TaskCategory
    _ = db.DB().Order("`order` ASC, name ASC").Find(&list).Error
    common.Ok(c, gin.H{"items": list})
}

func PutTaskCategories(c *gin.Context) {
    var payload TaskCategoriesPayload
    var arr []models.TaskCategory
    if err := c.ShouldBindJSON(&payload); err == nil && len(payload.Items) > 0 {
        arr = payload.Items
    } else {
        var tmp []models.TaskCategory
        if err := c.ShouldBindJSON(&tmp); err == nil {
            arr = tmp
        }
    }
    if len(arr) == 0 {
        common.Error(c, 40001, "参数错误")
        return
    }
    _ = db.DB().Transaction(func(tx *gorm.DB) error {
        if err := tx.Exec("DELETE FROM task_categories").Error; err != nil { return err }
        for i, it := range arr {
            item := models.TaskCategory{ Name: it.Name, Color: it.Color, Order: it.Order }
            if item.Order <= 0 { item.Order = i + 1 }
            if err := tx.Create(&item).Error; err != nil { return err }
        }
        return nil
    })
    common.Ok(c, gin.H{"saved": len(arr)})
}

type PatchOrderReq struct {
    Name  string `json:"name"`
    Order int    `json:"order"`
}

func PatchTaskCategoryOrder(c *gin.Context) {
    var req PatchOrderReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" || req.Order <= 0 {
        common.Error(c, 40001, "参数错误")
        return
    }
    var cat models.TaskCategory
    if err := db.DB().Where("name = ?", req.Name).First(&cat).Error; err != nil {
        common.Error(c, 40401, "分类不存在")
        return
    }
    cat.Order = req.Order
    if err := db.DB().Save(&cat).Error; err != nil {
        common.Error(c, 50012, "更新失败")
        return
    }
    common.Ok(c, gin.H{"name": cat.Name, "order": cat.Order})
}