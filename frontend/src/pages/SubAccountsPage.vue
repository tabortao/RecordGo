<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 pb-20 transition-colors duration-300">
    <!-- 顶部导航栏 -->
    <div class="sticky top-0 z-30 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800 px-4 h-14 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button 
          @click="router.push('/mine')"
          class="p-2 -ml-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-600 dark:text-gray-300">
          <el-icon :size="20"><ArrowLeft /></el-icon>
        </button>
        <span class="font-bold text-lg text-gray-800 dark:text-gray-100">子账号管理</span>
      </div>
      <el-button v-if="allowManage" type="primary" size="small" @click="openCreate" class="!rounded-lg !px-4">
        <el-icon class="mr-1"><Plus /></el-icon>创建
      </el-button>
    </div>

    <div class="max-w-4xl mx-auto p-4 md:p-6 animate-fade-in-up">
      <!-- 权限不足提示 -->
      <div v-if="!allowManage" class="flex flex-col items-center justify-center py-20 bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800">
        <div class="w-20 h-20 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mb-4 text-gray-400">
          <el-icon :size="40"><Lock /></el-icon>
        </div>
        <div class="text-gray-600 dark:text-gray-400 font-medium mb-4">当前权限不允许管理子账号</div>
        <el-button type="primary" @click="router.push('/mine')" class="!rounded-lg">返回个人中心</el-button>
      </div>

      <div v-else>
        <!-- 统计与提示 -->
        <div class="flex items-center justify-between mb-4 px-1">
          <div class="text-sm text-gray-500 dark:text-gray-400">共 {{ children.length }} 个子账号</div>
          <div class="text-xs text-orange-500 bg-orange-50 dark:bg-orange-900/20 px-2 py-1 rounded border border-orange-100 dark:border-orange-900/30">
            提示：令牌登录仅限子账号使用
          </div>
        </div>

        <!-- 子账号列表 (卡片式) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="row in children" :key="row.id" 
            class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 overflow-hidden hover:shadow-md transition-shadow duration-300 flex flex-col">
            
            <div class="p-5 flex items-center gap-4 border-b border-gray-50 dark:border-gray-800/50">
              <el-avatar :size="56" :src="avatarsMap[row.id] || defaultAvatar" class="border-2 border-white dark:border-gray-800 shadow-sm" />
              <div class="flex-1 min-w-0">
                <div class="font-bold text-lg text-gray-800 dark:text-gray-100 truncate">{{ row.nickname }}</div>
                <div class="text-xs text-gray-500 mt-1 flex items-center gap-1">
                  <el-icon><Key /></el-icon>
                  <span>ID: {{ row.id }}</span>
                </div>
              </div>
              <el-dropdown trigger="click">
                <button class="p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 text-gray-400">
                  <el-icon :size="20"><MoreFilled /></el-icon>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="openEdit(row)"><el-icon><Edit /></el-icon>编辑信息</el-dropdown-item>
                    <el-dropdown-item divided class="text-red-500" @click="onDelete(row)"><el-icon><Delete /></el-icon>删除账号</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>

            <div class="p-5 flex-1 flex flex-col gap-4">
              <!-- 权限概览 -->
              <div class="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-3">
                <div class="text-xs font-medium text-gray-500 mb-2">权限配置</div>
                <div class="text-xs text-gray-600 dark:text-gray-400 break-all font-mono leading-relaxed line-clamp-2">
                  {{ row.permissions?.trim() || '{"view_only":true}' }}
                </div>
              </div>

              <!-- 令牌管理 -->
              <div class="mt-auto">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs font-medium text-gray-500">登录令牌</span>
                  <div class="flex gap-2">
                     <el-select v-model="expiresMap[row.id]" size="small" style="width:90px" placeholder="有效期" class="custom-tiny-select">
                        <el-option label="永久" :value="0" />
                        <el-option label="1小时" :value="3600" />
                        <el-option label="24小时" :value="86400" />
                        <el-option label="7天" :value="604800" />
                      </el-select>
                      <button @click="refreshToken(row)" class="text-xs text-blue-500 hover:text-blue-600 font-medium">刷新</button>
                  </div>
                </div>
                
                <div class="flex items-center gap-2">
                  <div class="flex-1 bg-gray-50 dark:bg-gray-800 rounded-lg px-3 py-2 text-xs text-gray-600 dark:text-gray-300 font-mono truncate border border-gray-100 dark:border-gray-700">
                    {{ row.login_token || '未生成' }}
                  </div>
                  <button 
                    @click="copy(row.login_token)" 
                    :disabled="!row.login_token"
                    class="p-2 bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg hover:bg-blue-100 dark:hover:bg-blue-900/40 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                    <el-icon><CopyDocument /></el-icon>
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 创建卡片 (如果没有数据时显示) -->
          <button v-if="children.length === 0 && !loading" @click="openCreate" class="flex flex-col items-center justify-center h-64 border-2 border-dashed border-gray-200 dark:border-gray-800 rounded-2xl hover:border-blue-400 hover:bg-blue-50/50 dark:hover:border-blue-800 dark:hover:bg-blue-900/10 transition-all group">
             <div class="w-16 h-16 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center text-blue-500 mb-4 group-hover:scale-110 transition-transform">
                <el-icon :size="32"><Plus /></el-icon>
             </div>
             <span class="font-medium text-gray-600 dark:text-gray-400">创建第一个子账号</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 创建/编辑 子账号弹窗 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑子账号' : '创建子账号'" 
      :width="isMobile ? '95vw' : '680px'"
      class="rounded-2xl overflow-hidden custom-dialog"
      append-to-body
    >
      <el-form label-position="top" @submit.prevent class="max-h-[70vh] overflow-y-auto px-1">
        <!-- 基础信息 -->
        <div class="grid grid-cols-1 sm:grid-cols-[auto_1fr] gap-6 mb-6">
          <div class="flex flex-col items-center gap-2">
            <div class="relative group cursor-pointer" @click="triggerAvatarPick">
               <el-avatar :size="80" :src="avatarPreview || avatarDialogSrc" class="border shadow-sm" />
               <div class="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity text-white text-xs">更换</div>
            </div>
            <div v-if="avatarPreview" class="text-xs text-red-500 cursor-pointer" @click="clearAvatar">移除</div>
            <input ref="avatarInput" type="file" accept="image/*" class="hidden" @change="onAvatarPicked" />
          </div>
          
          <div class="space-y-4">
            <el-form-item label="昵称" required class="!mb-0">
              <el-input v-model="form.nickname" placeholder="请输入子账号昵称" size="large" />
            </el-form-item>
            
            <el-form-item label="快速模板" class="!mb-0">
              <el-select v-model="permTemplate" placeholder="选择权限模板" @change="applyTemplate" size="large">
                <el-option label="仅查看" value="view_only" />
                <el-option label="可完成任务（仅状态）" value="tasks_status" />
                <el-option label="任务增删改查" value="tasks_full" />
                <el-option label="全部权限" value="full" />
              </el-select>
            </el-form-item>
          </div>
        </div>

        <el-divider content-position="left">详细权限配置</el-divider>

        <!-- 权限配置 -->
        <div class="space-y-4">
          <div class="flex items-center p-3 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-100 dark:border-gray-700">
             <el-checkbox v-model="permModel.view_only" size="large">
               <span class="font-medium">仅查看模式</span>
               <span class="text-xs text-gray-500 ml-2">禁用所有增删改操作</span>
             </el-checkbox>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- 账户权限 -->
            <div class="p-4 border border-gray-100 dark:border-gray-700 rounded-xl hover:shadow-sm transition-shadow">
              <div class="font-bold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-2">
                <el-icon class="text-blue-500"><User /></el-icon>账户权限
              </div>
              <el-checkbox v-model="permModel.account.manage_children" :disabled="permModel.view_only">管理子账号</el-checkbox>
            </div>

            <!-- 任务权限 -->
            <div class="p-4 border border-gray-100 dark:border-gray-700 rounded-xl hover:shadow-sm transition-shadow">
              <div class="font-bold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-2">
                <el-icon class="text-green-500"><List /></el-icon>任务权限
              </div>
              <div class="grid grid-cols-2 gap-2">
                <el-checkbox v-model="permModel.tasks.create" :disabled="permModel.view_only">创建</el-checkbox>
                <el-checkbox v-model="permModel.tasks.edit" :disabled="permModel.view_only">编辑</el-checkbox>
                <el-checkbox v-model="permModel.tasks.delete" :disabled="permModel.view_only">删除</el-checkbox>
                <el-checkbox v-model="permModel.tasks.status" :disabled="permModel.view_only">状态</el-checkbox>
                <el-checkbox v-model="permModel.tasks.undo" :disabled="permModel.view_only">撤销</el-checkbox>
              </div>
            </div>

            <!-- 心愿权限 -->
            <div class="p-4 border border-gray-100 dark:border-gray-700 rounded-xl hover:shadow-sm transition-shadow">
              <div class="font-bold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-2">
                <el-icon class="text-purple-500"><Present /></el-icon>心愿权限
              </div>
              <div class="grid grid-cols-2 gap-2">
                <el-checkbox v-model="permModel.wishes.create" :disabled="permModel.view_only">创建</el-checkbox>
                <el-checkbox v-model="permModel.wishes.edit" :disabled="permModel.view_only">编辑</el-checkbox>
                <el-checkbox v-model="permModel.wishes.delete" :disabled="permModel.view_only">删除</el-checkbox>
                <el-checkbox v-model="permModel.wishes.exchange" :disabled="permModel.view_only">兑换</el-checkbox>
              </div>
            </div>

            <!-- 设置权限 -->
            <div class="p-4 border border-gray-100 dark:border-gray-700 rounded-xl hover:shadow-sm transition-shadow">
              <div class="font-bold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-2">
                <el-icon class="text-orange-500"><Setting /></el-icon>设置访问
              </div>
              <div class="grid grid-cols-2 gap-2">
                <el-checkbox v-model="permModel.settings.tomato">番茄钟</el-checkbox>
                <el-checkbox v-model="permModel.settings.tasks" :disabled="permModel.view_only">任务</el-checkbox>
                <el-checkbox v-model="permModel.settings.coins" :disabled="permModel.view_only">金币</el-checkbox>
                <el-checkbox v-model="permModel.settings.reading">朗读</el-checkbox>
              </div>
            </div>
          </div>
        </div>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-3 pt-4 border-t border-gray-100 dark:border-gray-800">
          <el-button @click="dialogVisible=false" class="!rounded-lg">取消</el-button>
          <el-button type="primary" @click="submit" class="!rounded-lg px-6">{{ isEdit ? '保存修改' : '立即创建' }}</el-button>
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
import { presignView } from '@/services/storage'
import defaultAvatar from '@/assets/avatars/default.png'
import { usePermissions } from '@/composables/permissions'
import { listSubAccounts, createSubAccount, updateSubAccount, deleteSubAccount, generateChildToken, type ChildAccount } from '@/services/subaccounts'
import { prepareUpload } from '@/utils/image'
import { uploadAvatar } from '@/services/user'
import { ArrowLeft, User, Plus, MoreFilled, Key, Edit, Delete, CopyDocument, Lock, List, Present, Setting } from '@element-plus/icons-vue'

