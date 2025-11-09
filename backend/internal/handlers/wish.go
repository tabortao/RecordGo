package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"recordgo/internal/common"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
)

// 中文注释：默认 6 个心愿占位返回（用于前端初始化与兜底显示）
func DefaultWishes(c *gin.Context) {
	// 中文注释：内置心愿默认所需金币为 1，并提供清晰的描述、单位与兑换数量
	// 说明：为保证前后端一致性，单位与兑换数量与描述对应。如“1金币可兑换玩10分钟游戏”→ unit="分钟"，exchange_amount=10
	data := []gin.H{
		{"name": "看电视", "coins": 1, "icon": "看电视.png", "content": "1金币可兑换看10分钟电视", "unit": "分钟", "exchange_amount": 10},
		{"name": "零花钱", "coins": 2, "icon": "零花钱.png", "content": "2金币可兑换1元零花钱", "unit": "元", "exchange_amount": 1},
		{"name": "玩平板", "coins": 1, "icon": "玩平板.png", "content": "1金币可兑换玩10分钟平板", "unit": "分钟", "exchange_amount": 10},
		{"name": "玩手机", "coins": 1, "icon": "玩手机.png", "content": "1金币可兑换玩10分钟手机", "unit": "分钟", "exchange_amount": 10},
		{"name": "玩游戏", "coins": 1, "icon": "玩游戏.png", "content": "1金币可兑换玩10分钟游戏", "unit": "分钟", "exchange_amount": 10},
		{"name": "自由活动", "coins": 1, "icon": "自由活动.png", "content": "1金币可兑换自由活动10分钟", "unit": "分钟", "exchange_amount": 10},
	}
	common.Ok(c, data)
}

