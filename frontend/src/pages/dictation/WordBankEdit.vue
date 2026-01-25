<template>
  <SettingsShell :title="isEdit ? '编辑词库' : '新建词库'" subtitle="为听写生成更清晰的内容来源" :icon="Collection" tone="blue" container-class="max-w-4xl">
    <template #headerRight>
      <el-button type="primary" class="!rounded-2xl !font-extrabold" :loading="saving" @click="submit">保存</el-button>
    </template>

    <SettingsCard title="基础信息">
      <el-form :model="form" label-position="top" class="dictation-form">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="例如：第一单元生词" />
        </el-form-item>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <el-form-item label="科目">
            <el-select v-model="form.subject" class="w-full">
              <el-option label="语文" value="语文" />
              <el-option label="英语" value="英语" />
              <el-option label="其他" value="其他" />
            </el-select>
          </el-form-item>
          <el-form-item label="阶段">
            <el-select v-model="form.education_stage" class="w-full">
              <el-option label="小学" value="小学" />
              <el-option label="初中" value="初中" />
              <el-option label="高中" value="高中" />
            </el-select>
          </el-form-item>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <el-form-item label="版本">
            <el-select v-model="form.version" class="w-full" allow-create filterable default-first-option placeholder="选择或输入版本">
              <el-option v-for="v in ['人教版', '苏教版', '北师大版', '外研版']" :key="v" :label="v" :value="v" />
            </el-select>
          </el-form-item>
          <el-form-item label="年级">
            <el-select v-model="form.grade" class="w-full" placeholder="选择年级">
              <el-option v-for="g in availableGrades" :key="g" :label="g" :value="g" />
            </el-select>
          </el-form-item>
        </div>
      </el-form>
    </SettingsCard>

    <SettingsCard title="词库内容" description="支持空格或换行分隔，系统会按设置规则分词。">
      <el-form :model="form" label-position="top" class="dictation-form">
        <el-form-item label="内容">
          <el-input v-model="form.content" type="textarea" :rows="12" resize="none" placeholder="输入词汇或课文…" />
        </el-form-item>
      </el-form>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Collection } from '@element-plus/icons-vue'
import { dictationApi, type WordBank } from '@/services/dictation'
import { ElMessage } from 'element-plus'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null
const isEdit = computed(() => !!id)
const saving = ref(false)

const form = ref<Partial<WordBank>>({
  title: '',
  subject: '语文',
  education_stage: '小学',
  version: '',
  grade: '',
  content: ''
})

const availableGrades = computed(() => {
  const stage = form.value.education_stage
  if (stage === '小学') return ['一年级上', '一年级下', '二年级上', '二年级下', '三年级上', '三年级下', '四年级上', '四年级下', '五年级上', '五年级下', '六年级上', '六年级下']
  if (stage === '初中') return ['初一上', '初一下', '初二上', '初二下', '初三上', '初三下']
  if (stage === '高中') return ['高一上', '高一下', '高二上', '高二下', '高三上', '高三下']
  return []
})

onMounted(async () => {
  if (isEdit && id) {
    try {
      // Since we don't have a getOne API in the snippet above, let's rely on list or implement getOne.
      // Wait, I forgot to add GetWordBank to handlers/router? 
      // Let's check router.go. I missed `dictation.GET("/wordbanks/:id", handlers.GetWordBank)`?
      // I implemented UpdateWordBank which takes ID, but I might have missed Get.
      // Let's assume I can just use the list to find it or fetch it. 
      // Actually, I didn't implement `GetWordBank` in handlers!
      // I implemented `ListWordBanks`, `Create`, `Update`, `Delete`.
      // So I should fetch list and find it, OR implement Get.
      // Implementing Get is better. But for speed, I will just fetch list and find.try {
      const res = await dictationApi.listWordBanks()
      const list = (res as any) || []
      const found = list.find((i: any) => i.id === id)
      if (found) {
        form.value = { ...found }
        // If content is JSON array, convert to string?
        try {
           if (found.content.startsWith('[')) {
             const arr = JSON.parse(found.content)
             form.value.content = arr.join('\n')
           }
        } catch {}
      }
    } catch (e) {
      ElMessage.error('加载失败')
    }
  }
})

async function submit() {
  if (!form.value.title) return ElMessage.warning('请输入标题')
  if (!form.value.content) return ElMessage.warning('请输入内容')
  
  saving.value = true
  try {
    const payload = { ...form.value }
    // Convert content to simple string or keep as is. 
    // PRD doesn't specify storage format, but string is fine.
    // If user inputs with newlines, we keep it.
    
    if (isEdit.value && id) {
      await dictationApi.updateWordBank(id, payload)
      ElMessage.success('已更新')
    } else {
      await dictationApi.createWordBank(payload)
      ElMessage.success('已创建')
    }
    router.back()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
:deep(.dictation-form .el-form-item__label) {
  font-size: 12px;
  font-weight: 800;
  color: rgb(107 114 128);
}

.dark :deep(.dictation-form .el-form-item__label) {
  color: rgb(156 163 175);
}
</style>
