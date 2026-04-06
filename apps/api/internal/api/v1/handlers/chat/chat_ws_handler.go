package chathandlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/ws"
)

// wsMessage cấu trúc JSON mà client gửi lên qua WebSocket.
type wsMessage struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

// HandleWS xử lý upgrade HTTP -> WebSocket.
// Token JWT được truyền qua Sec-WebSocket-Protocol header:
//   - Client gửi: Sec-WebSocket-Protocol: Bearer, <JWT>
//   - Server validate token và reply với cùng sub-protocol để browser chấp nhận
func (h *ChatHandler) HandleWS(c *gin.Context) {
	// Tạo upgrader với CheckOrigin theo allowlist
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// non-browser client (e.g. server-to-server) - cho phép
				return true
			}
			_, ok := h.allowedOrigins[origin]
			return ok
		},
	}

	// Đọc token từ Sec-WebSocket-Protocol (format: "Bearer, <JWT>")
	// Browser WebSocket API không cho phép set custom headers, nên dùng sub-protocol thay thế.
	// RFC 6455 §4.1: browser sẽ gửi lại Sec-WebSocket-Protocol trong response nếu server echo lại.
	protocols := websocket.Subprotocols(c.Request)
	var tokenStr string
	for i, p := range protocols {
		if p == "Bearer" && i+1 < len(protocols) {
			tokenStr = protocols[i+1]
			break
		}
	}

	if tokenStr == "" {
		response.Fail(c, http.StatusUnauthorized, "missing token in Sec-WebSocket-Protocol")
		return
	}

	claims, err := auth.Parse(h.jwtSecret, tokenStr)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "invalid token")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID in token")
		return
	}

	// Echo lại sub-protocol "Bearer" để browser chấp nhận handshake
	conn, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-WebSocket-Protocol": {"Bearer"},
	})
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		Hub:    h.hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}

	// Callback khi client gửi tin nhắn qua WS
	client.OnMessage = func(senderID uuid.UUID, data []byte) {
		var msg wsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("[WS] invalid message from user %s: %v", senderID, err)
			return
		}

		convID, err := uuid.Parse(msg.ConversationID)
		if err != nil {
			log.Printf("[WS] invalid conversation_id from user %s", senderID)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		saved, err := h.chatService.SendMessage(ctx, convID, senderID, msg.Content)
		if err != nil {
			log.Printf("[WS] sendMessage error: %v", err)
			return
		}

		// Lấy danh sách participant để broadcast
		participantIDs, err := h.chatService.GetParticipantIDs(ctx, convID)
		if err != nil {
			log.Printf("[WS] getParticipants error: %v", err)
			return
		}

		// Broadcast tin nhắn mới đến tất cả thành viên
		h.hub.BroadcastToUsers(participantIDs, ws.WSEvent{
			Type: "new_message",
			Data: saved,
		})
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
