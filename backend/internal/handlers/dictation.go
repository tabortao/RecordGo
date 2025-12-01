package handlers

import (
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Word Banks ---

func ListWordBanks(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var list []models.DictationWordBank
	if err := db.DB().Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error; err != nil {
		common.Error(c, 500, "查询失败")
		return
	}
	common.Ok(c, list)
}

func CreateWordBank(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var form models.DictationWordBank
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}
	form.ID = 0
	form.UserID = userID
	if err := db.DB().Create(&form).Error; err != nil {
		common.Error(c, 500, "创建失败")
		return
	}
	common.Ok(c, form)
}

func ListDictationHistory(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var list []models.DictationHistory
	if err := db.DB().Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error; err != nil {
		common.Error(c, 500, "查询失败")
		return
	}
	common.Ok(c, list)
}

func UpdateWordBank(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	id := c.Param("id")
	var wb models.DictationWordBank
	if err := db.DB().Where("id = ? AND user_id = ?", id, userID).First(&wb).Error; err != nil {
		common.Error(c, 404, "词库不存在")
		return
	}
	var form models.DictationWordBank
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}
	// Update fields
	wb.Title = form.Title
	wb.EducationStage = form.EducationStage
	wb.Subject = form.Subject
	wb.Version = form.Version
	wb.Grade = form.Grade
	wb.Content = form.Content

	if err := db.DB().Save(&wb).Error; err != nil {
		common.Error(c, 500, "更新失败")
		return
	}
	common.Ok(c, wb)
}

func DeleteWordBank(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	id := c.Param("id")
	if err := db.DB().Where("id = ? AND user_id = ?", id, userID).Delete(&models.DictationWordBank{}).Error; err != nil {
		common.Error(c, 500, "删除失败")
		return
	}
	common.Ok(c, nil)
}

// --- Settings ---

func GetDictationSettings(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var s models.DictationSettings
	if err := db.DB().Where("user_id = ?", userID).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default settings
			s = models.DictationSettings{
				UserID:          userID,
				SplitRule:       "newline",
				PlayMode:        "read",
				OrderMode:       "sequential",
				RepeatCount:     1,
				IntervalSeconds: 3,
				Speed:           1.0,
			}
			common.Ok(c, s)
			return
		}
		common.Error(c, 500, "查询失败")
		return
	}
	common.Ok(c, s)
}

func UpdateDictationSettings(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var form models.DictationSettings
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}
	form.UserID = userID
	// Use Save to update or create
	if err := db.DB().Save(&form).Error; err != nil {
		common.Error(c, 500, "保存失败")
		return
	}
	common.Ok(c, form)
}

// --- Mistakes ---

func ListMistakes(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var list []models.MistakeNotebook
	if err := db.DB().Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error; err != nil {
		common.Error(c, 500, "查询失败")
		return
	}
	common.Ok(c, list)
}

func AddMistake(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var form models.MistakeNotebook
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}
	form.ID = 0
	form.UserID = userID
	if err := db.DB().Create(&form).Error; err != nil {
		common.Error(c, 500, "添加失败")
		return
	}
	common.Ok(c, form)
}

func DeleteMistake(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	id := c.Param("id")
	if err := db.DB().Where("id = ? AND user_id = ?", id, userID).Delete(&models.MistakeNotebook{}).Error; err != nil {
		common.Error(c, 500, "删除失败")
		return
	}
	common.Ok(c, nil)
}

// --- History ---
func AddDictationHistory(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 401, "未授权")
		return
	}
	userID := cl.UserID

	var form models.DictationHistory
	if err := c.ShouldBindJSON(&form); err != nil {
		common.Error(c, 400, "参数错误")
		return
	}
	form.ID = 0
	form.UserID = userID
	if err := db.DB().Create(&form).Error; err != nil {
		common.Error(c, 500, "保存失败")
		return
	}
	common.Ok(c, form)
}
