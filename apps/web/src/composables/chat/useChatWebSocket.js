import { ref, onUnmounted } from 'vue'
import { getChatWsUrl } from '../../services/chatService'
import { getAuthToken } from '../../helpers/chatHelpers'

export function useChatWebSocket() {
  const websocket = ref(null)
  const isConnected = ref(false)
  let reconnectTimer = null
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
    clearTimeout(reconnectTimer)
    websocket.value?.close(1000, 'component unmount')
  }

  function sendMessage(conversationId, content) {
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

  function onMessage(callback) {
    if (!websocket.value) return

    const originalOnMessage = websocket.value.onmessage

    websocket.value.onmessage = (event) => {
      if (originalOnMessage) originalOnMessage(event)

      try {
        const wsEvent = JSON.parse(event.data)
        if (wsEvent.type === 'new_message') {
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
