package models

import "time"

// 中文注释：用户模型（简化版），后续可按 PRD 完善
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"uniqueIndex;size:64" json:"username"`
	PasswordSha string    `gorm:"size:128" json:"-"`
	Role        string    `gorm:"size:16" json:"role"` // admin/user
	Permissions string    `gorm:"type:text" json:"permissions"` // JSON
	// 中文注释：子账号登录令牌（随机生成，用于免密码登录子账号）
	LoginToken  string    `gorm:"size:128;index" json:"-"`
	Coins       int64     `json:"coins"`
	Tomatoes    int64     `json:"tomatoes"`
	Nickname    string    `gorm:"size:64" json:"nickname"`
	AvatarPath  string    `gorm:"size:256" json:"avatar_path"`
	Phone       string    `gorm:"size:32" json:"phone"`
	Email       string    `gorm:"size:128" json:"email"`
	ChildBirthday string  `gorm:"size:16" json:"child_birthday"`
	ChildGender   string  `gorm:"size:16" json:"child_gender"`
	ParentID    *uint     `json:"parent_id"`
	LoginTokenExpireAt *time.Time `json:"login_token_expire_at"`
	// 新增：最后登录时间
	LastLoginTime *time.Time `json:"last_login_time"`
	// 新增：VIP 状态与过期时间、终身 VIP 标记
	IsVIP          bool       `json:"is_vip"`
	VIPExpireTime  *time.Time `json:"vip_expire_time"`
	IsLifetimeVIP  bool       `json:"is_lifetime_vip"`
	IsDisabled     bool       `json:"is_disabled"`
	// 新增：是否必须修改密码（用于忘记密码临时密码登录后强制改密）
	MustChangePassword bool       `json:"must_change_password"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
