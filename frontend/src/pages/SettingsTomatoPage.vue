<template>
  <SettingsShell title="番茄钟设置" subtitle="专注时长、显示方式与常亮" :icon="Timer" tone="red">
    <SettingsCard title="计时模式" description="选择倒计时或正计时显示方式">
      <el-radio-group v-model="mode" class="!flex !flex-wrap gap-x-6 gap-y-2">
        <el-radio label="countdown">倒计时</el-radio>
        <el-radio label="countup">正计时</el-radio>
      </el-radio-group>
    </SettingsCard>

    <SettingsCard title="预设时长（分钟）" description="进入番茄钟页面时默认使用的计划时长">
      <div class="flex items-center justify-between gap-3">
        <div class="text-xs text-gray-500 dark:text-gray-400">范围 1～180 分钟</div>
        <el-input-number v-model="duration" :min="1" :max="180" controls-position="right" />
      </div>
    </SettingsCard>

    <SettingsCard title="固定番茄钟页面" description="开启后：进入番茄钟为固定全屏，适合专注模式">
      <template #right>
        <el-switch v-model="fixed" />
      </template>
    </SettingsCard>

    <SettingsCard title="保持常亮" description="移动端开始计时后保持屏幕常亮">
      <template #right>
        <el-switch v-model="keepAwake" />
      </template>
    </SettingsCard>

    <SettingsCard title="倒计时结束提示语" description="倒计时自然结束时将使用 TTS 播报该文本">
      <template #right>
        <el-switch v-model="countdownEndSpeakEnabled" />
      </template>
      <el-input
        v-model="countdownEndText"
        type="textarea"
        :rows="3"
        placeholder="例如：倒计时结束，宝贝，休息会儿吧！"
      />
    </SettingsCard>

    <SettingsCard>
      <div class="flex justify-end gap-2">
        <el-button class="!rounded-xl" @click="cancel">取消</el-button>
        <el-button type="primary" class="!rounded-xl" @click="confirm">确定</el-button>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { Timer } from '@element-plus/icons-vue'
import router from '@/router'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

// 读取并绑定番茄钟设置到本页面控件
const store = useAppState()
const mode = ref(store.tomato.mode)
const duration = ref(store.tomato.durationMinutes)
const fixed = ref(store.tomato.fixedTomatoPage)
const keepAwake = ref(store.tomato.keepAwakeEnabled)
const countdownEndText = ref(store.tomato.countdownEndText)
const countdownEndSpeakEnabled = ref(store.tomato.countdownEndSpeakEnabled)

function cancel() {
  // 中文注释：取消不保存设置，关闭页面
  router.back()
}

function confirm() {
  // 中文注释：确定按钮保存设置并关闭页面
  store.updateTomato({
    mode: mode.value as 'countdown' | 'countup',
    runningMode: null,
    durationMinutes: duration.value,
    fixedTomatoPage: fixed.value,
    keepAwakeEnabled: keepAwake.value,
    countdownEndText: countdownEndText.value,
    countdownEndSpeakEnabled: countdownEndSpeakEnabled.value
  })
  router.back()
}
</script>

<style scoped>
</style>
