package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"
	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"
	"recordgo/internal/storage"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListGrowthRecords 获取小成长记录列表
func ListGrowthRecords(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	uid := cl.UserID
	// 简单分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	var list []models.GrowthRecord
	var total int64

	q := db.DB().Model(&models.GrowthRecord{}).Where("user_id = ?", uid)
	q.Count(&total)
	if err := q.Order("date DESC, created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		common.Error(c, 50001, "查询失败")
		return
	}
	common.Ok(c, gin.H{"items": list, "total": total, "page": page, "page_size": size})
}

// CreateGrowthRecord 创建记录（支持 multipart/form-data 上传图片与音频）
func CreateGrowthRecord(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	uid := cl.UserID

	// Parse multipart
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		common.Error(c, 40001, "解析表单失败")
		return
	}

	content := c.PostForm("content")
	dateStr := c.PostForm("date")
	tagsStr := c.PostForm("tags") // Expecting JSON array of strings or IDs

	// Parse Date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now()
	}

	// Upload Images
	form, _ := c.MultipartForm()
	files := form.File["images"]
	var imagePaths []string
	st := storage.Get()

	for _, file := range files {
		path, err := uploadGrowthFile(st, file, uid, "image")
		if err == nil {
			imagePaths = append(imagePaths, path)
		} else {
			zap.L().Error("CreateGrowthRecord: upload image failed", zap.Error(err))
		}
	}

	// Upload Audio
	var audioPath string
	audioFile, _ := c.FormFile("audio")
	if audioFile != nil {
		path, err := uploadGrowthFile(st, audioFile, uid, "audio")
		if err == nil {
			audioPath = path
		} else {
			zap.L().Error("CreateGrowthRecord: upload audio failed", zap.Error(err))
		}
	}

	// JSON encode images
	imagesJSON, _ := json.Marshal(imagePaths)

	record := models.GrowthRecord{
		UserID:  uid,
		Date:    date,
		Content: content,
		Images:  string(imagesJSON),
		Audio:   audioPath,
		Tags:    tagsStr,
	}

	if err := db.DB().Create(&record).Error; err != nil {
		common.Error(c, 50002, "保存记录失败")
		return
	}

	common.Ok(c, record)
}

// DeleteGrowthRecord 删除记录
func DeleteGrowthRecord(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	id := c.Param("id")
	var r models.GrowthRecord
	if err := db.DB().First(&r, id).Error; err != nil {
		common.Error(c, 40401, "记录不存在")
		return
	}
	if r.UserID != cl.UserID {
		common.Error(c, 40301, "无权删除")
		return
	}

	// Delete files from storage
	st := storage.Get()
	ctx := context.Background()

	// Images
	var images []string
	if err := json.Unmarshal([]byte(r.Images), &images); err == nil {
		for _, imgPath := range images {
			filename := path.Base(imgPath)
			key := fmt.Sprintf("%d/little_growth/%s", r.UserID, filename)
			_ = st.Delete(ctx, key)
		}
	}

	// Audio
	if r.Audio != "" {
		filename := path.Base(r.Audio)
		key := fmt.Sprintf("%d/little_growth/%s", r.UserID, filename)
		_ = st.Delete(ctx, key)
	}

	if err := db.DB().Delete(&r).Error; err != nil {
		common.Error(c, 50003, "删除失败")
		return
	}
	common.Ok(c, gin.H{"id": id})
}

// Helper to upload file
func uploadGrowthFile(st storage.Storage, file *multipart.FileHeader, uid uint, ftype string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		if ftype == "audio" {
			ext = ".m4a"
		} else {
			ext = ".webp"
		}
	}

	// Key: {uid}/little_growth/{type}_{timestamp}_{random}{ext}
	filename := fmt.Sprintf("%s_%d%s", ftype, time.Now().UnixNano(), ext)
	key := fmt.Sprintf("%d/little_growth/%s", uid, filename)

	// We need to determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err := st.Put(context.Background(), key, src, contentType); err != nil {
		return "", err
	}

	// Return the path/URL
	url, isPublic := st.PublicURL(key)
	if isPublic {
		return url, nil
	}

	// Fallback for local: return relative path
	return "uploads/" + key, nil
}
