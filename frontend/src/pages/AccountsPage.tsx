import React, { useState, useEffect } from 'react'
import { useAppStore } from '../stores/useAppStore'
import { cn, getPlatformName, getPlatformIcon, PLATFORM_LIST, getStatusLabel, getStatusColor, getProgressColor } from '../utils/helpers'
import { Plus, Trash2, Edit3, CheckCircle, XCircle, Play, Square, Pause, Search, RefreshCw, Users, Settings, ChevronDown, ChevronUp, Mail, Upload, PlayCircle, CheckSquare, Square as EmptySquare, Loader2 } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

const PLATFORM_VIDEO_MODES: Record<string, { value: number; label: string; desc: string }[]> = {
  XUEXITONG: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '普通模式', desc: '标准播放' },
    { value: 2, label: '暴力模式', desc: '快速提交（可能回滚进度）' },
    { value: 3, label: '多任务点模式', desc: '并行执行多个任务点' },
  ],
  YINGHUA: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '普通模式', desc: '5秒间隔提交' },
    { value: 2, label: '暴力模式', desc: '8秒间隔快速提交' },
    { value: 3, label: '去红模式', desc: '处理并行播放标红' },
  ],
  CQIE: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '常规模式', desc: '3秒步进提交' },
    { value: 2, label: '秒刷模式', desc: '一次性提交完成' },
  ],
  ENAEA: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '常规模式', desc: '25秒间隔提交' },
    { value: 2, label: '暴力模式', desc: '60秒间隔快速提交' },
  ],
  ICVE: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '秒刷模式', desc: '直接完成' },
  ],
  QSXT: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '累计学时', desc: '60秒间隔提交' },
  ],
  WELEARN: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '累计学时', desc: '60秒间隔提交' },
    { value: 2, label: '刷完成度', desc: '直接完成' },
  ],
  KETANGX: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '秒刷模式', desc: '直接完成' },
  ],
  HQKJ: [
    { value: 0, label: '不刷', desc: '跳过视频' },
    { value: 1, label: '普通模式', desc: '30秒间隔提交' },
    { value: 2, label: '快速模式', desc: '并发秒刷' },
  ],
}

const getVideoModes = (platform: string) => PLATFORM_VIDEO_MODES[platform] || [
  { value: 0, label: '不刷', desc: '跳过视频' },
  { value: 1, label: '普通模式', desc: '标准播放' },
]

const EXAM_MODES = [
  { value: 0, label: '不考' },
  { value: 1, label: 'AI大模型自动答题' },
  { value: 2, label: '外置题库答题' },
  { value: 3, label: '内置AI答题' },
]

const SUBMIT_MODES = [
  { value: 0, label: '仅保存' },
  { value: 1, label: '自动提交' },
  { value: 2, label: '仅自动提交最后一题' },
]

interface AccountFormData {
  accountType: string
  url: string
  remarkName: string
  account: string
  password: string
  isProxy: boolean
  informEmails: string[]
  videoModel: number
  autoExam: number
  examAutoSubmit: number
  cxNode: number
  cxChapterTestSw: number
  cxWorkSw: number
  cxExamSw: number
  shuffleSw: number
  studyTime: string
  includeCourses: string[]
  excludeCourses: string[]
}

const defaultForm: AccountFormData = {
  accountType: 'XUEXITONG',
  url: '',
  remarkName: '',
  account: '',
  password: '',
  isProxy: false,
  informEmails: [],
  videoModel: 1,
  autoExam: 1,
  examAutoSubmit: 1,
  cxNode: 3,
  cxChapterTestSw: 1,
  cxWorkSw: 1,
  cxExamSw: 1,
  shuffleSw: 0,
  studyTime: '',
  includeCourses: [],
  excludeCourses: [],
}

