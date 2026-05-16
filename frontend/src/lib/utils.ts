import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import { format, formatDistanceToNow, parseISO } from "date-fns"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(dateString: string): string {
  try {
    return format(parseISO(dateString), "MMM d, yyyy h:mm a")
  } catch {
    return dateString
  }
}

export function formatDateShort(dateString: string): string {
  try {
    return format(parseISO(dateString), "MMM d, yyyy")
  } catch {
    return dateString
  }
}

export function timeAgo(dateString: string): string {
  try {
    return formatDistanceToNow(parseISO(dateString), { addSuffix: true })
  } catch {
    return dateString
  }
}

export function truncateAddress(address: string): string {
  if (!address || address.length < 10) return address
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}
