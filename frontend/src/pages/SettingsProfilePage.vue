<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 pb-20 transition-colors duration-300">
    <!-- 顶部导航栏 -->
    <div class="sticky top-0 z-30 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800 px-4 h-14 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button 
          @click="onCancel"
          class="p-2 -ml-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-600 dark:text-gray-300">
          <el-icon :size="20"><ArrowLeft /></el-icon>
        </button>
        <span class="font-bold text-lg text-gray-800 dark:text-gray-100">编辑个人信息</span>
      </div>
      <el-button type="primary" size="small" @click="onSave" class="!rounded-lg !px-4">保存</el-button>
    </div>

    <div class="max-w-3xl mx-auto p-4 md:p-6 space-y-6 animate-fade-in-up">
      <!-- 头像上传 -->
      <div class="flex flex-col items-center py-6">
        <div class="relative group cursor-pointer" @click="triggerAvatarPick">
          <div class="absolute -inset-1 bg-gradient-to-br from-blue-500 to-cyan-500 rounded-full blur opacity-30 group-hover:opacity-60 transition duration-500"></div>
          <el-avatar :size="100" :src="avatarPreview || avatarSrc" class="relative border-4 border-white dark:border-gray-900 shadow-xl transition-transform duration-300 group-hover:scale-105" />
          <div class="absolute bottom-0 right-0 bg-blue-500 text-white p-2 rounded-full shadow-lg border-2 border-white dark:border-gray-900 transition-transform duration-300 group-hover:scale-110">
            <el-icon :size="16"><Camera /></el-icon>
          </div>
          <input ref="avatarInput" type="file" accept="image/*" class="hidden" @change="onSelectAvatar" />
        </div>
        <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">点击更换头像</p>
      </div>

      <!-- 基本信息卡片 -->
      <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-50 dark:border-gray-800 flex items-center gap-2">
           <div class="w-1 h-4 bg-blue-500 rounded-full"></div>
           <span class="font-bold text-gray-800 dark:text-gray-200">基本资料</span>
        </div>
        <div class="p-5 grid grid-cols-1 md:grid-cols-2 gap-5">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">昵称</label>
            <el-input v-model="editNickname" placeholder="设置昵称" size="large" class="custom-input">
              <template #prefix><el-icon class="text-gray-400"><User /></el-icon></template>
            </el-input>
          </div>
          
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">性别</label>
            <el-select v-model="editChildGender" placeholder="选择性别" size="large" class="w-full custom-select">
              <template #prefix><el-icon class="text-gray-400"><Male /></el-icon></template>
              <el-option label="男" value="male" />
              <el-option label="女" value="female" />
            </el-select>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">生日</label>
            <el-date-picker 
              v-model="editChildBirthday" 
              type="date" 
              value-format="YYYY-MM-DD" 
              placeholder="选择生日" 
              size="large" 
              class="!w-full custom-date-picker" 
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">电话</label>
            <el-input v-model="editPhone" placeholder="绑定手机号" size="large" class="custom-input">
              <template #prefix><el-icon class="text-gray-400"><Iphone /></el-icon></template>
            </el-input>
          </div>

          <div class="space-y-1.5 md:col-span-2">
            <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">邮箱</label>
            <el-input v-model="editEmail" placeholder="绑定邮箱" size="large" class="custom-input">
              <template #prefix><el-icon class="text-gray-400"><Message /></el-icon></template>
            </el-input>
          </div>
        </div>
      </div>

      <!-- 安全设置卡片 -->
      <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-50 dark:border-gray-800 flex items-center gap-2">
           <div class="w-1 h-4 bg-red-500 rounded-full"></div>
           <span class="font-bold text-gray-800 dark:text-gray-200">修改密码</span>
           <span class="text-xs text-gray-400 font-normal ml-auto">留空则不修改</span>
        </div>
        <div class="p-5 space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">当前密码</label>
            <el-input v-model="oldPassword" type="password" show-password placeholder="请输入当前使用的密码" size="large" class="custom-input">
               <template #prefix><el-icon class="text-gray-400"><Lock /></el-icon></template>
            </el-input>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">新密码</label>
              <el-input v-model="newPassword" type="password" show-password placeholder="设置新密码" size="large" class="custom-input">
                <template #prefix><el-icon class="text-gray-400"><Key /></el-icon></template>
              </el-input>
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-gray-500 dark:text-gray-400 ml-1">确认新密码</label>
              <el-input v-model="confirmPassword" type="password" show-password placeholder="再次输入新密码" size="large" class="custom-input">
                <template #prefix><el-icon class="text-gray-400"><Check /></el-icon></template>
              </el-input>
            </div>
          </div>
        </div>
      </div>

      <!-- 管理员入口 -->
      <div v-if="isAdmin" class="flex justify-center pt-4">
        <button 
          @click="toAdmin" 
          class="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors font-medium text-sm">
          <el-icon><Setting /></el-icon>
          <span>进入后台管理系统</span>
        </button>
      </div>
      
      <!-- 底部留白 -->
      <div class="h-8"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '@/stores/auth'