const AccountFormFields: React.FC<{
  form: AccountFormData
  setForm: React.Dispatch<React.SetStateAction<AccountFormData>>
  showAdvanced: boolean
  setShowAdvanced: (v: boolean) => void
  isEdit?: boolean
}> = ({ form, setForm, showAdvanced, setShowAdvanced, isEdit }) => {
  const needUrl = ['YINGHUA', 'CANGHUI', 'HQKJ'].includes(form.accountType)
  const needStudyTime = form.accountType === 'WELEARN'

  return (
    <>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="text-xs font-medium text-dark-300 mb-1 block">平台类型</label>
          <select value={form.accountType} onChange={e => {
            const newType = e.target.value
            const modes = getVideoModes(newType)
            const validMode = modes.find(m => m.value === form.videoModel) ? form.videoModel : modes[Math.min(1, modes.length - 1)].value
            setForm({ ...form, accountType: newType, videoModel: validMode })
          }} className="select-field" disabled={isEdit}>
            {PLATFORM_LIST.map(p => <option key={p.id} value={p.id}>{p.icon} {p.name}</option>)}
          </select>
        </div>
        <div>
          <label className="text-xs font-medium text-dark-300 mb-1 block">备注名</label>
          <input type="text" value={form.remarkName} onChange={e => setForm({ ...form, remarkName: e.target.value })}
            placeholder="可选，方便识别" className="input-field" />
        </div>
      </div>

      {needUrl && (
        <div><label className="text-xs font-medium text-dark-300 mb-1 block">平台URL</label>
          <input type="text" value={form.url} onChange={e => setForm({ ...form, url: e.target.value })} placeholder="https://example.com" className="input-field" />
        </div>
      )}

      <div className="grid grid-cols-2 gap-3">
        <div><label className="text-xs font-medium text-dark-300 mb-1 block">账号</label>
          <input type="text" value={form.account} onChange={e => setForm({ ...form, account: e.target.value })} placeholder="请输入账号" className="input-field" />
        </div>
        <div><label className="text-xs font-medium text-dark-300 mb-1 block">密码 / Cookie</label>
          <input type="password" value={form.password} onChange={e => setForm({ ...form, password: e.target.value })} placeholder="请输入密码" className="input-field" />
        </div>
      </div>

      <div className="flex items-center justify-between p-2.5 rounded-lg bg-dark-800/30">
        <span className="text-sm text-dark-200">是否开启代理IP</span>
        <button onClick={() => setForm({ ...form, isProxy: !form.isProxy })}
          className={cn("w-10 h-5 rounded-full transition-colors relative", form.isProxy ? "bg-accent-600" : "bg-dark-600")}>
          <div className={cn("w-4 h-4 rounded-full bg-white absolute top-0.5 transition-transform", form.isProxy ? "translate-x-5" : "translate-x-0.5")} />
        </button>
      </div>

      <div>
        <label className="text-xs font-medium text-dark-300 mb-1 block">通知邮箱</label>
        <div className="flex gap-2">
          <input type="email" placeholder="email@example.com" className="input-field flex-1"
            onKeyDown={e => { if (e.key === 'Enter') { const target = e.target as HTMLInputElement; if (target.value && !form.informEmails.includes(target.value)) setForm({ ...form, informEmails: [...form.informEmails, target.value] }); target.value = '' } }} />
          <button onClick={() => {}}
            className="btn-secondary text-xs whitespace-nowrap"><Plus size={14} /> 新增邮箱</button>
        </div>
        {form.informEmails.length > 0 && (
          <div className="flex flex-wrap gap-1 mt-1.5">
            {form.informEmails.map((email, i) => (
              <span key={i} className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-accent-600/10 text-accent-400 text-xs">
                <Mail size={10} />{email}
                <button onClick={() => setForm({ ...form, informEmails: form.informEmails.filter((_, idx) => idx !== i) })}><XCircle size={12} /></button>
              </span>
            ))}
          </div>
        )}
      </div>

      <button onClick={() => setShowAdvanced(!showAdvanced)}
        className="flex items-center gap-2 w-full py-2 text-sm text-dark-400 hover:text-dark-200 transition-colors">
        <Settings size={14} />
        高级设置 ({showAdvanced ? '收起' : '展开'})
        <ChevronDown size={14} className={cn("transition-transform", showAdvanced && "rotate-180")} />
      </button>

      {showAdvanced && (
        <motion.div initial={{ height: 0, opacity: 0 }} animate={{ height: 'auto', opacity: 1 }} className="space-y-3 border-t border-dark-700 pt-3">
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-xs font-medium text-dark-300 mb-1 block">视频模式</label>
              <select value={form.videoModel} onChange={e => setForm({ ...form, videoModel: Number(e.target.value) })} className="select-field">
                {getVideoModes(form.accountType).map(m => <option key={m.value} value={m.value}>{m.label} - {m.desc}</option>)}
              </select>
              {form.videoModel === 2 && form.accountType !== 'CQIE' && form.accountType !== 'WELEARN' && form.accountType !== 'KETANGX' && form.accountType !== 'HQKJ' && <p className="text-xs text-orange-400 mt-0.5">注意: 暴力模式有概率打回进度</p>}
              {form.videoModel === 3 && form.accountType === 'XUEXITONG' && <div className="mt-1.5">
                <label className="text-xs text-dark-400 block mb-0.5">同时任务点数量</label>
                <input type="number" min={1} max={10} value={form.cxNode}
                  onChange={e => setForm({ ...form, cxNode: Number(e.target.value) })} className="input-field" />
              </div>}
            </div>
            <div>
              <label className="text-xs font-medium text-dark-300 mb-1 block">自动考试模式</label>
              <select value={form.autoExam} onChange={e => setForm({ ...form, autoExam: Number(e.target.value) })} className="select-field">
                {EXAM_MODES.map(m => <option key={m.value} value={m.value}>{m.label}</option>)}
              </select>
            </div>
          </div>

          {needStudyTime && (
            <div><label className="text-xs font-medium text-dark-300 mb-1 block">WeLearn学习时间范围</label>
              <input type="text" value={form.studyTime || ''} onChange={e => setForm({ ...form, studyTime: e.target.value })} placeholder="如: 08:00-22:00" className="input-field" />
            </div>
          )}

          {form.accountType === 'XUEXITONG' && (
            <div className="grid grid-cols-3 gap-3">
              <div><label className="text-xs font-medium text-dark-300 mb-1 block">章节测试</label>
                <select value={form.cxChapterTestSw} onChange={e => setForm({ ...form, cxChapterTestSw: Number(e.target.value) })} className="select-field">
                  <option value={0}>关闭</option><option value={1}>开启</option>
                </select></div>
              <div><label className="text-xs font-medium text-dark-300 mb-1 block">作业</label>
                <select value={form.cxWorkSw} onChange={e => setForm({ ...form, cxWorkSw: Number(e.target.value) })} className="select-field">
                  <option value={0}>关闭</option><option value={1}>开启</option>
                </select></div>
              <div><label className="text-xs font-medium text-dark-300 mb-1 block">考试</label>
                <select value={form.cxExamSw} onChange={e => setForm({ ...form, cxExamSw: Number(e.target.value) })} className="select-field">
                  <option value={0}>关闭</option><option value={1}>开启</option>
                </select></div>
            </div>
          )}

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-xs font-medium text-dark-300 mb-1 block">试卷提交方式</label>
              <select value={form.examAutoSubmit} onChange={e => setForm({ ...form, examAutoSubmit: Number(e.target.value) })} className="select-field">
                {SUBMIT_MODES.map(m => <option key={m.value} value={m.value}>{m.label}</option>)}
              </select>
            </div>
            <div>
              <label className="text-xs font-medium text-dark-300 mb-1 block">打乱顺序</label>
              <select value={form.shuffleSw} onChange={e => setForm({ ...form, shuffleSw: Number(e.target.value) })} className="select-field">
                <option value={0}>否</option><option value={1}>是</option>
              </select>
            </div>
          </div>

          <div>
            <label className="text-xs font-medium text-dark-300 mb-1 block">只刷课程设定</label>
            <div className="flex gap-2">
              <input type="text" placeholder="输入课程名后回车添加" className="input-field flex-1"
                onKeyDown={e => { if (e.key === 'Enter') { const target = e.target as HTMLInputElement; const v = target.value.trim(); if (v && !form.includeCourses.includes(v)) setForm({ ...form, includeCourses: [...form.includeCourses, v] }); target.value = '' } }} />
            </div>
            {form.includeCourses.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-1.5">
                {form.includeCourses.map((c, i) => (
                  <span key={i} className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-emerald-600/10 text-emerald-400 text-xs">
                    {c}
                    <button onClick={() => setForm({ ...form, includeCourses: form.includeCourses.filter((_, idx) => idx !== i) })}><XCircle size={12} /></button>
                  </span>
                ))}
              </div>
            )}
          </div>
          <div>
            <label className="text-xs font-medium text-dark-300 mb-1 block">排除课程设定</label>
            <div className="flex gap-2">
              <input type="text" placeholder="输入课程名后回车添加" className="input-field flex-1"
                onKeyDown={e => { if (e.key === 'Enter') { const target = e.target as HTMLInputElement; const v = target.value.trim(); if (v && !form.excludeCourses.includes(v)) setForm({ ...form, excludeCourses: [...form.excludeCourses, v] }); target.value = '' } }} />
            </div>
            {form.excludeCourses.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-1.5">
                {form.excludeCourses.map((c, i) => (
                  <span key={i} className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-red-600/10 text-red-400 text-xs">
                    {c}
                    <button onClick={() => setForm({ ...form, excludeCourses: form.excludeCourses.filter((_, idx) => idx !== i) })}><XCircle size={12} /></button>
                  </span>
                ))}
              </div>
            )}
          </div>
        </motion.div>
      )}
    </>
  )
}

