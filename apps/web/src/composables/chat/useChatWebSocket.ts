import { ref, onUnmounted } from 'vue'
import { getChatWsUrl } from '../../services/chatService'
import { getAuthToken } from '../../helpers/chatHelpers'
import type { Message } from '@/types'

interface ChatWebSocketEvent {
  type: string
  data?: unknown
}

function isNewMessageEvent(event: ChatWebSocketEvent): event is { type: 'new_message'; data: Message } {
  return event.type === 'new_message' && Boolean(event.data)
}

export function useChatWebSocket() {
  const websocket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let reconnectAttempts = 0

  function connect() {
    const token = getAuthToken()
    if (!token) return

    if (websocket.value) {
      websocket.value.close(1000, 'reconnect')
    }

    const ws = new WebSocket(getChatWsUrl(), ['Bearer', token])

    ws.onopen = () => {
      isConnected.value = true
      reconnectAttempts = 0
    }

    ws.onclose = (event) => {
      isConnected.value = false
      if (event.code === 1000 || reconnectAttempts >= 5) return
      reconnectAttempts += 1
      reconnectTimer = setTimeout(connect, Math.min(3000 * reconnectAttempts, 15000))
    }

    websocket.value = ws
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    websocket.value?.close(1000, 'component unmount')
  }

  function sendMessage(conversationId: string, content: string) {
    if (!content || !conversationId) return false
    if (websocket.value?.readyState !== WebSocket.OPEN) return false

    websocket.value.send(
      JSON.stringify({
        conversation_id: conversationId,
        content,
      }),
    )
    return true
  }

  function onMessage(callback: (message: Message) => void) {
    if (!websocket.value) return

    const originalOnMessage = websocket.value.onmessage

    websocket.value.onmessage = (event) => {
      if (originalOnMessage) originalOnMessage(event)

      try {
        const wsEvent = JSON.parse(event.data) as ChatWebSocketEvent
        if (isNewMessageEvent(wsEvent)) {
          callback(wsEvent.data)
        }
      } catch (error) {
        console.error('[chat] cannot parse websocket message', error)
      }
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    websocket,
    isConnected,
    connect,
    disconnect,
    sendMessage,
    onMessage,
  }
}
