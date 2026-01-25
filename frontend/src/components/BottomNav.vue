<template>
  <!-- 中文注释：底部导航，当前页面高亮 -->
  <nav class="liquid-nav-shell">
    <div class="liquid-nav-safe">
      <div class="liquid-nav-bar">
        <div class="liquid-active-pill" :style="{ transform: `translateX(${activeIndex * 100}%)` }" />
        <router-link
          v-for="item in items"
          :key="item.to"
          :to="item.to"
          class="liquid-nav-item"
          :class="$route.path === item.to ? 'is-active' : ''"
        >
          <div class="liquid-icon">
            <el-icon :size="20"><component :is="item.icon" /></el-icon>
          </div>
          <span class="liquid-label">{{ item.label }}</span>
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { List, Star, User, Reading } from '@element-plus/icons-vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

// 中文注释：底部导航四页：任务/心愿/作业家/我的
const items = [
  { to: '/tasks', label: '任务', icon: List },
  { to: '/wishes', label: '心愿', icon: Star },
  { to: '/homework', label: '作业家', icon: Reading },
  { to: '/mine', label: '我的', icon: User }
]

const route = useRoute()
const activeIndex = computed(() => {
  const idx = items.findIndex((x) => route.path === x.to)
  return idx >= 0 ? idx : 0
})
</script>

<style scoped>
:global(:root) {
  --liquid-nav-radius: 26px;
}

.liquid-nav-shell {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 50;
  pointer-events: none;
}

.liquid-nav-safe {
  width: calc(100vw - 32px);
  margin: 0 auto;
  padding: 12px 0 calc(env(safe-area-inset-bottom) + 12px);
}

.liquid-nav-bar {
  pointer-events: auto;
  position: relative;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  align-items: center;
  height: 68px;
  padding: 8px;
  border-radius: var(--liquid-nav-radius);
  background: rgb(255 255 255 / 0.68);
  border: 1px solid rgb(255 255 255 / 0.5);
  box-shadow:
    0 16px 50px rgb(0 0 0 / 0.14),
    0 1px 0 rgb(255 255 255 / 0.6) inset;
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  overflow: hidden;
}

.dark .liquid-nav-bar {
  background: rgb(17 24 39 / 0.7);
  border: 1px solid rgb(148 163 184 / 0.14);
  box-shadow:
    0 18px 60px rgb(0 0 0 / 0.45),
    0 1px 0 rgb(255 255 255 / 0.05) inset;
}

.liquid-nav-bar::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, rgb(255 255 255 / 0.38), transparent 44%);
  opacity: 0.6;
}

.dark .liquid-nav-bar::before {
  background: linear-gradient(180deg, rgb(255 255 255 / 0.08), transparent 40%);
  opacity: 1;
}

.liquid-active-pill {
  position: absolute;
  left: 8px;
  top: 8px;
  bottom: 8px;
  width: calc((100% - 16px) / 4);
  border-radius: calc(var(--liquid-nav-radius) - 10px);
  background: rgb(16 185 129 / 0.12);
  border: 1px solid rgb(16 185 129 / 0.16);
  box-shadow:
    0 16px 34px rgb(16 185 129 / 0.14),
    0 1px 0 rgb(255 255 255 / 0.6) inset;
  transition: transform 360ms cubic-bezier(.2,.9,.2,1), background-color 240ms ease;
  will-change: transform;
}

.dark .liquid-active-pill {
  background: rgb(16 185 129 / 0.14);
  border: 1px solid rgb(16 185 129 / 0.16);
  box-shadow:
    0 18px 40px rgb(0 0 0 / 0.5),
    0 1px 0 rgb(255 255 255 / 0.08) inset;
}

.liquid-nav-item {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  height: 52px;
  border-radius: 18px;
  color: rgb(55 65 81 / 0.9);
  text-decoration: none;
  transition: transform 220ms ease, color 220ms ease, opacity 220ms ease;
  user-select: none;
}

.dark .liquid-nav-item {
  color: rgb(229 231 235 / 0.88);
}

.liquid-nav-item:hover {
  transform: translateY(-1px);
}

.liquid-nav-item:active {
  transform: translateY(0);
}

.liquid-icon {
  width: 34px;
  height: 28px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  transition: transform 220ms ease;
}

.liquid-label {
  font-size: 11px;
  line-height: 1;
  opacity: 0.78;
}

.liquid-nav-item.is-active {
  color: rgb(4 120 87);
  font-weight: 700;
}

.dark .liquid-nav-item.is-active {
  color: rgb(167 243 208);
}

.liquid-nav-item.is-active .liquid-icon {
  transform: translateY(-0.5px) scale(1.02);
}
</style>
