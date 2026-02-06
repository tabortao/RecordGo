<template>
  <div class="h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col relative overflow-hidden">
    <!-- Header -->
    <div class="bg-white/80 backdrop-blur-xl fixed top-0 left-0 right-0 z-20 px-4 py-3 shadow-sm border-b border-gray-100 dark:bg-gray-900/80 dark:border-gray-800 transition-all duration-300">
      <div class="flex items-center justify-between max-w-6xl mx-auto w-full">
        <div class="flex items-center gap-4 flex-shrink-0">
          <button 
            class="w-9 h-9 rounded-full bg-gray-50 flex items-center justify-center cursor-pointer hover:bg-gray-100 transition-colors text-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 active:scale-95"
            @click="router.push('/homework')"
          >
            <el-icon :size="18"><ArrowLeft /></el-icon>
          </button>
          <button 
            class="w-9 h-9 rounded-full bg-gray-50 flex items-center justify-center cursor-pointer hover:bg-gray-100 transition-colors text-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 active:scale-95"
            @click="openSidebar"
          >
            <el-icon :size="18"><Menu /></el-icon>
          </button>
          <h1 class="font-black text-xl text-gray-800 tracking-tight dark:text-white whitespace-nowrap bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-blue-500">小成长</h1>
        </div>

        <div class="flex items-center gap-3 flex-nowrap overflow-x-auto hide-scrollbar pl-2">
          <div v-if="store.activeFilterTagId || store.onlyFavorites" class="flex items-center gap-2 bg-purple-50 px-3 py-1.5 rounded-full mr-1 whitespace-nowrap dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50">
            <span class="text-xs text-purple-600 font-bold dark:text-purple-300">筛选: {{ activeTagName }}</span>
            <el-icon class="text-purple-400 cursor-pointer hover:text-purple-600 dark:hover:text-purple-200 transition-colors" @click="clearFilter"><Close /></el-icon>
          </div>

          <!-- Search -->
          <div class="relative group">
            <el-input 
              v-model="searchQuery" 
              placeholder="搜索美好回忆..." 
              class="!w-32 sm:!w-48 transition-all focus:!w-56 dark:bg-gray-800 search-input" 
              size="default"
              clearable
              :prefix-icon="Search"
            />
          </div>

          <!-- Calendar -->
          <div class="relative flex items-center">
            <el-date-picker
              ref="datePickerRef"
              v-model="selectedDate"
              type="date"
              class="!w-0 !h-0 !border-0 !p-0 !overflow-hidden opacity-0 absolute top-0 left-0 -z-10"
              :popper-options="{ placement: 'bottom-end' }"
              @change="handleCalendarChange"
            >
              <template #cell="{ text, date, isCurrent }">
                <el-tooltip
                  :content="getRecordCountText(date)"
                  :disabled="getRecordCount(date) <= 0"
                  placement="top"
                  effect="dark"
                >
                  <div class="w-full h-full flex items-center justify-center">
                    <div
                      class="relative h-7 w-7 rounded-full flex items-center justify-center text-sm transition-all duration-300"
                      :class="hasRecord(date) ? 'bg-purple-600 text-white font-black shadow-md scale-110' : (isCurrent ? 'font-extrabold text-purple-600 bg-purple-50' : 'hover:bg-gray-100 dark:hover:bg-gray-700')"
                    >
                      {{ text }}
                      <div
                        v-if="getRecordCount(date) > 0"
                        class="absolute -top-1 -right-1 h-3.5 min-w-[14px] px-0.5 rounded-full bg-amber-400 text-white text-[9px] font-black shadow-sm border border-white dark:border-gray-800 flex items-center justify-center leading-none"
                      >
                        {{ getRecordCountBadge(date) }}
                      </div>
                    </div>
                  </div>
                </el-tooltip>
              </template>
            </el-date-picker>
            
            <button 
              class="w-9 h-9 rounded-full flex items-center justify-center cursor-pointer transition-all active:scale-95"
              :class="selectedDate ? 'bg-purple-100 text-purple-600 dark:bg-purple-900/50 dark:text-purple-300 ring-2 ring-purple-200 dark:ring-purple-800' : 'bg-gray-50 text-gray-600 hover:bg-purple-50 hover:text-purple-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'"
              @click="openCalendar"
            >
               <el-icon v-if="selectedDate" @click.stop="clearSelectedDate"><Close /></el-icon>
               <el-icon v-else :size="18"><Calendar /></el-icon>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content (Timeline) -->
    <div class="flex-1 overflow-y-auto p-4 pb-24 relative dark:bg-gray-900 bg-[#F5F7FA] px-0 sm:px-4 pt-16" ref="scrollContainer" @scroll="handleScroll">
      <div class="max-w-2xl mx-auto w-full">
        <!-- 快捷入口卡片 -->
        <div class="grid grid-cols-3 gap-3 px-4 mb-8 sm:px-0">
          <button
            class="group relative overflow-hidden rounded-2xl bg-white p-3 shadow-sm transition-all hover:-translate-y-1 hover:shadow-md dark:bg-gray-800 border border-gray-100 dark:border-gray-700"
            @click="router.push('/honors')"
          >
            <div class="absolute -right-6 -top-6 h-20 w-20 rounded-full bg-amber-50 dark:bg-amber-900/20 transition-transform group-hover:scale-110"></div>
            <div class="relative z-10 flex flex-col items-center gap-2">
              <div class="rounded-full bg-amber-100 p-2.5 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400 shadow-sm relative">
                <el-icon :size="20" class="text-amber-600 dark:text-amber-400"><TrophyBase /></el-icon>
              </div>
              <span class="text-xs font-bold text-gray-700 dark:text-gray-200 tracking-wide">荣誉墙</span>
            </div>
          </button>
          
          <button
            class="group relative overflow-hidden rounded-2xl bg-white p-3 shadow-sm transition-all hover:-translate-y-1 hover:shadow-md dark:bg-gray-800 border border-gray-100 dark:border-gray-700"
            @click="router.push('/growth')"
          >
            <div class="absolute -right-6 -top-6 h-20 w-20 rounded-full bg-emerald-50 dark:bg-emerald-900/20 transition-transform group-hover:scale-110"></div>
            <div class="relative z-10 flex flex-col items-center gap-2">
              <div class="rounded-full bg-emerald-100 p-2.5 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400 shadow-sm relative">
                <el-icon :size="20" class="text-emerald-600 dark:text-emerald-400 fill-current"><TrendCharts /></el-icon>
                <div class="absolute top-0 right-0 w-2.5 h-2.5 bg-emerald-500 rounded-full border-2 border-white dark:border-gray-800"></div>
              </div>
              <span class="text-xs font-bold text-gray-700 dark:text-gray-200 tracking-wide">生长记录</span>
            </div>
          </button>
          
          <button
            class="group relative overflow-hidden rounded-2xl bg-white p-3 shadow-sm transition-all hover:-translate-y-1 hover:shadow-md dark:bg-gray-800 border border-gray-100 dark:border-gray-700"
            @click="router.push('/grades')"
          >
            <div class="absolute -right-6 -top-6 h-20 w-20 rounded-full bg-blue-50 dark:bg-blue-900/20 transition-transform group-hover:scale-110"></div>
            <div class="relative z-10 flex flex-col items-center gap-2">
              <div class="rounded-full bg-blue-100 p-2.5 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400 shadow-sm">
                <el-icon :size="20"><DataAnalysis /></el-icon>
              </div>
              <span class="text-xs font-bold text-gray-700 dark:text-gray-200 tracking-wide">成绩分析</span>
            </div>
          </button>
        </div>

        <template v-if="hasRecords">
          <!-- 置顶记录 -->
          <div v-if="hasPinnedRecords" class="mb-8 px-4 sm:px-0">
             <div class="flex items-center gap-2 mb-3 text-purple-600 dark:text-purple-400">
                <el-icon><Top /></el-icon>
                <span class="text-sm font-bold">置顶回忆</span>
             </div>
             <div class="space-y-4">
                <div
                  v-for="record in (pinnedRecords as unknown as import('@/stores/littleGrowth').GrowthRecord[])"
                  :key="record.id"
                  :id="'record-' + String(record.id)"
                >
                  <TimelineCard 
                    :record="record"
                    :allTags="store.flattenedTags"
                    :searchQuery="getSearchQuery()"
                    @edit="handleEdit"
                    @delete="handleDelete"
                    @pin="handlePin"
                    @filter-tag="handleTagSelect"
                    @toggle-favorite="handleToggleFavorite"
                    @add-comment="handleAddComment"
                  />
                </div>
             </div>
          </div>

          <!-- 时间轴列表 -->
          <div class="relative">
            <!-- 桌面端左侧连接线 -->
            <div class="absolute left-[26px] top-12 bottom-0 w-0.5 bg-gradient-to-b from-purple-200 via-purple-100 to-transparent dark:from-purple-900/40 dark:via-purple-900/20 hidden sm:block"></div>

            <div v-for="year in sortedYears" :key="year" :id="'year-' + String(year)" class="relative mb-12">
                <!-- 年份标题 (Sticky) -->
                <div class="sticky top-[80px] z-10 mb-8 pl-4 sm:pl-16 transition-all duration-300">
                    <span class="inline-flex items-center gap-1.5 rounded-full bg-white/95 px-5 py-2 text-2xl font-black text-gray-800 shadow-lg shadow-purple-100/50 backdrop-blur-xl dark:bg-gray-800/95 dark:text-white dark:shadow-none border border-purple-50 dark:border-gray-700">
                        {{ year }}<span class="text-sm font-bold text-gray-400 mt-1">年</span>
                    </span>
                </div>

                <div
                  v-for="month in (sortedMonthsByYear[year] || [])"
                  :key="month"
                  :id="'month-' + String(year) + '-' + String(month)"
                  class="mb-10 relative"
                >
                  <!-- 月份标记 -->
                  <div class="mb-6 flex items-center gap-4 pl-4 sm:pl-0">
                      <!-- Desktop: 圆形月份标 -->
                      <div class="hidden sm:flex h-14 w-14 items-center justify-center rounded-full border-[4px] border-[#F5F7FA] bg-white text-purple-600 shadow-md z-10 dark:border-gray-900 dark:bg-gray-800 dark:text-purple-400 dark:shadow-none">
                          <span class="text-lg font-black">{{ month }}<span class="text-xs font-bold ml-0.5">月</span></span>
                      </div>
                      
                      <!-- Mobile: 胶囊月份标 -->
                      <div class="sm:hidden flex items-center gap-3">
                          <div class="h-8 w-1.5 rounded-full bg-gradient-to-b from-purple-500 to-blue-500 shadow-sm"></div>
                          <span class="text-xl font-black text-gray-800 dark:text-gray-100">{{ month }}月</span>
                      </div>
                  </div>

                  <!-- 卡片列表容器 -->
                  <div class="px-4 sm:px-0 sm:pl-20 space-y-6">
                      <div
                        v-for="record in (groupedRecords[year]?.[month] || [])"
                        :key="record.id"
                        :id="'record-' + String(record.id)"
                      >
                        <TimelineCard
                          :record="record"
                          :allTags="store.flattenedTags"
                          :searchQuery="getSearchQuery()"
                          @edit="handleEdit"
                          @delete="handleDelete"
                          @pin="handlePin"
                          @filter-tag="handleTagSelect"
                          @toggle-favorite="handleToggleFavorite"
                          @add-comment="handleAddComment"
                        />
                      </div>
                  </div>
                </div>
            </div>
          </div>
        </template>
        <div v-else class="flex flex-col items-center justify-center py-20 text-gray-400">
          <el-icon :size="48" class="mb-4 opacity-50"><Files /></el-icon>
          <p>暂无记录，快去添加第一条美好回忆吧~</p>
        </div>
      </div>
      
      <!-- Timeline Slider (Draggable) -->
      <transition name="fade">
        <div 
          v-show="isScrolling" 
          ref="sliderRef"
          class="fixed right-1 top-20 bottom-24 z-10 w-8 flex flex-col items-center justify-center select-none"
          @touchstart="handleSliderTouchStart"
          @touchmove="handleSliderTouchMove"
          @touchend="handleSliderTouchEnd"
          @click="handleSliderClick"
        >
           <div class="w-1 h-full bg-gray-200 dark:bg-gray-700 rounded-full relative">
              <!-- Handle -->
              <div 
                class="absolute w-7 h-7 bg-purple-600 rounded-full shadow-md -left-3 flex items-center justify-center text-white"
                :style="{ top: sliderTop + '%' }"
              >
                 <div class="flex flex-col items-center leading-none">
                   <span class="text-[12px] font-black">{{ currentSliderMonthLabel }}</span>
                   <span class="text-[8px] font-semibold opacity-90">{{ currentSliderYear }}</span>
                 </div>
                 <div
                   v-if="currentSliderYear && currentSliderMonth"
                   class="absolute right-8 whitespace-nowrap rounded-full bg-white/90 px-3 py-1 text-[12px] font-bold text-gray-800 shadow-sm backdrop-blur-md dark:bg-gray-800/90 dark:text-gray-100 border border-gray-100 dark:border-gray-700"
                 >
                   {{ currentSliderYear }}年{{ currentSliderMonthLabel }}月
                 </div>
              </div>
           </div>
        </div>
      </transition>
    </div>

    <!-- Back to Top -->
    <div 
      v-if="showBackToTop"
        ref="backToTopRef"
        class="fixed z-50 w-10 h-10 bg-purple-600 text-white rounded-full shadow-lg flex items-center justify-center cursor-pointer hover:bg-purple-700 active:scale-95 transition-all touch-none"
        :style="backToTopStyle"
        @click.stop="scrollToTop"
      >
        <el-icon><Top /></el-icon>
      </div>

    <!-- FAB (Add Button) -->
    <div 
      class="fixed right-6 bottom-8 z-20 w-14 h-14 bg-gradient-to-tr from-purple-400 to-blue-400 rounded-full shadow-lg shadow-purple-200 flex items-center justify-center text-white cursor-pointer hover:scale-110 transition-transform active:scale-95"
      @click="router.push('/little-growth/create')"
    >
      <el-icon :size="24"><Plus /></el-icon>
    </div>

    <!-- Sidebar Drawer -->
    <el-drawer
      v-model="showSidebar"
      direction="ltr"
      :with-header="false"
      :size="drawerSize"
      class="little-growth-drawer"
    >
      <TagSidebar 
        :tags="store.flattenedTags" 
        :active-tag-id="store.activeFilterTagId"
        :total-records="store.records.length"
        :show-favorites="store.onlyFavorites"
        :favorites-count="Number(favoritesCount)"
        @select="handleTagSelect"
        @select-favorites="handleFavoritesSelect"
      />
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Menu, Plus, Close, Files, Search, Calendar, Top, TrophyBase, TrendCharts, DataAnalysis } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useLittleGrowthStore } from '@/stores/littleGrowth'
import { useWindowSize, useDraggable, useStorage, useIntervalFn } from '@vueuse/core'
import TimelineCard from './components/TimelineCard.vue'
import TagSidebar from './components/TagSidebar.vue'
import dayjs from 'dayjs'
import http from '@/services/http'

