"use client"

import Link from "next/link"
import { Shield, FileText, Bell, Vote, ArrowRight, CheckCircle } from "lucide-react"

const features = [
  {
    icon: Shield,
    title: "Anonymous Reporting",
    description: "Report bad actors with zero-knowledge proofs. Your identity is never revealed.",
  },
  {
    icon: Bell,
    title: "Safety Check-Ins",
    description: "Scheduled check-ins with automatic escalation if you don't respond.",
  },
  {
    icon: FileText,
    title: "Client Registry",
    description: "Privacy-preserving registry of verified bad actor reports.",
  },
  {
    icon: Vote,
    title: "DAO Governance",
    description: "Community-driven moderation and protocol decisions.",
  },
]

export default function Home() {
  return (
    <div className="flex flex-col min-h-[calc(100vh-4rem)]">
      <section className="flex-1 flex flex-col items-center justify-center px-4 py-20 text-center">
        <div className="max-w-3xl mx-auto">
          <div className="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 mb-6">
            <Shield className="h-7 w-7 text-primary" />
          </div>
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold tracking-tight mb-6">
            Safety Through{" "}
            <span className="text-primary">Decentralization</span>
          </h1>
          <p className="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto">
            A privacy-first toolkit empowering sex workers with anonymous
            reporting, safety check-ins, and community-driven moderation.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Link
              href="/login"
              className="inline-flex items-center justify-center gap-2 rounded-lg bg-primary px-6 py-3 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
            >
              Get Started
              <ArrowRight className="h-4 w-4" />
            </Link>
            <Link
              href="/login"
              className="inline-flex items-center justify-center gap-2 rounded-lg border bg-background px-6 py-3 text-sm font-medium hover:bg-accent transition-colors"
            >
              <Shield className="h-4 w-4" />
              Connect Wallet
            </Link>
          </div>
        </div>
      </section>

      <section className="border-t py-16 px-4">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-2xl font-bold text-center mb-12">
            Everything you need to stay safe
          </h2>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature) => {
              const Icon = feature.icon
              return (
                <div
                  key={feature.title}
                  className="rounded-xl border bg-card p-6 hover:shadow-md transition-shadow"
                >
                  <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                    <Icon className="h-5 w-5 text-primary" />
                  </div>
                  <h3 className="font-semibold mb-2">{feature.title}</h3>
                  <p className="text-sm text-muted-foreground">
                    {feature.description}
                  </p>
                </div>
              )
            })}
          </div>
        </div>
      </section>

      <section className="border-t py-16 px-4 bg-muted/30">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-2xl font-bold mb-4">Powered by Zero-Knowledge Proofs</h2>
          <p className="text-muted-foreground mb-8 max-w-2xl mx-auto">
            All reports are encrypted and verified using zk-SNARKs, ensuring
            that your identity remains private while maintaining the integrity
            of the registry.
          </p>
          <div className="flex justify-center gap-8 flex-wrap">
            {["Noir", "RISC Zero", "Anon Aadhaar"].map((tech) => (
              <div key={tech} className="flex items-center gap-2 text-sm">
                <CheckCircle className="h-4 w-4 text-success" />
                <span className="font-medium">{tech}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      <footer className="border-t py-6 px-4 text-center text-sm text-muted-foreground">
        <p>Sex Worker Safety DAO &mdash; Building a safer industry, together.</p>
      </footer>
    </div>
  )
}
