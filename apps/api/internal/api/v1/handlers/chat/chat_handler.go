package chathandlers

import (
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/ws"
)

// ChatHandler xử lý các endpoint REST và WebSocket cho hệ thống chat.
type ChatHandler struct {
	chatService    *service.ChatService
	hub            *ws.Hub
	jwtSecret      string
	allowedOrigins map[string]struct{} // origin allowlist cho WS
}

// NewChatHandler tạo mới ChatHandler.
func NewChatHandler(chatService *service.ChatService, hub *ws.Hub, jwtSecret string, allowedOrigins []string) *ChatHandler {
	set := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		set[o] = struct{}{}
	}
	return &ChatHandler{
		chatService:    chatService,
		hub:            hub,
		jwtSecret:      jwtSecret,
		allowedOrigins: set,
	}
}
