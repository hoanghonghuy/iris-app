import { tokenStorage } from '@/helpers/auth'
import type { Conversation, Message, ConversationParticipant } from '@/types'

export function getConversationId(conversation: Conversation | any): string {
  return conversation?.conversation_id || conversation?.id || ''
}

export function getInitials(name: string | null | undefined): string {
  return (name || '?')
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase()
}

export function getConversationName(conversation: Conversation | any, currentUserId: string): string {
  if (conversation?.name) return conversation.name
  const other = conversation?.participants?.find(
    (participant: ConversationParticipant) => participant.user_id !== currentUserId,
  )
  return other?.full_name || other?.user_name || other?.email || 'Cuộc hội thoại'
}

/** Dòng phụ trong danh sách hội thoại: tin cuối (rút gọn) hoặc mô tả mặc định. */
export function getConversationListSubtitle(conversation: Conversation | any, currentUserId: string): string {
  const lm = conversation?.last_message
  if (!lm?.content) {
    return conversation?.type === 'direct'
      ? 'Trò chuyện trực tiếp'
      : `Nhóm ${conversation?.participants?.length || 0} thành viên`
  }
  const prefix = String(lm.sender_id) === String(currentUserId) ? 'Bạn: ' : ''
  const text = String(lm.content).replace(/\s+/g, ' ').trim()
  const max = 72
  const short = text.length > max ? `${text.slice(0, max)}…` : text
  return `${prefix}${short}`
}

export function isFirstInGroup(messages: Message[], index: number): boolean {
  const previous = messages[index - 1]
  const current = messages[index]
  return !previous || previous.sender_id !== current?.sender_id
}

export function isLastInGroup(messages: Message[], index: number): boolean {
  const next = messages[index + 1]
  const current = messages[index]
  return !next || next.sender_id !== current?.sender_id
}

export function shouldShowSenderName(
  message: Message,
  index: number,
  messages: Message[],
  conversationType: string,
  currentUserId: string,
): boolean {
  return (
    conversationType !== 'direct' &&
    message?.sender_id !== currentUserId &&
    isFirstInGroup(messages, index)
  )
}

export function getSenderName(message: Message | any): string {
  return message?.sender_name || message?.sender_email?.split('@')[0] || 'Người gửi'
}

export function parseJwtPayload(token: string): any | null {
  try {
    const [, payload] = token.split('.')
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return null
  }
}

export function getAuthToken(): string | null {
  return tokenStorage.getToken()
}
