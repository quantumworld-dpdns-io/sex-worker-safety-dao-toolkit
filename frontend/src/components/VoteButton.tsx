"use client"

import { useState } from "react"
import { ThumbsUp, ThumbsDown, Minus, Loader2 } from "lucide-react"
import { cn } from "@/lib/utils"

interface VoteButtonProps {
  proposalId: string
  onVote?: (proposalId: string, vote: "yes" | "no" | "abstain") => Promise<void>
  disabled?: boolean
}

export default function VoteButton({ proposalId, onVote, disabled }: VoteButtonProps) {
  const [selected, setSelected] = useState<"yes" | "no" | "abstain" | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleVote(vote: "yes" | "no" | "abstain") {
    setLoading(true)
    setSelected(vote)
    try {
      await onVote?.(proposalId, vote)
    } finally {
      setLoading(false)
    }
  }

  const options = [
    { value: "yes" as const, icon: ThumbsUp, label: "Yes", activeClass: "bg-success/10 text-success border-success/30" },
    { value: "no" as const, icon: ThumbsDown, label: "No", activeClass: "bg-destructive/10 text-destructive border-destructive/30" },
    { value: "abstain" as const, icon: Minus, label: "Abstain", activeClass: "bg-muted text-muted-foreground border-border" },
  ]

  return (
    <div className="flex gap-2">
      {options.map((opt) => {
        const Icon = opt.icon
        const isSelected = selected === opt.value
        return (
          <button
            key={opt.value}
            onClick={() => handleVote(opt.value)}
            disabled={disabled || loading}
            className={cn(
              "flex items-center gap-1.5 rounded-lg border px-3 py-2 text-sm font-medium transition-colors",
              isSelected
                ? opt.activeClass
                : "border-border bg-background hover:bg-accent",
              disabled && "opacity-50 cursor-not-allowed",
            )}
          >
            {loading && isSelected ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Icon className="h-4 w-4" />
            )}
            {opt.label}
          </button>
        )
      })}
    </div>
  )
}
