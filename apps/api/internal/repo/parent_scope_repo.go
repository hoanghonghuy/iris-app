package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ParentScopeRepo struct {
	pool *pgxpool.Pool
}

// scannableRows trừu tượng hóa kiểu rows để helper scan dùng lại cho nhiều query feed/post.
type scannableRows interface {
	Next() bool
	Scan(...any) error
	Err() error
}

func (r *ParentScopeRepo) CountMyChildren(ctx context.Context, parentUserID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(*)
		FROM student_parents sp
		JOIN parents p ON p.parent_id = sp.parent_id
		WHERE p.user_id = $1;
	`
	var count int
	err := r.pool.QueryRow(ctx, q, parentUserID).Scan(&count)
	return count, err
}

func (r *ParentScopeRepo) CountMyRecentPosts(ctx context.Context, parentUserID uuid.UUID, since time.Time) (int, error) {
	const q = `
		SELECT COUNT(DISTINCT p.post_id)
		FROM posts p
		JOIN student_parents sp ON (
			(p.scope_type = 'student' AND p.student_id = sp.student_id)
			OR (p.scope_type = 'class' AND p.class_id IN (
				SELECT s.current_class_id
				FROM students s
				WHERE s.student_id = sp.student_id
			))
		)
		JOIN parents pa ON pa.parent_id = sp.parent_id
		WHERE pa.user_id = $1
		  AND p.created_at >= $2;
	`
	var count int
	err := r.pool.QueryRow(ctx, q, parentUserID, since).Scan(&count)
	return count, err
}

func (r *ParentScopeRepo) CountMyRecentHealthAlerts(ctx context.Context, parentUserID uuid.UUID, since time.Time) (int, error) {
	const q = `
		SELECT COUNT(*)
		FROM health_logs h
		JOIN student_parents sp ON sp.student_id = h.student_id
		JOIN parents p ON p.parent_id = sp.parent_id
		WHERE p.user_id = $1
		  AND h.recorded_at >= $2
		  AND h.severity IN ('watch', 'urgent');
	`
	var count int
	err := r.pool.QueryRow(ctx, q, parentUserID, since).Scan(&count)
	return count, err
}

func NewParentScopeRepo(pool *pgxpool.Pool) *ParentScopeRepo {
	return &ParentScopeRepo{
		pool: pool,
	}
}

// ListMyChildren liệt kê các học sinh (con) của phụ huynh theo user_id
func (r *ParentScopeRepo) ListMyChildren(ctx context.Context, parentUserID uuid.UUID) ([]model.Student, error) {
	const q = `
		SELECT s.student_id, s.school_id, s.current_class_id, c.name AS current_class_name, s.full_name, s.dob, s.gender
		FROM students s
		JOIN classes c ON c.class_id = s.current_class_id
		JOIN student_parents sp ON sp.student_id = s.student_id
		JOIN parents p ON p.parent_id = sp.parent_id
		WHERE p.user_id = $1
		ORDER BY s.full_name;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.StudentID, &s.SchoolID, &s.CurrentClassID, &s.CurrentClassName, &s.FullName, &s.DOB, &s.Gender); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

