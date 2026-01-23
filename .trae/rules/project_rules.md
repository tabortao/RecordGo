# 任务与积分小程序项目规则（RecordGo）

版本：0.1.1  
更新日期：2025-11-04  
适用平台：Windows 10/11  
技术栈：前端 Vue 3 + Vite + Element Plus + Pinia；后端 Gin（Go 1.24.1）+ GORM；数据库 SQLite（纯 Go 驱动）；支持 Docker 部署

## 最新规则，必须遵守

- 前端具有热加载，修改无需刷新即可生效，无需再次启动前端，前端使用`pnpm dev` ，默认前端已经启动在 `http://localhost:5173`
- 请记住 不要pnpm build，项目已经启动了 pnpm dev，不需要再次启动前端了，具有热加载。
- 代码编写需满足`.trae\documents\常见错误.md`中的要求。
- 文档、注释使用中文编写。
- 前端使用pnpm代替npm。
- 每次执行任务，先在.trae\plans文件夹中按日期+任务名，创建好本次任务的计划，然后按照计划逐步执行任务并消项。

---

## 目录

1. 项目简介与架构约定  
2. 安装与执行步骤（Windows）  
3. 必须遵守的技术规范与约束  
4. 环境配置要求与依赖清单  
5. 执行流程中的关键检查点与验证标准  
6. 错误处理与调试指南  
附录 A：目录与路径约定  
附录 B：环境变量示例（Windows PowerShell）  
附录 C：接口一致性检查清单

---

## 1. 项目简介与架构约定

- 项目目标：面向学生的任务与积分管理工具，通过简洁界面与激励机制，帮助培养学习与时间管理习惯。
- 架构：前后端分离，RESTful API；前端缓存与 PWA 离线支持；支持 Docker 部署。
- 页面结构（底部导航三页）：任务、心愿、我的。
- 数据库：`storage/database/recordgo.db`
- 前端强约束：所有字段必须与后端接口保持一致；上传图片需前端压缩并转为 `webp`。

---

## 2. 安装与执行步骤（Windows）

### 2.1 前置条件

- 已安装 Node.js ≥ 18 LTS（推荐）与包管理器（pnpm ≥ 8 或 npm/yarn）
- 已安装 Go 1.24.1（64 位）
- 可选：Docker Desktop（用于容器化运行）
- Git（拉取代码与版本管理）
- 项目会推送到github，使用cloudflare部署，前端地址是https://recordgo.sdgarden.top，后端地址是https://recordgo-server.sdgarden.top

### 2.2 前端（Vue 3）安装与运行

```powershell
# 进入前端目录（示例：frontend）
cd frontend

# 安装依赖（推荐使用 pnpm）
pnpm install

# 开发模式启动（默认端口 5173）
pnpm dev

# 生产构建与本地预览
pnpm build
pnpm preview
```

说明：
- 如尚未生成前端工程，请按技术栈创建 `Vue 3 + Vite + TypeScript` 工程，并集成 `Element Plus`、`Pinia`、`Vue Router`、`Axios`、`Tailwind CSS`、`unplugin-auto-import`、`unplugin-vue-components`。
- 本地开发时可通过 `.env.development` 配置后端 API 地址（例如 `VITE_API_BASE=http://localhost:8080`）。

### 2.3 后端（Gin + GORM + SQLite）安装与运行

```powershell
# 进入后端目录（示例：backend）
cd backend

# 初始化/同步依赖
go mod tidy

# 设置必要环境变量（PowerShell 示例）
$env:SECRET_KEY = "replace-with-your-secret"   # JWT 加密密钥
$env:DB_PATH = "storage/database/recordgo.db"  # 数据库文件路径
$env:STORAGE_ROOT = "storage"                  # 存储根目录
$env:PORT = "8082"                             # 服务端口
$env:GIN_MODE = "release"                      # 或 "debug"

# 运行数据库迁移（如项目包含迁移命令）
# go run ./cmd/migrate

# 启动服务（示例）

cd g:\Code\RecordGo\backend ; $env:SECRET_KEY = 'dev-secret'; $env:DB_PATH = 'storage/database/recordgo.db'; $env:STORAGE_ROOT = 'storage'; $env:PORT = '8082'; $env:GIN_MODE = 'debug'; go run ./cmd/server

go run ./cmd/server

# 或构建可执行文件后运行
go build -o .\bin\server.exe ./cmd/server
./bin/server.exe
```

