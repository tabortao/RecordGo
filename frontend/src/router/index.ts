import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/stores/auth'

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
    {
      path: '/wishes',
      name: 'Wishes',
      component: () => import('@/pages/WishesPage.vue')
    },
    {
      path: '/mine',
      name: 'Mine',
      component: () => import('@/pages/MinePage.vue')
    }
    ,
    // 中文注释：设置页面独立子路由
    { path: '/settings/profile', name: 'SettingsProfile', component: () => import('@/pages/SettingsProfilePage.vue'), meta: { noNav: true } },
    { path: '/settings/tomato', name: 'SettingsTomato', component: () => import('@/pages/SettingsTomatoPage.vue'), meta: { noNav: true } },
    { path: '/settings/tasks', name: 'SettingsTasks', component: () => import('@/pages/SettingsTasksPage.vue'), meta: { noNav: true } },
    { path: '/settings/reading', name: 'SettingsReading', component: () => import('@/pages/SettingsReadingPage.vue'), meta: { noNav: true } },
    { path: '/settings/subjects', name: 'SettingsSubjects', component: () => import('@/pages/SettingsSubjectsPage.vue'), meta: { noNav: true } },
    { path: '/settings/coins', name: 'SettingsCoins', component: () => import('@/pages/SettingsCoinsPage.vue'), meta: { noNav: true } },
    { path: '/settings/about', name: 'SettingsAbout', component: () => import('@/pages/SettingsAboutPage.vue'), meta: { noNav: true } }
  ]
})

// 中文注释：全局前置守卫——未登录禁止进入业务页面
router.beforeEach((to) => {
  const auth = useAuth()
  const isPublic = to.meta?.public === true
  if (!isPublic && !auth.token) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
})

export default router
