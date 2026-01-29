<template>
  <SettingsShell title="用户管理后台" subtitle="仅管理员可见" :icon="DataAnalysis" tone="indigo" container-class="max-w-6xl">
    <SettingsCard>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <div class="rounded-2xl border border-sky-100/70 dark:border-gray-800 bg-sky-50/70 dark:bg-gray-900/50 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-sky-700/80 dark:text-sky-300/80">注册用户总数</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-sky-800 dark:text-sky-200">{{ overview.total_users }}</div>
        </div>
        <div class="rounded-2xl border border-emerald-100/70 dark:border-gray-800 bg-emerald-50/70 dark:bg-gray-900/50 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-emerald-700/80 dark:text-emerald-300/80">今日登录数</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-emerald-800 dark:text-emerald-200">{{ overview.today_logins }}</div>
        </div>
      </div>
      <div class="mt-3 flex justify-end">
        <el-button class="!rounded-xl" type="success" @click="router.push('/admin/login-stats')">查看登录统计</el-button>
      </div>
    </SettingsCard>

    <SettingsCard title="用户列表" description="支持搜索用户名 / 昵称 / 邮箱">
      <div class="flex flex-col sm:flex-row sm:items-center gap-2">
        <el-input v-model="query" placeholder="搜索..." class="w-full sm:w-72" @keyup.enter="loadUsers" />
        <div class="flex items-center gap-2">
          <el-button class="!rounded-xl" @click="loadUsers">查询</el-button>
          <el-button class="!rounded-xl" @click="loadOverview">刷新概览</el-button>
        </div>
      </div>

      <div class="mt-4 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800">
        <div class="w-full overflow-x-auto">
          <el-table
            :data="users"
            size="small"
            stripe
            class="min-w-[1180px]"
            :header-cell-style="{ whiteSpace: 'nowrap' }"
            :cell-style="{ whiteSpace: 'nowrap' }"
          >
          <el-table-column label="ID" prop="id" width="80" />
          <el-table-column label="用户名" prop="username" min-width="140" />
          <el-table-column label="昵称" prop="nickname" min-width="140" />
          <el-table-column label="邮箱" prop="email" min-width="190" />
          <el-table-column label="电话" prop="phone" min-width="130" />
          <el-table-column label="活跃标记" width="140">
            <template #default="{ row }">
              <div class="flex items-center gap-1">
                <el-tag v-if="row.login_today" type="success" effect="dark">今</el-tag>
                <el-tag v-if="row.inactive_30d" type="info">30天未登</el-tag>
                <span v-if="!row.login_today && !row.inactive_30d" class="text-xs text-gray-400">-</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="VIP" width="96">
            <template #default="{ row }">
              <el-tag v-if="row.is_lifetime_vip" type="success">终身</el-tag>
              <el-tag v-else-if="row.is_vip" type="warning">VIP</el-tag>
              <el-tag v-else type="info">普通</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="VIP到期" prop="vip_expire_time" width="180" />
          <el-table-column label="最后登录" prop="last_login_time" width="180" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.is_disabled ? 'danger' : 'success'">{{ row.is_disabled ? '禁用' : '正常' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="420">
            <template #default="{ row }">
              <div class="flex flex-wrap gap-2">
                <el-button size="small" class="!rounded-xl" @click="openVIP(row)">VIP设置</el-button>
                <el-button size="small" class="!rounded-xl" type="primary" @click="resetPwd(row)">重置密码</el-button>
                <el-button size="small" class="!rounded-xl" type="warning" @click="toggleDisabled(row)">{{ row.is_disabled ? '启用' : '禁用' }}</el-button>
                <el-button size="small" class="!rounded-xl" type="danger" @click="removeUser(row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
          </el-table>
        </div>
      </div>
    </SettingsCard>

    <el-dialog v-model="vipDialog.visible" width="420px" class="rounded-2xl overflow-hidden">
      <template #header><div class="font-extrabold">VIP设置（用户ID：{{ vipDialog.user?.id }}）</div></template>
      <div class="space-y-3">
        <el-switch v-model="vipDialog.is_vip" active-text="VIP" inactive-text="普通" />
        <el-switch v-model="vipDialog.is_lifetime_vip" active-text="终身VIP" />
        <el-date-picker v-model="vipDialog.vip_expire_time" type="datetime" placeholder="VIP到期时间" :disabled="vipDialog.is_lifetime_vip" />
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button class="!rounded-xl" @click="vipDialog.visible=false">取消</el-button>
          <el-button class="!rounded-xl" type="primary" @click="saveVIP">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { DataAnalysis } from '@element-plus/icons-vue'