// 中文注释：权限校验（家长允许；子账号需具备 account.manage_children）
const { isParent, manageChildren } = usePermissions()
const allowManage = computed(() => isParent.value || manageChildren.value)

const children = ref<ChildAccount[]>([])
const loading = ref(false)
const expiresMap = ref<Record<number, number>>({})
const isMobile = ref(false)
onMounted(() => {
  const update = () => { isMobile.value = window.innerWidth < 768 }
  update()
  window.addEventListener('resize', update)
})

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
  tasks: { create: boolean; edit: boolean; delete: boolean; status: boolean; undo: boolean }
  wishes: { create: boolean; edit: boolean; delete: boolean; exchange: boolean }
  // 中文注释：新增设置权限分组（用于控制“我的”页设置入口权限）
  settings: { tomato: boolean; tasks: boolean; coins: boolean; reading: boolean }
}
const defaultPerms: PermModel = {
  view_only: true,
  account: { manage_children: false },
  tasks: { create: false, edit: false, delete: false, status: false, undo: false },
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

const avatarsMap = ref<Record<number, string>>({})
const avatarDialogSrc = ref<string>(defaultAvatar)
function updateDialogAvatar() {
  const p = currentAvatarPath.value
  if (!p) { avatarDialogSrc.value = defaultAvatar; return }
  const s = String(p)
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { avatarDialogSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { avatarDialogSrc.value = s; return }
  const base = getStaticBase()
  if (/uploads\//i.test(s)) { avatarDialogSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { presignView(s).then(u => avatarDialogSrc.value = u).catch(() => { avatarDialogSrc.value = defaultAvatar }) } catch { avatarDialogSrc.value = defaultAvatar }
}
async function resolveListAvatars(list: ChildAccount[]) {
  const base = getStaticBase()
  for (const c of list) {
    const s = String(c.avatar_path || '')
    if (!s || /storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { avatarsMap.value[c.id] = defaultAvatar; continue }
    if (/^https?:\/\//i.test(s)) { avatarsMap.value[c.id] = s; continue }
    if (/uploads\//i.test(s)) { avatarsMap.value[c.id] = `${base}/api/${s.replace(/^\/+/, '')}`; continue }
    try { avatarsMap.value[c.id] = await presignView(s) } catch { avatarsMap.value[c.id] = defaultAvatar }
  }
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
      m.tasks = { create: true, edit: true, delete: true, status: true, undo: true }
      break
    case 'full':
      m.view_only = false
      m.account.manage_children = true
      m.tasks = { create: true, edit: true, delete: true, status: true, undo: true }
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
    if (!s) { permModel.value = JSON.parse(JSON.stringify(defaultPerms)); return }
    const obj = JSON.parse(s)
    const m = JSON.parse(JSON.stringify(defaultPerms))
    m.view_only = !!obj.view_only
    if (obj.account && typeof obj.account.manage_children === 'boolean') m.account.manage_children = obj.account.manage_children
    if (obj.tasks) {
      m.tasks.create = !!obj.tasks.create
      m.tasks.edit = !!obj.tasks.edit
      m.tasks.delete = !!obj.tasks.delete
      m.tasks.status = !!obj.tasks.status
      m.tasks.undo = !!obj.tasks.undo
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
    ;(['create','edit','delete','status','undo'] as const).forEach(k => { if ((m.tasks as any)[k]) t[k] = true })
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
    await resolveListAvatars(children.value)
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
  updateDialogAvatar()
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
  updateDialogAvatar()
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
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out forwards;
}

:deep(.custom-tiny-select .el-input__wrapper) {
  padding: 0 4px !important;
  box-shadow: none !important;
  background: transparent !important;
}

:deep(.custom-dialog .el-dialog__body) {
  padding: 20px 24px !important;
}

:deep(.el-checkbox.is-bordered) {
  border-radius: 8px;
  transition: all 0.2s;
}
:deep(.el-checkbox.is-bordered.is-checked) {
  background-color: #eff6ff;
  border-color: #3b82f6;
}
:deep(.dark .el-checkbox.is-bordered.is-checked) {
  background-color: rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
}
</style>
