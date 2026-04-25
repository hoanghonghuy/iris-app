<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { ArrowLeft, LoaderCircle, MessageSquare, Plus, Search, Send, X } from 'lucide-vue-next'
import { chatService, getChatWsUrl } from '../services/chatService'

const conversations = ref([])
const selectedConversation = ref(null)
const messages = ref([])
const input = ref('')
const loading = ref(true)
const loadingMessages = ref(false)
const loadingMore = ref(false)
const currentUserId = ref('')
const searchQuery = ref('')
const searchResults = ref([])
const showNewConversation = ref(false)
const hasMore = ref(false)
const nextCursor = ref(null)
const isConnected = ref(false)
const messagesContainer = ref(null)
const messagesEnd = ref(null)
const websocket = ref(null)
let reconnectTimer = null
let searchTimer = null
let reconnectAttempts = 0

const selectedConversationId = computed(() => getConversationId(selectedConversation.value))
const visibleConversations = computed(() => conversations.value.filter(Boolean))
const visibleSearchResults = computed(() => searchResults.value.filter(Boolean))

function normalizeList(value) {
  const data = value?.data ?? value
  return Array.isArray(data) ? data.filter(Boolean) : []
}

function getConversationId(conversation) {
  return conversation?.conversation_id || conversation?.id || ''
}

function getToken() {
  return sessionStorage.getItem('auth_token')
}

function parseJwtPayload(token) {
  try {
    const [, payload] = token.split('.')
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return null
  }
}

function getInitials(name) {
  return (name || '?')
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase()
}

function getConversationName(conversation) {
  if (conversation?.name) return conversation.name
  const other = conversation?.participants?.find((participant) => participant.user_id !== currentUserId.value)
  return other?.full_name || other?.email || 'Cuộc hội thoại'
}

function isFirstInGroup(index) {
  const previous = messages.value[index - 1]
  const current = messages.value[index]
  return !previous || previous.sender_id !== current?.sender_id
}

function isLastInGroup(index) {
  const next = messages.value[index + 1]
  const current = messages.value[index]
  return !next || next.sender_id !== current?.sender_id
}

function shouldShowSenderName(message, index) {
  return selectedConversation.value?.type !== 'direct'
    && message?.sender_id !== currentUserId.value
    && isFirstInGroup(index)
}

function getSenderName(message) {
  return message?.sender_name || message?.sender_email?.split('@')[0] || 'Người gửi'
}

function connectWebSocket() {
  const token = getToken()
  if (!token) return

  if (websocket.value) {
    websocket.value.close(1000, 'reconnect')
  }

  const ws = new WebSocket(getChatWsUrl(), ['Bearer', token])

  ws.onopen = () => {
    isConnected.value = true
    reconnectAttempts = 0
  }

  ws.onmessage = (event) => {
    try {
      const wsEvent = JSON.parse(event.data)
      if (wsEvent.type !== 'new_message') return
      const message = wsEvent.data
      if (message.conversation_id !== selectedConversationId.value) return
      if (messages.value.some((item) => item.message_id === message.message_id)) return
      messages.value = [...messages.value, message]
    } catch (error) {
      console.error('[chat] cannot parse websocket message', error)
    }
  }

  ws.onclose = (event) => {
    isConnected.value = false
    if (event.code === 1000 || reconnectAttempts >= 5) return
    reconnectAttempts += 1
    reconnectTimer = setTimeout(connectWebSocket, Math.min(3000 * reconnectAttempts, 15000))
  }

  websocket.value = ws
}

async function fetchConversations() {
  loading.value = true
  try {
    conversations.value = normalizeList(await chatService.listConversations())
  } catch (error) {
    console.error('[chat] cannot load conversations', error)
    conversations.value = []
  } finally {
    loading.value = false
  }
}

