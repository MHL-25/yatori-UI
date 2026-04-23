import React, { useEffect, useState } from 'react'
import { CheckCircle, AlertCircle, X } from 'lucide-react'
import { cn } from '../utils/helpers'
import { motion, AnimatePresence } from 'framer-motion'
import { EventsOn } from '../../wailsjs/runtime/runtime'

interface Toast {
  id: number
  type: 'success' | 'error'
  title: string
  message: string
}

let toastId = 0

const ToastContainer: React.FC = () => {
  const [toasts, setToasts] = useState<Toast[]>([])

  useEffect(() => {
    EventsOn('notification', (data: any) => {
      const id = ++toastId
      const toast: Toast = {
        id,
        type: data.type || 'success',
        title: data.title || '通知',
        message: data.message || '',
      }
      setToasts(prev => [...prev, toast])
      setTimeout(() => {
        setToasts(prev => prev.filter(t => t.id !== id))
      }, 5000)
    })
  }, [])

  const removeToast = (id: number) => {
    setToasts(prev => prev.filter(t => t.id !== id))
  }

  return (
    <div className="fixed top-12 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
      <AnimatePresence>
        {toasts.map(toast => (
          <motion.div
            key={toast.id}
            initial={{ opacity: 0, x: 80, scale: 0.9 }}
            animate={{ opacity: 1, x: 0, scale: 1 }}
            exit={{ opacity: 0, x: 80, scale: 0.9 }}
            transition={{ duration: 0.25, ease: 'easeOut' }}
            className={cn(
              "pointer-events-auto min-w-[320px] max-w-[420px] p-4 rounded-xl shadow-2xl border backdrop-blur-md",
              "flex items-start gap-3",
              toast.type === 'success'
                ? "bg-emerald-950/90 border-emerald-600/30"
                : "bg-red-950/90 border-red-600/30"
            )}
          >
            <div className={cn(
              "w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0",
              toast.type === 'success' ? "bg-emerald-600/20" : "bg-red-600/20"
            )}>
              {toast.type === 'success'
                ? <CheckCircle size={18} className="text-emerald-400" />
                : <AlertCircle size={18} className="text-red-400" />
              }
            </div>
            <div className="flex-1 min-w-0">
              <p className={cn(
                "text-sm font-semibold",
                toast.type === 'success' ? "text-emerald-300" : "text-red-300"
              )}>{toast.title}</p>
              <p className="text-xs text-dark-300 mt-0.5 break-words">{toast.message}</p>
            </div>
            <button onClick={() => removeToast(toast.id)}
              className="text-dark-500 hover:text-dark-300 flex-shrink-0 mt-0.5">
              <X size={14} />
            </button>
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  )
}

export default ToastContainer
