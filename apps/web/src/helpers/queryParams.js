/**
 * Build query parameters for date range filtering
 * @param {string} from - Start date (ISO format)
 * @param {string} to - End date (ISO format)
 * @returns {Object} Query parameters object
 */
export function buildDateRangeParams(from, to) {
  const params = {}
  if (from) params.from = from
  if (to) params.to = to
  return params
}
