import type { HealthStatus } from '@/types'

interface HealthSeverityOption {
  value: HealthStatus
  label: string
  badge: string
}

export const HEALTH_SEVERITY_OPTIONS: HealthSeverityOption[] = [
  { value: 'normal', label: 'Bình thường', badge: 'badge badge--outline' },
  { value: 'watch', label: 'Theo dõi', badge: 'badge badge--outline text-warning' },
  { value: 'urgent', label: 'Khẩn cấp', badge: 'badge badge--danger' },
]

export function getSeverityLabel(value: HealthStatus | null | undefined): string {
  return HEALTH_SEVERITY_OPTIONS.find((option) => option.value === value)?.label || value || '-'
}

export function getSeverityBadge(value: HealthStatus | null | undefined): string {
  return (
    HEALTH_SEVERITY_OPTIONS.find((option) => option.value === value)?.badge ||
    'badge badge--outline'
  )
}