async function loadMessages(conversation) {
  const conversationId = getConversationId(conversation)
  if (!conversationId) return

  selectedConversation.value = conversation
  messages.value = []
  nextCursor.value = null
  hasMore.value = false
  loadingMessages.value = true

  try {
    const response = await chatService.listMessages(conversationId)
    messages.value = normalizeList(response).reverse()
    nextCursor.value = response?.next_cursor ?? null
    hasMore.value = Boolean(response?.has_more)
    await nextTick()
    scrollToBottom('auto')
  } catch (error) {
    console.error('[chat] cannot load messages', error)
  } finally {
    loadingMessages.value = false
  }
}

async function loadOlderMessages() {
  const container = messagesContainer.value
  const conversationId = getConversationId(selectedConversation.value)
  if (!container || !conversationId || !hasMore.value || loadingMore.value || !nextCursor.value) return
  if (container.scrollTop >= 80) return

  loadingMore.value = true
  const previousHeight = container.scrollHeight
  try {
    const response = await chatService.listMessages(conversationId, 50, nextCursor.value)
    messages.value = normalizeList(response).reverse().concat(messages.value)
    nextCursor.value = response?.next_cursor ?? null
    hasMore.value = Boolean(response?.has_more)
    await nextTick()
    container.scrollTop = container.scrollHeight - previousHeight
  } catch (error) {
    console.error('[chat] cannot load older messages', error)
  } finally {
    loadingMore.value = false
  }
}

function scrollToBottom(behavior = 'smooth') {
  messagesEnd.value?.scrollIntoView({ behavior })
}

function sendMessage() {
  const content = input.value.trim()
  if (!content || !selectedConversationId.value) return
  if (websocket.value?.readyState !== WebSocket.OPEN) return

  websocket.value.send(
    JSON.stringify({
      conversation_id: selectedConversationId.value,
      content,
    }),
  )
  input.value = ''
}

function handleKeydown(event) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

function toggleNewConversation() {
  showNewConversation.value = !showNewConversation.value
  if (!showNewConversation.value) {
    searchQuery.value = ''
    searchResults.value = []
  }
}

async function startConversation(userId) {
  if (!userId) return
  try {
    const conversation = await chatService.createDirectConversation(userId)
    await fetchConversations()
    searchQuery.value = ''
    searchResults.value = []
    showNewConversation.value = false
    await loadMessages(conversation)
  } catch (error) {
    console.error('[chat] cannot start conversation', error)
  }
}

watch(searchQuery, (query) => {
  clearTimeout(searchTimer)
  if (!query.trim()) {
    searchResults.value = []
    return
  }

  searchTimer = setTimeout(async () => {
    try {
      searchResults.value = normalizeList(await chatService.searchUsers(query.trim()))
    } catch {
      searchResults.value = []
    }
  }, 500)
})

watch(messages, async () => {
  await nextTick()
  const container = messagesContainer.value
  if (!container) return
  const isNearBottom = container.scrollHeight - container.scrollTop - container.clientHeight < 120
  if (isNearBottom) scrollToBottom()
})

onMounted(async () => {
  const token = getToken()
  const payload = token ? parseJwtPayload(token) : null
  currentUserId.value = payload?.user_id || ''
  connectWebSocket()
  await fetchConversations()
})

onUnmounted(() => {
  clearTimeout(reconnectTimer)
  clearTimeout(searchTimer)
  websocket.value?.close(1000, 'component unmount')
})
</script>

