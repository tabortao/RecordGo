<template>
  <div
    class="relative min-h-screen overflow-x-hidden bg-gray-50 text-gray-900 dark:bg-gray-950 dark:text-gray-50 transition-colors duration-300"
    style="padding-bottom: calc(env(safe-area-inset-bottom) + 96px)"
  >
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div :class="tone.glowA" class="absolute -top-24 -right-24 h-72 w-72 rounded-full blur-3xl opacity-35 dark:opacity-25" />
      <div :class="tone.glowB" class="absolute -bottom-40 -left-28 h-80 w-80 rounded-full blur-3xl opacity-30 dark:opacity-20" />
      <div class="absolute inset-0 bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.65),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.45),transparent_55%)] dark:bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.07),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.06),transparent_55%)]" />
    </div>

    <div class="sticky top-0 z-30 px-4 pt-4">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/75 dark:bg-gray-900/70 backdrop-blur-xl shadow-sm">
        <div class="flex items-center gap-3 px-3 py-3">
          <button
            type="button"
            class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors"
            @click="props.backTo ? router.push(props.backTo) : router.back()"
          >
            <el-icon :size="20"><ArrowLeft /></el-icon>
          </button>

          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <div v-if="icon" :class="tone.iconWrap" class="flex h-10 w-10 items-center justify-center rounded-2xl border shadow-sm">
                <el-icon :size="18" :class="tone.iconText"><component :is="icon" /></el-icon>
              </div>
              <div class="min-w-0">
                <div class="text-[17px] font-extrabold tracking-tight truncate">{{ title }}</div>
                <div v-if="subtitle" class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 truncate">{{ subtitle }}</div>
              </div>
            </div>
          </div>

          <div class="shrink-0">
            <slot name="headerRight" />
          </div>
        </div>
      </div>
    </div>

    <div class="relative z-10 px-4 pt-5 pb-6">
      <div :class="containerClass" class="mx-auto space-y-4">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

type ToneKey = 'red' | 'emerald' | 'indigo' | 'blue' | 'violet' | 'amber' | 'sky'

const props = withDefaults(defineProps<{
  title: string
  subtitle?: string
  icon?: any
  tone?: ToneKey
  containerClass?: string
  backTo?: string
}>(), {
  tone: 'sky',
  containerClass: 'max-w-3xl',
  backTo: ''
})

const router = useRouter()

const tone = computed(() => {
  const map: Record<ToneKey, { glowA: string; glowB: string; iconWrap: string; iconText: string }> = {
    red: {
      glowA: 'bg-red-300/60 dark:bg-red-500/25',
      glowB: 'bg-rose-200/60 dark:bg-rose-500/20',
      iconWrap: 'border-red-200/70 dark:border-red-900/40 bg-red-50/80 dark:bg-red-900/25',
      iconText: 'text-red-600 dark:text-red-300'
    },
    emerald: {
      glowA: 'bg-emerald-300/60 dark:bg-emerald-500/25',
      glowB: 'bg-teal-200/60 dark:bg-teal-500/20',
      iconWrap: 'border-emerald-200/70 dark:border-emerald-900/40 bg-emerald-50/80 dark:bg-emerald-900/25',
      iconText: 'text-emerald-700 dark:text-emerald-300'
    },
    indigo: {
      glowA: 'bg-indigo-300/60 dark:bg-indigo-500/25',
      glowB: 'bg-violet-200/60 dark:bg-violet-500/20',
      iconWrap: 'border-indigo-200/70 dark:border-indigo-900/40 bg-indigo-50/80 dark:bg-indigo-900/25',
      iconText: 'text-indigo-700 dark:text-indigo-300'
    },
    blue: {
      glowA: 'bg-blue-300/60 dark:bg-blue-500/25',
      glowB: 'bg-sky-200/60 dark:bg-sky-500/20',
      iconWrap: 'border-blue-200/70 dark:border-blue-900/40 bg-blue-50/80 dark:bg-blue-900/25',
      iconText: 'text-blue-700 dark:text-blue-300'
    },
    violet: {
      glowA: 'bg-violet-300/60 dark:bg-violet-500/25',
      glowB: 'bg-fuchsia-200/60 dark:bg-fuchsia-500/20',
      iconWrap: 'border-violet-200/70 dark:border-violet-900/40 bg-violet-50/80 dark:bg-violet-900/25',
      iconText: 'text-violet-700 dark:text-violet-300'
    },
    amber: {
      glowA: 'bg-amber-300/60 dark:bg-amber-500/25',
      glowB: 'bg-orange-200/60 dark:bg-orange-500/20',
      iconWrap: 'border-amber-200/70 dark:border-amber-900/40 bg-amber-50/80 dark:bg-amber-900/25',
      iconText: 'text-amber-700 dark:text-amber-300'
    },
    sky: {
      glowA: 'bg-sky-300/60 dark:bg-sky-500/25',
      glowB: 'bg-cyan-200/60 dark:bg-cyan-500/20',
      iconWrap: 'border-sky-200/70 dark:border-sky-900/40 bg-sky-50/80 dark:bg-sky-900/25',
      iconText: 'text-sky-700 dark:text-sky-300'
    }
  }
  return map[props.tone]
})
</script>

<style scoped></style>
