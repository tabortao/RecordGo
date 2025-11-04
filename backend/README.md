# RecordGo 后端

> 中文注释：Gin + GORM + SQLite 后端服务基础骨架，统一响应、健康检查与默认心愿占位接口。

## 环境变量（Windows PowerShell 示例）

```powershell
$env:SECRET_KEY = "replace-with-your-secret"      # JWT 密钥
$env:DB_PATH = "storage/database/recordgo.db"     # 数据库路径
$env:STORAGE_ROOT = "storage"                     # 存储根目录
$env:PORT = "8080"                                # 服务端口
$env:GIN_MODE = "release"                         # 或 "debug"
```

## 运行

```powershell
cd backend
go mod tidy
go run ./cmd/server
```

## 接口

- `GET /api/health` 健康检查，返回 `{code:0,message:"ok",data:{status:"healthy"}}`
- `GET /api/wishes/default` 默认心愿占位数据

