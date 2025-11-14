/** 中文注释：Tailwind 配置，启用默认扫描路径 */
module.exports = {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{vue,ts,tsx}'
  ],
  theme: {
    extend: {
      colors: {
        primary: '#22c55e',
        brand: '#22c55e',
        sky: '#60a5fa',
        amber: '#f59e0b',
        pink: '#ec4899',
        teal: '#14b8a6'
      }
    }
  },
  plugins: []
}