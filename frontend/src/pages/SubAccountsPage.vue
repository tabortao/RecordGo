<template>
  <!-- 中文注释：子账号管理页面，支持列表、创建、编辑、删除与令牌生成 -->
  <div class="p-4 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer" title="返回" @click="router.push('/mine')"><ArrowLeft /></el-icon>
        <el-icon :size="18" style="color:#22c55e"><User /></el-icon>
        <span class="font-semibold">子账号管理</span>
      </div>
      <div class="text-sm text-gray-500">共 {{ children.length }} 个子账号</div>
    </div>
    <el-card shadow="never">

      <!-- 权限不足提示 -->
      <div v-if="!allowManage" class="py-6">
        <div class="text-gray-700 mb-3">当前权限不允许管理子账号。</div>
        <el-button type="primary" @click="router.push('/mine')">返回我的</el-button>
      </div>

      <div v-else>
        <div class="flex justify-between items-center mb-3">
          <div class="text-gray-600 text-sm">提示：令牌登录仅限子账号使用。</div>
          <el-button type="success" @click="openCreate">创建子账号</el-button>
        </div>

        <el-table :data="children" v-loading="loading" border>
          <el-table-column label="头像" width="80">
            <template #default="{ row }">
              <el-avatar :size="36" :src="resolveAvatarUrl(row.avatar_path)" />
            </template>
          </el-table-column>
          <el-table-column prop="nickname" label="昵称" min-width="140" />
          <el-table-column label="权限 JSON" min-width="220">
            <template #default="{ row }">
              <div class="text-xs text-gray-600 break-words">{{ row.permissions?.trim() || '{"view_only":true}' }}</div>
            </template>
          </el-table-column>
          <el-table-column label="令牌" min-width="360">
            <template #default="{ row }">
              <div class="flex items-center gap-2">
                <el-input v-model="row.login_token" placeholder="未生成" readonly class="flex-1" />
                <el-select v-model="expiresMap[row.id]" size="small" style="width:120px" placeholder="有效期">
                  <el-option label="永久" :value="0" />
                  <el-option label="1小时" :value="3600" />
                  <el-option label="24小时" :value="86400" />
                  <el-option label="7天" :value="604800" />
                </el-select>
                <el-button size="small" @click="copy(row.login_token)" :disabled="!row.login_token">复制</el-button>
                <el-button size="small" type="primary" @click="refreshToken(row)">刷新令牌</el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180">
            <template #default="{ row }">
              <div class="flex items-center gap-2">
                <el-button size="small" @click="openEdit(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="onDelete(row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 创建/编辑 子账号弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑子账号' : '创建子账号'" width="720px">
      <el-form label-position="top" @submit.prevent>
        <el-form-item label="昵称" required>
          <el-input v-model="form.nickname" placeholder="请输入子账号昵称" />
        </el-form-item>
        <!-- 中文注释：头像上传（前端转 webp），创建时先暂存，保存后再上传并写入子账号头像路径 -->
        <el-form-item label="头像">
          <div class="flex items-center gap-3">
            <el-avatar :size="48" :src="avatarPreview || resolveAvatarUrl(currentAvatarPath)" />
            <el-button size="small" @click="triggerAvatarPick">选择头像</el-button>
            <input ref="avatarInput" type="file" accept="image/*" class="hidden" @change="onAvatarPicked" />
            <el-button v-if="avatarPreview" size="small" @click="clearAvatar">移除</el-button>
          </div>
          <div class="text-xs text-gray-500 mt-1">前端自动转换为 webp；建议小于 2MB。</div>
        </el-form-item>
        <el-form-item label="权限模板（可选）">
          <div class="flex items-center gap-3">
            <el-select v-model="permTemplate" placeholder="选择模板以快速填充" style="width:240px" @change="applyTemplate">
              <el-option label="仅查看" value="view_only" />
              <el-option label="可完成任务（仅状态）" value="tasks_status" />
              <el-option label="任务增删改查" value="tasks_full" />
              <el-option label="全部权限" value="full" />
            </el-select>
            <span class="text-xs text-gray-500">选择后可在下方勾选微调</span>
          </div>
        </el-form-item>

        <!-- 中文注释：勾选式权限编辑（模块/动作粒度） -->
        <el-form-item label="基础权限">
          <el-checkbox v-model="permModel.view_only">仅查看（禁用所有操作）</el-checkbox>
        </el-form-item>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="rounded border p-3">
            <div class="font-medium mb-2">账户</div>
            <el-checkbox v-model="permModel.account.manage_children" :disabled="permModel.view_only">管理子账号</el-checkbox>
          </div>
          <div class="rounded border p-3">
            <div class="font-medium mb-2">任务</div>
            <div class="space-y-1">
              <el-checkbox v-model="permModel.tasks.create" :disabled="permModel.view_only">创建</el-checkbox>
              <el-checkbox v-model="permModel.tasks.edit" :disabled="permModel.view_only">编辑</el-checkbox>
              <el-checkbox v-model="permModel.tasks.delete" :disabled="permModel.view_only">删除</el-checkbox>
              <el-checkbox v-model="permModel.tasks.status" :disabled="permModel.view_only">更改状态</el-checkbox>
            </div>
          </div>
          <div class="rounded border p-3">
            <div class="font-medium mb-2">心愿</div>
            <div class="space-y-1">
              <el-checkbox v-model="permModel.wishes.create" :disabled="permModel.view_only">创建</el-checkbox>
              <el-checkbox v-model="permModel.wishes.edit" :disabled="permModel.view_only">编辑</el-checkbox>
              <el-checkbox v-model="permModel.wishes.delete" :disabled="permModel.view_only">删除</el-checkbox>
              <el-checkbox v-model="permModel.wishes.exchange" :disabled="permModel.view_only">兑换</el-checkbox>
            </div>
          </div>
        </div>

        <!-- 中文注释：设置权限分组（控制“我的”页设置入口）-->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-3">
          <div class="rounded border p-3">
            <div class="font-medium mb-2">设置（系统）</div>
            <div class="space-y-1">
              <el-checkbox v-model="permModel.settings.tomato">番茄钟设置</el-checkbox>
              <el-checkbox v-model="permModel.settings.tasks" :disabled="permModel.view_only">任务设置</el-checkbox>
              <el-checkbox v-model="permModel.settings.coins" :disabled="permModel.view_only">金币设置</el-checkbox>
              <el-checkbox v-model="permModel.settings.reading">朗读设置</el-checkbox>
            </div>
            <div class="text-xs text-gray-500 mt-2">提示：仅查看模式下“任务设置/金币设置”不可用。</div>
          </div>
        </div>

        <div class="text-xs text-gray-500 mt-3">提交时将自动生成与后端一致的权限 JSON。</div>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="dialogVisible=false">取消</el-button>
          <el-button type="primary" @click="submit">{{ isEdit ? '保存' : '创建' }}</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
 </template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import router from '@/router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getStaticBase } from '@/services/http'
import { usePermissions } from '@/composables/permissions'
import { listSubAccounts, createSubAccount, updateSubAccount, deleteSubAccount, generateChildToken, type ChildAccount } from '@/services/subaccounts'
import { prepareUpload } from '@/utils/image'
import { uploadAvatar } from '@/services/user'
import { ArrowLeft, User } from '@element-plus/icons-vue'

// 中文注释：权限校验（家长允许；子账号需具备 account.manage_children）
const { isParent, manageChildren } = usePermissions()
const allowManage = computed(() => isParent.value || manageChildren.value)

const children = ref<ChildAccount[]>([])
const loading = ref(false)
const expiresMap = ref<Record<number, number>>({})

// 创建/编辑弹窗状态
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const form = ref<{ nickname: string }>({ nickname: '' })
const permTemplate = ref<'view_only' | 'tasks_status' | 'tasks_full' | 'full' | ''>('')

// 中文注释：可视化权限模型（与后端字段一致）
type PermModel = {
  view_only: boolean
  account: { manage_children: boolean }
  tasks: { create: boolean; edit: boolean; delete: boolean; status: boolean }
  wishes: { create: boolean; edit: boolean; delete: boolean; exchange: boolean }
  // 中文注释：新增设置权限分组（用于控制“我的”页设置入口权限）
  settings: { tomato: boolean; tasks: boolean; coins: boolean; reading: boolean }
}
const defaultPerms: PermModel = {
  view_only: true,
  account: { manage_children: false },
  tasks: { create: false, edit: false, delete: false, status: false },
  wishes: { create: false, edit: false, delete: false, exchange: false },
  settings: { tomato: false, tasks: false, coins: false, reading: false }
}
const permModel = ref<PermModel>({ ...defaultPerms })

// 中文注释：头像选择/预览与上传（创建或编辑）
const avatarInput = ref<HTMLInputElement | null>(null)
const avatarPreview = ref<string>('')
const avatarFile = ref<File | null>(null)
const currentAvatarPath = ref<string>('')
function triggerAvatarPick() { avatarInput.value?.click() }
async function onAvatarPicked(e: Event) {
  const input = e.target as HTMLInputElement
  const f = (input.files || [])[0]
  if (!f) return
  // 大小限制：< 5MB（与后端一致）
  if (f.size > 5 * 1024 * 1024) { ElMessage.error('头像过大（需小于5MB）'); return }
  try {
    const webp = await prepareUpload(f)
    avatarFile.value = webp
    avatarPreview.value = URL.createObjectURL(webp)
  } catch (err: any) {
    ElMessage.error(err?.message || '头像处理失败')
  }
}
function clearAvatar() { avatarFile.value = null; avatarPreview.value = '' }

function resolveAvatarUrl(p?: string | null) {
  if (!p) return ''
  const s = String(p)
  if (/^https?:\/\//i.test(s)) return s
  if (!/uploads\//i.test(s)) return ''
  const base = getStaticBase()
  return `${base}/api/${s.replace(/^\/+/, '')}`
}

function applyTemplate(tpl: string) {
  // 中文注释：根据模板填充可视化权限
  const m = { ...defaultPerms }
  switch (tpl) {
    case 'view_only': m.view_only = true; break
    case 'tasks_status':
      m.view_only = false
      m.tasks.status = true
      break
    case 'tasks_full':
      m.view_only = false
      m.tasks = { create: true, edit: true, delete: true, status: true }
      break
    case 'full':
      m.view_only = false
      m.account.manage_children = true
      m.tasks = { create: true, edit: true, delete: true, status: true }
      m.wishes = { create: true, edit: true, delete: true, exchange: true }
      m.settings = { tomato: true, tasks: true, coins: true, reading: true }
      break
    default:
      break
  }
  permModel.value = m
}

// 中文注释：JSON <-> 模型转换，保证与后端一致
function parsePermsJSON(s: string | undefined) {
  try {
    if (!s) { permModel.value = { ...defaultPerms }; return }
    const obj = JSON.parse(s)
    const m = { ...defaultPerms }
    m.view_only = !!obj.view_only
    if (obj.account && typeof obj.account.manage_children === 'boolean') m.account.manage_children = obj.account.manage_children
    if (obj.tasks) {
      m.tasks.create = !!obj.tasks.create
      m.tasks.edit = !!obj.tasks.edit
      m.tasks.delete = !!obj.tasks.delete
      m.tasks.status = !!obj.tasks.status
    }
    if (obj.wishes) {
      m.wishes.create = !!obj.wishes.create
      m.wishes.edit = !!obj.wishes.edit
      m.wishes.delete = !!obj.wishes.delete
      m.wishes.exchange = !!obj.wishes.exchange
    }
    if (obj.settings) {
      m.settings.tomato = !!obj.settings.tomato
      m.settings.tasks = !!obj.settings.tasks
      m.settings.coins = !!obj.settings.coins
      m.settings.reading = !!obj.settings.reading
    }
    permModel.value = m
  } catch (_) {
    permModel.value = { ...defaultPerms }
  }
}

function toPermsJSON(): string {
  // 中文注释：生成最小化 JSON（仅包含启用的模块/动作）；view_only 优先
  const m = permModel.value
  const out: any = { view_only: !!m.view_only }
  if (!m.view_only) {
    if (m.account.manage_children) out.account = { manage_children: true }
    const t: any = {}
    ;(['create','edit','delete','status'] as const).forEach(k => { if ((m.tasks as any)[k]) t[k] = true })
    if (Object.keys(t).length) out.tasks = t
    const w: any = {}
    ;(['create','edit','delete','exchange'] as const).forEach(k => { if ((m.wishes as any)[k]) w[k] = true })
    if (Object.keys(w).length) out.wishes = w
    const s: any = {}
    ;(['tomato','tasks','coins','reading'] as const).forEach(k => { if ((m.settings as any)[k]) s[k] = true })
    if (Object.keys(s).length) out.settings = s
  }
  return JSON.stringify(out)
}

async function loadList() {
  if (!allowManage.value) return
  loading.value = true
  try {
    children.value = await listSubAccounts()
    // 中文注释：为每个子账号初始化令牌有效期（默认永久）
    children.value.forEach(c => { if (expiresMap.value[c.id] == null) expiresMap.value[c.id] = 0 })
  } catch (e: any) {
    ElMessage.error(e?.message || '加载子账号失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  if (!allowManage.value) { ElMessage.warning('无权限创建子账号'); return }
  isEdit.value = false
  editingId.value = null
  permTemplate.value = ''
  form.value = { nickname: '' }
  permModel.value = { ...defaultPerms }
  currentAvatarPath.value = ''
  clearAvatar()
  dialogVisible.value = true
}

function openEdit(row: ChildAccount) {
  if (!allowManage.value) { ElMessage.warning('无权限编辑子账号'); return }
  isEdit.value = true
  editingId.value = row.id
  permTemplate.value = ''
  form.value = { nickname: row.nickname || '' }
  parsePermsJSON(row.permissions)
  currentAvatarPath.value = row.avatar_path || ''
  clearAvatar()
  dialogVisible.value = true
}

async function submit() {
  if (!allowManage.value) { ElMessage.warning('无权限'); return }
  if (permTemplate.value) applyTemplate(permTemplate.value)
  const payload = {
    nickname: form.value.nickname.trim(),
    permissions: toPermsJSON()
  }
  if (!payload.nickname) { ElMessage.error('请输入昵称'); return }
  try {
    if (isEdit.value && editingId.value) {
      await updateSubAccount(editingId.value, { nickname: payload.nickname, permissions: payload.permissions })
      // 中文注释：如选择了新头像，上传并写入子账号头像路径
      if (avatarFile.value) {
        try {
          const { path } = await uploadAvatar(editingId.value, avatarFile.value)
          await updateSubAccount(editingId.value, { avatar_path: path })
        } catch (err: any) {
          console.error('上传子账号头像失败', err)
          ElMessage.error(err?.message || '头像上传失败')
        }
      }
      ElMessage.success('已保存')
    } else {
      const created = await createSubAccount(payload)
      // 中文注释：如选择了头像，创建后上传并保存路径
      if (avatarFile.value && created?.id) {
        try {
          const { path } = await uploadAvatar(created.id, avatarFile.value)
          await updateSubAccount(created.id, { avatar_path: path })
        } catch (err: any) {
          console.error('上传子账号头像失败', err)
          ElMessage.error(err?.message || '头像上传失败')
        }
      }
      ElMessage.success('已创建子账号')
    }
    dialogVisible.value = false
    await loadList()
  } catch (e: any) {
    ElMessage.error(e?.message || '提交失败')
  }
}

async function onDelete(row: ChildAccount) {
  if (!allowManage.value) { ElMessage.warning('无权限删除子账号'); return }
  try {
    await ElMessageBox.confirm('删除后不可恢复，确认删除该子账号？', '确认', { type: 'warning' })
    await deleteSubAccount(row.id)
    ElMessage.success('已删除')
    await loadList()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '删除失败')
  }
}

async function refreshToken(row: ChildAccount) {
  if (!allowManage.value) { ElMessage.warning('无权限生成令牌'); return }
  try {
    const sec = expiresMap.value[row.id] ?? 0
    const r = await generateChildToken(row.id, { expiresInSeconds: sec || undefined })
    row.login_token = r.token
    ElMessage.success('令牌已刷新')
  } catch (e: any) {
    ElMessage.error(e?.message || '生成令牌失败')
  }
}

async function copy(text?: string) {
  if (!text) return
  try { await navigator.clipboard.writeText(String(text)); ElMessage.success('已复制') } catch (_) { ElMessage.info('复制失败，请手动选择复制') }
}

onMounted(() => { loadList() })
</script>

<style scoped>
</style>