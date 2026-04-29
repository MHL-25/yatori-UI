import React, { useEffect, useState } from 'react'
import { useAppStore } from '../stores/useAppStore'
import { cn, AI_TYPE_LIST, APP_VERSION } from '../utils/helpers'
import { Settings, Save, Brain, Mail, Server, Info, ExternalLink, Eye, FileText } from 'lucide-react'
import { motion } from 'framer-motion'

const SettingsPage: React.FC = () => {
  const { config, fetchConfig, saveConfig } = useAppStore()
  const [activeTab, setActiveTab] = useState('ai')
  const [aiType, setAiType] = useState('TONGYI')
  const [aiUrl, setAiUrl] = useState('')
  const [aiModel, setAiModel] = useState('')
  const [aiApiKey, setAiApiKey] = useState('')
  const [visionAiType, setVisionAiType] = useState('TONGYI')
  const [visionAiUrl, setVisionAiUrl] = useState('')
  const [visionAiModel, setVisionAiModel] = useState('')
  const [visionAiApiKey, setVisionAiApiKey] = useState('')
  const [apiQueUrl, setApiQueUrl] = useState('http://localhost:8083')
  const [completionTone, setCompletionTone] = useState(1)
  const [logLevel, setLogLevel] = useState('INFO')
  const [emailSw, setEmailSw] = useState(0)
  const [smtpHost, setSmtpHost] = useState('')
  const [smtpPort, setSmtpPort] = useState(465)
  const [emailUser, setEmailUser] = useState('')
  const [emailPass, setEmailPass] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => { fetchConfig() }, [])

  useEffect(() => {
    if (config) {
      setAiType(config.setting?.aiSetting?.aiType || 'TONGYI')
      setAiUrl(config.setting?.aiSetting?.aiUrl || '')
      setAiModel(config.setting?.aiSetting?.model || '')
      setAiApiKey(config.setting?.aiSetting?.API_KEY || '')
      setVisionAiType(config.setting?.visionAiSetting?.aiType || 'TONGYI')
      setVisionAiUrl(config.setting?.visionAiSetting?.aiUrl || '')
      setVisionAiModel(config.setting?.visionAiSetting?.model || '')
      setVisionAiApiKey(config.setting?.visionAiSetting?.API_KEY || '')
      setApiQueUrl(config.setting?.apiQueSetting?.url || 'http://localhost:8083')
      setCompletionTone(config.setting?.basicSetting?.completionTone ?? 1)
      setLogLevel(config.setting?.basicSetting?.logLevel || 'INFO')
      setEmailSw(config.setting?.emailInform?.sw ?? 0)
      setSmtpHost(config.setting?.emailInform?.smtpHost || '')
      setSmtpPort(config.setting?.emailInform?.smtpPort || 465)
      setEmailUser(config.setting?.emailInform?.userName || '')
      setEmailPass(config.setting?.emailInform?.password || '')
    }
  }, [config])

  const handleSave = async () => {
    setSaving(true)
    try {
      const cfg = {
        setting: {
          basicSetting: { completionTone, logLevel, colorLog: 1, logOutFileSw: 1, logModel: 0 },
          aiSetting: { aiType, aiUrl, model: aiModel, API_KEY: aiApiKey },
          visionAiSetting: { aiType: visionAiType, aiUrl: visionAiUrl, model: visionAiModel, API_KEY: visionAiApiKey },
          apiQueSetting: { url: apiQueUrl },
          emailInform: { sw: emailSw, smtpHost, smtpPort, userName: emailUser, password: emailPass },
        },
        users: config?.users || [],
      }
      await saveConfig(cfg)
    } catch (e) { console.error('保存配置失败:', e) }
    setSaving(false)
  }

  const tabs = [
    { id: 'ai', label: 'AI答题配置', icon: Brain },
    { id: 'basic', label: '基本设置', icon: Settings },
    { id: 'email', label: '邮件通知', icon: Mail },
    { id: 'api', label: 'API配置', icon: Server },
    { id: 'about', label: '关于软件', icon: Info },
  ]

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-dark-50">系统设置</h1>
          <p className="text-sm text-dark-400 mt-1">AI答题 / 通知 / 日志</p>
        </div>
        <button onClick={handleSave} disabled={saving} className={cn("btn-primary flex items-center gap-2", saving && "opacity-50")}>
          <Save size={14} />{saving ? '保存中...' : '保存配置'}
        </button>
      </div>

      <div className="flex gap-6">
        <div className="w-48 space-y-1">
          {tabs.map(tab => (
            <button key={tab.id} onClick={() => setActiveTab(tab.id)} className={cn("sidebar-item w-full", activeTab === tab.id && "sidebar-item-active")}>
              <tab.icon size={16} />
              <span>{tab.label}</span>
            </button>
          ))}
        </div>

        <div className="flex-1">
          <motion.div key={activeTab} initial={{ opacity:0, x:10 }} animate={{ opacity:1, x:0 }} className="glass-card p-6 space-y-5">
            {activeTab === 'ai' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-semibold text-dark-100 flex items-center gap-2">
                    <FileText size={18} className="text-accent-400" />
                    AI答题配置
                  </h3>
                  <p className="text-xs text-dark-500 mt-1">配置两个AI模型：纯文本模型用于普通题目，识图模型用于含图片的题目</p>
                </div>

                <div className="p-4 rounded-lg bg-dark-800/30 border border-dark-700/50 space-y-4">
                  <div className="flex items-center gap-2 mb-2">
                    <FileText size={16} className="text-blue-400" />
                    <h4 className="text-sm font-semibold text-dark-200">纯文本模型（必填）</h4>
                  </div>
                  <p className="text-xs text-dark-500 -mt-2 ml-6">用于处理纯文字题目，如选择题、判断题等</p>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">AI平台类型</label>
                    <select value={aiType} onChange={e => setAiType(e.target.value)} className="select-field">
                      {AI_TYPE_LIST.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
                    </select>
                  </div>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">AI API地址</label>
                    <input type="text" value={aiUrl} onChange={e => setAiUrl(e.target.value)} placeholder="自定义API地址（选填）" className="input-field" />
                    <p className="text-xs text-dark-500 mt-1">仅在使用"其他"或自定义URL时填写</p>
                  </div>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">模型名称</label>
                    <input type="text" value={aiModel} onChange={e => setAiModel(e.target.value)} placeholder="如：qwen-turbo, gpt-3.5-turbo" className="input-field" />
                  </div>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">API密钥</label>
                    <input type="password" value={aiApiKey} onChange={e => setAiApiKey(e.target.value)} placeholder="请输入API KEY" className="input-field" />
                    <p className="text-xs text-dark-500 mt-1">密钥将保存在本地配置文件中</p>
                  </div>
                </div>

                <div className="p-4 rounded-lg bg-dark-800/30 border border-dark-700/50 space-y-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Eye size={16} className="text-purple-400" />
                    <h4 className="text-sm font-semibold text-dark-200">识图模型（选填）</h4>
                  </div>
                  <p className="text-xs text-dark-500 -mt-2 ml-6">用于处理含图片的题目，如数学公式图、图表题等。若未配置，含图片的题目将降级使用纯文本模型</p>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">AI平台类型</label>
                    <select value={visionAiType} onChange={e => setVisionAiType(e.target.value)} className="select-field">
                      {AI_TYPE_LIST.filter(t => t.id !== 'METAAI').map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
                    </select>
                    <p className="text-xs text-dark-500 mt-1">推荐使用支持视觉能力的模型，如通义千问(qwen-vl)、GPT-4o等</p>
                  </div>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">AI API地址</label>
                    <input type="text" value={visionAiUrl} onChange={e => setVisionAiUrl(e.target.value)} placeholder="自定义API地址（选填）" className="input-field" />
                  </div>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">模型名称</label>
                    <input type="text" value={visionAiModel} onChange={e => setVisionAiModel(e.target.value)} placeholder="如：qwen-vl-plus, gpt-4o" className="input-field" />
                    <p className="text-xs text-dark-500 mt-1">请填写支持图片识别的视觉模型名称</p>
                  </div>
                  <div>
                    <label className="text-xs font-medium text-dark-300 mb-1.5 block">API密钥</label>
                    <input type="password" value={visionAiApiKey} onChange={e => setVisionAiApiKey(e.target.value)} placeholder="请输入API KEY（选填）" className="input-field" />
                    <p className="text-xs text-dark-500 mt-1">可与纯文本模型使用同一密钥（若平台相同）</p>
                  </div>
                </div>

                <div className="p-3 rounded-lg bg-accent-900/20 border border-accent-700/30">
                  <p className="text-xs text-accent-300">
                    💡 提示：识图模型为可选配置。若不配置识图模型，遇到含图片的题目时将自动使用纯文本模型作答（图片内容将被忽略）。配置识图模型后可显著提升含图题目的准确率。
                  </p>
                </div>
              </div>
            )}

            {activeTab === 'basic' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-dark-100">基本设置</h3>
                <div className="flex items-center justify-between p-3 rounded-lg bg-dark-800/30">
                  <div>
                    <p className="text-sm font-medium text-dark-200">完成提示音</p>
                    <p className="text-xs text-dark-500">任务完成时播放提示音</p>
                  </div>
                  <button onClick={() => setCompletionTone(completionTone === 1 ? 0 : 1)} className={cn("w-10 h-5 rounded-full transition-colors relative", completionTone === 1 ? "bg-accent-600" : "bg-dark-600")}>
                    <div className={cn("w-4 h-4 rounded-full bg-white absolute top-0.5 transition-transform", completionTone === 1 ? "translate-x-5" : "translate-x-0.5")} />
                  </button>
                </div>
                <div>
                  <label className="text-xs font-medium text-dark-300 mb-1.5 block">日志级别</label>
                  <select value={logLevel} onChange={e => setLogLevel(e.target.value)} className="select-field">
                    <option value="DEBUG">调试 (DEBUG)</option>
                    <option value="INFO">信息 (INFO)</option>
                    <option value="WARN">警告 (WARN)</option>
                    <option value="ERROR">错误 (ERROR)</option>
                  </select>
                </div>
              </div>
            )}

            {activeTab === 'email' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-dark-100">邮件通知</h3>
                <div className="flex items-center justify-between p-3 rounded-lg bg-dark-800/30">
                  <div>
                    <p className="text-sm font-medium text-dark-200">启用邮件通知</p>
                    <p className="text-xs text-dark-500">任务完成后发送邮件提醒</p>
                  </div>
                  <button onClick={() => setEmailSw(emailSw === 1 ? 0 : 1)} className={cn("w-10 h-5 rounded-full transition-colors relative", emailSw === 1 ? "bg-accent-600" : "bg-dark-600")}>
                    <div className={cn("w-4 h-4 rounded-full bg-white absolute top-0.5 transition-transform", emailSw === 1 ? "translate-x-5" : "translate-x-0.5")} />
                  </button>
                </div>
                {emailSw === 1 && (
                  <>
                    <div>
                      <label className="text-xs font-medium text-dark-300 mb-1.5 block">SMTP服务器</label>
                      <input type="text" value={smtpHost} onChange={e => setSmtpHost(e.target.value)} placeholder="如：smtp.qq.com" className="input-field" />
                    </div>
                    <div>
                      <label className="text-xs font-medium text-dark-300 mb-1.5 block">SMTP端口</label>
                      <input type="number" value={smtpPort} onChange={e => setSmtpPort(Number(e.target.value))} className="input-field" />
                    </div>
                    <div>
                      <label className="text-xs font-medium text-dark-300 mb-1.5 block">邮箱账号</label>
                      <input type="text" value={emailUser} onChange={e => setEmailUser(e.target.value)} placeholder="your@email.com" className="input-field" />
                    </div>
                    <div>
                      <label className="text-xs font-medium text-dark-300 mb-1.5 block">授权码</label>
                      <input type="password" value={emailPass} onChange={e => setEmailPass(e.target.value)} placeholder="邮箱授权码或密码" className="input-field" />
                    </div>
                  </>
                )}
              </div>
            )}

            {activeTab === 'api' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-dark-100">API配置</h3>
                <div>
                  <label className="text-xs font-medium text-dark-300 mb-1.5 block">外置题库服务URL</label>
                  <input type="text" value={apiQueUrl} onChange={e => setApiQueUrl(e.target.value)} placeholder="http://localhost:8083" className="input-field" />
                  <p className="text-xs text-dark-500 mt-1">外置题库服务的访问地址，用于自动考试功能</p>
                </div>
              </div>
            )}

            {activeTab === 'about' && (
              <div className="space-y-6">
                <div className="flex items-center gap-4">
                  <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-accent-500 to-neon-purple flex items-center justify-center shadow-lg shadow-accent-600/20">
                    <span className="text-3xl">⚡</span>
                  </div>
                  <div>
                    <h3 className="text-xl font-bold text-dark-50">Yatori-UI</h3>
                    <p className="text-sm text-dark-400">智能网课助手 · 桌面端</p>
                    <p className="text-xs text-accent-400 mt-1">版本号：{APP_VERSION}</p>
                  </div>
                </div>

                <div className="border-t border-dark-700 pt-4 space-y-4">
                  <div className="p-4 rounded-lg bg-dark-800/30 space-y-3">
                    <h4 className="text-sm font-semibold text-dark-200">核心代码来源</h4>
                    <p className="text-xs text-dark-400">
                      本软件核心代码由 <span className="text-accent-400 font-medium">Yatori-Dev</span> 提供
                    </p>
                    <a
                      href="https://github.com/yatori-dev/yatori-go-console"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 text-xs text-blue-400 hover:text-blue-300 transition-colors"
                    >
                      <ExternalLink size={12} />
                      github.com/yatori-dev/yatori-go-console
                    </a>
                  </div>

                  <div className="p-4 rounded-lg bg-dark-800/30 space-y-3">
                    <h4 className="text-sm font-semibold text-dark-200">作者</h4>
                    <p className="text-sm text-neon-purple font-medium">❦Angelic 音乐</p>
                    <p className="text-xs text-dark-400">联系方式：QQ 2844189228</p>
                  </div>

                  <div className="p-4 rounded-lg bg-dark-800/30 space-y-3">
                    <h4 className="text-sm font-semibold text-dark-200">项目地址</h4>
                    <a
                      href="https://github.com/MHL-25/yatori-UI"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 text-xs text-blue-400 hover:text-blue-300 transition-colors"
                    >
                      <ExternalLink size={12} />
                      github.com/MHL-25/yatori-UI
                    </a>
                  </div>
                </div>

                <div className="border-t border-dark-700 pt-4">
                  <p className="text-xs text-dark-500 text-center">
                    本程序仅供学习、研究与技术交流使用，严禁用于任何商业用途或违法活动。
                  </p>
                </div>
              </div>
            )}
          </motion.div>
        </div>
      </div>
    </div>
  )
}

export default SettingsPage
