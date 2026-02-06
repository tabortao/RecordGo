<template>
  <SettingsShell title="领取记录" subtitle="记录每一次兑换" :icon="List" tone="emerald" container-class="max-w-5xl">
    <SettingsCard title="兑换记录" description="按时间记录每一次兑换明细">
      <template #right>
        <div class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-3 py-1 text-xs font-extrabold text-gray-700 dark:text-gray-200">
          共 {{ total }} 条
        </div>
      </template>

      <div v-if="records.length === 0" class="py-10">
        <el-empty description="暂无领取记录" />
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
          <div
            v-for="item in records"
            :key="item.id"
            class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition"
          >
            <div class="flex items-start gap-3">
              <div class="relative">
                <div class="absolute -inset-2 rounded-2xl bg-emerald-200/35 dark:bg-emerald-500/10 blur-xl opacity-0 group-hover:opacity-100 transition-opacity" />
                <div class="relative w-10 h-10 rounded-2xl ring-1 ring-black/5 dark:ring-white/10 overflow-hidden bg-white/60 dark:bg-gray-950/25 grid place-items-center">
                  <div v-if="getRecordIcon(item).kind === 'emoji'" class="text-[20px] leading-none">{{ getRecordIcon(item).emoji }}</div>
                  <img v-else :src="getRecordIcon(item).src" class="w-full h-full object-contain" alt="icon" @error="onIconError" />
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0">
                    <div class="truncate text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">{{ item.wish_name }}</div>
                    <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ formatTime(item.created_at) }}</div>
                  </div>
                  <el-tag :type="statusType(item.status)" class="!rounded-full">{{ statusText(item.status) }}</el-tag>
                </div>

                <div class="mt-3 grid grid-cols-2 gap-2">
                  <div class="rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-950/30 px-3 py-2">
                    <div class="text-[11px] font-extrabold uppercase tracking-wider text-gray-500 dark:text-gray-400">数量</div>
                    <div class="mt-0.5 text-base font-extrabold text-gray-900 dark:text-gray-50">{{ item.amount }} {{ item.unit }}</div>
                  </div>
                  <div class="rounded-2xl border border-amber-100/70 dark:border-amber-900/40 bg-amber-50/60 dark:bg-amber-900/15 px-3 py-2">
                    <div class="text-[11px] font-extrabold uppercase tracking-wider text-amber-700/70 dark:text-amber-300/70">消耗金币</div>
                    <div class="mt-0.5 text-base font-extrabold text-amber-800 dark:text-amber-200">{{ item.coins_used }}</div>
                  </div>
                </div>

                <div v-if="item.remark && item.remark.trim()" class="mt-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/55 dark:bg-gray-950/25 px-3 py-2 text-sm text-gray-600 dark:text-gray-300">
                  备注：{{ item.remark }}
                </div>
              </div>
            </div>
          </div>
        </div>

      <div v-if="total > pageSize" class="flex justify-end mt-4">
        <el-pagination
          layout="prev, pager, next"
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          @current-change="onPageChange"
        />
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { List } from '@element-plus/icons-vue'
import { listWishRecords, listWishes, type WishRecord, type Wish } from '@/services/wishes'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/stores/auth'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
import { emojiChar, isEmojiIcon, useWishIconResolver } from '@/composables/useWishIcon'

const auth = useAuth()
const userId = computed(() => auth.user?.id ?? 0)

const records = ref<WishRecord[]>([])
const page = ref(1)
const pageSize = ref(9)
const total = ref(0)
const wishMap = ref<Record<number, Wish>>({})
const { resolveIconSrc, resolveBuiltinByName, prefetchPresigned } = useWishIconResolver({ fallbackBuiltin: '领取记录.png' })

async function load() {
  try {
    const resp = await listWishRecords(userId.value, page.value, pageSize.value)
    records.value = resp.items
    total.value = resp.total
    const wishes = await listWishes(userId.value)
    const map: Record<number, Wish> = {}
    for (const w of wishes) map[w.id] = w
    wishMap.value = map
    await prefetchPresigned(wishes.map(w => w.icon))
  } catch (e: any) {
    const msg = e?.message || '加载失败'
    ElMessage.error(`兑换记录加载失败：${msg}`)
  }
}

function onPageChange(p: number) { page.value = p; load() }

onMounted(async () => { await load() })

function getRecordIcon(r: WishRecord) {
  const w = wishMap.value[r.wish_id]
  const icon = w?.icon || ''
  if (isEmojiIcon(icon)) return { kind: 'emoji' as const, emoji: emojiChar(icon), src: '' }
  if (icon) return { kind: 'img' as const, emoji: '', src: resolveIconSrc(icon) }
  return { kind: 'img' as const, emoji: '', src: resolveBuiltinByName(w?.name || r.wish_name) }
}

function onIconError(e: Event) {
  const img = e.target as HTMLImageElement
  img.src = resolveBuiltinByName('领取记录')
}

function formatTime(ts: string) {
  const d = new Date(ts)
  if (Number.isNaN(d.getTime())) return ts
  const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`)
  return `${d.getMonth() + 1}月${d.getDate()}日 ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function statusType(s: string) {
  const v = String(s || '').toLowerCase()
  if (v === 'success' || v === 'ok') return 'success'
  if (v === 'failed' || v === 'fail' || v === 'error') return 'danger'
  if (v === 'pending' || v === 'processing') return 'warning'
  return 'info'
}

function statusText(s: string) {
  const v = String(s || '').toLowerCase()
  if (v === 'success' || v === 'ok') return '成功'
  if (v === 'failed' || v === 'fail' || v === 'error') return '失败'
  if (v === 'pending' || v === 'processing') return '处理中'
  return s || '未知'
}
</script>

<style scoped>
/* 中文注释：布局主要使用 Tailwind 与 Element Plus 组件 */
</style>
