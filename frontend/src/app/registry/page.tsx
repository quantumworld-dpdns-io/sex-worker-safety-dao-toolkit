"use client"

import { useState } from "react"
import { Search, BookOpen, Filter, SlidersHorizontal } from "lucide-react"
import ReportCard from "@/components/ReportCard"
import type { BadClientReport } from "@/types"

const mockReports: BadClientReport[] = [
  { id: "1", encryptedDetails: btoa("Client became aggressive after negotiation. Refused to stop when asked. I felt unsafe and had to leave through the fire escape."), locationRegion: "North America", reportCategory: "Harassment", status: "verified", confidenceScore: 0.87, createdAt: new Date(Date.now() - 3600000 * 2).toISOString() },
  { id: "2", encryptedDetails: btoa("Left without payment after services were rendered. Blocked all contact methods immediately."), locationRegion: "Europe", reportCategory: "Non-payment", status: "verified", confidenceScore: 0.92, createdAt: new Date(Date.now() - 86400000).toISOString() },
  { id: "3", encryptedDetails: btoa("Threatened to harm me if I didn't comply with requests outside our agreement."), locationRegion: "North America", reportCategory: "Violence", status: "pending", confidenceScore: 0.65, createdAt: new Date(Date.now() - 86400000 * 2).toISOString() },
  { id: "4", encryptedDetails: btoa("Stole personal items from my apartment while I was in the bathroom."), locationRegion: "Europe", reportCategory: "Theft", status: "dismissed", confidenceScore: 0.35, createdAt: new Date(Date.now() - 86400000 * 5).toISOString() },
  { id: "5", encryptedDetails: btoa("Repeatedly contacted me after being told not to. Created fake profiles."), locationRegion: "Asia Pacific", reportCategory: "Harassment", status: "verified", confidenceScore: 0.94, createdAt: new Date(Date.now() - 86400000 * 7).toISOString() },
]

const categories = ["All", "Violence", "Theft", "Harassment", "Non-payment", "Boundary violation"]
const statuses = ["All", "verified", "pending", "dismissed"]

export default function RegistryPage() {
  const [search, setSearch] = useState("")
  const [category, setCategory] = useState("All")
  const [status, setStatus] = useState("All")

  const filtered = mockReports.filter((r) => {
    const matchesSearch = search === "" ||
      r.reportCategory.toLowerCase().includes(search.toLowerCase()) ||
      atob(r.encryptedDetails).toLowerCase().includes(search.toLowerCase())
    const matchesCategory = category === "All" || r.reportCategory === category
    const matchesStatus = status === "All" || r.status === status
    return matchesSearch && matchesCategory && matchesStatus
  })

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Bad Client Registry</h1>
        <p className="text-sm text-muted-foreground">
          Privacy-preserving registry of verified reports
        </p>
      </div>

      <div className="flex flex-col sm:flex-row gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search reports..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-lg border bg-background pl-10 pr-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
          />
        </div>
        <select
          value={category}
          onChange={(e) => setCategory(e.target.value)}
          className="rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
        >
          {categories.map((c) => (
            <option key={c} value={c}>{c}</option>
          ))}
        </select>
        <select
          value={status}
          onChange={(e) => setStatus(e.target.value)}
          className="rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
        >
          {statuses.map((s) => (
            <option key={s} value={s}>{s.charAt(0).toUpperCase() + s.slice(1)}</option>
          ))}
        </select>
      </div>

      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Showing {filtered.length} of {mockReports.length} reports
        </p>
      </div>

      {filtered.length === 0 ? (
        <div className="text-center py-12">
          <BookOpen className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="font-semibold mb-1">No reports found</h3>
          <p className="text-sm text-muted-foreground">
            Try adjusting your search or filters
          </p>
        </div>
      ) : (
        <div className="grid gap-3 sm:grid-cols-2">
          {filtered.map((r) => (
            <ReportCard key={r.id} report={r} />
          ))}
        </div>
      )}
    </div>
  )
}
