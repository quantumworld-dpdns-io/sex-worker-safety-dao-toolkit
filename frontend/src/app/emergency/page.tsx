"use client"

import { useState } from "react"
import { AlertTriangle, CheckCircle, Clock, History } from "lucide-react"
import EmergencyButton from "@/components/EmergencyButton"
import { timeAgo } from "@/lib/utils"

interface AlertRecord {
  id: string
  triggerType: string
  resolvedAt?: string
  createdAt: string
}

const mockAlerts: AlertRecord[] = [
  { id: "1", triggerType: "manual", resolvedAt: new Date(Date.now() - 86400000).toISOString(), createdAt: new Date(Date.now() - 86400000 * 2).toISOString() },
  { id: "2", triggerType: "missed_checkin", resolvedAt: new Date(Date.now() - 86400000 * 5).toISOString(), createdAt: new Date(Date.now() - 86400000 * 5).toISOString() },
]

export default function EmergencyPage() {
  const [alertTriggered, setAlertTriggered] = useState(false)

  async function handleTrigger() {
    await new Promise((r) => setTimeout(r, 2000))
    setAlertTriggered(true)
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold">Emergency</h1>
        <p className="text-sm text-muted-foreground">
          Immediate panic alert with location sharing
        </p>
      </div>

      <div className="flex justify-center py-8">
        <EmergencyButton onTrigger={handleTrigger} />
      </div>

      {alertTriggered && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-6 text-center shadow-sm">
          <AlertTriangle className="h-8 w-8 text-destructive mx-auto mb-2" />
          <h3 className="font-semibold text-destructive mb-1">Alert Active</h3>
          <p className="text-sm text-muted-foreground">
            Your emergency contacts have been notified. Stay safe.
          </p>
        </div>
      )}

      <div className="rounded-xl border bg-card p-6 shadow-sm">
        <div className="flex items-center gap-2 mb-4">
          <ShieldIcon className="h-5 w-5 text-muted-foreground" />
          <h2 className="font-semibold">Safety Tips</h2>
        </div>
        <ul className="space-y-2 text-sm text-muted-foreground">
          <li className="flex items-start gap-2">
            <CheckCircle className="h-4 w-4 text-success mt-0.5 shrink-0" />
            <span>Always share your location with a trusted contact before sessions</span>
          </li>
          <li className="flex items-start gap-2">
            <CheckCircle className="h-4 w-4 text-success mt-0.5 shrink-0" />
            <span>Set up regular check-in schedules that match your work hours</span>
          </li>
          <li className="flex items-start gap-2">
            <CheckCircle className="h-4 w-4 text-success mt-0.5 shrink-0" />
            <span>Use the SOS button at the first sign of danger</span>
          </li>
          <li className="flex items-start gap-2">
            <CheckCircle className="h-4 w-4 text-success mt-0.5 shrink-0" />
            <span>Memorize a safe word that triggers emergency from a phone call</span>
          </li>
        </ul>
      </div>

      <div className="space-y-3">
        <div className="flex items-center gap-2">
          <History className="h-5 w-5 text-muted-foreground" />
          <h2 className="text-lg font-semibold">Alert History</h2>
        </div>
        <div className="rounded-xl border bg-card shadow-sm divide-y">
          {mockAlerts.map((a) => (
            <div key={a.id} className="flex items-center gap-4 p-4">
              <div className={`h-8 w-8 rounded-full flex items-center justify-center ${
                a.resolvedAt ? "bg-success/10" : "bg-destructive/10"
              }`}>
                {a.resolvedAt ? (
                  <CheckCircle className="h-4 w-4 text-success" />
                ) : (
                  <AlertTriangle className="h-4 w-4 text-destructive" />
                )}
              </div>
              <div className="flex-1">
                <p className="text-sm font-medium capitalize">{a.triggerType.replace(/_/g, " ")}</p>
                <p className="text-xs text-muted-foreground">
                  {a.resolvedAt ? "Resolved" : "Active"} &middot; {timeAgo(a.createdAt)}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function ShieldIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
    </svg>
  )
}
