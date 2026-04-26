<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { ArrowLeft, LoaderCircle, MessageSquare, Plus, Search, Send, X } from 'lucide-vue-next'
import { useChatWebSocket, useChatConversations, useChatMessages, useChatSearch } from '../composables/chat'
import {
  getConversationId,
  getInitials,
  getConversationName,
  isFirstInGroup,
  isLastInGroup,
  shouldShowSenderName,
  getSenderName,
  parseJwtPayload,
  getAuthToken,
} from '../helpers/chatHelpers'

const { isConnected, connect, sendMessage: wsSendMessage, onMessage } = useChatWebSocket()
const { conversations, selectedConversation, loading, fetchConversations, selectConversation, createDirectConversation, getSelectedConversationId } = useChatConversations()
const { messages, loadingMessages, loadingMore, hasMore, loadMessages, loadOlderMessages, addMessage } = useChatMessages()
const { searchQuery, searchResults, showNewConversation, toggleNewConversation, clearSearch } = useChatSearch()

const input = ref('')
const currentUserId = ref('')
const messagesContainer = ref(null)
const messagesEnd = ref(null)

const selectedConversationId = computed(() => getSelectedConversationId())
const visibleConversations = computed(() => conversations.value.filter(Boolean))
const visibleSearchResults = computed(() => searchResults.value.filter(Boolean))

function scrollToBottom(behavior = 'smooth') {
  messagesEnd.value?.scrollIntoView({ behavior })
}

function handleSendMessage() {
  const content = input.value.trim()
  if (!content || !selectedConversationId.value) return

  const sent = wsSendMessage(selectedConversationId.value, content)
  if (sent) {
    input.value = ''
  }
}

function handleKeydown(event) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSendMessage()
  }
}

async function handleLoadMessages(conversation) {
  selectConversation(conversation)
  await loadMessages(getConversationId(conversation))
  await nextTick()
  scrollToBottom('auto')
}

async function handleStartConversation(userId) {
  const conversation = await createDirectConversation(userId)
  if (conversation) {
    clearSearch()
    await handleLoadMessages(conversation)
  }
}

function handleScroll() {
  loadOlderMessages(selectedConversationId.value, messagesContainer.value)
}

watch(messages, async () => {
  await nextTick()
  const container = messagesContainer.value
  if (!container) return
  const isNearBottom = container.scrollHeight - container.scrollTop - container.clientHeight < 120
  if (isNearBottom) scrollToBottom()
})

onMounted(async () => {
  const token = getAuthToken()
  const payload = token ? parseJwtPayload(token) : null
  currentUserId.value = payload?.user_id || ''
  
  connect()
  onMessage((message) => {
    if (message.conversation_id === selectedConversationId.value) {
      addMessage(message)
    }
  })
  
  await fetchConversations()
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
            type="button"
            class="conversation-item"
            @click="handleStartConversation(user.user_id)"
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
          type="button"
          class="conversation-item"
          :class="{ 'conversation-item--active': selectedConversationId === getConversationId(conversation) }"
          @click="handleLoadMessages(conversation)"
        >
          <div class="avatar">{{ getInitials(getConversationName(conversation, currentUserId)) }}</div>
          <div class="conversation-copy">
            <p>{{ getConversationName(conversation, currentUserId) }}</p>
            <span>{{ conversation.type === 'direct' ? 'Trò chuyện trực tiếp' : `Nhóm ${conversation.participants?.length || 0} thành viên` }}</span>
          </div>
        </button>
      </div>
    </aside>

    <section class="chat-area" :class="{ 'chat-area--active': selectedConversation }">
      <template v-if="selectedConversation">
        <header class="chat-area__header">
          <button type="button" class="back-button" @click="selectConversation(null)">
            <ArrowLeft :size="24" />
          </button>
          <div class="avatar avatar--primary">{{ getInitials(getConversationName(selectedConversation, currentUserId)) }}</div>
          <div class="conversation-copy">
            <p>{{ getConversationName(selectedConversation, currentUserId) }}</p>
            <span>{{ selectedConversation.type === 'direct' ? 'Đang trực tuyến' : `${selectedConversation.participants?.length || 0} thành viên` }}</span>
          </div>
        </header>

        <div ref="messagesContainer" class="message-list" @scroll="handleScroll">
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
                'message-bubble--first': isFirstInGroup(messages, index),
                'message-bubble--last': isLastInGroup(messages, index),
              }"
            >
              <span v-if="shouldShowSenderName(message, index, messages, selectedConversation?.type, currentUserId)" class="message-sender">
                {{ getSenderName(message) }}
              </span>
              <p>{{ message.content }}</p>
              <span>{{ new Date(message.created_at).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) }}</span>
            </div>
          </div>
          <div ref="messagesEnd"></div>
        </div>

        <footer class="composer">
          <button type="button" class="attach-button">
            <Plus :size="24" />
          </button>
          <textarea
            v-model="input"
            rows="1"
            placeholder="Nhắn tin..."
            @keydown="handleKeydown"
          ></textarea>
          <button type="button" class="send-button" :disabled="!input.trim()" @click="handleSendMessage">
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
