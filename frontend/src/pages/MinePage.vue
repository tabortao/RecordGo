<template>
  <!-- 中文注释：我的页面骨架，展示昵称、ID、头像占位 -->
  <div class="p-4 space-y-4">
    <el-card shadow="never">
      <div class="flex items-center gap-3">
        <!-- 中文注释：头像优先显示用户自定义头像，未设置则使用默认头像；修复相对路径未加 API 前缀导致无法显示的问题 -->
        <el-avatar :size="48" :src="avatarSrc" />
        <div>
          <!-- 中文注释：昵称优先显示真实昵称，未设置则显示用户名 -->
          <div class="font-semibold">{{ displayName }}</div>
          <!-- 中文注释：用户ID格式化为6位数，不足左侧补0 -->
          <div class="text-gray-500 text-sm">用户 ID：{{ displayId }}</div>
        </div>
      </div>
    </el-card>

    <!-- 中文注释：账号管理 UI（shadcn 风格 + tailwind），包含编辑个人信息/子账号管理/退出登录 -->
    <div class="rounded-lg border bg-white">
      <div class="px-4 py-3 flex items-center gap-2">
        <el-icon :size="18" style="color:#3b82f6"><User /></el-icon>
        <span class="font-semibold">账号管理</span>
      </div>
      <div class="px-2 py-2 space-y-1">
        <!-- 编辑个人信息 -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition" @click="profileDialogVisible = true">
          <el-icon :size="18" style="color:#60a5fa"><Edit /></el-icon>
          <span class="text-gray-800">编辑个人信息</span>
        </button>
        <!-- 子账号管理（占位） -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition" @click="onChildManage">
          <el-icon :size="18" style="color:#22c55e"><User /></el-icon>
          <span class="text-gray-800">子账号管理</span>
        </button>
        <!-- 退出登录 -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition" @click="onLogout">
          <el-icon :size="18" style="color:#ef4444"><SwitchButton /></el-icon>
          <span class="text-gray-800">退出登录</span>
        </button>
      </div>
    </div>

    <!-- 中文注释：编辑个人信息对话框（包含昵称、头像上传、修改密码），样式采用 tailwind，遵循接口一致性与图片 webp 约束 -->
    <el-dialog v-model="profileDialogVisible" :width="dialogWidth" :show-close="true" title="编辑个人信息">
      <div class="space-y-4">
        <!-- 昵称 -->
        <div class="space-y-2">
          <label class="text-sm text-gray-600">昵称</label>
          <el-input v-model="editNickname" placeholder="请输入昵称" />
        </div>
        <!-- 头像上传与预览 -->
        <div class="space-y-2">
          <label class="text-sm text-gray-600">头像</label>
          <div class="flex items-center gap-4">
            <el-avatar :size="56" :src="avatarPreview || avatarSrc" />
            <input type="file" accept="image/*" @change="onSelectAvatar" />
          </div>
          <div class="text-xs text-gray-500">将自动压缩并转换为 webp，失败则回退原格式</div>
        </div>
        <!-- 修改密码 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <div class="space-y-2">
            <label class="text-sm text-gray-600">当前密码</label>
            <el-input v-model="oldPassword" type="password" placeholder="请输入当前密码" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600">新密码</label>
            <el-input v-model="newPassword" type="password" placeholder="请输入新密码" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600">确认新密码</label>
            <el-input v-model="confirmPassword" type="password" placeholder="请再次输入新密码" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="profileDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveProfile">保存</el-button>
        </div>
      </template>
    </el-dialog>

    <el-card shadow="never">
      <template #header>
        <span>系统设置</span>
      </template>
      <el-form label-width="160px">
        <el-form-item label="固定番茄钟页面">
          <el-switch v-model="fixed" @change="onFixedChange" />
        </el-form-item>
        
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { computed, ref, watchEffect } from 'vue'
import defaultAvatar from '@/assets/avatars/default.png'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { prepareUpload } from '@/utils/image'
import { updateNickname, changePassword, uploadAvatar } from '@/services/user'
import { User, Edit, SwitchButton } from '@element-plus/icons-vue'

