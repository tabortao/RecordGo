package handlers

import (
	"net/http"
	"path/filepath"
	"recordgo/internal/common"
	"recordgo/internal/config"
	"recordgo/internal/db"
	"recordgo/internal/models"
	s "recordgo/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type presignReq struct {
	ResourceType string `json:"resource_type"`
	UserID       uint   `json:"user_id"`
	TaskID       *uint  `json:"task_id"`
	ContentType  string `json:"content_type"`
	Ext          string `json:"ext"`
}

func StorageInfo(c *gin.Context) {
	st := s.Get()
	drv := ""
	pub := ""
	if st != nil {
		drv = st.DriverName()
	}
	cfg, _ := config.Load()
	if cfg != nil {
		pub = cfg.S3PublicBaseURL
	}
	common.Ok(c, gin.H{"driver": drv, "public_base_url": pub})
}

func StoragePresign(c *gin.Context) {
	var req presignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40000, "参数错误")
		return
	}
	if req.UserID == 0 || strings.TrimSpace(req.ResourceType) == "" || strings.TrimSpace(req.ContentType) == "" {
		common.Error(c, 40000, "参数错误")
		return
	}
	// 备注附件与任务图片上传需校验 VIP 权限
	if req.ResourceType == "task_attachment_img" || req.ResourceType == "task_attachment_audio" || req.ResourceType == "task_image" {
		if !requireVIP(c, req.UserID) {
			return
		}
	}
	key, err := buildKey(req)
	if err != nil {
		common.Error(c, 40000, "参数错误")
		return
	}
	st := s.Get()
	if st == nil {
		common.Error(c, 50010, "存储未初始化")
		return
	}
	url, headers, err := st.PresignPut(c.Request.Context(), key, req.ContentType, 300)
	if err != nil {
		common.Error(c, 50011, "生成上传地址失败")
		return
	}
	common.Ok(c, gin.H{"upload_url": url, "object_key": key, "headers": headers, "expires": 300})
}

func StoragePresignView(c *gin.Context) {
	key := strings.TrimSpace(c.Query("object_key"))
	if key == "" {
		common.Error(c, 40000, "参数错误")
		return
	}
	st := s.Get()
	if st == nil {
		common.Error(c, 50010, "存储未初始化")
		return
	}
	url, err := st.PresignGet(c.Request.Context(), key, 300)
	if err != nil {
		common.Error(c, 50012, "生成访问地址失败")
		return
	}
	common.Ok(c, gin.H{"url": url})
}

func StoragePut(c *gin.Context) {
	key := strings.TrimSpace(c.Query("key"))
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}
	st := s.Get()
	if st == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "存储未初始化", "data": nil})
		return
	}
	if st.DriverName() != "local" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": 405, "message": "不支持", "data": nil})
		return
	}
	// 针对备注附件（图片/音频）本地 PUT，解析 key 中的用户 ID 并强制 VIP 校验
	// 约定：images/{user_id}/task_images/img_* 或 audio_*
	{
		parts := strings.Split(strings.Trim(key, "/"), "/")
		// 兼容可能带前缀，如 blog/images/...，故扫描到 "images" 位置
		idx := -1
		for i, p := range parts {
			if p == "images" {
				idx = i
				break
			}
		}
		if idx >= 0 && len(parts) > idx+3 {
			uidStr := parts[idx+1]
			folder := parts[idx+2]
			base := parts[len(parts)-1]
			if folder == "task_images" && (strings.HasPrefix(base, "img_") || strings.HasPrefix(base, "audio_") || strings.HasPrefix(base, "task_")) {
				if uid, err := strconv.Atoi(uidStr); err == nil {
					if !requireVIP(c, uint(uid)) {
						return
					}
				}
			}
		}
	}
	ct := c.GetHeader("Content-Type")
	r := c.Request.Body
	defer r.Close()
	if err := st.Put(c.Request.Context(), key, r, ct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"path": filepath.ToSlash("uploads/" + key)}})
}

// requireVIP 校验用户 VIP 权限：终身 VIP 或在有效期内的 VIP；失败返回 403 并记录审计日志
func requireVIP(c *gin.Context, uid uint) bool {
	var u models.User
	if err := db.DB().First(&u, uid).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "需要 VIP 会员权限", "data": nil})
		zap.L().Warn("VIP 校验失败：用户不存在", zap.Uint("uid", uid))
		return false
	}
	now := time.Now()
	ok := u.IsLifetimeVIP || (u.IsVIP && u.VIPExpireTime != nil && u.VIPExpireTime.After(now))
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "该功能需要开通VIP会员", "data": nil})
		zap.L().Warn("备注功能未授权访问", zap.Uint("uid", uid), zap.Any("vip_expire_time", u.VIPExpireTime), zap.Bool("is_vip", u.IsVIP), zap.Bool("is_lifetime_vip", u.IsLifetimeVIP))
		return false
	}
	return true
}

func buildKey(req presignReq) (string, error) {
	ts := time.Now().Unix()
	ns := time.Now().UnixNano()
	ext := strings.TrimLeft(strings.ToLower(req.Ext), ".")
	if ext == "" {
		ext = "webp"
	}
	// 可选路径前缀（例如：blog），由环境变量 S3_KEY_PREFIX 控制
	cfg, _ := config.Load()
	prefix := strings.Trim(strings.TrimSpace(cfg.S3KeyPrefix), "/")
	if prefix != "" {
		prefix = prefix + "/"
	}
	switch req.ResourceType {
	case "avatar":
		return filepath.ToSlash(prefix + "images/" + strconv.Itoa(int(req.UserID)) + "/avatars/" + "avatar_" + strconv.Itoa(int(req.UserID)) + "_" + strconv.FormatInt(ts, 10) + "_" + strconv.FormatInt(ns, 10) + "." + ext), nil
	case "task_image":
		if req.TaskID == nil || *req.TaskID == 0 {
			return "", nil
		}
		return filepath.ToSlash(prefix + "images/" + strconv.Itoa(int(req.UserID)) + "/task_images/" + "task_" + strconv.Itoa(int(*req.TaskID)) + "_" + strconv.FormatInt(ts, 10) + "_" + strconv.FormatInt(ns, 10) + "." + ext), nil
	case "task_attachment_img":
		if req.TaskID == nil || *req.TaskID == 0 {
			return "", nil
		}
		return filepath.ToSlash(prefix + "images/" + strconv.Itoa(int(req.UserID)) + "/task_images/" + "img_" + strconv.Itoa(int(req.UserID)) + "_" + strconv.Itoa(int(*req.TaskID)) + "_" + strconv.FormatInt(ts, 10) + "_" + strconv.FormatInt(ns, 10) + "." + ext), nil
	case "task_attachment_audio":
		if req.TaskID == nil || *req.TaskID == 0 {
			return "", nil
		}
		aext := ext
		if aext == "" {
			aext = "mp3"
		}
		return filepath.ToSlash(prefix + "images/" + strconv.Itoa(int(req.UserID)) + "/task_images/" + "audio_" + strconv.Itoa(int(req.UserID)) + "_" + strconv.Itoa(int(*req.TaskID)) + "_" + strconv.FormatInt(ts, 10) + "_" + strconv.FormatInt(ns, 10) + "." + aext), nil
	case "wish_image":
		return filepath.ToSlash(prefix + "images/" + strconv.Itoa(int(req.UserID)) + "/wish/" + "wish_" + strconv.Itoa(int(req.UserID)) + "_" + strconv.FormatInt(ts, 10) + "_" + strconv.FormatInt(ns, 10) + "." + ext), nil
	default:
		return "", nil
	}
}