<template>
  <div class="chat-page">
    <aside class="chat-sidebar" :class="{ 'chat-sidebar--hidden': selectedConversation }">
      <div class="chat-sidebar__header">
        <h1>Tin nhắn</h1>
        <div class="chat-header-actions">
          <div class="connection-pill" :class="isConnected ? 'connection-pill--online' : 'connection-pill--offline'">
            <span class="connection-dot"></span>
            <span>{{ isConnected ? 'Online' : 'Offline' }}</span>
          </div>
          <button class="new-chat-button" type="button" :title="showNewConversation ? 'Đóng' : 'Tạo mới'" @click="toggleNewConversation">
            <X v-if="showNewConversation" :size="20" />
            <Plus v-else :size="20" />
          </button>
        </div>
      </div>

      <div v-if="showNewConversation" class="new-conversation-panel">
        <div class="search-box">
          <Search :size="16" />
          <input v-model="searchQuery" placeholder="Tìm kiếm email hoặc tên..." />
        </div>

        <div v-if="visibleSearchResults.length > 0" class="search-results">
          <button
            v-for="user in visibleSearchResults"
            :key="user.user_id"
            class="conversation-item"
            @click="startConversation(user.user_id)"
          >
            <div class="avatar">{{ getInitials(user.full_name || user.email) }}</div>
            <div class="conversation-copy">
              <p>{{ user.full_name || user.email }}</p>
              <span>{{ user.email }}</span>
            </div>
          </button>
        </div>
      </div>

      <div v-if="loading" class="chat-loading">
        <LoaderCircle class="spin" :size="24" />
      </div>

      <div v-else-if="visibleConversations.length === 0" class="chat-empty">
        <MessageSquare :size="32" />
        <p>Chưa có tin nhắn nào</p>
      </div>

      <div v-else class="conversation-list">
        <button
          v-for="conversation in visibleConversations"
          :key="getConversationId(conversation)"
          class="conversation-item"
          :class="{ 'conversation-item--active': selectedConversationId === getConversationId(conversation) }"
          @click="loadMessages(conversation)"
        >
          <div class="avatar">{{ getInitials(getConversationName(conversation)) }}</div>
          <div class="conversation-copy">
            <p>{{ getConversationName(conversation) }}</p>
            <span>{{ conversation.type === 'direct' ? 'Trò chuyện trực tiếp' : `Nhóm ${conversation.participants?.length || 0} thành viên` }}</span>
          </div>
        </button>
      </div>
    </aside>

    <section class="chat-area" :class="{ 'chat-area--active': selectedConversation }">
      <template v-if="selectedConversation">
        <header class="chat-area__header">
          <button class="back-button" @click="selectedConversation = null">
            <ArrowLeft :size="24" />
          </button>
          <div class="avatar avatar--primary">{{ getInitials(getConversationName(selectedConversation)) }}</div>
          <div class="conversation-copy">
            <p>{{ getConversationName(selectedConversation) }}</p>
            <span>{{ selectedConversation.type === 'direct' ? 'Đang trực tuyến' : `${selectedConversation.participants?.length || 0} thành viên` }}</span>
          </div>
        </header>

        <div ref="messagesContainer" class="message-list" @scroll="loadOlderMessages">
          <div v-if="loadingMore" class="history-loading">
            <LoaderCircle class="spin" :size="16" />
            Đang tải lịch sử...
          </div>

          <div v-if="!hasMore && messages.length > 0" class="conversation-start">
            Khởi đầu cuộc trò chuyện
          </div>

          <div v-if="loadingMessages" class="chat-loading">
            <LoaderCircle class="spin" :size="24" />
          </div>

          <div v-else-if="messages.length === 0" class="message-empty">
            <MessageSquare :size="40" />
            <p>Gửi lời chào đầu tiên!</p>
          </div>

          <div v-else class="messages">
            <div
              v-for="(message, index) in messages"
              :key="message.message_id"
              class="message-bubble"
              :class="{
                'message-bubble--mine': message.sender_id === currentUserId,
                'message-bubble--first': isFirstInGroup(index),
                'message-bubble--last': isLastInGroup(index),
              }"
            >
              <span v-if="shouldShowSenderName(message, index)" class="message-sender">
                {{ getSenderName(message) }}
              </span>
              <p>{{ message.content }}</p>
              <span>{{ new Date(message.created_at).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) }}</span>
            </div>
          </div>
          <div ref="messagesEnd"></div>
        </div>

        <footer class="composer">
          <button class="attach-button">
            <Plus :size="24" />
          </button>
          <textarea
            v-model="input"
            rows="1"
            placeholder="Nhắn tin..."
            @keydown="handleKeydown"
          ></textarea>
          <button class="send-button" :disabled="!input.trim()" @click="sendMessage">
            <Send :size="20" />
          </button>
        </footer>
      </template>

      <div v-else class="chat-placeholder">
        <div class="placeholder-icon">
          <MessageSquare :size="40" />
        </div>
        <h2>Ứng dụng Nhắn tin IRIS</h2>
        <p>Chọn một cuộc trò chuyện để bắt đầu</p>
      </div>
    </section>
  </div>
