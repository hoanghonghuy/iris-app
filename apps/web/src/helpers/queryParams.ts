/**
 * Build query parameters for date range filtering
 */
export function buildDateRangeParams(from: string | null | undefined, to: string | null | undefined): Record<string, string> {
  const params: Record<string, string> = {}
  if (from) params.from = from
  if (to) params.to = to
  return params
}
