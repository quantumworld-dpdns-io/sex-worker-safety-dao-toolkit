"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  FileText,
  BookOpen,
  Bell,
  AlertTriangle,
  Vote,
  CheckCircle,
} from "lucide-react"
import { cn } from "@/lib/utils"

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/report", label: "Report Client", icon: FileText },
  { href: "/registry", label: "Client Registry", icon: BookOpen },
  { href: "/checkin", label: "Check-In", icon: CheckCircle },
  { href: "/emergency", label: "Emergency", icon: AlertTriangle },
  { href: "/dao", label: "DAO", icon: Vote },
]

export default function Sidebar() {
  const pathname = usePathname()

  if (pathname === "/login" || pathname === "/") return null

  return (
    <aside className="hidden md:flex w-56 flex-col border-r bg-card">
      <nav className="flex-1 space-y-1 p-3">
        {navItems.map((item) => {
          const Icon = item.icon
          const active = pathname === item.href || pathname.startsWith(item.href + "/")
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                active
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
              )}
            >
              <Icon className="h-4 w-4" />
              {item.label}
            </Link>
          )
        })}
      </nav>
      <div className="border-t p-4">
        <p className="text-xs text-muted-foreground text-center">
          Safety DAO v0.1
        </p>
      </div>
    </aside>
  )
}