import http from '@/services/http'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const auth = useAuth()
// 仅允许用户ID为1访问
if (!auth.user || Number(auth.user.id) !== 1) {
  router.replace('/tasks')
}

const overview = reactive({ total_users: 0, today_logins: 0 })
const users = ref<any[]>([])
const query = ref('')

async function loadOverview() {
  try {
    const data = await http.get('/admin/overview')
    overview.total_users = Number((data as any)?.total_users || 0)
    overview.today_logins = Number((data as any)?.today_logins || 0)
  } catch (e: any) { ElMessage.error(e?.message || '概览加载失败') }
}
async function loadUsers() {
  try {
    const data = await http.get('/admin/users', { params: { q: query.value } }) as any
    users.value = Array.isArray(data.items) ? data.items : []
  } catch (e: any) { ElMessage.error(e?.message || '用户列表加载失败') }
}
onMounted(async () => { await loadOverview(); await loadUsers() })

const vipDialog = reactive<{ visible: boolean; user: any|null; is_vip: boolean; is_lifetime_vip: boolean; vip_expire_time: string|Date|null }>({ visible: false, user: null, is_vip: false, is_lifetime_vip: false, vip_expire_time: null })
function openVIP(u: any) {
  vipDialog.user = u
  vipDialog.visible = true
  vipDialog.is_vip = !!u.is_vip
  vipDialog.is_lifetime_vip = !!u.is_lifetime_vip
  vipDialog.vip_expire_time = u.vip_expire_time || null
}
async function saveVIP() {
  if (!vipDialog.user) return
  try {
    // 前端表单校验：终身VIP时不允许设置到期时间；非VIP且设置了到期时间提示错误
    if (vipDialog.is_lifetime_vip) {
      vipDialog.vip_expire_time = null
    }
    if (!vipDialog.is_vip && !vipDialog.is_lifetime_vip && vipDialog.vip_expire_time) {
      ElMessage.error('非VIP不应设置到期时间')
      return
    }
    await http.post(`/admin/users/${vipDialog.user.id}/vip`, {
      is_vip: vipDialog.is_vip,
      is_lifetime_vip: vipDialog.is_lifetime_vip,
      vip_expire_time: vipDialog.vip_expire_time ? new Date(vipDialog.vip_expire_time).toISOString() : null
    })
    ElMessage.success('VIP设置已保存')
    vipDialog.visible = false
    await loadUsers()
  } catch (e: any) { console.error('VIP保存失败', e?.response?.data || e); ElMessage.error(e?.message || '保存失败') }
}
async function toggleDisabled(u: any) {
  try {
    const disabled = !u.is_disabled
    await ElMessageBox.confirm(`确认${disabled ? '禁用' : '启用'}该用户？`, '确认', { type: 'warning' })
    await http.post(`/admin/users/${u.id}/disabled`, { disabled })
    ElMessage.success('已更新用户状态')
    await loadUsers()
  } catch {}
}
async function removeUser(u: any) {
  try {
    await ElMessageBox.confirm('确认删除该用户？此操作不可恢复', '警告', { type: 'warning' })
    await http.delete(`/admin/users/${u.id}`)
    ElMessage.success('已删除用户')
    await loadUsers()
  } catch {}
}

async function resetPwd(u: any) {
  try {
    await ElMessageBox.confirm(`确认重置用户 ${u.username} 的密码？`, '确认', { type: 'warning' })
    const data = await http.post(`/admin/users/${u.id}/reset-password`, {}) as any
    const pw = String(data?.temp_password || '')
    if (!pw) { ElMessage.error('重置失败'); return }
    await ElMessageBox.alert(`临时密码：${pw}\n该用户下次登录将被要求修改密码。`, '临时密码', { type: 'success' })
  } catch {}
}
</script>

<style scoped>
</style>
