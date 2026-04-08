package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
)

type fakeChatRepo struct {
	findDirectConversationCalls  int
	findDirectConversationUserA  uuid.UUID
	findDirectConversationUserB  uuid.UUID
	findDirectConversationResult *model.Conversation
	findDirectConversationErr    error

	createConversationCalls          int
	createConversationType           string
	createConversationName           *string
	createConversationParticipantIDs []uuid.UUID
	createConversationResult         *model.Conversation
	createConversationErr            error

	canSuperAdminCalls       int
	canSuperAdminRequesterID uuid.UUID
	canSuperAdminTargetID    uuid.UUID
	canSuperAdminResult      bool
	canSuperAdminErr         error

	canSchoolAdminCalls       int
	canSchoolAdminRequesterID uuid.UUID
	canSchoolAdminTargetID    uuid.UUID
	canSchoolAdminResult      bool
	canSchoolAdminErr         error

	canTeacherCalls       int
	canTeacherRequesterID uuid.UUID
	canTeacherTargetID    uuid.UUID
	canTeacherResult      bool
	canTeacherErr         error

	canParentCalls       int
	canParentRequesterID uuid.UUID
	canParentTargetID    uuid.UUID
	canParentResult      bool
	canParentErr         error

	listConversationsCalls  int
	listConversationsUserID uuid.UUID
	listConversationsResult []model.Conversation
	listConversationsErr    error

	listParticipantsCalls           int
	listParticipantsConversationIDs []uuid.UUID
	listParticipantsResult          map[uuid.UUID][]model.ParticipantInfo
	listParticipantsErr             error

	isParticipantCalls          int
	isParticipantConversationID uuid.UUID
	isParticipantUserID         uuid.UUID
	isParticipantResult         bool
	isParticipantErr            error

	createMessageCalls          int
	createMessageConversationID uuid.UUID
	createMessageSenderID       uuid.UUID
	createMessageContent        string
	createMessageResult         *model.Message
	createMessageErr            error

	listMessagesCalls          int
	listMessagesConversationID uuid.UUID
	listMessagesBefore         *uuid.UUID
	listMessagesLimit          int
	listMessagesResult         []model.MessageWithSender
	listMessagesErr            error

	getParticipantsCalls          int
	getParticipantsConversationID uuid.UUID
	getParticipantsResult         []model.ParticipantInfo
	getParticipantsErr            error

	searchGlobalCalls   int
	searchGlobalKeyword string
	searchGlobalLimit   int
	searchGlobalResult  []model.ParticipantInfo
	searchGlobalErr     error

	searchSchoolAdminCalls   int
	searchSchoolAdminAdminID uuid.UUID
	searchSchoolAdminKeyword string
	searchSchoolAdminLimit   int
	searchSchoolAdminResult  []model.ParticipantInfo
	searchSchoolAdminErr     error

	searchTeacherCalls   int
	searchTeacherUserID  uuid.UUID
	searchTeacherKeyword string
	searchTeacherLimit   int
	searchTeacherResult  []model.ParticipantInfo
	searchTeacherErr     error

	searchParentCalls   int
	searchParentUserID  uuid.UUID
	searchParentKeyword string
	searchParentLimit   int
	searchParentResult  []model.ParticipantInfo
	searchParentErr     error
}

func (f *fakeChatRepo) FindDirectConversation(_ context.Context, userA, userB uuid.UUID) (*model.Conversation, error) {
	f.findDirectConversationCalls++
	f.findDirectConversationUserA = userA
	f.findDirectConversationUserB = userB
	return f.findDirectConversationResult, f.findDirectConversationErr
}

func (f *fakeChatRepo) CreateConversation(_ context.Context, convType string, name *string, participantIDs []uuid.UUID) (*model.Conversation, error) {
	f.createConversationCalls++
	f.createConversationType = convType
	f.createConversationName = name
	f.createConversationParticipantIDs = append([]uuid.UUID{}, participantIDs...)
	return f.createConversationResult, f.createConversationErr
}

