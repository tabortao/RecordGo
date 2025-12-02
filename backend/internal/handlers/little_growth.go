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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Records ---

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

// UploadGrowthFile 上传文件 (图片或音频)
func UploadGrowthFile(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.Error(c, 40001, "请选择文件")
		return
	}

	// Determine type from extension or Content-Type or query param
	ftype := "image"
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".m4a" || ext == ".mp3" || ext == ".wav" {
		ftype = "audio"
	} else if strings.HasPrefix(file.Header.Get("Content-Type"), "audio/") {
		ftype = "audio"
	}

	st := storage.Get()
	path, err := uploadGrowthFile(st, file, cl.UserID, ftype)
	if err != nil {
		common.Error(c, 50001, "上传失败")
		return
	}

	common.Ok(c, gin.H{"url": path, "type": ftype})
}

// CreateGrowthRecord 创建记录
func CreateGrowthRecord(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	uid := cl.UserID

	// Parse JSON body
	var req struct {
		Content string   `json:"content"`
		Date    string   `json:"date"`
		Images  []string `json:"images"`
		Audio   string   `json:"audio"`
		Tags    []string `json:"tags"` // IDs as strings or numbers
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}

	// Parse Date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	imagesJSON, _ := json.Marshal(req.Images)
	tagsJSON, _ := json.Marshal(req.Tags)

	r := models.GrowthRecord{
		UserID:  uid,
		Date:    date,
		Content: req.Content,
		Images:  string(imagesJSON),
		Audio:   req.Audio,
		Tags:    string(tagsJSON),
	}

	if err := db.DB().Create(&r).Error; err != nil {
		common.Error(c, 50002, "创建失败")
		return
	}

	common.Ok(c, r)
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
			// Handle full URL or relative path
			// If full URL (http...), we need to extract key.
			// But storage.Delete expects key.
			// Our uploadGrowthFile returns relative path "uploads/..." or full URL.
			// If it's S3, it's likely "images/..." key if we strip host?
			// Actually, uploadGrowthFile returns what `st.PublicURL` returns.
			// If using local storage, it returns "uploads/images/..."
			// If S3, it returns public URL.
			// This is tricky. Deleting by URL is hard if we don't know the key mapping.
			// But wait, `uploadGrowthFile` constructs key: `images/{uid}/little_growth/{filename}`
			// We can reconstruct key from filename if we assume structure.

			// Let's try to extract filename
			filename := path.Base(imgPath)
			// Reconstruct key
			key := fmt.Sprintf("images/%d/little_growth/%s", r.UserID, filename)
			_ = st.Delete(ctx, key)
		}
	}

	// Audio
	if r.Audio != "" {
		filename := path.Base(r.Audio)
		key := fmt.Sprintf("images/%d/little_growth/%s", r.UserID, filename)
		_ = st.Delete(ctx, key)
	}

	if err := db.DB().Delete(&r).Error; err != nil {
		common.Error(c, 50003, "删除失败")
		return
	}
	common.Ok(c, gin.H{"id": id})
}

// --- Tags ---

// ListGrowthTags 获取标签列表
func ListGrowthTags(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	var tags []models.GrowthTag
	if err := db.DB().Where("user_id = ?", cl.UserID).Find(&tags).Error; err != nil {
		common.Error(c, 50001, "查询标签失败")
		return
	}
	common.Ok(c, tags)
}

// CreateGrowthTag 创建标签
func CreateGrowthTag(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	var req struct {
		Name     string `json:"name"`
		Color    string `json:"color"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.Name == "" {
		common.Error(c, 40001, "标签名不能为空")
		return
	}

	// Check duplicate
	var count int64
	db.DB().Model(&models.GrowthTag{}).Where("user_id = ? AND name = ?", cl.UserID, req.Name).Count(&count)
	if count > 0 {
		common.Error(c, 40002, "标签已存在")
		return
	}

	// Generate random color if not provided
	color := req.Color
	if color == "" {
		// Palette of nice colors
		colors := []string{
			"#FFB6C1", "#FF69B4", "#FFD700", "#FFA07A", "#90EE90",
			"#20B2AA", "#87CEFA", "#9370DB", "#FF6347", "#40E0D0",
			"#EE82EE", "#F0E68C", "#E6E6FA", "#DDA0DD", "#B0C4DE",
		}
		// Simple random based on time
		idx := time.Now().UnixNano() % int64(len(colors))
		color = colors[idx]
	}

	t := models.GrowthTag{
		UserID:   cl.UserID,
		Name:     req.Name,
		Color:    color,
		ParentID: req.ParentID,
	}
	if err := db.DB().Create(&t).Error; err != nil {
		common.Error(c, 50002, "创建失败")
		return
	}
	common.Ok(c, t)
}

// UpdateGrowthTag 更新标签
func UpdateGrowthTag(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	id := c.Param("id")
	var t models.GrowthTag
	if err := db.DB().First(&t, id).Error; err != nil {
		common.Error(c, 40401, "标签不存在")
		return
	}
	if t.UserID != cl.UserID {
		common.Error(c, 40301, "无权修改")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Color    string `json:"color"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}

	t.Name = req.Name
	if req.Color != "" {
		t.Color = req.Color
	}
	t.ParentID = req.ParentID
	if err := db.DB().Save(&t).Error; err != nil {
		common.Error(c, 50002, "更新失败")
		return
	}
	common.Ok(c, t)
}

// DeleteGrowthTag 删除标签
func DeleteGrowthTag(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var t models.GrowthTag
	if err := db.DB().First(&t, id).Error; err != nil {
		common.Error(c, 40401, "标签不存在")
		return
	}
	if t.UserID != cl.UserID {
		common.Error(c, 40301, "无权删除")
		return
	}

	// Transaction to delete tag and update records
	err := db.DB().Transaction(func(tx *gorm.DB) error {
		// 1. Delete tag
		if err := tx.Delete(&t).Error; err != nil {
			return err
		}

		// 2. Find records containing this tag
		// Since tags are stored as JSON string ["id1", "id2"], we need to use LIKE or fetch all and process.
		// SQLite JSON support varies, and we are using "pure go sqlite".
		// Simplest: Fetch all records for user, decode tags, remove, save if changed.
		var records []models.GrowthRecord
		if err := tx.Where("user_id = ?", cl.UserID).Find(&records).Error; err != nil {
			return err
		}

		targetID := strconv.Itoa(id)
		for _, r := range records {
			var tags []string
			if err := json.Unmarshal([]byte(r.Tags), &tags); err == nil {
				newTags := make([]string, 0, len(tags))
				changed := false
				for _, tid := range tags {
					if tid == targetID {
						changed = true
					} else {
						newTags = append(newTags, tid)
					}
				}
				if changed {
					tagsJSON, _ := json.Marshal(newTags)
					r.Tags = string(tagsJSON)
					if err := tx.Save(&r).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
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

	// Key: images/{uid}/little_growth/{type}_{timestamp}_{random}{ext}
	filename := fmt.Sprintf("%s_%d%s", ftype, time.Now().UnixNano(), ext)
	key := fmt.Sprintf("images/%d/little_growth/%s", uid, filename)

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