const router = useRouter()
const store = useLittleGrowthStore()

const { pause } = useIntervalFn(() => {
  store.fetchRecords({ is_favorite: store.onlyFavorites, silent: true })
}, 30000) // Poll every 30s

onUnmounted(() => {
  pause()
})

const showSidebar = ref(false)

const searchQuery = ref('')
const selectedDate = ref<Date | null>(null)
const datePickerRef = ref()
const recordDateCountMap = ref<Record<string, number>>({})

const { width, height } = useWindowSize()
const drawerSize = computed(() => width.value < 768 ? '50%' : '25%')

// Scroll & Back to Top
const scrollContainer = ref<HTMLElement | null>(null)
const isScrolling = ref(false)
let scrollTimeout: any = null
const showBackToTop = ref(false)

// Draggable Back to Top
const backToTopRef = ref<HTMLElement | null>(null)
// Initialize position. Note: we don't know window size on SSR or early init, so we default to bottom right.
// We'll update it in onMounted if needed or let useDraggable handle it.
// To ensure it shows on desktop, we avoid complex logic that might be hiding it based on screen width unless intended.
const initialPos = useStorage('back-to-top-pos', { x: 0, y: 0 }) 

const { style: backToTopStyle } = useDraggable(backToTopRef as any, {
  initialValue: initialPos,
  onEnd: (position) => {
    initialPos.value = position
  },
  preventDefault: true
})

