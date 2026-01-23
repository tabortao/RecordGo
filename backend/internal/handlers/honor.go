package handlers

import (
	"time"

	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
)

type HonorCreateReq struct {
	Title     string `json:"title"`
	Issuer    string `json:"issuer"`
	AwardedAt string `json:"awarded_at"`
	PhotoURL  string `json:"photo_url"`
	Remark    string `json:"remark"`
}

// 中文注释：更新荣誉请求体，字段可选
type HonorUpdateReq struct {
	Title     *string `json:"title"`
	Issuer    *string `json:"issuer"`
	AwardedAt *string `json:"awarded_at"`
	PhotoURL  *string `json:"photo_url"`
	Remark    *string `json:"remark"`
}

func ListHonors(c *gin.Context) {
	// 中文注释：列出当前用户荣誉记录
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	uid := cl.UserID
	if cl.ParentID != nil {
		uid = *cl.ParentID
	}
	var items []models.Honor
	if err := db.DB().Where("user_id = ?", uid).Order("awarded_at DESC, created_at DESC").Find(&items).Error; err != nil {
		common.Error(c, 50001, "查询失败")
		return
	}
	common.Ok(c, gin.H{"items": items})
}

func CreateHonor(c *gin.Context) {
	// 中文注释：创建荣誉记录
	if !hasPermission(c, "honors", "create") {
		deny(c, "无权限创建荣誉")
		return
	}
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	var req HonorCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.Title == "" || req.AwardedAt == "" || req.PhotoURL == "" {
		common.Error(c, 40002, "标题/颁发日期/照片必填")
		return
	}
	d, err := time.ParseInLocation("2006-01-02", req.AwardedAt, time.Local)
	if err != nil {
		common.Error(c, 40003, "日期格式错误")
		return
	}
	uid := cl.UserID
	if cl.ParentID != nil {
		uid = *cl.ParentID
	}
	h := models.Honor{
		UserID:    uid,
		Title:     req.Title,
		Issuer:    req.Issuer,
		AwardedAt: d,
		PhotoURL:  req.PhotoURL,
		Remark:    req.Remark,
	}
	if err := db.DB().Create(&h).Error; err != nil {
		common.Error(c, 50002, "创建失败")
		return
	}
	common.Ok(c, h)
}

func GetHonor(c *gin.Context) {
	// 中文注释：查询荣誉详情
	id := c.Param("id")
	var h models.Honor
	if err := db.DB().First(&h, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if !canAccessUser(c, h.UserID) {
		common.Error(c, 40301, "无权查看")
		return
	}
	common.Ok(c, h)
}

func UpdateHonor(c *gin.Context) {
	// 中文注释：更新荣誉记录
	if !hasPermission(c, "honors", "edit") {
		deny(c, "无权限编辑荣誉")
		return
	}
	id := c.Param("id")
	var req HonorUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	var h models.Honor
	if err := db.DB().First(&h, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if !canAccessUser(c, h.UserID) {
		common.Error(c, 40301, "无权操作")
		return
	}
	if req.Title != nil {
		h.Title = *req.Title
	}
	if req.Issuer != nil {
		h.Issuer = *req.Issuer
	}
	if req.AwardedAt != nil {
		d, err := time.ParseInLocation("2006-01-02", *req.AwardedAt, time.Local)
		if err != nil {
			common.Error(c, 40003, "日期格式错误")
			return
		}
		h.AwardedAt = d
	}
	if req.PhotoURL != nil {
		h.PhotoURL = *req.PhotoURL
	}
	if req.Remark != nil {
		h.Remark = *req.Remark
	}
	if err := db.DB().Save(&h).Error; err != nil {
		common.Error(c, 50004, "更新失败")
		return
	}
	common.Ok(c, h)
}

func DeleteHonor(c *gin.Context) {
	// 中文注释：删除荣誉记录
	if !hasPermission(c, "honors", "delete") {
		deny(c, "无权限删除荣誉")
		return
	}
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	id := c.Param("id")
	var h models.Honor
	if err := db.DB().First(&h, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if h.UserID != cl.UserID {
		common.Error(c, 40301, "无权操作")
		return
	}
	if err := db.DB().Delete(&h).Error; err != nil {
		common.Error(c, 50003, "删除失败")
		return
	}
	common.Ok(c, gin.H{"id": h.ID})
}
