import React, { useEffect, useState } from 'react'
import AppLayout from './components/Layout'
import ErrorBoundary from './components/ErrorBoundary'
import ToastContainer from './components/ToastContainer'
import DisclaimerModal from './components/DisclaimerModal'
import SponsorsModal from './components/SponsorsModal'
import DashboardPage from './pages/DashboardPage'
import AccountsPage from './pages/AccountsPage'
import CheckPage from './pages/CheckPage'
import CheckDetailPage from './pages/CheckDetailPage'
import MonitorPage from './pages/MonitorPage'
import SettingsPage from './pages/SettingsPage'
import { useAppStore } from './stores/useAppStore'
import { parseResponse } from './utils/helpers'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { IsDisclaimerAccepted, AcceptDisclaimer, GetDonatePrompt, DisableDonatePrompt } from '../wailsjs/go/main/App'

const App: React.FC = () => {
  const { currentPage, fetchAccounts, fetchProgress, updateProgress, addLog } = useAppStore()
  const [disclaimerAccepted, setDisclaimerAccepted] = useState<boolean | null>(null)
  const [sponsorCount, setSponsorCount] = useState<number | null>(null)

  useEffect(() => {
    IsDisclaimerAccepted().then(accepted => {
      setDisclaimerAccepted(accepted)
    }).catch(() => {
      setDisclaimerAccepted(false)
    })
  }, [])

  const handleAccept = async () => {
    try {
      await AcceptDisclaimer()
      setDisclaimerAccepted(true)
    } catch (e) {
      console.error('Accept disclaimer failed:', e)
    }
  }

  useEffect(() => {
    if (!disclaimerAccepted) return
    GetDonatePrompt().then((raw) => {
      const resp = parseResponse<{ show: boolean; count: number }>(raw)
      if (resp.code === 200 && resp.data?.show) {
        setSponsorCount(resp.data.count)
      }
    }).catch(() => {})
  }, [disclaimerAccepted])

  const handleDisableSponsor = async () => {
    try { await DisableDonatePrompt() } catch {}
    setSponsorCount(null)
  }

  useEffect(() => {
    if (!disclaimerAccepted) return

    fetchAccounts()
    fetchProgress()

    const interval = setInterval(() => {
      fetchProgress()
      fetchAccounts()
    }, 3000)

    EventsOn('monitor:update', (data: string) => {
      try {
        const parsed = JSON.parse(data)
        if (parsed.event === 'log_update' && parsed.data) {
          addLog(parsed.data.uid, parsed.data.log)
        } else if (parsed.event === 'progress_update' && parsed.data) {
          updateProgress(parsed.data)
        }
      } catch (e) {
        console.error('Event parse error:', e)
      }
    })

    EventsOn('config:reloaded', () => {
      fetchAccounts()
    })

    return () => { clearInterval(interval) }
  }, [disclaimerAccepted])

  if (disclaimerAccepted === null) {
    return (
      <div className="w-full h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-dark-400 text-sm">加载中...</div>
      </div>
    )
  }

  if (!disclaimerAccepted) {
    return <DisclaimerModal onAccept={handleAccept} />
  }

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard': return <DashboardPage />
      case 'accounts': return <AccountsPage />
      case 'check': return <CheckPage />
      case 'check-detail': return <CheckDetailPage />
      case 'monitor': return <MonitorPage />
      case 'settings': return <SettingsPage />
      default: return <DashboardPage />
    }
  }

  return (
    <ErrorBoundary>
      <AppLayout>
        {renderPage()}
      </AppLayout>
      <ToastContainer />
      {sponsorCount !== null && (
        <SponsorsModal
          count={sponsorCount}
          onClose={() => setSponsorCount(null)}
          onDisable={handleDisableSponsor}
        />
      )}
    </ErrorBoundary>
  )
}

export default App
