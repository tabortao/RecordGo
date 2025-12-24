import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useAppState } from '@/stores/appState'
import http from '@/services/http'
import { ElMessage } from 'element-plus'

// 中文注释：简单的三页面路由结构，默认进入任务页；恢复无登录守卫的版本
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/tasks' },
    // 中文注释：新增登录与注册占位页路由，便于用户返回登录/注册页面
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/pages/LoginPage.vue'),
      meta: { public: true, noNav: true }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/pages/RegisterPage.vue'),
      meta: { public: true, noNav: true }
    },
    {
      path: '/tasks',
      name: 'Tasks',
      component: () => import('@/pages/TasksPage.vue')
    },
    // 中文注释：为规避别名解析诊断问题，这里改用相对路径导入统计页面
    { path: '/tasks/stats', name: 'TasksStats', component: () => import('../pages/TasksStatsPage.vue'), meta: { noNav: true } },
    // 中文注释：任务独立页面路由（创建、编辑、番茄钟），隐藏底部导航
    { path: '/tasks/create', name: 'TaskCreate', component: () => import('@/pages/TaskCreatePage.vue'), meta: { noNav: true } },
    { path: '/data/records', name: 'DataRecords', component: () => import('@/pages/DataRecordsPage.vue'), meta: { noNav: true } },
    { path: '/tasks/:id/edit', name: 'TaskEdit', component: () => import('@/pages/TaskEditPage.vue'), meta: { noNav: true } },
    { path: '/tasks/:id/tomato', name: 'TaskTomato', component: () => import('@/pages/TomatoPage.vue'), meta: { noNav: true } },
    // 中文注释：任务备注独立页面，支持多条备注与附件（图片/音频），隐藏底部导航
    // 中文注释：为规避别名解析的偶发诊断错误，这里改用相对路径导入
    { path: '/tasks/:id/notes', name: 'TaskNotes', component: () => import('../pages/TaskNotesPage.vue'), meta: { noNav: true } },
    {
      path: '/wishes',
      name: 'Wishes',
      component: () => import('@/pages/WishesPage.vue')
    },
    // 中文注释：心愿独立页面（创建、兑换记录），隐藏底部导航
    { path: '/wishes/create', name: 'WishCreate', component: () => import('@/pages/WishCreatePage.vue'), meta: { noNav: true } },
    { path: '/wishes/records', name: 'WishRecords', component: () => import('@/pages/WishRecordsPage.vue'), meta: { noNav: true } },
    // 中文注释：心愿编辑独立页面
    { path: '/wishes/:id/edit', name: 'WishEdit', component: () => import('@/pages/WishEditPage.vue'), meta: { noNav: true } },
    
    // --- 作业家模块 ---
    {
      path: '/homework',
      name: 'Homework',
      component: () => import('@/pages/HomeworkPage.vue')
    },
    { path: '/homework/placeholder', name: 'HomeworkPlaceholder', component: () => import('@/pages/homework/HomeworkPlaceholder.vue'), meta: { noNav: true } },
    { path: '/homework/chinese', name: 'HomeworkChinese', component: () => import('@/pages/HomeworkChinesePage.vue'), meta: { noNav: true } },
    { path: '/homework/chinese/pinyin', name: 'PinyinMenu', component: () => import('@/pages/PinyinMenuPage.vue'), meta: { noNav: true } },
    { path: '/homework/chinese/pinyin/wheel', name: 'PinyinWheel', component: () => import('@/pages/PinyinWheelPage.vue'), meta: { noNav: true } },
    
    // --- 听写大师模块 ---
    {
      path: '/dictation',
      name: 'DictationHome',
      component: () => import('@/pages/dictation/DictationHome.vue'),
      meta: { noNav: true }
    },
    { path: '/dictation/banks', name: 'WordBankList', component: () => import('@/pages/dictation/WordBankList.vue'), meta: { noNav: true } },
    { path: '/dictation/banks/create', name: 'WordBankCreate', component: () => import('@/pages/dictation/WordBankEdit.vue'), meta: { noNav: true } },
    { path: '/dictation/banks/:id', name: 'WordBankEdit', component: () => import('@/pages/dictation/WordBankEdit.vue'), meta: { noNav: true } },
    { path: '/dictation/player', name: 'DictationPlayer', component: () => import('@/pages/dictation/DictationPlayer.vue'), meta: { noNav: true } },
    { path: '/dictation/settings', name: 'DictationSettings', component: () => import('@/pages/dictation/DictationSettings.vue'), meta: { noNav: true } },
    { path: '/dictation/mistakes', name: 'MistakeBook', component: () => import('@/pages/dictation/MistakeBook.vue'), meta: { noNav: true } },
    { path: '/dictation/history', name: 'DictationHistory', component: () => import('@/pages/dictation/DictationHistory.vue'), meta: { noNav: true } },

    // --- 小成长模块 ---
    {
      path: '/little-growth',
      name: 'LittleGrowthHome',
      component: () => import('@/pages/little-growth/LittleGrowthHome.vue'),
      meta: { noNav: true }
    },
    {
      path: '/little-growth/create',
      name: 'LittleGrowthCreate',
      component: () => import('@/pages/little-growth/LittleGrowthEdit.vue'),
      meta: { noNav: true }
    },
    {
      path: '/little-growth/edit/:id',
      name: 'LittleGrowthEdit',
      component: () => import('@/pages/little-growth/LittleGrowthEdit.vue'),
      meta: { noNav: true }
    },
    {
      path: '/little-growth/tags',
      name: 'LittleGrowthTags',
      component: () => import('@/pages/little-growth/LittleGrowthTags.vue'),
      meta: { noNav: true }
    },

    {
      path: '/mine',
      name: 'Mine',
      component: () => import('@/pages/MinePage.vue')
    }
    ,
    // 中文注释：子账号管理页面（隐藏底部导航）
    { path: '/mine/subaccounts', name: 'SubAccounts', component: () => import('../pages/SubAccountsPage.vue'), meta: { noNav: true } },
    // 中文注释：设置页面独立子路由
    { path: '/settings/profile', name: 'SettingsProfile', component: () => import('@/pages/SettingsProfilePage.vue'), meta: { noNav: true } },
    { path: '/settings/tomato', name: 'SettingsTomato', component: () => import('@/pages/SettingsTomatoPage.vue'), meta: { noNav: true } },
    { path: '/settings/tasks', name: 'SettingsTasks', component: () => import('@/pages/SettingsTasksPage.vue'), meta: { noNav: true } },
    { path: '/settings/reading', name: 'SettingsReading', component: () => import('@/pages/SettingsReadingPage.vue'), meta: { noNav: true } },
    { path: '/settings/subjects', name: 'SettingsSubjects', component: () => import('@/pages/SettingsSubjectsPage.vue'), meta: { noNav: true } },
    { path: '/settings/coins', name: 'SettingsCoins', component: () => import('@/pages/SettingsCoinsPage.vue'), meta: { noNav: true } },
    { path: '/settings/about', name: 'SettingsAbout', component: () => import('@/pages/SettingsAboutPage.vue'), meta: { noNav: true } },
    { path: '/settings/help', name: 'Help', component: () => import('@/pages/HelpPage.vue'), meta: { noNav: true } },
    { path: '/settings/support', name: 'Support', component: () => import('@/pages/SupportPage.vue'), meta: { noNav: true } },
    { path: '/settings', name: 'Settings', component: () => import('@/pages/SettingsPage.vue'), meta: { noNav: true } },
    { path: '/settings/appearance', name: 'SettingsAppearance', component: () => import('@/pages/SettingsAppearancePage.vue'), meta: { noNav: true } }
    ,
    // 管理后台页面，仅允许用户ID=1访问
    { path: '/admin', name: 'Admin', component: () => import('@/pages/AdminPage.vue'), meta: { noNav: true } }
  ]
})

