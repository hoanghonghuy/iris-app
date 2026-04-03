package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/jackc/pgx/v5"
)

// ChatService xử lý business logic cho hệ thống chat
type ChatService struct {
	chatRepo *repo.ChatRepo
}

// NewChatService tạo mới ChatService
func NewChatService(chatRepo *repo.ChatRepo) *ChatService {
	return &ChatService{chatRepo: chatRepo}
}

// GetOrCreateDirectConversation tìm hoặc tạo cuộc hội thoại direct giữa 2 user.
// Nếu đã tồn tại, trả về cuộc hội thoại cũ. Nếu chưa, tạo mới.
func (s *ChatService) GetOrCreateDirectConversation(ctx context.Context, userA, userB uuid.UUID) (*model.Conversation, error) {
	if userA == userB {
		return nil, ErrChatCannotMessageSelf
	}

	// Tìm cuộc hội thoại direct đã tồn tại
	conv, err := s.chatRepo.FindDirectConversation(ctx, userA, userB)
	if err == nil {
		return conv, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Chưa có → tạo mới
	return s.chatRepo.CreateConversation(ctx, "direct", nil, []uuid.UUID{userA, userB})
}

// CanCreateDirectConversation kiểm tra requester có được phép tạo direct conversation với target không.
// Rule được đồng bộ với visibility logic của SearchUsers theo từng role.
func (s *ChatService) CanCreateDirectConversation(ctx context.Context, requesterID uuid.UUID, roles []string, targetID uuid.UUID) (bool, error) {
	if requesterID == targetID {
		return false, ErrChatCannotMessageSelf
	}

	hasRole := func(r string) bool {
		for _, role := range roles {
			if role == r {
				return true
			}
		}
		return false
	}

	if hasRole("SUPER_ADMIN") {
		return s.chatRepo.CanSuperAdminMessageTarget(ctx, requesterID, targetID)
	}
	if hasRole("SCHOOL_ADMIN") {
		return s.chatRepo.CanSchoolAdminMessageTarget(ctx, requesterID, targetID)
	}
	if hasRole("TEACHER") {
		return s.chatRepo.CanTeacherMessageTarget(ctx, requesterID, targetID)
	}
	if hasRole("PARENT") {
		return s.chatRepo.CanParentMessageTarget(ctx, requesterID, targetID)
	}

	return false, nil
}

// CreateGroupConversation tạo cuộc hội thoại nhóm mới
func (s *ChatService) CreateGroupConversation(ctx context.Context, name string, participantIDs []uuid.UUID) (*model.Conversation, error) {
	if len(participantIDs) < 2 {
		return nil, ErrChatGroupNeedMembers
	}
	return s.chatRepo.CreateConversation(ctx, "group", &name, participantIDs)
}

// ListConversations lấy danh sách cuộc hội thoại của user
func (s *ChatService) ListConversations(ctx context.Context, userID uuid.UUID) ([]model.ConversationWithParticipants, error) {
	convs, err := s.chatRepo.ListConversationsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []model.ConversationWithParticipants
	for _, c := range convs {
		parts, err := s.chatRepo.GetParticipants(ctx, c.ConversationID)
		if err != nil {
			return nil, err
		}
		result = append(result, model.ConversationWithParticipants{
			Conversation: c,
			Participants: parts,
		})
	}
	return result, nil
}

// SendMessage gửi tin nhắn vào cuộc hội thoại (kiểm tra quyền truy cập)
func (s *ChatService) SendMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string) (*model.Message, error) {
	if content == "" {
		return nil, ErrChatEmptyMessage
	}

	// Kiểm tra user có thuộc cuộc hội thoại không
	ok, err := s.chatRepo.IsParticipant(ctx, conversationID, senderID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrChatNotParticipant
	}

	return s.chatRepo.CreateMessage(ctx, conversationID, senderID, content)
}

// ListMessages lấy danh sách tin nhắn theo cursor (kiểm tra quyền truy cập)
func (s *ChatService) ListMessages(ctx context.Context, conversationID, userID uuid.UUID, before *uuid.UUID, limit int) ([]model.MessageWithSender, error) {
	// Kiểm tra user có thuộc cuộc hội thoại không
	ok, err := s.chatRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrChatNotParticipant
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.chatRepo.ListMessages(ctx, conversationID, before, limit)
}

// GetParticipantIDs lấy danh sách user_id trong cuộc hội thoại (dùng cho broadcast WS)
func (s *ChatService) GetParticipantIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	parts, err := s.chatRepo.GetParticipants(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(parts))
	for i, p := range parts {
		ids[i] = p.UserID
	}
	return ids, nil
}

// SearchUsers tìm kiếm user qua email hoặc tên tùy theo vai trò của người gọi
func (s *ChatService) SearchUsers(ctx context.Context, requesterID uuid.UUID, roles []string, keyword string) ([]model.ParticipantInfo, error) {
	if keyword == "" {
		return []model.ParticipantInfo{}, nil
	}

	limit := 10

	hasRole := func(r string) bool {
		for _, role := range roles {
			if role == r {
				return true
			}
		}
		return false
	}

	if hasRole("SUPER_ADMIN") {
		return s.chatRepo.SearchUsersGlobal(ctx, keyword, limit)
	}
	if hasRole("SCHOOL_ADMIN") {
		return s.chatRepo.SearchUsersForSchoolAdmin(ctx, requesterID, keyword, limit)
	}
	if hasRole("TEACHER") {
		return s.chatRepo.SearchUsersForTeacher(ctx, requesterID, keyword, limit)
	}
	if hasRole("PARENT") {
		return s.chatRepo.SearchUsersForParent(ctx, requesterID, keyword, limit)
	}

	return []model.ParticipantInfo{}, nil
}
