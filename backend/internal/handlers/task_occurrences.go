package handlers

import (
    "encoding/json"
    "time"
    "strings"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/db"
    "recordgo/internal/models"
    "recordgo/internal/config"
    "fmt"
    "go.uber.org/zap"
)

func normalizeDate(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func parseDate(s string) (time.Time, error) {
    if s == "" { return time.Time{}, gorm.ErrInvalidData }
    t, err := time.Parse("2006-01-02", s)
    if err != nil { return time.Time{}, err }
    return t, nil
}

func weeklySetFromJSON(s string) map[int]struct{} {
    out := map[int]struct{}{}
    s = strings.TrimSpace(s)
    if s == "" { return out }
    var arr []int
    if err := json.Unmarshal([]byte(s), &arr); err == nil {
        for _, v := range arr { if v >= 1 && v <= 7 { out[v] = struct{}{} } }
    }
    return out
}

func matchRepeat(t models.Task, d time.Time) bool {
    d = normalizeDate(d)
    s := normalizeDate(t.StartDate)
    if d.Before(s) { return false }
    if t.EndDate != nil {
        e := normalizeDate(*t.EndDate)
        if d.After(e) { return false }
    }
    typ := strings.ToLower(strings.TrimSpace(t.Repeat))
    if typ == "none" || typ == "" { return d.Equal(s) }
    if typ == "daily" { return true }
    if typ == "weekdays" {
        wd := d.Weekday()
        return wd >= time.Monday && wd <= time.Friday
    }
    if typ == "weekly" {
        set := weeklySetFromJSON(t.RepeatDaysJSON)
        w := int(d.Weekday()); if w == 0 { w = 7 }
        if len(set) == 0 {
            // 默认按开始日星期
            sw := int(s.Weekday()); if sw == 0 { sw = 7 }
            return w == sw
        }
        _, ok := set[w]
        return ok
    }
    if typ == "monthly" {
        dom := s.Day()
        return d.Day() == dom
    }
    return false
}

type OccurrenceCompleteReq struct {
    Date    string `json:"date"`
    Minutes int    `json:"minutes"`
}

func CompleteTaskOccurrence(c *gin.Context) {
    id := c.Param("id")
    if !hasPermission(c, "tasks", "status") {
        deny(c, "无权限修改任务状态")
        return
    }
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil { common.Error(c, 40401, "任务不存在"); return }
    if !canAccessUser(c, t.UserID) { deny(c, "无权限"); return }
    var req OccurrenceCompleteReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Date) == "" { common.Error(c, 40001, "参数错误"); return }
    d, err := parseDate(req.Date)
    if err != nil { common.Error(c, 40001, "日期格式错误"); return }
    if !matchRepeat(t, d) {
        zap.L().Warn("CompleteTaskOccurrence: repeat mismatch",
            zap.Uint("task_id", t.ID),
            zap.String("repeat", strings.ToLower(strings.TrimSpace(t.Repeat))),
            zap.Time("start", t.StartDate),
            zap.String("date", req.Date),
            zap.String("weekly_days", t.RepeatDaysJSON),
            zap.Any("end", t.EndDate),
        )
        common.Error(c, 40003, "日期不在重复规则范围内")
        return
    }
    m := req.Minutes
    if m <= 0 { m = t.PlanMinutes }
    occ := models.TaskOccurrence{}
    if err := db.DB().Where("task_id = ? AND date = ?", t.ID, normalizeDate(d)).First(&occ).Error; err != nil {
        occ.TaskID = t.ID
        occ.Date = normalizeDate(d)
    }
    occ.Status = 2
    occ.Minutes = m
    if occ.ID == 0 { _ = db.DB().Create(&occ).Error } else { _ = db.DB().Save(&occ).Error }
    var owner models.User
    if err := db.DB().First(&owner, t.UserID).Error; err != nil {
        owner = models.User{ID: t.UserID, Username: fmt.Sprintf("user%d", t.UserID), Role: "user", Coins: 0, Nickname: "测试用户"}
        if ierr := db.DB().Create(&owner).Error; ierr != nil { common.Error(c, 50005, "初始化用户失败"); return }
    }
    target := owner
    if cfg2, _ := config.Load(); cfg2 != nil && cfg2.ParentCoinsSync && owner.ParentID != nil {
        var parent models.User
        if err := db.DB().First(&parent, *owner.ParentID).Error; err == nil { target = parent }
    }
    target.Coins += int64(t.Score)
    if err := db.DB().Save(&target).Error; err != nil { common.Error(c, 50005, "更新用户金币失败"); return }
    common.Ok(c, gin.H{"task_id": t.ID, "date": req.Date, "status": 2, "minutes": m, "user_coins": target.Coins})
}

