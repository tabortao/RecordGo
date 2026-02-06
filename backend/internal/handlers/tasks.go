package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"recordgo/internal/common"
	"recordgo/internal/config"
	"recordgo/internal/db"
	"recordgo/internal/models"
)

// 中文注释：任务创建请求结构体（前端字段需与后端一致）
type CreateTaskReq struct {
	UserID           uint       `json:"user_id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	Priority         string     `json:"priority"` // 中文注释：占位字段，实际存入 Category/Remark
	Category         string     `json:"category"`
	Score            int        `json:"score"`
	ScoreMode        string     `json:"score_mode"`
	CustomScoreMax   int        `json:"custom_score_max"`
	DailyMaxCheckins int        `json:"daily_max_checkins"`
	PlanMinutes      int        `json:"plan_minutes"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	ImageJSON        string     `json:"image_json"` // 中文注释：任务图片 JSON 数组字符串
	RepeatType       string     `json:"repeat_type"`
	WeeklyDays       []int      `json:"weekly_days"`
}

// 中文注释：编辑请求结构体（允许全部可变更字段）
type UpdateTaskReq struct {
	Name             *string     `json:"name"`
	Description      *string     `json:"description"`
	Category         *string     `json:"category"`
	Score            *int        `json:"score"`
	ScoreMode        *string     `json:"score_mode"`
	CustomScoreMax   *int        `json:"custom_score_max"`
	DailyMaxCheckins *int        `json:"daily_max_checkins"`
	PlanMinutes      *int        `json:"plan_minutes"`
	ActualMinutes    *int        `json:"actual_minutes"`
	StartDate        *time.Time  `json:"start_date"`
	EndDate          **time.Time `json:"end_date"`
	Remark           *string     `json:"remark"`
	ImageJSON        *string     `json:"image_json"`
	RepeatType       *string     `json:"repeat_type"`
	WeeklyDays       []int       `json:"weekly_days"`
}

// 中文注释：状态变更请求结构体
type UpdateStatusReq struct {
	Status int `json:"status"` // 0待完成 1进行中 2已完成
	// 中文注释：番茄钟完成放行标记（仅用于在 view_only 下允许标记为“已完成”）
	AllowByTomato bool `json:"allow_by_tomato"`
	CustomCoins   *int `json:"custom_coins"`
	Date          string `json:"date"`
}