func (f *fakeChatRepo) CanSuperAdminMessageTarget(_ context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	f.canSuperAdminCalls++
	f.canSuperAdminRequesterID = requesterID
	f.canSuperAdminTargetID = targetID
	return f.canSuperAdminResult, f.canSuperAdminErr
}

func (f *fakeChatRepo) CanSchoolAdminMessageTarget(_ context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	f.canSchoolAdminCalls++
	f.canSchoolAdminRequesterID = requesterID
	f.canSchoolAdminTargetID = targetID
	return f.canSchoolAdminResult, f.canSchoolAdminErr
}

func (f *fakeChatRepo) CanTeacherMessageTarget(_ context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	f.canTeacherCalls++
	f.canTeacherRequesterID = requesterID
	f.canTeacherTargetID = targetID
	return f.canTeacherResult, f.canTeacherErr
}

func (f *fakeChatRepo) CanParentMessageTarget(_ context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	f.canParentCalls++
	f.canParentRequesterID = requesterID
	f.canParentTargetID = targetID
	return f.canParentResult, f.canParentErr
}

func (f *fakeChatRepo) ListConversationsByUser(_ context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	f.listConversationsCalls++
	f.listConversationsUserID = userID
	return f.listConversationsResult, f.listConversationsErr
}

func (f *fakeChatRepo) ListParticipantsByConversationIDs(_ context.Context, conversationIDs []uuid.UUID) (map[uuid.UUID][]model.ParticipantInfo, error) {
	f.listParticipantsCalls++
	f.listParticipantsConversationIDs = append([]uuid.UUID{}, conversationIDs...)
	return f.listParticipantsResult, f.listParticipantsErr
}

func (f *fakeChatRepo) IsParticipant(_ context.Context, conversationID, userID uuid.UUID) (bool, error) {
	f.isParticipantCalls++
	f.isParticipantConversationID = conversationID
	f.isParticipantUserID = userID
	return f.isParticipantResult, f.isParticipantErr
}

func (f *fakeChatRepo) CreateMessage(_ context.Context, conversationID, senderID uuid.UUID, content string) (*model.Message, error) {
	f.createMessageCalls++
	f.createMessageConversationID = conversationID
	f.createMessageSenderID = senderID
	f.createMessageContent = content
	return f.createMessageResult, f.createMessageErr
}

func (f *fakeChatRepo) ListMessages(_ context.Context, conversationID uuid.UUID, before *uuid.UUID, limit int) ([]model.MessageWithSender, error) {
	f.listMessagesCalls++
	f.listMessagesConversationID = conversationID
	f.listMessagesBefore = before
	f.listMessagesLimit = limit
	return f.listMessagesResult, f.listMessagesErr
}

func (f *fakeChatRepo) GetParticipants(_ context.Context, conversationID uuid.UUID) ([]model.ParticipantInfo, error) {
	f.getParticipantsCalls++
	f.getParticipantsConversationID = conversationID
	return f.getParticipantsResult, f.getParticipantsErr
}

func (f *fakeChatRepo) SearchUsersGlobal(_ context.Context, keyword string, limit int) ([]model.ParticipantInfo, error) {
	f.searchGlobalCalls++
	f.searchGlobalKeyword = keyword
	f.searchGlobalLimit = limit
	return f.searchGlobalResult, f.searchGlobalErr
}

func (f *fakeChatRepo) SearchUsersForSchoolAdmin(_ context.Context, adminID uuid.UUID, keyword string, limit int) ([]model.ParticipantInfo, error) {
	f.searchSchoolAdminCalls++
	f.searchSchoolAdminAdminID = adminID
	f.searchSchoolAdminKeyword = keyword
	f.searchSchoolAdminLimit = limit
	return f.searchSchoolAdminResult, f.searchSchoolAdminErr
}

