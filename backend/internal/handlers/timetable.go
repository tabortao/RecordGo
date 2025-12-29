package handlers

import (
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 默认课程列表 (马卡龙配色)
var defaultCourses = []models.CourseDict{
	{Name: "语文", Color: "#FF9AA2"},   // 柔和红
	{Name: "数学", Color: "#C7CEEA"},   // 浅紫蓝
	{Name: "英语", Color: "#FFF5BA"},   // 奶油黄
	{Name: "物理", Color: "#B5EAD7"},   // 薄荷绿
	{Name: "化学", Color: "#E0BBE4"},   // 薰衣草紫
	{Name: "生物", Color: "#A2E4B8"},   // 清新绿
	{Name: "政治", Color: "#FFDAC1"},   // 蜜桃色
	{Name: "历史", Color: "#FFB7B2"},   // 浅鲑鱼红
	{Name: "地理", Color: "#95F9E3"},   // 浅青色
	{Name: "体育", Color: "#D4F0F0"},   // 淡蓝
	{Name: "音乐", Color: "#FFC8A2"},   // 浅橙
	{Name: "美术", Color: "#D7B3E4"},   // 浅锦葵紫
	{Name: "信息技术", Color: "#D4D8E7"}, // 凉爽灰
	{Name: "科学", Color: "#E6F9AF"},   // 嫩黄绿
	{Name: "自习", Color: "#F2F2F2"},   // 浅灰
}

// ensureDefaultCourses 确保默认课程存在，并更新颜色
func ensureDefaultCourses() {
	var count int64
	db.DB().Model(&models.CourseDict{}).Where("user_id IS NULL").Count(&count)
	if count == 0 {
		for _, course := range defaultCourses {
			db.DB().Create(&course)
		}
	} else {
		// 简单的检查：如果名字不对，可能是旧数据，这里暂不自动修正，以免覆盖用户可能的修改（虽然user_id=null不应修改）
		// 但为了满足需求“信息技术”，我们可以尝试更新一下
		db.DB().Model(&models.CourseDict{}).Where("user_id IS NULL AND name = ?", "信息").Update("name", "信息技术")

		// 强制更新系统课程颜色为最新的浅色系
		for _, course := range defaultCourses {
			db.DB().Model(&models.CourseDict{}).
				Where("user_id IS NULL AND name = ?", course.Name).
				Update("color", course.Color)
		}
	}
}

// getEffectiveTimetableUserID 获取实际的课表归属用户ID（子账号继承父账号）
func getEffectiveTimetableUserID(userID uint) uint {
	var u models.User
	if err := db.DB().Select("id, parent_id").First(&u, userID).Error; err != nil {
		return userID
	}
	if u.ParentID != nil {
		return *u.ParentID
	}
	return userID
}

// GetTimetableConfig 获取课表配置
func GetTimetableConfig(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		uid, err := strconv.Atoi(userIDStr)
		if err != nil {
			common.Error(c, 400, "无效的用户ID")
			return
		}
		userID = uint(uid)
	} else {
		cl := extractClaims(c)
		if cl == nil {
			common.Error(c, 401, "未登录")
			return
		}
		userID = cl.UserID
	}

	if !canAccessUser(c, userID) {
		deny(c, "无权限查看该用户配置")
		return
	}

	// 获取实际归属用户ID（子账号 -> 父账号）
	effectiveUserID := getEffectiveTimetableUserID(userID)

	var config models.TimetableConfig
	err := db.DB().Where("user_id = ?", effectiveUserID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，创建默认配置
			config = models.TimetableConfig{
				UserID:          effectiveUserID,
				ShowSaturday:    false,
				ShowSunday:      false,
				CurrentGrade:    "一年级",
				CurrentSemester: "上学期",
			}
			db.DB().Create(&config)
		} else {
			zap.L().Error("GetTimetableConfig error", zap.Error(err))
			common.Error(c, 500, "获取配置失败")
			return
		}
	}

	common.Ok(c, config)
}

