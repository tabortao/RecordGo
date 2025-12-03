**已实施**

- 在前端工程新增可执行脚本：
  - `lint`：`eslint .`
  - `lint:fix`：`eslint . --fix`
- 调整 ESLint 配置，确保仅对 `.vue` 文件启用 Vue 解析，并忽略 `dist` 与 `node_modules`，降低不必要的噪声并符合项目风格。

**改动文件**

- `frontend/package.json`：新增脚本
- `frontend/eslint.config.ts`：精简与约束规则范围，避免解析错误，放宽高噪声规则
  - 仅 `.vue` 启用 `vue-eslint-parser` 与 `eslint-plugin-vue`（`vue/multi-word-component-names` 关闭）
  - 忽略目录：`dist/**`、`node_modules/**`
  - 放宽规则：`@typescript-eslint/no-explicit-any`、`@typescript-eslint/no-empty-object-type`、`@typescript-eslint/no-unused-expressions`、`no-empty` 关闭；`@typescript-eslint/no-unused-vars` 设为 warn 并放行 `_` 开头的占位
  - `.d.ts` 文件关闭 `@typescript-eslint/triple-slash-reference`

**执行与结果**

- 执行检查：`pnpm run lint`（退出码 0）
- 自动修复：`pnpm run lint:fix`（退出码 0）
- 当前状态：仅剩 30 条 warning（主要是 `_` 或 `e` 未使用），不影响运行；已符合“不要过度设计”“尽量不修改其他模块代码”的要求。

**你可以直接使用**

- 检查：`pnpm run lint`
- 自动修复：`pnpm run lint:fix`

**说明与取舍**

- 我们将高噪声的规则（如 `no-explicit-any`、空块、空表达式）下调或关闭，以适配现有代码风格与工程约束，避免大量无意义变更。
- Vue 规则限定到 `.vue` 文件，修复了之前 `eslint-plugin-vue` 对非 Vue 文件（例如 `.jsonc`）误用导致的崩溃。
- 忽略 `dist` 防止针对构建产物的无效检查。

**后续建议**

- 如需严格化类型与未使用变量校验，可逐步将 `@typescript-eslint/no-unused-vars` 提升为 `error` 并在关键模块启用，而保持默认全局为 `warn`。
- 若你希望将 ESLint 脚本纳入固定流程，我可以把执行规范写入 `.trae/rules/project_rules.md`，并在每次修改后自动运行 `lint` 与 `lint:fix`。
