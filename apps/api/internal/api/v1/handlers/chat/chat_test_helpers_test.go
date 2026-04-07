package chathandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeChatService struct {
	canCreateDirectConversationFn func(context.Context, uuid.UUID, []string, uuid.UUID) (bool, error)
	getOrCreateDirectFn           func(context.Context, uuid.UUID, uuid.UUID) (*model.Conversation, error)
	listConversationsFn           func(context.Context, uuid.UUID) ([]model.ConversationWithParticipants, error)
	listMessagesFn                func(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID, int) ([]model.MessageWithSender, error)
	searchUsersFn                 func(context.Context, uuid.UUID, []string, string) ([]model.ParticipantInfo, error)
	sendMessageFn                 func(context.Context, uuid.UUID, uuid.UUID, string) (*model.Message, error)
	getParticipantIDsFn           func(context.Context, uuid.UUID) ([]uuid.UUID, error)
}

func (f *fakeChatService) CanCreateDirectConversation(ctx context.Context, requesterID uuid.UUID, roles []string, targetID uuid.UUID) (bool, error) {
	if f.canCreateDirectConversationFn == nil {
		return false, errors.New("unexpected CanCreateDirectConversation call")
	}
	return f.canCreateDirectConversationFn(ctx, requesterID, roles, targetID)
}

func (f *fakeChatService) GetOrCreateDirectConversation(ctx context.Context, userA, userB uuid.UUID) (*model.Conversation, error) {
	if f.getOrCreateDirectFn == nil {
		return nil, errors.New("unexpected GetOrCreateDirectConversation call")
	}
	return f.getOrCreateDirectFn(ctx, userA, userB)
}

func (f *fakeChatService) ListConversations(ctx context.Context, userID uuid.UUID) ([]model.ConversationWithParticipants, error) {
	if f.listConversationsFn == nil {
		return nil, errors.New("unexpected ListConversations call")
	}
	return f.listConversationsFn(ctx, userID)
}

func (f *fakeChatService) ListMessages(ctx context.Context, conversationID, userID uuid.UUID, before *uuid.UUID, limit int) ([]model.MessageWithSender, error) {
	if f.listMessagesFn == nil {
		return nil, errors.New("unexpected ListMessages call")
	}
	return f.listMessagesFn(ctx, conversationID, userID, before, limit)
}

func (f *fakeChatService) SearchUsers(ctx context.Context, requesterID uuid.UUID, roles []string, keyword string) ([]model.ParticipantInfo, error) {
	if f.searchUsersFn == nil {
		return nil, errors.New("unexpected SearchUsers call")
	}
	return f.searchUsersFn(ctx, requesterID, roles, keyword)
}

func (f *fakeChatService) SendMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string) (*model.Message, error) {
	if f.sendMessageFn == nil {
		return nil, errors.New("unexpected SendMessage call")
	}
	return f.sendMessageFn(ctx, conversationID, senderID, content)
}

func (f *fakeChatService) GetParticipantIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	if f.getParticipantIDsFn == nil {
		return nil, errors.New("unexpected GetParticipantIDs call")
	}
	return f.getParticipantIDsFn(ctx, conversationID)
}

func withChatClaims(userID uuid.UUID, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{UserID: userID.String(), Roles: roles})
		c.Next()
	}
}

func decodeChatError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
