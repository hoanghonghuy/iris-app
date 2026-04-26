/**
 * Health severity configuration and utilities
 */

export const HEALTH_SEVERITY_OPTIONS = [
  { value: 'normal', label: 'Bình thường', badge: 'badge badge--outline' },
  { value: 'watch', label: 'Theo dõi', badge: 'badge badge--outline text-warning' },
  { value: 'urgent', label: 'Khẩn cấp', badge: 'badge badge--danger' },
]

export function getSeverityLabel(value) {
  return HEALTH_SEVERITY_OPTIONS.find((option) => option.value === value)?.label || value || '-'
}

export function getSeverityBadge(value) {
  return HEALTH_SEVERITY_OPTIONS.find((option) => option.value === value)?.badge || 'badge badge--outline'
}
