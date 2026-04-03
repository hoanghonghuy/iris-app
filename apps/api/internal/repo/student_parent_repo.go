package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentParentRepo struct {
	pool *pgxpool.Pool
}

func NewStudentParentRepo(pool *pgxpool.Pool) *StudentParentRepo {
	return &StudentParentRepo{
		pool: pool,
	}
}

func (r *StudentParentRepo) Assign(ctx context.Context, studentID, parentID uuid.UUID, relationship string) error {
	const q = `
		INSERT INTO student_parents (student_id, parent_id, relationship)
		VALUES ($1, $2, $3)
		ON CONFLICT (student_id, parent_id) DO UPDATE
		SET relationship = EXCLUDED.relationship;
	`
	_, err := r.pool.Exec(ctx, q, studentID, parentID, relationship)
	return err
}

func (r *StudentParentRepo) Unassign(ctx context.Context, studentID, parentID uuid.UUID) error {
	const q = `DELETE FROM student_parents WHERE student_id=$1 AND parent_id=$2;`
	_, err := r.pool.Exec(ctx, q, studentID, parentID)
	return err
}