type OccurrenceUncompleteReq struct { Date string `json:"date"` }

func UncompleteTaskOccurrence(c *gin.Context) {
    id := c.Param("id")
    if !hasPermission(c, "tasks", "status") { deny(c, "无权限修改任务状态"); return }
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil { common.Error(c, 40401, "任务不存在"); return }
    if !canAccessUser(c, t.UserID) { deny(c, "无权限"); return }
    var req OccurrenceUncompleteReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Date) == "" { common.Error(c, 40001, "参数错误"); return }
    d, err := parseDate(req.Date)
    if err != nil { common.Error(c, 40001, "日期格式错误"); return }
    _ = db.DB().Where("task_id = ? AND date = ?", t.ID, normalizeDate(d)).Delete(&models.TaskOccurrence{}).Error
    var owner models.User
    if err := db.DB().First(&owner, t.UserID).Error; err != nil {
        owner = models.User{ID: t.UserID, Username: fmt.Sprintf("user%d", t.UserID), Role: "user", Coins: 0, Nickname: "测试用户"}
        if ierr := db.DB().Create(&owner).Error; ierr != nil { common.Error(c, 50005, "初始化用户失败"); return }
    }
    target := owner
    if cfg2, _ := config.Load(); cfg2 != nil && cfg2.ParentCoinsSync && owner.ParentID != nil {
        var parent models.User
        if err := db.DB().First(&parent, *owner.ParentID).Error; err == nil { target = parent }
    }
    target.Coins -= int64(t.Score)
    if err := db.DB().Save(&target).Error; err != nil { common.Error(c, 50005, "更新用户金币失败"); return }
    common.Ok(c, gin.H{"task_id": t.ID, "date": req.Date, "status": 0, "user_coins": target.Coins})
}

type OccurrenceDeleteReq struct { Date string `json:"date"` }

func DeleteTaskOccurrence(c *gin.Context) {
    id := c.Param("id")
    if !hasPermission(c, "tasks", "status") { deny(c, "无权限修改任务状态"); return }
    var t models.Task
    if err := db.DB().First(&t, id).Error; err != nil { common.Error(c, 40401, "任务不存在"); return }
    if !canAccessUser(c, t.UserID) { deny(c, "无权限"); return }
    var req OccurrenceDeleteReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Date) == "" { common.Error(c, 40001, "参数错误"); return }
    d, err := parseDate(req.Date)
    if err != nil { common.Error(c, 40001, "日期格式错误"); return }
    if !matchRepeat(t, d) { common.Error(c, 40003, "日期不在重复规则范围内"); return }
    occ := models.TaskOccurrence{}
    if err := db.DB().Where("task_id = ? AND date = ?", t.ID, normalizeDate(d)).First(&occ).Error; err != nil {
        occ.TaskID = t.ID
        occ.Date = normalizeDate(d)
    }
    occ.Status = -1
    occ.Minutes = 0
    if occ.ID == 0 { _ = db.DB().Create(&occ).Error } else { _ = db.DB().Save(&occ).Error }
    common.Ok(c, gin.H{"task_id": t.ID, "date": req.Date, "status": -1})
}

func ListTaskOccurrences(c *gin.Context) {
    cl := extractClaims(c)
    if cl == nil { common.Error(c, 40100, "未登录"); return }
    uid := cl.UserID
    if cl.ParentID != nil { uid = *cl.ParentID }
    startStr := c.Query("start")
    endStr := c.Query("end")
    dateStr := c.Query("date")
    var start, end time.Time
    var err error
    if dateStr != "" {
        d, e := parseDate(dateStr); if e != nil { common.Error(c, 40001, "日期格式错误"); return }
        start = d; end = d
    } else {
        if startStr == "" || endStr == "" { common.Error(c, 40001, "缺少日期范围"); return }
        start, err = parseDate(startStr); if err != nil { common.Error(c, 40001, "开始日期错误"); return }
        end, err = parseDate(endStr); if err != nil { common.Error(c, 40001, "结束日期错误"); return }
    }
    start = normalizeDate(start); end = normalizeDate(end)
    var tasks []models.Task
    if err := db.DB().Where("user_id = ?", uid).Find(&tasks).Error; err != nil { common.Error(c, 50002, "查询任务失败"); return }
    ids := make([]uint, 0, len(tasks))
    for _, t := range tasks { ids = append(ids, t.ID) }
    var list []models.TaskOccurrence
    if len(ids) > 0 {
        _ = db.DB().Where("task_id IN ? AND date BETWEEN ? AND ?", ids, start, end).Find(&list).Error
    }
    common.Ok(c, gin.H{"items": list})
}