// ListMyChildClassPosts liệt kê bài đăng của lớp con của phụ huynh đang học
func (r *ParentScopeRepo) ListMyChildClassPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content,
			COALESCE(lc.like_count, 0) AS like_count,
			COALESCE(cc.comment_count, 0) AS comment_count,
			COALESCE(sc.share_count, 0) AS share_count,
			EXISTS(
				SELECT 1
				FROM post_interactions self_like
				WHERE self_like.post_id = p.post_id AND self_like.user_id = $1 AND self_like.action_type = 'like'
			) AS liked_by_me,
			p.created_at, p.updated_at,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN students s ON s.current_class_id = p.class_id
		JOIN student_parents sp ON sp.student_id = s.student_id
		JOIN parents pa ON pa.parent_id = sp.parent_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS like_count
			FROM post_interactions
			WHERE action_type = 'like'
			GROUP BY post_id
		) lc ON lc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM post_comments
			GROUP BY post_id
		) cc ON cc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS share_count
			FROM post_interactions
			WHERE action_type = 'share'
			GROUP BY post_id
		) sc ON sc.post_id = p.post_id
		WHERE pa.user_id = $1
			AND s.student_id = $2
			AND p.scope_type = 'class'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return scanPostsWithTotal(rows)
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con phụ huynh (student scope)
func (r *ParentScopeRepo) ListMyChildStudentPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content,
			COALESCE(lc.like_count, 0) AS like_count,
			COALESCE(cc.comment_count, 0) AS comment_count,
			COALESCE(sc.share_count, 0) AS share_count,
			EXISTS(
				SELECT 1
				FROM post_interactions self_like
				WHERE self_like.post_id = p.post_id AND self_like.user_id = $1 AND self_like.action_type = 'like'
			) AS liked_by_me,
			p.created_at, p.updated_at,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN student_parents sp ON sp.student_id = p.student_id
		JOIN parents pa ON pa.parent_id = sp.parent_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS like_count
			FROM post_interactions
			WHERE action_type = 'like'
			GROUP BY post_id
		) lc ON lc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM post_comments
			GROUP BY post_id
		) cc ON cc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS share_count
			FROM post_interactions
			WHERE action_type = 'share'
			GROUP BY post_id
		) sc ON sc.post_id = p.post_id
		WHERE pa.user_id = $1
			AND p.student_id = $2
			AND p.scope_type = 'student'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return scanPostsWithTotal(rows)
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con phụ huynh (cả class và student scope)
func (r *ParentScopeRepo) ListAllMyChildPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content,
			COALESCE(lc.like_count, 0) AS like_count,
			COALESCE(cc.comment_count, 0) AS comment_count,
			COALESCE(sc.share_count, 0) AS share_count,
			EXISTS(
				SELECT 1
				FROM post_interactions self_like
				WHERE self_like.post_id = p.post_id AND self_like.user_id = $1 AND self_like.action_type = 'like'
			) AS liked_by_me,
			p.created_at, p.updated_at,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN student_parents sp ON sp.student_id = $2
		JOIN parents pa ON pa.parent_id = sp.parent_id
		JOIN students s ON s.student_id = sp.student_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS like_count
			FROM post_interactions
			WHERE action_type = 'like'
			GROUP BY post_id
		) lc ON lc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM post_comments
			GROUP BY post_id
		) cc ON cc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS share_count
			FROM post_interactions
			WHERE action_type = 'share'
			GROUP BY post_id
		) sc ON sc.post_id = p.post_id
		WHERE pa.user_id = $1
			AND (
				(p.scope_type = 'class' AND p.class_id = s.current_class_id)
				OR (p.scope_type = 'student' AND p.student_id = $2)
			)
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return scanPostsWithTotal(rows)
}

// GetMyFeed lấy tất cả bài đăng liên quan đến con phụ huynh, sắp xếp theo thời gian, có phân trang.
func (r *ParentScopeRepo) GetMyFeed(ctx context.Context, parentUserID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	const q = `
		SELECT DISTINCT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content,
			COALESCE(lc.like_count, 0) AS like_count,
			COALESCE(cc.comment_count, 0) AS comment_count,
			COALESCE(sc.share_count, 0) AS share_count,
			EXISTS(
				SELECT 1
				FROM post_interactions self_like
				WHERE self_like.post_id = p.post_id AND self_like.user_id = $1 AND self_like.action_type = 'like'
			) AS liked_by_me,
			p.created_at, p.updated_at,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN student_parents sp ON (
			(p.scope_type = 'student' AND p.student_id = sp.student_id)
			OR (p.scope_type = 'class' AND p.class_id IN (
				SELECT s.current_class_id
				FROM students s
				WHERE s.student_id = sp.student_id
			))
		)
		JOIN parents pa ON pa.parent_id = sp.parent_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS like_count
			FROM post_interactions
			WHERE action_type = 'like'
			GROUP BY post_id
		) lc ON lc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM post_comments
			GROUP BY post_id
		) cc ON cc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS share_count
			FROM post_interactions
			WHERE action_type = 'share'
			GROUP BY post_id
		) sc ON sc.post_id = p.post_id
		WHERE pa.user_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return scanPostsWithTotal(rows)
}

// scanPostsWithTotal gom logic Scan post + total_count cho các endpoint feed/list của parent.
func scanPostsWithTotal(rows scannableRows) ([]model.Post, int, error) {
	posts := make([]model.Post, 0)
	total := 0
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content,
			&p.LikeCount, &p.CommentCount, &p.ShareCount, &p.LikedByMe,
			&p.CreatedAt, &p.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}

	return posts, total, rows.Err()
}

// IsParentOfStudent kiểm tra xem user có phải là parent của student không
func (r *ParentScopeRepo) IsParentOfStudent(ctx context.Context, parentUserID, studentID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1 FROM student_parents sp
			JOIN parents p ON p.parent_id = sp.parent_id
			WHERE p.user_id = $1 AND sp.student_id = $2
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, q, parentUserID, studentID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
