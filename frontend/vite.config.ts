import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// 中文注释：Vite 配置，集成 Vue、自动导入、Element Plus 组件解析，支持根据环境代理后端 API
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())
  const apiBase = env.VITE_API_BASE || 'http://localhost:8080'

  return {
    plugins: [
      vue(),
      AutoImport({
        // 中文注释：自动导入常用的 Vue API 与 Pinia
        imports: ['vue', 'vue-router', 'pinia'],
        resolvers: [ElementPlusResolver()],
        dts: 'src/auto-imports.d.ts'
      }),
      Components({
        // 中文注释：自动解析 Element Plus 组件
        resolvers: [ElementPlusResolver()],
        dts: 'src/components.d.ts'
      })
    ],
    // 中文注释：将 favicon 静态资源目录作为公共目录，便于通过根路径引用图标与 manifest
    publicDir: 'src/assets/favicon',
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: apiBase,
          changeOrigin: true
        }
      }
    },
    resolve: {
      alias: {
        '@': '/src'
      }
    }
  }
})