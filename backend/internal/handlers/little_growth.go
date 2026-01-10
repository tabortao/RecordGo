package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"
	"recordgo/internal/common"
	"recordgo/internal/config"
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

	// Determine family scope
	var familyIDs []uint
	familyIDs = append(familyIDs, uid)

	// Logic:
	// If I am parent (ParentID == nil): Include all my children (ParentID = MyID)
	// If I am child (ParentID != nil): Include my parent (ParentID) and siblings (ParentID = MyParentID)
	// Simplified: Find "Root" ID.
	var rootID uint
	if cl.ParentID != nil && *cl.ParentID > 0 {
		rootID = *cl.ParentID
	} else {
		rootID = uid
	}

	// Add root if not self
	if rootID != uid {
		familyIDs = append(familyIDs, rootID)
	}

	// Find all children of root
	var children []models.User
	if err := db.DB().Select("id").Where("parent_id = ?", rootID).Find(&children).Error; err == nil {
		for _, child := range children {
			if child.ID != uid { // avoid duplicate
				familyIDs = append(familyIDs, child.ID)
			}
		}
	}

	// Filter by favorite
	isFavorite := c.Query("is_favorite") == "true"

	var list []models.GrowthRecord
	var total int64

	q := db.DB().Model(&models.GrowthRecord{}).Where("user_id IN ?", familyIDs)

	// Visibility check:
	// If I am the author, I can see everything (Private + Family).
	// If I am not the author, I can only see Family (0).
	// So: (user_id = uid) OR (visibility = 0)
	q = q.Where("user_id = ? OR visibility = 0", uid)

	if isFavorite {
		q = q.Where("is_favorite = ?", true)
	}

	q.Count(&total)
	// Sort by is_pinned desc, date desc, created_at desc
	// Preload Comments and Comments.User
	// Also preload Author (User)
	if err := q.Order("is_pinned DESC, date DESC, created_at DESC").
		Preload("User").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Comments.User").
		Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		common.Error(c, 50001, "查询失败")
		return
	}
	common.Ok(c, gin.H{"items": list, "total": total, "page": page, "page_size": size})
}

// ToggleFavoriteGrowthRecord 切换收藏状态
func ToggleFavoriteGrowthRecord(c *gin.Context) {
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
	// Permission check (can favorite own or family?)
	allowed := false
	if r.UserID == cl.UserID {
		allowed = true
	} else {
		// Check family relationship
		var rootID uint
		if cl.ParentID != nil && *cl.ParentID > 0 {
			rootID = *cl.ParentID
		} else {
			rootID = cl.UserID
		}

		// If author is root (and I am child of root) -> OK
		if r.UserID == rootID {
			allowed = true
		} else {
			// Check if author is child of root
			var count int64
			db.DB().Model(&models.User{}).Where("id = ? AND parent_id = ?", r.UserID, rootID).Count(&count)
			if count > 0 {
				allowed = true
			}
		}
	}

	if !allowed {
		common.Error(c, 40301, "无权操作")
		return
	}

	r.IsFavorite = !r.IsFavorite
	if err := db.DB().Save(&r).Error; err != nil {
		common.Error(c, 50002, "操作失败")
		return
	}
	common.Ok(c, gin.H{"id": r.ID, "is_favorite": r.IsFavorite})
}

// AddGrowthComment 添加评论
func AddGrowthComment(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		common.Error(c, 40001, "无效ID")
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		common.Error(c, 40001, "评论内容不能为空")
		return
	}

	// Verify record exists
	var r models.GrowthRecord
	if err := db.DB().First(&r, id).Error; err != nil {
		common.Error(c, 40404, "记录不存在")
		return
	}

	// Permission check (can comment on own or family?)
	allowed := false
	if r.UserID == cl.UserID {
		allowed = true
	} else {
		// Check family relationship
		var rootID uint
		if cl.ParentID != nil && *cl.ParentID > 0 {
			rootID = *cl.ParentID
		} else {
			rootID = cl.UserID
		}

		// If author is root (and I am child of root) -> OK
		if r.UserID == rootID {
			allowed = true
		} else {
			// Check if author is child of root
			var count int64
			db.DB().Model(&models.User{}).Where("id = ? AND parent_id = ?", r.UserID, rootID).Count(&count)
			if count > 0 {
				allowed = true
			}
		}
	}

	if !allowed {
		common.Error(c, 40301, "无权评论")
		return
	}

	comment := models.GrowthComment{
		RecordID: uint(id),
		UserID:   cl.UserID,
		Content:  req.Content,
	}
	if err := db.DB().Create(&comment).Error; err != nil {
		common.Error(c, 50002, "评论失败")
		return
	}

	// Load user info for response
	db.DB().Preload("User").First(&comment, comment.ID)

	common.Ok(c, comment)
}

