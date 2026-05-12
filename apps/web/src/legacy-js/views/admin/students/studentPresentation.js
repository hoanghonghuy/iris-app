export const GENDER_OPTIONS = [
  { value: 'male', label: 'Nam' },
  { value: 'female', label: 'Nữ' },
  { value: 'other', label: 'Khác' },
]

export function normalizeGender(gender) {
  const normalized = String(gender || '').toLowerCase()
  return normalized === 'female' || normalized === 'other' ? normalized : 'male'
}

export function getGenderLabel(gender) {
  const normalizedGender = normalizeGender(gender)
  return (
    GENDER_OPTIONS.find((option) => option.value === normalizedGender)?.label || normalizedGender
  )
}

export function getCodeExpiryText(dateString) {
  if (!dateString) {
    return ''
  }

  const diff = new Date(dateString).getTime() - Date.now()
  if (diff <= 0) {
    return 'Hết hạn'
  }

  return `Còn ${Math.ceil(diff / (1000 * 60 * 60 * 24))} ngày`
}
