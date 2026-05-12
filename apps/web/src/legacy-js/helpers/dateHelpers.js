export function daysAgo(count) {
  const date = new Date()
  date.setDate(date.getDate() - count)
  return date
}

export function offsetDate(days) {
  const date = new Date()
  date.setDate(date.getDate() + days)
  return date
}

export function getDateInputValue(date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

export function getTodayDateString() {
  return getDateInputValue(new Date())
}
