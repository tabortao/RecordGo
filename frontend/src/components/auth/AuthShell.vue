<template>
  <div class="min-h-screen relative overflow-hidden bg-gradient-to-br from-slate-950 via-slate-950 to-indigo-950">
    <div class="pointer-events-none absolute inset-0">
      <div class="absolute -top-32 -left-28 h-[420px] w-[420px] rounded-full bg-indigo-600/30 blur-3xl" />
      <div class="absolute -bottom-48 -right-44 h-[560px] w-[560px] rounded-full bg-cyan-500/18 blur-3xl" />
      <div class="absolute left-1/3 top-10 h-56 w-56 rounded-full bg-fuchsia-500/10 blur-3xl" />
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,rgba(255,255,255,0.09),transparent_42%),radial-gradient(circle_at_80%_80%,rgba(255,255,255,0.06),transparent_45%)]" />
      <div class="absolute inset-0 opacity-[0.12] [background-image:linear-gradient(to_right,rgba(255,255,255,0.10)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.10)_1px,transparent_1px)] [background-size:64px_64px]" />
    </div>

    <div class="relative mx-auto flex min-h-screen w-full max-w-5xl items-center px-4 py-10">
      <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-2">
        <div class="hidden md:flex flex-col justify-between rounded-3xl border border-white/10 bg-white/5 p-8 backdrop-blur-xl shadow-2xl shadow-black/30">
          <div>
            <div class="inline-flex items-center gap-2">
              <div class="h-10 w-10 rounded-2xl bg-white/10 border border-white/10 flex items-center justify-center text-white font-extrabold tracking-tight">RG</div>
              <div class="text-3xl font-extrabold tracking-tight text-white">{{ brandTitle }}</div>
            </div>
            <div class="mt-2 text-sm text-white/70">{{ brandSubtitle }}</div>
            <div class="mt-8 space-y-3 text-sm text-white/80">
              <div v-for="(t, idx) in tips" :key="idx" class="rounded-2xl border border-white/10 bg-white/5 p-3">
                {{ t }}
              </div>
            </div>
          </div>
          <div class="pt-8 text-xs text-white/40">登录即表示同意遵守应用的使用规范与隐私约定</div>
        </div>

        <div class="rounded-3xl border border-white/10 bg-white/95 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl md:p-8">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  brandTitle?: string
  brandSubtitle?: string
  tips?: string[]
}>(), {
  brandTitle: 'RecordGo',
  brandSubtitle: '用更轻的操作，养成更稳的习惯',
  tips: () => [
    '任务、心愿、积分一体化管理',
    '子账号令牌登录，家庭协作更方便',
    '登录后支持在“我的-设置”里修改密码'
  ]
})
</script>

<style scoped></style>

