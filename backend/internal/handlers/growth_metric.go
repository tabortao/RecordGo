package handlers

import (
	"strconv"
	"time"

	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
)

var growthMetricTypes = map[string]bool{
	"height": true,
	"weight": true,
	"vision": true,
	"foot":   true,
}

type GrowthMetricCreateReq struct {
	MetricType string  `json:"metric_type"`
	RecordDate string  `json:"record_date"`
	Value      float64 `json:"value"`
	LeftValue  float64 `json:"left_value"`
	RightValue float64 `json:"right_value"`
}

func ListGrowthMetricRecords(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	metricType := c.Query("type")
	if metricType != "" && !growthMetricTypes[metricType] {
		common.Error(c, 40001, "指标类型错误")
		return
	}
	q := db.DB().Model(&models.GrowthMetricRecord{}).Where("user_id = ?", cl.UserID)
	if metricType != "" {
		q = q.Where("metric_type = ?", metricType)
	}
	var total int64
	q.Count(&total)
	var list []models.GrowthMetricRecord
	if err := q.Order("record_date DESC, created_at DESC").Find(&list).Error; err != nil {
		common.Error(c, 50001, "查询失败")
		return
	}
	common.Ok(c, gin.H{"items": list, "total": total})
}

func CreateGrowthMetricRecord(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	var req GrowthMetricCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if !growthMetricTypes[req.MetricType] {
		common.Error(c, 40002, "指标类型错误")
		return
	}
	if req.RecordDate == "" {
		common.Error(c, 40004, "日期不能为空")
		return
	}
	d, err := time.ParseInLocation("2006-01-02", req.RecordDate, time.Local)
	if err != nil {
		common.Error(c, 40005, "日期格式错误")
		return
	}
	item := models.GrowthMetricRecord{
		UserID:     cl.UserID,
		MetricType: req.MetricType,
		RecordDate: d,
	}
	if req.MetricType == "vision" {
		if req.LeftValue <= 0 || req.RightValue <= 0 {
			common.Error(c, 40003, "左右眼数值必须大于0")
			return
		}
		item.LeftValue = req.LeftValue
		item.RightValue = req.RightValue
	} else {
		if req.Value <= 0 {
			common.Error(c, 40003, "数值必须大于0")
			return
		}
		item.Value = req.Value
	}
	if err := db.DB().Create(&item).Error; err != nil {
		common.Error(c, 50002, "创建失败")
		return
	}
	common.Ok(c, item)
}

func DeleteGrowthMetricRecord(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		common.Error(c, 40001, "无效ID")
		return
	}
	var item models.GrowthMetricRecord
	if err := db.DB().First(&item, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if item.UserID != cl.UserID {
		common.Error(c, 40301, "无权操作")
		return
	}
	if err := db.DB().Delete(&item).Error; err != nil {
		common.Error(c, 50003, "删除失败")
		return
	}
	common.Ok(c, gin.H{"id": item.ID})
}
