export function daysAgo(count: number): Date {
  const date = new Date()
  date.setDate(date.getDate() - count)
  return date
}

export function offsetDate(days: number): Date {
  const date = new Date()
  date.setDate(date.getDate() + days)
  return date
}

export function getDateInputValue(date: Date): string {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

export function getTodayDateString(): string {
  return getDateInputValue(new Date())
}
