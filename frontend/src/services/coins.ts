import http from './http'

export async function setCoins(total: number, reason: string): Promise<{ coins: number }> {
  return await http.post('/coins/set', { coins: Math.max(0, Number(total || 0)), reason }) as any
}