// 中文注释：番茄钟完成上报请求结构体
type TomatoCompleteReq struct {
	Minutes int `json:"minutes"` // 中文注释：一次番茄钟的工作分钟数（默认25）
}

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	// 中文注释：子账号权限校验——需要具备 tasks.create 权限；家长默认放行
	if !hasPermission(c, "tasks", "create") {
		deny(c, "无权限创建任务")
		return
	}
	var req CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("CreateTask: bind failed", zap.Error(err))
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.Name == "" || req.PlanMinutes <= 0 {
		zap.L().Warn("CreateTask: invalid fields", zap.String("name", req.Name), zap.Int("plan_minutes", req.PlanMinutes))
		common.Error(c, 40002, "任务名称与计划时长必填且合法")
		return
	}
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录或令牌无效")
		return
	}
	// 中文注释：若为子账号，强制归属到父账号（实现数据共享）
	if cl.ParentID != nil {
		req.UserID = *cl.ParentID
	} else if req.UserID == 0 {
		req.UserID = cl.UserID
	}
	if !canAccessUser(c, req.UserID) {
		deny(c, "无权限为该用户创建任务")
		return
	}
	zap.L().Info("CreateTask: creating",
		zap.Uint("user_id", req.UserID),
		zap.String("name", req.Name),
		zap.Int("score", req.Score),
		zap.Int("plan_minutes", req.PlanMinutes),
		zap.String("category", req.Category),
		zap.Time("start_date", req.StartDate),
	)
	s := time.Date(req.StartDate.Year(), req.StartDate.Month(), req.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	var e *time.Time
	if req.EndDate != nil {
		ez := time.Date(req.EndDate.Year(), req.EndDate.Month(), req.EndDate.Day(), 0, 0, 0, 0, time.UTC)
		e = &ez
	}
	if e != nil && e.Before(s) {
		common.Error(c, 40003, "截止日期不能早于开始日期")
		return
	}
	maxCheckins := req.DailyMaxCheckins
	if maxCheckins < 1 {
		maxCheckins = 1
	}
	maxCustom := req.CustomScoreMax
	if maxCustom <= 0 {
		maxCustom = 5
	}
	if maxCustom > 10 {
		maxCustom = 10
	}
	t := models.Task{
		UserID:           req.UserID,
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		Score:            req.Score,
		CustomScoreMax:   maxCustom,
		DailyMaxCheckins: maxCheckins,
		PlanMinutes:      req.PlanMinutes,
		StartDate:        s,
		EndDate:          e,
		ImageJSON:        req.ImageJSON,
		Status:           0,
	}
	if strings.ToLower(strings.TrimSpace(req.ScoreMode)) == "custom" {
		t.ScoreMode = "custom"
		t.Score = 0
	} else {
		t.ScoreMode = "fixed"
	}
	// 规范 repeat 字段
	switch strings.ToLower(strings.TrimSpace(req.RepeatType)) {
	case "daily":
		t.Repeat = "daily"
	case "weekdays":
		t.Repeat = "weekdays"
	case "weekly":
		t.Repeat = "weekly"
	case "monthly":
		t.Repeat = "monthly"
	default:
		t.Repeat = "none"
	}
	if len(req.WeeklyDays) > 0 {
		b, _ := json.Marshal(req.WeeklyDays)
		t.RepeatDaysJSON = string(b)
	}
	zap.L().Info(
		"CreateTask: normalized",
		zap.Uint("user_id", req.UserID),
		zap.String("repeat", t.Repeat),
		zap.String("start_raw", req.StartDate.Format(time.RFC3339)),
		zap.String("start_norm", t.StartDate.Format(time.RFC3339)),
		zap.String("end_raw", func() string {
			if req.EndDate != nil {
				return req.EndDate.Format(time.RFC3339)
			}
			return ""
		}()),
		zap.String("end_norm", func() string {
			if t.EndDate != nil {
				return t.EndDate.Format(time.RFC3339)
			}
			return ""
		}()),
	)
	if err := db.DB().Create(&t).Error; err != nil {
		zap.L().Error("CreateTask: db create failed", zap.Error(err))
		common.Error(c, 50001, "创建任务失败")
		return
	}
	zap.L().Info("CreateTask: created", zap.Uint("task_id", t.ID))
	// 历史记录：create
	snapshot, _ := json.Marshal(t)
	_ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "create", SnapshotJSON: string(snapshot)}).Error
	common.Ok(c, t)
}

// ListTasks 列表查询（支持状态与分页）
func ListTasks(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录或令牌无效")
		return
	}
	uid := cl.UserID
	if cl.ParentID != nil {
		uid = *cl.ParentID
	}
	if s := strings.TrimSpace(c.Query("user_id")); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			if canAccessUser(c, uint(v)) {
				uid = uint(v)
			} else {
				deny(c, "无权限查看该用户任务")
				return
			}
		}
	}
	var tasks []models.Task
	q := db.DB().Model(&models.Task{}).Where("user_id = ?", uid)

	// 过滤状态
	if s := c.Query("status"); s != "" {
		if si, err := strconv.Atoi(s); err == nil {
			q = q.Where("status = ?", si)
		}
	}
	// 简单分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var total int64
	q.Count(&total)
	if err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&tasks).Error; err != nil {
		common.Error(c, 50002, "查询任务失败")
		return
	}
	for i := range tasks {
		s := strings.TrimSpace(tasks[i].RepeatDaysJSON)
		if s != "" {
			var arr []int
			if jsonErr := json.Unmarshal([]byte(s), &arr); jsonErr == nil {
				tasks[i].WeeklyDays = arr
			}
		}
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
	if !canAccessUser(c, t.UserID) {
		deny(c, "无权限查看该任务")
		return
	}
	s := strings.TrimSpace(t.RepeatDaysJSON)
	if s != "" {
		var arr []int
		if jsonErr := json.Unmarshal([]byte(s), &arr); jsonErr == nil {
			t.WeeklyDays = arr
		}
	}
	common.Ok(c, t)
}

