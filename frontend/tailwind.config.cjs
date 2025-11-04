/** 中文注释：Tailwind 配置，启用默认扫描路径 */
module.exports = {
  content: [
    './index.html',
    './src/**/*.{vue,ts,tsx}'
  ],
  theme: {
    extend: {
      colors: {
        primary: '#22c55e' // 绿色主题色
      }
    }
  },
  plugins: []
}