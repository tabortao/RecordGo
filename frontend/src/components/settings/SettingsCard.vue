<template>
  <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/80 dark:bg-gray-900/65 backdrop-blur-xl shadow-sm">
    <div class="px-4 py-4">
      <div v-if="title || $slots.header || $slots.right" class="flex items-start justify-between gap-4">
        <div class="min-w-0">
          <slot name="header">
            <div v-if="title" class="text-sm font-bold text-gray-900 dark:text-gray-50">{{ title }}</div>
            <div v-if="description" class="mt-1 text-xs leading-relaxed text-gray-500 dark:text-gray-400">{{ description }}</div>
          </slot>
        </div>
        <div class="shrink-0">
          <slot name="right" />
        </div>
      </div>

      <div v-if="$slots.default" :class="[title || $slots.header || $slots.right ? 'mt-4' : '']">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title?: string
  description?: string
}>(), {
  title: '',
  description: ''
})
</script>

<style scoped></style>

