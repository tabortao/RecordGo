import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginVue from "eslint-plugin-vue";
import vueParser from "vue-eslint-parser";
import json from "@eslint/json";
import { defineConfig } from "eslint/config";

export default defineConfig([
  // 忽略目录
  { ignores: ["dist/**", "node_modules/**"] },
  // JS/TS 基础规则
  { files: ["**/*.{js,mjs,cjs,ts,mts,cts}"], plugins: { js }, extends: ["js/recommended"], languageOptions: { globals: { ...globals.browser, ...globals.node } } },
  tseslint.configs.recommended,
  // 仅对 .vue 文件启用 Vue 解析与规则，避免应用到非 Vue 文件导致错误
  { files: ["**/*.vue"], plugins: { vue: pluginVue }, languageOptions: { parser: vueParser, parserOptions: { parser: tseslint.parser } }, rules: { "vue/multi-word-component-names": "off" } },
  // JSONC 文件
  { files: ["**/*.jsonc"], plugins: { json }, language: "json/jsonc", extends: ["json/recommended"] },
  // 项目容忍策略：允许 any 与空块，忽略下划线占位参数
  { files: ["**/*.{ts,js,vue}"], rules: { "@typescript-eslint/no-explicit-any": "off", "@typescript-eslint/no-empty-object-type": "off", "@typescript-eslint/no-unused-expressions": "off", "no-empty": "off", "@typescript-eslint/no-unused-vars": ["warn", { argsIgnorePattern: "^_", varsIgnorePattern: "^_" }] } },
  // d.ts 特例：关闭三斜线引用规则
  { files: ["**/*.d.ts"], rules: { "@typescript-eslint/triple-slash-reference": "off" } },
]);
