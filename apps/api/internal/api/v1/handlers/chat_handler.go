package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/ws"
)

// upgrader cấu hình WebSocket upgrader cho Gin
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: restrict origins in production
	},
}

// ChatHandler xử lý các endpoint REST và WebSocket cho hệ thống chat
type ChatHandler struct {
	chatService *service.ChatService
	hub         *ws.Hub
	jwtSecret   string
}

// NewChatHandler tạo mới ChatHandler
func NewChatHandler(chatService *service.ChatService, hub *ws.Hub, jwtSecret string) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		hub:         hub,
		jwtSecret:   jwtSecret,
	}
}

// --- REST API endpoints ---

// CreateDirectConversationRequest chứa thông tin tạo cuộc hội thoại direct
type CreateDirectConversationRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

// CreateDirectConversation tạo hoặc tìm cuộc hội thoại direct giữa 2 user
func (h *ChatHandler) CreateDirectConversation(c *gin.Context) {
	claimsAny, ok := c.Get(middleware.CtxClaims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req CreateDirectConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	conv, err := h.chatService.GetOrCreateDirectConversation(ctx, userID, req.TargetUserID)
	if err != nil {
		if errors.Is(err, service.ErrChatCannotMessageSelf) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create conversation")
		return
	}

	response.OK(c, conv)
}

// ListConversations lấy danh sách cuộc hội thoại của user hiện tại
func (h *ChatHandler) ListConversations(c *gin.Context) {
	claimsAny, ok := c.Get(middleware.CtxClaims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	convs, err := h.chatService.ListConversations(ctx, userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to list conversations")
		return
	}

	response.OK(c, convs)
}

// ListMessages lấy danh sách tin nhắn theo cursor.
// Query params: ?before=<message_uuid>&limit=<int>
// Response: { data, has_more, next_cursor }
func (h *ChatHandler) ListMessages(c *gin.Context) {
	claimsAny, ok := c.Get(middleware.CtxClaims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid conversation_id format")
		return
	}

	// Parse limit (default 50, max 100)
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	// Parse optional cursor: before=<message_uuid>
	var before *uuid.UUID
	if v := c.Query("before"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid before cursor (must be a valid message UUID)")
			return
		}
		before = &id
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	msgs, err := h.chatService.ListMessages(ctx, conversationID, userID, before, limit)
	if err != nil {
		if errors.Is(err, service.ErrChatNotParticipant) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list messages")
		return
	}

	// next_cursor là message_id của tin nhắn cũ nhất trong batch (phần tử cuối, vì DESC)
	// fe dùng giá trị này cho lần fetch tiếp theo khi user cuộn lên
	var nextCursor *uuid.UUID
	hasMore := len(msgs) == limit
	if hasMore {
		id := msgs[len(msgs)-1].MessageID
		nextCursor = &id
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        msgs,
		"has_more":    hasMore,
		"next_cursor": nextCursor,
	})
}


// --- WebSocket endpoint ---

// wsMessage cấu trúc JSON mà client gửi lên qua WebSocket
type wsMessage struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

// HandleWS xử lý upgrade HTTP → WebSocket.
// Token JWT được truyền qua query string: /api/v1/chat/ws?token=<JWT>
func (h *ChatHandler) HandleWS(c *gin.Context) {
	// Lấy token từ query string (WebSocket browser không hỗ trợ custom headers)
	tokenStr := c.Query("token")
	if tokenStr == "" {
		response.Fail(c, http.StatusUnauthorized, "missing token query param")
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
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

// SearchUsers tìm kiếm user qua query param "q"
func (h *ChatHandler) SearchUsers(c *gin.Context) {
	claimsAny, ok := c.Get(middleware.CtxClaims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Lấy keyword
	q := c.Query("q")
	if q == "" {
		response.OK(c, []model.ParticipantInfo{})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	users, err := h.chatService.SearchUsers(ctx, userID, claims.Roles, q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to search users")
		return
	}

	response.OK(c, users)
}