// Initialize position if it's 0,0 (first run) to bottom right
onMounted(() => {
  store.fetchRecords()
  if (initialPos.value.x === 0 && initialPos.value.y === 0) {
      initialPos.value = { x: width.value - 60, y: height.value - 100 }
  }
})

function openSidebar() { showSidebar.value = true }

function clearSelectedDate() { selectedDate.value = null }

const handlePin = async (id: string) => {
    try {
        await store.togglePin(id)
        ElMessage.success('操作成功')
    } catch {
        // error
    }
}

// Slider Logic
const sliderRef = ref<HTMLElement | null>(null)
const sliderTop = ref(0)
const currentSliderYear = ref('')
const currentSliderMonth = ref('')
const currentSliderMonthLabel = computed(() => currentSliderMonth.value ? String(Number(currentSliderMonth.value)) : '')
const isDraggingSlider = ref(false)

type MonthAnchor = { year: string; month: string; key: string }
type MonthAnchorOffset = MonthAnchor & { top: number }

const monthAnchors = computed<MonthAnchor[]>(() => {
  const anchors: MonthAnchor[] = []
  for (const year of sortedYears.value) {
    const months = sortedMonthsByYear.value[year] || []
    for (const month of months) {
      anchors.push({ year, month, key: `${year}-${month}` })
    }
  }
  return anchors
})

