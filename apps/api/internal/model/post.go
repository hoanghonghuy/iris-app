package model

import (
	"time"

	"github.com/google/uuid"
)

// Post = một bài đăng trong hệ thống bảng tin.
type Post struct {
	PostID       uuid.UUID  `json:"post_id"`
	AuthorUserID uuid.UUID  `json:"author_user_id"`
	ScopeType    string     `json:"scope_type"` // school|class|student
	SchoolID     *uuid.UUID `json:"school_id,omitempty"`
	ClassID      *uuid.UUID `json:"class_id,omitempty"`
	StudentID    *uuid.UUID `json:"student_id,omitempty"`
	Type         string     `json:"type"` // announcement|activity|daily_note|health_note
	Content      string     `json:"content"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// PostAttachment = một tệp đính kèm được liên kết với bài đăng
type PostAttachment struct {
	AttachmentID uuid.UUID `json:"attachment_id"`
	PostID       uuid.UUID `json:"post_id"`
	URL          string    `json:"url"`
	MimeType     *string   `json:"mime_type,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