// UpdateTask 编辑任务（记录历史）
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	// 中文注释：子账号权限校验——需要具备 tasks.edit 权限；家长默认放行
	if !hasPermission(c, "tasks", "edit") {
		deny(c, "无权限编辑任务")
		return
	}
	var req UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("UpdateTask: bind failed", zap.Error(err), zap.String("id", id))
		common.Error(c, 40001, "参数错误")
		return
	}
	var t models.Task
	if err := db.DB().First(&t, id).Error; err != nil {
		zap.L().Warn("UpdateTask: not found", zap.Error(err), zap.String("id", id))
		common.Error(c, 40401, "任务不存在")
		return
	}
	if !canAccessUser(c, t.UserID) {
		deny(c, "无权限编辑该任务")
		return
	}
	// 修改前快照
	before, _ := json.Marshal(t)

	// 应用变更
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.Category != nil {
		t.Category = *req.Category
	}
	if req.Score != nil {
		t.Score = *req.Score
	}
	if req.ScoreMode != nil {
		if strings.ToLower(strings.TrimSpace(*req.ScoreMode)) == "custom" {
			t.ScoreMode = "custom"
		} else {
			t.ScoreMode = "fixed"
		}
	}
	if req.CustomScoreMax != nil {
		v := *req.CustomScoreMax
		if v <= 0 {
			v = 5
		}
		if v > 10 {
			v = 10
		}
		t.CustomScoreMax = v
	}
	if t.ScoreMode == "custom" {
		t.Score = 0
	}
	if req.DailyMaxCheckins != nil {
		v := *req.DailyMaxCheckins
		if v < 1 {
			v = 1
		}
		t.DailyMaxCheckins = v
	}
	if req.PlanMinutes != nil {
		t.PlanMinutes = *req.PlanMinutes
	}
	if req.ActualMinutes != nil {
		t.ActualMinutes = *req.ActualMinutes
	}
	if req.StartDate != nil {
		t.StartDate = time.Date(req.StartDate.Year(), req.StartDate.Month(), req.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	}
	if req.EndDate != nil {
		if *req.EndDate == nil {
			t.EndDate = nil
		} else {
			ez := time.Date((*req.EndDate).Year(), (*req.EndDate).Month(), (*req.EndDate).Day(), 0, 0, 0, 0, time.UTC)
			t.EndDate = &ez
		}
	}
	if req.Remark != nil {
		t.Remark = *req.Remark
	}
	if req.ImageJSON != nil {
		t.ImageJSON = *req.ImageJSON
	}
	if req.RepeatType != nil {
		switch strings.ToLower(strings.TrimSpace(*req.RepeatType)) {
		case "daily":
			t.Repeat = "daily"
		case "weekdays":
			t.Repeat = "weekdays"
		case "weekly":
			t.Repeat = "weekly"
		case "monthly":
			t.Repeat = "monthly"
		default:
			t.Repeat = "none"
		}
	}
	if len(req.WeeklyDays) > 0 {
		b, _ := json.Marshal(req.WeeklyDays)
		t.RepeatDaysJSON = string(b)
	}
	if t.EndDate != nil && t.EndDate.Before(t.StartDate) {
		common.Error(c, 40003, "截止日期不能早于开始日期")
		return
	}

	if err := db.DB().Save(&t).Error; err != nil {
		zap.L().Error("UpdateTask: db save failed", zap.Error(err), zap.String("id", id))
		common.Error(c, 50004, "更新任务失败")
		return
	}
	zap.L().Info("UpdateTask: updated", zap.Uint("task_id", t.ID))
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
	if !canAccessUser(c, t.UserID) {
		deny(c, "无权限更改该任务状态")
		return
	}
	// 快照
	before, _ := json.Marshal(t)
	prev := t.Status
	next := req.Status
	var completedAt *time.Time
	if strings.TrimSpace(req.Date) != "" {
		if d, err := time.Parse("2006-01-02", strings.TrimSpace(req.Date)); err == nil {
			dd := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
			completedAt = &dd
		}
	}
	if completedAt == nil {
		now := time.Now().UTC()
		dd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		completedAt = &dd
	}
	// 中文注释：权限拆分——完成/状态变更用 tasks.status；撤销完成用 tasks.undo；番茄钟完成可通过 allow_by_tomato 放行
	if prev == 2 && next != 2 {
		if !hasPermission(c, "tasks", "undo") {
			deny(c, "无权限撤销任务")
			return
		}
	} else if next == 2 {
		if !hasPermission(c, "tasks", "status") {
			if !(req.AllowByTomato && next == 2) {
				deny(c, "无权限更改任务状态")
				return
			}
		}
	} else {
		if !hasPermission(c, "tasks", "status") {
			deny(c, "无权限更改任务状态")
			return
		}
	}
	var targetCoins int64
	if err := db.DB().Transaction(func(tx *gorm.DB) error {
		var owner models.User
		if err := tx.First(&owner, t.UserID).Error; err != nil {
			owner = models.User{ID: t.UserID, Username: fmt.Sprintf("user%d", t.UserID), Role: "user", Coins: 0, Nickname: "测试用户"}
			if ierr := tx.Create(&owner).Error; ierr != nil {
				return ierr
			}
		}
		target := owner
		if cfg2, _ := config.Load(); cfg2 != nil && cfg2.ParentCoinsSync && owner.ParentID != nil {
			var parent models.User
			if err := tx.First(&parent, *owner.ParentID).Error; err == nil {
				target = parent
			}
		}
		if prev == 2 && next != 2 {
			revert := t.Score
			if strings.ToLower(strings.TrimSpace(t.ScoreMode)) == "custom" || t.CompletedScore != 0 {
				revert = t.CompletedScore
			}
			target.Coins -= int64(revert)
			t.CompletedScore = 0
			t.CompletedAt = nil
		} else if prev != 2 && next == 2 {
			award := t.Score
			if strings.ToLower(strings.TrimSpace(t.ScoreMode)) == "custom" {
				if req.CustomCoins == nil {
					return newBizError(40003, "缺少自定义金币")
				}
				max := t.CustomScoreMax
				if max <= 0 {
					max = 5
				}
				if max > 10 {
					max = 10
				}
				if *req.CustomCoins < 1 || *req.CustomCoins > max {
					return newBizError(40003, "自定义金币超出上限")
				}
				award = *req.CustomCoins
			}
			target.Coins += int64(award)
			t.CompletedScore = award
			t.CompletedAt = completedAt
		}
		if err := tx.Save(&target).Error; err != nil {
			return err
		}
		t.Status = next
		if err := tx.Save(&t).Error; err != nil {
			return err
		}
		targetCoins = target.Coins
		return nil
	}); err != nil {
		common.Error(c, 50005, "更新状态失败")
		return
	}
	_ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "status", SnapshotJSON: string(before)}).Error
	// 中文注释：返回 user_coins 为金币实际归属用户（父或本人）的最新值，便于前端统一展示
	common.Ok(c, gin.H{"task": t, "user_coins": targetCoins})
}

