"use client"

import Link from "next/link"
import { Clock, CheckCircle, XCircle, AlertCircle } from "lucide-react"
import { cn, formatDateShort } from "@/lib/utils"
import type { DAOProposal } from "@/types"

const statusConfig: Record<string, { icon: typeof Clock; color: string; bg: string; label: string }> = {
  active: { icon: Clock, color: "text-primary", bg: "bg-primary/10", label: "Active" },
  passed: { icon: CheckCircle, color: "text-success", bg: "bg-success/10", label: "Passed" },
  rejected: { icon: XCircle, color: "text-destructive", bg: "bg-destructive/10", label: "Rejected" },
  pending: { icon: AlertCircle, color: "text-warning", bg: "bg-warning/10", label: "Pending" },
}

export default function ProposalCard({ proposal }: { proposal: DAOProposal }) {
  const config = statusConfig[proposal.status] || statusConfig.pending
  const StatusIcon = config.icon

  return (
    <Link href={`/dao/proposals/${proposal.id}`}>
      <div className="rounded-xl border bg-card p-5 shadow-sm hover:shadow-md hover:border-primary/30 transition-all cursor-pointer">
        <div className="flex items-start justify-between mb-2">
          <span className={cn("inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium", config.bg, config.color)}>
            <StatusIcon className="h-3 w-3 mr-1" />
            {config.label}
          </span>
          <span className="text-xs text-muted-foreground">
            {proposal.proposalType}
          </span>
        </div>

        <h3 className="font-semibold mb-1 line-clamp-1">{proposal.title}</h3>
        <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
          {proposal.description}
        </p>

        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <span>Voting: {proposal.votingType}</span>
          <span>Ends: {formatDateShort(proposal.endsAt)}</span>
        </div>
      </div>
    </Link>
  )
}
