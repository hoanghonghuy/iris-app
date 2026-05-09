<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { ArrowLeft, LoaderCircle, MessageSquare, Plus, Search, Send, Users, X } from 'lucide-vue-next'
import {
  useChatWebSocket,
  useChatConversations,
  useChatMessages,
  useChatSearch,
} from '../composables/chat'
import {
  getConversationId,
  getInitials,
  getConversationName,
  getConversationListSubtitle,
  isFirstInGroup,
  isLastInGroup,
  shouldShowSenderName,
  getSenderName,
  parseJwtPayload,
  getAuthToken,
} from '../helpers/chatHelpers'
import { extractErrorMessage } from '../helpers/errorHandler'
import { normalizeListResponse } from '../helpers/collectionUtils'
import { chatService } from '../services/chatService'

const { isConnected, connect, sendMessage: wsSendMessage, onMessage } = useChatWebSocket()
const {
  conversations,
  selectedConversation,
  loading,
  fetchConversations,
  selectConversation,
  createDirectConversation,
  createGroupConversation,
  renameGroup,
  addGroupParticipants,
  removeGroupParticipant,
  getSelectedConversationId,
  handleIncomingWsMessage,
  syncUnreadAfterOpen,
} = useChatConversations()
const {
  messages,
  loadingMessages,
  loadingMore,
  hasMore,
  loadMessages,
  loadOlderMessages,
  addMessage,
} = useChatMessages()
const { searchQuery, searchResults, showNewConversation, toggleNewConversation, clearSearch } =
  useChatSearch()

const newChatMode = ref('direct')
const groupName = ref('')
const groupMembers = ref([])
const groupSubmitting = ref(false)
const groupError = ref('')
const directError = ref('')

const showGroupManage = ref(false)
const groupRenameDraft = ref('')
const groupManageError = ref('')
const groupManageBusy = ref(false)
const manageSearchQuery = ref('')
const manageSearchResults = ref([])
let manageSearchTimer = null

const input = ref('')
const currentUserId = ref('')
const messagesContainer = ref(null)
const messagesEnd = ref(null)

const selectedConversationId = computed(() => getSelectedConversationId())
const visibleConversations = computed(() => conversations.value.filter(Boolean))
const visibleSearchResults = computed(() => searchResults.value.filter(Boolean))
const composerDisabled = computed(() => !selectedConversationId.value || !isConnected.value)
const composerHint = computed(() => {
  if (!selectedConversationId.value) return 'Chọn một cuộc trò chuyện để bắt đầu nhắn tin.'
  if (!isConnected.value) return 'Mất kết nối. Vui lòng chờ hệ thống tự kết nối lại.'
  return ''
})

const canRemoveGroupMembers = computed(
  () => (selectedConversation.value?.participants?.length || 0) > 2,
)

const visibleManageSearchResults = computed(() => {
  const ids = new Set(
    (selectedConversation.value?.participants || []).map((p) => p.user_id),
  )
  return manageSearchResults.value.filter((u) => u.user_id && !ids.has(u.user_id))
})

function scrollToBottom(behavior = 'smooth') {
  messagesEnd.value?.scrollIntoView({ behavior })
}

function handleSendMessage() {
  const content = input.value.trim()
  if (!content || composerDisabled.value) return

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
  const cid = getConversationId(conversation)
  await loadMessages(cid)
  syncUnreadAfterOpen(cid)
  await nextTick()
  scrollToBottom('auto')
}

async function handleStartConversation(userId) {
  directError.value = ''
  try {
    const conversation = await createDirectConversation(userId)
    if (conversation) {
      clearSearch()
      if (showNewConversation.value) toggleNewConversation()
      await handleLoadMessages(conversation)
    }
  } catch (err) {
    directError.value = extractErrorMessage(err)
  }
}

function addGroupMember(user) {
  if (!user?.user_id || user.user_id === currentUserId.value) return
  if (groupMembers.value.some((m) => m.user_id === user.user_id)) return
  groupMembers.value = [...groupMembers.value, user]
}

function removeGroupMember(userId) {
  groupMembers.value = groupMembers.value.filter((m) => m.user_id !== userId)
}

function toggleGroupManage() {
  if (selectedConversation.value?.type !== 'group') return
  showGroupManage.value = !showGroupManage.value
  if (showGroupManage.value) {
    groupManageError.value = ''
    manageSearchQuery.value = ''
    manageSearchResults.value = []
    groupRenameDraft.value = selectedConversation.value?.name ?? ''
  }
}

