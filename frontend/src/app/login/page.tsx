"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Shield, Loader2, Wallet, ArrowRight } from "lucide-react"
import { storeToken, storeUser } from "@/lib/auth"

export default function LoginPage() {
  const router = useRouter()
  const [address, setAddress] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  async function handleConnect(e: React.FormEvent) {
    e.preventDefault()
    if (!address.trim()) {
      setError("Please enter a wallet address")
      return
    }
    setLoading(true)
    setError("")

    // Simulated auth
    await new Promise((r) => setTimeout(r, 1000))

    const mockToken = `eyJhbGciOiJIUzI1NiJ9.${btoa(JSON.stringify({ sub: address, exp: Math.floor(Date.now() / 1000) + 86400 }))}.mock`
    storeToken(mockToken)
    storeUser({
      id: crypto.randomUUID(),
      walletAddress: address,
      role: "member",
      displayName: address.slice(0, 6),
    })
    router.push("/dashboard")
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <div className="inline-flex h-14 w-14 items-center justify-center rounded-2xl bg-primary/10 mb-4">
            <Shield className="h-7 w-7 text-primary" />
          </div>
          <h1 className="text-2xl font-bold">Welcome to Safety DAO</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Connect your wallet to continue
          </p>
        </div>

        <form onSubmit={handleConnect} className="rounded-xl border bg-card p-6 shadow-sm space-y-4">
          <div>
            <label className="text-sm font-medium mb-1.5 block" htmlFor="address">
              Wallet Address
            </label>
            <div className="relative">
              <Wallet className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <input
                id="address"
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                placeholder="0x..."
                className="w-full rounded-lg border bg-background pl-10 pr-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              />
            </div>
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
                Connect
                <ArrowRight className="h-4 w-4" />
              </>
            )}
          </button>

          <p className="text-xs text-center text-muted-foreground">
            This is a simulated wallet connection. In production, MetaMask or
            WalletConnect would be used.
          </p>
        </form>
      </div>
    </div>
  )
}
