package handlers

import (
	"time"

	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
)

type ScoreCreateReq struct {
	Subject   string   `json:"subject"`
	ExamName  string   `json:"exam_name"`
	ExamType  string   `json:"exam_type"`
	Grade     string   `json:"grade"`
	Term      string   `json:"term"`
	Topic     string   `json:"topic"`
	Difficulty string  `json:"difficulty"`
	ExamDate  string   `json:"exam_date"`
	Score     float64  `json:"score"`
	FullScore float64  `json:"full_score"`
	ClassRank *int     `json:"class_rank"`
	RankType  string   `json:"rank_type"`
	GradeRank *int     `json:"grade_rank"`
	ClassAvg  *float64 `json:"class_avg"`
	ClassHighest *float64 `json:"class_highest"`
	PhotoURL  string   `json:"photo_url"`
	Remark    string   `json:"remark"`
}

type ScoreUpdateReq struct {
	Subject   *string  `json:"subject"`
	ExamName  *string  `json:"exam_name"`
	ExamType  *string  `json:"exam_type"`
	Grade     *string  `json:"grade"`
	Term      *string  `json:"term"`
	Topic     *string  `json:"topic"`
	Difficulty *string `json:"difficulty"`
	ExamDate  *string  `json:"exam_date"`
	Score     *float64 `json:"score"`
	FullScore *float64 `json:"full_score"`
	ClassRank *int     `json:"class_rank"`
	RankType  *string  `json:"rank_type"`
	GradeRank *int     `json:"grade_rank"`
	ClassAvg  *float64 `json:"class_avg"`
	ClassHighest *float64 `json:"class_highest"`
	PhotoURL  *string  `json:"photo_url"`
	Remark    *string  `json:"remark"`
}

func ListScores(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	uid := cl.UserID
	if cl.ParentID != nil {
		uid = *cl.ParentID
	}
	subject := c.Query("subject")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	q := db.DB().Model(&models.ScoreRecord{}).Where("user_id = ?", uid)
	if subject != "" {
		q = q.Where("subject = ?", subject)
	}
	if startDate != "" {
		d, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
		if err != nil {
			common.Error(c, 40001, "开始日期格式错误")
			return
		}
		q = q.Where("exam_date >= ?", d)
	}
	if endDate != "" {
		d, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
		if err != nil {
			common.Error(c, 40001, "结束日期格式错误")
			return
		}
		q = q.Where("exam_date <= ?", d)
	}
	var items []models.ScoreRecord
	if err := q.Order("exam_date DESC, created_at DESC").Find(&items).Error; err != nil {
		common.Error(c, 50001, "查询失败")
		return
	}
	common.Ok(c, gin.H{"items": items})
}

func CreateScore(c *gin.Context) {
	if !hasPermission(c, "scores", "create") {
		deny(c, "无权限创建成绩")
		return
	}
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	var req ScoreCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.Subject == "" || req.ExamName == "" || req.ExamDate == "" {
		common.Error(c, 40002, "学科/考试名称/日期必填")
		return
	}
	if req.FullScore <= 0 {
		common.Error(c, 40003, "满分必须大于0")
		return
	}
	if req.Score < 0 || req.Score > req.FullScore {
		common.Error(c, 40004, "分数范围错误")
		return
	}
	d, err := time.ParseInLocation("2006-01-02", req.ExamDate, time.Local)
	if err != nil {
		common.Error(c, 40005, "日期格式错误")
		return
	}
	uid := cl.UserID
	if cl.ParentID != nil {
		uid = *cl.ParentID
	}
	item := models.ScoreRecord{
		UserID:    uid,
		Subject:   req.Subject,
		ExamName:  req.ExamName,
		ExamType:  req.ExamType,
		Grade:     req.Grade,
		Term:      req.Term,
		Topic:     req.Topic,
		Difficulty: req.Difficulty,
		ExamDate:  d,
		Score:     req.Score,
		FullScore: req.FullScore,
		ClassRank: req.ClassRank,
		RankType:  req.RankType,
		GradeRank: req.GradeRank,
		ClassAvg:  req.ClassAvg,
		ClassHighest: req.ClassHighest,
		PhotoURL:  req.PhotoURL,
		Remark:    req.Remark,
	}
	if err := db.DB().Create(&item).Error; err != nil {
		common.Error(c, 50002, "创建失败")
		return
	}
	common.Ok(c, item)
}

func GetScore(c *gin.Context) {
	id := c.Param("id")
	var item models.ScoreRecord
	if err := db.DB().First(&item, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if !canAccessUser(c, item.UserID) {
		common.Error(c, 40301, "无权查看")
		return
	}
	common.Ok(c, item)
}

func UpdateScore(c *gin.Context) {
	if !hasPermission(c, "scores", "edit") {
		deny(c, "无权限编辑成绩")
		return
	}
	id := c.Param("id")
	var req ScoreUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	var item models.ScoreRecord
	if err := db.DB().First(&item, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if !canAccessUser(c, item.UserID) {
		common.Error(c, 40301, "无权操作")
		return
	}
	if req.Subject != nil {
		item.Subject = *req.Subject
	}
	if req.ExamName != nil {
		item.ExamName = *req.ExamName
	}
	if req.ExamType != nil {
		item.ExamType = *req.ExamType
	}
	if req.Grade != nil {
		item.Grade = *req.Grade
	}
	if req.Term != nil {
		item.Term = *req.Term
	}
	if req.Topic != nil {
		item.Topic = *req.Topic
	}
	if req.Difficulty != nil {
		item.Difficulty = *req.Difficulty
	}
	if req.ExamDate != nil {
		d, err := time.ParseInLocation("2006-01-02", *req.ExamDate, time.Local)
		if err != nil {
			common.Error(c, 40005, "日期格式错误")
			return
		}
		item.ExamDate = d
	}
	if req.Score != nil {
		if *req.Score < 0 {
			common.Error(c, 40004, "分数范围错误")
			return
		}
		item.Score = *req.Score
	}
	if req.FullScore != nil {
		if *req.FullScore <= 0 {
			common.Error(c, 40003, "满分必须大于0")
			return
		}
		item.FullScore = *req.FullScore
	}
	if item.Score > item.FullScore {
		common.Error(c, 40004, "分数范围错误")
		return
	}
	if req.ClassRank != nil {
		item.ClassRank = req.ClassRank
	}
	if req.RankType != nil {
		item.RankType = *req.RankType
	}
	if req.GradeRank != nil {
		item.GradeRank = req.GradeRank
	}
	if req.ClassAvg != nil {
		item.ClassAvg = req.ClassAvg
	}
	if req.ClassHighest != nil {
		item.ClassHighest = req.ClassHighest
	}
	if req.PhotoURL != nil {
		item.PhotoURL = *req.PhotoURL
	}
	if req.Remark != nil {
		item.Remark = *req.Remark
	}
	if err := db.DB().Save(&item).Error; err != nil {
		common.Error(c, 50004, "更新失败")
		return
	}
	common.Ok(c, item)
}

func DeleteScore(c *gin.Context) {
	if !hasPermission(c, "scores", "delete") {
		deny(c, "无权限删除成绩")
		return
	}
	id := c.Param("id")
	var item models.ScoreRecord
	if err := db.DB().First(&item, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if !canAccessUser(c, item.UserID) {
		common.Error(c, 40301, "无权操作")
		return
	}
	if err := db.DB().Delete(&item).Error; err != nil {
		common.Error(c, 50003, "删除失败")
		return
	}
	common.Ok(c, gin.H{"id": item.ID})
}