// 中文注释：固定番茄设置持久化到 appState（后续联动后端 UserSettings）
const store = useAppState()
const fixed = computed({
  get: () => store.tomato.fixedTomatoPage,
  set: (v: boolean) => store.updateTomato({ fixedTomatoPage: v })
})

function onFixedChange() {
  // 中文注释：此处可调用后端保存设置
}

// 中文注释：退出登录，清除认证信息并跳转到登录页
const auth = useAuth()

// 中文注释：展示名称逻辑——优先真实昵称，缺省则回退到用户名
const displayName = computed(() => {
  const u = auth.user
  if (!u) return '未登录'
  const n = (u.nickname || '').trim()
  return n ? n : u.username
})

// 中文注释：展示ID为6位数字字符串（不足位数左侧补0）
const displayId = computed(() => {
  const u = auth.user
  if (!u) return '------'
  return String(u.id).padStart(6, '0')
})

// 中文注释：头像地址优先使用用户头像路径，未设置则用默认头像
// 中文注释：头像地址解析——支持完整 URL；后端返回的相对路径（如 uploads/images/avatars/...）统一加上 /api 前缀通过 Vite 代理
function resolveAvatarUrl(p?: string | null) {
  if (!p) return defaultAvatar
  const s = String(p)
  // 中文注释：仅当为完整 URL 或包含 uploads 路径时才走后端；否则回退到内置默认头像
  if (/^https?:\/\//i.test(s)) return s
  if (!/uploads\//i.test(s)) return defaultAvatar
  return `/api/${s}`.replace(/\/+/g, '/').replace(/\/$/, '')
}

const avatarSrc = computed(() => resolveAvatarUrl(auth.user?.avatar_path))

// 中文注释：编辑个人信息弹窗状态与字段
const profileDialogVisible = ref(false)
const editNickname = ref('')
const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string>('')
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const isMobile = ref(false)
const dialogWidth = computed(() => (isMobile.value ? '96vw' : '560px'))

// 中文注释：打开编辑弹窗时，初始化字段
watchEffect(() => {
  if (profileDialogVisible.value) {
    editNickname.value = auth.user?.nickname || ''
    avatarFile.value = null
    avatarPreview.value = ''
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  }
})

function onSelectAvatar(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files && input.files[0]
  if (!f) return
  avatarFile.value = f
  avatarPreview.value = URL.createObjectURL(f)
}

async function saveProfile() {
  try {
    // 1) 更新昵称（非空且有变更时）
    const nicknameTrim = (editNickname.value || '').trim()
    if (nicknameTrim && nicknameTrim !== (auth.user?.nickname || '')) {
      await updateNickname(nicknameTrim)
      auth.updateUser({ nickname: nicknameTrim })
    }

    // 2) 上传头像（有选择时），前端先转换为 webp
    if (avatarFile.value && auth.user) {
      const webp = await prepareUpload(avatarFile.value)
      const resp = await uploadAvatar(auth.user.id, webp)
      const path = resp.path
      auth.updateUser({ avatar_path: path })
    }

    // 3) 修改密码（若填写）
    if (oldPassword.value || newPassword.value || confirmPassword.value) {
      if (!oldPassword.value || !newPassword.value || !confirmPassword.value) {
        ElMessage.error('请完整填写修改密码字段')
        return
      }
      if (newPassword.value !== confirmPassword.value) {
        ElMessage.error('两次输入的新密码不一致')
        return
      }
      await changePassword(oldPassword.value, newPassword.value)
    }

    ElMessage.success('资料已保存')
    profileDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  }
}

function onChildManage() {
  ElMessage.info('子账号管理功能将在后续版本提供')
}

function onLogout() {
  try {
    auth.logout()
    store.reset()
    router.replace('/login')
  } catch (_) {
    location.reload()
  }
}

// 中文注释：导航到登录与注册占位页面
function toLogin() {
  router.push('/login')
}

function toRegister() {
  router.push('/register')
}
</script>

<style scoped>
</style>