func (f *fakeChatRepo) SearchUsersForTeacher(_ context.Context, teacherUserID uuid.UUID, keyword string, limit int) ([]model.ParticipantInfo, error) {
	f.searchTeacherCalls++
	f.searchTeacherUserID = teacherUserID
	f.searchTeacherKeyword = keyword
	f.searchTeacherLimit = limit
	return f.searchTeacherResult, f.searchTeacherErr
}

func (f *fakeChatRepo) SearchUsersForParent(_ context.Context, parentUserID uuid.UUID, keyword string, limit int) ([]model.ParticipantInfo, error) {
	f.searchParentCalls++
	f.searchParentUserID = parentUserID
	f.searchParentKeyword = keyword
	f.searchParentLimit = limit
	return f.searchParentResult, f.searchParentErr
}

func TestChatServiceGetOrCreateDirectConversation(t *testing.T) {
	userA := uuid.New()
	userB := uuid.New()
	conv := &model.Conversation{ConversationID: uuid.New(), Type: "direct", CreatedAt: time.Now().UTC()}
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		userA           uuid.UUID
		userB           uuid.UUID
		findResult      *model.Conversation
		findErr         error
		createResult    *model.Conversation
		createErr       error
		wantErr         error
		wantConv        *model.Conversation
		wantFindCalls   int
		wantCreateCalls int
	}{
		{name: "cannot message self", userA: userA, userB: userA, wantErr: ErrChatCannotMessageSelf},
		{name: "existing direct conversation", userA: userA, userB: userB, findResult: conv, wantConv: conv, wantFindCalls: 1},
		{name: "find direct returns non-notfound error", userA: userA, userB: userB, findErr: sentinelErr, wantErr: sentinelErr, wantFindCalls: 1},
		{name: "not found creates new conversation", userA: userA, userB: userB, findErr: pgx.ErrNoRows, createResult: conv, wantConv: conv, wantFindCalls: 1, wantCreateCalls: 1},
		{name: "create conversation fails", userA: userA, userB: userB, findErr: pgx.ErrNoRows, createErr: sentinelErr, wantErr: sentinelErr, wantFindCalls: 1, wantCreateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeChatRepo{findDirectConversationResult: tc.findResult, findDirectConversationErr: tc.findErr, createConversationResult: tc.createResult, createConversationErr: tc.createErr}
			svc := &ChatService{chatRepo: repo}

			got, err := svc.GetOrCreateDirectConversation(context.Background(), tc.userA, tc.userB)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("GetOrCreateDirectConversation() error = %v", err)
			}

			if got != tc.wantConv {
				t.Fatalf("conversation = %#v, want %#v", got, tc.wantConv)
			}
			if repo.findDirectConversationCalls != tc.wantFindCalls {
				t.Fatalf("find calls = %d, want %d", repo.findDirectConversationCalls, tc.wantFindCalls)
			}
			if repo.createConversationCalls != tc.wantCreateCalls {
				t.Fatalf("create calls = %d, want %d", repo.createConversationCalls, tc.wantCreateCalls)
			}
			if tc.wantCreateCalls > 0 {
				if repo.createConversationType != "direct" {
					t.Fatalf("create type = %q, want %q", repo.createConversationType, "direct")
				}
				if repo.createConversationName != nil {
					t.Fatalf("create name should be nil for direct conversation")
				}
				if len(repo.createConversationParticipantIDs) != 2 || repo.createConversationParticipantIDs[0] != tc.userA || repo.createConversationParticipantIDs[1] != tc.userB {
					t.Fatalf("create participants = %#v", repo.createConversationParticipantIDs)
				}
			}
		})
	}
}

