<template>
  <div class="fixed inset-0 flex flex-col overflow-hidden bg-gray-50 dark:bg-gray-950 transition-colors duration-300">
    <div class="shrink-0 px-4 pt-4">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/75 dark:bg-gray-900/65 backdrop-blur-xl shadow-sm">
        <div class="flex items-center justify-between gap-3 px-3 py-3">
          <div class="flex items-center gap-3 min-w-0">
            <div class="h-10 w-10 rounded-2xl border border-emerald-200/70 dark:border-emerald-900/40 bg-emerald-50/80 dark:bg-emerald-900/25 flex items-center justify-center text-emerald-700 dark:text-emerald-300">
              <el-icon :size="18"><Coin /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-[17px] font-extrabold tracking-tight text-gray-900 dark:text-gray-50 truncate">心愿</div>
              <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 truncate">可用金币：<span class="font-extrabold text-amber-700 dark:text-amber-300">{{ coins }}</span></div>
            </div>
          </div>

          <div class="flex items-center gap-2 shrink-0">
            <el-button size="small" class="!rounded-2xl" @click="openRecords()">领取记录</el-button>
            <el-button v-if="isParent || canWishCreate" size="small" type="primary" class="!rounded-2xl" @click="openCreate()">创建心愿</el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto relative z-10 px-4 pt-5" style="padding-bottom: calc(env(safe-area-inset-bottom) + 96px)">

      <div class="max-w-6xl mx-auto">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
          <div
            v-for="w in wishList"
            :key="w.id"
            class="group relative rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition cursor-pointer flex flex-col min-h-[168px]"
            @click="toggleOps(w.id)"
          >
            <div class="absolute -top-2 -right-2 h-20 w-20 rounded-full bg-emerald-200/35 dark:bg-emerald-500/15 blur-2xl opacity-0 group-hover:opacity-100 transition-opacity" />

            <div class="flex items-start justify-between gap-3">
              <div class="flex items-center gap-3 min-w-0">
                <div class="relative shrink-0">
                  <div class="absolute -inset-2 rounded-2xl bg-emerald-200/40 dark:bg-emerald-500/10 blur-xl opacity-0 group-hover:opacity-100 transition-opacity" />
                  <div class="relative w-10 h-10 rounded-2xl ring-1 ring-black/5 dark:ring-white/10 overflow-hidden bg-white/60 dark:bg-gray-950/25 grid place-items-center">
                    <img :src="resolveIcon(w.icon)" alt="icon" class="w-full h-full object-contain" @error="onIconError" />
                  </div>
                </div>
                <div class="min-w-0">
                  <div class="truncate text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">{{ w.name }}</div>
                  <div class="mt-1 line-clamp-2 text-sm leading-relaxed text-gray-600 dark:text-gray-300">{{ w.content }}</div>
                </div>
              </div>

              <div class="flex items-center gap-2 shrink-0">
                <div class="rounded-full border border-amber-100/70 dark:border-amber-900/40 bg-amber-50/70 dark:bg-amber-900/15 px-2.5 py-1 text-xs font-extrabold text-amber-800 dark:text-amber-200 flex items-center gap-1">
                  <el-icon :size="14"><Coin /></el-icon>
                  <span>{{ w.need_coins }}</span>
                </div>
                <div v-if="opsVisibleId === w.id" class="flex items-center gap-1">
                  <el-tooltip content="编辑">
                    <el-button v-if="isParent || canWishEdit" circle size="small" class="!rounded-2xl" @click.stop="openEdit(w)"><el-icon><Edit /></el-icon></el-button>
                  </el-tooltip>
                  <el-tooltip content="删除">
                    <el-button v-if="isParent || canWishDelete" circle size="small" type="danger" class="!rounded-2xl" @click.stop="onDelete(w)"><el-icon><Delete /></el-icon></el-button>
                  </el-tooltip>
                </div>
              </div>
            </div>

            <div class="mt-auto flex items-center justify-between gap-2 pt-4 border-t border-gray-100/70 dark:border-gray-800/70">
              <div class="text-xs font-semibold text-gray-500 dark:text-gray-400">已兑换：{{ w.exchanged }} 次</div>
              <el-button
                type="warning"
                size="small"
                class="!rounded-2xl !font-bold"
                :disabled="!(isParent || canWishExchange)"
                @click.stop="onExchange(w)"
              >
                兑换
              </el-button>
            </div>
          </div>
        </div>
      </div>

    <!-- 中文注释：已改为独立页面，不再使用内置创建/编辑弹窗 -->

    <!-- 兑换记录抽屉 -->
    <el-drawer v-model="showRecords" title="兑换记录" :with-header="true" size="380px">
      <div class="space-y-3">
        <el-card v-for="r in records.items" :key="r.id" shadow="hover">
          <div class="flex items-center gap-3">
            <img :src="resolveRecordIcon(r.wish_name)" class="w-8 h-8 rounded" />
            <div class="flex-1">
              <div class="flex justify-between">
                <div>心愿：{{ r.wish_name }}</div>
                <div class="text-green-600">{{ r.status }}</div>
              </div>
              <div class="text-gray-500 text-xs">使用{{ r.coins_used }}金币 · {{ formatTime(r.created_at) }}</div>
              <!-- 中文注释：显示兑换备注（若有填写） -->
              <div class="text-gray-500 text-xs mt-1" v-if="r.remark && r.remark.trim()">备注：{{ r.remark }}</div>
            </div>
          </div>
        </el-card>
        <div class="text-center text-xs text-gray-500">共 {{ records.total }} 条记录，页码 {{ records.page }}/{{ Math.ceil(records.total / records.page_size) || 1 }}</div>
      </div>
    </el-drawer>

    <!-- 兑换心愿弹窗：支持输入数量与备注 -->
    <el-dialog
      v-model="showExchange"
      width="420px"
      class="pretty-exchange-dialog rounded-xl overflow-hidden shadow-xl"
    >
      <!-- 自定义标题栏，使用 Tailwind 增加背景与圆角 -->
      <template #header>
        <div class="relative overflow-hidden border-b border-gray-100 dark:border-gray-800 bg-white/78 dark:bg-gray-950/35 backdrop-blur-xl px-4 pt-4 pb-3">
          <div class="relative flex items-start gap-3">
            <div class="relative shrink-0">
              <div class="relative h-12 w-12 rounded-2xl border border-white/60 dark:border-gray-800/60 bg-white/70 dark:bg-gray-950/35 ring-1 ring-black/5 dark:ring-white/10 overflow-hidden">
                <img v-if="exchangeTarget" :src="resolveIcon(exchangeTarget.icon)" class="h-full w-full object-cover" />
              </div>
            </div>

            <div class="min-w-0 flex-1">
              <div class="text-[11px] font-extrabold uppercase tracking-wider text-gray-600 dark:text-gray-400">兑换心愿</div>
              <div class="mt-1 text-lg font-extrabold tracking-tight text-gray-900 dark:text-gray-50 truncate">
                {{ exchangeTarget?.name || '兑换心愿' }}
              </div>
              <div v-if="exchangeTarget?.content" class="mt-1 text-xs text-gray-600 dark:text-gray-300 line-clamp-1">
                {{ exchangeTarget.content }}
              </div>
            </div>

            <div class="shrink-0 flex flex-col items-end gap-1">
              <div class="rounded-full border border-amber-100/80 dark:border-amber-900/40 bg-amber-50/80 dark:bg-amber-900/15 px-3 py-1 text-xs font-extrabold text-amber-800 dark:text-amber-200">
                合计 {{ totalCost }} 金币
              </div>
              <div class="text-[11px] font-semibold text-gray-500 dark:text-gray-400">余额 {{ coins }}</div>
            </div>
          </div>

          <div v-if="insufficientCoins" class="relative mt-3 rounded-2xl border border-red-200/70 dark:border-red-900/40 bg-red-50/80 dark:bg-red-900/15 px-3 py-2 text-xs font-semibold text-red-700 dark:text-red-300">
            金币不足，先去完成任务赚取金币吧
          </div>
        </div>
      </template>

      <el-form label-position="top" class="pretty-exchange-form">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <el-form-item label="兑换数量">
            <el-input-number v-model="exchangeCount" :min="1" controls-position="right" class="w-full" />
          </el-form-item>
          <el-form-item label="单次金币">
            <el-input :model-value="String(exchangeTarget?.need_coins ?? '-')" disabled />
          </el-form-item>
        </div>
        <el-form-item label="备注">
          <el-input v-model="exchangeRemark" type="textarea" :rows="3" placeholder="可选，记录具体兑换情况" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between bg-gray-50/80 dark:bg-gray-950/40 border-t border-gray-100 dark:border-gray-800 px-4 py-3">
          <div class="text-xs font-semibold text-gray-500 dark:text-gray-400">
            将扣除 <span class="font-extrabold text-amber-700 dark:text-amber-300">{{ totalCost }}</span> 金币
          </div>
          <div class="flex justify-end gap-2">
            <el-button class="!rounded-2xl" @click="cancelExchange">取消</el-button>
            <el-button type="warning" class="!rounded-2xl !font-extrabold" :disabled="insufficientCoins" @click="confirmExchange">
              确认兑换
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：引入必要的 Vue API、全局状态、服务方法与图标组件
import { computed, onMounted, ref, watch } from 'vue'
import router from '@/router'
import { ElMessage } from 'element-plus'
import { Edit, Delete, Coin } from '@element-plus/icons-vue'
import { useAppState } from '@/stores/appState'
import { useAuth } from '@/stores/auth'
import { usePermissions } from '@/composables/permissions'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import { listWishes, deleteWish, exchangeWish, listWishRecords, type Wish, type WishRecord } from '@/services/wishes'