// UpdateTimetableConfig 更新课表配置
func UpdateTimetableConfig(c *gin.Context) {
	var req models.TimetableConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}

	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未登录")
		return
	}

	// 如果请求中没有 UserID，使用当前登录用户
	targetUserID := req.UserID
	if targetUserID == 0 {
		targetUserID = cl.UserID
		req.UserID = targetUserID
	}

	if !canAccessUser(c, targetUserID) {
		deny(c, "无权限修改该用户配置")
		return
	}

	// 获取实际归属用户ID
	effectiveUserID := getEffectiveTimetableUserID(targetUserID)
	req.UserID = effectiveUserID // 确保保存到正确的 ID

	// 更新或创建
	var config models.TimetableConfig
	err := db.DB().Where("user_id = ?", effectiveUserID).First(&config).Error
	switch err {
	case gorm.ErrRecordNotFound:
		config = req
		if createErr := db.DB().Create(&config).Error; createErr != nil {
			common.Error(c, 500, "创建配置失败")
			return
		}
	case nil:
		// 更新字段
		config.ShowSaturday = req.ShowSaturday
		config.ShowSunday = req.ShowSunday
		config.CurrentGrade = req.CurrentGrade
		config.CurrentSemester = req.CurrentSemester
		config.PeriodSettingsJSON = req.PeriodSettingsJSON
		config.BackgroundEmojis = req.BackgroundEmojis
		if err := db.DB().Save(&config).Error; err != nil {
			common.Error(c, 500, "更新配置失败")
			return
		}
	default:
		common.Error(c, 500, "数据库错误")
		return
	}

	common.Ok(c, config)
}

// GetCourses 获取课程列表（系统默认+用户自定义）
func GetCourses(c *gin.Context) {
	ensureDefaultCourses()

	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		uid, err := strconv.Atoi(userIDStr)
		if err != nil {
			common.Error(c, 400, "无效的用户ID")
			return
		}
		userID = uint(uid)
	} else {
		cl := extractClaims(c)
		if cl == nil {
			common.Error(c, 401, "未登录")
			return
		}
		userID = cl.UserID
	}

	if !canAccessUser(c, userID) {
		deny(c, "无权限查看该用户课程")
		return
	}

	// 获取实际归属用户ID
	effectiveUserID := getEffectiveTimetableUserID(userID)

	var courses []models.CourseDict
	// 查询 UserID 为 NULL (系统) 或 UserID = effectiveUserID (自定义)
	if err := db.DB().Where("user_id IS NULL OR user_id = ?", effectiveUserID).Find(&courses).Error; err != nil {
		common.Error(c, 500, "获取课程失败")
		return
	}

	common.Ok(c, courses)
}

// AddCourse 添加自定义课程
func AddCourse(c *gin.Context) {
	var req models.CourseDict
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}

	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未登录")
		return
	}

	// 强制设置为当前操作的目标用户（如果有权限）
	if req.UserID == nil || *req.UserID == 0 {
		req.UserID = &cl.UserID
	}

	if !canAccessUser(c, *req.UserID) {
		deny(c, "无权限为该用户添加课程")
		return
	}

	// 获取实际归属用户ID
	effectiveUserID := getEffectiveTimetableUserID(*req.UserID)
	req.UserID = &effectiveUserID

	// 检查是否已存在同名课程
	var existCount int64
	db.DB().Model(&models.CourseDict{}).
		Where("(user_id = ? OR user_id IS NULL) AND name = ?", effectiveUserID, req.Name).
		Count(&existCount)

	if existCount > 0 {
		common.Error(c, 400, "该课程已存在")
		return
	}

	if err := db.DB().Create(&req).Error; err != nil {
		common.Error(c, 500, "添加课程失败")
		return
	}

	common.Ok(c, req)
}

// UpdateCourse 更新自定义课程
func UpdateCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.Error(c, 400, "无效的课程ID")
		return
	}

	var req models.CourseDict
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}

	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未登录")
		return
	}

	var course models.CourseDict
	if err := db.DB().First(&course, id).Error; err != nil {
		common.Error(c, 404, "课程不存在")
		return
	}

	// 系统课程不可修改
	if course.UserID == nil {
		common.Error(c, 403, "系统内置课程不可修改")
		return
	}

	// 权限检查：注意这里的 course.UserID 已经是数据库中的实际拥有者（父账号ID）
	// 我们需要检查当前操作者是否有权限操作这个“实际拥有者”的数据
	// 但 canAccessUser 检查的是当前用户能否操作 targetID。
	// 这里比较复杂：子账号想修改父账号的课程（如果它是父账号创建的）。
	// 简单处理：如果 course.UserID == cl.UserID (本人)，或者 cl.UserID 是 course.UserID 的子账号 (子改父)，或者 cl.UserID 是 course.UserID 的父账号 (父改子，虽然子不存课程)。
	// 更好的逻辑：检查当前用户是否关联到 course.UserID。
	// getEffectiveTimetableUserID(cl.UserID) == *course.UserID

	effectiveOpUserID := getEffectiveTimetableUserID(cl.UserID)
	if effectiveOpUserID != *course.UserID {
		// 也许是父账号在操作子账号？不，所有数据都在父账号下。
		// 也许是管理员？
		if !canAccessUser(c, *course.UserID) { // 回退到标准检查
			deny(c, "无权限修改该课程")
			return
		}
	}

	// 更新字段
	course.Name = req.Name
	course.Color = req.Color
	if err := db.DB().Save(&course).Error; err != nil {
		common.Error(c, 500, "更新课程失败")
		return
	}

	common.Ok(c, course)
}

