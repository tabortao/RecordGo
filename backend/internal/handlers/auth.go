package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"recordgo/internal/common"
	"recordgo/internal/config"
	"recordgo/internal/db"
	"recordgo/internal/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// 中文注释：登录与注册请求结构体（字段与前端一致）
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 中文注释：JWT Claims 结构，包含必要的用户信息
type Claims struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	Permissions string `json:"permissions"`
	ParentID    *uint  `json:"parent_id"`
	// 中文注释：子账号登录令牌（用于令牌刷新后使旧 JWT 失效）；父账号或普通登录可为空
	LoginToken string `json:"login_token"`
	jwt.RegisteredClaims
}

// 中文注释：Register 用户注册；首次注册的用户为 admin，后续为 user
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, 40001, "参数错误")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Password == "" || req.Email == "" {
		common.Error(c, 40001, "参数错误：请填写用户名、密码与邮箱")
		return
	}
	if len([]rune(req.Username)) < 4 {
		common.Error(c, 40004, "用户名不少于4个字符")
		return
	}
	if !isStrongPassword(req.Password) {
		common.Error(c, 40005, "密码需包含数字、大写字母与小写字母")
		return
	}
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		common.Error(c, 40006, "邮箱格式不正确")
		return
	}

	// 用户名唯一性校验
	var exists models.User
	if err := db.DB().Where("username = ?", req.Username).First(&exists).Error; err == nil {
		common.Error(c, 40002, "用户名已存在")
		return
	}

	// 判断是否首个管理员
	var cnt int64
	_ = db.DB().Model(&models.User{}).Where("role = ?", "admin").Count(&cnt).Error
	role := "user"
	if cnt == 0 {
		role = "admin"
	}

	// 密码 SHA256 加密存储（与 PRD 一致）
	h := sha256.Sum256([]byte(req.Password))
	u := models.User{
		Username:    req.Username,
		PasswordSha: hex.EncodeToString(h[:]),
		Role:        role,
		Permissions: `{"view_only": false}`,
		Nickname:    ifEmpty(req.Nickname, req.Username),
		AvatarPath:  "storage/images/avatars/default.png",
		Email:       req.Email,
		Coins:       0,
		Tomatoes:    0,
	}
	if err := db.DB().Create(&u).Error; err != nil {
		common.Error(c, 50001, "注册失败")
		return
	}
	common.Ok(c, gin.H{"id": u.ID, "username": u.Username})
}

// 中文注释：Login 用户登录，成功返回 JWT 与用户信息
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		common.Error(c, 40001, "参数错误：请填写用户名与密码")
		return
	}

	var u models.User
	if err := db.DB().Where("username = ?", req.Username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			common.Error(c, 40401, "用户不存在")
			return
		}
		common.Error(c, 50002, "查询用户失败")
		return
	}
	// 校验密码
	h := sha256.Sum256([]byte(req.Password))
	if u.PasswordSha != hex.EncodeToString(h[:]) {
		common.Error(c, 40003, "密码错误")
		return
	}

	// 记录最后登录时间
	now := time.Now()
	_ = db.DB().Model(&u).Updates(map[string]any{"last_login_time": now}).Error
	// 生成 JWT
	cfg, _ := config.Load()
	// 中文注释：若为子账号，将当前 LoginToken 注入 JWT 用于后续校验；父账号为空字符串
	lt := ""
	if u.ParentID != nil {
		lt = strings.TrimSpace(u.LoginToken)
	}
	claims := Claims{
		UserID:      u.ID,
		Username:    u.Username,
		Role:        u.Role,
		Permissions: ifEmpty(u.Permissions, `{"view_only": false}`),
		ParentID:    u.ParentID,
		// 中文注释：若为子账号，附带当前登录令牌；用于服务端校验令牌是否已刷新
		LoginToken: lt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTExpireDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		common.Error(c, 50003, "签发令牌失败")
		return
	}
	// 中文注释：计算返回的金币值——若启用了父子金币同步且为子账号，则使用父账号金币
	coinsToReturn := u.Coins
	cfg2, _ := config.Load()
	if cfg2 != nil && cfg2.ParentCoinsSync && u.ParentID != nil {
		var parent models.User
		if err := db.DB().First(&parent, *u.ParentID).Error; err == nil {
			coinsToReturn = parent.Coins
		}
	}
	// 返回与前端一致的用户字段
	common.Ok(c, gin.H{
		"token": signed,
		"user": gin.H{
			"id":                   u.ID,
			"username":             u.Username,
			"nickname":             u.Nickname,
			"role":                 u.Role,
			"permissions":          u.Permissions,
			"parent_id":            u.ParentID,
			"coins":                coinsToReturn,
			"tomatoes":             u.Tomatoes,
			"avatar_path":          u.AvatarPath,
			"phone":                u.Phone,
			"email":                u.Email,
			"child_birthday":       u.ChildBirthday,
			"child_gender":         u.ChildGender,
			"must_change_password": u.MustChangePassword,
			// VIP 字段
			"last_login_time": now.Format(time.RFC3339),
			"is_vip":          u.IsVIP,
			"vip_expire_time": func() any {
				if u.VIPExpireTime != nil {
					return u.VIPExpireTime.Format(time.RFC3339)
				}
				return nil
			}(),
			"is_lifetime_vip": u.IsLifetimeVIP,
		},
	})
}

// 中文注释：空字符串兼容处理
func ifEmpty(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}

func isStrongPassword(pw string) bool {
	if len([]rune(pw)) < 6 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, r := range pw {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		} else if r >= 'a' && r <= 'z' {
			hasLower = true
		} else if r >= '0' && r <= '9' {
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func generateTempPassword() (string, error) {
	const digits = "0123456789"
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	out := make([]byte, 0, 8)
	out = append(out, 'R', 'g')
	for i := 0; i < len(b); i++ {
		out = append(out, digits[int(b[i])%len(digits)])
	}
	return fmt.Sprintf("%s", string(out)), nil
}