// 中文注释：全局状态与本页状态
const store = useAppState()
const coins = computed(() => store.coins)
// 中文注释：从认证状态读取当前登录用户ID；未登录时回退为 0（将触发默认心愿回退逻辑）
const auth = useAuth()
const userId = computed(() => auth.user?.id ?? 0)
// 中文注释：权限解析与常用布尔值
const { isParent, canWishCreate, canWishEdit, canWishDelete, canWishExchange } = usePermissions()

// 心愿列表与表单/记录状态
const wishList = ref<Wish[]>([])
const showRecords = ref(false)
const records = ref<{ items: WishRecord[]; total: number; page: number; page_size: number }>({ items: [], total: 0, page: 1, page_size: 10 })

// 兑换弹窗状态
const showExchange = ref(false)
const exchangeTarget = ref<Wish | null>(null)
const exchangeCount = ref<number>(1)
const exchangeRemark = ref<string>('')

const totalCost = computed(() => {
  const w = exchangeTarget.value
  if (!w) return 0
  const count = Math.max(1, Number(exchangeCount.value || 1))
  return Number(w.need_coins || 0) * count
})
const insufficientCoins = computed(() => {
  return totalCost.value > Number(coins.value || 0)
})

// 操作区显示切换：点击卡片右上角才显示编辑/删除
const opsVisibleId = ref<number | null>(null)
function toggleOps(id: number) {
  opsVisibleId.value = opsVisibleId.value === id ? null : id
}