const monthAnchorOffsets = ref<MonthAnchorOffset[]>([])

async function recomputeMonthAnchorOffsets() {
  await nextTick()
  const container = scrollContainer.value
  if (!container) {
    monthAnchorOffsets.value = []
    return
  }

  const containerTop = container.getBoundingClientRect().top
  const baseScrollTop = container.scrollTop
  const offsets: MonthAnchorOffset[] = []
  for (const a of monthAnchors.value) {
    const el = document.getElementById(`month-${a.year}-${a.month}`)
    if (!el) continue
    const rect = el.getBoundingClientRect()
    offsets.push({ ...a, top: rect.top - containerTop + baseScrollTop })
  }
  offsets.sort((a, b) => a.top - b.top)
  monthAnchorOffsets.value = offsets
  updateSliderFromScroll()
}

function updateSliderFromScroll() {
  const container = scrollContainer.value
  const offsets = monthAnchorOffsets.value
  if (!container || offsets.length === 0) return

  const y = container.scrollTop + 140
  let idx = 0
  for (let i = 0; i < offsets.length; i++) {
    if (offsets[i].top <= y) idx = i
    else break
  }
  const cur = offsets[idx]
  currentSliderYear.value = cur.year
  currentSliderMonth.value = cur.month
  sliderTop.value = offsets.length === 1 ? 0 : (idx / (offsets.length - 1)) * 100
}

