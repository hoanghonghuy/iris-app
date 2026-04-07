package chathandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestCreateDirectConversation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	targetID := uuid.New()

	t.Run("unauthorized", func(t *testing.T) {
		h := &ChatHandler{}
		r := gin.New()
		r.POST("/conversations/direct", h.CreateDirectConversation)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{"target_user_id":"`+targetID.String()+`"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := &ChatHandler{}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.POST("/conversations/direct", h.CreateDirectConversation)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{bad-json`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if got := decodeChatError(t, rec); got != "invalid request body" {
			t.Fatalf("error = %q, want %q", got, "invalid request body")
		}
	})

	t.Run("cannot create with self", func(t *testing.T) {
		h := &ChatHandler{}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.POST("/conversations/direct", h.CreateDirectConversation)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{"target_user_id":"`+userID.String()+`"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if got := decodeChatError(t, rec); got != service.ErrChatCannotMessageSelf.Error() {
			t.Fatalf("error = %q, want %q", got, service.ErrChatCannotMessageSelf.Error())
		}
	})

	t.Run("permission check failed", func(t *testing.T) {
		h := &ChatHandler{chatService: &fakeChatService{canCreateDirectConversationFn: func(context.Context, uuid.UUID, []string, uuid.UUID) (bool, error) {
			return false, errors.New("boom")
		}}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.POST("/conversations/direct", h.CreateDirectConversation)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{"target_user_id":"`+targetID.String()+`"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})

	t.Run("target not allowed", func(t *testing.T) {
		h := &ChatHandler{chatService: &fakeChatService{canCreateDirectConversationFn: func(context.Context, uuid.UUID, []string, uuid.UUID) (bool, error) {
			return false, nil
		}}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.POST("/conversations/direct", h.CreateDirectConversation)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{"target_user_id":"`+targetID.String()+`"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
		if got := decodeChatError(t, rec); got != service.ErrChatTargetNotAllowed.Error() {
			t.Fatalf("error = %q, want %q", got, service.ErrChatTargetNotAllowed.Error())
		}
	})

	t.Run("create conversation mappings and success", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
		}{
			{name: "cannot message self", err: service.ErrChatCannotMessageSelf, wantStatus: http.StatusBadRequest},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ChatHandler{chatService: &fakeChatService{
					canCreateDirectConversationFn: func(context.Context, uuid.UUID, []string, uuid.UUID) (bool, error) { return true, nil },
					getOrCreateDirectFn:           func(context.Context, uuid.UUID, uuid.UUID) (*model.Conversation, error) { return nil, tc.err },
				}}
				r := gin.New()
				r.Use(withChatClaims(userID, "PARENT"))
				r.POST("/conversations/direct", h.CreateDirectConversation)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{"target_user_id":"`+targetID.String()+`"}`))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
			})
		}

		h := &ChatHandler{chatService: &fakeChatService{
			canCreateDirectConversationFn: func(context.Context, uuid.UUID, []string, uuid.UUID) (bool, error) { return true, nil },
			getOrCreateDirectFn: func(_ context.Context, callerID, gotTargetID uuid.UUID) (*model.Conversation, error) {
				if callerID != userID || gotTargetID != targetID {
					t.Fatalf("unexpected ids")
				}
				return &model.Conversation{ConversationID: uuid.New(), Type: "direct", CreatedAt: time.Now()}, nil
			},
		}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.POST("/conversations/direct", h.CreateDirectConversation)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/conversations/direct", strings.NewReader(`{"target_user_id":"`+targetID.String()+`"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

func TestListConversationsAndSearchUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()

	t.Run("list conversations success and error", func(t *testing.T) {
		h := &ChatHandler{chatService: &fakeChatService{listConversationsFn: func(_ context.Context, gotUserID uuid.UUID) ([]model.ConversationWithParticipants, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id")
			}
			return []model.ConversationWithParticipants{}, nil
		}}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/conversations", h.ListConversations)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/conversations", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		h = &ChatHandler{chatService: &fakeChatService{listConversationsFn: func(context.Context, uuid.UUID) ([]model.ConversationWithParticipants, error) {
			return nil, errors.New("boom")
		}}}
		r = gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/conversations", h.ListConversations)

		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/conversations", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})

	t.Run("search users empty keyword", func(t *testing.T) {
		h := &ChatHandler{}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/users/search", h.SearchUsers)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/search", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("search users success and error", func(t *testing.T) {
		h := &ChatHandler{chatService: &fakeChatService{searchUsersFn: func(_ context.Context, gotUserID uuid.UUID, roles []string, keyword string) ([]model.ParticipantInfo, error) {
			if gotUserID != userID || keyword != "ann" || len(roles) != 1 || roles[0] != "PARENT" {
				t.Fatalf("unexpected search params")
			}
			return []model.ParticipantInfo{{UserID: uuid.New(), Email: "ann@example.com"}}, nil
		}}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/users/search", h.SearchUsers)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/search?q=ann", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		h = &ChatHandler{chatService: &fakeChatService{searchUsersFn: func(context.Context, uuid.UUID, []string, string) ([]model.ParticipantInfo, error) {
			return nil, errors.New("boom")
		}}}
		r = gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/users/search", h.SearchUsers)

		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/users/search?q=ann", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}

func TestListMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	conversationID := uuid.New()
	beforeID := uuid.New()

	t.Run("invalid conversation id", func(t *testing.T) {
		h := &ChatHandler{}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/conversations/:conversation_id/messages", h.ListMessages)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/conversations/not-a-uuid/messages", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid before cursor", func(t *testing.T) {
		h := &ChatHandler{}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/conversations/:conversation_id/messages", h.ListMessages)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/conversations/"+conversationID.String()+"/messages?before=bad-cursor", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("service mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
		}{
			{name: "not participant", err: service.ErrChatNotParticipant, wantStatus: http.StatusForbidden},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ChatHandler{chatService: &fakeChatService{listMessagesFn: func(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID, int) ([]model.MessageWithSender, error) {
					return nil, tc.err
				}}}
				r := gin.New()
				r.Use(withChatClaims(userID, "PARENT"))
				r.GET("/conversations/:conversation_id/messages", h.ListMessages)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/conversations/"+conversationID.String()+"/messages", nil)
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
			})
		}
	})

	t.Run("success with cursor and has_more", func(t *testing.T) {
		oldestID := uuid.New()
		h := &ChatHandler{chatService: &fakeChatService{listMessagesFn: func(_ context.Context, gotConversationID, gotUserID uuid.UUID, before *uuid.UUID, limit int) ([]model.MessageWithSender, error) {
			if gotConversationID != conversationID || gotUserID != userID {
				t.Fatalf("unexpected ids")
			}
			if before == nil || *before != beforeID {
				t.Fatalf("expected before cursor")
			}
			if limit != 2 {
				t.Fatalf("limit = %d, want %d", limit, 2)
			}
			return []model.MessageWithSender{
				{MessageID: uuid.New(), ConversationID: conversationID, SenderID: userID, Content: "new"},
				{MessageID: oldestID, ConversationID: conversationID, SenderID: userID, Content: "old"},
			}, nil
		}}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/conversations/:conversation_id/messages", h.ListMessages)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/conversations/"+conversationID.String()+"/messages?limit=2&before="+beforeID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var body struct {
			HasMore    bool      `json:"has_more"`
			NextCursor uuid.UUID `json:"next_cursor"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}
		if !body.HasMore {
			t.Fatalf("has_more = false, want true")
		}
		if body.NextCursor != oldestID {
			t.Fatalf("next_cursor = %s, want %s", body.NextCursor, oldestID)
		}
	})

	t.Run("success with nil list", func(t *testing.T) {
		h := &ChatHandler{chatService: &fakeChatService{listMessagesFn: func(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID, int) ([]model.MessageWithSender, error) {
			return nil, nil
		}}}
		r := gin.New()
		r.Use(withChatClaims(userID, "PARENT"))
		r.GET("/conversations/:conversation_id/messages", h.ListMessages)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/conversations/"+conversationID.String()+"/messages", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Body.String(), `"data":[]`) {
			t.Fatalf("expected empty data array, got %s", rec.Body.String())
		}
	})
}
