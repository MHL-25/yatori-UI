import React, { useMemo, useState } from 'react'
import { useAppStore } from '../stores/useAppStore'
import { CourseCard } from './CheckPage'
import { cn, getPlatformName, getPlatformIcon } from '../utils/helpers'
import { ArrowLeft, LogOut, RefreshCw, Search, ShieldCheck, Loader2 } from 'lucide-react'
import { motion } from 'framer-motion'

const CheckDetailPage: React.FC = () => {
  const {
    accounts, checkDetailUid, checkReports, checkingUids,
    checkAccount, closeCheckDetail, setCurrentPage, fetchCheckReports,
  } = useAppStore()
  const [keyword, setKeyword] = useState('')

  const uid = checkDetailUid || ''
  const report = uid ? checkReports[uid] : undefined
  const account = accounts.find((a) => a.uid === uid)

  const filteredCourses = useMemo(() => {
    if (!report) return []
    const kw = keyword.trim().toLowerCase()
    if (!kw) return report.courses
    return report.courses.filter((c) => c.courseName.toLowerCase().includes(kw))
  }, [report, keyword])

  const stats = report ? [
    { label: '课程总数', value: report.totalCourses, tone: 'text-dark-100' },
    { label: '已完成', value: report.completedCourses, tone: 'text-accent-400' },
    { label: '待刷', value: report.pendingCourses, tone: 'text-dark-200' },
    { label: '受限', value: report.lockedCourses, tone: 'text-warn' },
  ] : []

  const handleRecheck = async () => {
    if (!uid) return
    const r = await checkAccount(uid)
    if (!r.ok) alert(r.msg)
  }

  return (
    <div className="h-full overflow-y-auto">
      <div className="sticky top-0 z-10 bg-dark-950/95 backdrop-blur border-b border-dark-700 px-6 py-4">
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 min-w-0">
            <button onClick={closeCheckDetail} className="btn-secondary flex items-center gap-1.5 shrink-0">
              <ArrowLeft size={15} /> 返回
            </button>
            <div className="min-w-0">
              <h1 className="text-lg font-semibold text-dark-50 flex items-center gap-2 truncate">
                <ShieldCheck size={18} className="text-accent-400 shrink-0" />
                {report?.accountName || account?.remarkName || '课程质检详情'}
              </h1>
              <p className="text-xs text-dark-400 mt-0.5 truncate">
                {account ? `${getPlatformIcon(account.accountType)} ${getPlatformName(account.accountType)} · ${account.account}` : ''}
                {report ? ` · 检测于 ${new Date(report.checkedAt).toLocaleString('zh-CN')}` : ''}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <button onClick={handleRecheck} disabled={!uid || checkingUids.has(uid)} className="btn-secondary flex items-center gap-1.5">
              {checkingUids.has(uid) ? <Loader2 size={14} className="animate-spin" /> : <RefreshCw size={14} />}
              重新检测
            </button>
            <button onClick={() => setCurrentPage('dashboard')} className="btn-ghost flex items-center gap-1.5">
              <LogOut size={14} /> 退出
            </button>
          </div>
        </div>
      </div>

      <div className="p-6 space-y-4">
        {!report && (
          <div className="surface-card p-14 text-center">
            <ShieldCheck size={30} className="mx-auto text-dark-600 mb-3" />
            <p className="text-sm text-dark-300">该账号还没有质检结果</p>
            <button onClick={handleRecheck} className="btn-primary mt-4">立即检测</button>
          </div>
        )}

        {report && (
          <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.24 }} className="space-y-4">
            <div className="grid grid-cols-4 gap-3">
              {stats.map((s) => (
                <div key={s.label} className="surface-card px-4 py-3">
                  <div className={cn('text-2xl font-semibold font-mono', s.tone)}>{s.value}</div>
                  <div className="text-xs text-dark-400 mt-0.5">{s.label}</div>
                </div>
              ))}
            </div>

            <div className="flex items-center justify-between gap-3">
              <div className="relative flex-1 max-w-sm">
                <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-dark-400" />
                <input
                  value={keyword}
                  onChange={(e) => setKeyword(e.target.value)}
                  placeholder="搜索课程名称"
                  className="input-field pl-9"
                />
              </div>
              <span className="text-xs text-dark-400">显示 {filteredCourses.length} / {report.courses.length} 门</span>
            </div>

            <div className="space-y-2.5">
              {filteredCourses.map((c) => <CourseCard key={c.courseId + c.courseName} course={c} />)}
              {filteredCourses.length === 0 && (
                <div className="surface-card p-10 text-center text-sm text-dark-400">没有匹配的课程</div>
              )}
            </div>
          </motion.div>
        )}
      </div>
    </div>
  )
}

export default CheckDetailPage