// TogglePinGrowthRecord 切换置顶状态
func TogglePinGrowthRecord(c *gin.Context) {
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
		common.Error(c, 40301, "无权操作")
		return
	}

	r.IsPinned = !r.IsPinned
	if err := db.DB().Save(&r).Error; err != nil {
		common.Error(c, 50002, "操作失败")
		return
	}
	common.Ok(c, gin.H{"id": r.ID, "is_pinned": r.IsPinned})
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

// GetGrowthRecord 获取单个记录
func GetGrowthRecord(c *gin.Context) {
	cl := extractClaims(c)
	if cl == nil {
		common.Error(c, 40100, "未登录")
		return
	}
	id := c.Param("id")
	uid := cl.UserID

	// Determine family scope (similar to ListGrowthRecords)
	var rootID uint
	if cl.ParentID != nil && *cl.ParentID > 0 {
		rootID = *cl.ParentID
	} else {
		rootID = uid
	}

	var familyIDs []uint
	familyIDs = append(familyIDs, rootID)
	var children []models.User
	if err := db.DB().Select("id").Where("parent_id = ?", rootID).Find(&children).Error; err == nil {
		for _, child := range children {
			familyIDs = append(familyIDs, child.ID)
		}
	}

	var r models.GrowthRecord
	// Preload author and comments
	q := db.DB().Model(&models.GrowthRecord{}).
		Preload("User").
		Preload("Comments").
		Preload("Comments.User").
		Where("id = ?", id).
		Where("user_id IN ?", familyIDs)

	// Visibility Check
	// (user_id = uid) OR (visibility = 0)
	q = q.Where("user_id = ? OR visibility = 0", uid)

	if err := q.First(&r).Error; err != nil {
		common.Error(c, 40401, "记录不存在或无权查看")
		return
	}

	common.Ok(c, r)
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
		Content    string   `json:"content"`
		Date       string   `json:"date"`
		Images     []string `json:"images"`
		Audio      string   `json:"audio"`
		Tags       []string `json:"tags"` // IDs as strings or numbers
		Visibility int      `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}

	// Parse Date (support full datetime or date only)
	var date time.Time
	var err error
	if len(req.Date) > 10 {
		date, err = time.Parse("2006-01-02 15:04", req.Date)
		if err != nil {
			date, err = time.Parse("2006-01-02 15:04:05", req.Date)
		}
	} else {
		date, err = time.Parse("2006-01-02", req.Date)
	}

	if err != nil {
		date = time.Now()
	}

	imagesJSON, _ := json.Marshal(req.Images)
	tagsJSON, _ := json.Marshal(req.Tags)

	r := models.GrowthRecord{
		UserID:     uid,
		Date:       date,
		Content:    req.Content,
		Images:     string(imagesJSON),
		Audio:      req.Audio,
		Tags:       string(tagsJSON),
		Visibility: req.Visibility,
	}

	if err := db.DB().Create(&r).Error; err != nil {
		common.Error(c, 50002, "创建失败")
		return
	}

	common.Ok(c, r)
}

// UpdateGrowthRecord 更新记录
func UpdateGrowthRecord(c *gin.Context) {
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
		common.Error(c, 40301, "无权修改")
		return
	}

	// Parse JSON body
	var req struct {
		Content    string   `json:"content"`
		Date       string   `json:"date"`
		Images     []string `json:"images"`
		Audio      string   `json:"audio"`
		Tags       []string `json:"tags"`
		Visibility int      `json:"visibility"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}

	// Update fields
	r.Content = req.Content
	r.Visibility = req.Visibility

	// Parse Date
	if req.Date != "" {
		var date time.Time
		var err error
		if len(req.Date) > 10 {
			date, err = time.Parse("2006-01-02 15:04", req.Date)
			if err != nil {
				date, err = time.Parse("2006-01-02 15:04:05", req.Date)
			}
		} else {
			date, err = time.Parse("2006-01-02", req.Date)
		}
		if err == nil {
			r.Date = date
		}
	}

	imagesJSON, _ := json.Marshal(req.Images)
	tagsJSON, _ := json.Marshal(req.Tags)
	r.Images = string(imagesJSON)
	r.Tags = string(tagsJSON)
	if req.Audio != "" {
		r.Audio = req.Audio
	}

	if err := db.DB().Save(&r).Error; err != nil {
		common.Error(c, 50002, "更新失败")
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
	
	// Apply S3_KEY_PREFIX if configured
	cfg, _ := config.Load()
	prefix := ""
	if cfg != nil {
		prefix = strings.Trim(strings.TrimSpace(cfg.S3KeyPrefix), "/")
		if prefix != "" {
			prefix = prefix + "/"
		}
	}
	
	key := fmt.Sprintf("%simages/%d/little_growth/%s", prefix, uid, filename)

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
