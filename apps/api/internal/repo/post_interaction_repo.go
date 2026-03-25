package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostInteractionRepo struct {
	pool *pgxpool.Pool
}

func NewPostInteractionRepo(pool *pgxpool.Pool) *PostInteractionRepo {
	return &PostInteractionRepo{pool: pool}
}

func (r *PostInteractionRepo) TeacherCanAccessPost(ctx context.Context, teacherUserID, postID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM posts p
			WHERE p.post_id = $2
			  AND (
				(p.scope_type = 'class' AND EXISTS (
					SELECT 1
					FROM teacher_classes tc
					JOIN teachers t ON t.teacher_id = tc.teacher_id
					WHERE t.user_id = $1 AND tc.class_id = p.class_id
				))
				OR
				(p.scope_type = 'student' AND EXISTS (
					SELECT 1
					FROM students s
					JOIN teacher_classes tc ON tc.class_id = s.current_class_id
					JOIN teachers t ON t.teacher_id = tc.teacher_id
					WHERE t.user_id = $1 AND s.student_id = p.student_id
				))
			  )
		)
	`

	var ok bool
	err := r.pool.QueryRow(ctx, q, teacherUserID, postID).Scan(&ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *PostInteractionRepo) ParentCanAccessPost(ctx context.Context, parentUserID, postID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM posts p
			WHERE p.post_id = $2
			  AND (
				(p.scope_type = 'class' AND EXISTS (
					SELECT 1
					FROM students s
					JOIN student_parents sp ON sp.student_id = s.student_id
					JOIN parents pa ON pa.parent_id = sp.parent_id
					WHERE pa.user_id = $1 AND s.current_class_id = p.class_id
				))
				OR
				(p.scope_type = 'student' AND EXISTS (
					SELECT 1
					FROM student_parents sp
					JOIN parents pa ON pa.parent_id = sp.parent_id
					WHERE pa.user_id = $1 AND sp.student_id = p.student_id
				))
				OR
				(p.scope_type = 'school' AND EXISTS (
					SELECT 1
					FROM students s
					JOIN student_parents sp ON sp.student_id = s.student_id
					JOIN parents pa ON pa.parent_id = sp.parent_id
					WHERE pa.user_id = $1 AND s.school_id = p.school_id
				))
			  )
		)
	`

	var ok bool
	err := r.pool.QueryRow(ctx, q, parentUserID, postID).Scan(&ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *PostInteractionRepo) ToggleLike(ctx context.Context, userID, postID uuid.UUID) (bool, int, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const qDelete = `
		DELETE FROM post_interactions
		WHERE post_id = $1 AND user_id = $2 AND action_type = 'like'
	`
	deleteTag, err := tx.Exec(ctx, qDelete, postID, userID)
	if err != nil {
		return false, 0, err
	}

	liked := false
	if deleteTag.RowsAffected() == 0 {
		const qInsert = `
			INSERT INTO post_interactions (post_id, user_id, action_type)
			VALUES ($1, $2, 'like')
			ON CONFLICT (post_id, user_id, action_type) WHERE action_type = 'like' DO NOTHING
		`
		if _, err := tx.Exec(ctx, qInsert, postID, userID); err != nil {
			return false, 0, err
		}
		liked = true
	}

	const qCount = `
		SELECT COUNT(*)
		FROM post_interactions
		WHERE post_id = $1 AND action_type = 'like'
	`
	var likeCount int
	if err := tx.QueryRow(ctx, qCount, postID).Scan(&likeCount); err != nil {
		return false, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, 0, err
	}

	return liked, likeCount, nil
}

func (r *PostInteractionRepo) AddComment(ctx context.Context, userID, postID uuid.UUID, content string) (model.PostComment, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO post_comments (post_id, author_user_id, content)
			VALUES ($1, $2, $3)
			RETURNING comment_id, post_id, author_user_id, content, created_at
		)
		SELECT i.comment_id, i.post_id, i.author_user_id,
			COALESCE(t.full_name, pa.full_name, u.email) AS author_display,
			i.content, i.created_at
		FROM inserted i
		JOIN users u ON u.user_id = i.author_user_id
		LEFT JOIN teachers t ON t.user_id = i.author_user_id
		LEFT JOIN parents pa ON pa.user_id = i.author_user_id
	`

	var out model.PostComment
	err := r.pool.QueryRow(ctx, q, postID, userID, strings.TrimSpace(content)).Scan(
		&out.CommentID,
		&out.PostID,
		&out.AuthorUserID,
		&out.AuthorDisplay,
		&out.Content,
		&out.CreatedAt,
	)
	if err != nil {
		return model.PostComment{}, err
	}

	return out, nil
}

func (r *PostInteractionRepo) ListComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.PostComment, int, error) {
	const q = `
		SELECT c.comment_id, c.post_id, c.author_user_id,
			COALESCE(t.full_name, pa.full_name, u.email) AS author_display,
			c.content, c.created_at,
			COUNT(*) OVER() AS total_count
		FROM post_comments c
		JOIN users u ON u.user_id = c.author_user_id
		LEFT JOIN teachers t ON t.user_id = c.author_user_id
		LEFT JOIN parents pa ON pa.user_id = c.author_user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, q, postID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.PostComment, 0, limit)
	total := 0
	for rows.Next() {
		var item model.PostComment
		if err := rows.Scan(
			&item.CommentID,
			&item.PostID,
			&item.AuthorUserID,
			&item.AuthorDisplay,
			&item.Content,
			&item.CreatedAt,
			&total,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, rows.Err()
}

func (r *PostInteractionRepo) AddShare(ctx context.Context, userID, postID uuid.UUID) (int, error) {
	const qInsert = `
		INSERT INTO post_interactions (post_id, user_id, action_type)
		VALUES ($1, $2, 'share')
	`
	if _, err := r.pool.Exec(ctx, qInsert, postID, userID); err != nil {
		return 0, err
	}

	const qCount = `
		SELECT COUNT(*)
		FROM post_interactions
		WHERE post_id = $1 AND action_type = 'share'
	`
	var shareCount int
	if err := r.pool.QueryRow(ctx, qCount, postID).Scan(&shareCount); err != nil {
		return 0, err
	}

	return shareCount, nil
}

func (r *PostInteractionRepo) CountComments(ctx context.Context, postID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(*)
		FROM post_comments
		WHERE post_id = $1
	`
	var count int
	if err := r.pool.QueryRow(ctx, q, postID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostInteractionRepo) EnsurePostExists(ctx context.Context, postID uuid.UUID) error {
	const q = `
		SELECT 1
		FROM posts
		WHERE post_id = $1
	`
	var one int
	err := r.pool.QueryRow(ctx, q, postID).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRowsUpdated
		}
		return err
	}
	return nil
}
