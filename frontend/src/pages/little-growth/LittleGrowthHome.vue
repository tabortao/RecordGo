<template>
  <div class="h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col relative overflow-hidden">
    <!-- Header -->
    <div class="bg-white/80 backdrop-blur-md fixed top-0 left-0 right-0 z-20 px-4 py-2 shadow-sm flex items-center justify-between dark:bg-gray-900/80 dark:border-b dark:border-gray-800">
      <div class="flex items-center gap-3 flex-shrink-0">
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center cursor-pointer hover:bg-gray-200 transition-colors text-gray-600 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
          @click="router.push('/homework')"
        >
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center cursor-pointer hover:bg-gray-200 transition-colors text-gray-600 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
          @click="openSidebar"
        >
          <el-icon><Menu /></el-icon>
        </div>
        <h1 class="font-bold text-lg text-gray-800 tracking-tight dark:text-white whitespace-nowrap">小成长</h1>
      </div>

      <div class="flex items-center gap-2 flex-nowrap overflow-x-auto hide-scrollbar">
        <div v-if="store.activeFilterTagId || store.onlyFavorites" class="flex items-center gap-2 bg-purple-100 px-3 py-1 rounded-full mr-2 whitespace-nowrap dark:bg-purple-900/30">
          <span class="text-xs text-purple-600 font-medium dark:text-purple-300">正在筛选: {{ activeTagName }}</span>
          <el-icon class="text-purple-400 cursor-pointer hover:text-purple-600" @click="clearFilter"><Close /></el-icon>
        </div>

        <!-- Search -->
        <el-input 
          v-model="searchQuery" 
          placeholder="搜索..." 
          class="!w-32 sm:!w-40 transition-all focus:!w-48 dark:bg-gray-800" 
          size="small"
          clearable
          :prefix-icon="Search"
        />

        <!-- Calendar -->
        <div class="relative flex items-center">
          <el-date-picker
            ref="datePickerRef"
            v-model="selectedDate"
            type="date"
            class="!w-0 !h-0 !border-0 !p-0 !overflow-hidden opacity-0 absolute top-0 left-0 -z-10"
            :popper-options="{ placement: 'bottom-end' }"
          >
            <template #cell="{ text, date, isCurrent }">
              <div class="w-full h-full flex flex-col items-center justify-center relative">
                <span :class="isCurrent ? 'font-bold text-purple-600' : ''">{{ text }}</span>
                <span v-if="hasRecord(date)" class="w-1.5 h-1.5 bg-purple-500 rounded-full absolute bottom-1"></span>
              </div>
            </template>
          </el-date-picker>
          
          <div 
            class="w-8 h-8 rounded-full flex items-center justify-center cursor-pointer transition-colors"
            :class="selectedDate ? 'bg-purple-100 text-purple-600 dark:bg-purple-900/50 dark:text-purple-300' : 'bg-gray-100 text-gray-600 hover:bg-purple-50 hover:text-purple-600 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600'"
            @click="openCalendar"
          >
             <el-icon v-if="selectedDate" @click.stop="clearSelectedDate"><Close /></el-icon>
             <el-icon v-else><Calendar /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content (Timeline) -->
    <div class="flex-1 overflow-y-auto p-4 pb-24 relative dark:bg-gray-900 bg-[#F5F7FA] px-0 sm:px-4 pt-16" ref="scrollContainer" @scroll="handleScroll">
      <div class="max-w-2xl mx-auto w-full">
        <!-- 中文注释：快捷入口按钮组（荣誉/生长/成绩） -->
        <div class="grid grid-cols-3 gap-2 px-4 mb-3 sm:px-0">
          <button
            class="flex items-center justify-center gap-1 rounded-xl border border-amber-100 bg-amber-50 text-amber-700 px-3 py-2 text-xs font-semibold shadow-sm transition hover:shadow-md dark:border-amber-500/30 dark:bg-amber-500/15 dark:text-amber-200"
            @click="router.push('/honors')"
          >
            <el-icon :size="16"><Trophy /></el-icon>
            <span>荣誉</span>
          </button>
          <button
            class="flex items-center justify-center gap-1 rounded-xl border border-emerald-100 bg-emerald-50 text-emerald-700 px-3 py-2 text-xs font-semibold shadow-sm transition hover:shadow-md dark:border-emerald-500/30 dark:bg-emerald-500/15 dark:text-emerald-200"
            @click="router.push('/growth')"
          >
            <el-icon :size="16"><TrendCharts /></el-icon>
            <span>生长</span>
          </button>
          <button
            class="flex items-center justify-center gap-1 rounded-xl border border-blue-100 bg-blue-50 text-blue-700 px-3 py-2 text-xs font-semibold shadow-sm transition hover:shadow-md dark:border-blue-500/30 dark:bg-blue-500/15 dark:text-blue-200"
            @click="router.push('/grades')"
          >
            <el-icon :size="16"><DataAnalysis /></el-icon>
            <span>成绩</span>
          </button>
        </div>
        <template v-if="hasRecords">
          <div v-if="hasPinnedRecords" class="mb-2">
            <TimelineCard 
              v-for="record in (pinnedRecords as unknown as import('@/stores/littleGrowth').GrowthRecord[])"
              :key="record.id"
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
          <div v-for="year in sortedYears" :key="year" :id="'year-' + String(year)">
            <div
              v-for="month in (sortedMonthsByYear[year] || [])"
              :key="month"
              :id="'month-' + String(year) + '-' + String(month)"
            >
              <TimelineCard
                v-for="record in (groupedRecords[year]?.[month] || [])"
                :key="record.id"
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
          class="fixed right-1 top-20 bottom-24 z-10 w-6 flex flex-col items-center justify-center"
          @touchstart="handleSliderTouchStart"
          @touchmove="handleSliderTouchMove"
          @touchend="handleSliderTouchEnd"
        >
           <div class="w-1 h-full bg-gray-200 dark:bg-gray-700 rounded-full relative">
              <!-- Handle -->
              <div 
                class="absolute w-6 h-6 bg-purple-600 rounded-full shadow-md -left-2.5 flex items-center justify-center text-[10px] text-white font-bold"
                :style="{ top: sliderTop + '%' }"
              >
                 {{ currentSliderYear }}
              </div>
           </div>
        </div>
      </transition>
    </div>

    <!-- Back to Top -->
    <div 
      v-if="showBackToTop"
        ref="backToTopRef"
        class="fixed z-50 w-10 h-10 bg-purple-600 text-white rounded-full shadow-lg flex items-center justify-center cursor-pointer hover:bg-purple-700 active:scale-95 transition-all"
        :style="backToTopStyle"
        @click="scrollToTop"
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Menu, Plus, Close, Files, Search, Calendar, Top, Trophy, TrendCharts, DataAnalysis } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useLittleGrowthStore } from '@/stores/littleGrowth'
import { useWindowSize, useDraggable, useStorage, useIntervalFn } from '@vueuse/core'
import TimelineCard from './components/TimelineCard.vue'
import TagSidebar from './components/TagSidebar.vue'
import dayjs from 'dayjs'

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

