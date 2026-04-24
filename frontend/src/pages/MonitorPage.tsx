import React, { useEffect, useState, useMemo } from 'react'
import { useAppStore } from '../stores/useAppStore'
import { cn, getPlatformName, getPlatformIcon, getStatusColor, getStatusLabel, getProgressColor } from '../utils/helpers'
import { Activity, RefreshCw, ChevronDown, ChevronUp, Terminal } from 'lucide-react'
import { AnimatePresence } from 'framer-motion'

const MonitorCard: React.FC<{
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
}> = ({ uid, accountName, platform, platformName, status, progress, currentTask, totalCourses, doneCourses, logs, errorMessage }) => {
  const [showLogs, setShowLogs] = useState(false)

  return (
    <div className="glass-card overflow-hidden">
      <div className="p-4">
        <div className="flex items-center gap-4">
          <div className="w-10 h-10 rounded-lg bg-dark-700 flex items-center justify-center text-lg">
            {getPlatformIcon(platform)}
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="font-medium text-dark-100">{accountName}</span>
              <span className="text-xs text-dark-500">{platformName}</span>
              <span className={cn("text-xs px-2 py-0.5 rounded-full", getStatusColor(status),
                status === 'running' ? 'bg-accent-600/10' :
                status === 'completed' ? 'bg-emerald-600/10' :
                status === 'error' ? 'bg-red-600/10' : 'bg-dark-700'
              )}>
                {getStatusLabel(status)}
              </span>
            </div>

            <div className="mt-2">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-dark-400">{currentTask}</span>
                <span className="text-xs text-dark-400">
                  {doneCourses ?? 0}/{totalCourses ?? 0} 门课程 · {(progress ?? 0).toFixed(1)}%
                </span>
              </div>
              <div className="w-full h-2 bg-dark-700 rounded-full overflow-hidden">
                <div
                  className={cn(
                    "h-full rounded-full transition-all duration-500 bg-gradient-to-r",
                    getProgressColor(progress),
                    status === 'running' && "progress-striped"
                  )}
                  style={{ width: `${Math.min(progress ?? 0, 100)}%` }}
                />
              </div>
            </div>

            {errorMessage && (
              <p className="text-xs text-red-400 mt-2">❌ {errorMessage}</p>
            )}
          </div>

          <button onClick={() => setShowLogs(!showLogs)} className="btn-secondary text-xs px-2 py-1.5">
            {showLogs ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
          </button>
        </div>
      </div>

      <AnimatePresence>
        {showLogs && (
          <div className="border-t border-dark-800/50">
            <div className="p-3 bg-dark-950/50 max-h-48 overflow-y-auto">
              <div className="flex items-center gap-2 mb-2">
                <Terminal size={12} className="text-dark-500" />
                <span className="text-xs text-dark-500 font-medium">运行日志</span>
              </div>
              <div className="space-y-0.5 font-mono text-xs">
                {(!logs || logs.length === 0) ? (
                  <p className="text-dark-600">暂无日志</p>
                ) : (
                  logs.map((log, i) => (
                    <p key={i} className="text-dark-400 leading-relaxed">
                      <span className="text-dark-600 mr-2">[{String(i + 1).padStart(3, '0')}]</span>
                      {log}
                    </p>
                  ))
                )}
              </div>
            </div>
          </div>
        )}
      </AnimatePresence>
    </div>
  )
}

const MonitorPage: React.FC = () => {
  const { progressMap, fetchProgress, accounts, fetchAccounts } = useAppStore()
  const [filter, setFilter] = useState<string>('all')

  useEffect(() => {
    fetchAccounts()
    fetchProgress()
    const interval = setInterval(() => {
      fetchProgress()
      fetchAccounts()
    }, 3000)
    return () => clearInterval(interval)
  }, [])

  const mergedList = useMemo(() => {
    const progressList = Object.values(progressMap)
    const runningAccounts = accounts.filter(a => a.isRunning && !progressMap[a.uid])
    const list = [
      ...progressList,
      ...runningAccounts.map(a => ({
        uid: a.uid,
        accountName: a.remarkName || a.account,
        platform: a.accountType,
        platformName: getPlatformName(a.accountType),
        status: 'running' as const,
        progress: 0,
        currentTask: '正在启动...',
        totalCourses: 0,
        doneCourses: 0,
        logs: [],
      })),
    ]
    list.sort((a, b) => a.uid.localeCompare(b.uid))
    return list
  }, [progressMap, accounts])

  const filtered = useMemo(() => {
    if (filter === 'all') return mergedList
    return mergedList.filter(p => p.status === filter)
  }, [mergedList, filter])

  const runningCount = mergedList.filter(p => p.status === 'running').length
  const completedCount = mergedList.filter(p => p.status === 'completed').length
  const errorCount = mergedList.filter(p => p.status === 'error').length

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-dark-50">任务监控</h1>
          <p className="text-sm text-dark-400 mt-1">实时监控所有用户任务进度</p>
        </div>
        <button onClick={fetchProgress} className="btn-secondary flex items-center gap-2">
          <RefreshCw size={14} />刷新
        </button>
      </div>

      <div className="flex items-center gap-2">
        {[
          { id: 'all', label: '全部', count: mergedList.length },
          { id: 'running', label: '运行中', count: runningCount },
          { id: 'completed', label: '已完成', count: completedCount },
          { id: 'error', label: '异常', count: errorCount },
        ].map(f => (
          <button key={f.id}
            onClick={() => setFilter(f.id)}
            className={cn(
              "px-3 py-1.5 rounded-lg text-xs font-medium transition-all",
              filter === f.id
                ? "bg-accent-600/20 text-accent-400 border border-accent-600/30"
                : "bg-dark-800 text-dark-400 hover:text-dark-200 border border-dark-700"
            )}
          >
            {f.label} ({f.count})
          </button>
        ))}
      </div>

      {filtered.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <Activity size={48} className="mx-auto text-dark-600 mb-4" />
          <p className="text-dark-400 text-sm">
            {mergedList.length === 0 ? '暂无监控数据，请先启动刷课任务' : '当前筛选条件下没有任务'}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map((p) => (
            <MonitorCard key={p.uid} {...p} />
          ))}
        </div>
      )}
    </div>
  )
}

export default MonitorPage
