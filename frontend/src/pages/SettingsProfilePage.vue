<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <!-- 中文注释：返回按钮不再自动保存，仅关闭页面 -->
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="onCancel"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#60a5fa"><Edit /></el-icon>
      <h2 class="font-semibold">编辑个人信息</h2>
    </div>

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
      <!-- 中文注释：取消 webp 自动压缩与回退提示，按原格式上传 -->
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

    <!-- 底部右侧：取消 / 保存 按钮 -->
    <div class="flex justify-end mt-6 gap-2">
      <el-button @click="onCancel">取消</el-button>
      <el-button type="primary" @click="onSave">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '@/stores/auth'
import { updateNickname, changePassword, uploadAvatar } from '@/services/user'
import { ElMessage } from 'element-plus'
import { Edit, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
import defaultAvatar from '@/assets/avatars/default.png'

const auth = useAuth()

function resolveAvatarUrl(p?: string | null) {
  if (!p) return defaultAvatar
  const s = String(p)
  if (/^https?:\/\//i.test(s)) return s
  if (!/uploads\//i.test(s)) return defaultAvatar
  // 中文注释：为相对路径统一加上 /api 前缀并规范化多余斜杠
  return `/api/${s}`.replace(/\/+/g, '/').replace(/\/$/, '')
}

const avatarSrc = computed(() => resolveAvatarUrl(auth.user?.avatar_path))

// 字段状态
const editNickname = ref(auth.user?.nickname || '')
const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string>('')
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

function onSelectAvatar(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files && input.files[0]
  if (!f) return
  avatarFile.value = f
  avatarPreview.value = URL.createObjectURL(f)
}

// 中文注释：取消返回（不保存，仅关闭页面）
function onCancel() {
  router.back()
}

// 中文注释：显式保存按钮逻辑，保存成功后关闭页面
async function onSave() {
  try {
    // 1) 昵称更新（非空且有变更）
    const nicknameTrim = (editNickname.value || '').trim()
    if (nicknameTrim && nicknameTrim !== (auth.user?.nickname || '')) {
      await updateNickname(nicknameTrim)
      auth.updateUser({ nickname: nicknameTrim })
    }

    // 2) 头像上传（有选择时），按原文件格式上传（取消 webp 自动转换）
    if (avatarFile.value && auth.user) {
      const resp = await uploadAvatar(auth.user.id, avatarFile.value)
      const path = resp.path
      auth.updateUser({ avatar_path: path })
    }

    // 3) 修改密码（若填写）
    if (oldPassword.value || newPassword.value || confirmPassword.value) {
      if (!oldPassword.value || !newPassword.value || !confirmPassword.value) {
        ElMessage.error('请完整填写修改密码字段');
        return
      }
      if (newPassword.value !== confirmPassword.value) {
        ElMessage.error('两次输入的新密码不一致');
        return
      }
      await changePassword(oldPassword.value, newPassword.value)
    }

    ElMessage.success('个人信息已保存')
    router.back()
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  }
}
</script>

<style scoped>
</style>