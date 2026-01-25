<template>
  <SettingsShell title="难点收藏" subtitle="把易错词收集起来，反复巩固" :icon="Notebook" tone="red" container-class="max-w-4xl">
    <template #headerRight>
      <el-button type="primary" class="!rounded-2xl !font-extrabold" @click="addManual">手动添加</el-button>
    </template>

    <SettingsCard>
      <el-empty v-if="list.length === 0" description="暂无错题" />

      <div v-else class="space-y-3">
        <div
          v-for="item in list"
          :key="item.id"
          class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition flex items-start justify-between gap-3"
        >
          <div class="min-w-0">
            <div class="text-lg font-extrabold tracking-tight text-gray-900 dark:text-gray-50">{{ item.word }}</div>
            <div v-if="item.context" class="mt-1 text-sm text-gray-600 dark:text-gray-300">来源：{{ item.context }}</div>
            <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ new Date(item.created_at).toLocaleDateString() }}</div>
          </div>
          <el-button circle type="danger" class="!rounded-2xl" :icon="Delete" size="small" @click="removeItem(item.id)" />
        </div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Delete, Notebook } from '@element-plus/icons-vue'
import { dictationApi, type Mistake } from '@/services/dictation'
import { ElMessage, ElMessageBox } from 'element-plus'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const router = useRouter()
const list = ref<Mistake[]>([])

onMounted(loadData)

async function loadData() {
  try {
    const res = await dictationApi.listMistakes()
    list.value = (res as any) || []
  } catch (e) {
    ElMessage.error('加载失败')
  }
}

async function removeItem(id: number) {
  try {
    await dictationApi.deleteMistake(id)
    list.value = list.value.filter(i => i.id !== id)
    ElMessage.success('已移除')
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

async function addManual() {
  try {
    const { value } = await ElMessageBox.prompt('请输入错词', '手动添加', {
      confirmButtonText: '添加',
      cancelButtonText: '取消',
    })
    if (value) {
      await dictationApi.addMistake({ word: value, context: '手动添加' })
      loadData()
    }
  } catch {}
}
</script>
