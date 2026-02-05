package handlers

import (
    "encoding/json"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

type BatchCreateReq struct {
    UserID      uint      `json:"user_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Score       int       `json:"score"`
    ScoreMode   string    `json:"score_mode"`
    PlanMinutes int       `json:"plan_minutes"`
    StartDate   time.Time `json:"start_date"`
    EndDate     *time.Time `json:"end_date"`
    RepeatType  string    `json:"repeat_type"` // none/daily/weekdays/weekly/monthly
    WeeklyDays  []int     `json:"weekly_days"` // 1..7 (周一..周日)
}

func normalizeRepeatType(s string) string {
    t := strings.ToLower(strings.TrimSpace(s))
    switch t {
    case "daily":
        return "daily"
    case "weekdays":
        return "weekdays"
    case "weekly":
        return "weekly"
    case "monthly":
        return "monthly"
    default:
        return "none"
    }
}

func genDates(start time.Time, end *time.Time, repeatType string, weeklyDays []int) []time.Time {
    if end == nil || repeatType == "none" { return []time.Time{time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)} }
    s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
    e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
    if e.Before(s) { return []time.Time{s} }
    out := make([]time.Time, 0, 32)
    if repeatType == "daily" {
        for d := s; !d.After(e); d = d.AddDate(0, 0, 1) { out = append(out, d) }
        return out
    }
    if repeatType == "weekdays" {
        for d := s; !d.After(e); d = d.AddDate(0, 0, 1) {
            wd := d.Weekday()
            if wd >= time.Monday && wd <= time.Friday { out = append(out, d) }
        }
        return out
    }
    if repeatType == "weekly" {
        set := map[int]struct{}{}
        if len(weeklyDays) == 0 {
            w := int(s.Weekday())
            if w == 0 { w = 7 }
            set[w] = struct{}{}
        } else {
            for _, v := range weeklyDays { if v >= 1 && v <= 7 { set[v] = struct{}{} } }
        }
        for d := s; !d.After(e); d = d.AddDate(0, 0, 1) {
            w := int(d.Weekday()); if w == 0 { w = 7 }
            if _, ok := set[w]; ok { out = append(out, d) }
        }
        return out
    }
    if repeatType == "monthly" {
        dom := s.Day()
        cur := time.Date(s.Year(), s.Month(), 1, 0, 0, 0, 0, s.Location())
        for !cur.After(e) {
            cand := time.Date(cur.Year(), cur.Month(), dom, 0, 0, 0, 0, s.Location())
            if cand.Month() == cur.Month() && !cand.After(e) && !cand.Before(s) { out = append(out, cand) }
            cur = cur.AddDate(0, 1, 0)
        }
        return out
    }
    return []time.Time{s}
}

// CreateTasksBatch 批量创建重复任务
func CreateTasksBatch(c *gin.Context) {
    if !hasPermission(c, "tasks", "create") {
        deny(c, "无权限创建任务")
        return
    }
    var req BatchCreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    if strings.TrimSpace(req.Name) == "" || req.PlanMinutes <= 0 {
        common.Error(c, 40002, "任务名称与计划时长必填且合法")
        return
    }
    cl := extractClaims(c)
    if cl == nil {
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    uid := req.UserID
    if uid == 0 { if cl.ParentID != nil { uid = *cl.ParentID } else { uid = cl.UserID } }
    if !canAccessUser(c, uid) {
        deny(c, "无权限为该用户创建任务")
        return
    }
    rtype := normalizeRepeatType(req.RepeatType)
    dates := genDates(req.StartDate, req.EndDate, rtype, req.WeeklyDays)
    if len(dates) == 0 { dates = []time.Time{req.StartDate} }
    if len(dates) > 500 {
        common.Error(c, 40006, "重复范围过大，请缩短截止日期")
        return
    }
    zap.L().Info("CreateTasksBatch: creating",
        zap.Uint("user_id", uid),
        zap.String("name", req.Name),
        zap.Int("count", len(dates)),
        zap.String("repeat_type", rtype),
        zap.Time("start_date_raw", req.StartDate),
        zap.String("weekly_days", func() string { b, _ := json.Marshal(req.WeeklyDays); return string(b) }()),
    )
    var created []models.Task
    err := db.DB().Transaction(func(tx *gorm.DB) error {
        // 构造任务切片并批量写入
        batch := make([]models.Task, 0, len(dates))
        for _, d := range dates {
            t := models.Task{
                UserID: uid,
                Name: req.Name,
                Description: req.Description,
                Category: req.Category,
                Score: req.Score,
                PlanMinutes: req.PlanMinutes,
                StartDate: d,
                EndDate: nil,
                Status: 0,
            }
            if strings.ToLower(strings.TrimSpace(req.ScoreMode)) == "custom" {
                t.ScoreMode = "custom"
                t.Score = 0
            } else {
                t.ScoreMode = "fixed"
            }
            batch = append(batch, t)
        }
        if err := tx.Create(&batch).Error; err != nil { return err }
        // 同一系列：用首个任务ID作为 series_id
        if len(batch) > 0 {
            sid := batch[0].ID
            ids := make([]uint, 0, len(batch))
            for _, t := range batch { ids = append(ids, t.ID) }
            if err := tx.Model(&models.Task{}).Where("id IN ?", ids).Update("series_id", strconv.Itoa(int(sid))).Error; err != nil {
                return err
            }
            // 更新返回切片中的 SeriesID
            for i := range batch { batch[i].SeriesID = &sid }
        }
        // 记录历史
        for _, t := range batch {
            snap, _ := json.Marshal(t)
            _ = tx.Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "create", SnapshotJSON: string(snap)}).Error
        }
        created = batch
        return nil
    })
    if err != nil {
        zap.L().Error("CreateTasksBatch: db failed", zap.Error(err))
        common.Error(c, 50010, "批量创建失败")
        return
    }
    common.Ok(c, gin.H{"items": created, "count": len(created)})
}
