package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
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

func (r ClassRepo) Create(ctx context.Context, schoolID uuid.UUID, name, schoolYear string) (uuid.UUID, error) {
	const q = `
		INSERT INTO classes (school_id, name, school_year)
		VALUES ($1, $2, $3)
		RETURNING class_id;
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, schoolID, name, schoolYear).Scan(&id)
	return id, err
}

func (r ClassRepo) List(ctx context.Context, schoolID uuid.UUID) ([]model.Class, error) {
	const q = `
		SELECT class_id, school_id, name, school_year
		FROM classes
		WHERE school_id = $1
		ORDER BY created_at DESC;
	`
	rows, err := r.pool.Query(ctx, q, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var classes []model.Class
	for rows.Next() {
		var c model.Class
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &c.SchoolYear); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, rows.Err()
} 