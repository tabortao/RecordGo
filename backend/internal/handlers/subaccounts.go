package handlers

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "strings"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "recordgo/internal/common"
    "recordgo/internal/db"
    "recordgo/internal/models"
)

// 中文注释：子账号响应体（避免暴露密码等字段）
type childResp struct {
    ID         uint   `json:"id"`
    Nickname   string `json:"nickname"`
    AvatarPath string `json:"avatar_path"`
    Permissions string `json:"permissions"`
    LoginToken string `json:"login_token"`
}

func toChildResp(u *models.User) childResp {
    return childResp{ID: u.ID, Nickname: u.Nickname, AvatarPath: u.AvatarPath, Permissions: strings.TrimSpace(u.Permissions), LoginToken: strings.TrimSpace(u.LoginToken)}
}

// ListSubAccounts 列出当前家长（或拥有权限的子账号）下的子账号
func ListSubAccounts(c *gin.Context) {
    if !hasPermission(c, "account", "manage_children") {
        deny(c, "无权限查看子账号")
        return
    }
    // 取 claims：父账号使用自己的 ID；子账号使用 claims.ParentID
    cl := extractClaims(c)
    if cl == nil {
        common.Error(c, 40100, "未登录或令牌无效")
        return
    }
    var parentID uint
    if cl.ParentID == nil { parentID = cl.UserID } else { parentID = *cl.ParentID }
    var users []models.User
    if err := db.DB().Where("parent_id = ?", parentID).Find(&users).Error; err != nil {
        common.Error(c, 50050, "查询子账号失败")
        return
    }
    res := make([]childResp, 0, len(users))
    for i := range users { res = append(res, toChildResp(&users[i])) }
    common.Ok(c, res)
}

// CreateChildReq 创建子账号请求
type CreateChildReq struct {
    Nickname    string `json:"nickname"`
    Permissions string `json:"permissions"` // JSON 字符串
    AvatarPath  string `json:"avatar_path"` // 可选，上传后设置
}

// CreateSubAccount 创建子账号（用户名自动生成，密码为空，LoginToken 初始为空）
func CreateSubAccount(c *gin.Context) {
    if !hasPermission(c, "account", "manage_children") {
        deny(c, "无权限创建子账号")
        return
    }
    var req CreateChildReq
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Nickname) == "" {
        common.Error(c, 40001, "参数错误：昵称必填")
        return
    }
    cl := extractClaims(c)
    if cl == nil { common.Error(c, 40100, "未登录或令牌无效"); return }
    var parentID uint
    if cl.ParentID == nil { parentID = cl.UserID } else { parentID = *cl.ParentID }
    // 生成唯一用户名（child_父ID_随机）
    uname := fmt.Sprintf("child_%d_%s", parentID, randHex(8))
    u := &models.User{
        Username:    uname,
        Role:        "user",
        Permissions: strings.TrimSpace(req.Permissions),
        Nickname:    strings.TrimSpace(req.Nickname),
        AvatarPath:  strings.TrimSpace(req.AvatarPath),
        ParentID:    &parentID,
    }
    if err := db.DB().Create(u).Error; err != nil {
        common.Error(c, 50051, "创建子账号失败")
        return
    }
    common.Ok(c, toChildResp(u))
}

// UpdateChildReq 更新子账号请求
type UpdateChildReq struct {
    Nickname    *string `json:"nickname"`
    Permissions *string `json:"permissions"` // JSON 字符串
    AvatarPath  *string `json:"avatar_path"` // 可选，头像上传后回填路径
}

func UpdateSubAccount(c *gin.Context) {
    if !hasPermission(c, "account", "manage_children") {
        deny(c, "无权限修改子账号")
        return
    }
    id := strings.TrimSpace(c.Param("id"))
    var child models.User
    if err := db.DB().First(&child, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            common.Error(c, 40404, "子账号不存在")
        } else {
            common.Error(c, 50052, "查询子账号失败")
        }
        return
    }
    var req UpdateChildReq
    if err := c.ShouldBindJSON(&req); err != nil {
        common.Error(c, 40001, "参数错误")
        return
    }
    if req.Nickname != nil { child.Nickname = strings.TrimSpace(*req.Nickname) }
    if req.Permissions != nil { child.Permissions = strings.TrimSpace(*req.Permissions) }
    if req.AvatarPath != nil { child.AvatarPath = strings.TrimSpace(*req.AvatarPath) }
    if err := db.DB().Save(&child).Error; err != nil {
        common.Error(c, 50053, "更新子账号失败")
        return
    }
    common.Ok(c, toChildResp(&child))
}

func DeleteSubAccount(c *gin.Context) {
    if !hasPermission(c, "account", "manage_children") {
        deny(c, "无权限删除子账号")
        return
    }
    id := strings.TrimSpace(c.Param("id"))
    var child models.User
    if err := db.DB().First(&child, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            common.Error(c, 40404, "子账号不存在")
        } else {
            common.Error(c, 50054, "查询子账号失败")
        }
        return
    }
    if err := db.DB().Delete(&child).Error; err != nil {
        common.Error(c, 50055, "删除子账号失败")
        return
    }
    common.Ok(c, gin.H{"id": child.ID})
}

// GenerateChildToken 生成或刷新子账号登录令牌
func GenerateChildToken(c *gin.Context) {
    if !hasPermission(c, "account", "manage_children") {
        deny(c, "无权限生成令牌")
        return
    }
    id := strings.TrimSpace(c.Param("id"))
    var child models.User
    if err := db.DB().First(&child, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            common.Error(c, 40404, "子账号不存在")
        } else {
            common.Error(c, 50056, "查询子账号失败")
        }
        return
    }
    token := randHex(32)
    child.LoginToken = token
    if err := db.DB().Save(&child).Error; err != nil {
        common.Error(c, 50057, "保存令牌失败")
        return
    }
    common.Ok(c, gin.H{"token": token})
}

func randHex(n int) string {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil { return "" }
    return hex.EncodeToString(b)
}