package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// quan ly cac ket noi den db
type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// FindByEmail truy vấn thông tin người dùng theo email.
// Trả về con trỏ tới User nếu tìm thấy, hoặc lỗi nếu không tìm thấy hoặc có lỗi truy vấn.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const q = `SELECT user_id, email, password_hash, status FROM users WHERE email=$1 LIMIT 1;`
	u := &model.User{}
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

func (r *UserRepo) FindByID(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	// Get user info
	const userQuery = `SELECT user_id, email, status FROM users WHERE user_id = $1;`

	info := &model.UserInfo{}
	if err := r.pool.QueryRow(ctx, userQuery, userID).Scan(&info.ID, &info.Email, &info.Status); err != nil {
		return nil, err
	}

	// Get roles
	const rolesQuery = `
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON r.role_id = ur.role_id
		WHERE ur.user_id = $1;
	`

	rows, err := r.pool.Query(ctx, rolesQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		info.Roles = append(info.Roles, role)
	}

	return info, rows.Err()
}

// List lấy danh sách tất cả users kèm roles
func (r *UserRepo) List(ctx context.Context) ([]model.UserInfo, error) {
	const q = `
		SELECT u.user_id, u.email, u.status, ARRAY_AGG(r.name ORDER BY r.name) as roles
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		GROUP BY u.user_id, u.email, u.status
		ORDER BY u.created_at DESC;
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserInfo
	for rows.Next() {
		var u model.UserInfo
		if err := rows.Scan(&u.ID, &u.Email, &u.Status, &u.Roles); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// Create tạo mới user trong cơ sở dữ liệu
func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (uuid.UUID, error) {
	const q = `
		INSERT INTO users (email, password_hash, status)
		VALUES ($1, $2, 'active')
		RETURNING user_id;
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, email, passwordHash).Scan(&id)
	return id, err
}

// AssignRole gán role cho user
func (r *UserRepo) AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	const q = `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, role_id FROM roles WHERE name = $2
		ON CONFLICT (user_id, role_id) DO NOTHING;
	`
	_, err := r.pool.Exec(ctx, q, userID, roleName)
	return err
}

// Update cập nhật thông tin user
func (r *UserRepo) Update(ctx context.Context, userID uuid.UUID, email, passwordHash string) error {
	const q = `
		UPDATE users
		SET email = $1, password_hash = $2, updated_at = now()
		WHERE user_id = $3;
	`
	_, err := r.pool.Exec(ctx, q, email, passwordHash, userID)
	return err
}

// Delete xóa user (hard delete)
func (r *UserRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	const q = `DELETE FROM users WHERE user_id = $1;`
	_, err := r.pool.Exec(ctx, q, userID)
	return err
}

// Lock khóa user
func (r *UserRepo) Lock(ctx context.Context, userID uuid.UUID) error {
	const q = `
		UPDATE users
		SET status = 'locked', updated_at = now()
		WHERE user_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, userID)
	return err
}

// Unlock mở khóa user
func (r *UserRepo) Unlock(ctx context.Context, userID uuid.UUID) error {
	const q = `
		UPDATE users
		SET status = 'active', updated_at = now()
		WHERE user_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, userID)
	return err
}
