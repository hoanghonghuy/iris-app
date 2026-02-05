package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ParentScopeRepo struct {
	pool *pgxpool.Pool
}

func NewParentScopeRepo(pool *pgxpool.Pool) *ParentScopeRepo {
	return &ParentScopeRepo{
		pool: pool,
	}
}

// ListMyChildren liệt kê các học sinh (con) của phụ huynh theo user_id
func (r *ParentScopeRepo) ListMyChildren(ctx context.Context, parentUserID uuid.UUID) ([]model.Student, error) {
	const q = `
		SELECT s.student_id, s.school_id, s.current_class_id, s.full_name, s.dob, s.gender
		FROM students s
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
		if err := rows.Scan(&s.ID, &s.SchoolID, &s.CurrentClassID, &s.FullName, &s.DOB, &s.Gender); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

// ListMyChildClassPosts liệt kê bài đăng của lớp con mình đang học
func (r *ParentScopeRepo) ListMyChildClassPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content, p.created_at, p.updated_at
		FROM posts p
		JOIN students s ON s.current_class_id = p.class_id
		JOIN student_parents sp ON sp.student_id = s.student_id
		JOIN parents pa ON pa.parent_id = sp.parent_id
		WHERE pa.user_id = $1
			AND s.student_id = $2
			AND p.scope_type = 'class'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con mình (student scope)
func (r *ParentScopeRepo) ListMyChildStudentPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content, p.created_at, p.updated_at
		FROM posts p
		JOIN student_parents sp ON sp.student_id = p.student_id
		JOIN parents pa ON pa.parent_id = sp.parent_id
		WHERE pa.user_id = $1
			AND p.student_id = $2
			AND p.scope_type = 'student'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con mình (cả class và student scope)
func (r *ParentScopeRepo) ListAllMyChildPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content, p.created_at, p.updated_at
		FROM posts p
		JOIN student_parents sp ON sp.student_id = $2
		JOIN parents pa ON pa.parent_id = sp.parent_id
		JOIN students s ON s.student_id = sp.student_id
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
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *ParentScopeRepo) GetMyFeed(ctx context.Context, parentUserID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT DISTINCT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content, p.created_at, p.updated_at
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
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, q, parentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
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
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
