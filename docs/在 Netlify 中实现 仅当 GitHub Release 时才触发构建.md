# 在 Netlify 中实现 **仅当 GitHub Release 时才触发构建**

在 Netlify 中实现 **仅当 GitHub Release 时才触发构建**，而默认关闭自动构建的配置方法如下：

---

### 1. **禁用 Netlify 默认的自动构建**
- 进入 Netlify 站点管理界面 → `Site settings` → `Build & deploy` → `Continuous Deployment`。
- 取消勾选 **`Build on push`**（或找到对应的 GitHub 仓库禁用“自动构建”）。
- 这将阻止 Netlify 在每次 `git push` 时触发构建。

---

### 2. **使用 GitHub Release 触发构建**
需要通过 **GitHub Webhook** 或 **Netlify Build Hooks** 实现：

#### 方法一：通过 **Netlify Build Hook**（推荐）
1. **创建 Build Hook**：
   - 在 Netlify 的 `Site settings` → `Build & deploy` → `Build hooks` 中点击 **`Add build hook`**。
   - 命名为 `Release Trigger`，选择分支（如 `main`），生成 Hook URL。

2. **配置 GitHub Release 触发 Webhook**：
   - 进入你的 GitHub 仓库 → `Settings` → `Webhooks` → `Add webhook`。
   - 填写以下参数：
     - **Payload URL**: 粘贴 Netlify 提供的 Build Hook URL。
     - **Content type**: `application/json`。
     - **Secret**: 可留空（Netlify Build Hook 无需密钥）。

3.**选择触发事件**：

   - 在 **“Which events would you like to trigger this webhook?”** 部分：
     - ❌ 不要选择 `Just the push event`（这是默认的 `push` 触发）。
     - ✅ 选择 **`Let me select individual events`**，然后勾选：
       - `Releases`（关键选项）

4. 点击 **Add webhook** 保存。
   

#### 方法二：通过 **GitHub Actions**（灵活控制）
1. 在项目中创建 `.github/workflows/netlify-release.yml` 文件：
   ```yaml
   name: Trigger Netlify Build on Release
   on:
     release:
       types: [published]
    
   jobs:
     deploy:
       runs-on: ubuntu-latest
       steps:
         - name: Trigger Netlify Build
           run: |
             curl -X POST -d {} "YOUR_NETLIFY_BUILD_HOOK_URL"
   ```
   - 替换 `YOUR_NETLIFY_BUILD_HOOK_URL` 为实际的 Build Hook URL。

---

### 3. **验证配置**

1. 在 GitHub 上创建一个测试 Release：
   - 点击 `Releases` → `Draft a new release` → 发布。
2. 观察效果：
   - Netlify 的构建列表（`Deploys`）应出现新构建。
   - 在 GitHub Webhook 设置的 **Recent Deliveries** 中，能看到 Release 事件已触发 Webhook。

这样即可精准控制 Netlify 仅在 GitHub Release 时构建，避免不必要的资源消耗。