func TestChatServiceCanCreateDirectConversation(t *testing.T) {
	requesterID := uuid.New()
	targetID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name                 string
		requesterID          uuid.UUID
		targetID             uuid.UUID
		roles                []string
		superAdminResult     bool
		superAdminErr        error
		schoolAdminResult    bool
		schoolAdminErr       error
		teacherResult        bool
		teacherErr           error
		parentResult         bool
		parentErr            error
		wantResult           bool
		wantErr              error
		wantSuperCalls       int
		wantSchoolAdminCalls int
		wantTeacherCalls     int
		wantParentCalls      int
	}{
		{name: "cannot message self", requesterID: requesterID, targetID: requesterID, roles: []string{"SUPER_ADMIN"}, wantErr: ErrChatCannotMessageSelf},
		{name: "super admin allowed", requesterID: requesterID, targetID: targetID, roles: []string{"SUPER_ADMIN"}, superAdminResult: true, wantResult: true, wantSuperCalls: 1},
		{name: "school admin error", requesterID: requesterID, targetID: targetID, roles: []string{"SCHOOL_ADMIN"}, schoolAdminErr: sentinelErr, wantErr: sentinelErr, wantSchoolAdminCalls: 1},
		{name: "teacher denied", requesterID: requesterID, targetID: targetID, roles: []string{"TEACHER"}, teacherResult: false, wantResult: false, wantTeacherCalls: 1},
		{name: "parent allowed", requesterID: requesterID, targetID: targetID, roles: []string{"PARENT"}, parentResult: true, wantResult: true, wantParentCalls: 1},
		{name: "role priority uses super admin first", requesterID: requesterID, targetID: targetID, roles: []string{"TEACHER", "SUPER_ADMIN"}, superAdminResult: true, wantResult: true, wantSuperCalls: 1},
		{name: "unknown role returns false", requesterID: requesterID, targetID: targetID, roles: []string{"STUDENT"}, wantResult: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeChatRepo{
				canSuperAdminResult:  tc.superAdminResult,
				canSuperAdminErr:     tc.superAdminErr,
				canSchoolAdminResult: tc.schoolAdminResult,
				canSchoolAdminErr:    tc.schoolAdminErr,
				canTeacherResult:     tc.teacherResult,
				canTeacherErr:        tc.teacherErr,
				canParentResult:      tc.parentResult,
				canParentErr:         tc.parentErr,
			}
			svc := &ChatService{chatRepo: repo}

			ok, err := svc.CanCreateDirectConversation(context.Background(), tc.requesterID, tc.roles, tc.targetID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("CanCreateDirectConversation() error = %v", err)
			}
			if ok != tc.wantResult {
				t.Fatalf("allowed = %v, want %v", ok, tc.wantResult)
			}
			if repo.canSuperAdminCalls != tc.wantSuperCalls {
				t.Fatalf("super admin calls = %d, want %d", repo.canSuperAdminCalls, tc.wantSuperCalls)
			}
			if repo.canSchoolAdminCalls != tc.wantSchoolAdminCalls {
				t.Fatalf("school admin calls = %d, want %d", repo.canSchoolAdminCalls, tc.wantSchoolAdminCalls)
			}
			if repo.canTeacherCalls != tc.wantTeacherCalls {
				t.Fatalf("teacher calls = %d, want %d", repo.canTeacherCalls, tc.wantTeacherCalls)
			}
			if repo.canParentCalls != tc.wantParentCalls {
				t.Fatalf("parent calls = %d, want %d", repo.canParentCalls, tc.wantParentCalls)
			}
		})
	}
}

func TestChatServiceCreateGroupConversation(t *testing.T) {
	repo := &fakeChatRepo{createConversationResult: &model.Conversation{ConversationID: uuid.New(), Type: "group", CreatedAt: time.Now().UTC()}}
	svc := &ChatService{chatRepo: repo}

	_, err := svc.CreateGroupConversation(context.Background(), "group", []uuid.UUID{uuid.New()})
	if !errors.Is(err, ErrChatGroupNeedMembers) {
		t.Fatalf("error = %v, want %v", err, ErrChatGroupNeedMembers)
	}
	if repo.createConversationCalls != 0 {
		t.Fatalf("create calls = %d, want 0", repo.createConversationCalls)
	}

	name := "parents"
	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	conv, err := svc.CreateGroupConversation(context.Background(), name, ids)
	if err != nil {
		t.Fatalf("CreateGroupConversation() error = %v", err)
	}
	if conv == nil || conv.Type != "group" {
		t.Fatalf("unexpected conversation returned: %#v", conv)
	}
	if repo.createConversationCalls != 1 {
		t.Fatalf("create calls = %d, want 1", repo.createConversationCalls)
	}
	if repo.createConversationType != "group" {
		t.Fatalf("create type = %q, want group", repo.createConversationType)
	}
	if repo.createConversationName == nil || *repo.createConversationName != name {
		t.Fatalf("group name was not forwarded")
	}
	if len(repo.createConversationParticipantIDs) != len(ids) {
		t.Fatalf("participants len = %d, want %d", len(repo.createConversationParticipantIDs), len(ids))
	}
}

