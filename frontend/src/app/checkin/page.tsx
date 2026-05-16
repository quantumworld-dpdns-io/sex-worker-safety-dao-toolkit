"use client"

import { useState } from "react"
import { CheckCircle, Clock, Calendar, Plus, History } from "lucide-react"
import CheckInWidget from "@/components/CheckInWidget"
import { formatDate, timeAgo } from "@/lib/utils"
import type { CheckIn } from "@/types"

const mockCheckIns: CheckIn[] = [
  { id: "1", userId: "0xmock", scheduledAt: new Date(Date.now() + 3600000).toISOString(), windowMinutes: 30, status: "pending", createdAt: new Date(Date.now() - 3600000 * 24).toISOString() },
  { id: "2", userId: "0xmock", scheduledAt: new Date(Date.now() - 86400000).toISOString(), windowMinutes: 30, status: "completed", completedAt: new Date(Date.now() - 86400000 + 1800000).toISOString(), createdAt: new Date(Date.now() - 86400000 * 2).toISOString() },
  { id: "3", userId: "0xmock", scheduledAt: new Date(Date.now() - 86400000 * 2).toISOString(), windowMinutes: 30, status: "completed", completedAt: new Date(Date.now() - 86400000 * 2 + 1200000).toISOString(), createdAt: new Date(Date.now() - 86400000 * 3).toISOString() },
]

export default function CheckInPage() {
  const [checkIns, setCheckIns] = useState(mockCheckIns)
  const [scheduleTime, setScheduleTime] = useState("")
  const [scheduleMinutes, setScheduleMinutes] = useState(30)
  const [showSchedule, setShowSchedule] = useState(false)

  const currentCheckIn = checkIns.find((c) => c.status === "pending")

  async function handleComplete() {
    await new Promise((r) => setTimeout(r, 1000))
    setCheckIns((prev) =>
      prev.map((c) =>
        c.status === "pending"
          ? { ...c, status: "completed" as const, completedAt: new Date().toISOString() }
          : c,
      ),
    )
  }

  async function handleSchedule() {
    if (!scheduleTime) return
    await new Promise((r) => setTimeout(r, 1000))
    const newCheckIn: CheckIn = {
      id: crypto.randomUUID(),
      userId: "0xmock",
      scheduledAt: new Date(scheduleTime).toISOString(),
      windowMinutes: scheduleMinutes,
      status: "pending",
      createdAt: new Date().toISOString(),
    }
    setCheckIns((prev) => [newCheckIn, ...prev])
    setShowSchedule(false)
    setScheduleTime("")
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Safety Check-Ins</h1>
        <p className="text-sm text-muted-foreground">
          Regular check-ins ensure you're safe. Missed check-ins trigger emergency protocols.
        </p>
      </div>

      {currentCheckIn ? (
        <CheckInWidget
          status={currentCheckIn.status}
          scheduledAt={currentCheckIn.scheduledAt}
          onComplete={handleComplete}
        />
      ) : (
        <div className="rounded-xl border bg-card p-8 text-center shadow-sm">
          <CheckCircle className="h-12 w-12 text-success mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">All Checked In</h3>
          <p className="text-sm text-muted-foreground mb-4">
            You have no pending check-ins. Schedule your next one.
          </p>
          <button
            onClick={() => setShowSchedule(true)}
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
          >
            <Plus className="h-4 w-4" />
            Schedule Check-In
          </button>
        </div>
      )}

      {showSchedule && (
        <div className="rounded-xl border bg-card p-6 shadow-sm space-y-4">
          <h3 className="font-semibold">Schedule a Check-In</h3>
          <div className="grid sm:grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium mb-1.5 block" htmlFor="time">
                Date & Time
              </label>
              <input
                id="time"
                type="datetime-local"
                value={scheduleTime}
                onChange={(e) => setScheduleTime(e.target.value)}
                className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              />
            </div>
            <div>
              <label className="text-sm font-medium mb-1.5 block" htmlFor="window">
                Response Window (minutes)
              </label>
              <select
                id="window"
                value={scheduleMinutes}
                onChange={(e) => setScheduleMinutes(Number(e.target.value))}
                className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              >
                <option value={15}>15 minutes</option>
                <option value={30}>30 minutes</option>
                <option value={60}>1 hour</option>
                <option value={120}>2 hours</option>
              </select>
            </div>
          </div>
          <div className="flex gap-3">
            <button
              onClick={() => setShowSchedule(false)}
              className="flex-1 rounded-lg border bg-background px-4 py-2 text-sm font-medium hover:bg-accent transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSchedule}
              disabled={!scheduleTime}
              className="flex-1 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
            >
              Schedule
            </button>
          </div>
        </div>
      )}

      <div className="space-y-3">
        <div className="flex items-center gap-2">
          <History className="h-5 w-5 text-muted-foreground" />
          <h2 className="text-lg font-semibold">Check-In History</h2>
        </div>
        {checkIns.length === 0 ? (
          <p className="text-sm text-muted-foreground text-center py-8">
            No check-in history yet
          </p>
        ) : (
          <div className="rounded-xl border bg-card shadow-sm divide-y">
            {checkIns.map((ci) => (
              <div key={ci.id} className="flex items-center gap-4 p-4">
                <div className={`h-8 w-8 rounded-full flex items-center justify-center ${
                  ci.status === "completed" ? "bg-success/10" :
                  ci.status === "missed" ? "bg-destructive/10" :
                  ci.status === "emergency" ? "bg-destructive/10" :
                  "bg-warning/10"
                }`}>
                  {ci.status === "completed" ? (
                    <CheckCircle className="h-4 w-4 text-success" />
                  ) : (
                    <Clock className="h-4 w-4 text-warning" />
                  )}
                </div>
                <div className="flex-1">
                  <p className="text-sm font-medium capitalize">{ci.status}</p>
                  <p className="text-xs text-muted-foreground">
                    {formatDate(ci.scheduledAt)}
                  </p>
                </div>
                <span className="text-xs text-muted-foreground">
                  {timeAgo(ci.createdAt)}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
