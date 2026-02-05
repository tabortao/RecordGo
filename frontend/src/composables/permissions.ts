// 中文注释：权限组合式，统一解析与判断当前用户的权限
import { computed } from 'vue'
import { useAuth } from '@/stores/auth'

// 中文注释：规范化后的权限结构类型（后端返回为 JSON 字符串，这里解析为对象）
export interface NormalizedPerms {
  view_only?: boolean
  manage_children?: boolean
  account?: { manage_children?: boolean }
  tasks?: { create?: boolean; edit?: boolean; delete?: boolean; status?: boolean; undo?: boolean }
  wishes?: { create?: boolean; edit?: boolean; delete?: boolean; exchange?: boolean }
  // 中文注释：新增“设置”权限分组（番茄钟/任务/金币/朗读），用于控制“我的”页设置入口
  settings?: { tomato?: boolean; tasks?: boolean; coins?: boolean; reading?: boolean }
}

// 中文注释：解析 JSON 字符串为权限对象，失败时返回默认空对象
function parsePerms(raw: string | object | undefined | null): NormalizedPerms {
  if (!raw) return {}
  if (typeof raw === 'object') return raw as NormalizedPerms
  try {
    const obj = JSON.parse(String(raw))
    if (obj && typeof obj === 'object') return obj as NormalizedPerms
    return {}
  } catch { return {} }
}

// 中文注释：统一的权限判断：父账号默认放行；子账号按 view_only 与具体模块动作判断
export function usePermissions() {
  const auth = useAuth()
  const user = computed(() => auth.user)
  const isParent = computed(() => {
    const u = user.value
    // 无 parent_id 视为主账号；0 或 null 也视为主账号
    return !u || !u.parent_id || Number(u.parent_id) === 0
  })
  const perms = computed(() => parsePerms(user.value?.permissions))
  const viewOnly = computed(() => Boolean(perms.value.view_only))
  const manageChildren = computed(() => isParent.value || Boolean(perms.value.manage_children) || Boolean(perms.value.account?.manage_children))

  function has(category: keyof NormalizedPerms, action?: string): boolean {
    if (isParent.value) return true
    if (viewOnly.value) return false
    if (!action) return Boolean((perms.value as any)[category])
    const cat = (perms.value as any)[category]
    if (!cat || typeof cat !== 'object') return false
    return Boolean(cat[action])
  }

  // 中文注释：常用动作的便捷布尔值，避免在页面多次拼接字符串
  const canTaskCreate = computed(() => has('tasks', 'create'))
  const canTaskEdit = computed(() => has('tasks', 'edit'))
  const canTaskDelete = computed(() => has('tasks', 'delete'))
  const canTaskStatus = computed(() => has('tasks', 'status'))
  const canTaskUndo = computed(() => has('tasks', 'undo'))

  const canWishCreate = computed(() => has('wishes', 'create'))
  const canWishEdit = computed(() => has('wishes', 'edit'))
  const canWishDelete = computed(() => has('wishes', 'delete'))
  const canWishExchange = computed(() => has('wishes', 'exchange'))

  // 中文注释：设置权限的判断逻辑：
  // - 父账号始终允许；
  // - 子账号：
  //   * 若 view_only=true，则“任务设置”“金币设置”不可用（即使 settings 中为 true）；
  //   * “番茄钟设置”“朗读设置”在仅查看模式下若 settings 对应为 true，则允许打开；
  //   * 非仅查看模式下，按 settings 对应布尔值控制。
  const settingsGroup = computed(() => perms.value.settings || {})
  const canSettingTomato = computed(() => {
    if (isParent.value) return true
    const allow = Boolean(settingsGroup.value.tomato)
    if (viewOnly.value) return allow
    return allow
  })
  const canSettingTasks = computed(() => {
    if (isParent.value) return true
    const allow = Boolean(settingsGroup.value.tasks)
    if (viewOnly.value) return false // 仅查看模式任务设置不可用
    return allow
  })
  const canSettingCoins = computed(() => {
    if (isParent.value) return true
    const allow = Boolean(settingsGroup.value.coins)
    if (viewOnly.value) return false // 仅查看模式金币设置不可用
    return allow
  })
  const canSettingReading = computed(() => {
    if (isParent.value) return true
    const allow = Boolean(settingsGroup.value.reading)
    if (viewOnly.value) return allow
    return allow
  })

  return {
    isParent,
    perms,
    viewOnly,
    manageChildren,
    has,
    // 任务相关
    canTaskCreate,
    canTaskEdit,
    canTaskDelete,
    canTaskStatus,
    canTaskUndo,
    // 心愿相关
    canWishCreate,
    canWishEdit,
    canWishDelete,
    canWishExchange,
    // 设置相关
    canSettingTomato,
    canSettingTasks,
    canSettingCoins,
    canSettingReading,
  }
}
