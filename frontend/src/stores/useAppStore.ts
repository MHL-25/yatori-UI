import { create } from 'zustand'
import { parseResponse } from '../utils/helpers'
import { GetAccountList, AddAccount, DeleteAccount, UpdateAccount, LoginCheck, StartBrush, StopBrush, PauseBrush, GetAllProgress, GetConfig, SaveConfig, ImportConfigFile, StartAllBrush, StartBatchBrush } from '../../wailsjs/go/main/App'

interface Account {
  uid: string
  accountType: string
  url: string
  remarkName: string
  account: string
  password: string
  userConfigJson: string
  isRunning: boolean
}

interface TaskProgress {
  uid: string
  accountName: string
  platform: string
  platformName: string
  status: string
  progress: number
  currentTask: string
  totalCourses: number
  doneCourses: number
  logs: string[]
  errorMessage?: string
}

interface AppConfig {
  setting: {
    basicSetting: { completionTone: number; colorLog: number; logOutFileSw: number; logLevel: string; logModel: number }
    aiSetting: { aiType: string; aiUrl: string; model: string; API_KEY: string }
    visionAiSetting: { aiType: string; aiUrl: string; model: string; API_KEY: string }
    apiQueSetting: { url: string }
    emailInform: { sw: number; smtpHost: string; smtpPort: number; userName: string; password: string }
  }
  users: any[]
}

interface AppStore {
  accounts: Account[]
  progressMap: Record<string, TaskProgress>
  config: AppConfig | null
  currentPage: string
  loading: boolean
  sidebarCollapsed: boolean

  setCurrentPage: (page: string) => void
  toggleSidebar: () => void
  fetchAccounts: () => Promise<void>
  addAccount: (data: any) => Promise<{ ok: boolean; msg: string }>
  deleteAccount: (uid: string) => Promise<{ ok: boolean; msg: string }>
  updateAccount: (data: any) => Promise<{ ok: boolean; msg: string }>
  loginCheck: (uid: string) => Promise<{ ok: boolean; msg: string }>
  startBrush: (uid: string) => Promise<{ ok: boolean; msg: string }>
  stopBrush: (uid: string) => Promise<{ ok: boolean; msg: string }>
  pauseBrush: (uid: string) => Promise<{ ok: boolean; msg: string }>
  startAllBrush: () => Promise<{ ok: boolean; msg: string; success: number; total: number }>
  startBatchBrush: (uids: string[]) => Promise<{ ok: boolean; msg: string; success: number; total: number }>
  fetchProgress: () => Promise<void>
  fetchConfig: () => Promise<void>
  saveConfig: (cfg: any) => Promise<{ ok: boolean; msg: string }>
  importConfig: () => Promise<{ ok: boolean; msg: string; count: number }>
  updateProgress: (data: TaskProgress) => void
  addLog: (uid: string, log: string) => void
}

export const useAppStore = create<AppStore>((set, get) => ({
  accounts: [],
  progressMap: {},
  config: null,
  currentPage: 'dashboard',
  loading: false,
  sidebarCollapsed: false,

  setCurrentPage: (page) => set({ currentPage: page }),
  toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),

  fetchAccounts: async () => {
    set({ loading: true })
    try {
      const result = await GetAccountList()
      const resp = parseResponse<Account[]>(result)
      if (resp.code === 200) {
        set({ accounts: resp.data || [], loading: false })
      } else {
        set({ loading: false })
      }
    } catch (e) {
      console.error('fetchAccounts error:', e)
      set({ loading: false })
    }
  },

  addAccount: async (data) => {
    try {
      const result = await AddAccount(JSON.stringify(data))
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK' }
      }
      return { ok: false, msg: resp.message || 'Failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  deleteAccount: async (uid) => {
    try {
      const result = await DeleteAccount(uid)
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK' }
      }
      return { ok: false, msg: resp.message || 'Failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  updateAccount: async (data) => {
    try {
      const result = await UpdateAccount(JSON.stringify(data))
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK' }
      }
      return { ok: false, msg: resp.message || 'Failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  loginCheck: async (uid) => {
    try {
      const result = await LoginCheck(uid)
      const resp = parseResponse<any>(result)
      return resp.code === 200 ? { ok: true, msg: 'OK' } : { ok: false, msg: resp.message || 'Login failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  startBrush: async (uid) => {
    try {
      const result = await StartBrush(uid)
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK' }
      }
      return { ok: false, msg: resp.message || 'Failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  stopBrush: async (uid) => {
    try {
      const result = await StopBrush(uid)
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK' }
      }
      return { ok: false, msg: resp.message || 'Failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  pauseBrush: async (uid) => {
    try {
      const result = await PauseBrush(uid)
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK' }
      }
      return { ok: false, msg: resp.message || 'Failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  startAllBrush: async () => {
    try {
      const result = await StartAllBrush()
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK', success: resp.data?.success || 0, total: resp.data?.total || 0 }
      }
      return { ok: false, msg: resp.message || 'Failed', success: 0, total: 0 }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error', success: 0, total: 0 }
    }
  },

  startBatchBrush: async (uids) => {
    try {
      const result = await StartBatchBrush(JSON.stringify(uids))
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        return { ok: true, msg: 'OK', success: resp.data?.success || 0, total: resp.data?.total || 0 }
      }
      return { ok: false, msg: resp.message || 'Failed', success: 0, total: 0 }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error', success: 0, total: 0 }
    }
  },

  fetchProgress: async () => {
    try {
      const result = await GetAllProgress()
      const resp = parseResponse<TaskProgress[]>(result)
      if (resp.code === 200 && resp.data) {
        const map: Record<string, TaskProgress> = {}
        resp.data.forEach((p) => { map[p.uid] = p })
        set({ progressMap: map })
      }
    } catch (e) {
      console.error('fetchProgress error:', e)
    }
  },

  fetchConfig: async () => {
    try {
      const result = await GetConfig()
      const resp = parseResponse<AppConfig>(result)
      if (resp.code === 200) {
        set({ config: resp.data })
      }
    } catch (e) {
      console.error('fetchConfig error:', e)
    }
  },

  saveConfig: async (cfg) => {
    try {
      const result = await SaveConfig(JSON.stringify(cfg))
      const resp = parseResponse<any>(result)
      return resp.code === 200 ? { ok: true, msg: 'OK' } : { ok: false, msg: resp.message || 'Save failed' }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error' }
    }
  },

  importConfig: async () => {
    try {
      const result = await ImportConfigFile()
      const resp = parseResponse<any>(result)
      if (resp.code === 200) {
        await get().fetchAccounts()
        await get().fetchConfig()
        return { ok: true, msg: `成功导入 ${resp.data?.count || 0} 个账号`, count: resp.data?.count || 0 }
      }
      return { ok: false, msg: resp.message || '导入失败', count: 0 }
    } catch (e: any) {
      return { ok: false, msg: e?.message || 'Network error', count: 0 }
    }
  },

  updateProgress: (data) => {
    set((s) => ({
      progressMap: { ...s.progressMap, [data.uid]: { ...s.progressMap[data.uid], ...data } }
    }))
  },

  addLog: (uid: string, log: string) => {
    set((s) => {
      const existing = s.progressMap[uid]
      if (!existing) return s
      const logs = [...(existing.logs || []), log]
      if (logs.length > 200) logs.splice(0, logs.length - 200)
      return {
        progressMap: { ...s.progressMap, [uid]: { ...existing, logs } }
      }
    })
  },
}))
