package handlers

import (
    "net/http"
    "path/filepath"
    "strconv"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
    "recordgo/internal/common"
    "recordgo/internal/config"
    s "recordgo/internal/storage"
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
    if st != nil { drv = st.DriverName() }
    cfg, _ := config.Load()
    if cfg != nil { pub = cfg.S3PublicBaseURL }
    common.Ok(c, gin.H{"driver": drv, "public_base_url": pub})
}

func StoragePresign(c *gin.Context) {
    var req presignReq
    if err := c.ShouldBindJSON(&req); err != nil { common.Error(c, 40000, "参数错误"); return }
    if req.UserID == 0 || strings.TrimSpace(req.ResourceType) == "" || strings.TrimSpace(req.ContentType) == "" { common.Error(c, 40000, "参数错误"); return }
    key, err := buildKey(req)
    if err != nil { common.Error(c, 40000, "参数错误"); return }
    st := s.Get()
    if st == nil { common.Error(c, 50010, "存储未初始化"); return }
    url, headers, err := st.PresignPut(c.Request.Context(), key, req.ContentType, 300)
    if err != nil { common.Error(c, 50011, "生成上传地址失败"); return }
    common.Ok(c, gin.H{"upload_url": url, "object_key": key, "headers": headers, "expires": 300})
}

func StoragePresignView(c *gin.Context) {
    key := strings.TrimSpace(c.Query("object_key"))
    if key == "" { common.Error(c, 40000, "参数错误"); return }
    st := s.Get()
    if st == nil { common.Error(c, 50010, "存储未初始化"); return }
    url, err := st.PresignGet(c.Request.Context(), key, 300)
    if err != nil { common.Error(c, 50012, "生成访问地址失败"); return }
    common.Ok(c, gin.H{"url": url})
}

func StoragePut(c *gin.Context) {
    key := strings.TrimSpace(c.Query("key"))
    if key == "" { c.JSON(http.StatusBadRequest, gin.H{"code":400,"message":"参数错误","data":nil}); return }
    st := s.Get()
    if st == nil { c.JSON(http.StatusInternalServerError, gin.H{"code":500,"message":"存储未初始化","data":nil}); return }
    if st.DriverName() != "local" { c.JSON(http.StatusMethodNotAllowed, gin.H{"code":405,"message":"不支持","data":nil}); return }
    ct := c.GetHeader("Content-Type")
    r := c.Request.Body
    defer r.Close()
    if err := st.Put(c.Request.Context(), key, r, ct); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":500,"message":"保存失败","data":nil}); return }
    c.JSON(http.StatusOK, gin.H{"code":0,"message":"ok","data":gin.H{"path": filepath.ToSlash("uploads/"+key)}})
}

func buildKey(req presignReq) (string, error) {
    ts := time.Now().Unix()
    ns := time.Now().UnixNano()
    ext := strings.TrimLeft(strings.ToLower(req.Ext), ".")
    if ext == "" { ext = "webp" }
    // 可选路径前缀（例如：blog），由环境变量 S3_KEY_PREFIX 控制
    cfg, _ := config.Load()
    prefix := strings.Trim(strings.TrimSpace(cfg.S3KeyPrefix), "/")
    if prefix != "" { prefix = prefix + "/" }
    switch req.ResourceType {
    case "avatar":
        return filepath.ToSlash(prefix+"images/avatars/"+strconv.Itoa(int(req.UserID))+"/avatar_"+strconv.Itoa(int(req.UserID))+"_"+strconv.FormatInt(ts,10)+"_"+strconv.FormatInt(ns,10)+"."+ext), nil
    case "task_image":
        if req.TaskID == nil || *req.TaskID == 0 { return "", nil }
        return filepath.ToSlash(prefix+"images/task_images/"+strconv.Itoa(int(req.UserID))+"/task_"+strconv.Itoa(int(*req.TaskID))+"_"+strconv.FormatInt(ts,10)+"_"+strconv.FormatInt(ns,10)+"."+ext), nil
    case "task_attachment_img":
        if req.TaskID == nil || *req.TaskID == 0 { return "", nil }
        return filepath.ToSlash(prefix+"task_attachments/"+strconv.Itoa(int(req.UserID))+"/"+strconv.Itoa(int(*req.TaskID))+"/img_"+strconv.Itoa(int(req.UserID))+"_"+strconv.Itoa(int(*req.TaskID))+"_"+strconv.FormatInt(ts,10)+"_"+strconv.FormatInt(ns,10)+"."+ext), nil
    case "task_attachment_audio":
        if req.TaskID == nil || *req.TaskID == 0 { return "", nil }
        aext := ext
        if aext == "" { aext = "mp3" }
        return filepath.ToSlash(prefix+"task_attachments/"+strconv.Itoa(int(req.UserID))+"/"+strconv.Itoa(int(*req.TaskID))+"/audio_"+strconv.Itoa(int(req.UserID))+"_"+strconv.Itoa(int(*req.TaskID))+"_"+strconv.FormatInt(ts,10)+"_"+strconv.FormatInt(ns,10)+"."+aext), nil
    case "wish_image":
        return filepath.ToSlash(prefix+"images/wish/"+strconv.Itoa(int(req.UserID))+"/wish_"+strconv.Itoa(int(req.UserID))+"_"+strconv.FormatInt(ts,10)+"_"+strconv.FormatInt(ns,10)+"."+ext), nil
    default:
        return "", nil
    }
}