// 中文注释：创建心愿请求结构体（字段需与前端一致）
type CreateWishReq struct {
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`
	Content        string `json:"content"`
	Icon           string `json:"icon"` // 中文注释：图标路径或文件名（前端上传后返回保存路径）
	NeedCoins      int    `json:"need_coins"`
	ExchangeAmount int    `json:"exchange_amount"` // 中文注释：如 10 表示 1金币兑换10分钟/个/次
	Unit           string `json:"unit"`
}

// 中文注释：编辑心愿请求结构体
type UpdateWishReq struct {
	Name           *string `json:"name"`
	Content        *string `json:"content"`
	Icon           *string `json:"icon"`
	NeedCoins      *int    `json:"need_coins"`
	ExchangeAmount *int    `json:"exchange_amount"`
	Unit           *string `json:"unit"`
}

// 中文注释：列出某用户心愿，若该用户首次访问，则自动创建 6 个内置心愿
func ListWishes(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		common.Error(c, 40001, "缺少 user_id")
		return
	}
	userID, _ := strconv.Atoi(userIDStr)
	var items []models.Wish
	if err := db.DB().Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error; err != nil {
		common.Error(c, 50010, "查询心愿失败")
		return
	}
	if len(items) == 0 {
		// 首次访问：创建 6 个内置心愿（图标文件名用于前端映射 src/assets/wishs）
		type defWish struct {
			name, content, unit string
			coins, amount       int
		}
		defaults := []defWish{
			{"看电视", "1金币可兑换看10分钟电视", "分钟", 1, 10},
			{"零花钱", "2金币可兑换1元零花钱", "元", 2, 1},
			{"玩平板", "1金币可兑换玩10分钟平板", "分钟", 1, 10},
			{"玩手机", "1金币可兑换玩10分钟手机", "分钟", 1, 10},
			{"玩游戏", "1金币可兑换玩10分钟游戏", "分钟", 1, 10},
			{"自由活动", "1金币可兑换自由活动10分钟", "分钟", 1, 10},
		}
		for _, d := range defaults {
			w := models.Wish{
				// 中文注释：内置心愿采用清晰的描述与单位匹配，确保记录展示一致
				UserID: uint(userID), Name: d.name, Content: d.content, Icon: d.name + ".png",
				NeedCoins: d.coins, ExchangeAmount: d.amount, Unit: d.unit, BuiltIn: true,
			}
			_ = db.DB().Create(&w).Error
		}
		_ = db.DB().Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	}
	common.Ok(c, items)
}

// CreateWish 创建心愿
func CreateWish(c *gin.Context) {
	var req CreateWishReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.UserID == 0 || req.Name == "" || req.NeedCoins <= 0 {
		common.Error(c, 40002, "用户、名称与所需金币必填且合法")
		return
	}
	w := models.Wish{
		UserID:         req.UserID,
		Name:           req.Name,
		Content:        req.Content,
		Icon:           req.Icon,
		NeedCoins:      req.NeedCoins,
		ExchangeAmount: req.ExchangeAmount,
		Unit:           req.Unit,
		BuiltIn:        false,
	}
	if err := db.DB().Create(&w).Error; err != nil {
		common.Error(c, 50011, "创建心愿失败")
		return
	}
	common.Ok(c, w)
}

// GetWish 查询单个心愿详情
func GetWish(c *gin.Context) {
    id := c.Param("id")
    var w models.Wish
    if err := db.DB().First(&w, id).Error; err != nil {
        common.Error(c, 40401, "心愿不存在")
        return
    }
    common.Ok(c, w)
}

// UpdateWish 编辑心愿
func UpdateWish(c *gin.Context) {
	id := c.Param("id")
	var req UpdateWishReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	var w models.Wish
	if err := db.DB().First(&w, id).Error; err != nil {
		common.Error(c, 40401, "心愿不存在")
		return
	}
	if req.Name != nil {
		w.Name = *req.Name
	}
	if req.Content != nil {
		w.Content = *req.Content
	}
	if req.Icon != nil {
		w.Icon = *req.Icon
	}
	if req.NeedCoins != nil {
		w.NeedCoins = *req.NeedCoins
	}
	if req.ExchangeAmount != nil {
		w.ExchangeAmount = *req.ExchangeAmount
	}
	if req.Unit != nil {
		w.Unit = *req.Unit
	}
	if err := db.DB().Save(&w).Error; err != nil {
		common.Error(c, 50012, "更新心愿失败")
		return
	}
	common.Ok(c, w)
}

// DeleteWish 删除心愿（允许删除内置心愿）
func DeleteWish(c *gin.Context) {
	// 中文注释：删除心愿时，若为自定义心愿且图标为上传到 storage/uploads 的文件，则同步删除对应图片
	id := c.Param("id")
	var w models.Wish
	if err := db.DB().First(&w, id).Error; err != nil {
		common.Error(c, 40401, "心愿不存在")
		return
	}
	// 条件：非内置且 Icon 指向 uploads 路径（相对路径）
	// 注意：内置心愿的图标是内置 PNG 文件名（如 “看电视.png”），不应删除；
	// 前端上传的自定义图标存储为 “uploads/images/wish/{用户id}/xxx.webp”。
	if !w.BuiltIn {
		icon := strings.TrimSpace(w.Icon)
		if icon != "" && strings.HasPrefix(icon, "uploads/") {
			root := os.Getenv("STORAGE_ROOT")
			if strings.TrimSpace(root) == "" {
				root = "storage"
			}
			// 拼接绝对路径并删除文件（忽略删除错误，避免影响主流程）
			full := filepath.Join(root, filepath.FromSlash(icon))
			_ = os.Remove(full)
		}
	}
	if err := db.DB().Delete(&w).Error; err != nil {
		common.Error(c, 50013, "删除心愿失败")
		return
	}
	common.Ok(c, gin.H{"id": w.ID})
}

// 中文注释：心愿兑换请求结构体
type ExchangeReq struct {
	UserID uint `json:"user_id"`
	Count  int  `json:"count"` // 中文注释：一次兑换的份数，默认 1
}

// ExchangeWish 兑换心愿：扣减金币、累计心愿兑换次数、写入兑换记录
func ExchangeWish(c *gin.Context) {
	id := c.Param("id")
	var req ExchangeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	if req.Count <= 0 {
		req.Count = 1
	}
	var w models.Wish
	if err := db.DB().First(&w, id).Error; err != nil {
		common.Error(c, 40401, "心愿不存在")
		return
	}
	var u models.User
	if err := db.DB().First(&u, req.UserID).Error; err != nil {
		// 中文注释：开发环境兜底——如果用户不存在，自动创建一个测试用户并给予初始金币
		// 说明：正式环境应通过登录/注册获得用户，这里仅为便于联调前端心愿功能
		u = models.User{ID: req.UserID, Username: fmt.Sprintf("user%d", req.UserID), Role: "user", Coins: 100, Nickname: "测试用户"}
		_ = db.DB().Create(&u).Error
	}
	// 中文注释：按数量计算总金币需求
	totalCost := int64(w.NeedCoins * req.Count)
	if u.Coins < totalCost {
		common.Error(c, 40006, "金币不足")
		return
	}
	// 扣减金币与心愿计数
	u.Coins -= totalCost
	w.Exchanged += req.Count
	if err := db.DB().Save(&u).Error; err != nil {
		common.Error(c, 50014, "扣减金币失败")
		return
	}
	if err := db.DB().Save(&w).Error; err != nil {
		common.Error(c, 50015, "更新心愿失败")
		return
	}
	// 写入兑换记录
	rec := models.WishRecord{
		UserID:    req.UserID,
		WishID:    w.ID,
		WishName:  w.Name,
		CoinsUsed: w.NeedCoins * req.Count,
		Amount:    w.ExchangeAmount * req.Count,
		Unit:      w.Unit,
		Status:    "成功",
		CreatedAt: time.Now(),
	}
	if err := db.DB().Create(&rec).Error; err != nil {
		common.Error(c, 50016, "写入兑换记录失败")
		return
	}
	common.Ok(c, gin.H{"wish": w, "user_coins": u.Coins, "record": rec})
}

// ListWishRecords 分页查询兑换记录（支持 user_id 过滤）
func ListWishRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	q := db.DB().Model(&models.WishRecord{})
	if uid := c.Query("user_id"); uid != "" {
		q = q.Where("user_id = ?", uid)
	}
	var total int64
	q.Count(&total)
	var items []models.WishRecord
	if err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		common.Error(c, 50017, "查询兑换记录失败")
		return
	}
	common.Ok(c, gin.H{"items": items, "total": total, "page": page, "page_size": size})
}

// UploadWishIcon 上传心愿图标（前端需先压缩并转换为 webp）
func UploadWishIcon(c *gin.Context) {
	userID := c.PostForm("user_id")
	if userID == "" {
		common.Error(c, 40001, "缺少 user_id")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		common.Error(c, 40002, "缺少文件")
		return
	}
	// 保存到 storage/uploads/images/wish/{用户id}
	path, err := saveWishIcon(file, userID)
	if err != nil {
		common.Error(c, 50018, "保存图标失败")
		return
	}
	common.Ok(c, gin.H{"path": path})
}

// saveWishIcon 将上传文件保存到约定目录并按规范命名
func saveWishIcon(file *multipart.FileHeader, userID string) (string, error) {
	root := os.Getenv("STORAGE_ROOT")
	if strings.TrimSpace(root) == "" {
		root = "storage"
	}
	dir := filepath.Join(root, "uploads", "images", "wish", userID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}
	// 命名：wish_用户ID_时间戳_uuid.webp（此处 uuid 用纳秒替代以减少依赖）
	filename := fmt.Sprintf("wish_%s_%d_%d.webp", userID, time.Now().Unix(), time.Now().UnixNano())
	full := filepath.Join(dir, filename)
	if err := cSaveUploadedFile(file, full); err != nil {
		return "", err
	}
	// 返回相对路径，前端与后端均可记录
	rel := filepath.ToSlash(filepath.Join("uploads", "images", "wish", userID, filename))
	return rel, nil
}

// cSaveUploadedFile 兼容性的保存封装（避免跨平台路径问题）
func cSaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// 中文注释：手动复制上传内容到目标文件，避免依赖 gin 上下文
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return err
	}
	return nil
}
