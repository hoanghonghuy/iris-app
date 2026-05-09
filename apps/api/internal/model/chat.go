package model

import (
	"time"

	"github.com/google/uuid"
)

// Conversation đại diện cho một cuộc hội thoại (direct hoặc group)
type Conversation struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	Type           string    `json:"type"` // direct | group
	Name           *string   `json:"name,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ConversationParticipant đại diện cho một thành viên trong cuộc hội thoại
type ConversationParticipant struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
	JoinedAt       time.Time `json:"joined_at"`
}

// Message đại diện cho một tin nhắn trong cuộc hội thoại
type Message struct {
	MessageID      uuid.UUID `json:"message_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// MessageWithSender chứa thông tin tin nhắn kèm email người gửi (dùng cho API response)
type MessageWithSender struct {
	MessageID      uuid.UUID `json:"message_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	SenderEmail    string    `json:"sender_email"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// LastMessagePreview tin nhắn mới nhất (sidebar / list conversations).
type LastMessagePreview struct {
	MessageID   uuid.UUID `json:"message_id"`
	SenderID    uuid.UUID `json:"sender_id"`
	SenderEmail string    `json:"sender_email"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

// ConversationListSummary một hàng từ DB trước khi gắn participants (internal/repo → service).
type ConversationListSummary struct {
	Conversation
	LastMessage *LastMessagePreview
	UnreadCount int
}

// ConversationWithParticipants chứa thông tin cuộc hội thoại kèm danh sách thành viên
type ConversationWithParticipants struct {
	Conversation
	Participants []ParticipantInfo `json:"participants"`
	// LastMessage nil nếu chưa có tin nhắn.
	LastMessage *LastMessagePreview `json:"last_message,omitempty"`
	// UnreadCount số tin từ người khác chưa đọc (theo last_read_at / joined_at).
	UnreadCount int `json:"unread_count"`
}

// ParticipantInfo chứa thông tin cơ bản của thành viên trong cuộc hội thoại
type ParticipantInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
}