async function saveGroupName() {
  const id = selectedConversationId.value
  if (!id || selectedConversation.value?.type !== 'group') return
  groupManageError.value = ''
  groupManageBusy.value = true
  try {
    const updated = await renameGroup(id, groupRenameDraft.value)
    if (updated) selectConversation(updated)
  } catch (err) {
    groupManageError.value = extractErrorMessage(err)
  } finally {
    groupManageBusy.value = false
  }
}

async function addMemberFromManage(user) {
  const id = selectedConversationId.value
  if (!id || !user?.user_id) return
  groupManageError.value = ''
  groupManageBusy.value = true
  try {
    const updated = await addGroupParticipants(id, [user.user_id])
    if (updated) selectConversation(updated)
    manageSearchQuery.value = ''
    manageSearchResults.value = []
  } catch (err) {
    groupManageError.value = extractErrorMessage(err)
  } finally {
    groupManageBusy.value = false
  }
}

async function removeMemberFromManage(userId) {
  const id = selectedConversationId.value
  if (!id || !userId) return
  groupManageError.value = ''
  groupManageBusy.value = true
  try {
    const updated = await removeGroupParticipant(id, userId)
    if (updated) {
      selectConversation(updated)
    } else {
      showGroupManage.value = false
      selectConversation(null)
    }
  } catch (err) {
    groupManageError.value = extractErrorMessage(err)
  } finally {
    groupManageBusy.value = false
  }
}