import { updateProfile, changePassword, uploadAvatar } from '@/services/user'
import { ElMessage } from 'element-plus'
import { Edit, ArrowLeft, Camera, User, Iphone, Message, Lock, Key, Check, Male, Setting } from '@element-plus/icons-vue'
import router from '@/router'
import { computed, ref, onMounted, watch } from 'vue'
import defaultAvatar from '@/assets/avatars/default.png'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import { prepareUpload } from '@/utils/image'

const auth = useAuth()
const isAdmin = computed(() => Number(auth.user?.id || 0) === 1)
function toAdmin() { router.push('/admin') }

const avatarSrc = ref<string>(defaultAvatar)
async function updateAvatarSrc() {
  const p = auth.user?.avatar_path
  if (!p) { avatarSrc.value = defaultAvatar; return }
  const s = String(p)
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { avatarSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { avatarSrc.value = s; return }
  const base = getStaticBase()
  if (/uploads\//i.test(s)) { avatarSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { avatarSrc.value = await presignView(s) } catch { avatarSrc.value = defaultAvatar }
}
onMounted(updateAvatarSrc)
watch(() => auth.user?.avatar_path, () => { updateAvatarSrc() })

// 字段状态
const editNickname = ref(auth.user?.nickname || '')
const editPhone = ref(auth.user?.phone || '')
const editEmail = ref(auth.user?.email || '')
const editChildBirthday = ref(auth.user?.child_birthday || '')
const editChildGender = ref(auth.user?.child_gender || '')
const avatarInput = ref<HTMLInputElement | null>(null)
const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string>('')
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

function triggerAvatarPick() {
  avatarInput.value?.click()
}

async function onSelectAvatar(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files && input.files[0]
  if (!f) return
  const webp = await prepareUpload(f)
  avatarFile.value = webp
  avatarPreview.value = URL.createObjectURL(webp)
}

// 中文注释：取消返回（不保存，仅关闭页面）
function onCancel() {
  router.back()
}

// 中文注释：显式保存按钮逻辑，保存成功后关闭页面
async function onSave() {
  try {
    // 1) 资料更新：昵称/电话/邮箱（有变更时一起提交）
    const nicknameTrim = (editNickname.value || '').trim()
    const phoneTrim = (editPhone.value || '').trim()
    const emailTrim = (editEmail.value || '').trim()
    const payload: { nickname?: string; phone?: string; email?: string; child_birthday?: string; child_gender?: string } = {}
    if (nicknameTrim && nicknameTrim !== (auth.user?.nickname || '')) payload.nickname = nicknameTrim
    if (phoneTrim !== (auth.user?.phone || '')) payload.phone = phoneTrim
    if (emailTrim !== (auth.user?.email || '')) payload.email = emailTrim
    if ((editChildBirthday.value || '') !== (auth.user?.child_birthday || '')) payload.child_birthday = editChildBirthday.value || ''
    if ((editChildGender.value || '') !== (auth.user?.child_gender || '')) payload.child_gender = editChildGender.value || ''
    if (Object.keys(payload).length > 0) {
      await updateProfile(payload)
      auth.updateUser(payload)
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
      auth.updateUser({ must_change_password: false })
    }

    ElMessage.success('个人信息已保存')
    router.back()
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  }
}
</script>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out forwards;
}

:deep(.el-input__wrapper) {
  box-shadow: none !important;
  background-color: #f9fafb !important;
  border: 1px solid #e5e7eb !important;
  border-radius: 0.75rem !important;
  transition: all 0.2s;
}
:deep(.dark .el-input__wrapper) {
  background-color: #1f2937 !important;
  border-color: #374151 !important;
}
:deep(.el-input__wrapper:hover) {
  border-color: #d1d5db !important;
}
:deep(.el-input__wrapper.is-focus) {
  border-color: #3b82f6 !important;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1) !important;
}
:deep(.el-select__wrapper) {
  box-shadow: none !important;
  background-color: #f9fafb !important;
  border: 1px solid #e5e7eb !important;
  border-radius: 0.75rem !important;
}
:deep(.dark .el-select__wrapper) {
  background-color: #1f2937 !important;
  border-color: #374151 !important;
}
</style>
