# 任务计划：更新支持页面图片

日期：2026-02-06

## 任务目标
将 `SupportPage.vue` 中的远程图片链接替换为本地资源文件。

## 执行步骤
1. [x] 检查本地图片文件是否存在
   - 检查 `frontend/src/assets/images/wechatpay.jpg`
   - 检查 `frontend/src/assets/images/alipay.jpg`
2. [x] 修改 `SupportPage.vue`
   - 将 `https://img.sdgarden.top/blog/wechat/wechatpay.jpg` 替换为 `@/assets/images/wechatpay.jpg`
   - 将 `https://img.sdgarden.top/blog/wechat/alipay.jpg` 替换为 `@/assets/images/alipay.jpg`
3. [x] 验证页面显示情况
   - 通过 `pnpm run typecheck` 验证
   - 通过 `pnpm run lint` 验证