async function handleCreateGroup() {
  groupError.value = ''
  if (groupMembers.value.length < 1) {
    groupError.value = 'Thêm ít nhất một thành viên để tạo nhóm.'
    return
  }
  groupSubmitting.value = true
  try {
    const conv = await createGroupConversation(
      groupName.value,
      groupMembers.value.map((m) => m.user_id),
    )
    if (conv) {
      clearSearch()
      if (showNewConversation.value) toggleNewConversation()
      await handleLoadMessages(conv)
    }
  } catch (err) {
    groupError.value = extractErrorMessage(err)
  } finally {
    groupSubmitting.value = false
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

watch(showNewConversation, (open) => {
  if (!open) {
    newChatMode.value = 'direct'
    groupName.value = ''
    groupMembers.value = []
    groupError.value = ''
    directError.value = ''
    groupSubmitting.value = false
  }
})

watch(newChatMode, () => {
  groupError.value = ''
  directError.value = ''
})

watch(manageSearchQuery, (q) => {
  clearTimeout(manageSearchTimer)
  if (!q?.trim()) {
    manageSearchResults.value = []
    return
  }
  manageSearchTimer = setTimeout(async () => {
    try {
      manageSearchResults.value = normalizeListResponse(
        await chatService.searchUsers(q.trim()),
      )
    } catch {
      manageSearchResults.value = []
    }
  }, 400)
})

watch(selectedConversation, (conv) => {
  if (!conv || conv.type !== 'group') {
    showGroupManage.value = false
    return
  }
  if (showGroupManage.value) {
    groupRenameDraft.value = conv.name ?? ''
  }
})

onMounted(async () => {
  const token = getAuthToken()
  const payload = token ? parseJwtPayload(token) : null
  currentUserId.value = payload?.user_id || ''

  connect()
  onMessage((message) => {
    handleIncomingWsMessage(message, currentUserId.value)
    if (message.conversation_id === selectedConversationId.value) {
      addMessage(message)
    }
  })

  await fetchConversations()
})

onUnmounted(() => {
  clearTimeout(manageSearchTimer)
})
</script>

<template>
  <div class="chat-page">
    <aside class="chat-sidebar" :class="{ 'chat-sidebar--hidden': selectedConversation }">
      <div class="chat-sidebar__header">
        <h1>Tin nhắn</h1>
        <div class="chat-header-actions">
          <div
            class="connection-pill"
            :class="isConnected ? 'connection-pill--online' : 'connection-pill--offline'"
          >
            <span class="connection-dot"></span>
            <span>{{ isConnected ? 'Online' : 'Offline' }}</span>
          </div>
          <button
            class="new-chat-button"
            type="button"
            :title="showNewConversation ? 'Đóng' : 'Tạo mới'"
            @click="toggleNewConversation"
          >
            <X v-if="showNewConversation" :size="20" />
            <Plus v-else :size="20" />
          </button>
        </div>
      </div>

      <div v-if="showNewConversation" class="new-conversation-panel">
        <div class="new-chat-tabs" role="tablist">
          <button
            type="button"
            class="new-chat-tab"
            role="tab"
            :aria-selected="newChatMode === 'direct'"
            @click="newChatMode = 'direct'"
          >
            Trực tiếp
          </button>
          <button
            type="button"
            class="new-chat-tab"
            role="tab"
            :aria-selected="newChatMode === 'group'"
            @click="newChatMode = 'group'"
          >
            Nhóm
          </button>
        </div>

        <div class="search-box">
          <Search :size="16" />
          <input v-model="searchQuery" placeholder="Tìm kiếm email hoặc tên..." />
        </div>

        <template v-if="newChatMode === 'direct'">
          <p v-if="directError" class="group-inline-error">{{ directError }}</p>
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
        </template>

        <template v-else>
          <input
            v-model="groupName"
            type="text"
            class="group-name-input"
            maxlength="255"
            placeholder="Tên nhóm (tuỳ chọn)"
            aria-label="Tên nhóm"
          />
          <div v-if="groupMembers.length > 0" class="member-chips">
            <span v-for="m in groupMembers" :key="m.user_id" class="member-chip">
              {{ m.full_name || m.email }}
              <button type="button" class="chip-remove" @click="removeGroupMember(m.user_id)">
                <X :size="14" />
              </button>
            </span>
          </div>
          <p v-if="groupError" class="group-inline-error">{{ groupError }}</p>
          <button
            type="button"
            class="create-group-submit"
            :disabled="groupSubmitting || groupMembers.length < 1"
            @click="handleCreateGroup"
          >
            <LoaderCircle v-if="groupSubmitting" class="spin" :size="18" />
            <span>{{ groupSubmitting ? 'Đang tạo...' : 'Tạo nhóm' }}</span>
          </button>
          <div v-if="visibleSearchResults.length > 0" class="search-results">
            <button
              v-for="user in visibleSearchResults"
              :key="user.user_id"
              type="button"
              class="conversation-item"
              :disabled="
                groupMembers.some((m) => m.user_id === user.user_id) || user.user_id === currentUserId
              "
              @click="addGroupMember(user)"
            >
              <div class="avatar">{{ getInitials(user.full_name || user.email) }}</div>
              <div class="conversation-copy">
                <p>{{ user.full_name || user.email }}</p>
                <span>{{
                  groupMembers.some((m) => m.user_id === user.user_id)
                    ? 'Đã thêm'
                    : 'Chạm để thêm vào nhóm'
                }}</span>
              </div>
            </button>
          </div>
        </template>
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
          :class="{
            'conversation-item--active': selectedConversationId === getConversationId(conversation),
          }"
          @click="handleLoadMessages(conversation)"
        >
          <div class="avatar">
            {{ getInitials(getConversationName(conversation, currentUserId)) }}
          </div>
          <div class="conversation-copy">
            <p>{{ getConversationName(conversation, currentUserId) }}</p>
            <span>{{ getConversationListSubtitle(conversation, currentUserId) }}</span>
          </div>
          <span
            v-if="(conversation.unread_count || 0) > 0"
            class="conversation-unread"
            :aria-label="`${conversation.unread_count} tin chưa đọc`"
          >
            {{ conversation.unread_count > 99 ? '99+' : conversation.unread_count }}
          </span>
        </button>
      </div>
    </aside>

    <section class="chat-area" :class="{ 'chat-area--active': selectedConversation }">
      <template v-if="selectedConversation">
        <header class="chat-area__header">
          <button type="button" class="back-button" @click="selectConversation(null)">
            <ArrowLeft :size="24" />
          </button>
          <div class="avatar avatar--primary">
            {{ getInitials(getConversationName(selectedConversation, currentUserId)) }}
          </div>
          <div class="conversation-copy">
            <p>{{ getConversationName(selectedConversation, currentUserId) }}</p>
            <span>{{
              selectedConversation.type === 'direct'
                ? 'Đang trực tuyến'
                : `${selectedConversation.participants?.length || 0} thành viên`
            }}</span>
          </div>
          <button
            v-if="selectedConversation.type === 'group'"
            type="button"
            class="group-manage-toggle"
            :class="{ 'group-manage-toggle--active': showGroupManage }"
            title="Quản lý nhóm"
            :aria-expanded="showGroupManage"
            @click="toggleGroupManage"
          >
            <Users :size="22" />
          </button>
        </header>

        <div
          v-if="showGroupManage && selectedConversation.type === 'group'"
          class="group-manage-panel"
        >
          <div class="group-manage-panel__head">
            <span class="group-manage-title">Quản lý nhóm</span>
            <button type="button" class="group-manage-close" @click="showGroupManage = false">
              <X :size="18" />
            </button>
          </div>

          <label class="group-manage-label" for="groupRenameInput">Tên nhóm</label>
          <div class="group-manage-row">
            <input
              id="groupRenameInput"
              v-model="groupRenameDraft"
              type="text"
              class="group-manage-input"
              maxlength="255"
              placeholder="Đặt tên hoặc để trống"
              :disabled="groupManageBusy"
            />
            <button
              type="button"
              class="btn-save-name"
              :disabled="groupManageBusy"
              @click="saveGroupName"
            >
              <LoaderCircle v-if="groupManageBusy" class="spin" :size="16" />
              <span v-else>Lưu</span>
            </button>
          </div>

          <p class="group-manage-hint">Để trống tên để hiển thị theo thành viên.</p>
          <p v-if="groupManageError" class="group-inline-error">{{ groupManageError }}</p>

          <p class="group-manage-section-title">Thành viên</p>
          <ul class="group-member-list">
            <li
              v-for="p in selectedConversation.participants"
              :key="p.user_id"
              class="group-member-row"
            >
              <span class="group-member-name">{{ p.full_name || p.email }}</span>
              <button
                type="button"
                class="group-member-remove"
                :disabled="!canRemoveGroupMembers || groupManageBusy"
                :title="
                  canRemoveGroupMembers ? 'Xóa khỏi nhóm' : 'Nhóm cần ít nhất 2 người — không thể xóa'
                "
                @click="removeMemberFromManage(p.user_id)"
              >
                <X :size="16" />
              </button>
            </li>
          </ul>

          <p class="group-manage-section-title">Thêm thành viên</p>
          <div class="search-box search-box--manage">
            <Search :size="16" />
            <input v-model="manageSearchQuery" placeholder="Tìm theo email hoặc tên..." />
          </div>
          <div v-if="visibleManageSearchResults.length > 0" class="search-results search-results--manage">
            <button
              v-for="user in visibleManageSearchResults"
              :key="user.user_id"
              type="button"
              class="conversation-item"
              :disabled="groupManageBusy"
              @click="addMemberFromManage(user)"
            >
              <div class="avatar">{{ getInitials(user.full_name || user.email) }}</div>
              <div class="conversation-copy">
                <p>{{ user.full_name || user.email }}</p>
                <span>{{ user.email }} · Chạm để thêm</span>
              </div>
            </button>
          </div>
        </div>

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
              <span
                v-if="
                  shouldShowSenderName(
                    message,
                    index,
                    messages,
                    selectedConversation?.type,
                    currentUserId,
                  )
                "
                class="message-sender"
              >
                {{ getSenderName(message) }}
              </span>
              <p>{{ message.content }}</p>
              <span>{{
                new Date(message.created_at).toLocaleTimeString('vi-VN', {
                  hour: '2-digit',
                  minute: '2-digit',
                })
              }}</span>
            </div>
          </div>
          <div ref="messagesEnd"></div>
        </div>

        <footer class="composer">
          <button type="button" class="attach-button" disabled title="Sắp hỗ trợ gửi tệp">
            <Plus :size="24" />
          </button>
          <textarea
            v-model="input"
            rows="1"
            placeholder="Nhắn tin..."
            :disabled="composerDisabled"
            @keydown="handleKeydown"
          ></textarea>
          <button
            type="button"
            class="send-button"
            :disabled="!input.trim() || composerDisabled"
            @click="handleSendMessage"
          >
            <Send :size="20" />
          </button>
        </footer>
        <p v-if="composerHint" class="composer-hint">{{ composerHint }}</p>
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

.new-chat-tabs {
  display: flex;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-3);
}