// DeleteTask 软删除（进入回收站）
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	// 中文注释：子账号权限校验——需要具备 tasks.delete 权限；家长默认放行
	if !hasPermission(c, "tasks", "delete") {
		deny(c, "无权限删除任务")
		return
	}
	var t models.Task
	if err := db.DB().First(&t, id).Error; err != nil {
		common.Error(c, 40401, "任务不存在")
		return
	}
	if !canAccessUser(c, t.UserID) {
		deny(c, "无权限删除该任务")
		return
	}
	if err := db.DB().Delete(&t).Error; err != nil {
		common.Error(c, 50006, "删除失败")
		return
	}
	// 历史记录：delete
	snapshot, _ := json.Marshal(t)
	_ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "delete", SnapshotJSON: string(snapshot)}).Error
	// 中文注释：删除任务时同步删除关联图片文件
	if strings.TrimSpace(t.ImageJSON) != "" {
		var paths []string
		_ = json.Unmarshal([]byte(t.ImageJSON), &paths)
		if len(paths) > 0 {
			root := os.Getenv("STORAGE_ROOT")
			if strings.TrimSpace(root) == "" {
				root = "storage"
			}
			for _, rel := range paths {
				full := filepath.Join(root, filepath.FromSlash(rel))
				if rmErr := os.Remove(full); rmErr == nil {
					zap.L().Info("DeleteTask: removed image file", zap.String("path", full))
				} else {
					zap.L().Warn("DeleteTask: remove image file failed", zap.String("path", full), zap.Error(rmErr))
				}
			}
		}
	}
	// 中文注释：额外清理“任务批注”上传的图片（不在 image_json 中）——按文件名前缀 task_任务ID_*
	// 目录：storage/uploads/images/task_images/{用户id}
	{
		root := os.Getenv("STORAGE_ROOT")
		if strings.TrimSpace(root) == "" {
			root = "storage"
		}
		dir := filepath.Join(root, "uploads", "images", "task_images", fmt.Sprintf("%d", t.UserID))
		entries, readErr := os.ReadDir(dir)
		if readErr == nil {
			prefix := fmt.Sprintf("task_%d_", t.ID)
			for _, ent := range entries {
				if ent.IsDir() {
					continue
				}
				name := ent.Name()
				if strings.HasPrefix(name, prefix) {
					full := filepath.Join(dir, name)
					if rmErr := os.Remove(full); rmErr == nil {
						zap.L().Info("DeleteTask: removed note image", zap.Uint("task_id", t.ID), zap.String("path", full))
					} else {
						zap.L().Warn("DeleteTask: remove note image failed", zap.Uint("task_id", t.ID), zap.String("path", full), zap.Error(rmErr))
					}
				}
			}
		} else {
			zap.L().Debug("DeleteTask: skip note image cleanup (dir missing)", zap.String("dir", dir), zap.Error(readErr))
		}
	}
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
	// 中文注释：删除前查询图片路径并尝试删除文件
	var tasks []models.Task
	_ = db.DB().Where("id IN ?", ids).Find(&tasks).Error
	kept := make([]uint, 0, len(ids))
	for _, t := range tasks {
		if canAccessUser(c, t.UserID) {
			kept = append(kept, t.ID)
		}
	}
	if len(kept) == 0 {
		common.Error(c, 40301, "无可删除任务")
		return
	}
	if err := db.DB().Where("id IN ?", kept).Delete(&models.Task{}).Error; err != nil {
		common.Error(c, 50007, "批量删除失败")
		return
	}
	if len(tasks) > 0 {
		root := os.Getenv("STORAGE_ROOT")
		if strings.TrimSpace(root) == "" {
			root = "storage"
		}
		for _, t := range tasks {
			if !containsID(kept, t.ID) {
				continue
			}
			if strings.TrimSpace(t.ImageJSON) == "" {
				continue
			}
			var paths []string
			_ = json.Unmarshal([]byte(t.ImageJSON), &paths)
			for _, rel := range paths {
				full := filepath.Join(root, filepath.FromSlash(rel))
				if rmErr := os.Remove(full); rmErr == nil {
					zap.L().Info("BatchDelete: removed image file", zap.Uint("task_id", t.ID), zap.String("path", full))
				} else {
					zap.L().Warn("BatchDelete: remove image file failed", zap.Uint("task_id", t.ID), zap.String("path", full), zap.Error(rmErr))
				}
			}
			// 中文注释：额外清理“任务批注”上传的图片（不在 image_json 中），按文件名前缀 task_任务ID_*，目录按用户ID
			dir := filepath.Join(root, "uploads", "images", "task_images", fmt.Sprintf("%d", t.UserID))
			entries, readErr := os.ReadDir(dir)
			if readErr == nil {
				prefix := fmt.Sprintf("task_%d_", t.ID)
				for _, ent := range entries {
					if ent.IsDir() {
						continue
					}
					name := ent.Name()
					if strings.HasPrefix(name, prefix) {
						full := filepath.Join(dir, name)
						if rmErr := os.Remove(full); rmErr == nil {
							zap.L().Info("BatchDelete: removed note image", zap.Uint("task_id", t.ID), zap.String("path", full))
						} else {
							zap.L().Warn("BatchDelete: remove note image failed", zap.Uint("task_id", t.ID), zap.String("path", full), zap.Error(rmErr))
						}
					}
				}
			} else {
				zap.L().Debug("BatchDelete: skip note image cleanup (dir missing)", zap.Uint("task_id", t.ID), zap.String("dir", dir), zap.Error(readErr))
			}
		}
	}
	common.Ok(c, gin.H{"deleted": len(ids)})
}

