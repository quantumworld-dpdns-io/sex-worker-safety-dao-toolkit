"use client"

import { useState } from "react"
import { Loader2, Send } from "lucide-react"
import { api } from "@/lib/api"

const categories = [
  "Violence",
  "Theft",
  "Harassment",
  "Non-payment",
  "Boundary violation",
  "Other",
]

const regions = [
  "North America",
  "South America",
  "Europe",
  "Asia Pacific",
  "Africa",
  "Middle East",
  "Other",
]

interface AttestationFormProps {
  onSuccess?: () => void
}

export default function AttestationForm({ onSuccess }: AttestationFormProps) {
  const [category, setCategory] = useState("")
  const [region, setRegion] = useState("")
  const [details, setDetails] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!category || !details) {
      setError("Category and details are required")
      return
    }
    setLoading(true)
    setError("")
    try {
      await api.submitReport({
        reportCategory: category,
        locationRegion: region || undefined,
        encryptedDetails: btoa(details),
      })
      setSuccess(true)
      onSuccess?.()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Submission failed")
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="rounded-xl border bg-card p-8 text-center shadow-sm">
        <div className="h-12 w-12 rounded-full bg-success/20 flex items-center justify-center mx-auto mb-4">
          <Send className="h-6 w-6 text-success" />
        </div>
        <h3 className="text-lg font-semibold mb-2">Report Submitted</h3>
        <p className="text-sm text-muted-foreground mb-4">
          Your report has been encrypted and submitted anonymously.
        </p>
        <button
          onClick={() => {
            setSuccess(false)
            setCategory("")
            setRegion("")
            setDetails("")
          }}
          className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          Submit Another
        </button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="rounded-xl border bg-card p-6 shadow-sm space-y-4">
      <div>
        <label className="text-sm font-medium mb-1.5 block" htmlFor="category">
          Incident Category
        </label>
        <select
          id="category"
          value={category}
          onChange={(e) => setCategory(e.target.value)}
          className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="">Select category...</option>
          {categories.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="text-sm font-medium mb-1.5 block" htmlFor="region">
          Location Region
        </label>
        <select
          id="region"
          value={region}
          onChange={(e) => setRegion(e.target.value)}
          className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="">Prefer not to say</option>
          {regions.map((r) => (
            <option key={r} value={r}>
              {r}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="text-sm font-medium mb-1.5 block" htmlFor="details">
          Encrypted Details
        </label>
        <textarea
          id="details"
          value={details}
          onChange={(e) => setDetails(e.target.value)}
          rows={5}
          placeholder="Describe what happened (this will be encrypted before submission)..."
          className="w-full rounded-lg border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring resize-y min-h-[100px]"
        />
      </div>

      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
      >
        {loading ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <>
            <Send className="h-4 w-4" />
            Submit Anonymous Report
          </>
        )}
      </button>
    </form>
  )
}