func TestChatServiceListConversations(t *testing.T) {
	userID := uuid.New()
	convA := model.Conversation{ConversationID: uuid.New(), Type: "direct", CreatedAt: time.Now().UTC()}
	convB := model.Conversation{ConversationID: uuid.New(), Type: "group", CreatedAt: time.Now().UTC()}
	sentinelErr := errors.New("repo failed")

	t.Run("list conversations error", func(t *testing.T) {
		repo := &fakeChatRepo{listConversationsErr: sentinelErr}
		svc := &ChatService{chatRepo: repo}
		_, err := svc.ListConversations(context.Background(), userID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("list participants error", func(t *testing.T) {
		repo := &fakeChatRepo{listConversationsResult: []model.Conversation{convA}, listParticipantsErr: sentinelErr}
		svc := &ChatService{chatRepo: repo}
		_, err := svc.ListConversations(context.Background(), userID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeChatRepo{
			listConversationsResult: []model.Conversation{convA, convB},
			listParticipantsResult: map[uuid.UUID][]model.ParticipantInfo{
				convA.ConversationID: {{UserID: uuid.New(), Email: "a@example.com", FullName: "A"}},
			},
		}
		svc := &ChatService{chatRepo: repo}
		got, err := svc.ListConversations(context.Background(), userID)
		if err != nil {
			t.Fatalf("ListConversations() error = %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("conversations len = %d, want 2", len(got))
		}
		if len(got[0].Participants) != 1 {
			t.Fatalf("first conversation participants = %d, want 1", len(got[0].Participants))
		}
		if len(got[1].Participants) != 0 {
			t.Fatalf("second conversation participants = %d, want 0", len(got[1].Participants))
		}
		if len(repo.listParticipantsConversationIDs) != 2 || repo.listParticipantsConversationIDs[0] != convA.ConversationID || repo.listParticipantsConversationIDs[1] != convB.ConversationID {
			t.Fatalf("conversation IDs forwarded mismatch: %#v", repo.listParticipantsConversationIDs)
		}
	})
}

func TestChatServiceSendMessage(t *testing.T) {
	conversationID := uuid.New()
	senderID := uuid.New()
	msg := &model.Message{MessageID: uuid.New(), ConversationID: conversationID, SenderID: senderID, Content: "hello", CreatedAt: time.Now().UTC()}
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name             string
		content          string
		isParticipant    bool
		isParticipantErr error
		createErr        error
		wantErr          error
		wantCreateCalls  int
	}{
		{name: "empty content", content: "", wantErr: ErrChatEmptyMessage},
		{name: "is participant repo error", content: "hello", isParticipantErr: sentinelErr, wantErr: sentinelErr},
		{name: "not participant", content: "hello", isParticipant: false, wantErr: ErrChatNotParticipant},
		{name: "create message error", content: "hello", isParticipant: true, createErr: sentinelErr, wantErr: sentinelErr, wantCreateCalls: 1},
		{name: "success", content: "hello", isParticipant: true, wantCreateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeChatRepo{isParticipantResult: tc.isParticipant, isParticipantErr: tc.isParticipantErr, createMessageResult: msg, createMessageErr: tc.createErr}
			svc := &ChatService{chatRepo: repo}
			got, err := svc.SendMessage(context.Background(), conversationID, senderID, tc.content)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("SendMessage() error = %v", err)
			}
			if tc.wantErr == nil && got == nil {
				t.Fatalf("expected message on success")
			}
			if repo.createMessageCalls != tc.wantCreateCalls {
				t.Fatalf("create calls = %d, want %d", repo.createMessageCalls, tc.wantCreateCalls)
			}
		})
	}
}

func TestChatServiceListMessages(t *testing.T) {
	conversationID := uuid.New()
	userID := uuid.New()
	beforeID := uuid.New()
	messages := []model.MessageWithSender{{MessageID: uuid.New(), ConversationID: conversationID, SenderID: userID, SenderEmail: "a@example.com", Content: "hi", CreatedAt: time.Now().UTC()}}
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name             string
		before           *uuid.UUID
		limit            int
		isParticipant    bool
		isParticipantErr error
		listErr          error
		wantErr          error
		wantLimit        int
		wantListCalls    int
	}{
		{name: "participant check error", limit: 20, isParticipantErr: sentinelErr, wantErr: sentinelErr},
		{name: "not participant", limit: 20, isParticipant: false, wantErr: ErrChatNotParticipant},
		{name: "default limit when zero", limit: 0, isParticipant: true, wantLimit: 50, wantListCalls: 1},
		{name: "default limit when too high", limit: 101, isParticipant: true, wantLimit: 50, wantListCalls: 1},
		{name: "valid limit preserved", before: &beforeID, limit: 30, isParticipant: true, wantLimit: 30, wantListCalls: 1},
		{name: "list messages error", limit: 10, isParticipant: true, listErr: sentinelErr, wantErr: sentinelErr, wantLimit: 10, wantListCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeChatRepo{isParticipantResult: tc.isParticipant, isParticipantErr: tc.isParticipantErr, listMessagesResult: messages, listMessagesErr: tc.listErr}
			svc := &ChatService{chatRepo: repo}
			got, err := svc.ListMessages(context.Background(), conversationID, userID, tc.before, tc.limit)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ListMessages() error = %v", err)
			}
			if tc.wantErr == nil && len(got) != len(messages) {
				t.Fatalf("messages len = %d, want %d", len(got), len(messages))
			}
			if repo.listMessagesCalls != tc.wantListCalls {
				t.Fatalf("list calls = %d, want %d", repo.listMessagesCalls, tc.wantListCalls)
			}
			if tc.wantListCalls > 0 {
				if repo.listMessagesLimit != tc.wantLimit {
					t.Fatalf("limit forwarded = %d, want %d", repo.listMessagesLimit, tc.wantLimit)
				}
				if tc.before == nil && repo.listMessagesBefore != nil {
					t.Fatalf("before should be nil")
				}
				if tc.before != nil {
					if repo.listMessagesBefore == nil || *repo.listMessagesBefore != *tc.before {
						t.Fatalf("before pointer was not forwarded")
					}
				}
			}
		})
	}
}

