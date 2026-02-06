import { ref, computed } from 'vue'
import dayjs from 'dayjs'
import { listTasks, listTaskOccurrences, type TaskItem } from '@/services/tasks'

export function useTodayStats() {
  const loading = ref(false)
  const total = ref(0)
  const completed = ref(0)

  // Duplicate logic from TasksPage to generate repeat dates
  function generateRepeatDates(start: Date, end: Date | undefined, type: 'none' | 'daily' | 'weekdays' | 'weekly' | 'monthly', weeklyDays: number[]) {
    const out: Date[] = []
    if (!end || type === 'none') return out
    const s = dayjs(start).startOf('day')
    const e = dayjs(end).startOf('day')
    if (e.isBefore(s)) return out
    
    if (type === 'daily') {
      let d = s.clone()
      while (!d.isAfter(e)) { out.push(d.toDate()); d = d.add(1, 'day') }
    } else if (type === 'weekdays') {
      let d = s.clone()
      while (!d.isAfter(e)) {
        const w = d.day()
        if (w >= 1 && w <= 5) out.push(d.toDate())
        d = d.add(1, 'day')
      }
    } else if (type === 'weekly') {
      const set = new Set(weeklyDays || [])
      let d = s.clone()
      while (!d.isAfter(e)) {
        const w = d.day() === 0 ? 7 : d.day()
        if (set.has(w)) out.push(d.toDate())
        d = d.add(1, 'day')
      }
    } else if (type === 'monthly') {
      let d = s.clone()
      const dayOfMonth = s.date()
      while (!d.isAfter(e)) {
        const candidate = d.date(dayOfMonth)
        if (candidate.month() === d.month() && !candidate.isAfter(e)) out.push(candidate.toDate())
        d = d.add(1, 'month')
      }
    }
    return out
  }

  function isRepeatTask(t: TaskItem) {
    const rep = String((t as any).repeat || (t as any).repeat_type || 'none').trim().toLowerCase()
    const type = /^(daily|weekdays|weekly|monthly)$/.test(rep) ? rep : 'none'
    return type !== 'none'
  }

  function getDailyMaxCheckins(t: TaskItem) {
    const v = Number((t as any).daily_max_checkins ?? 1)
    return Number.isFinite(v) && v > 1 ? Math.floor(v) : 1
  }

  function useOccurrenceTracking(t: TaskItem) {
    return isRepeatTask(t) || getDailyMaxCheckins(t) > 1
  }

  function getCompletedKey(t: TaskItem): string | null {
    const raw = (t as any).completed_at ?? (t as any).CompletedAt ?? (t as any).completedAt
    if (!raw) return null
    const key = dayjs(String(raw)).format('YYYY-MM-DD')
    return key && key !== 'Invalid Date' ? key : null
  }

  async function refresh() {
    loading.value = true
    try {
      const todayStr = dayjs().format('YYYY-MM-DD')
      
      // 1. Fetch all tasks
      const resTasks = await listTasks({ page_size: 1000 })
      const tasks = resTasks.items || []

      // 2. Fetch occurrences for today
      const resOccur = await listTaskOccurrences({ date: todayStr })
      const occurMap: Record<number, { status: number; checkins_count: number }> = {}
      for (const it of resOccur.items || []) {
        occurMap[it.task_id] = { 
          status: Number(it.status || 0), 
          checkins_count: Number((it as any).checkins_count ?? 0) 
        }
      }

      // 3. Filter for today
      const todayTasks = tasks.filter(t => {
        // Exclude deleted occurrences
        const occStatus = occurMap[t.id]?.status
        if (occStatus === -1) return false

        const sDateStr = t.start_date ? String(t.start_date) : ''
        const eDateStr = t.end_date ? String(t.end_date) : ''
        if (!sDateStr) return false
        
        const sDate = dayjs(sDateStr).toDate()
        const eDate = eDateStr ? dayjs(eDateStr).toDate() : undefined
        const rawRepeat = (t as any).repeat || 'none'
        const rep = String(rawRepeat).toLowerCase()
        const type: 'none'|'daily'|'weekdays'|'weekly'|'monthly' =
          /none|无|^$/i.test(rep) ? 'none' :
          /daily|每天/i.test(rep) ? 'daily' :
          /weekdays|工作日/i.test(rep) ? 'weekdays' :
          /weekly|每周/i.test(rep) ? 'weekly' :
          /monthly|每月/i.test(rep) ? 'monthly' : 'none'

        // Logic matched from TasksPage
        if (type === 'none') {
          if (eDate) {
            const sKey = dayjs(sDate).format('YYYY-MM-DD')
            const eKey = dayjs(eDate).format('YYYY-MM-DD')
            if (Number((t as any).status ?? 0) === 2) {
              const cKey = getCompletedKey(t) || sKey
              return todayStr === cKey
            }
            return todayStr >= sKey && todayStr <= eKey
          }
          return dayjs(sDate).format('YYYY-MM-DD') === todayStr
        }

        // Repeating logic
        const dow = dayjs(sDate).day() === 0 ? 7 : dayjs(sDate).day()
        const weeklyDays: number[] = Array.isArray((t as any).weekly_days) ? ((t as any).weekly_days as number[]) : [dow]

        if (!eDate) {
           const d = dayjs()
           const w = d.day() === 0 ? 7 : d.day()
           const sd = dayjs(sDate)
           if (d.isBefore(sd.startOf('day'))) return false
           if (type === 'daily') return true
           if (type === 'weekdays') return w >= 1 && w <= 5
           if (type === 'weekly') {
             const set = new Set(weeklyDays || [])
             const sw = sd.day() === 0 ? 7 : sd.day()
             return set.size ? set.has(w) : (w === sw)
           }
           if (type === 'monthly') return d.date() === sd.date()
           return false
        }

        const dates = generateRepeatDates(sDate, eDate, type, weeklyDays)
        const keys = new Set(dates.map((d) => dayjs(d).format('YYYY-MM-DD')))
        return keys.has(todayStr)
      })

      // 4. Calculate stats
      total.value = todayTasks.length
      completed.value = todayTasks.filter(t => {
        if (useOccurrenceTracking(t)) {
          const occ = occurMap[t.id]
          const max = getDailyMaxCheckins(t)
          // Default checkins count is 0 if not found
          let cnt = occ?.checkins_count ?? 0
          // Fallback if checkins_count is 0 but status is 2 (legacy)
          if (cnt === 0 && occ?.status === 2) cnt = 1
          
          return cnt >= max
        }
        return t.status === 2
      }).length

    } catch (e) {
      console.error('Failed to fetch stats', e)
    } finally {
      loading.value = false
    }
  }

  return {
    total,
    completed,
    loading,
    refresh
  }
}
