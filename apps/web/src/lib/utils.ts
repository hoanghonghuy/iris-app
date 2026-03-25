import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDateVN(input?: string | Date | null): string {
  if (!input) {
    return ""
  }

  if (typeof input === "string") {
    const dateMatch = input.match(/^(\d{4})-(\d{2})-(\d{2})/)
    if (dateMatch) {
      const year = Number(dateMatch[1])
      const month = Number(dateMatch[2]) - 1
      const day = Number(dateMatch[3])
      const parsed = new Date(year, month, day)
      return Number.isNaN(parsed.getTime()) ? input : parsed.toLocaleDateString("vi-VN")
    }
  }

  const parsed = input instanceof Date ? input : new Date(input)
  return Number.isNaN(parsed.getTime()) ? String(input) : parsed.toLocaleDateString("vi-VN")
}
