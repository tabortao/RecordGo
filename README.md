# RecordGo - 任务与积分管理小程序

RecordGo 是一个面向学生的任务与积分管理工具，旨在通过简洁的界面与激励机制，帮助用户培养良好的学习与时间管理习惯。

## ✨ 项目特性

- **任务管理**：支持创建、编辑、完成任务，设置重复规则（每天/每周/每月等）。
- **积分激励**：完成任务获得金币，未完成任务扣除金币（支持负分）。
- **心愿商店**：使用赚取的金币兑换心愿奖励。
- **番茄钟**：内置番茄钟专注工具，支持悬浮球、全屏模式，自动记录专注时长。
- **多账号支持**：支持家长/主账号与子账号体系，积分池共享，权限可控。
- **数据统计**：提供详细的任务完成率、金币收支、专注时长等统计图表。
- **跨平台支持**：响应式设计，支持 PWA 离线访问。

## 🛠️ 技术栈

### 前端 (Frontend)
- **核心框架**: Vue 3 + Vite + TypeScript
- **UI 组件库**: Element Plus
- **状态管理**: Pinia
- **样式库**: Tailwind CSS
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **其他**: Day.js, ECharts, PWA

### 后端 (Backend)
- **Web 框架**: Gin (Go 1.24.1)
- **ORM**: GORM
- **数据库**: SQLite (纯 Go 驱动 `glebarez/sqlite`)
- **认证**: JWT
- **配置管理**: Viper
- **日志**: Zap

## 📂 目录结构

```
RecordGo/
├── frontend/           # 前端 Vue 3 项目源码
├── backend/            # 后端 Go 项目源码
├── docker/             # Docker 部署配置
├── docs/               # 项目文档与 PRD
└── storage/            # (运行时生成) 数据存储目录
    ├── database/       # SQLite 数据库文件
    └── uploads/        # 用户上传的图片资源
```

## 🚀 快速开始 (Windows)

### 前置要求
- Node.js ≥ 18 LTS
- Go ≥ 1.24
- Git
- (可选) Docker Desktop

### 1. 后端启动

```powershell
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 设置环境变量 (开发环境示例)
$env:SECRET_KEY = "dev-secret"
$env:DB_PATH = "../storage/database/recordgo.db"
$env:STORAGE_ROOT = "../storage"
$env:PORT = "8082"
$env:GIN_MODE = "debug"

# 运行服务
go run ./cmd/server
```

服务默认运行在 `http://localhost:8082`。

### 2. 前端启动

```powershell
# 进入前端目录
cd frontend

# 安装依赖
pnpm install
# 或 npm install

# 启动开发服务器
pnpm dev
# 或 npm run dev
```

前端默认运行在 `http://localhost:5173`。
请确保前端 `.env.development` 中的 `VITE_API_BASE` 配置指向正确的后端地址。

### 3. Docker 部署 (可选)

项目提供了 Docker Compose 配置，可快速启动完整环境。

```powershell
# 进入 docker 目录
cd docker

# 启动服务
docker-compose up -d
```

## ⚙️ 环境配置

### 后端环境变量
| 变量名 | 说明 | 示例 |
|--------|------|------|
| `PORT` | 服务监听端口 | `8082` |
| `DB_PATH` | SQLite 数据库路径 | `storage/database/recordgo.db` |
| `STORAGE_ROOT` | 文件存储根目录 | `storage` |
| `SECRET_KEY` | JWT 签名密钥 | `your-secret-key` |
| `GIN_MODE` | Gin 运行模式 | `debug` / `release` |

### 前端环境变量
| 变量名 | 说明 |
|--------|------|
| `VITE_API_BASE` | 后端 API 接口地址 |

## 📄 许可证

本项目采用 MIT 许可证。
