# 任务计划：input/textarea 字号与 span 对齐

日期：2026-02-06

## 任务目标
- 调整移动端下 `input` / `textarea` 的字体大小策略，使其与同区域文本（如 `span`）字号一致。

## 执行步骤
1. [x] 修改全局样式中对 `input/textarea` 的强制字号规则
2. [x] 在任务创建页 `/tasks/create` 验证 `input`、`textarea` 与 `span` 字号一致（均为 14px）
