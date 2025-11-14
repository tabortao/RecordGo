## 目标与原则
- 不改动业务与数据逻辑，仅调整样式与配色，提升童趣与可读性
- 引入深色/浅色（夜间/白天）样式，默认跟随系统 `prefers-color-scheme`
- 统一使用 Tailwind 工具类与少量 `:deep()` 覆盖，保持现有组件结构

## 深浅色方案
- 保持 Tailwind 默认 `darkMode: 'media'`，通过 `dark:` 变体实现自动夜间模式
- 可选增强（若需要 Element Plus 同步深色）：
  - 在 `src/main.ts:4` 引入 `element-plus/theme-chalk/dark/css-vars.css`
  - 在应用初始化时根据 `matchMedia('(prefers-color-scheme: dark)')` 为 `document.documentElement` 添加/移除 `dark` 类，使 EP 使用深色变量
  - 不涉及业务，仅 DOM 类名切换；若不启用该增强，EP 组件将保持浅色，但页面容器会是深色

## 主题与配色（儿童友好）
- 扩展 Tailwind 主题色：
  - `brand`（主色/绿）：`#22c55e`（已存在 `primary`）
  - `sky`（辅助/蓝）：`#60a5fa`
  - `amber`（强调/橙）：`#f59e0b`
  - `pink`（奖励/粉）：`#ec4899`
  - `teal`（完成/青）：`#14b8a6`
  - `surface`：浅色 `#f9fafb`，深色 `#0f172a`（通过 `dark:` 应用）
- 统一圆角与阴影：组件外层采用 `rounded-xl`、`shadow-sm/hover:shadow-md`
- 字体与字号：标题 `text-lg font-semibold`，正文 `text-sm~base`，移动端触控尺寸不低于 `44px`

## 全局样式
- 在 `src/styles/index.css:2-4` 保持 Tailwind 导入；追加：
  - 基础背景与文字：`body` 浅色 `bg-gray-50 text-gray-900`，深色 `dark:bg-gray-900 dark:text-gray-100`
  - 平滑滚动与抗锯齿、输入在移动端保持 16px 最小字号（现有规则保留）

## 页面与组件美化清单
- `src/App.vue:3` 顶层容器：
  - `bg-gray-50` → `bg-gray-50 dark:bg-gray-900`
  - 统一文本色：`text-gray-900 dark:text-gray-100`
- `src/components/BottomNav.vue:3` 底部导航：
  - 背景与边框：`bg-white border-t` → `bg-white/90 dark:bg-gray-800/90 backdrop-blur border-t border-gray-200 dark:border-gray-700`
  - 文字：`text-gray-600` → `text-gray-600 dark:text-gray-300`
  - 高亮：保持 `text-primary`，并增加粗体与轻微缩放 `font-semibold`
- 顶部栏（固定导航头）：
  - `TasksPage.vue:8`、`WishesPage.vue:2`、`MinePage.vue:2`：`bg-white border-b` → `bg-white/80 dark:bg-gray-800/80 backdrop-blur border-b border-gray-200 dark:border-gray-700`
- 卡片与列表：
  - 任务卡片 `TasksPage.vue:153-161`：增加 `dark:border-gray-700 dark:hover:ring-blue-200/30`
  - 统计卡片文字：使用 Tailwind 图标色类替代内联样式：`text-emerald-500/600`、`text-sky-500`、`text-amber-600`、`text-teal-600`
  - 空状态与标签周边文字统一为 `text-gray-500 dark:text-gray-400`
- 心愿卡片 `WishesPage.vue:20-49`：
  - 卡片容器 `rounded-xl overflow-hidden shadow-sm hover:shadow-md`，边框加 `dark:border-gray-700`
  - 数值与按钮色使用 `amber/emerald` 系列，文字 `dark:` 适配
- 我的页分区卡 `MinePage.vue:21,51`：
  - 外层容器 `rounded-lg border bg-white` → `rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800`
  - 子项 hover：`hover:bg-gray-50 dark:hover:bg-gray-700`

## 图标与颜色一致性
- 将内联 `style="color:#xxxxxx"` 替换为 Tailwind 类：
  - 绿：`text-emerald-500/600`；蓝：`text-sky-500/600`；橙：`text-amber-500/600`；粉：`text-pink-500/600`；青：`text-teal-500/600`
- 状态文字：完成 `text-teal-600`、金币 `text-amber-600`、计划用时 `text-green-600`，并在深色下保持对比度

## 交互与可读性优化（仅样式）
- 增大点击热区与间距：主要按钮与列表项 `py-2 px-3`，移动端卡片间距 `gap-3`
- 使用 `backdrop-blur` 提升固定头/底栏在滚动时的层次感
- 保持任务卡片左侧“完成”圆按钮在深色中的边框对比度（未完成 `border-gray-400 dark:border-gray-500`）

## 验证与兼容
- 在 Windows 浏览器切换系统主题，检查三大页面与底部导航、统计卡片在深色/浅色下的对比与可读性
- 移动端浏览器（Chrome/Edge）模拟深色外观，确认固定头/底栏透明度与模糊效果
- 若启用 EP 深色增强，核对对话框/抽屉/表单在深色变量下的背景与边框可读性

## 改动范围与安全性
- 仅修改样式相关文件与模板类名：`tailwind.config.cjs`、`src/styles/index.css`、`src/App.vue`、`src/components/BottomNav.vue`、`src/pages/*`
- 保留现有逻辑与 API 交互；不改动 Pinia/服务层/路由逻辑

## 交付
- 提交一组小步变更，按页面依次美化与验证，确保热更新生效
- 如需手动主题切换，后续可在“设置”页追加开关，不影响本次自动深色实现