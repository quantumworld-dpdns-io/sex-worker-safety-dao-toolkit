"use client"

import { Shield } from "lucide-react"
import AttestationForm from "@/components/AttestationForm"

export default function ReportPage() {
  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div className="text-center">
        <div className="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 mb-4">
          <Shield className="h-6 w-6 text-primary" />
        </div>
        <h1 className="text-2xl font-bold">Anonymous Client Report</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Your report is encrypted and verified using zero-knowledge proofs.
          Your identity remains private.
        </p>
      </div>

      <AttestationForm />

      <div className="rounded-xl border bg-card p-5 shadow-sm">
        <h3 className="font-semibold mb-2">How it works</h3>
        <ul className="space-y-2 text-sm text-muted-foreground">
          <li className="flex items-start gap-2">
            <Shield className="h-4 w-4 text-primary mt-0.5 shrink-0" />
            <span>Your details are encrypted before submission</span>
          </li>
          <li className="flex items-start gap-2">
            <Shield className="h-4 w-4 text-primary mt-0.5 shrink-0" />
            <span>A zero-knowledge proof is generated to verify authenticity without revealing identity</span>
          </li>
          <li className="flex items-start gap-2">
            <Shield className="h-4 w-4 text-primary mt-0.5 shrink-0" />
            <span>Reports are reviewed by the DAO community for verification</span>
          </li>
          <li className="flex items-start gap-2">
            <Shield className="h-4 w-4 text-primary mt-0.5 shrink-0" />
            <span>Verified reports contribute to the shared registry</span>
          </li>
        </ul>
      </div>
    </div>
  )
}