func TestChatServiceGetParticipantIDs(t *testing.T) {
	conversationID := uuid.New()
	sentinelErr := errors.New("repo failed")

	t.Run("repo error", func(t *testing.T) {
		repo := &fakeChatRepo{getParticipantsErr: sentinelErr}
		svc := &ChatService{chatRepo: repo}
		_, err := svc.GetParticipantIDs(context.Background(), conversationID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		idA := uuid.New()
		idB := uuid.New()
		repo := &fakeChatRepo{getParticipantsResult: []model.ParticipantInfo{{UserID: idA}, {UserID: idB}}}
		svc := &ChatService{chatRepo: repo}
		ids, err := svc.GetParticipantIDs(context.Background(), conversationID)
		if err != nil {
			t.Fatalf("GetParticipantIDs() error = %v", err)
		}
		if len(ids) != 2 || ids[0] != idA || ids[1] != idB {
			t.Fatalf("participant IDs = %#v", ids)
		}
	})
}

func TestChatServiceSearchUsers(t *testing.T) {
	requesterID := uuid.New()
	result := []model.ParticipantInfo{{UserID: uuid.New(), Email: "x@example.com", FullName: "X"}}
	sentinelErr := errors.New("repo failed")

	t.Run("empty keyword returns empty result without repo calls", func(t *testing.T) {
		repo := &fakeChatRepo{}
		svc := &ChatService{chatRepo: repo}
		got, err := svc.SearchUsers(context.Background(), requesterID, []string{"SUPER_ADMIN"}, "")
		if err != nil {
			t.Fatalf("SearchUsers() error = %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("result len = %d, want 0", len(got))
		}
		if repo.searchGlobalCalls != 0 && repo.searchSchoolAdminCalls != 0 && repo.searchTeacherCalls != 0 && repo.searchParentCalls != 0 {
			t.Fatalf("search repo should not be called for empty keyword")
		}
	})

	t.Run("super admin", func(t *testing.T) {
		repo := &fakeChatRepo{searchGlobalResult: result}
		svc := &ChatService{chatRepo: repo}
		got, err := svc.SearchUsers(context.Background(), requesterID, []string{"SUPER_ADMIN"}, "alice")
		if err != nil {
			t.Fatalf("SearchUsers() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("result len = %d, want 1", len(got))
		}
		if repo.searchGlobalCalls != 1 || repo.searchGlobalKeyword != "alice" || repo.searchGlobalLimit != 10 {
			t.Fatalf("super admin search arguments mismatch")
		}
	})

	t.Run("school admin error", func(t *testing.T) {
		repo := &fakeChatRepo{searchSchoolAdminErr: sentinelErr}
		svc := &ChatService{chatRepo: repo}
		_, err := svc.SearchUsers(context.Background(), requesterID, []string{"SCHOOL_ADMIN"}, "alice")
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("teacher", func(t *testing.T) {
		repo := &fakeChatRepo{searchTeacherResult: result}
		svc := &ChatService{chatRepo: repo}
		got, err := svc.SearchUsers(context.Background(), requesterID, []string{"TEACHER"}, "alice")
		if err != nil {
			t.Fatalf("SearchUsers() error = %v", err)
		}
		if len(got) != 1 || repo.searchTeacherCalls != 1 || repo.searchTeacherUserID != requesterID || repo.searchTeacherLimit != 10 {
			t.Fatalf("teacher search mismatch")
		}
	})

	t.Run("parent", func(t *testing.T) {
		repo := &fakeChatRepo{searchParentResult: result}
		svc := &ChatService{chatRepo: repo}
		got, err := svc.SearchUsers(context.Background(), requesterID, []string{"PARENT"}, "alice")
		if err != nil {
			t.Fatalf("SearchUsers() error = %v", err)
		}
		if len(got) != 1 || repo.searchParentCalls != 1 || repo.searchParentUserID != requesterID || repo.searchParentLimit != 10 {
			t.Fatalf("parent search mismatch")
		}
	})

	t.Run("unknown role", func(t *testing.T) {
		repo := &fakeChatRepo{}
		svc := &ChatService{chatRepo: repo}
		got, err := svc.SearchUsers(context.Background(), requesterID, []string{"STUDENT"}, "alice")
		if err != nil {
			t.Fatalf("SearchUsers() error = %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("result len = %d, want 0", len(got))
		}
	})

	t.Run("role priority uses super admin", func(t *testing.T) {
		repo := &fakeChatRepo{searchGlobalResult: result}
		svc := &ChatService{chatRepo: repo}
		_, err := svc.SearchUsers(context.Background(), requesterID, []string{"TEACHER", "SUPER_ADMIN"}, "alice")
		if err != nil {
			t.Fatalf("SearchUsers() error = %v", err)
		}
		if repo.searchGlobalCalls != 1 || repo.searchTeacherCalls != 0 {
			t.Fatalf("role priority mismatch: super=%d teacher=%d", repo.searchGlobalCalls, repo.searchTeacherCalls)
		}
	})
}
