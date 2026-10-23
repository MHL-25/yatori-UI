import React from 'react'
import { Heart } from 'lucide-react'

const AfdianBanner: React.FC = () => {
  return (
    <a
      href="https://afdian.com/a/yatori-UI"
      target="_blank"
      rel="noopener noreferrer"
      className="group flex items-center gap-3 rounded-2xl border border-accent-500/40 bg-accent-500/10 px-4 py-3 transition-colors hover:bg-accent-500/15"
    >
      <div className="w-9 h-9 rounded-xl bg-accent-500 flex items-center justify-center shrink-0">
        <Heart size={17} className="text-dark-950" />
      </div>
      <div className="min-w-0">
        <p className="text-sm font-medium text-dark-100">
          如果觉得本项目做的不错，可以赞助一下吗（双手合十）
        </p>
        <p className="text-xs text-dark-300 mt-0.5">
          你的赞助将会是我维护的最大动力，当然也可以不赞助.....哦
        </p>
      </div>
      <span className="ml-auto text-xs font-semibold text-accent-400 shrink-0 group-hover:text-accent-300">
        爱发电 →
      </span>
    </a>
  )
}

export default AfdianBanner
