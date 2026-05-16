"use client"

import { Shield, ShieldCheck, ShieldX } from "lucide-react"
import { cn, formatDateShort, timeAgo } from "@/lib/utils"
import type { BadClientReport } from "@/types"

const statusConfig = {
  pending: { icon: Shield, color: "text-warning", bg: "bg-warning/10", label: "Pending" },
  verified: { icon: ShieldCheck, color: "text-success", bg: "bg-success/10", label: "Verified" },
  dismissed: { icon: ShieldX, color: "text-muted-foreground", bg: "bg-muted", label: "Dismissed" },
}

export default function ReportCard({ report }: { report: BadClientReport }) {
  const config = statusConfig[report.status]
  const StatusIcon = config.icon

  return (
    <div className="rounded-xl border bg-card p-5 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className={cn("inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium", config.bg, config.color)}>
            <StatusIcon className="h-3 w-3 mr-1" />
            {config.label}
          </span>
          <span className="inline-flex items-center rounded-full bg-secondary px-2.5 py-0.5 text-xs font-medium text-secondary-foreground">
            {report.reportCategory}
          </span>
        </div>
        <span className="text-xs text-muted-foreground">
          {timeAgo(report.createdAt)}
        </span>
      </div>

      <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
        {atob(report.encryptedDetails).slice(0, 120)}...
      </p>

      <div className="flex items-center justify-between text-xs">
        <span className="text-muted-foreground">
          {report.locationRegion || "Region unknown"}
        </span>
        <span className="font-medium">
          Confidence: {Math.round(report.confidenceScore * 100)}%
        </span>
      </div>
    </div>
  )
}
