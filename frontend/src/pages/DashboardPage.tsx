import React, { useEffect } from 'react'
import { useAppStore } from '../stores/useAppStore'
import { cn, getPlatformName, getPlatformIcon, getStatusColor, getStatusLabel, getProgressColor } from '../utils/helpers'
import { Users, Activity, Play, Square, CheckCircle, AlertCircle, Clock, Zap, Pause } from 'lucide-react'
import { motion } from 'framer-motion'

const StatCard: React.FC<{ icon: React.ReactNode; label: string; value: string | number; color: string; sub?: string }> = ({ icon, label, value, color, sub }) => (
  <div className="glass-card p-4 flex items-center gap-4">
    <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center", color)}>
      {icon}
    </div>
    <div>
      <p className="text-xs text-dark-400 font-medium">{label}</p>
      <p className="text-xl font-bold text-dark-100">{value}</p>
      {sub && <p className="text-xs text-dark-500">{sub}</p>}
    </div>
  </div>
)

const DashboardPage: React.FC = () => {
  const { accounts, progressMap, setCurrentPage, startBrush, stopBrush, pauseBrush, fetchAccounts, fetchProgress } = useAppStore()

  useEffect(() => {
    fetchAccounts()
    fetchProgress()
    const interval = setInterval(() => {
      fetchAccounts()
      fetchProgress()
    }, 3000)
    return () => clearInterval(interval)
  }, [])

  const runningFromProgress = Object.values(progressMap).filter(p => p.status === 'running').length
  const runningFromAccounts = accounts.filter(a => a.isRunning).length
  const runningCount = Math.max(runningFromProgress, runningFromAccounts)

  const completedCount = Object.values(progressMap).filter(p => p.status === 'completed').length
  const errorCount = Object.values(progressMap).filter(p => p.status === 'error').length

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-dark-50">仪表盘</h1>
          <p className="text-sm text-dark-400 mt-1">实时概览所有网课任务状态</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          <span className="text-xs text-dark-400">系统运行中</span>
        </div>
      </div>

      <div className="grid grid-cols-4 gap-4">
        <StatCard
          icon={<Users size={20} className="text-accent-400" />}
          label="总账号数"
          value={accounts.length}
          color="bg-accent-600/10"
          sub="已添加的平台账号"
        />
        <StatCard
          icon={<Play size={20} className="text-neon-green" />}
          label="运行中"
          value={runningCount}
          color="bg-emerald-600/10"
          sub="正在刷课的账号"
        />
        <StatCard
          icon={<CheckCircle size={20} className="text-emerald-400" />}
          label="已完成"
          value={completedCount}
          color="bg-emerald-600/10"
          sub="已完成所有课程"
        />
        <StatCard
          icon={<AlertCircle size={20} className="text-red-400" />}
          label="异常"
          value={errorCount}
          color="bg-red-600/10"
          sub="需要关注的账号"
        />
      </div>

      <div className="glass-card p-5">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-dark-100">账号状态概览</h2>
          <button
            onClick={() => setCurrentPage('accounts')}
            className="text-xs text-accent-400 hover:text-accent-300 transition-colors"
          >
            管理账号 →
          </button>
        </div>

        {accounts.length === 0 ? (
          <div className="text-center py-12">
            <Zap size={48} className="mx-auto text-dark-600 mb-4" />
            <p className="text-dark-400 text-sm">暂无账号，请先添加账号</p>
            <button
              onClick={() => setCurrentPage('accounts')}
              className="btn-primary mt-4"
            >
              添加账号
            </button>
          </div>
        ) : (
          <div className="space-y-3">
            {accounts.map((account, index) => {
              const progress = progressMap[account.uid]
              const status = progress?.status || (account.isRunning ? 'running' : 'idle')
              const progressValue = progress?.progress || 0

              return (
                <motion.div
                  key={account.uid}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.05 }}
                  className="flex items-center gap-4 p-3 rounded-lg bg-dark-800/30 hover:bg-dark-800/50 transition-colors"
                >
                  <div className="w-8 h-8 rounded-lg bg-dark-700 flex items-center justify-center text-sm">
                    {getPlatformIcon(account.accountType)}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium text-dark-100 truncate">
                        {account.remarkName || account.account}
                      </span>
                      <span className="text-xs text-dark-500">{getPlatformName(account.accountType)}</span>
                    </div>
                    {progress && (
                      <div className="mt-1.5">
                        <div className="flex items-center justify-between mb-1">
                          <span className={cn("text-xs", getStatusColor(status))}>
                            {getStatusLabel(status)}
                          </span>
                          <span className="text-xs text-dark-500">{progress.currentTask}</span>
                        </div>
                        <div className="w-full h-1.5 bg-dark-700 rounded-full overflow-hidden">
                          <div
                            className={cn(
                              "h-full rounded-full transition-all duration-500 bg-gradient-to-r",
                              getProgressColor(progressValue),
                              status === 'running' && "progress-striped"
                            )}
                            style={{ width: `${Math.min(progressValue, 100)}%` }}
                          />
                        </div>
                      </div>
                    )}
                  </div>
                  <div className="flex items-center gap-2">
                    {status === 'running' ? (
                      <>
                        <button onClick={() => pauseBrush(account.uid)} className="btn-secondary text-xs px-3 py-1.5">
                          <Pause size={12} className="mr-1 inline" />暂停
                        </button>
                        <button onClick={() => stopBrush(account.uid)} className="btn-danger text-xs px-3 py-1.5">
                          <Square size={12} className="mr-1 inline" />停止
                        </button>
                      </>
                    ) : (
                      <button onClick={() => startBrush(account.uid)} className="btn-success text-xs px-3 py-1.5">
                        <Play size={12} className="mr-1 inline" />启动
                      </button>
                    )}
                  </div>
                </motion.div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}

export default DashboardPage