// DeleteCourse 删除自定义课程
func DeleteCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.Error(c, 400, "无效的课程ID")
		return
	}

	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未登录")
		return
	}

	var course models.CourseDict
	if err := db.DB().First(&course, id).Error; err != nil {
		common.Error(c, 404, "课程不存在")
		return
	}

	// 系统课程不可删除
	if course.UserID == nil {
		common.Error(c, 403, "系统内置课程不可删除")
		return
	}

	// 权限检查 (同 UpdateCourse)
	effectiveOpUserID := getEffectiveTimetableUserID(cl.UserID)
	if effectiveOpUserID != *course.UserID {
		if !canAccessUser(c, *course.UserID) {
			deny(c, "无权限删除该课程")
			return
		}
	}

	// 检查是否被课表引用
	var count int64
	db.DB().Model(&models.Timetable{}).Where("course_id = ?", id).Count(&count)
	if count > 0 {
		common.Error(c, 400, "该课程已被课表使用，无法删除，请先移除相关课表安排")
		return
	}

	if err := db.DB().Delete(&course).Error; err != nil {
		common.Error(c, 500, "删除课程失败")
		return
	}

	common.Ok(c, "删除成功")
}

// GetTimetable 获取课表
func GetTimetable(c *gin.Context) {
	userIDStr := c.Query("user_id")
	grade := c.Query("grade")
	semester := c.Query("semester")

	if grade == "" || semester == "" {
		common.Error(c, 400, "缺少年级或学期参数")
		return
	}

	var userID uint
	if userIDStr != "" {
		uid, err := strconv.Atoi(userIDStr)
		if err != nil {
			common.Error(c, 400, "无效的用户ID")
			return
		}
		userID = uint(uid)
	} else {
		cl := extractClaims(c)
		if cl == nil {
			common.Error(c, 401, "未登录")
			return
		}
		userID = cl.UserID
	}

	if !canAccessUser(c, userID) {
		deny(c, "无权限查看该用户课表")
		return
	}

	// 获取实际归属用户ID
	effectiveUserID := getEffectiveTimetableUserID(userID)

	var timetable []models.Timetable
	if err := db.DB().Preload("Course").Where("user_id = ? AND grade = ? AND semester = ?", effectiveUserID, grade, semester).Find(&timetable).Error; err != nil {
		common.Error(c, 500, "获取课表失败")
		return
	}

	common.Ok(c, timetable)
}

// SaveTimetable 保存课表（批量）
type SaveTimetableReq struct {
	UserID   uint               `json:"user_id"`
	Grade    string             `json:"grade"`
	Semester string             `json:"semester"`
	Items    []models.Timetable `json:"items"`
}

func SaveTimetable(c *gin.Context) {
	var req SaveTimetableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}

	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未登录")
		return
	}

	if req.UserID == 0 {
		req.UserID = cl.UserID
	}

	if !canAccessUser(c, req.UserID) {
		deny(c, "无权限修改该用户课表")
		return
	}

	// 获取实际归属用户ID
	effectiveUserID := getEffectiveTimetableUserID(req.UserID)
	req.UserID = effectiveUserID // 修正为实际ID

	tx := db.DB().Begin()

	// 1. 删除该用户该年级该学期的原有数据
	if err := tx.Where("user_id = ? AND grade = ? AND semester = ?", req.UserID, req.Grade, req.Semester).Delete(&models.Timetable{}).Error; err != nil {
		tx.Rollback()
		common.Error(c, 500, "清理旧数据失败")
		return
	}

	// 2. 插入新数据
	if len(req.Items) > 0 {
		for i := range req.Items {
			req.Items[i].UserID = req.UserID
			req.Items[i].Grade = req.Grade
			req.Items[i].Semester = req.Semester
			// ID 应该置零以创建新记录
			req.Items[i].ID = 0
		}
		if err := tx.Create(&req.Items).Error; err != nil {
			tx.Rollback()
			common.Error(c, 500, "保存数据失败")
			return
		}
	}

	tx.Commit()
	common.Ok(c, "保存成功")
}