.new-chat-tab {
  flex: 1;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  padding: var(--spacing-2) var(--spacing-3);
  background: transparent;
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
  font-weight: 700;
  cursor: pointer;
}

.new-chat-tab[aria-selected='true'] {
  background: color-mix(in srgb, var(--color-primary) 14%, transparent);
  color: var(--color-primary);
  border-color: transparent;
}

.group-name-input {
  width: 100%;
  box-sizing: border-box;
  margin-top: var(--spacing-3);
  margin-bottom: var(--spacing-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  padding: var(--spacing-2) var(--spacing-3);
  font-size: var(--font-size-sm);
  background: var(--color-background);
  color: var(--color-text);
}

.member-chips {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-2);
}

.member-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-1);
  padding: var(--spacing-1) var(--spacing-2);
  border-radius: var(--radius-full);
  background: var(--color-background);
  font-size: var(--font-size-xs);
  font-weight: 600;
}

.chip-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  padding: 0;
  cursor: pointer;
  border-radius: var(--radius-full);
}

.chip-remove:hover {
  color: var(--color-text);
}

.group-inline-error {
  margin: 0 0 var(--spacing-2);
  font-size: var(--font-size-xs);
  color: var(--color-danger);
}

.create-group-submit {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-2);
  border: 0;
  border-radius: var(--radius-full);
  padding: var(--spacing-3);
  background: var(--color-primary);
  color: var(--color-on-primary);
  font-weight: 700;
  font-size: var(--font-size-sm);
  cursor: pointer;
}

