<template>
  <!-- 中文注释：心愿页面，展示内置/自定义心愿卡片，支持编辑/删除/兑换与记录查看 -->
  <div class="p-4 space-y-4">
    <!-- 顶部：金币与操作按钮 -->
    <div class="flex items-center justify-between">
      <el-tag type="success">可用金币：{{ coins }}</el-tag>
      <div class="flex items-center gap-2">
        <el-button size="small" type="primary" @click="showRecords = true">领取记录</el-button>
        <el-button size="small" type="success" @click="openCreate()">创建心愿</el-button>
      </div>
    </div>

    <!-- 心愿列表：新用户自动显示 6 个内置心愿图标（参考/复制到 assets/wishs） -->
    <div>
      <div class="font-semibold mb-2">心愿</div>
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
                <el-button circle size="small" @click.stop="openEdit(w)"><el-icon><Edit /></el-icon></el-button>
              </el-tooltip>
              <el-tooltip content="删除">
                <el-button circle size="small" type="danger" @click.stop="onDelete(w)"><el-icon><Delete /></el-icon></el-button>
              </el-tooltip>
            </div>
          </div>

          <!-- 图标与名称 -->
          <div class="flex items-center gap-2">
            <img :src="resolveIcon(w.icon)" alt="icon" class="w-8 h-8 rounded" />
            <span class="font-medium">{{ w.name }}</span>
          </div>
          <div class="text-gray-500 text-sm mt-1">{{ w.content }}</div>

          <!-- 底部：左兑换次数，右兑换按钮 -->
          <div class="mt-2 flex items-center justify-between">
            <div class="text-xs text-gray-600">已兑换：{{ w.exchanged }} 次</div>
            <el-button type="warning" size="small" @click.stop="onExchange(w)">兑换</el-button>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 创建/编辑心愿弹窗 -->
    <el-dialog v-model="showForm" :title="formMode === 'create' ? '创建心愿' : '编辑心愿'" width="480px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="心愿图标">
          <div class="flex items-center gap-2">
            <!-- 中文注释：优先显示本地预览，其次显示已上传路径 -->
            <img v-if="form.icon_preview || form.icon" :src="form.icon_preview || resolveIcon(form.icon!)" class="w-10 h-10 rounded" />
            <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" @change="onPickIcon">
              <el-button type="primary" size="small">选择图片</el-button>
            </el-upload>
          </div>
          <!-- <div class="text-xs text-gray-500 mt-1">前端会压缩并转换为 webp，失败将回退原图。</div> -->
        </el-form-item>
        <el-form-item label="心愿名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="心愿描述"><el-input type="textarea" v-model="form.content" /></el-form-item>
        <el-form-item label="所需金币"><el-input-number v-model="form.need_coins" :min="1" /></el-form-item>
        <el-form-item label="单位"><el-select v-model="form.unit"><el-option label="个" value="个" /><el-option label="次" value="次" /><el-option label="分钟" value="分钟" /><el-option label="元" value="元" /></el-select></el-form-item>
        <el-form-item label="兑换数量"><el-input-number v-model="form.exchange_amount" :min="1" /></el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="showForm = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </div>
      </template>
    </el-dialog>

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
            </div>
          </div>
        </el-card>
        <div class="text-center text-xs text-gray-500">共 {{ records.total }} 条记录，页码 {{ records.page }}/{{ Math.ceil(records.total / records.page_size) || 1 }}</div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
// 中文注释：引入必要的 Vue API、全局状态、服务方法与图标组件
import { computed, onMounted, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete, Coin } from '@element-plus/icons-vue'
import { useAppState } from '@/stores/appState'
import { listWishes, createWish, updateWish, deleteWish, exchangeWish, listWishRecords, uploadWishIcon, toWebp, type Wish, type WishRecord } from '@/services/wishes'

// 中文注释：全局状态与本页状态
const store = useAppState()
const coins = computed(() => store.coins)
const userId = 1 // 中文注释：示例用户ID（真实项目中从登录信息获取）