</template>

<style scoped>
.chat-page {
  position: relative;
  display: flex;
  height: calc(100dvh - var(--header-height) - 2rem);
  min-height: 520px;
  width: 100%;
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-background);
}

.chat-sidebar {
  width: 100%;
  max-width: 380px;
  border-right: 1px solid var(--color-border);
  background: var(--color-surface);
  display: flex;
  flex-direction: column;
}

.chat-sidebar__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
  border-bottom: 1px solid var(--color-border);
}

.chat-sidebar__header h1,
.conversation-copy p,
.conversation-copy span,
.message-bubble p,
.message-bubble span,
.chat-placeholder h2,
.chat-placeholder p {
  margin: 0;
}

.chat-sidebar__header h1 {
  font-size: var(--font-size-2xl);
}

.chat-header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}

.connection-pill {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  border-radius: var(--radius-full);
  padding: var(--spacing-1) var(--spacing-2);
  background: var(--color-background);
  color: var(--color-text-muted);
  font-size: 0.65rem;
  font-weight: 800;
  letter-spacing: 0;
  text-transform: uppercase;
}

.connection-pill--online {
  color: var(--color-success);
  background: color-mix(in srgb, var(--color-success) 12%, transparent);
}

.connection-dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: var(--radius-full);
  background: currentColor;
}

.new-chat-button {
  width: 2.25rem;
  height: 2.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: var(--radius-full);
  background: color-mix(in srgb, var(--color-primary) 12%, transparent);
  color: var(--color-primary);
}

.new-chat-button:hover {
  background: color-mix(in srgb, var(--color-primary) 20%, transparent);
}

.new-conversation-panel {
  padding: var(--spacing-3);
  border-bottom: 1px solid var(--color-border);
}

.search-box {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-full);
  border: 1px solid var(--color-border);
  color: var(--color-text-muted);
}

.search-box input {
  min-width: 0;
  flex: 1;
  border: 0;
  background: transparent;
  color: var(--color-text);
  outline: none;
}

.search-results {
  overflow: hidden;
  margin-top: var(--spacing-2);
  border: 1px solid var(--color-border);
  border-radius: 1rem;
  background: var(--color-surface);
}

.conversation-list {
  overflow-y: auto;
  padding: var(--spacing-2);
}

.conversation-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  border: 0;
  border-radius: 1rem;
  background: transparent;
  color: var(--color-text);
  padding: var(--spacing-3);
  text-align: left;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background: var(--color-background);
}

.conversation-item--active {
  background: var(--color-primary);
  color: var(--color-on-primary);
}

.conversation-copy {
  min-width: 0;
  flex: 1;
}

.conversation-copy p {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 700;
  font-size: var(--font-size-sm);
}

.conversation-copy span {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.conversation-item--active .conversation-copy span {
  color: var(--color-on-primary-muted);
}

.avatar {
  width: 3rem;
  height: 3rem;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  background: var(--color-background);
  color: var(--color-text-muted);
  font-weight: 800;
}

.avatar--primary {
  width: 2.5rem;
  height: 2.5rem;
  background: var(--color-primary);
  color: var(--color-on-primary);
}

.chat-area {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: var(--color-background);
}

.chat-area__header {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-3);
  border-bottom: 1px solid var(--color-border);
  background: color-mix(in srgb, var(--color-background) 82%, transparent);
  backdrop-filter: blur(12px);
}

