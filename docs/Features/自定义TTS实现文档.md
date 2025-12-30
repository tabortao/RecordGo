# 自定义 TTS 功能实现文档

## 1. 功能概述
本项目新增了自定义 TTS（Text-to-Speech）功能，允许用户接入第三方语音合成接口（如 OpenAI TTS 或兼容接口），以获得比浏览器内置语音更优质的朗读体验。

## 2. 架构设计

### 2.1 状态管理 (Pinia)
在 `frontend/src/stores/appState.ts` 中扩展了 `SpeechSettings` 接口：
- `engine`: `'system' | 'custom'` - 切换 TTS 引擎。
- `customProfiles`: `CustomTTSProfile[]` - 存储用户配置的接口列表。
- `activeCustomId`: `string | null` - 当前选中的配置 ID。

### 2.2 核心逻辑
- **frontend/src/utils/customTTS.ts**: 封装了 `createCustomAudio` 函数，负责根据配置发起 HTTP 请求获取音频 Blob，并创建 `HTMLAudioElement`。支持 GET 和 POST 请求（OpenAI 格式）。
- **frontend/src/utils/speech.ts**: 改造了全局 `speak` 函数为异步函数 (`Promise<boolean>`)。
  - 调用 `speak` 时，优先检查 `appState` 中的 `engine` 设置。
  - 若启用自定义 TTS 且配置有效，则调用 `createCustomAudio` 并播放音频。
  - 若失败或未启用，自动降级回退到浏览器原生 `speechSynthesis`。
  - 增加了 `currentAudio` 引用，确保调用 `stop()` 时能立即停止正在播放的音频。

### 2.3 界面实现
- **frontend/src/pages/SettingsReadingPage.vue**:
  - 重构为双 Tab 布局（系统默认 / 自定义接口）。
  - **自定义接口 Tab**: 提供配置的新增、编辑、删除、选择功能。
  - **测试连接**: 提供“测试连接”按钮，智能猜测接口格式（GET/POST）并验证可用性。
  - **持久化**: 配置数据通过 Pinia 自动持久化到 `localStorage`。

## 3. 使用说明

1. 进入 **我的 -> 朗读设置**。
2. 切换到 **自定义接口** Tab。
3. 点击 **新增配置**，填写：
   - **配置名称**: 任意备注。
   - **接口地址**: 例如 `https://api.openai.com/v1/audio/speech` 或其他兼容 API。
   - **API Key**: 若接口需要鉴权则填写。
   - **音色ID**: 第三方接口所需的 voice 参数（如 `alloy`）。
4. 点击 **测试连接** 验证。
5. 选中该配置，系统即刻生效。

## 4. 兼容性与降级
- **并行运行**: 系统 TTS 设置与自定义 TTS 设置独立存储，切换引擎即可无缝切换。
- **全局生效**: 任务详情页、任务列表页的朗读功能均已适配异步调用，支持自定义 TTS。
- **降级策略**: 若自定义接口请求失败（网络错误、Key 无效等），`speak` 函数会捕获异常并静默失败（或降级，目前设计为返回 false 由调用方提示），不会导致应用崩溃。

## 5. 测试
- **单元测试**: `frontend/src/__tests__/customTTS.spec.ts` 覆盖了 GET/POST 请求构建、参数传递、错误处理等核心逻辑。
- **运行测试**: `cd frontend; npx vitest run src/__tests__/customTTS.spec.ts`