const scrollToMonth = (year: string, month: string) => {
  const el = document.getElementById(`month-${year}-${month}`)
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function setSliderPercent(percent: number, shouldScroll: boolean) {
  const p = Math.max(0, Math.min(100, percent))
  sliderTop.value = p

  const anchors = monthAnchorOffsets.value.length > 0 ? monthAnchorOffsets.value : monthAnchors.value
  if (anchors.length === 0) return

  const index = Math.floor((p / 100) * anchors.length)
  const safeIndex = Math.min(anchors.length - 1, Math.max(0, index))
  const a: any = anchors[safeIndex]
  currentSliderYear.value = a.year
  currentSliderMonth.value = a.month
  if (shouldScroll) scrollToMonth(a.year, a.month)
}

const handleSliderTouchStart = (e: TouchEvent) => {
  isDraggingSlider.value = true
  handleSliderTouchMove(e)
}

const handleSliderTouchMove = (e: TouchEvent) => {
  e.preventDefault()
  const touch = e.touches[0]
  const sliderContainer = sliderRef.value
  if (!sliderContainer) return

  const rect = sliderContainer.getBoundingClientRect()
  const y = touch.clientY - rect.top
  const height = rect.height
  let percent = (y / height) * 100
  setSliderPercent(percent, true)
}

const handleSliderTouchEnd = () => {
  isDraggingSlider.value = false
}

const handleSliderClick = (e: MouseEvent) => {
  const sliderContainer = sliderRef.value
  if (!sliderContainer) return
  const rect = sliderContainer.getBoundingClientRect()
  const y = e.clientY - rect.top
  const height = rect.height
  const percent = (y / height) * 100
  setSliderPercent(percent, true)
}

const handleScroll = () => {
  isScrolling.value = true
  clearTimeout(scrollTimeout)
  scrollTimeout = setTimeout(() => {
    if (!isDraggingSlider.value) isScrolling.value = false
  }, 1500)

  if (scrollContainer.value) {
    showBackToTop.value = scrollContainer.value.scrollTop > 300
  }
  if (!isDraggingSlider.value) updateSliderFromScroll()
}

const scrollToTop = () => {
  scrollContainer.value?.scrollTo({ top: 0, behavior: 'smooth' })
}


const activeTagName = computed(() => {
  if (store.onlyFavorites) return '我的收藏'
  const tag = store.flattenedTags.find(t => t.id === store.activeFilterTagId)
  return tag ? tag.name : ''
})

// 中文注释：收藏数量计算，实时统计当前已加载记录中的收藏条数
const favoritesCount = computed(() => {
  return store.records.filter(r => !!r.is_favorite).length
})

// 中文注释：父子标签映射，用于父标签筛选时包含其子标签
const childrenMap = computed<Record<string, { id: string }[]>>(() => {
  const m: Record<string, { id: string }[]> = {}
  for (const t of store.flattenedTags) {
    if ((t as any).parentId) {
      const pid = String((t as any).parentId)
      ;(m[pid] = m[pid] || []).push({ id: String(t.id) })
    }
  }
  return m
})

const clearFilter = () => {
    store.activeFilterTagId = null
    if (store.onlyFavorites) {
        store.onlyFavorites = false
        store.fetchRecords()
    }
}

const normalizeRecordDate = (v: any) => String(v || '').slice(0, 10)

const filteredList = computed(() => {
  let list = store.records

  // 1. 标签筛选：父标签包含其所有子标签
  if (store.activeFilterTagId) {
    const activeId = String(store.activeFilterTagId)
    const activeTag = store.flattenedTags.find(t => String(t.id) === activeId)

    if (activeTag && !(activeTag as any).parentId) {
      const childIds = (childrenMap.value[activeId] || []).map(x => x.id)
      const allIds = [activeId, ...childIds]
      list = list.filter(r => {
        let rTags: string[] = []
        if (Array.isArray(r.tags)) {
          rTags = r.tags
        } else if (typeof r.tags === 'string') {
          try { rTags = JSON.parse(r.tags) } catch {}
        }
        const ids = rTags.map(String)
        return ids.some(id => allIds.includes(id))
      })
    } else {
      list = list.filter(r => {
        let rTags: string[] = []
        if (Array.isArray(r.tags)) {
          rTags = r.tags
        } else if (typeof r.tags === 'string') {
          try { rTags = JSON.parse(r.tags) } catch {}
        }
        return rTags.map(String).includes(activeId)
      })
    }
  }

  // 2. Search Filter
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(r => r.content.toLowerCase().includes(q))
  }

  // 3. Date Filter
  if (selectedDate.value) {
    const d = dayjs(selectedDate.value).format('YYYY-MM-DD')
    list = list.filter(r => normalizeRecordDate(r.date) === d)
  }
  
  return list
})