// ListRecycleBin 回收站任务
func ListRecycleBin(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录或令牌无效")
		return
	}
	uid := cl.UserID
	if cl.ParentID != nil {
		uid = *cl.ParentID
	}
	if s := strings.TrimSpace(c.Query("user_id")); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			if canAccessUser(c, uint(v)) {
				uid = uint(v)
			} else {
				deny(c, "无权限查看该用户回收站")
				return
			}
		}
	}
	var tasks []models.Task
	if err := db.DB().Unscoped().Where("deleted_at IS NOT NULL AND user_id = ?", uid).Order("deleted_at DESC").Find(&tasks).Error; err != nil {
		common.Error(c, 50008, "查询回收站失败")
		return
	}
	for i := range tasks {
		s := strings.TrimSpace(tasks[i].RepeatDaysJSON)
		if s != "" {
			var arr []int
			_ = json.Unmarshal([]byte(s), &arr)
			tasks[i].WeeklyDays = arr
		}
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
	restored := 0
	for _, id := range ids {
		var t models.Task
		if err := db.DB().Unscoped().First(&t, id).Error; err == nil {
			if !canAccessUser(c, t.UserID) {
				continue
			}
			t.DeletedAt = gorm.DeletedAt{}
			if err := db.DB().Unscoped().Save(&t).Error; err == nil {
				snapshot, _ := json.Marshal(t)
				_ = db.DB().Create(&models.TaskHistory{TaskID: t.ID, ChangeType: "restore", SnapshotJSON: string(snapshot)}).Error
				restored++
			}
		}
	}
	common.Ok(c, gin.H{"restored": restored})
}

// CompleteTomato 番茄钟完成上报（累加任务番茄数与实际耗时）
func CompleteTomato(c *gin.Context) {
	id := c.Param("id")
	var req TomatoCompleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.Minutes <= 0 {
		req.Minutes = 25
	}
	var t models.Task
	if err := db.DB().First(&t, id).Error; err != nil {
		common.Error(c, 40401, "任务不存在")
		return
	}
	if !canAccessUser(c, t.UserID) {
		deny(c, "无权限操作该任务")
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
			if seg != "" {
				out = append(out, seg)
			}
			start = i + 1
		}
	}
	return out
}

func containsID(arr []uint, id uint) bool {
	for _, v := range arr {
		if v == id {
			return true
		}
	}
	return false
}
