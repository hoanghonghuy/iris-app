export const POST_TYPE_META = {
  announcement: { label: 'Thông báo', badgeClass: 'badge--primary' },
  activity: { label: 'Hoạt động', badgeClass: 'badge--info' },
  daily_note: { label: 'Nhận xét ngày', badgeClass: 'badge--outline' },
  health_note: { label: 'Sức khỏe', badgeClass: 'badge--danger' },
}

export const POST_TYPE_OPTIONS = Object.entries(POST_TYPE_META).map(([value, meta]) => ({
  value,
  label: meta.label,
}))

export const POST_SCOPE_LABELS = {
  school: 'Toàn trường',
  class: 'Cả lớp',
  student: 'Từng học sinh',
}