const AddAccountModal: React.FC<{ open: boolean; onClose: () => void }> = ({ open, onClose }) => {
  const { addAccount } = useAppStore()
  const [form, setForm] = useState<AccountFormData>(defaultForm)
  const [loading, setLoading] = useState(false)
  const [showAdvanced, setShowAdvanced] = useState(false)

  useEffect(() => {
    if (open) setForm(defaultForm)
  }, [open])

  const handleSubmit = async () => {
    if (!form.account || !form.password) return
    setLoading(true)
    const data = {
      ...form,
      coursesCustom: {
        videoModel: form.videoModel, autoExam: form.autoExam, examAutoSubmit: form.examAutoSubmit,
        cxNode: form.cxNode, cxChapterTestSw: form.cxChapterTestSw, cxWorkSw: form.cxWorkSw, cxExamSw: form.cxExamSw,
        shuffleSw: form.shuffleSw, studyTime: form.studyTime,
        includeCourses: form.includeCourses.filter(e => e),
        excludeCourses: form.excludeCourses.filter(e => e),
      },
      isProxy: form.isProxy ? 1 : 0,
      informEmails: form.informEmails.filter(e => e),
    }
    const result = await addAccount(data)
    setLoading(false)
    if (result.ok) { onClose() } else { alert('添加失败: ' + result.msg) }
  }

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }}
        className="glass-card w-[560px] max-h-[90vh] overflow-y-auto p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-dark-50">添加账号</h3>
          <button onClick={onClose} className="text-dark-400 hover:text-dark-200"><XCircle size={20} /></button>
        </div>
        <AccountFormFields form={form} setForm={setForm} showAdvanced={showAdvanced} setShowAdvanced={setShowAdvanced} />
        <div className="flex justify-end gap-3 pt-2">
          <button onClick={onClose} className="btn-secondary">取消</button>
          <button onClick={handleSubmit} disabled={loading || !form.account || !form.password}
            className={cn("btn-primary", (loading || !form.account || !form.password) && "opacity-50 cursor-not-allowed")}>
            {loading ? '添加中...' : '确认添加'}
          </button>
        </div>
      </motion.div>
    </div>
  )
}