说明：
- SQLite 驱动需使用纯 Go 版本（推荐 GORM 适配器：`github.com/glebarez/sqlite`）。
- 首次运行会自动创建数据库与必要目录（如代码已实现自动化）。

---

## 3. 必须遵守的技术规范与约束

- 一致性与状态管理：
  - 前端状态统一使用 Pinia，维护 `appState`（番茄钟、权限、缓存、统计等）。
  - 前端所有页面数据与后端字段保持一致，严禁前端私自更改字段名称或结构。
  - 子账号与主账号金币池共享，权限通过 `permissions` JSON 控制（例如：`{"view_only": false}`）。
- 页面与交互：
  - 底部导航固定为「任务 / 心愿 / 我的」，当前页高亮。
  - 图片上传统一前端压缩并转换为 `webp`，失败需有回退与失败提示。
  - 番茄钟支持悬浮球、固定全屏、状态持久化，设置保存到 `UserSettings`。
- 接口与数据：
  - API 采用 RESTful 设计，响应统一为 `{ code, message, data }` 结构。
  - 关键接口添加超时与错误重试策略；大数据列表采用分页（如操作记录、兑换记录）。
  - 任务状态：0-待完成，1-进行中，2-已完成；支持负金币（惩罚）。
- 代码规范：
  - 前端 TypeScript 严格模式；统一模块化与常量化；逻辑可复用抽象在 `composables/`、`services/`、`stores/`。
  - 后端使用分层架构（router → handler → service → repository），统一错误处理与日志（Zap）。
  - 所有新增代码必须包含中文注释，保证可读性与维护性。

---

## 4. 环境配置要求与依赖清单

### 4.1 前端主要依赖

- 运行时：`vue@3`、`vue-router@4`、`pinia`、`element-plus`
- 网络与工具：`axios`、`@vueuse/core`、`dayjs`、`lodash-es`
- 构建与工程：`vite`、`typescript`、`unplugin-auto-import`、`unplugin-vue-components`
- 样式：`tailwindcss`
- PWA：`vite-plugin-pwa`（启用离线支持）

### 4.2 后端主要依赖

- Web 框架：`github.com/gin-gonic/gin`
- ORM：`gorm.io/gorm` + 纯 Go SQLite 驱动 `github.com/glebarez/sqlite`
- 配置：`github.com/spf13/viper`
- 认证：`github.com/golang-jwt/jwt/v5`
- 权限（可选）：`github.com/casbin/casbin/v2`
- 日志：`go.uber.org/zap`
- 数据库迁移：`github.com/golang-migrate/migrate/v4`

### 4.3 运行环境要求

- Windows 10/11，PowerShell 7（推荐）
- Node.js ≥ 18 LTS，Go 1.24.1
- 开发网络：可访问本地 `http://localhost:8080`（后端）与 `http://localhost:5173`（前端）

---

## 5. 执行流程中的关键检查点与验证标准

- 目录与路径：确认 `storage/` 目录存在且可写。数据库文件位于 `storage/database/recordgo.db`。
- 接口健康检查：后端提供 `GET /api/health`（或等价）返回 `code=0`。
- 认证与权限：
  - 登录注册流程可用；JWT 正常签发与校验；子账号返回包含 `parent_id` 与 `permissions`。
  - 前端根据 `permissions.view_only` 动态控制按钮显示与禁用逻辑。
- 任务页验证：
  - 顶部统计（日时长、任务数、日金币、完成率、统计图标）随任务变更自动更新。
  - 日期与周视图切换正确；今日回到“本周”按钮可用；任务筛选（全部/已完成/待完成）正确。
- 番茄钟：
  - 计时模式可切换（倒计时/正计时）；预设时间与自定义可用；悬浮球动画正常；完成后自动标记任务完成并写入实际耗时。
  - `fixed_tomato_page` 设置生效并持久化到 `UserSettings`。
