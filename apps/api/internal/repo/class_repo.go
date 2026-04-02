package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClassRepo struct {
	pool *pgxpool.Pool
}

func NewClassRepo(pool *pgxpool.Pool) *ClassRepo {
	return &ClassRepo{
		pool: pool,
	}
}

func (r *ClassRepo) Create(ctx context.Context, schoolID uuid.UUID, name, schoolYear string) (uuid.UUID, error) {
	const q = `
		INSERT INTO classes (school_id, name, school_year)
		VALUES ($1, $2, $3)
		RETURNING class_id;
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, schoolID, name, schoolYear).Scan(&id)
	return id, err
}

func (r *ClassRepo) List(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]model.Class, int, error) {
	const q = `
		SELECT class_id, school_id, name, school_year,
		       COUNT(*) OVER() AS total_count
		FROM classes
		WHERE school_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	rows, err := r.pool.Query(ctx, q, schoolID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var classes []model.Class
	var total int
	for rows.Next() {
		var c model.Class
		if err := rows.Scan(&c.ClassID, &c.SchoolID, &c.Name, &c.SchoolYear, &total); err != nil {
			return nil, 0, err
		}
		classes = append(classes, c)
	}
	return classes, total, rows.Err()
}

// GetByClassID lấy thông tin lớp học theo class_id
func (r *ClassRepo) GetByClassID(ctx context.Context, classID uuid.UUID) (*model.Class, error) {
	const q = `
		SELECT class_id, school_id, name, school_year
		FROM classes
		WHERE class_id = $1;
	`
	var c model.Class
	err := r.pool.QueryRow(ctx, q, classID).Scan(&c.ClassID, &c.SchoolID, &c.Name, &c.SchoolYear)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CountBySchool đếm tổng số lớp học (nếu schoolID rỗng thì đếm toàn hệ thống)
func (r *ClassRepo) CountBySchool(ctx context.Context, schoolID *uuid.UUID) (int, error) {
	var q string
	var err error
	var count int

	if schoolID != nil {
		q = `SELECT COUNT(*) FROM classes WHERE school_id = $1;`
		err = r.pool.QueryRow(ctx, q, *schoolID).Scan(&count)
	} else {
		q = `SELECT COUNT(*) FROM classes;`
		err = r.pool.QueryRow(ctx, q).Scan(&count)
	}

	return count, err
}

// Update cập nhật thông tin lớp học
func (r *ClassRepo) Update(ctx context.Context, classID uuid.UUID, name, schoolYear string) error {
	const q = `
		UPDATE classes
		SET name = $2, school_year = $3, updated_at = now()
		WHERE class_id = $1;
	`
	ct, err := r.pool.Exec(ctx, q, classID, name, schoolYear)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

// Delete xóa lớp học
func (r *ClassRepo) Delete(ctx context.Context, classID uuid.UUID) error {
	const q = `DELETE FROM classes WHERE class_id = $1;`
	ct, err := r.pool.Exec(ctx, q, classID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}
