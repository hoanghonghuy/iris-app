export function getConversationId(conversation) {
  return conversation?.conversation_id || conversation?.id || ''
}

export function getInitials(name) {
  return (name || '?')
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase()
}

export function getConversationName(conversation, currentUserId) {
  if (conversation?.name) return conversation.name
  const other = conversation?.participants?.find(
    (participant) => participant.user_id !== currentUserId,
  )
  return other?.full_name || other?.email || 'Cuộc hội thoại'
}

export function isFirstInGroup(messages, index) {
  const previous = messages[index - 1]
  const current = messages[index]
  return !previous || previous.sender_id !== current?.sender_id
}

export function isLastInGroup(messages, index) {
  const next = messages[index + 1]
  const current = messages[index]
  return !next || next.sender_id !== current?.sender_id
}

export function shouldShowSenderName(message, index, messages, conversationType, currentUserId) {
  return (
    conversationType !== 'direct' &&
    message?.sender_id !== currentUserId &&
    isFirstInGroup(messages, index)
  )
}

export function getSenderName(message) {
  return message?.sender_name || message?.sender_email?.split('@')[0] || 'Người gửi'
}

export function parseJwtPayload(token) {
  try {
    const [, payload] = token.split('.')
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return null
  }
}

export function getAuthToken() {
  return sessionStorage.getItem('auth_token')
}
