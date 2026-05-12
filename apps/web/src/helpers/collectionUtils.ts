export function normalizeListResponse(value: any): any[] {
  if (Array.isArray(value?.data)) return value.data.filter(Boolean)
  if (Array.isArray(value)) return value.filter(Boolean)
  return []
}

interface Pagination {
  total: number
  limit: number
  offset: number
  has_more: boolean
}

interface PaginatedResult<T = any> {
  items: T[]
  pagination: Pagination
}

export function normalizePaginatedResponse<T = any>(
  value: any,
  fallbackLimit = 100,
): PaginatedResult<T> {
  const items = normalizeListResponse(value)

  return {
    items,
    pagination: value?.pagination ?? {
      total: items.length,
      limit: fallbackLimit,
      offset: 0,
      has_more: false,
    },
  }
}

interface FetchAllOptions {
  limit?: number
  initialOffset?: number
}

export async function fetchAllPaginated<T = any>(
  fetchPage: (params: { limit: number; offset: number }) => Promise<any>,
  options: FetchAllOptions = {},
): Promise<{ items: T[]; total: number }> {
  const { limit = 100, initialOffset = 0 } = options

  let offset = initialOffset
  let hasMore = true
  let total = 0
  const items: T[] = []

  while (hasMore) {
    const response = await fetchPage({ limit, offset })
    const normalized = normalizePaginatedResponse<T>(response, limit)

    items.push(...normalized.items)
    total = normalized.pagination.total ?? items.length

    hasMore = Boolean(normalized.pagination.has_more) && normalized.items.length > 0
    offset += normalized.pagination.limit || limit

    if (!hasMore || offset >= total) {
      break
    }
  }

  return { items, total: total || items.length }
}
