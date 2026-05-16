"use client"

import { LucideIcon, TrendingUp, TrendingDown } from "lucide-react"
import { cn } from "@/lib/utils"

interface StatsCardProps {
  icon: LucideIcon
  label: string
  value: string | number
  trend?: "up" | "down"
  trendValue?: string
  className?: string
}

export default function StatsCard({
  icon: Icon,
  label,
  value,
  trend,
  trendValue,
  className,
}: StatsCardProps) {
  return (
    <div
      className={cn(
        "rounded-xl border bg-card p-5 text-card-foreground shadow-sm",
        className,
      )}
    >
      <div className="flex items-center justify-between">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
          <Icon className="h-5 w-5 text-primary" />
        </div>
        {trend && (
          <span
            className={cn(
              "inline-flex items-center gap-1 text-xs font-medium",
              trend === "up" ? "text-success" : "text-destructive",
            )}
          >
            {trend === "up" ? (
              <TrendingUp className="h-3 w-3" />
            ) : (
              <TrendingDown className="h-3 w-3" />
            )}
            {trendValue}
          </span>
        )}
      </div>
      <div className="mt-3">
        <p className="text-sm text-muted-foreground">{label}</p>
        <p className="text-2xl font-bold">{value}</p>
      </div>
    </div>
  )
}
