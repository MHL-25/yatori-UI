import React from 'react'
import { Heart, X, ExternalLink } from 'lucide-react'
import { motion } from 'framer-motion'

interface Props {
  count: number
  onClose: () => void
  onDisable: () => void
}

const SponsorsModal: React.FC<Props> = ({ count, onClose, onDisable }) => {
  return (
    <div className="fixed inset-0 z-[9998] flex items-center justify-center bg-black/70 backdrop-blur-sm px-4">
      <motion.div
        initial={{ opacity: 0, scale: 0.94, y: 12 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        transition={{ duration: 0.28, ease: [0.22, 1, 0.36, 1] }}
        className="relative w-[520px] max-w-full rounded-2xl border border-accent-500/50 bg-dark-900 p-6 text-center"
      >
        <button onClick={onClose} className="absolute right-4 top-4 text-dark-400 hover:text-dark-200" title="关闭">
          <X size={18} />
        </button>

        <div className="mx-auto w-16 h-16 rounded-2xl bg-accent-500 flex items-center justify-center">
          <Heart size={30} className="text-dark-950" />
        </div>

        <h2 className="mt-4 text-xl font-bold text-dark-50">
          您已经使用了 <span className="text-accent-400 font-mono text-2xl">{count}</span> 次本软件
        </h2>
        <p className="mt-3 text-sm text-dark-200 leading-relaxed">
          如果觉得不错的话，可以为我发电吗（双手合十）
        </p>
        <p className="mt-1 text-sm text-dark-300">
          你的赞助将会是我维护的最大动力，当然也可以不赞助.....哦
        </p>

        <a
          href="https://afdian.com/a/yatori-UI"
          target="_blank"
          rel="noopener noreferrer"
          className="btn-primary mt-5 inline-flex items-center gap-2 px-6 py-2.5"
        >
          <ExternalLink size={15} /> 前往爱发电支持
        </a>
        <p className="mt-2 text-xs font-mono text-accent-400">https://afdian.com/a/yatori-UI</p>

        <div className="mt-5 flex items-center justify-center gap-3">
          <button onClick={onClose} className="btn-secondary">暂不</button>
          <button onClick={onDisable} className="btn-ghost text-xs text-dark-400 hover:text-dark-200">
            不再提示
          </button>
        </div>
      </motion.div>
    </div>
  )
}

export default SponsorsModal
