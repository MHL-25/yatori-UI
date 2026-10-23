import React, { useEffect, useMemo, useState } from 'react'
import { useAppStore, CheckCourse, CheckChapter } from '../stores/useAppStore'
import { cn, getPlatformName, getPlatformIcon, getStatusLabel, getStatusColor } from '../utils/helpers'
import {
  ShieldCheck, RefreshCw, ScanSearch, Clock, Lock, CheckCircle2,
  AlertTriangle, Ban, ChevronDown, ChevronUp, Loader2, Play, Users, Search, ExternalLink,
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

const LOCK_META: Record<string, { label: string; icon: React.ElementType; cls: string }> = {
  open: { label: '可学习', icon: CheckCircle2, cls: 'badge-success' },
  time: { label: '定时开放', icon: Clock, cls: 'badge-warning' },
  prereq: { label: '前置未完成', icon: Lock, cls: 'badge-warning' },
  expired: { label: '已截止', icon: Ban, cls: 'badge-danger' },
  closed: { label: '已关闭', icon: Lock, cls: 'badge-neutral' },
}

export const LockBadge: React.FC<{ type: string; unlockAt?: string }> = ({ type, unlockAt }) => {
  const meta = LOCK_META[type] || LOCK_META.closed
  const Icon = meta.icon
  return (
    <span className={meta.cls}>
      <Icon size={11} />
      {meta.label}
      {type === 'time' && unlockAt ? ` · ${unlockAt}` : ''}
    </span>
  )
}

export const ChapterRow: React.FC<{ chapter: CheckChapter }> = ({ chapter }) => {
  const total = chapter.jobTotal || 0
  const done = chapter.jobFinished || 0
  const pct = total > 0 ? Math.min(100, Math.round((done / total) * 100)) : 0
  return (
    <div className="flex items-center gap-3 py-2 px-3 rounded-lg hover:bg-dark-800/70 transition-colors">
      <span className="w-10 text-[11px] font-mono text-dark-400 shrink-0">{chapter.label}</span>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-sm text-dark-100 truncate">{chapter.name}</span>
          <LockBadge type={chapter.lockType} unlockAt={chapter.unlockAt} />
        </div>
        <div className="mt-1 flex items-center gap-2">
          <div className="h-1.5 flex-1 rounded-full bg-dark-700 overflow-hidden">
            <div
              className={cn('h-full rounded-full transition-all duration-500',
                chapter.lockType === 'open' ? 'bg-accent-500' : 'bg-dark-500')}
              style={{ width: `${pct}%` }}
            />
          </div>
          <span className="text-[11px] text-dark-400 font-mono shrink-0">{done}/{total || chapter.jobTotal || 0}</span>
        </div>
        {chapter.message && (
          <p className="text-[11px] text-dark-400 mt-0.5 truncate">{chapter.message}</p>
        )}
      </div>
    </div>
  )
}

export const CourseCard: React.FC<{ course: CheckCourse }> = ({ course }) => {
  const [open, setOpen] = useState(false)
  const pct = Math.min(100, Math.max(0, course.jobRate || 0))
  const finished = course.jobTotal > 0 && course.jobFinished >= course.jobTotal
  return (
    <div className="surface-sub overflow-hidden">
      <button onClick={() => setOpen(!open)} className="w-full flex items-center gap-3 p-3.5 text-left hover:bg-dark-700/30 transition-colors">
        <div className={cn('w-9 h-9 rounded-xl flex items-center justify-center shrink-0',
          finished ? 'bg-accent-500/15 text-accent-400' : course.canStudy ? 'bg-dark-700 text-dark-200' : 'bg-dark-700 text-dark-400')}>
          {finished ? <CheckCircle2 size={18} /> : course.canStudy ? <Play size={18} /> : <Lock size={18} />}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-dark-100 truncate">{course.courseName}</span>
            {finished && <span className="badge-success">已完成</span>}
            {!finished && course.canStudy && <span className="badge-neutral">可继续</span>}
            {!finished && !course.canStudy && <span className="badge-warning">受限</span>}
          </div>
          <div className="mt-1.5 flex items-center gap-2">
            <div className="h-1.5 flex-1 rounded-full bg-dark-700 overflow-hidden">
              <div className={cn('h-full rounded-full transition-all duration-500',
                finished ? 'bg-accent-500' : 'bg-accent-700')} style={{ width: `${pct}%` }} />
            </div>
            <span className="text-[11px] font-mono text-dark-400 shrink-0">
              {course.jobFinished}/{course.jobTotal} · {pct.toFixed(1)}%
            </span>
          </div>
          {course.reason && <p className="text-[11px] text-dark-400 mt-1 truncate">{course.reason}</p>}
        </div>
        {course.chapters?.length > 0 && (
          <span className="text-dark-400 shrink-0">{open ? <ChevronUp size={16} /> : <ChevronDown size={16} />}</span>
        )}
      </button>
      <AnimatePresence initial={false}>
        {open && course.chapters?.length > 0 && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.24, ease: [0.22, 1, 0.36, 1] }}
            className="border-t border-dark-700/70 bg-dark-950/40"
          >
            <div className="p-2 max-h-72 overflow-y-auto">
              {course.chapters.map((ch) => <ChapterRow key={ch.nodeId} chapter={ch} />)}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

const CheckPage: React.FC = () => {
  const {
    accounts, progressMap, checkReports, checkingUids,
    checkAccount, checkAccounts, checkAllAccounts, fetchCheckReports, fetchAccounts,
    openCheckDetail,
  } = useAppStore()
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [activeUid, setActiveUid] = useState<string | null>(null)
  const [batchLoading, setBatchLoading] = useState(false)
  const [keyword, setKeyword] = useState('')

  useEffect(() => {
    fetchAccounts()
    fetchCheckReports()
  }, [])

  useEffect(() => {
    if (!activeUid && accounts.length > 0) {
      const firstWithReport = accounts.find((a) => checkReports[a.uid])
      setActiveUid(firstWithReport?.uid || accounts[0].uid)
    }
  }, [accounts, checkReports, activeUid])

  const report = activeUid ? checkReports[activeUid] : undefined

  const filteredAccounts = useMemo(() => {
    const kw = keyword.trim().toLowerCase()
    if (!kw) return accounts
    return accounts.filter((a) =>
      (a.remarkName || '').toLowerCase().includes(kw) ||
      (a.account || '').toLowerCase().includes(kw) ||
      getPlatformName(a.accountType).toLowerCase().includes(kw)
    )
  }, [accounts, keyword])

  const toggleSelect = (uid: string) => {
    setSelected((prev) => {
      const next = new Set(prev)
      next.has(uid) ? next.delete(uid) : next.add(uid)
      return next
    })
  }

  const handleCheckOne = async (uid: string) => {
    setActiveUid(uid)
    const result = await checkAccount(uid)
    if (!result.ok) alert(result.msg)
  }

  const handleBatch = async () => {
    if (selected.size === 0) return
    setBatchLoading(true)
    const result = await checkAccounts(Array.from(selected))
    setBatchLoading(false)
    if (!result.ok) alert(result.msg)
  }

  const handleAll = async () => {
    setBatchLoading(true)
    const result = await checkAllAccounts()
    setBatchLoading(false)
    if (!result.ok) alert(result.msg)
  }

  const stats = useMemo(() => {
    if (!report) return null
    return [
      { label: '课程总数', value: report.totalCourses, tone: 'text-dark-100' },
      { label: '已完成', value: report.completedCourses, tone: 'text-accent-400' },
      { label: '待刷', value: report.pendingCourses, tone: 'text-dark-200' },
      { label: '受限', value: report.lockedCourses, tone: 'text-warn' },
    ]
  }, [report])

  return (
    <div className="p-6 space-y-5">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-xl font-semibold text-dark-50 flex items-center gap-2">
            <ShieldCheck size={20} className="text-accent-400" />
            课程质检
          </h1>
          <p className="text-sm text-dark-400 mt-1">检测各账号课程完成度、前置任务锁与定时开放时间</p>
        </div>
        <div className="flex items-center gap-2">
          <button onClick={() => { fetchAccounts(); fetchCheckReports() }} className="btn-secondary flex items-center gap-2">
            <RefreshCw size={15} /> 刷新
          </button>
          <button onClick={handleAll} disabled={batchLoading} className={cn('btn-primary flex items-center gap-2', batchLoading && 'opacity-60 cursor-not-allowed')}>
            {batchLoading ? <Loader2 size={15} className="animate-spin" /> : <ScanSearch size={15} />}
            全部检测
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-[360px_1fr] gap-5">
        {/* 账号列表 */}
        <div className="space-y-3">
          <div className="surface-card p-3">
            <div className="px-1 pb-2">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-medium text-dark-300">账号 ({filteredAccounts.length}/{accounts.length})</span>
                {selected.size > 0 && (
                  <button onClick={handleBatch} disabled={batchLoading}
                    className={cn('text-xs font-semibold text-accent-400 hover:text-accent-300', batchLoading && 'opacity-50')}>
                    检测所选 ({selected.size})
                  </button>
                )}
              </div>
              <div className="relative">
                <Search size={13} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-dark-400" />
                <input
                  value={keyword}
                  onChange={(e) => setKeyword(e.target.value)}
                  placeholder="搜索备注 / 账号 / 平台"
                  className="w-full pl-7 pr-2 py-1.5 rounded-lg bg-dark-950 border border-dark-600 text-xs text-dark-100 placeholder-dark-400 focus:outline-none focus:border-accent-500"
                />
              </div>
            </div>
            <div className="space-y-2 max-h-[calc(100vh-300px)] overflow-y-auto">
              {accounts.length === 0 && (
                <div className="py-8 text-center text-sm text-dark-400">暂无账号，请先添加</div>
              )}
              {accounts.length > 0 && filteredAccounts.length === 0 && (
                <div className="py-8 text-center text-sm text-dark-400">没有匹配的账号</div>
              )}
              {filteredAccounts.map((acc) => {
                const isActive = activeUid === acc.uid
                const isChecking = checkingUids.has(acc.uid)
                const reportItem = checkReports[acc.uid]
                const progress = progressMap[acc.uid]
                return (
                  <div key={acc.uid}
                    className={cn(
                      'rounded-xl border p-3 transition-all duration-200 cursor-pointer',
                      isActive ? 'bg-dark-800 border-dark-500' : 'bg-dark-950/60 border-dark-700 hover:border-dark-600'
                    )}
                    onClick={() => setActiveUid(acc.uid)}
                  >
                    <div className="flex items-start gap-2">
                      <button
                        onClick={(e) => { e.stopPropagation(); toggleSelect(acc.uid) }}
                        className={cn('w-4 h-4 rounded border mt-0.5 shrink-0 flex items-center justify-center transition-colors',
                          selected.has(acc.uid) ? 'bg-accent-500 border-accent-500 text-dark-950' : 'border-dark-500')}
                      >
                        {selected.has(acc.uid) && <CheckCircle2 size={11} />}
                      </button>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-1.5">
                          <span className="text-sm font-medium text-dark-100 truncate">
                            {acc.remarkName || acc.account}
                          </span>
                          <span className="text-[11px] text-dark-400 shrink-0">
                            {getPlatformIcon(acc.accountType)} {getPlatformName(acc.accountType)}
                          </span>
                        </div>
                        <div className="mt-1 flex items-center gap-1.5">
                          <span className={cn('text-[11px] px-1.5 py-0.5 rounded', getStatusColor(progress?.status || 'idle'))}>
                            {getStatusLabel(progress?.status || 'idle')}
                          </span>
                          {reportItem && (
                            <span className="text-[11px] text-dark-400">
                              完成 {reportItem.completedCourses}/{reportItem.totalCourses} · 受限 {reportItem.lockedCourses}
                            </span>
                          )}
                        </div>
                      </div>
                      <button
                        onClick={(e) => { e.stopPropagation(); handleCheckOne(acc.uid) }}
                        disabled={isChecking}
                        className={cn('btn-ghost px-2 py-1 text-xs flex items-center gap-1 shrink-0', isChecking && 'opacity-50')}
                      >
                        {isChecking ? <Loader2 size={12} className="animate-spin" /> : <ScanSearch size={12} />}
                        {isChecking ? '检测中' : '检测'}
                      </button>
                      <button
                        onClick={(e) => { e.stopPropagation(); openCheckDetail(acc.uid) }}
                        className="btn-ghost px-2 py-1 text-xs flex items-center gap-1 shrink-0"
                        title="打开独立详情页"
                      >
                        <ExternalLink size={12} />详情
                      </button>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
          <p className="text-[11px] text-dark-500 px-1 flex items-center gap-1">
            <AlertTriangle size={11} /> 质检会登录账号并读取课程，不会修改学习进度
          </p>
        </div>

        {/* 质检结果 */}
        <div className="space-y-4">
          {!report && (
            <div className="surface-card p-12 text-center">
              <Users size={28} className="text-dark-600 mx-auto mb-3" />
              <p className="text-sm text-dark-300">选择账号后点击「检测」查看课程质检结果</p>
              <p className="text-xs text-dark-500 mt-1">支持单个、批量或全部账号检测</p>
            </div>
          )}

          <AnimatePresence mode="wait">
            {report && (
              <motion.div key={report.uid}
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -8 }}
                transition={{ duration: 0.24 }}
                className="space-y-4"
              >
                <div className="surface-card p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h2 className="text-base font-semibold text-dark-50">{report.accountName}</h2>
                      <p className="text-xs text-dark-400 mt-0.5">
                        检测时间 {new Date(report.checkedAt).toLocaleString('zh-CN')}
                      </p>
                    </div>
                    <button onClick={() => handleCheckOne(report.uid)} disabled={checkingUids.has(report.uid)}
                      className="btn-secondary text-xs flex items-center gap-1.5">
                      <RefreshCw size={13} className={cn(checkingUids.has(report.uid) && 'animate-spin')} /> 重新检测
                    </button>
                    <button onClick={() => openCheckDetail(report.uid)} className="btn-primary text-xs flex items-center gap-1.5 ml-2">
                      <ExternalLink size={13} /> 独立详情
                    </button>
                  </div>
                  {stats && (
                    <div className="grid grid-cols-4 gap-2 mt-3">
                      {stats.map((s) => (
                        <div key={s.label} className="rounded-xl bg-dark-950/60 border border-dark-700 px-3 py-2">
                          <div className={cn('text-lg font-semibold font-mono', s.tone)}>{s.value}</div>
                          <div className="text-[11px] text-dark-400">{s.label}</div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                <div className="space-y-2.5">
                  {report.courses.map((c) => <CourseCard key={c.courseId + c.courseName} course={c} />)}
                  {report.courses.length === 0 && (
                    <div className="surface-card p-8 text-center text-sm text-dark-400">该账号暂无课程数据</div>
                  )}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  )
}

export default CheckPage
