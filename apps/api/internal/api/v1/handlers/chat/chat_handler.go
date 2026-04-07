package chathandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/ws"
)

type chatService interface {
	CanCreateDirectConversation(context.Context, uuid.UUID, []string, uuid.UUID) (bool, error)
	GetOrCreateDirectConversation(context.Context, uuid.UUID, uuid.UUID) (*model.Conversation, error)
	ListConversations(context.Context, uuid.UUID) ([]model.ConversationWithParticipants, error)
	ListMessages(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID, int) ([]model.MessageWithSender, error)
	SearchUsers(context.Context, uuid.UUID, []string, string) ([]model.ParticipantInfo, error)
	SendMessage(context.Context, uuid.UUID, uuid.UUID, string) (*model.Message, error)
	GetParticipantIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
}

// ChatHandler xử lý các endpoint REST và WebSocket cho hệ thống chat.
type ChatHandler struct {
	chatService    chatService
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
