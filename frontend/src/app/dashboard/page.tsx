"use client"

import { useState, useEffect } from "react"
import {
  FileText,
  Bell,
  AlertTriangle,
  Vote,
  ArrowRight,
  Clock,
  CheckCircle,
} from "lucide-react"
import Link from "next/link"
import StatsCard from "@/components/StatsCard"
import ReportCard from "@/components/ReportCard"
import type { BadClientReport, CheckIn } from "@/types"

const mockReports: BadClientReport[] = [
  {
    id: "1",
    encryptedDetails: btoa("Client became aggressive after negotiation. Refused to stop when asked."),
    locationRegion: "North America",
    reportCategory: "Harassment",
    status: "verified",
    confidenceScore: 0.87,
    createdAt: new Date(Date.now() - 3600000 * 2).toISOString(),
  },
  {
    id: "2",
    encryptedDetails: btoa("Left without payment after services were rendered. Blocked all contact."),
    locationRegion: "Europe",
    reportCategory: "Non-payment",
    status: "pending",
    confidenceScore: 0.72,
    createdAt: new Date(Date.now() - 86400000).toISOString(),
  },
]

const mockCheckIns: CheckIn[] = [
  {
    id: "1",
    userId: "0xmock",
    scheduledAt: new Date(Date.now() + 3600000).toISOString(),
    windowMinutes: 30,
    status: "pending",
    createdAt: new Date().toISOString(),
  },
]

const activities = [
  { icon: FileText, text: "New bad client report verified", time: "2h ago", color: "text-success" },
  { icon: Vote, text: "Proposal #12 passed community vote", time: "5h ago", color: "text-primary" },
  { icon: Bell, text: "Check-in reminder sent to 3 members", time: "1d ago", color: "text-warning" },
  { icon: AlertTriangle, text: "Emergency alert resolved successfully", time: "2d ago", color: "text-destructive" },
]

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-sm text-muted-foreground">
          Overview of your safety network
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          icon={Bell}
          label="Active Check-Ins"
          value={mockCheckIns.filter((c) => c.status === "pending").length}
          trend="up"
          trendValue="+2"
        />
        <StatsCard
          icon={FileText}
          label="Reports This Week"
          value={12}
          trend="up"
          trendValue="+3"
        />
        <StatsCard
          icon={AlertTriangle}
          label="Emergency Alerts"
          value={1}
          trend="down"
          trendValue="-2"
        />
        <StatsCard
          icon={Vote}
          label="Active Proposals"
          value={3}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Recent Reports</h2>
            <Link
              href="/registry"
              className="text-sm text-primary hover:underline inline-flex items-center gap-1"
            >
              View all <ArrowRight className="h-3 w-3" />
            </Link>
          </div>
          <div className="space-y-3">
            {mockReports.map((r) => (
              <ReportCard key={r.id} report={r} />
            ))}
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Activity Feed</h2>
          </div>
          <div className="rounded-xl border bg-card shadow-sm">
            {activities.map((a, i) => {
              const Icon = a.icon
              return (
                <div
                  key={i}
                  className="flex items-start gap-3 p-4 border-b last:border-b-0"
                >
                  <div className="h-8 w-8 rounded-full bg-secondary flex items-center justify-center shrink-0">
                    <Icon className={`h-4 w-4 ${a.color}`} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm">{a.text}</p>
                    <p className="text-xs text-muted-foreground">{a.time}</p>
                  </div>
                </div>
              )
            })}
          </div>

          <div className="rounded-xl border bg-card p-5 shadow-sm">
            <h3 className="font-semibold mb-3">Quick Actions</h3>
            <div className="grid grid-cols-2 gap-2">
              <Link
                href="/checkin"
                className="flex items-center gap-2 rounded-lg border p-3 text-sm hover:bg-accent transition-colors"
              >
                <CheckCircle className="h-4 w-4 text-success" />
                Check In
              </Link>
              <Link
                href="/report"
                className="flex items-center gap-2 rounded-lg border p-3 text-sm hover:bg-accent transition-colors"
              >
                <FileText className="h-4 w-4 text-primary" />
                Report
              </Link>
              <Link
                href="/emergency"
                className="flex items-center gap-2 rounded-lg border p-3 text-sm hover:bg-accent transition-colors"
              >
                <AlertTriangle className="h-4 w-4 text-destructive" />
                Emergency
              </Link>
              <Link
                href="/dao"
                className="flex items-center gap-2 rounded-lg border p-3 text-sm hover:bg-accent transition-colors"
              >
                <Vote className="h-4 w-4 text-primary" />
                Vote
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
