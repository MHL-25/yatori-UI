import React, { useState, useEffect } from 'react'

interface DisclaimerModalProps {
  onAccept: () => void
  readOnly?: boolean
  onClose?: () => void
}

const DisclaimerModal: React.FC<DisclaimerModalProps> = ({ onAccept, readOnly, onClose }) => {
  const [countdown, setCountdown] = useState(readOnly ? 0 : 10)
  const [canAccept, setCanAccept] = useState(readOnly)

  useEffect(() => {
    if (readOnly) return
    if (countdown <= 0) {
      setCanAccept(true)
      return
    }
    const timer = setInterval(() => {
      setCountdown(prev => prev - 1)
    }, 1000)
    return () => clearInterval(timer)
  }, [countdown, readOnly])

  const handleDecline = async () => {
    try {
      const { QuitApp } = await import('../../wailsjs/go/main/App')
      await QuitApp()
    } catch {
      window.close()
    }
  }

  return (
    <div className="fixed inset-0 z-[9999] flex items-center justify-center bg-black/80 backdrop-blur-sm">
      <div className="relative w-[640px] max-h-[85vh] mx-4 bg-dark-900 border border-dark-700 rounded-xl shadow-2xl overflow-hidden">
        <div className="bg-gradient-to-r from-accent-600 to-neon-purple p-4">
          <h2 className="text-lg font-bold text-white text-center">⚠️ 免责声明与使用须知</h2>
        </div>

        <div className="p-6 overflow-y-auto max-h-[55vh] text-sm text-dark-200 leading-relaxed space-y-4">
          <div className="space-y-2">
            <p>
              本软件由大佬：<span className="text-accent-400 font-semibold">"Yatori-Dev"</span> 提供核心代码，原项目地址：
              <a href="https://github.com/yatori-dev/yatori-go-console" target="_blank" rel="noopener noreferrer" className="text-blue-400 underline hover:text-blue-300">
                https://github.com/yatori-dev/yatori-go-console
              </a>
            </p>
            <p>
              本软件作者为：<span className="text-neon-purple font-semibold">"❦Angelic 音乐"</span> 联系方式：<span className="text-accent-400">"2844189228"</span>，项目地址：
              <a href="https://github.com/MHL-25/yatori-UI" target="_blank" rel="noopener noreferrer" className="text-blue-400 underline hover:text-blue-300">
                https://github.com/MHL-25/yatori-UI
              </a>
              如有问题欢迎联系我会进行修改。
            </p>
            <p>
              群聊：<a href="https://qun.qq.com/universal-share/share?ac=1&authKey=Wwi3jpZXkgm6wkjnxzOxAHp94fekV16hiwB513WmVeWntdyYMCmZcF7MAxRegbFR&busi_data=eyJncm91cENvZGUiOiIyMTMyNDYxNjQiLCJ0b2tlbiI6InVyMUowRHdaeHMrNHR0cDI5OG1lb1FtT1RWZUpYaVZGcDZRdU5laXVIQnVNS1hJN0g3SEswQUhwNk5qSGR1SmUiLCJ1aW4iOiIyODQ0MTg5MjI4In0%3D&data=oM6AeBGUD-qHiWshlmElnCBhYbDN33jrBXSRaNlOBmztBy7Bx0Sa7dbeoKSRdKzcgUoQjrjSrq5qQVnBA4_W8g&svctype=4&tempid=h5_group_info" target="_blank" rel="noopener noreferrer" className="text-blue-400 underline hover:text-blue-300">213246164</a>（点击加入交流群）
            </p>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-green-400 font-bold text-base mb-3">📋 v1.0.5 更新内容</h3>
            <ul className="list-disc list-inside space-y-1.5 text-dark-300">
              <li>新增识图模型配置：支持配置视觉AI模型，自动识别题目中的图片进行答题</li>
              <li>新增双模型配置界面：纯文本模型（必填）+ 识图模型（选填），含图片题目自动使用识图模型</li>
              <li>修复多任务点模式下视频不观看的问题：增加视频卡片获取重试机制和视频信息获取重试机制</li>
              <li>优化视频播放日志：增加视频开始观看、已完成、跳过等状态提示</li>
              <li>视觉AI降级策略：识图模型失败时自动降级为纯文本模型答题</li>
            </ul>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-yellow-400 font-bold text-base mb-3">⚠️ 免责声明</h3>
            <ol className="list-decimal list-inside space-y-2 text-dark-300">
              <li>本程序仅供学习、研究与技术交流使用，严禁用于任何商业用途或违法活动。</li>
              <li>本程序开源免费，严禁贩卖。若因使用本程序对相关平台或机构造成任何损失，请立即停止使用并删除本程序。</li>
              <li>任何个人或组织使用本程序所从事的一切违法行为，均与作者无关，作者不承担任何法律责任。</li>
              <li>使用本程序即表示您已阅读、理解并同意遵守上述声明。如不同意，请立即删除本程序。</li>
              <li>本程序涉及的部分功能可能违反相关平台的使用条款，使用者需自行承担由此产生的一切风险与后果。</li>
            </ol>
          </div>
        </div>

        <div className="border-t border-dark-700 p-4 flex flex-col items-center gap-3">
          {!readOnly && !canAccept && (
            <p className="text-dark-400 text-sm">
              请仔细阅读声明，<span className="text-accent-400 font-bold">{countdown}</span> 秒后可点击接受
            </p>
          )}
          <div className="flex gap-3">
            {!readOnly && (
              <button
                onClick={handleDecline}
                className="px-6 py-2 rounded-lg bg-dark-700 text-dark-300 hover:bg-dark-600 transition-colors text-sm"
              >
                拒绝并退出
              </button>
            )}
            {readOnly ? (
              <button
                onClick={onClose}
                className="px-6 py-2 rounded-lg bg-gradient-to-r from-accent-500 to-neon-purple text-white hover:opacity-90 transition-opacity text-sm font-medium"
              >
                我已知晓
              </button>
            ) : (
              <button
                onClick={onAccept}
                disabled={!canAccept}
                className={`px-6 py-2 rounded-lg text-sm font-medium transition-all ${
                  canAccept
                    ? 'bg-gradient-to-r from-accent-500 to-neon-purple text-white hover:opacity-90 cursor-pointer'
                    : 'bg-dark-700 text-dark-500 cursor-not-allowed'
                }`}
              >
                {canAccept ? '我已阅读并接受' : `请等待 ${countdown} 秒`}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default DisclaimerModal