const { width, height } = useWindowSize()
const drawerSize = computed(() => width.value < 768 ? '50%' : '25%')

// Scroll & Back to Top
const scrollContainer = ref<HTMLElement | null>(null)
const isScrolling = ref(false)
let scrollTimeout: any = null
const showBackToTop = ref(false)

// Draggable Back to Top
const backToTopRef = ref<HTMLElement | null>(null)
const initialPos = useStorage('back-to-top-pos', { x: width.value - 60, y: height.value - 100 })
const { style: backToTopStyle } = useDraggable(backToTopRef as any, {
  initialValue: initialPos,
  onEnd: (position) => {
    initialPos.value = position
  }
})

onMounted(() => {
  store.fetchRecords()
})

function openSidebar() { showSidebar.value = true }

function clearSelectedDate() { selectedDate.value = null }

const handlePin = async (id: string) => {
    try {
        await store.togglePin(id)
        ElMessage.success('操作成功')
    } catch (e) {
        // error
    }
}

// Slider Logic
const sliderTop = ref(0)
const currentSliderYear = ref('')
const isDraggingSlider = ref(false)

const handleSliderTouchStart = (e: TouchEvent) => {
  isDraggingSlider.value = true
  handleSliderTouchMove(e)
}

const handleSliderTouchMove = (e: TouchEvent) => {
  e.preventDefault()
  const touch = e.touches[0]
  const sliderContainer = document.querySelector('.fixed.right-1') as HTMLElement
  if (!sliderContainer) return

  const rect = sliderContainer.getBoundingClientRect()
  const y = touch.clientY - rect.top
  const height = rect.height
  let percent = (y / height) * 100
  percent = Math.max(0, Math.min(100, percent))
  
  sliderTop.value = percent
  
  // Calculate year based on percent
  const years = sortedYears.value
  if (years.length > 0) {
    const index = Math.floor((percent / 100) * years.length)
    const safeIndex = Math.min(years.length - 1, Math.max(0, index))
    const year = years[safeIndex]
    currentSliderYear.value = year
    scrollToYear(year)
  }
}

const handleSliderTouchEnd = () => {
  isDraggingSlider.value = false
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
    list = list.filter(r => r.date === d)
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

const scrollToYear = (year: string) => {
  const el = document.getElementById('year-' + year)
  if (el) el.scrollIntoView({ behavior: 'smooth' })
}

const openCalendar = () => {
    if (selectedDate.value) {
        selectedDate.value = null
    } else {
        datePickerRef.value?.focus()
    }
}

const hasRecord = (date: Date) => {
  const d = dayjs(date).format('YYYY-MM-DD')
  return store.records.some(r => r.date === d)
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
    } catch (e) {
        ElMessage.error('操作失败')
    }
}

const handleAddComment = async (id: string, content: string) => {
    try {
        await store.addComment(id, content)
        ElMessage.success('评论成功')
    } catch (e) {
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
