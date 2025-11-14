<template>
  <div class="fixed top-0 left-0 right-0 bg-white z-40 border-b">
    <div class="px-4 py-2 font-semibold">兑换心愿</div>
    <div class="px-4 pb-2 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <el-button size="small" type="success" plain disabled>可用金币：{{ coins }}</el-button>
      </div>
      <div class="flex items-center gap-2">
        <el-button size="small" type="primary" @click="openRecords()">领取记录</el-button>
        <el-button v-if="isParent || canWishCreate" size="small" type="success" @click="openCreate()">创建心愿</el-button>
      </div>
    </div>
  </div>
  <div class="h-20"></div>
  <div class="p-4 space-y-4">

    <div>
      <!-- 中文注释：改为响应式网格布局，移动端单列，桌面端多列 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
        <el-card v-for="w in wishList" :key="w.id" shadow="hover" class="relative" @click="toggleOps(w.id)">
          <!-- 顶部右侧：先显示所需金币（图标+数量），点击卡片后显示编辑/删除图标 -->
          <div class="absolute top-2 right-2 flex items-center gap-2">
            <!-- 中文注释：所需金币常显，位于编辑/删除图标左侧 -->
            <div class="flex items-center gap-1 text-amber-600" title="所需金币">
              <el-icon :size="16"><Coin /></el-icon>
              <span class="text-sm font-medium">{{ w.need_coins }}</span>
            </div>
            <div v-if="opsVisibleId === w.id" class="flex items-center gap-2">
              <el-tooltip content="编辑">
                <el-button v-if="isParent || canWishEdit" circle size="small" @click.stop="openEdit(w)"><el-icon><Edit /></el-icon></el-button>
              </el-tooltip>
              <el-tooltip content="删除">
                <el-button v-if="isParent || canWishDelete" circle size="small" type="danger" @click.stop="onDelete(w)"><el-icon><Delete /></el-icon></el-button>
              </el-tooltip>
            </div>
          </div>

          <!-- 图标与名称 -->
          <div class="flex items-center gap-2">
            <img :src="resolveIcon(w.icon)" alt="icon" class="w-8 h-8 rounded" @error="onIconError" />
            <span class="font-medium">{{ w.name }}</span>
          </div>
          <div class="text-gray-500 text-sm mt-1">{{ w.content }}</div>

          <!-- 底部：左兑换次数，右兑换按钮 -->
          <div class="mt-2 flex items-center justify-between">
            <div class="text-xs text-gray-600">已兑换：{{ w.exchanged }} 次</div>
            <el-button type="warning" size="small" :disabled="!(isParent || canWishExchange)" @click.stop="onExchange(w)">兑换</el-button>
          </div>
        </el-card>
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
        <div class="bg-emerald-50 px-4 py-3 border-b border-emerald-100">
          <div class="flex items-center gap-2">
            <el-icon class="text-emerald-600"><Coin /></el-icon>
            <h3 class="text-emerald-800 font-semibold text-sm">{{ exchangeTarget ? `兑换：${exchangeTarget.name}` : '兑换心愿' }}</h3>
          </div>
        </div>
      </template>
      <el-form label-width="90px">
        <el-form-item label="兑换数量">
          <el-input-number v-model="exchangeCount" :min="1" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="exchangeRemark" type="textarea" placeholder="可选，记录具体兑换情况" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2 bg-gray-50 px-4 py-3">
          <el-button @click="cancelExchange">取消</el-button>
          <el-button type="primary" @click="confirmExchange">确认兑换</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
// 中文注释：引入必要的 Vue API、全局状态、服务方法与图标组件
import { computed, onMounted, ref } from 'vue'
import router from '@/router'
import { ElMessage } from 'element-plus'
import { Edit, Delete, Coin } from '@element-plus/icons-vue'
import { useAppState } from '@/stores/appState'
import { useAuth } from '@/stores/auth'
import { usePermissions } from '@/composables/permissions'
import { listWishes, deleteWish, exchangeWish, listWishRecords, type Wish, type WishRecord } from '@/services/wishes'

// 中文注释：全局状态与本页状态
const store = useAppState()
const coins = computed(() => store.coins)
// 中文注释：从认证状态读取当前登录用户ID；未登录时回退为 0（将触发默认心愿回退逻辑）
const auth = useAuth()
const userId = computed(() => auth.user?.id ?? 0)
// 中文注释：权限解析与常用布尔值
const { isParent, viewOnly, canWishCreate, canWishEdit, canWishDelete, canWishExchange } = usePermissions()

// 心愿列表与表单/记录状态
const wishList = ref<Wish[]>([])
const showRecords = ref(false)
const records = ref<{ items: WishRecord[]; total: number; page: number; page_size: number }>({ items: [], total: 0, page: 1, page_size: 10 })

// 兑换弹窗状态
const showExchange = ref(false)
const exchangeTarget = ref<Wish | null>(null)
const exchangeCount = ref<number>(1)
const exchangeRemark = ref<string>('')

// 操作区显示切换：点击卡片右上角才显示编辑/删除
const opsVisibleId = ref<number | null>(null)
function toggleOps(id: number) {
  opsVisibleId.value = opsVisibleId.value === id ? null : id
}

// 中文注释：解析图标路径，内置图标在 assets/wishs 中，用户上传图标走后端 uploads 路径
function resolveIcon(icon: string | undefined) {
  // 中文注释：默认图标改为相对路径，避免 @ 别名在 new URL 中解析异常
  if (!icon) return new URL('../assets/wishs/看电视.png', import.meta.url).href
  // 内置：文件名.png；自定义：uploads/... 路径
  if (/\.(png|jpg|jpeg|webp)$/i.test(icon) && !icon.includes('/')) {
    return new URL(`../assets/wishs/${icon}`, import.meta.url).href
  }
  // 中文注释：静态文件走后端基址，无需 /api 前缀；若未设置基址则回退到本机 8080
  let base = ((import.meta as any).env.VITE_API_BASE || '').replace(/\/+$/, '')
  if (!base) {
    try {
      const url = new URL(window.location.href)
      const host = url.hostname || 'localhost'
      base = `${url.protocol}//${host}:8080`
    } catch {
      base = 'http://localhost:8080'
    }
  }
  const path = String(icon).replace(/^\/+/, '')
  // 中文注释：后端静态文件映射为 /api/uploads，需要加上 /api 前缀
  return `${base}/api/${path}`
}
// 记录图标可复用 wish 名称对应文件
function resolveRecordIcon(name: string) {
  const f = `${name}.png`
  try { return new URL(`../assets/wishs/${f}`, import.meta.url).href } catch { return new URL(`../assets/wishs/领取记录.png`, import.meta.url).href }
}

// 中文注释：图标加载失败时回退到默认占位图，避免破图影响体验
function onIconError(e: Event) {
  const img = e.target as HTMLImageElement
  try { img.src = new URL(`../assets/wishs/领取记录.png`, import.meta.url).href } catch {}
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
</style>
