package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Status       string
}

// quan ly cac ket noi den db
type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// FindByEmail truy vấn thông tin người dùng theo email.
// Trả về con trỏ tới User nếu tìm thấy, hoặc lỗi nếu không tìm thấy hoặc có lỗi truy vấn.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	const q = `SELECT user_id, email, password_hash, status FROM users WHERE email=$1 LIMIT 1;`
	u := &User{}
	if err := r.pool.QueryRow(ctx, q, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Status); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) RolesOfUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	const q = `
			SELECT r.name
			FROM roles r
			JOIN user_roles ur ON ur.role_id = r.role_id
			WHERE ur.user_id = $1
			ORDER BY r.name;
			`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		roles = append(roles, s)
	}
	return roles, rows.Err()
}
