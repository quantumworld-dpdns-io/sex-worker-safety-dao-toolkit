"use client"

import { useParams } from "next/navigation"
import Link from "next/link"
import { ArrowLeft, Calendar, User, Vote as VoteIcon } from "lucide-react"
import VoteButton from "@/components/VoteButton"
import { formatDate, cn } from "@/lib/utils"
import type { DAOProposal } from "@/types"

const mockProposal: DAOProposal = {
  id: "1",
  title: "Increase verification threshold for bad client reports",
  description: "Proposal to increase the minimum confidence score from 0.7 to 0.85 for automatic verification of bad client reports to reduce false positives.\n\nCurrently, reports with a confidence score of 0.7 or above are automatically verified by the system. While this has helped build the registry quickly, there have been concerns about false positives.\n\nThis proposal suggests raising the threshold to 0.85, which would mean:\n- Fewer auto-verified reports\n- More manual review by moderators\n- Higher overall quality of the registry\n\nAdditional funding for moderator compensation may be needed if this passes.",
  proposalType: "Parameter Change",
  status: "active",
  votingType: "Quadratic",
  startsAt: new Date(Date.now() - 86400000 * 3).toISOString(),
  endsAt: new Date(Date.now() + 86400000 * 4).toISOString(),
  createdBy: "0x1234...5678",
}

export default function ProposalDetailPage() {
  const params = useParams()
  const proposal = mockProposal

  const results = {
    yes: 125000,
    no: 42000,
    abstain: 8000,
    total: 175000,
  }

  async function handleVote(proposalId: string, vote: "yes" | "no" | "abstain") {
    await new Promise((r) => setTimeout(r, 1000))
    console.log("Voted", vote, "on proposal", proposalId)
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <Link
        href="/dao"
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to proposals
      </Link>

      <div className="rounded-xl border bg-card p-6 shadow-sm">
        <div className="flex items-start justify-between mb-4">
          <div>
            <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary mb-2">
              {proposal.status}
            </span>
            <span className="inline-flex items-center rounded-full bg-secondary px-2.5 py-0.5 text-xs font-medium text-secondary-foreground ml-2 mb-2">
              {proposal.proposalType}
            </span>
            <h1 className="text-xl font-bold">{proposal.title}</h1>
          </div>
        </div>

        <div className="flex items-center gap-4 text-sm text-muted-foreground mb-6">
          <span className="flex items-center gap-1">
            <User className="h-3.5 w-3.5" />
            {proposal.createdBy}
          </span>
          <span className="flex items-center gap-1">
            <Calendar className="h-3.5 w-3.5" />
            Ends {formatDate(proposal.endsAt)}
          </span>
          <span className="flex items-center gap-1">
            <VoteIcon className="h-3.5 w-3.5" />
            {proposal.votingType} voting
          </span>
        </div>

        <div className="prose prose-sm max-w-none text-muted-foreground whitespace-pre-line mb-6">
          {proposal.description}
        </div>

        <div className="rounded-lg border p-4">
          <h3 className="font-semibold mb-4">Cast Your Vote</h3>
          <VoteButton proposalId={proposal.id} onVote={handleVote} />
        </div>
      </div>

      <div className="rounded-xl border bg-card p-6 shadow-sm">
        <h2 className="text-lg font-semibold mb-4">Current Results</h2>
        <div className="space-y-4">
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-success font-medium">Yes</span>
              <span className="text-muted-foreground">
                {formatVotes(results.yes)}
              </span>
            </div>
            <div className="h-2 rounded-full bg-muted overflow-hidden">
              <div
                className="h-full rounded-full bg-success transition-all"
                style={{
                  width: `${(results.yes / results.total) * 100}%`,
                }}
              />
            </div>
          </div>
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-destructive font-medium">No</span>
              <span className="text-muted-foreground">
                {formatVotes(results.no)}
              </span>
            </div>
            <div className="h-2 rounded-full bg-muted overflow-hidden">
              <div
                className="h-full rounded-full bg-destructive transition-all"
                style={{
                  width: `${(results.no / results.total) * 100}%`,
                }}
              />
            </div>
          </div>
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-muted-foreground font-medium">Abstain</span>
              <span className="text-muted-foreground">
                {formatVotes(results.abstain)}
              </span>
            </div>
            <div className="h-2 rounded-full bg-muted overflow-hidden">
              <div
                className="h-full rounded-full bg-muted-foreground transition-all"
                style={{
                  width: `${(results.abstain / results.total) * 100}%`,
                }}
              />
            </div>
          </div>
          <div className="text-sm text-muted-foreground pt-2 border-t">
            Total voting power: {formatVotes(results.total)}
          </div>
        </div>
      </div>
    </div>
  )
}

function formatVotes(n: number): string {
  if (n >= 1000000) return `${(n / 1000000).toFixed(1)}M`
  if (n >= 1000) return `${(n / 1000).toFixed(1)}K`
  return n.toString()
}
