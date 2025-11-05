package handlers

import (
    "encoding/json"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

// 中文注释：任务创建请求结构体（前端字段需与后端一致）
type CreateTaskReq struct {
    UserID       uint      `json:"user_id"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Priority     string    `json:"priority"` // 中文注释：占位字段，实际存入 Category/Remark
    Category     string    `json:"category"`
    Score        int       `json:"score"`
    PlanMinutes  int       `json:"plan_minutes"`
    StartDate    time.Time `json:"start_date"`
    EndDate      *time.Time `json:"end_date"`
    ImageJSON    string    `json:"image_json"` // 中文注释：任务图片 JSON 数组字符串
}

// 中文注释：编辑请求结构体（允许全部可变更字段）
type UpdateTaskReq struct {
    Name         *string    `json:"name"`
    Description  *string    `json:"description"`
    Category     *string    `json:"category"`
    Score        *int       `json:"score"`
    PlanMinutes  *int       `json:"plan_minutes"`
    ActualMinutes *int      `json:"actual_minutes"`
    StartDate    *time.Time `json:"start_date"`
    EndDate      **time.Time `json:"end_date"`
    Remark       *string    `json:"remark"`
    ImageJSON    *string    `json:"image_json"`
}

// 中文注释：状态变更请求结构体
type UpdateStatusReq struct {
    Status int `json:"status"` // 0待完成 1进行中 2已完成
}

// 中文注释：番茄钟完成上报请求结构体
type TomatoCompleteReq struct {
    Minutes int `json:"minutes"` // 中文注释：一次番茄钟的工作分钟数（默认25）
}

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
    var req CreateTaskReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    if req.Name == "" || req.PlanMinutes <= 0 {
        common.Error(c, 40002, "任务名称与计划时长必填且合法")
        return
    }
    t := models.Task{
        UserID: req.UserID,
        Name: req.Name,
        Description: req.Description,
        Category: req.Category,
        Score: req.Score,
        PlanMinutes: req.PlanMinutes,
        StartDate: req.StartDate,
        EndDate: req.EndDate,
        ImageJSON: req.ImageJSON,
        Status: 0,
    }
    if err := db.DB().Create(&t).Error; err != nil {
        common.Error(c, 50001, "创建任务失败")
        return
    }
    // 历史记录：create
    snapshot, _ := json.Marshal(t)
    _ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "create", SnapshotJSON: string(snapshot)}).Error
    common.Ok(c, t)
}

// ListTasks 列表查询（支持状态与分页）
func ListTasks(c *gin.Context) {
    var tasks []models.Task
    q := db.DB().Model(&models.Task{})

    // 过滤状态
    if s := c.Query("status"); s != "" {
        if si, err := strconv.Atoi(s); err == nil {
            q = q.Where("status = ?", si)
        }
    }
    // 简单分页
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    if page < 1 { page = 1 }
    if size < 1 || size > 100 { size = 20 }

    var total int64
    q.Count(&total)
    if err := q.Order("created_at DESC").Offset((page-1)*size).Limit(size).Find(&tasks).Error; err != nil {
        common.Error(c, 50002, "查询任务失败")
        return
    }
    common.Ok(c, gin.H{"items": tasks, "total": total, "page": page, "page_size": size})
}

// GetTask 获取单个任务
func GetTask(c *gin.Context) {
    id := c.Param("id")
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            common.Error(c, 40401, "任务不存在")
            return
        }
        common.Error(c, 50003, "查询任务失败")
        return
    }
    common.Ok(c, t)
}

// UpdateTask 编辑任务（记录历史）
func UpdateTask(c *gin.Context) {
    id := c.Param("id")
    var req UpdateTaskReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil {
        common.Error(c, 40401, "任务不存在")
        return
    }
    // 修改前快照
    before, _ := json.Marshal(t)

    // 应用变更
    if req.Name != nil { t.Name = *req.Name }
    if req.Description != nil { t.Description = *req.Description }
    if req.Category != nil { t.Category = *req.Category }
    if req.Score != nil { t.Score = *req.Score }
    if req.PlanMinutes != nil { t.PlanMinutes = *req.PlanMinutes }
    if req.ActualMinutes != nil { t.ActualMinutes = *req.ActualMinutes }
    if req.StartDate != nil { t.StartDate = *req.StartDate }
    if req.EndDate != nil { t.EndDate = *req.EndDate }
    if req.Remark != nil { t.Remark = *req.Remark }
    if req.ImageJSON != nil { t.ImageJSON = *req.ImageJSON }

    if err := db.DB().Save(&t).Error; err != nil {
        common.Error(c, 50004, "更新任务失败")
        return
    }
    // 历史记录：update，保存修改前快照
    _ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "update", SnapshotJSON: string(before)}).Error
    common.Ok(c, t)
}

// UpdateStatus 状态变更（记录历史）
func UpdateStatus(c *gin.Context) {
    id := c.Param("id")
    var req UpdateStatusReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    if req.Status < 0 || req.Status > 2 {
        common.Error(c, 40003, "非法状态")
        return
    }
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil {
        common.Error(c, 40401, "任务不存在")
        return
    }
    // 快照
    before, _ := json.Marshal(t)
    // 中文注释：当任务状态从“已完成”变更为其他状态时，扣减用户金币；当从非已完成变更为“已完成”时，增加用户金币
    prev := t.Status
    t.Status = req.Status
    // 同步用户金币
    var u models.User
    _ = db.DB().First(&u, t.UserID).Error
    if prev == 2 && req.Status != 2 {
        // 取消完成：扣减对应金币（允许为负数，实际效果为增加金币）
        u.Coins -= int64(t.Score)
    } else if prev != 2 && req.Status == 2 {
        // 标记完成：增加对应金币
        u.Coins += int64(t.Score)
    }
    if err := db.DB().Save(&u).Error; err != nil {
        common.Error(c, 50005, "更新用户金币失败")
        return
    }
    if err := db.DB().Save(&t).Error; err != nil {
        common.Error(c, 50005, "更新状态失败")
        return
    }
    _ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "status", SnapshotJSON: string(before)}).Error
    common.Ok(c, gin.H{"task": t, "user_coins": u.Coins})
}

// DeleteTask 软删除（进入回收站）
func DeleteTask(c *gin.Context) {
    id := c.Param("id")
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil {
        common.Error(c, 40401, "任务不存在")
        return
    }
    if err := db.DB().Delete(&t).Error; err != nil {
        common.Error(c, 50006, "删除失败")
        return
    }
    // 历史记录：delete
    snapshot, _ := json.Marshal(t)
    _ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "delete", SnapshotJSON: string(snapshot)}).Error
    common.Ok(c, gin.H{"id": t.ID})
}

// BatchDelete 批量软删除
func BatchDelete(c *gin.Context) {
    idsStr := c.Query("ids") // 逗号分隔
    if idsStr == "" {
        common.Error(c, 40004, "缺少 ids 参数")
        return
    }
    ids := []uint{}
    for _, s := range splitComma(idsStr) {
        if v, err := strconv.Atoi(s); err == nil {
            ids = append(ids, uint(v))
        }
    }
    if len(ids) == 0 {
        common.Error(c, 40005, "无有效 ID")
        return
    }
    if err := db.DB().Where("id IN ?", ids).Delete(&models.Task{}).Error; err != nil {
        common.Error(c, 50007, "批量删除失败")
        return
    }
    common.Ok(c, gin.H{"deleted": len(ids)})
}

// ListRecycleBin 回收站任务
func ListRecycleBin(c *gin.Context) {
    var tasks []models.Task
    if err := db.DB().Unscoped().Where("deleted_at IS NOT NULL").Order("deleted_at DESC").Find(&tasks).Error; err != nil {
        common.Error(c, 50008, "查询回收站失败")
        return
    }
    common.Ok(c, tasks)
}

// RestoreTasks 从回收站恢复
func RestoreTasks(c *gin.Context) {
    idsStr := c.Query("ids")
    if idsStr == "" {
        common.Error(c, 40004, "缺少 ids 参数")
        return
    }
    ids := []uint{}
    for _, s := range splitComma(idsStr) {
        if v, err := strconv.Atoi(s); err == nil {
            ids = append(ids, uint(v))
        }
    }
    for _, id := range ids {
        var t models.Task
        if err := db.DB().Unscoped().First(&t, id).Error; err == nil {
            // 恢复软删除
            t.DeletedAt = gorm.DeletedAt{}
            if err := db.DB().Unscoped().Save(&t).Error; err == nil {
                snapshot, _ := json.Marshal(t)
                _ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "restore", SnapshotJSON: string(snapshot)}).Error
            }
        }
    }
    common.Ok(c, gin.H{"restored": len(ids)})
}

// CompleteTomato 番茄钟完成上报（累加任务番茄数与实际耗时）
func CompleteTomato(c *gin.Context) {
    id := c.Param("id")
    var req TomatoCompleteReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    if req.Minutes <= 0 { req.Minutes = 25 }
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil {
        common.Error(c, 40401, "任务不存在")
        return
    }
    // 更新番茄钟次数与实际耗时
    before, _ := json.Marshal(t)
    t.TomatoCount += 1
    t.ActualMinutes += req.Minutes
    // 若完成一次番茄，状态可自动改为进行中或完成由前端控制，这里不强制
    if err := db.DB().Save(&t).Error; err != nil {
        common.Error(c, 50009, "番茄上报失败")
        return
    }
    _ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "tomato", SnapshotJSON: string(before)}).Error
    common.Ok(c, t)
}

// 工具函数：按逗号分割并去空
func splitComma(s string) []string {
    out := []string{}
    start := 0
    for i := 0; i <= len(s); i++ {
        if i == len(s) || s[i] == ',' {
            seg := s[start:i]
            if seg != "" { out = append(out, seg) }
            start = i+1
        }
    }
    return out
}