// 中文注释：解析图标路径，内置图标在 assets/wishs 中，用户上传图标走后端 uploads 路径
const iconResolvedMap = ref<Record<string, string>>({})
function resolveIcon(icon: string | undefined) {
  if (!icon) return new URL('../assets/wishs/看电视.png', import.meta.url).href
  if (/\.(png|jpg|jpeg|webp)$/i.test(icon) && !icon.includes('/')) {
    return new URL(`../assets/wishs/${icon}`, import.meta.url).href
  }
  const base = getStaticBase()
  const p = String(icon).replace(/^\/+/, '')
  if (p.startsWith('uploads/')) return `${base}/api/${p}`
  return iconResolvedMap.value[p] || new URL('../assets/wishs/零花钱.png', import.meta.url).href
}
async function resolveIconsForList(list: Wish[]) {
  const targets = list
    .map(w => String(w.icon || ''))
    .filter(p => !!p && !p.startsWith('uploads/') && p.includes('/'))
  await Promise.all(targets.map(async (k) => {
    try { iconResolvedMap.value[k] = await presignView(k) } catch {}
  }))
}
// 记录图标可复用 wish 名称对应文件
function resolveRecordIcon(name: string) {
  const f = `${name}.png`
  try { return new URL(`../assets/wishs/${f}`, import.meta.url).href } catch { return new URL(`../assets/wishs/零花钱.png`, import.meta.url).href }
}

// 中文注释：图标加载失败时回退到默认占位图，避免破图影响体验
function onIconError(e: Event) {
  const img = e.target as HTMLImageElement
  try { img.src = new URL(`../assets/wishs/零花钱.png`, import.meta.url).href } catch {}
}

