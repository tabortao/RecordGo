# RecordGo 项目骨架

> 中文注释：本仓库按照 `/docs/PRD/项目需求.md` 与 `.trae/rules/project_rules.md` 规范，提供前后端基础架构与核心占位模块，便于后续迭代完整需求。

## 目录结构

```
frontend/  # Vue 3 + Vite + TS + Pinia + Element Plus
backend/   # Gin + GORM + SQLite（纯 Go 驱动）
docs/      # 需求文档
```

## 前端运行（Windows）

```powershell
cd frontend
npm install
npm run dev
```

> 如使用 pnpm：`pnpm install && pnpm dev`

## 后端运行（Windows）

```powershell
cd backend
$env:SECRET_KEY = "replace-with-your-secret"
$env:DB_PATH = "storage/database/recordgo.db"
$env:STORAGE_ROOT = "storage"
$env:PORT = "8080"
$env:GIN_MODE = "release"
go mod tidy
go run ./cmd/server
```

## 关键规范落实

- 统一响应结构：`{code,message,data}`（后端 `common.Ok/Error`）
- 状态管理：前端 `Pinia` 的 `appState`，持久化到 `localStorage`
- 图片处理：前端提供 `utils/image.ts` 将上传图片转换为 `webp`，失败回退原格式
- 底部导航：固定为「任务 / 心愿 / 我的」，当前页面高亮

## 下一步建议

- 接入用户注册/登录与 JWT 鉴权
- 任务 CRUD 与番茄钟状态联动；愿望兑换与金币扣减
- 图片上传 API 与服务端保存到 `storage/uploads` 对应路径
- 前端缓存与 PWA 支持（`vite-plugin-pwa`）