- 心愿页：
  - 默认 6 个心愿显示；兑换流程扣减金币并增加兑换次数；兑换记录可查看。
  - 图片上传为 `webp`，命名格式 `wish_用户ID_时间戳_uuid.webp`。
- 我的页：
  - 昵称/ID/头像显示正确；荣誉墙交互符合要求；数据导出/清除/操作记录分页可用。

---

## 6. 错误处理与调试指南

### 6.1 前端

- Axios 拦截器：统一处理网络错误、超时与业务错误码；弹出友好提示并记录日志。
- 重试与超时：关键请求设置合理超时与指数退避重试（最多 3 次）。
- 本地缓存：优先读取缓存，失败时回退到 API；注意缓存失效策略。
- 图片处理：上传失败时提示并允许重试；处理 `webp` 转换失败的回退（保留原格式或重新压缩）。
- 调试建议：
  - 检查 `.env` 中 `VITE_API_BASE` 是否正确。
  - 打开 DevTools Network 观察请求与响应结构 `{ code, message, data }`。
  - 排查权限控制逻辑（按钮禁用/隐藏与路由守卫）。

### 6.2 后端

- 统一错误响应：
  - 成功：`{ code: 0, message: "ok", data: ... }`
  - 失败：`{ code: 非0, message: 错误描述, data: null }`
- 日志与链路：使用 Zap 记录关键操作与错误；输出包含请求路径、用户 ID、子账号标识、原因等。
- 数据迁移与清理：
  - 自动迁移（如启用）遇到冲突需写入操作日志并回滚或标记跳过。
  - 操作记录只保留近一个月，定期清理。
- 调试建议：
  - 运行时设置 `GIN_MODE=debug` 并打印 SQL（谨慎在生产开启）。
  - 使用 Postman/Bruno 逐个验证注册、登录、任务、心愿、兑换、操作记录等 API。
  - 检查 SQLite 文件权限与路径；确认使用纯 Go 驱动（避免 cgo 依赖）。

---

## 附录 A：目录与路径约定

- `storage/database/recordgo.db`：SQLite 数据库文件。

---

## 附录 B：环境变量示例（Windows PowerShell）

```powershell
$env:SECRET_KEY = "replace-with-your-secret"      # JWT 密钥
$env:DB_PATH = "storage/database/recordgo.db"     # 数据库路径
$env:STORAGE_ROOT = "storage"                     # 存储根目录
$env:PORT = "8080"                                # 服务端口
$env:GIN_MODE = "release"                         # 或 "debug"
$env:PARENT_COINS_SYNC = "true"                   # 子账号金币与父账号同步

# 前端（示例）
$env:VITE_API_BASE = "http://localhost:8080"      # 后端 API 地址
```

---

## 附录 C：接口一致性检查清单

- 字段与命名：确保前端字段与后端响应完全一致（大小写、类型、必填）。
- 状态枚举：任务状态使用 0/1/2；重复设置枚举与后端一致（无、每天、每周、每月、自定义星期）。
- 权限：`permissions.view_only` 控制 UI 行为；子账号返回包含 `parent_id`。
- 图片：上传前压缩为 `webp`；命名遵循规范；失败有回退与提示。
- 分页与查询：大数据列表使用分页；提供时间范围过滤（操作记录）。
- 统计与缓存：统计数据在前端缓存并在页面切换/任务变更时刷新；调用 API 校验实时性。

---

> 说明：本规则文件基于项目需求文档提炼，旨在确保团队成员以统一规范完成安装、开发与运行。执行时请以本文件为准，严格落实「状态管理与数据一致性、拆分与复用逻辑、代码注释与可读性、常量化与模块化、数据管理层抽象」等要求。

---

## 7. 代码检查与类型验证流程（前端）

- 固定流程：每次修改后统一执行以下命令，保证语法与类型安全。
- 在前端目录执行：

```powershell
pnpm run lint && pnpm run lint:fix && pnpm run typecheck
```

- 脚本说明：
  - `lint`：执行 ESLint 检查。
  - `lint:fix`：自动修复可修复问题。
  - `typecheck`：TypeScript 严格类型检查（`tsc --noEmit`）。

- 注意：前端采用热加载（`pnpm dev`），无需重启服务；仅在必要时关注控制台输出的警告与错误。
