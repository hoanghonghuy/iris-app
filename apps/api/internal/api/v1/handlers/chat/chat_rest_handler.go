package chathandlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// CreateDirectConversationRequest chứa thông tin tạo cuộc hội thoại direct.
type CreateDirectConversationRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

// CreateDirectConversation tạo hoặc tìm cuộc hội thoại direct giữa 2 user.
func (h *ChatHandler) CreateDirectConversation(c *gin.Context) {
	userID, claims, ok := shared.RequireCurrentUser(c)
	if !ok {
		return
	}

	var req CreateDirectConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TargetUserID == userID {
		response.Fail(c, http.StatusBadRequest, service.ErrChatCannotMessageSelf.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	allowed, err := h.chatService.CanCreateDirectConversation(ctx, userID, claims.Roles, req.TargetUserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to verify conversation permission")
		return
	}
	if !allowed {
		response.Fail(c, http.StatusForbidden, service.ErrChatTargetNotAllowed.Error())
		return
	}

	conv, created, err := h.chatService.GetOrCreateDirectConversation(ctx, userID, req.TargetUserID)
	if err != nil {
		if errors.Is(err, service.ErrChatCannotMessageSelf) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create conversation")
		return
	}

	if created {
		response.Created(c, conv)
		return
	}
	response.OK(c, conv)
}

// CreateGroupConversationRequest chứa thông tin tạo cuộc hội thoại nhóm.
// participant_user_ids: danh sách user khác (không bắt buộc gửi chính mình); server luôn thêm requester.
type CreateGroupConversationRequest struct {
	Name               string      `json:"name"`
	ParticipantUserIDs []uuid.UUID `json:"participant_user_ids" binding:"required"`
}

// CreateGroupConversation tạo cuộc hội thoại nhóm; requester luôn là thành viên.
func (h *ChatHandler) CreateGroupConversation(c *gin.Context) {
	userID, claims, ok := shared.RequireCurrentUser(c)
	if !ok {
		return
	}

	var req CreateGroupConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.ParticipantUserIDs) == 0 {
		response.Fail(c, http.StatusBadRequest, service.ErrChatGroupNeedMembers.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	conv, err := h.chatService.CreateGroupConversationAsRequester(ctx, userID, claims.Roles, req.Name, req.ParticipantUserIDs)
	if err != nil {
		if errors.Is(err, service.ErrChatGroupNeedMembers) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatTargetNotAllowed) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatCannotMessageSelf) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatGroupNameTooLong) || errors.Is(err, service.ErrChatGroupTooManyMembers) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create group conversation")
		return
	}

	response.Created(c, conv)
}

// ListConversations lấy danh sách cuộc hội thoại của user hiện tại.
func (h *ChatHandler) ListConversations(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
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

// MarkConversationRead đánh dấu đã đọc tới tin mới nhất (dùng khi đang mở hội thoại qua WebSocket).
func (h *ChatHandler) MarkConversationRead(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid conversation_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.chatService.MarkConversationRead(ctx, conversationID, userID); err != nil {
		if errors.Is(err, service.ErrChatNotParticipant) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to mark conversation read")
		return
	}

	response.OK(c, gin.H{"conversation_id": conversationID.String(), "ok": true})
}

// ListMessages lấy danh sách tin nhắn theo cursor.
// Query params: ?before=<message_uuid>&limit=<int>
// Response: { data, has_more, next_cursor }
func (h *ChatHandler) ListMessages(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
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
	if limit > 100 {
		limit = 100
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

	// Đã tải (một phần) lịch sử → coi như đã đọc tới tin mới nhất trong hội thoại.
	_ = h.chatService.MarkConversationRead(ctx, conversationID, userID)

	// next_cursor là message_id của tin nhắn cũ nhất trong batch (phần tử cuối, vì DESC)
	// fe dùng giá trị này cho lần fetch tiếp theo khi user cuộn lên
	var nextCursor *uuid.UUID
	hasMore := len(msgs) == limit
	if hasMore {
		id := msgs[len(msgs)-1].MessageID
		nextCursor = &id
	}

	if msgs == nil {
		msgs = []model.MessageWithSender{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        msgs,
		"has_more":    hasMore,
		"next_cursor": nextCursor,
	})
}

// PatchGroupConversationRequest đổi tên nhóm (name rỗng → xóa tên hiển thị).
type PatchGroupConversationRequest struct {
	Name string `json:"name"`
}

// PatchGroupConversation cập nhật tên nhóm chat.
func (h *ChatHandler) PatchGroupConversation(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid conversation_id format")
		return
	}

	var req PatchGroupConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cw, err := h.chatService.RenameGroupConversation(ctx, userID, conversationID, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrChatNotParticipant) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatNotGroup) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatConversationNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatGroupNameTooLong) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update group")
		return
	}

	response.OK(c, cw)
}

// AddGroupParticipantsRequest thêm thành viên vào nhóm.
type AddGroupParticipantsRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
}

// AddGroupParticipants POST thêm thành viên.
func (h *ChatHandler) AddGroupParticipants(c *gin.Context) {
	userID, claims, ok := shared.RequireCurrentUser(c)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid conversation_id format")
		return
	}

	var req AddGroupParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.UserIDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "user_ids must not be empty")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cw, err := h.chatService.AddGroupParticipants(ctx, userID, claims.Roles, conversationID, req.UserIDs)
	if err != nil {
		if errors.Is(err, service.ErrChatNotParticipant) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatNotGroup) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatConversationNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatTargetNotAllowed) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatGroupTooManyMembers) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to add participants")
		return
	}

	response.OK(c, cw)
}

// RemoveGroupParticipant DELETE một thành viên khỏi nhóm.
func (h *ChatHandler) RemoveGroupParticipant(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid conversation_id format")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cw, err := h.chatService.RemoveGroupParticipant(ctx, userID, conversationID, targetUserID)
	if err != nil {
		if errors.Is(err, service.ErrChatNotParticipant) {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatNotGroup) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatConversationNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatParticipantNotInGroup) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrChatCannotRemoveWouldDropBelowMin) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to remove participant")
		return
	}

	response.OK(c, cw)
}

// SearchUsers tìm kiếm user qua query param "q".
func (h *ChatHandler) SearchUsers(c *gin.Context) {
	userID, claims, ok := shared.RequireCurrentUser(c)
	if !ok {
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