.create-group-submit:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.conversation-item:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.search-box {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-full);
  border: 1px solid var(--color-border);
  color: var(--color-text-muted);
  margin-bottom: var(--spacing-3);
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

.conversation-unread {
  flex-shrink: 0;
  min-width: 1.25rem;
  padding: 0.125rem 0.45rem;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  color: var(--color-on-primary);
  font-size: var(--font-size-xs);
  font-weight: 700;
  line-height: 1.25;
}

.conversation-item--active .conversation-unread {
  background: var(--color-on-primary);
  color: var(--color-primary);
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

.chat-area__header .conversation-copy {
  flex: 1;
  min-width: 0;
}

.group-manage-toggle {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  background: var(--color-surface);
  color: var(--color-text-muted);
  cursor: pointer;
}

.group-manage-toggle--active {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: color-mix(in srgb, var(--color-primary) 10%, transparent);
}

.group-manage-panel {
  border-bottom: 1px solid var(--color-border);
  background: var(--color-surface);
  padding: var(--spacing-3) var(--spacing-4);
  max-height: min(50vh, 22rem);
  overflow-y: auto;
}

.group-manage-panel__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-3);
}

.group-manage-title {
  font-weight: 800;
  font-size: var(--font-size-sm);
}

.group-manage-close {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--spacing-1);
  border-radius: var(--radius-full);
}

.group-manage-close:hover {
  color: var(--color-text);
  background: var(--color-background);
}

.group-manage-label {
  display: block;
  font-size: var(--font-size-xs);
  font-weight: 700;
  margin-bottom: var(--spacing-1);
  color: var(--color-text-muted);
}

.group-manage-row {
  display: flex;
  gap: var(--spacing-2);
  align-items: center;
  margin-bottom: var(--spacing-1);
}

.group-manage-input {
  flex: 1;
  min-width: 0;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  padding: var(--spacing-2) var(--spacing-3);
  font-size: var(--font-size-sm);
  background: var(--color-background);
  color: var(--color-text);
}

.btn-save-name {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 4rem;
  border: 0;
  border-radius: var(--radius-full);
  padding: var(--spacing-2) var(--spacing-3);
  background: var(--color-primary);
  color: var(--color-on-primary);
  font-weight: 700;
  font-size: var(--font-size-xs);
  cursor: pointer;
}

.btn-save-name:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.group-manage-hint {
  margin: 0 0 var(--spacing-3);
  font-size: 0.68rem;
  color: var(--color-text-muted);
}

.group-manage-section-title {
  margin: 0 0 var(--spacing-2);
  font-size: var(--font-size-xs);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.02em;
  color: var(--color-text-muted);
}

.group-member-list {
  list-style: none;
  margin: 0 0 var(--spacing-3);
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.group-member-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-md);
  background: var(--color-background);
  font-size: var(--font-size-sm);
}

.group-member-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-member-remove {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--spacing-1);
  border-radius: var(--radius-full);
}

.group-member-remove:hover:not(:disabled) {
  color: var(--color-danger);
  background: var(--color-danger-soft-bg, color-mix(in srgb, var(--color-danger) 12%, transparent));
}

.group-member-remove:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.search-box--manage {
  margin-bottom: var(--spacing-2);
}

.search-results--manage {
  margin-top: 0;
  max-height: 10rem;
  overflow-y: auto;
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

.composer-hint {
  margin: 0;
  padding: 0 var(--spacing-3) var(--spacing-3);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
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

.attach-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
