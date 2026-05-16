"use client"

import { useState } from "react"
import { CheckCircle, Clock, AlertTriangle, Loader2 } from "lucide-react"
import { cn } from "@/lib/utils"

interface CheckInWidgetProps {
  status?: "pending" | "completed" | "missed" | "emergency"
  scheduledAt?: string
  onSchedule?: () => Promise<void>
  onComplete?: () => Promise<void>
}

export default function CheckInWidget({
  status = "pending",
  scheduledAt,
  onSchedule,
  onComplete,
}: CheckInWidgetProps) {
  const [completing, setCompleting] = useState(false)
  const [scheduling, setScheduling] = useState(false)

  const statusConfig = {
    pending: {
      icon: Clock,
      color: "text-warning",
      bg: "bg-warning/10",
      label: "Pending",
    },
    completed: {
      icon: CheckCircle,
      color: "text-success",
      bg: "bg-success/10",
      label: "Completed",
    },
    missed: {
      icon: AlertTriangle,
      color: "text-destructive",
      bg: "bg-destructive/10",
      label: "Missed",
    },
    emergency: {
      icon: AlertTriangle,
      color: "text-destructive",
      bg: "bg-destructive/10",
      label: "Emergency",
    },
  }

  const config = statusConfig[status]
  const StatusIcon = config.icon

  async function handleComplete() {
    setCompleting(true)
    try {
      await onComplete?.()
    } finally {
      setCompleting(false)
    }
  }

  async function handleSchedule() {
    setScheduling(true)
    try {
      await onSchedule?.()
    } finally {
      setScheduling(false)
    }
  }

  return (
    <div className="rounded-xl border bg-card p-6 shadow-sm">
      <div className="flex items-center gap-4 mb-4">
        <div className={cn("h-12 w-12 rounded-full flex items-center justify-center", config.bg)}>
          <StatusIcon className={cn("h-6 w-6", config.color)} />
        </div>
        <div>
          <p className="text-sm text-muted-foreground">Current Status</p>
          <p className={cn("text-lg font-semibold", config.color)}>
            {config.label}
          </p>
        </div>
      </div>

      {scheduledAt && (
        <p className="text-sm text-muted-foreground mb-4">
          Next check-in scheduled for{" "}
          <span className="font-medium text-foreground">
            {new Date(scheduledAt).toLocaleString()}
          </span>
        </p>
      )}

      <div className="flex gap-3">
        {status === "pending" && (
          <button
            onClick={handleComplete}
            disabled={completing}
            className="flex-1 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            {completing ? (
              <Loader2 className="h-4 w-4 animate-spin mx-auto" />
            ) : (
              "Complete Check-In"
            )}
          </button>
        )}
        {status === "completed" && (
          <button
            onClick={handleSchedule}
            disabled={scheduling}
            className="flex-1 rounded-lg border bg-background px-4 py-2 text-sm font-medium hover:bg-accent transition-colors disabled:opacity-50"
          >
            {scheduling ? (
              <Loader2 className="h-4 w-4 animate-spin mx-auto" />
            ) : (
              "Schedule Next"
            )}
          </button>
        )}
      </div>
    </div>
  )
}
