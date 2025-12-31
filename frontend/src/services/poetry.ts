import http from './http'

// 中文注释：古诗词 API 接口
// 包含获取和保存所有古诗词数据的接口

export const poetryService = {
  // 获取所有古诗词数据
  async getPoetryData() {
    return await http.get<{ data: string }>('/poetry/data')
  },

  // 保存所有古诗词数据
  // data 为 JSON 字符串
  async savePoetryData(data: string) {
    return await http.post('/poetry/data', { data })
  }
}
