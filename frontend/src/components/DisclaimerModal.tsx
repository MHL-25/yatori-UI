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
            <h3 className="text-red-400 font-bold text-base mb-3">🔥 v1.1.3 严重Bug修复</h3>
            <ul className="list-disc list-inside space-y-1.5 text-dark-300">
              <li>【关键修复】多任务点模式和正常刷课模式无法使用：移除nodeRun中的调试代码块，恢复正确视频处理逻辑</li>
              <li>【修复】日志中大量乱码中文字符：修复约50处乱码字符，提升用户体验</li>
              <li>【修复】暂停后状态标记错误：为所有Activity添加isPaused()检查，确保暂停状态正确</li>
            </ul>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-red-400 font-bold text-base mb-3">🔥 v1.1.2 全面对齐原项目</h3>
            <ul className="list-disc list-inside space-y-1.5 text-dark-300">
              <li>【关键修复】课程级并发逻辑：VideoModel==1时串行，VideoModel!=1时并发，与原项目完全对齐</li>
              <li>【关键修复】节点级并发队列机制：VideoModel==3时使用队列管理并发资源，与原项目一致</li>
              <li>【关键修复】视频被错误跳过为"非任务点"问题：Go语言range循环变量副本问题已修复</li>
              <li>【关键修复】model3Caches初始化逻辑：先添加到数组再重新登录，与原项目一致</li>
              <li>【修复】所有任务点类型（视频/音频/文档/作业/外链/直播/讨论）的指针问题全部修复</li>
              <li>【修复】作业/考试状态值语言错误：使用中文状态值（待做/未交/待重做/待重考）</li>
              <li>【修复】章节测试缺失题型：名词解释、论述题、连线题</li>
              <li>【修复】ExamAutoSubmit==2模式空答案检测</li>
              <li>【修复】考试限时提交处理</li>
              <li>【修复】暂停后状态标记错误</li>
            </ul>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-orange-400 font-bold text-base mb-3">📋 v1.1.1 核心修复</h3>
            <ul className="list-disc list-inside space-y-1.5 text-dark-300">
              <li>【关键修复】视频被错误跳过为"非任务点"问题：Go语言range循环变量是副本，AttachmentsDetection修改副本而非原数组，导致IsJob始终为false，现已修复</li>
              <li>【关键修复】model3Caches初始化逻辑错误：原项目先添加到数组再重新登录，新项目错误地修改临时变量，现已修复</li>
              <li>【修复】所有任务点类型（视频/音频/文档/作业/外链/直播/讨论）的指针问题全部修复</li>
            </ul>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-orange-400 font-bold text-base mb-3">📋 v1.1.0 紧急修复</h3>
            <ul className="list-disc list-inside space-y-1.5 text-dark-300">
              <li>【紧急修复】作业/考试状态值语言错误：原项目使用中文状态值（待做/未交/待重做/待重考），新项目错误使用英文导致作业考试无法识别，现已修复</li>
              <li>【重要修复】章节测试新增缺失题型支持：名词解释、论述题、连线题三种题型自动答题</li>
              <li>【重要修复】ExamAutoSubmit==2模式空答案检测：有空答案时不自动提交，避免0分</li>
              <li>【新增】考试限时提交处理：自动检测限制并延时重新提交</li>
              <li>【修复】暂停后状态标记错误：暂停后不再误显示为"已完成"</li>
              <li>【修复】日志中多处乱码中文字符修复</li>
            </ul>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-green-400 font-bold text-base mb-3">📋 v1.0.9 更新内容</h3>
            <ul className="list-disc list-inside space-y-1.5 text-dark-300">
              <li>【重要修复】多任务点模式(VModel2/3)课程并发处理：对齐原项目逻辑，VideoModel!=1时课程级并发执行，大幅提升刷课效率</li>
              <li>【重要修复】章节测试新增缺失题型支持：名词解释(TermExplanation)、论述题(Essay)、连线题(Matching)三种题型自动答题</li>
              <li>【重要修复】ExamAutoSubmit==2模式空答案检测：答题后自动检测是否存在空答案，有空答案时不自动提交，避免0分</li>
              <li>【新增】考试限时提交处理：自动检测"考试N分钟内不允许提交"限制，延时后自动重新提交</li>
              <li>【新增】考试时间已用完检测：增加"考试时间已用完"中文检测，避免超时提交报错</li>
              <li>【修复】暂停(Pause)后状态标记错误：暂停后不再误显示为"已完成"，正确保持"已暂停"状态</li>
              <li>【修复】所有Activity类型的Start/Pause方法增加isPaused状态管理</li>
              <li>【修复】日志中多处乱码中文字符修复</li>
            </ul>
          </div>

          <div className="border-t border-dark-700 pt-4">
            <h3 className="text-blue-400 font-bold text-base mb-3">📋 v1.0.5 更新内容</h3>
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