const hasRecords = computed(() => filteredList.value.length > 0)

function getSearchQuery() { return searchQuery.value }

type YearMonthGroups = Record<string, Record<string, import('@/stores/littleGrowth').GrowthRecord[]>>
const pinnedRecords = computed<import('@/stores/littleGrowth').GrowthRecord[]>(() => {
  return [...filteredList.value]
    .filter(r => !!r.is_pinned)
    .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
})

const hasPinnedRecords = computed(() => pinnedRecords.value.length > 0)

const groupedRecords = computed<YearMonthGroups>(() => {
  const groups: YearMonthGroups = {}
  
  // Sort list first by date desc (using full date string which now includes time)
  const list = [...filteredList.value]
    .filter(r => !r.is_pinned)
    .sort((a, b) => {
    // Use timestamp for accurate comparison including time
    const tA = new Date(a.date).getTime()
    const tB = new Date(b.date).getTime()
    // If times are equal, use ID or created_at as fallback
    if (tA === tB) {
        return String(b.id).localeCompare(String(a.id))
    }
    return tB - tA
  })

  list.forEach(r => {
    const d = dayjs(r.date)
    const y = d.format('YYYY')
    const m = d.format('MM')
    if (!groups[y]) groups[y] = {}
    if (!groups[y][m]) groups[y][m] = []
    groups[y][m].push(r)
  })
  return groups
})