.back-button {
  display: none;
  border: 0;
  background: transparent;
  color: var(--color-primary);
  border-radius: var(--radius-full);
  padding: var(--spacing-2);
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-4);
}

.messages {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.message-bubble {
  max-width: min(78%, 36rem);
  align-self: flex-start;
  border-radius: 1.25rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  padding: var(--spacing-2) var(--spacing-3);
}

.message-bubble--first {
  margin-top: var(--spacing-3);
}

.message-bubble--mine {
  align-self: flex-end;
  background: var(--color-primary);
  color: var(--color-on-primary);
  border-color: var(--color-primary);
}

.message-sender {
  display: block;
  margin-bottom: var(--spacing-1);
  color: var(--color-primary);
  font-size: 0.68rem;
  font-weight: 800;
}

.message-bubble p {
  white-space: pre-line;
  font-size: var(--font-size-sm);
}

.message-bubble span {
  display: block;
  margin-top: var(--spacing-1);
  font-size: 0.68rem;
  opacity: 0.72;
}

.composer {
  display: flex;
  align-items: flex-end;
  gap: var(--spacing-2);
  padding: var(--spacing-3);
  border-top: 1px solid var(--color-border);
  background: var(--color-surface);
}

.composer textarea {
  min-height: 2.75rem;
  max-height: 7.5rem;
  flex: 1;
  resize: none;
  border: 1px solid transparent;
  border-radius: 1.4rem;
  background: var(--color-background);
  color: var(--color-text);
  padding: var(--spacing-3) var(--spacing-4);
  outline: none;
}

.composer textarea:focus {
  border-color: var(--color-primary);
  background: var(--color-surface);
}

.attach-button,
.send-button {
  width: 2.75rem;
  height: 2.75rem;
  border: 0;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
}

.attach-button {
  background: transparent;
  color: var(--color-text-muted);
}

.send-button {
  background: var(--color-primary);
  color: var(--color-on-primary);
}

.send-button:disabled {
  background: var(--color-background);
  color: var(--color-text-muted);
}

.chat-loading,
.chat-empty,
.message-empty,
.chat-placeholder,
.history-loading {
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-loading,
.chat-empty,
.message-empty,
.chat-placeholder {
  flex-direction: column;
  gap: var(--spacing-3);
  color: var(--color-text-muted);
  padding: 2rem;
  text-align: center;
}

.chat-placeholder {
  height: 100%;
}

.placeholder-icon {
  width: 6rem;
  height: 6rem;
  border-radius: var(--radius-full);
  background: var(--color-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-primary);
}

.history-loading {
  gap: var(--spacing-2);
  margin: 0 auto var(--spacing-3);
  width: fit-content;
  border-radius: var(--radius-full);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  padding: var(--spacing-1) var(--spacing-3);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.conversation-start {
  width: fit-content;
  margin: 0 auto var(--spacing-4);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  color: var(--color-text-muted);
  padding: var(--spacing-1) var(--spacing-3);
  font-size: 0.68rem;
  font-weight: 700;
  text-transform: uppercase;
}

.spin {
  animation: spin 1s linear infinite;
}

@media (max-width: 767px) {
  .chat-page {
    height: calc(100dvh - var(--header-height) - 1rem);
    border-radius: 0;
    border-left: 0;
    border-right: 0;
  }

  .chat-sidebar {
    max-width: none;
    transition: transform 0.25s ease;
  }

  .chat-sidebar--hidden {
    transform: translateX(-100%);
  }

  .chat-area {
    position: absolute;
    inset: 0;
    transform: translateX(100%);
    transition: transform 0.25s ease;
  }

  .chat-area--active {
    transform: translateX(0);
  }

  .back-button {
    display: flex;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
