"use client"

import { useState } from "react"
import { Vote, Plus, Loader2 } from "lucide-react"
import ProposalCard from "@/components/ProposalCard"
import type { DAOProposal } from "@/types"

const mockProposals: DAOProposal[] = [
  {
    id: "1",
    title: "Increase verification threshold for bad client reports",
    description: "Proposal to increase the minimum confidence score from 0.7 to 0.85 for automatic verification of bad client reports to reduce false positives.",
    proposalType: "Parameter Change",
    status: "active",
    votingType: "Quadratic",
    startsAt: new Date(Date.now() - 86400000 * 3).toISOString(),
    endsAt: new Date(Date.now() + 86400000 * 4).toISOString(),
    createdBy: "0x1234...5678",
  },
  {
    id: "2",
    title: "Add two new moderator positions",
    description: "Due to growing community size, we need two additional moderators to help review incoming reports. Candidates will be vetted by existing moderators.",
    proposalType: "Governance",
    status: "active",
    votingType: "Token",
    startsAt: new Date(Date.now() - 86400000 * 2).toISOString(),
    endsAt: new Date(Date.now() + 86400000 * 5).toISOString(),
    createdBy: "0x8765...4321",
  },
  {
    id: "3",
    title: "Fund community safety workshop initiative",
    description: "Allocate 5000 USDC from the treasury to fund quarterly online safety workshops for community members.",
    proposalType: "Treasury",
    status: "passed",
    votingType: "Token",
    startsAt: new Date(Date.now() - 86400000 * 10).toISOString(),
    endsAt: new Date(Date.now() - 86400000 * 3).toISOString(),
    createdBy: "0xabcd...ef01",
  },
  {
    id: "4",
    title: "Update privacy policy for data retention",
    description: "Proposal to reduce data retention period from 2 years to 1 year for inactive accounts.",
    proposalType: "Policy",
    status: "rejected",
    votingType: "Quadratic",
    startsAt: new Date(Date.now() - 86400000 * 15).toISOString(),
    endsAt: new Date(Date.now() - 86400000 * 8).toISOString(),
    createdBy: "0x2468...1357",
  },
]

export default function DAOPage() {
  const [showCreate, setShowCreate] = useState(false)
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const [proposalType, setProposalType] = useState("Governance")
  const [votingType, setVotingType] = useState("Token")
  const [creating, setCreating] = useState(false)

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!title || !description) return
    setCreating(true)
    await new Promise((r) => setTimeout(r, 1000))
    setCreating(false)
    setShowCreate(false)
    setTitle("")
    setDescription("")
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">DAO Governance</h1>
          <p className="text-sm text-muted-foreground">
            Community-driven decision making
          </p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Proposal
        </button>
      </div>

      {showCreate && (
        <div className="rounded-xl border bg-card p-6 shadow-sm">
          <form onSubmit={handleCreate} className="space-y-4">
            <h3 className="font-semibold">Create Proposal</h3>
            <div>
              <label className="text-sm font-medium mb-1.5 block">Title</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Proposal title..."
                className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              />
            </div>
            <div>
              <label className="text-sm font-medium mb-1.5 block">Description</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={4}
                placeholder="Describe your proposal..."
                className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring resize-y"
              />
            </div>
            <div className="grid sm:grid-cols-2 gap-4">
              <div>
                <label className="text-sm font-medium mb-1.5 block">Type</label>
                <select
                  value={proposalType}
                  onChange={(e) => setProposalType(e.target.value)}
                  className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                >
                  <option>Governance</option>
                  <option>Treasury</option>
                  <option>Parameter Change</option>
                  <option>Policy</option>
                </select>
              </div>
              <div>
                <label className="text-sm font-medium mb-1.5 block">Voting Type</label>
                <select
                  value={votingType}
                  onChange={(e) => setVotingType(e.target.value)}
                  className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                >
                  <option>Token</option>
                  <option>Quadratic</option>
                </select>
              </div>
            </div>
            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => setShowCreate(false)}
                className="flex-1 rounded-lg border bg-background px-4 py-2 text-sm font-medium hover:bg-accent transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={creating || !title || !description}
                className="flex-1 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
              >
                {creating ? (
                  <Loader2 className="h-4 w-4 animate-spin mx-auto" />
                ) : (
                  "Submit Proposal"
                )}
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2">
        {mockProposals.map((p) => (
          <ProposalCard key={p.id} proposal={p} />
        ))}
      </div>
    </div>
  )
}