const sortedYears = computed(() => {
  return Object.keys(groupedRecords.value).sort((a, b) => Number(b) - Number(a))
})

const sortedMonthsByYear = computed<Record<string, string[]>>(() => {
  const res: Record<string, string[]> = {}
  for (const [year, months] of Object.entries(groupedRecords.value)) {
    res[year] = Object.keys(months).sort((a, b) => Number(b) - Number(a))
  }
  return res
})

watch(groupedRecords, () => {
  recomputeMonthAnchorOffsets()
}, { flush: 'post' })

watch([width, height], () => {
  recomputeMonthAnchorOffsets()
}, { flush: 'post' })

async function ensureRecordDatesLoaded() {
  if (Object.keys(recordDateCountMap.value).length > 0) return
  try {
    const res = await http.get('/little-growth/records', { params: { page_size: 5000 } }) as any
    const items = Array.isArray(res?.items) ? res.items : []
    const map: Record<string, number> = {}
    for (const r of items) {
      const k = String((r as any)?.date || '').slice(0, 10)
      if (!k) continue
      map[k] = (map[k] || 0) + 1
    }
    recordDateCountMap.value = map
  } catch {}
}

const getRecordCount = (date: Date) => {
  const d = dayjs(date).format('YYYY-MM-DD')
  const v = recordDateCountMap.value[d]
  if (typeof v === 'number') return v
  return store.records.reduce((acc, r) => acc + (normalizeRecordDate(r.date) === d ? 1 : 0), 0)
}

