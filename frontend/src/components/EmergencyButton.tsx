"use client"

import { useState } from "react"
import { AlertTriangle, X, Loader2 } from "lucide-react"

interface EmergencyButtonProps {
  onTrigger?: () => Promise<void>
}

export default function EmergencyButton({ onTrigger }: EmergencyButtonProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [sent, setSent] = useState(false)

  async function handleConfirm() {
    setLoading(true)
    try {
      await onTrigger?.()
      setSent(true)
    } finally {
      setLoading(false)
    }
  }

  if (sent) {
    return (
      <div className="flex flex-col items-center gap-4 p-8">
        <div className="h-24 w-24 rounded-full bg-success/20 flex items-center justify-center">
          <AlertTriangle className="h-12 w-12 text-success" />
        </div>
        <h2 className="text-2xl font-bold text-success">Alert Sent</h2>
        <p className="text-muted-foreground text-center max-w-md">
          Help is on the way. Stay where you are and remain safe.
        </p>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center gap-6">
      <button
        onClick={() => setOpen(true)}
        className="h-48 w-48 rounded-full bg-destructive hover:bg-destructive/90 text-destructive-foreground shadow-xl hover:shadow-2xl transition-all duration-200 flex flex-col items-center justify-center gap-2 animate-pulse-slow"
        aria-label="Emergency SOS"
      >
        <AlertTriangle className="h-16 w-16" />
        <span className="text-2xl font-bold tracking-wide">SOS</span>
      </button>
      <p className="text-sm text-muted-foreground">Tap to send emergency alert</p>

      {open && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
          <div className="rounded-xl bg-card p-8 max-w-md mx-4 shadow-2xl border">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">Confirm Emergency Alert</h3>
              <button
                onClick={() => setOpen(false)}
                className="h-8 w-8 rounded-md hover:bg-accent flex items-center justify-center"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
            <p className="text-muted-foreground mb-6">
              This will immediately notify emergency contacts and share your
              location. Only use in a genuine emergency.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setOpen(false)}
                className="flex-1 rounded-lg border bg-background px-4 py-2 text-sm font-medium hover:bg-accent transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleConfirm}
                disabled={loading}
                className="flex-1 rounded-lg bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground hover:bg-destructive/90 transition-colors disabled:opacity-50"
              >
                {loading ? (
                  <Loader2 className="h-4 w-4 animate-spin mx-auto" />
                ) : (
                  "Send Alert"
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