// 心愿列表与表单/记录状态
const wishList = ref<Wish[]>([])
const showForm = ref(false)
const formMode = ref<'create' | 'edit'>('create')
// 中文注释：使用 reactive 而非 ref，确保 Element Plus 表单的 :model 与 v-model 正常工作
// 中文注释：表单类型扩展 icon_preview 字段用于本地预览
type WishForm = Partial<Wish> & { icon_preview?: string }
const form = reactive<WishForm>({ need_coins: 1, unit: '次', exchange_amount: 1, icon_preview: '' })
const showRecords = ref(false)
const records = ref<{ items: WishRecord[]; total: number; page: number; page_size: number }>({ items: [], total: 0, page: 1, page_size: 10 })

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
  return `${(import.meta as any).env.VITE_API_BASE || ''}/api/${icon}`.replace(/\/api\/$/, '/')
}
// 记录图标可复用 wish 名称对应文件
function resolveRecordIcon(name: string) {
  const f = `${name}.png`
  try { return new URL(`../assets/wishs/${f}`, import.meta.url).href } catch { return new URL(`../assets/wishs/领取记录.png`, import.meta.url).href }
}

function formatTime(ts: string) {
  const d = new Date(ts)
  const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// 加载心愿列表（新用户自动生成内置心愿）
async function loadWishes() {
  try {
    wishList.value = await listWishes(userId)
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
      user_id: userId,
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
  formMode.value = 'create'
  Object.assign(form, { user_id: userId, name: '', content: '', need_coins: 1, unit: '次', exchange_amount: 1, icon: '' })
  showForm.value = true
}
function openEdit(w: Wish) {
  formMode.value = 'edit'
  Object.assign(form, { ...w })
  showForm.value = true
}

// 选择图标并转换为 webp 上传到后端
async function onPickIcon(fileEvent: any) {
  // 中文注释：兼容 Element Plus Upload 的 onChange(uploadFile) 与原生 input 的 change 事件
  const raw: File | undefined = fileEvent?.raw || fileEvent?.target?.files?.[0] || fileEvent?.file
  if (!raw) return
  // 中文注释：先本地预览，再上传
  try { form.icon_preview = URL.createObjectURL(raw) } catch {}
  const webp = await toWebp(raw)
  try {
    const { path } = await uploadWishIcon(userId, webp)
    form.icon = path // 中文注释：保存后端返回的相对路径
    // 上传成功后清理临时预览
    try { form.icon_preview && URL.revokeObjectURL(form.icon_preview as any) } catch {}
    form.icon_preview = ''
  } catch (e) {
    // 失败回退：暂用本地文件名（不会持久化到服务器）
    form.icon = raw.name
  }
}

// 提交创建/编辑
async function submitForm() {
  try {
    if (formMode.value === 'create') {
      await createWish(form as any)
    } else {
      await updateWish(Number(form.id), form as any)
    }
    showForm.value = false
    await loadWishes()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

// 删除心愿
async function onDelete(w: Wish) {
  try {
    await deleteWish(w.id)
    await loadWishes()
  } catch (e: any) {
    ElMessage.error(e.message || '删除失败')
  }
}

// 兑换心愿
async function onExchange(w: Wish) {
  try {
    // 中文注释：弹出输入框选择兑换份数（默认为 1），后端将按 count 乘以所需金币扣减
    const { value } = await ElMessageBox.prompt(
      `输入兑换份数（每份消耗 ${w.need_coins} 金币，得到 ${w.exchange_amount} ${w.unit}）`,
      '选择兑换数量', { inputType: 'number', inputValue: '1', confirmButtonText: '确认', cancelButtonText: '取消' }
    )
    const count = Math.max(1, Number(value || 1))
    const resp = await exchangeWish(w.id, userId, count)
    // 中文注释：同步前端金币，与列表和记录刷新
    if (resp && typeof resp.user_coins !== 'undefined') {
      store.setCoins(Number(resp.user_coins))
    }
    await loadWishes()
    await loadRecords()
    ElMessage.success(`兑换成功「${w.name}」x${count}，扣除 ${w.need_coins * count} 金币！`)
  } catch (e: any) {
    if (e && e.message && /取消/.test(e.message)) return
    ElMessage.error(e.message || '兑换失败')
  }
}

// 加载记录
async function loadRecords(page = 1) {
  try {
    records.value = await listWishRecords(userId, page, records.value.page_size)
  } catch { records.value = { items: [], total: 0, page: 1, page_size: 10 } }
}
</script>

<style scoped>
</style>