const getRecordCountBadge = (date: Date) => {
  const n = getRecordCount(date)
  if (n <= 0) return ''
  return n > 9 ? '9+' : String(n)
}

const getRecordCountText = (date: Date) => {
  const n = getRecordCount(date)
  return n > 0 ? `当天 ${n} 篇` : ''
}

const handleCalendarChange = async (val: Date | null) => {
  if (!val) return
  if (getRecordCount(val) <= 0) {
    selectedDate.value = null
    return
  }
  await nextTick()
  const target = pinnedRecords.value[0] || filteredList.value[0]
  if (!target) return
  await nextTick()
  const el = document.getElementById(`record-${String((target as any).id)}`)
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const openCalendar = async () => {
  if (selectedDate.value) {
    selectedDate.value = null
    return
  }
  await ensureRecordDatesLoaded()
  datePickerRef.value?.focus()
}

const hasRecord = (date: Date) => {
  return getRecordCount(date) > 0
}

const handleTagSelect = (tagId: string | null) => {
  store.activeFilterTagId = tagId
  // If we were in favorite mode, exit it and fetch all (or just let filteredList handle it if we fetched all)
  // But we fetch page 1 by default.
  // Ideally we should re-fetch "All" if we switch context.
  if (store.onlyFavorites) {
      store.onlyFavorites = false
      store.fetchRecords()
  }
  showSidebar.value = false
}

const handleFavoritesSelect = () => {
    store.activeFilterTagId = null
    store.onlyFavorites = true
    store.fetchRecords({ is_favorite: true })
    showSidebar.value = false
}

const handleToggleFavorite = async (id: string) => {
    try {
        await store.toggleFavorite(id)
    } catch {
        ElMessage.error('操作失败')
    }
}

const handleAddComment = async (id: string, content: string) => {
    try {
        await store.addComment(id, content)
        ElMessage.success('评论成功')
    } catch {
        ElMessage.error('评论失败')
    }
}

const handleEdit = (id: string) => {
  router.push(`/little-growth/edit/${id}`)
}

  const handleDelete = async (id: string) => {
    try {
      await ElMessageBox.confirm(
        '确定要删除这条记录吗？删除后无法恢复。',
        '删除确认',
        {
          confirmButtonText: '删除',
          cancelButtonText: '取消',
          type: 'warning',
          center: true,
          draggable: true,
          lockScroll: false,
          customClass: 'center-message-box custom-delete-dialog',
          showClose: false,
          closeOnClickModal: false,
          closeOnPressEscape: true
        }
      )
      await store.deleteRecord(id)
      await store.fetchRecords() // Refresh list from server
      ElMessage.success('删除成功')
  } catch {
    // cancelled
  }
}
</script>

<style>
.little-growth-drawer .el-drawer__body {
  padding: 0;
}

/* CSS for hiding scrollbar in header */
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

/* Search Input Custom Style */
.search-input .el-input__wrapper {
  border-radius: 999px;
  background-color: #f9fafb;
  box-shadow: none !important;
  border: 1px solid transparent;
  padding-left: 12px;
  transition: all 0.3s;
}
.dark .search-input .el-input__wrapper {
  background-color: #1f2937;
}
.search-input .el-input__wrapper:hover,
.search-input .el-input__wrapper.is-focus {
  background-color: #fff;
  border-color: #e5e7eb;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06) !important;
}
.dark .search-input .el-input__wrapper:hover,
.dark .search-input .el-input__wrapper.is-focus {
  background-color: #374151;
  border-color: #4b5563;
  box-shadow: none !important;
}

/* Fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 对话框按钮居中 */
.center-message-box .el-message-box__btns {
  display: flex;
  justify-content: center;
}
.center-message-box .el-message-box__btns .el-button {
  min-width: 88px;
}
</style>
