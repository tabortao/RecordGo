
https://juejin.cn/post/7480604184564465700

**部署目标**
- 将 `frontend` 作为纯前端静态站点部署到 Cloudflare Pages，前端通过 `VITE_API_BASE` 直连后端 `https://recordgo-server.a.top/api`。

**你需要准备**
- 一个 Cloudflare 账号与绑定的域名（可选）。
- 后端 API 能公网访问，且开启 CORS 允许你的前端域名。
- 源代码托管在 GitHub/GitLab（Cloudflare Pages 支持直接拉取构建）。

**关键点**
- 生产环境必须配置 `VITE_API_BASE`，否则构建后的前端会把接口当作静态路径请求，导致 404 或返回 HTML。
- 这是由 `src/services/http.ts:5` 的逻辑决定：生产环境会读取 `import.meta.env.VITE_API_BASE` 作为 `axios` 的 `baseURL`。
- 本项目是 History 模式的 SPA（`src/router/index.ts:8`），需要一个重写规则把任意路径回退到 `index.html`。我已把 `_redirects` 文件放进 Vite 的 `publicDir`，便于 Cloudflare Pages正确识别并生效。

**需要改代码吗**
- 不需要改业务代码。
- 我已在 `src/assets/favicon/_redirects` 添加了 SPA 回退规则，构建后会被复制到 `dist` 根目录，Cloudflare 会按此处理前端路由：
  - `frontend/src/assets/favicon/_redirects` 内容：`/* /index.html 200`

**Cloudflare Pages 部署步骤**
- 在 Cloudflare 控制台选择 Pages → Create a project。
- 选择 “Connect to Git”，绑定你的代码仓库。
- Build 设置：
  - Root directory：`frontend`
  - Build command：`npm run build`
  - Build output directory：`dist`
  - Node version：选择 18 或以上
- 环境变量（非常重要）：
  - 在 Pages → Settings → Environment variables 中添加
    - `VITE_API_BASE=https://recordgo-server.a.top/api`
  - 建议同时在 Production 与 Preview 环境都设置。
- 域名绑定（可选）：
  - Pages → Custom domains，绑定你的前端自定义域名（例如 `recordgo.a.top`）。
- 提交代码或触发构建，等待 Cloudflare 完成构建与发布。

**后端与跨域**
- 确保后端允许来自你的 Pages 域名的跨域请求（CORS）：
  - 允许 `Origin` 为你的前端域名（例如 `https://recordgo.a.top` 或 Pages 的默认子域）。
  - 允许 `Authorization` 头（前端会加 `Bearer <token>`）。
- 如果后端未配置 CORS，浏览器会在控制台报跨域错误。

**如何验证部署**
- 打开 Pages 生成的访问地址，进入 `/tasks`、`/wishes`、`/mine` 等路由，确认不刷新也能直接打开子路由（说明 SPA 回退生效）。
- 登录后在 Network 面板检查请求发往 `https://recordgo-server.a.top/api/...`，响应结构为 `{ code, message, data }`，且没有跨域错误。
- 如果看到错误“API响应异常：收到HTML内容”，说明没有正确设置 `VITE_API_BASE` 或后端代理混淆了静态与接口：
  - 参考 `src/services/http.ts:36`、`src/services/http.ts:97` 的错误提示，修正 `VITE_API_BASE`。

**关于 VITE_API_BASE 配置**
- 你可以填 `https://recordgo-server.a.top` 或 `https://recordgo-server.a.top/api` 均可，代码会自动补上 `/api`（`src/services/http.ts:6–9`）。
- 推荐用 Cloudflare Pages 的环境变量方式设置，不必在仓库里新增 `.env.production` 文件；后续改动域名只需在 Cloudflare 控制台调整即可。

**常见问题与排查**
- 前端 404：检查构建产物 `dist` 根目录是否包含 `_redirects`；Pages 项目路由应回退到 `index.html`。
- 接口 404 或返回 HTML：检查 `VITE_API_BASE` 是否设置到 Pages 的环境变量；确认后端 `/api/*` 路由对外可访问。
- CORS 报错：确认后端允许你的前端域名；浏览器控制台的报错会显示被拦截的请求与原因。

**代码位置参考**
- 生产环境接口基址读取：`frontend/src/services/http.ts:5`
- 自动补 `/api` 的处理：`frontend/src/services/http.ts:6-9`
- SPA 路由模式：`frontend/src/router/index.ts:8`
- Cloudflare Pages 路由回退规则文件：`frontend/src/assets/favicon/_redirects`

如果你愿意，我可以继续为你在 Cloudflare Pages 中创建项目的具体表单填写清单、截图步骤，或帮你把 `.env.production` 方案也准备好作为备选。