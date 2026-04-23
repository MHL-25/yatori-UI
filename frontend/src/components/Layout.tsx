import React from 'react'
import { useAppStore } from '../stores/useAppStore'
import {
  LayoutDashboard, Users, Activity, Settings,
  Minus, Square, X, ChevronLeft, ChevronRight, Zap, MinusSquare
} from 'lucide-react'
import { cn } from '../utils/helpers'

const navItems = [
  { id: 'dashboard', label: '仪表盘', icon: LayoutDashboard },
  { id: 'accounts', label: '账号管理', icon: Users },
  { id: 'monitor', label: '任务监控', icon: Activity },
  { id: 'settings', label: '系统设置', icon: Settings },
]

const TitleBar: React.FC = () => {
  const handleMinimize = async () => {
    try {
      const { MinimizeWindow: fn } = await import('../../wailsjs/go/main/App')
      await fn()
    } catch {}
  }
  const handleMaximize = async () => {
    try {
      const { ToggleMaximizeWindow: fn } = await import('../../wailsjs/go/main/App')
      await fn()
    } catch {}
  }
  const handleQuit = async () => {
    try {
      const { QuitApp: fn } = await import('../../wailsjs/go/main/App')
      await fn()
    } catch {}
  }
  const handleClose = async () => {
    try {
      const { CloseWindow: fn } = await import('../../wailsjs/go/main/App')
      await fn()
    } catch {}
  }

  return (
    <div className="wails-drag flex items-center justify-between h-9 bg-dark-950/80 border-b border-dark-800/50 px-3 select-none">
      <div className="flex items-center gap-2">
        <div className="w-5 h-5 rounded-md bg-gradient-to-br from-accent-500 to-neon-purple flex items-center justify-center">
          <Zap size={12} className="text-white" />
        </div>
        <span className="text-xs font-medium text-dark-300">Yatori-UI - 智能网课助手</span>
      </div>
      <div className="wails-no-drag flex items-center gap-0.5">
        <button
          onClick={handleMinimize}
          title="最小化"
          className="w-8 h-7 flex items-center justify-center rounded hover:bg-dark-700/50 transition-colors"
        >
          <Minus size={14} className="text-dark-400" />
        </button>
        <button
          onClick={handleMaximize}
          title="最大化"
          className="w-8 h-7 flex items-center justify-center rounded hover:bg-dark-700/50 transition-colors"
        >
          <Square size={10} className="text-dark-400" />
        </button>
        <button
          onClick={handleClose}
          title="最小化到托盘"
          className="w-8 h-7 flex items-center justify-center rounded hover:bg-dark-700/50 transition-colors"
        >
          <MinusSquare size={12} className="text-dark-400" />
        </button>
        <button
          onClick={handleQuit}
          title="退出程序"
          className="w-8 h-7 flex items-center justify-center rounded hover:bg-red-600/80 transition-colors"
        >
          <X size={14} className="text-dark-200" />
        </button>
      </div>
    </div>
  )
}

const Sidebar: React.FC = () => {
  const { currentPage, setCurrentPage, sidebarCollapsed, toggleSidebar } = useAppStore()

  return (
    <div className={cn(
      "flex flex-col h-full bg-dark-950/60 border-r border-dark-800/50 overflow-hidden",
      "transition-[width] duration-300 ease-in-out",
      sidebarCollapsed ? "w-16" : "w-56"
    )}>
      <div className="flex-1 py-3 px-2 space-y-1">
        {navItems.map((item) => (
          <button
            key={item.id}
            onClick={() => setCurrentPage(item.id)}
            className={cn(
              "sidebar-item w-full",
              currentPage === item.id && "sidebar-item-active"
            )}
          >
            <item.icon size={18} className="flex-shrink-0" />
            <span className={cn(
              "whitespace-nowrap overflow-hidden transition-opacity duration-200",
              sidebarCollapsed ? "opacity-0 w-0" : "opacity-100"
            )}>{item.label}</span>
          </button>
        ))}
      </div>

      <div className="p-2 border-t border-dark-800/50">
        <button
          onClick={toggleSidebar}
          className="sidebar-item w-full"
        >
          {sidebarCollapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
          <span className={cn(
            "whitespace-nowrap overflow-hidden transition-opacity duration-200",
            sidebarCollapsed ? "opacity-0 w-0" : "opacity-100"
          )}>收起侧栏</span>
        </button>
      </div>
    </div>
  )
}

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <div className="flex flex-col h-screen w-screen bg-dark-950 overflow-hidden">
      <TitleBar />
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <main className="flex-1 overflow-auto bg-dark-950">
          {children}
        </main>
      </div>
    </div>
  )
}

export default AppLayout
