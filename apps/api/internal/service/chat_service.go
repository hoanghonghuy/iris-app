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
	summaries, err := s.chatRepo.ListConversationSummariesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	conversationIDs := make([]uuid.UUID, 0, len(summaries))
	for _, row := range summaries {
		conversationIDs = append(conversationIDs, row.ConversationID)
	}

	participantsByConversation, err := s.chatRepo.ListParticipantsByConversationIDs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}

	var result []model.ConversationWithParticipants
	for _, row := range summaries {
		parts := participantsByConversation[row.ConversationID]
		result = append(result, model.ConversationWithParticipants{
			Conversation: row.Conversation,
			Participants: parts,
			LastMessage:  row.LastMessage,
			UnreadCount:  row.UnreadCount,
		})
	}
	return result, nil
}

// MarkConversationRead đánh dấu đã xem tới tin mới nhất (sidebar unread về 0 cho viewer).
func (s *ChatService) MarkConversationRead(ctx context.Context, conversationID, userID uuid.UUID) error {
	ok, err := s.chatRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrChatNotParticipant
	}
	return s.chatRepo.MarkConversationRead(ctx, conversationID, userID)
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

func (s *ChatService) conversationWithParticipants(ctx context.Context, convID uuid.UUID) (*model.ConversationWithParticipants, error) {
	conv, err := s.chatRepo.GetConversationByID(ctx, convID)
	if err != nil {
		return nil, err
	}
	parts, err := s.chatRepo.GetParticipants(ctx, convID)
	if err != nil {
		return nil, err
	}
	return &model.ConversationWithParticipants{
		Conversation: *conv,
		Participants: parts,
	}, nil
}

func groupNamePtrFromInput(name string) (*string, error) {
	t := strings.TrimSpace(name)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > 255 {
		return nil, ErrChatGroupNameTooLong
	}
	return &t, nil
}

// ensureGroupAndParticipant: actor thuộc cuộc hội thoại và cuộc là nhóm.
func (s *ChatService) ensureGroupAndParticipant(ctx context.Context, convID, actorID uuid.UUID) (*model.Conversation, error) {
	ok, err := s.chatRepo.IsParticipant(ctx, convID, actorID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrChatNotParticipant
	}
	conv, err := s.chatRepo.GetConversationByID(ctx, convID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChatConversationNotFound
		}
		return nil, err
	}
	if conv.Type != "group" {
		return nil, ErrChatNotGroup
	}
	return conv, nil
}

// RenameGroupConversation đổi tên nhóm (chuỗi rỗng sau trim → name NULL).
func (s *ChatService) RenameGroupConversation(ctx context.Context, actorID, convID uuid.UUID, name string) (*model.ConversationWithParticipants, error) {
	if _, err := s.ensureGroupAndParticipant(ctx, convID, actorID); err != nil {
		return nil, err
	}
	namePtr, err := groupNamePtrFromInput(name)
	if err != nil {
		return nil, err
	}
	if err := s.chatRepo.UpdateConversationName(ctx, convID, namePtr); err != nil {
		return nil, err
	}
	return s.conversationWithParticipants(ctx, convID)
}

// AddGroupParticipants thêm thành viên; mỗi user mới phải thỏa CanCreateDirectConversation với actor.
func (s *ChatService) AddGroupParticipants(ctx context.Context, actorID uuid.UUID, roles []string, convID uuid.UUID, newUserIDs []uuid.UUID) (*model.ConversationWithParticipants, error) {
	if _, err := s.ensureGroupAndParticipant(ctx, convID, actorID); err != nil {
		return nil, err
	}
	existing, err := s.chatRepo.GetParticipants(ctx, convID)
	if err != nil {
		return nil, err
	}
	existingSet := make(map[uuid.UUID]struct{}, len(existing))
	for _, p := range existing {
		existingSet[p.UserID] = struct{}{}
	}
	var toAdd []uuid.UUID
	seen := make(map[uuid.UUID]struct{})
	for _, id := range newUserIDs {
		if _, ok := existingSet[id]; ok {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		toAdd = append(toAdd, id)
	}
	if len(toAdd) == 0 {
		return s.conversationWithParticipants(ctx, convID)
	}
	if len(existing)+len(toAdd) > maxChatGroupParticipants {
		return nil, ErrChatGroupTooManyMembers
	}
	for _, id := range toAdd {
		allowed, err := s.CanCreateDirectConversation(ctx, actorID, roles, id)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, ErrChatTargetNotAllowed
		}
	}
	if err := s.chatRepo.AddConversationParticipants(ctx, convID, toAdd); err != nil {
		return nil, err
	}
	return s.conversationWithParticipants(ctx, convID)
}

// RemoveGroupParticipant gỡ thành viên; nhóm phải còn ít nhất 2 người sau khi gỡ.
func (s *ChatService) RemoveGroupParticipant(ctx context.Context, actorID, convID, targetUserID uuid.UUID) (*model.ConversationWithParticipants, error) {
	if _, err := s.ensureGroupAndParticipant(ctx, convID, actorID); err != nil {
		return nil, err
	}
	ok, err := s.chatRepo.IsParticipant(ctx, convID, targetUserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrChatParticipantNotInGroup
	}
	n, err := s.chatRepo.CountParticipants(ctx, convID)
	if err != nil {
		return nil, err
	}
	if n <= 2 {
		return nil, ErrChatCannotRemoveWouldDropBelowMin
	}
	rows, err := s.chatRepo.RemoveConversationParticipant(ctx, convID, targetUserID)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrChatParticipantNotInGroup
	}
	// Nếu actor tự xóa chính mình, không trả conversation (actor đã mất quyền truy cập).
	if actorID == targetUserID {
		return nil, nil
	}
	return s.conversationWithParticipants(ctx, convID)
}
