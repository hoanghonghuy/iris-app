package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo struct {
	pool *pgxpool.Pool
}

func NewChatRepo(pool *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{pool: pool}
}

// CreateConversation tạo cuộc hội thoại mới và thêm danh sách thành viên
func (r *ChatRepo) CreateConversation(ctx context.Context, convType string, name *string, participantIDs []uuid.UUID) (*model.Conversation, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const qConv = `
		INSERT INTO conversations (type, name)
		VALUES ($1, $2)
		RETURNING conversation_id, type, name, created_at;
	`
	var conv model.Conversation
	err = tx.QueryRow(ctx, qConv, convType, name).Scan(
		&conv.ConversationID, &conv.Type, &conv.Name, &conv.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	const qPart = `
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2);
	`
	for _, uid := range participantIDs {
		if _, err := tx.Exec(ctx, qPart, conv.ConversationID, uid); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &conv, nil
}

// FindDirectConversation tìm cuộc hội thoại direct giữa 2 user (nếu đã tồn tại)
func (r *ChatRepo) FindDirectConversation(ctx context.Context, userA, userB uuid.UUID) (*model.Conversation, error) {
	const q = `
		SELECT c.conversation_id, c.type, c.name, c.created_at
		FROM conversations c
		WHERE c.type = 'direct'
		  AND c.conversation_id IN (
		    SELECT cp1.conversation_id
		    FROM conversation_participants cp1
		    JOIN conversation_participants cp2 ON cp1.conversation_id = cp2.conversation_id
		    WHERE cp1.user_id = $1 AND cp2.user_id = $2
		  );
	`
	var conv model.Conversation
	err := r.pool.QueryRow(ctx, q, userA, userB).Scan(
		&conv.ConversationID, &conv.Type, &conv.Name, &conv.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// ListConversationsByUser lấy danh sách cuộc hội thoại mà user tham gia
func (r *ChatRepo) ListConversationsByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	const q = `
		SELECT c.conversation_id, c.type, c.name, c.created_at
		FROM conversations c
		JOIN conversation_participants cp ON c.conversation_id = cp.conversation_id
		WHERE cp.user_id = $1
		ORDER BY c.created_at DESC;
	`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []model.Conversation
	for rows.Next() {
		var c model.Conversation
		if err := rows.Scan(&c.ConversationID, &c.Type, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, rows.Err()
}

// GetParticipants lấy danh sách thành viên của cuộc hội thoại kèm full_name
func (r *ChatRepo) GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]model.ParticipantInfo, error) {
	const q = `
		SELECT cp.user_id, u.email,
		       COALESCE(t.full_name, p.full_name, sa.full_name, 'Admin/Unknown') as full_name
		FROM conversation_participants cp
		JOIN users u ON cp.user_id = u.user_id
		LEFT JOIN teachers t ON t.user_id = u.user_id
		LEFT JOIN parents p ON p.user_id = u.user_id
		LEFT JOIN school_admins sa ON sa.user_id = u.user_id
		WHERE cp.conversation_id = $1
		ORDER BY cp.joined_at;
	`
	rows, err := r.pool.Query(ctx, q, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []model.ParticipantInfo
	for rows.Next() {
		var p model.ParticipantInfo
		if err := rows.Scan(&p.UserID, &p.Email, &p.FullName); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}

// IsParticipant kiểm tra user có thuộc cuộc hội thoại không
func (r *ChatRepo) IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1 FROM conversation_participants
			WHERE conversation_id = $1 AND user_id = $2
		);
	`
	var exists bool
	err := r.pool.QueryRow(ctx, q, conversationID, userID).Scan(&exists)
	return exists, err
}

// CanSuperAdminMessageTarget cho phép SUPER_ADMIN nhắn trực tiếp đến mọi active user (trừ chính mình).
func (r *ChatRepo) CanSuperAdminMessageTarget(ctx context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1
			FROM users u
			WHERE u.user_id = $2
			  AND u.status = 'active'
			  AND u.user_id <> $1
		);
	`
	var allowed bool
	err := r.pool.QueryRow(ctx, q, requesterID, targetID).Scan(&allowed)
	return allowed, err
}

// CanSchoolAdminMessageTarget cho phép SCHOOL_ADMIN nhắn user active trong cùng school.
func (r *ChatRepo) CanSchoolAdminMessageTarget(ctx context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1
			FROM users u
			JOIN school_admins sa_req ON sa_req.user_id = $1
			LEFT JOIN teachers t ON t.user_id = u.user_id
			LEFT JOIN parents p ON p.user_id = u.user_id
			LEFT JOIN school_admins sa2 ON sa2.user_id = u.user_id
			WHERE u.user_id = $2
			  AND u.status = 'active'
			  AND u.user_id <> $1
			  AND (
				t.school_id = sa_req.school_id
				OR p.school_id = sa_req.school_id
				OR sa2.school_id = sa_req.school_id
			  )
		);
	`
	var allowed bool
	err := r.pool.QueryRow(ctx, q, requesterID, targetID).Scan(&allowed)
	return allowed, err
}

// CanTeacherMessageTarget cho phép TEACHER nhắn teacher/admin cùng school hoặc parent của học sinh mình dạy.
func (r *ChatRepo) CanTeacherMessageTarget(ctx context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1
			FROM users u
			JOIN teachers t_req ON t_req.user_id = $1
			LEFT JOIN teachers t2 ON t2.user_id = u.user_id
			LEFT JOIN parents p ON p.user_id = u.user_id
			LEFT JOIN school_admins sa ON sa.user_id = u.user_id
			WHERE u.user_id = $2
			  AND u.status = 'active'
			  AND u.user_id <> $1
			  AND (
				t2.school_id = t_req.school_id
				OR sa.school_id = t_req.school_id
				OR p.parent_id IN (
					SELECT sp.parent_id
					FROM student_parents sp
					JOIN students s ON s.student_id = sp.student_id
					JOIN teacher_classes tc ON tc.class_id = s.current_class_id
					WHERE tc.teacher_id = t_req.teacher_id
				)
			  )
		);
	`
	var allowed bool
	err := r.pool.QueryRow(ctx, q, requesterID, targetID).Scan(&allowed)
	return allowed, err
}

// CanParentMessageTarget cho phép PARENT nhắn teacher của con mình hoặc school admin cùng school.
func (r *ChatRepo) CanParentMessageTarget(ctx context.Context, requesterID, targetID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1
			FROM users u
			JOIN parents p_req ON p_req.user_id = $1
			LEFT JOIN teachers t ON t.user_id = u.user_id
			LEFT JOIN school_admins sa ON sa.user_id = u.user_id
			WHERE u.user_id = $2
			  AND u.status = 'active'
			  AND u.user_id <> $1
			  AND (
				sa.school_id = p_req.school_id
				OR t.teacher_id IN (
					SELECT tc.teacher_id
					FROM teacher_classes tc
					JOIN students s ON s.current_class_id = tc.class_id
					JOIN student_parents sp ON sp.student_id = s.student_id
					WHERE sp.parent_id = p_req.parent_id
				)
			  )
		);
	`
	var allowed bool
	err := r.pool.QueryRow(ctx, q, requesterID, targetID).Scan(&allowed)
	return allowed, err
}

// CreateMessage lưu tin nhắn mới vào database
func (r *ChatRepo) CreateMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string) (*model.Message, error) {
	const q = `
		INSERT INTO messages (conversation_id, sender_id, content)
		VALUES ($1, $2, $3)
		RETURNING message_id, conversation_id, sender_id, content, created_at;
	`
	var msg model.Message
	err := r.pool.QueryRow(ctx, q, conversationID, senderID, content).Scan(
		&msg.MessageID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ListMessages lấy danh sách tin nhắn theo cursor (before message_id).
// Trả về tối đa `limit` tin nhắn, sắp xếp DESC (mới nhất trước).
// Frontend đảo ngược slice trước khi hiển thị, dùng message_id cuối làm next_cursor.
func (r *ChatRepo) ListMessages(ctx context.Context, conversationID uuid.UUID, before *uuid.UUID, limit int) ([]model.MessageWithSender, error) {
	if before == nil {
		return r.listLatestMessages(ctx, conversationID, limit)
	}
	return r.listMessagesBefore(ctx, conversationID, *before, limit)
}

func (r *ChatRepo) listLatestMessages(ctx context.Context, conversationID uuid.UUID, limit int) ([]model.MessageWithSender, error) {
	const q = `
		SELECT m.message_id, m.conversation_id, m.sender_id, u.email, m.content, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.user_id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC, m.message_id DESC
		LIMIT $2;
	`
	rows, err := r.pool.Query(ctx, q, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessages(rows)
}

func (r *ChatRepo) listMessagesBefore(ctx context.Context, conversationID, beforeID uuid.UUID, limit int) ([]model.MessageWithSender, error) {
	// Composite cursor: lấy tin nhắn cũ hơn cursor (created_at, message_id)
	const q = `
		SELECT m.message_id, m.conversation_id, m.sender_id, u.email, m.content, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.user_id
		WHERE m.conversation_id = $1
		  AND (m.created_at, m.message_id) < (
		      SELECT created_at, message_id FROM messages WHERE message_id = $2
		  )
		ORDER BY m.created_at DESC, m.message_id DESC
		LIMIT $3;
	`
	rows, err := r.pool.Query(ctx, q, conversationID, beforeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessages(rows)
}

// scanMessages đọc rows từ truy vấn messages
func scanMessages(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]model.MessageWithSender, error) {
	var msgs []model.MessageWithSender
	for rows.Next() {
		var m model.MessageWithSender
		if err := rows.Scan(&m.MessageID, &m.ConversationID, &m.SenderID, &m.SenderEmail, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// SearchUsersGlobal tìm kiếm toàn bộ hệ thống (dành cho SUPER_ADMIN)
func (r *ChatRepo) SearchUsersGlobal(ctx context.Context, keyword string, limit int) ([]model.ParticipantInfo, error) {
	const q = `
		SELECT u.user_id, u.email,
		       COALESCE(t.full_name, p.full_name, sa.full_name, 'Unknown') as full_name
		FROM users u
		LEFT JOIN teachers t ON t.user_id = u.user_id
		LEFT JOIN parents p ON p.user_id = u.user_id
		LEFT JOIN school_admins sa ON sa.user_id = u.user_id
		WHERE u.status = 'active'
		  AND (
		      u.email ILIKE '%' || $1 || '%'
		      OR t.full_name ILIKE '%' || $1 || '%'
		      OR p.full_name ILIKE '%' || $1 || '%'
		      OR sa.full_name ILIKE '%' || $1 || '%'
		  )
		LIMIT $2;
	`
	rows, err := r.pool.Query(ctx, q, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []model.ParticipantInfo
	for rows.Next() {
		var p model.ParticipantInfo
		if err := rows.Scan(&p.UserID, &p.Email, &p.FullName); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}

// SearchUsersForSchoolAdmin tìm kiếm giáo viên và phụ huynh trong cùng trường
func (r *ChatRepo) SearchUsersForSchoolAdmin(ctx context.Context, adminID uuid.UUID, keyword string, limit int) ([]model.ParticipantInfo, error) {
	const q = `
		SELECT u.user_id, u.email,
		       COALESCE(t.full_name, p.full_name, sa2.full_name, 'Unknown') as full_name
		FROM users u
		LEFT JOIN teachers t ON t.user_id = u.user_id
		LEFT JOIN parents p ON p.user_id = u.user_id
		LEFT JOIN school_admins sa2 ON sa2.user_id = u.user_id
		JOIN school_admins sa_req ON sa_req.user_id = $1
		WHERE u.status = 'active'
		  AND u.user_id != $1
		  AND (
		      t.school_id = sa_req.school_id
		      OR p.school_id = sa_req.school_id
		      OR sa2.school_id = sa_req.school_id
		  )
		  AND (
		      u.email ILIKE '%' || $2 || '%'
		      OR t.full_name ILIKE '%' || $2 || '%'
		      OR p.full_name ILIKE '%' || $2 || '%'
		      OR sa2.full_name ILIKE '%' || $2 || '%'
		  )
		LIMIT $3;
	`
	rows, err := r.pool.Query(ctx, q, adminID, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []model.ParticipantInfo
	for rows.Next() {
		var p model.ParticipantInfo
		if err := rows.Scan(&p.UserID, &p.Email, &p.FullName); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}

// SearchUsersForTeacher tìm kiếm phụ huynh (của hs mà gv hiện tại dạy), giáo viên & admin cùng trường
func (r *ChatRepo) SearchUsersForTeacher(ctx context.Context, teacherID uuid.UUID, keyword string, limit int) ([]model.ParticipantInfo, error) {
	const q = `
		SELECT u.user_id, u.email,
		       COALESCE(t2.full_name, p.full_name, sa.full_name, 'Unknown') as full_name
		FROM users u
		LEFT JOIN teachers t2 ON t2.user_id = u.user_id
		LEFT JOIN parents p ON p.user_id = u.user_id
		LEFT JOIN school_admins sa ON sa.user_id = u.user_id
		JOIN teachers t_req ON t_req.user_id = $1
		WHERE u.status = 'active'
		  AND u.user_id != $1
		  AND (
		      u.email ILIKE '%' || $2 || '%'
		      OR t2.full_name ILIKE '%' || $2 || '%'
		      OR p.full_name ILIKE '%' || $2 || '%'
		      OR sa.full_name ILIKE '%' || $2 || '%'
		      OR EXISTS (
		         SELECT 1 FROM students s 
		         JOIN student_parents sp ON s.student_id = sp.student_id 
		         WHERE sp.parent_id = p.parent_id AND s.full_name ILIKE '%' || $2 || '%'
		      )
		  )
		  AND (
		      t2.school_id = t_req.school_id
		      OR sa.school_id = t_req.school_id
		      OR p.parent_id IN (
		          SELECT sp.parent_id 
		          FROM student_parents sp
		          JOIN students s ON s.student_id = sp.student_id
		          JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		          WHERE tc.teacher_id = t_req.teacher_id
		      )
		  )
		LIMIT $3;
	`
	rows, err := r.pool.Query(ctx, q, teacherID, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []model.ParticipantInfo
	for rows.Next() {
		var p model.ParticipantInfo
		if err := rows.Scan(&p.UserID, &p.Email, &p.FullName); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}

// SearchUsersForParent tìm kiếm giáo viên (dạy con của parent hiện tại) & admin cùng trường
func (r *ChatRepo) SearchUsersForParent(ctx context.Context, parentID uuid.UUID, keyword string, limit int) ([]model.ParticipantInfo, error) {
	const q = `
		SELECT u.user_id, u.email,
		       COALESCE(t.full_name, sa.full_name, 'Unknown') as full_name
		FROM users u
		LEFT JOIN teachers t ON t.user_id = u.user_id
		LEFT JOIN school_admins sa ON sa.user_id = u.user_id
		JOIN parents p_req ON p_req.user_id = $1
		WHERE u.status = 'active'
		  AND u.user_id != $1
		  AND (
		      u.email ILIKE '%' || $2 || '%'
		      OR t.full_name ILIKE '%' || $2 || '%'
		      OR sa.full_name ILIKE '%' || $2 || '%'
		      OR EXISTS (
		         SELECT 1 FROM students s 
		         JOIN teacher_classes tc ON s.current_class_id = tc.class_id 
		         WHERE tc.teacher_id = t.teacher_id 
		           AND s.full_name ILIKE '%' || $2 || '%'
		           AND EXISTS (
		               SELECT 1 FROM student_parents sp WHERE sp.student_id = s.student_id AND sp.parent_id = p_req.parent_id
		           )
		      )
		  )
		  AND (
		      sa.school_id = p_req.school_id
		      OR t.teacher_id IN (
		          SELECT tc.teacher_id
		          FROM teacher_classes tc
		          JOIN students s ON s.current_class_id = tc.class_id
		          JOIN student_parents sp ON sp.student_id = s.student_id
		          WHERE sp.parent_id = p_req.parent_id
		      )
		  )
		LIMIT $3;
	`
	rows, err := r.pool.Query(ctx, q, parentID, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []model.ParticipantInfo
	for rows.Next() {
		var p model.ParticipantInfo
		if err := rows.Scan(&p.UserID, &p.Email, &p.FullName); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}