// 中文注释：全局前置守卫——未登录禁止进入业务页面
router.beforeEach(async (to) => {
  const auth = useAuth()
  const store = useAppState()
  const isPublic = to.meta?.public === true
  // 中文注释：登录后首次进入业务页面时，确保全局金币从用户信息恢复
  try {
    const userCoins = Number((auth.user as any)?.coins ?? NaN)
    if (!isNaN(userCoins)) {
      // 仅当值不一致时才更新，避免不必要的持久化
      if (Number(store.coins) !== userCoins) {
        store.setCoins(userCoins)
      }
    }
  } catch {}
  // 中文注释：若为子账号，且启用父子金币同步（由后端控制），在进入业务页时主动请求一次 /api/coins
  try {
    const isChild = !!(auth.user && (auth.user as any).parent_id)
    // 中文注释：取消“一次性检查”与 0 值条件，改为每次进入业务页都拉取最新金币，确保父子同步即时生效
    if (isChild && !isPublic) {
      await http.get('/coins')
    }
  } catch {}
  // 中文注释：已登录访问登录/注册页面时自动跳转任务页
  if (auth.token && (to.path === '/login' || to.path === '/register')) {
    return { path: '/tasks' }
  }
  // 备注页面需 VIP 权限：非VIP访问时提示并跳转任务页
  const isNotesPage = /^\/tasks\/[0-9]+\/notes$/.test(to.path)
  const u = auth.user as any
  const lifetime = !!(u?.is_lifetime_vip)
  const vip = !!(u?.is_vip)
  const expire = u?.vip_expire_time ? new Date(u.vip_expire_time as string) : null
  const isVIP = !!u && (lifetime || (vip && !!expire && expire.getTime() > Date.now()))
  if (isNotesPage && !isVIP) {
    try { ElMessage.warning('该功能需要开通VIP会员才能使用，添加微信：tabor2024，备注“任务家”') } catch {}
    return { path: '/tasks' }
  }
  // 管理后台访问控制：仅允许用户ID=1
  if (to.path === '/admin') {
    const id = Number((auth.user as any)?.id || 0)
    if (id !== 1) {
      try { ElMessage.warning('仅管理员可访问该页面') } catch {}
      return { path: '/tasks' }
    }
  }
  // 中文注释：未登录访问私有页面时跳转至登录，并记录重定向
  if (!isPublic && !auth.token) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router