interface EditAccount {
  uid: string
  accountType: string
  url: string
  remarkName: string
  account: string
  password: string
  userConfigJson: string
}

const EditAccountModal: React.FC<{ open: boolean; onClose: () => void; account: EditAccount | null }> = ({ open, onClose, account }) => {
  const { updateAccount } = useAppStore()
  const [form, setForm] = useState<AccountFormData>(defaultForm)
  const [loading, setLoading] = useState(false)
  const [showAdvanced, setShowAdvanced] = useState(true)

  useEffect(() => {
    if (open && account) {
      let parsed: any = {}
      try { parsed = JSON.parse(account.userConfigJson || '{}') } catch {}
      const cc = parsed.coursesCustom || {}
      setForm({
        accountType: account.accountType || parsed.accountType || 'XUEXITONG',
        url: account.url || parsed.url || '',
        remarkName: account.remarkName || parsed.remarkName || '',
        account: account.account || parsed.account || '',
        password: account.password || parsed.password || '',
        isProxy: (parsed.isProxy || 0) === 1,
        informEmails: parsed.informEmails || [],
        videoModel: cc.videoModel ?? 1,
        autoExam: cc.autoExam ?? 1,
        examAutoSubmit: cc.examAutoSubmit ?? 1,
        cxNode: cc.cxNode ?? 3,
        cxChapterTestSw: cc.cxChapterTestSw ?? 1,
        cxWorkSw: cc.cxWorkSw ?? 1,
        cxExamSw: cc.cxExamSw ?? 1,
        shuffleSw: cc.shuffleSw ?? 0,
        studyTime: cc.studyTime || '',
        includeCourses: Array.isArray(cc.includeCourses) ? cc.includeCourses : [],
        excludeCourses: Array.isArray(cc.excludeCourses) ? cc.excludeCourses : [],
      })
    }
  }, [open, account])

  const handleSubmit = async () => {
    if (!account) return
    setLoading(true)
    const data = {
      uid: account.uid,
      accountType: form.accountType,
      url: form.url,
      remarkName: form.remarkName,
      account: form.account,
      password: form.password,
      coursesCustom: {
        videoModel: form.videoModel, autoExam: form.autoExam, examAutoSubmit: form.examAutoSubmit,
        cxNode: form.cxNode, cxChapterTestSw: form.cxChapterTestSw, cxWorkSw: form.cxWorkSw, cxExamSw: form.cxExamSw,
        shuffleSw: form.shuffleSw, studyTime: form.studyTime,
        includeCourses: form.includeCourses.filter(e => e),
        excludeCourses: form.excludeCourses.filter(e => e),
      },
      isProxy: form.isProxy ? 1 : 0,
      informEmails: form.informEmails.filter(e => e),
    }
    const result = await updateAccount(data)
    setLoading(false)
    if (result.ok) { onClose() } else { alert('保存失败: ' + result.msg) }
  }

  if (!open || !account) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }}
        className="glass-card w-[560px] max-h-[90vh] overflow-y-auto p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-dark-50">编辑账号</h3>
          <button onClick={onClose} className="text-dark-400 hover:text-dark-200"><XCircle size={20} /></button>
        </div>
        <AccountFormFields form={form} setForm={setForm} showAdvanced={showAdvanced} setShowAdvanced={setShowAdvanced} isEdit />
        <div className="flex justify-end gap-3 pt-2">
          <button onClick={onClose} className="btn-secondary">取消</button>
          <button onClick={handleSubmit} disabled={loading}
            className={cn("btn-primary", loading && "opacity-50 cursor-not-allowed")}>
            {loading ? '保存中...' : '保存修改'}
          </button>
        </div>
      </motion.div>
    </div>
  )
}

