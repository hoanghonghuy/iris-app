package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

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
// Nếu đã tồn tại, trả về cuộc hội thoại cũ và created=false. Nếu chưa, tạo mới và created=true.
func (s *ChatService) GetOrCreateDirectConversation(ctx context.Context, userA, userB uuid.UUID) (*model.Conversation, bool, error) {
	if userA == userB {
		return nil, false, ErrChatCannotMessageSelf
	}

	conv, err := s.chatRepo.FindDirectConversation(ctx, userA, userB)
	if err == nil {
		return conv, false, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, false, err
	}

	// Chưa có → tạo mới
	conv, err = s.chatRepo.CreateConversation(ctx, "direct", nil, []uuid.UUID{userA, userB})
	if err != nil {
		return nil, false, err
	}
	return conv, true, nil
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

const maxChatGroupParticipants = 50

// CreateGroupConversation tạo cuộc hội thoại nhóm mới (internal; participantIDs gồm cả chủ phòng).
func (s *ChatService) CreateGroupConversation(ctx context.Context, name string, participantIDs []uuid.UUID) (*model.Conversation, error) {
	if len(participantIDs) < 2 {
		return nil, ErrChatGroupNeedMembers
	}
	if len(participantIDs) > maxChatGroupParticipants {
		return nil, ErrChatGroupTooManyMembers
	}
	var namePtr *string
	if t := strings.TrimSpace(name); t != "" {
		if utf8.RuneCountInString(t) > 255 {
			return nil, ErrChatGroupNameTooLong
		}
		namePtr = &t
	}
	return s.chatRepo.CreateConversation(ctx, "group", namePtr, participantIDs)
}

// CreateGroupConversationAsRequester tạo nhóm: luôn có requester trong danh sách; otherUserIDs là các user khác (ít nhất 1).
// Mỗi thành viên khác phải thỏa cùng quy tắc nhắn direct với SearchUsers/CanCreateDirectConversation.
func (s *ChatService) CreateGroupConversationAsRequester(ctx context.Context, requesterID uuid.UUID, roles []string, name string, otherUserIDs []uuid.UUID) (*model.Conversation, error) {
	seen := map[uuid.UUID]struct{}{requesterID: {}}
	var orderedOthers []uuid.UUID
	for _, id := range otherUserIDs {
		if id == requesterID {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		allowed, err := s.CanCreateDirectConversation(ctx, requesterID, roles, id)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, ErrChatTargetNotAllowed
		}
		seen[id] = struct{}{}
		orderedOthers = append(orderedOthers, id)
	}
	if len(orderedOthers) < 1 {
		return nil, ErrChatGroupNeedMembers
	}

	participants := make([]uuid.UUID, 0, 1+len(orderedOthers))
	participants = append(participants, requesterID)
	participants = append(participants, orderedOthers...)

	return s.CreateGroupConversation(ctx, name, participants)
}

// ListConversations lấy danh sách cuộc hội thoại của user
func (s *ChatService) ListConversations(ctx context.Context, userID uuid.UUID) ([]model.ConversationWithParticipants, error) {
	convs, err := s.chatRepo.ListConversationsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	conversationIDs := make([]uuid.UUID, 0, len(convs))
	for _, c := range convs {
		conversationIDs = append(conversationIDs, c.ConversationID)
	}

	participantsByConversation, err := s.chatRepo.ListParticipantsByConversationIDs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}

	var result []model.ConversationWithParticipants
	for _, c := range convs {
		parts := participantsByConversation[c.ConversationID]
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
