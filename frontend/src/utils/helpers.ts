import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export const APP_VERSION = '1.0.2'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function parseResponse<T>(jsonStr: string): { code: number; message: string; data: T } {
  try {
    return JSON.parse(jsonStr)
  } catch {
    return { code: 400, message: '解析响应失败', data: null as T }
  }
}

export const PLATFORM_LIST = [
  { id: 'XUEXITONG', name: '学习通', icon: '📖' },
  { id: 'YINGHUA', name: '英华学堂', icon: '🎓' },
  { id: 'CANGHUI', name: '仓辉实训', icon: '🔧' },
  { id: 'ENAEA', name: '学习公社', icon: '📚' },
  { id: 'CQIE', name: '重庆工程学院', icon: '🏫' },
  { id: 'KETANGX', name: '码上研训', icon: '💻' },
  { id: 'ICVE', name: '智慧职教', icon: '🎯' },
  { id: 'QSXT', name: '青书学堂', icon: '📗' },
  { id: 'WELEARN', name: 'WeLearn', icon: '🌐' },
  { id: 'HQKJ', name: '海旗科技', icon: '🚀' },
  { id: 'GONGXUE', name: '工学云', icon: '☁️' },
  { id: 'WEIBAN', name: '安全微伴', icon: '🛡️' },
]

export const AI_TYPE_LIST = [
  { id: 'TONGYI', name: '通义千问' },
  { id: 'CHATGLM', name: '智谱ChatGLM' },
  { id: 'XINGHUO', name: '讯飞星火' },
  { id: 'DOUBAO', name: '豆包' },
  { id: 'OPENAI', name: 'OpenAI' },
  { id: 'DEEPSEEK', name: 'DeepSeek' },
  { id: 'SILICON', name: '硅基流动' },
  { id: 'METAAI', name: '秘塔AI' },
  { id: 'OTHER', name: '其他兼容' },
]

export function getPlatformName(id: string): string {
  return PLATFORM_LIST.find(p => p.id === id)?.name ?? id
}

export function getPlatformIcon(id: string): string {
  return PLATFORM_LIST.find(p => p.id === id)?.icon ?? '📱'
}

export function getAiTypeName(id: string): string {
  return AI_TYPE_LIST.find(a => a.id === id)?.name ?? id
}

export function getStatusColor(status: string): string {
  switch (status) {
    case 'idle': return 'text-dark-400'
    case 'logging': return 'text-yellow-400'
    case 'running': return 'text-accent-400'
    case 'paused': return 'text-orange-400'
    case 'completed': return 'text-emerald-400'
    case 'error': return 'text-red-400'
    case 'stopped': return 'text-dark-400'
    default: return 'text-dark-400'
  }
}

export function getStatusLabel(status: string): string {
  switch (status) {
    case 'idle': return '空闲'
    case 'logging': return '登录中'
    case 'running': return '运行中'
    case 'paused': return '已暂停'
    case 'completed': return '已完成'
    case 'error': return '错误'
    case 'stopped': return '已停止'
    default: return '未知'
  }
}

export function getProgressColor(progress: number): string {
  if (progress >= 100) return 'from-emerald-500 to-emerald-400'
  if (progress >= 60) return 'from-accent-500 to-accent-400'
  if (progress >= 30) return 'from-yellow-500 to-yellow-400'
  return 'from-orange-500 to-orange-400'
}