function formatTime(ts: string) {
  const d = new Date(ts)
  const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// 加载心愿列表（新用户自动生成内置心愿）
async function loadWishes() {
  try {
    wishList.value = await listWishes(userId.value)
  } catch (e) {
    // 兜底：调用默认心愿接口（增加 JSON 安全解析，避免空响应导致异常）
    let d: any = { data: [] }
    try {
      const resp = await fetch('/api/wishes/default')
      // 中文注释：确保响应为 JSON；若失败则使用空数据
      d = await resp.json()
    } catch (_) {
      d = { data: [] }
    }
    wishList.value = (d.data || []).map((x: any, i: number) => ({
      // 中文注释：内置图标强制与 assets/wishs 文件名一致：<心愿名称>.png
      id: i + 1,
      user_id: userId.value,
      name: x.name,
      content: x.content || `默认内置心愿：${x.name}`,
      icon: `${x.name}.png`,
      need_coins: x.coins ?? 1,
      // 中文注释：优先使用后端提供的单位与数量，保证与描述一致
      exchange_amount: x.exchange_amount ?? 1,
      unit: x.unit ?? '次',
      exchanged: 0,
      built_in: true
    }))
  }
}

onMounted(async () => {
  await loadWishes()
  await resolveIconsForList(wishList.value)
  await loadRecords()
})

watch(userId, async () => {
  await loadWishes()
  await resolveIconsForList(wishList.value)
  await loadRecords()
})

// 打开创建/编辑
function openCreate() {
  // 中文注释：权限校验：父账号允许；子账号需具备 wishes.create 权限
  if (!isParent.value && !canWishCreate.value) { ElMessage.warning('当前权限不允许创建心愿'); return }
  // 中文注释：改为进入独立创建页面，提升移动端体验
  router.push('/wishes/create')
}
function openEdit(w: Wish) {
  // 中文注释：权限校验：父账号允许；子账号需具备 wishes.edit 权限
  if (!isParent.value && !canWishEdit.value) { ElMessage.warning('当前权限不允许编辑心愿'); return }
  // 中文注释：改为进入独立编辑页面
  router.push(`/wishes/${w.id}/edit`)
}

// 删除心愿
async function onDelete(w: Wish) {
  // 中文注释：权限校验：父账号允许；子账号需具备 wishes.delete 权限
  if (!isParent.value && !canWishDelete.value) { ElMessage.warning('当前权限不允许删除心愿'); return }
  try {
    await deleteWish(w.id)
    await loadWishes()
  } catch (e: any) {
    ElMessage.error(e.message || '删除失败')
  }
}

// 兑换心愿：打开弹窗
function onExchange(w: Wish) {
  // 中文注释：权限校验：父账号允许；子账号需具备 wishes.exchange 权限
  if (!isParent.value && !canWishExchange.value) { ElMessage.warning('当前权限不允许兑换心愿'); return }
  exchangeTarget.value = w
  exchangeCount.value = 1
  exchangeRemark.value = ''
  showExchange.value = true
}

// 确认兑换：调用服务层并可选传备注
async function confirmExchange() {
  const w = exchangeTarget.value
  if (!w) { showExchange.value = false; return }
  const count = Math.max(1, Number(exchangeCount.value || 1))
  try {
    await exchangeWish(w.id, userId.value, count, exchangeRemark.value)
    // 中文注释：移除局部金币更新，统一交由响应拦截器处理最新金币同步
    await loadWishes()
    await loadRecords()
    const remarkText = (exchangeRemark.value || '').trim()
    ElMessage.success(`兑换成功「${w.name}」x${count}，扣除 ${w.need_coins * count} 金币！${remarkText ? '备注：' + remarkText : ''}`)
    showExchange.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '兑换失败')
  }
}

function cancelExchange() {
  showExchange.value = false
}

// 加载记录
async function loadRecords(page = 1) {
  try {
    records.value = await listWishRecords(userId.value, page, records.value.page_size)
  } catch { records.value = { items: [], total: 0, page: 1, page_size: 10 } }
}

// 中文注释：独立记录页面入口
function openRecords() {
  router.push('/wishes/records')
}
</script>

<style scoped>
/* 中文注释：美化兑换弹窗的主体内边距，限定作用于当前组件 */
:deep(.pretty-exchange-dialog .el-dialog__body) {
  padding: 16px;
}
/* 自适应移动端对话框宽度 */
:deep(.pretty-exchange-dialog .el-dialog) {
  max-width: calc(100vw - 24px);
}

:deep(.pretty-exchange-dialog .el-dialog__header) {
  padding: 0;
  margin: 0;
}

:deep(.pretty-exchange-dialog .el-dialog__footer) {
  padding: 0;
}

:deep(.pretty-exchange-dialog .el-input__wrapper),
:deep(.pretty-exchange-dialog .el-textarea__inner),
:deep(.pretty-exchange-dialog .el-input-number) {
  border-radius: 14px;
}

:deep(.pretty-exchange-dialog .el-form-item__label) {
  font-weight: 800;
  color: rgb(55 65 81);
}

.dark :deep(.pretty-exchange-dialog .el-form-item__label) {
  color: rgb(209 213 219);
}

.dark :deep(.pretty-exchange-dialog .el-dialog) {
  background: rgb(17 24 39 / 0.92);
  border: 1px solid rgb(31 41 55 / 0.9);
}

.dark :deep(.pretty-exchange-dialog .el-dialog__body),
.dark :deep(.pretty-exchange-dialog .el-dialog__footer) {
  background: transparent;
}

.dark :deep(.pretty-exchange-dialog .el-input__wrapper) {
  background: rgb(3 7 18 / 0.55);
  box-shadow: none;
}

.dark :deep(.pretty-exchange-dialog .el-textarea__inner) {
  background: rgb(3 7 18 / 0.55);
  color: rgb(243 244 246);
  border-color: rgb(31 41 55 / 0.9);
}

.dark :deep(.pretty-exchange-dialog .el-input__inner) {
  color: rgb(243 244 246);
}
</style>
