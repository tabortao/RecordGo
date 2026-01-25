<template>
  <SettingsShell title="领取记录" subtitle="记录每一次兑换" :icon="List" tone="emerald" container-class="max-w-5xl">
    <SettingsCard>
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
        <div
          v-for="item in records"
          :key="item.id"
          class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition"
        >
          <div class="flex items-start gap-3">
            <div class="relative">
              <div class="absolute -inset-2 rounded-2xl bg-emerald-200/35 dark:bg-emerald-500/10 blur-xl opacity-0 group-hover:opacity-100 transition-opacity" />
              <img :src="getIconSrc(item)" class="relative w-10 h-10 rounded-2xl ring-1 ring-black/5 dark:ring-white/10" />
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <div class="truncate text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">{{ item.wish_name }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ item.created_at }}</div>
                </div>
                <el-tag :type="item.status === 'success' ? 'success' : 'info'" class="!rounded-full">{{ item.status }}</el-tag>
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

      <div class="flex justify-end mt-4">
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
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/stores/auth'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const auth = useAuth()
const userId = computed(() => auth.user?.id ?? 0)

const records = ref<WishRecord[]>([])
const page = ref(1)
const pageSize = ref(9)
const total = ref(0)
const wishMap = ref<Record<number, Wish>>({})

async function load() {
  try {
    const resp = await listWishRecords(userId.value, page.value, pageSize.value)
    records.value = resp.items
    total.value = resp.total
    const wishes = await listWishes(userId.value)
    const map: Record<number, Wish> = {}
    for (const w of wishes) map[w.id] = w
    wishMap.value = map
    await resolveIconsForWishes(wishes)
  } catch (e: any) {
    const msg = e?.message || '加载失败'
    ElMessage.error(`兑换记录加载失败：${msg}`)
  }
}

function onPageChange(p: number) { page.value = p; load() }

onMounted(async () => { await load() })

// 中文注释：解析心愿图标路径（兼容内置 assets 与后端 uploads 相对路径）
const iconResolvedMap = ref<Record<string,string>>({})
function resolveIcon(icon?: string) {
  if (!icon) return new URL('../assets/wishs/领取记录.png', import.meta.url).href
  if (/\.(png|jpg|jpeg|webp)$/i.test(icon) && !icon.includes('/')) {
    return new URL(`../assets/wishs/${icon}`, import.meta.url).href
  }
  const base = getStaticBase()
  const path = String(icon).replace(/^\/+/, '')
  if (path.startsWith('uploads/')) return `${base}/api/${path}`
  return iconResolvedMap.value[path] || new URL('../assets/wishs/领取记录.png', import.meta.url).href
}
async function resolveIconsForWishes(list: Wish[]) {
  const targets = list
    .map(w => String(w.icon || ''))
    .filter(p => !!p && !p.startsWith('uploads/') && p.includes('/'))
  await Promise.all(targets.map(async (k) => { try { iconResolvedMap.value[k] = await presignView(k) } catch {} }))
}

function getIconSrc(r: WishRecord) {
  const w = wishMap.value[r.wish_id]
  if (w && w.icon) return resolveIcon(w.icon)
  // 中文注释：找不到映射时，回退到名称同名的内置图标
  try { return new URL(`../assets/wishs/${r.wish_name}.png`, import.meta.url).href } catch { return new URL(`../assets/wishs/领取记录.png`, import.meta.url).href }
}
</script>

<style scoped>
/* 中文注释：布局主要使用 Tailwind 与 Element Plus 组件 */
</style>
