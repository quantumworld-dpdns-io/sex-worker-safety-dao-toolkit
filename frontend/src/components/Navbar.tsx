"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { Shield, LogOut, User, Moon, Sun } from "lucide-react"
import { cn } from "@/lib/utils"
import { getUser, logout } from "@/lib/auth"
import { useState, useEffect } from "react"

export default function Navbar() {
  const pathname = usePathname()
  const [dark, setDark] = useState(false)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    const stored = localStorage.getItem("sws_theme")
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches
    const isDark = stored ? stored === "dark" : prefersDark
    setDark(isDark)
    document.documentElement.classList.toggle("dark", isDark)
  }, [])

  function toggleTheme() {
    const next = !dark
    setDark(next)
    localStorage.setItem("sws_theme", next ? "dark" : "light")
    document.documentElement.classList.toggle("dark", next)
  }

  const user = mounted ? getUser() : null

  if (pathname === "/login" || pathname === "/") return null

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-14 items-center px-4 gap-4">
        <Link href="/dashboard" className="flex items-center gap-2 font-semibold">
          <Shield className="h-5 w-5 text-primary" />
          <span className="hidden sm:inline-block">Safety DAO</span>
        </Link>

        <div className="flex-1" />

        {mounted && (
          <button
            onClick={toggleTheme}
            className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors h-9 w-9 hover:bg-accent hover:text-accent-foreground"
            aria-label="Toggle theme"
          >
            {dark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </button>
        )}

        {user && (
          <div className="flex items-center gap-2">
            <div className="hidden sm:flex items-center gap-2 text-sm text-muted-foreground">
              <User className="h-4 w-4" />
              <span>{user.displayName || user.walletAddress.slice(0, 6)}</span>
            </div>
            <button
              onClick={() => { logout(); window.location.href = "/login" }}
              className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors h-9 w-9 hover:bg-accent hover:text-accent-foreground"
              aria-label="Logout"
            >
              <LogOut className="h-4 w-4" />
            </button>
          </div>
        )}
      </div>
    </header>
  )
}
