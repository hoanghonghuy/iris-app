export function normalizeListResponse(value) {
  if (Array.isArray(value?.data)) return value.data.filter(Boolean)
  if (Array.isArray(value)) return value.filter(Boolean)
  return []
}

export function normalizePaginatedResponse(value, fallbackLimit = 100) {
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

export async function fetchAllPaginated(fetchPage, options = {}) {
  const { limit = 100, initialOffset = 0 } = options

  let offset = initialOffset
  let hasMore = true
  let total = 0
  const items = []

  while (hasMore) {
    const response = await fetchPage({ limit, offset })
    const normalized = normalizePaginatedResponse(response, limit)

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