const AccountsPage: React.FC = () => {
  const { accounts, fetchAccounts, deleteAccount, loginCheck, startBrush, stopBrush, pauseBrush, startAllBrush, startBatchBrush, progressMap, importConfig } = useAppStore()
  const [showAdd, setShowAdd] = useState(false)
  const [editAccount, setEditAccount] = useState<EditAccount | null>(null)
  const [search, setSearch] = useState('')
  const [checking, setChecking] = useState<string | null>(null)
  const [actionMsg, setActionMsg] = useState<{ uid: string; msg: string; type: 'ok' | 'err' } | null>(null)
  const [importing, setImporting] = useState(false)
  const [batchMode, setBatchMode] = useState(false)
  const [selectedUids, setSelectedUids] = useState<Set<string>>(new Set())
  const [batchStarting, setBatchStarting] = useState(false)
  const [batchProgress, setBatchProgress] = useState<{ current: number; total: number; currentName: string } | null>(null)

  useEffect(() => {
    fetchAccounts()
    const interval = setInterval(fetchAccounts, 3000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (actionMsg) {
      const timer = setTimeout(() => setActionMsg(null), 2500)
      return () => clearTimeout(timer)
    }
  }, [actionMsg])

  const filtered = accounts.filter(a =>
    a.account.includes(search) ||
    (a.remarkName && a.remarkName.includes(search)) ||
    getPlatformName(a.accountType).includes(search)
  )

  const handleLoginCheck = async (uid: string) => {
    setChecking(uid)
    const result = await loginCheck(uid)
    setChecking(null)
    setActionMsg({ uid, msg: result.ok ? '验证成功!' : result.msg, type: result.ok ? 'ok' : 'err' })
  }

  const handleDelete = async (uid: string) => {
    const result = await deleteAccount(uid)
    setActionMsg({ uid, msg: result.ok ? '删除成功!' : result.msg, type: result.ok ? 'ok' : 'err' })
  }

  const handleStart = async (uid: string) => {
    const result = await startBrush(uid)
    if (!result.ok) setActionMsg({ uid, msg: result.msg, type: 'err' })
  }

  const handleStop = async (uid: string) => {
    const result = await stopBrush(uid)
    if (!result.ok) setActionMsg({ uid, msg: result.msg, type: 'err' })
  }

  const handleImport = async () => {
    setImporting(true)
    const result = await importConfig()
    setImporting(false)
    if (result.ok) { alert(result.msg) } else { alert('导入失败: ' + result.msg) }
  }

  const handleStartAll = async () => {
    const idleAccounts = accounts.filter(a => !a.isRunning)
    if (idleAccounts.length === 0) { alert('没有可启动的账号'); return }
    setBatchStarting(true)
    setBatchProgress({ current: 0, total: idleAccounts.length, currentName: '' })
    let success = 0
    for (let i = 0; i < idleAccounts.length; i++) {
      const a = idleAccounts[i]
      setBatchProgress({ current: i + 1, total: idleAccounts.length, currentName: a.remarkName || a.account })
      const result = await startBrush(a.uid)
      if (result.ok) success++
      await new Promise(r => setTimeout(r, 500))
    }
    setBatchStarting(false)
    setBatchProgress(null)
    alert(`批量启动完成: ${success}/${idleAccounts.length} 个账号成功启动`)
  }

  const handleBatchStart = async () => {
    if (selectedUids.size === 0) return
    const uids = Array.from(selectedUids)
    setBatchStarting(true)
    setBatchProgress({ current: 0, total: uids.length, currentName: '' })
    let success = 0
    for (let i = 0; i < uids.length; i++) {
      const uid = uids[i]
      const acc = accounts.find(a => a.uid === uid)
      setBatchProgress({ current: i + 1, total: uids.length, currentName: acc?.remarkName || acc?.account || uid })
      const result = await startBrush(uid)
      if (result.ok) success++
      await new Promise(r => setTimeout(r, 500))
    }
    setBatchStarting(false)
    setBatchProgress(null)
    setSelectedUids(new Set())
    setBatchMode(false)
    alert(`批量启动完成: ${success}/${uids.length} 个账号成功启动`)
  }

  const toggleSelect = (uid: string) => {
    setSelectedUids(prev => {
      const next = new Set(prev)
      if (next.has(uid)) next.delete(uid)
      else next.add(uid)
      return next
    })
  }

  const selectAll = () => {
    const allUids = filtered.filter(a => !a.isRunning).map(a => a.uid)
    setSelectedUids(new Set(allUids))
  }

  const deselectAll = () => setSelectedUids(new Set())

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-dark-50">账号管理</h1>
          <p className="text-sm text-dark-400 mt-1">管理你的网课平台账号</p>
        </div>
        <div className="flex items-center gap-3">
          <button onClick={() => setShowAdd(true)} className="btn-primary flex items-center gap-2">
            <Plus size={16} />添加账号
          </button>
          <button onClick={handleImport} disabled={importing} className={cn("btn-secondary flex items-center gap-2", importing && "opacity-50")}>
            <Upload size={16} />{importing ? '导入中...' : '导入配置'}
          </button>
          <button onClick={handleStartAll} disabled={batchStarting} className={cn("btn-success flex items-center gap-2", batchStarting && "opacity-50")}>
            <PlayCircle size={16} />{batchStarting && batchProgress ? `启动中 ${batchProgress.current}/${batchProgress.total}` : '一键启动全部'}
          </button>
          <button onClick={() => { setBatchMode(!batchMode); setSelectedUids(new Set()) }}
            className={cn("btn-secondary flex items-center gap-2", batchMode && "bg-accent-600/20 text-accent-400 border-accent-600/30")}>
            <CheckSquare size={16} />{batchMode ? '取消批量' : '批量启动'}
          </button>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <div className="relative flex-1 max-w-md">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-dark-500" />
          <input type="text" value={search} onChange={e => setSearch(e.target.value)} placeholder="搜索账号..." className="input-field pl-9" />
        </div>
        <button onClick={() => fetchAccounts()} className="btn-secondary flex items-center gap-2">
          <RefreshCw size={14} />刷新
        </button>
        {batchMode && (
          <>
            <button onClick={selectAll} className="btn-secondary text-xs flex items-center gap-1">
              <CheckSquare size={14} />全选
            </button>
            <button onClick={deselectAll} className="btn-secondary text-xs flex items-center gap-1">
              <EmptySquare size={14} />全不选
            </button>
            <button onClick={handleBatchStart} disabled={batchStarting || selectedUids.size === 0}
              className={cn("btn-success text-xs flex items-center gap-1", (batchStarting || selectedUids.size === 0) && "opacity-50")}>
              <Play size={14} />{batchStarting && batchProgress ? `${batchProgress.current}/${batchProgress.total}` : `启动选中 (${selectedUids.size})`}
            </button>
          </>
        )}
      </div>

      {batchStarting && batchProgress && (
        <div className="glass-card p-3 flex items-center gap-3">
          <Loader2 size={16} className="text-accent-400 animate-spin flex-shrink-0" />
          <div className="flex-1">
            <div className="flex items-center justify-between mb-1">
              <span className="text-xs text-dark-300">正在启动: {batchProgress.currentName}</span>
              <span className="text-xs text-dark-400">{batchProgress.current}/{batchProgress.total}</span>
            </div>
            <div className="w-full h-1.5 bg-dark-700 rounded-full overflow-hidden">
              <div className="h-full rounded-full bg-gradient-to-r from-accent-500 to-accent-400 transition-all duration-300"
                style={{ width: `${(batchProgress.current / batchProgress.total) * 100}%` }} />
            </div>
          </div>
        </div>
      )}

      {filtered.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <Users size={48} className="mx-auto text-dark-600 mb-4" />
          <p className="text-dark-400 text-sm">{accounts.length === 0 ? '暂无账号，点击上方按钮添加' : '没有匹配的搜索结果'}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-3">
          <AnimatePresence>
            {filtered.map((account, index) => {
              const progress = progressMap[account.uid]
              const isRunning = progress?.status === 'running' || account.isRunning
              const actionInfo = actionMsg?.uid === account.uid ? actionMsg : null

              return (
                <motion.div key={account.uid} initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }} transition={{ delay: index * 0.03 }}
                  className="glass-card p-4 flex items-center gap-4 relative overflow-hidden">
                  {actionInfo && (
                    <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }}
                      className={cn("absolute right-4 top-4 px-3 py-1 rounded-md text-xs font-medium z-10",
                        actionInfo.type === 'ok' ? "bg-emerald-600/20 text-emerald-400" : "bg-red-600/20 text-red-400"
                      )}>
                      {actionInfo.msg}
                    </motion.div>
                  )}

                  {batchMode && (
                    <button onClick={() => toggleSelect(account.uid)}
                      className={cn("w-5 h-5 rounded border-2 flex items-center justify-center flex-shrink-0 transition-colors",
                        selectedUids.has(account.uid)
                          ? "bg-accent-600 border-accent-600"
                          : "border-dark-600 hover:border-dark-400"
                      )}>
                      {selectedUids.has(account.uid) && <CheckCircle size={14} className="text-white" />}
                    </button>
                  )}

                  <div className="w-10 h-10 rounded-lg bg-dark-700 flex items-center justify-center text-lg flex-shrink-0">
                    {getPlatformIcon(account.accountType)}
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-dark-100 truncate">{account.remarkName || account.account}</span>
                      <span className="text-xs px-2 py-0.5 rounded-full bg-dark-700 text-dark-300 shrink-0">{getPlatformName(account.accountType)}</span>
                      {isRunning && (
                        <span className="text-xs px-2 py-0.5 rounded-full bg-accent-600/20 text-accent-400 flex items-center gap-1 shrink-0">
                          <span className="w-1.5 h-1.5 rounded-full bg-accent-400 animate-pulse" />运行中
                        </span>
                      )}
                      {progress && progress.status && !isRunning && progress.status !== 'idle' && (
                        <span className={cn("text-xs px-2 py-0.5 rounded-full shrink-0", getStatusColor(progress.status),
                          progress.status === 'completed' ? 'bg-emerald-600/10' :
                          progress.status === 'error' ? 'bg-red-600/10' :
                          progress.status === 'paused' ? 'bg-orange-600/10' : 'bg-dark-700'
                        )}>
                          {getStatusLabel(progress.status)}
                        </span>
                      )}
                    </div>
                    <p className="text-xs text-dark-500 mt-0.5 truncate">账号: {account.account}</p>
                    {progress && (progress.progress > 0 || progress.currentTask) && (
                      <div className="mt-1.5">
                        <div className="flex items-center justify-between mb-0.5">
                          <span className="text-[10px] text-dark-500 truncate max-w-[200px]">{progress.currentTask || ''}</span>
                          <span className="text-[10px] text-dark-500 shrink-0">{(progress.progress ?? 0).toFixed(0)}%</span>
                        </div>
                        <div className="w-full h-1 bg-dark-700 rounded-full overflow-hidden">
                          <div className={cn("h-full rounded-full transition-all duration-500 bg-gradient-to-r", getProgressColor(progress.progress ?? 0))}
                            style={{ width: `${Math.min(progress.progress ?? 0, 100)}%` }} />
                        </div>
                      </div>
                    )}
                  </div>

                  <div className="flex items-center gap-2 flex-shrink-0">
                    <button onClick={() => setEditAccount(account)}
                      className="p-1.5 rounded-lg hover:bg-accent-600/20 text-dark-500 hover:text-accent-400 transition-colors" title="编辑配置">
                      <Edit3 size={14} />
                    </button>
                    <button onClick={() => handleLoginCheck(account.uid)} disabled={checking === account.uid}
                      className={cn("btn-secondary text-xs px-3 py-1.5", checking === account.uid && "opacity-50")}>
                      {checking === account.uid ? '验证中...' : '验证登录'}
                    </button>
                    {isRunning ? (
                      <>
                        <button onClick={() => pauseBrush(account.uid)} className="btn-secondary text-xs px-3 py-1.5 flex items-center gap-1">
                          <Pause size={10} />暂停
                        </button>
                        <button onClick={() => handleStop(account.uid)} className="btn-danger text-xs px-3 py-1.5 flex items-center gap-1">
                          <Square size={10} />停止
                        </button>
                      </>
                    ) : (
                      <button onClick={() => handleStart(account.uid)} className="btn-success text-xs px-3 py-1.5 flex items-center gap-1">
                        <Play size={10} />启动
                      </button>
                    )}
                    <button onClick={() => handleDelete(account.uid)}
                      className="p-1.5 rounded-lg hover:bg-red-600/20 text-dark-500 hover:text-red-400 transition-colors" title="删除">
                      <Trash2 size={14} />
                    </button>
                  </div>
                </motion.div>
              )
            })}
          </AnimatePresence>
        </div>
      )}

      <AddAccountModal open={showAdd} onClose={() => setShowAdd(false)} />
      <EditAccountModal open={!!editAccount} onClose={() => setEditAccount(null)} account={editAccount} />
    </div>
  )
}

export default AccountsPage
