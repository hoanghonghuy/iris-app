package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
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
	const q = `
		SELECT user_id, email, password_hash, status, COALESCE(google_sub, '') AS google_sub
		FROM users
		WHERE email = $1
		LIMIT 1;
	`
	u := &model.User{}
	if err := r.pool.QueryRow(ctx, q, email).Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Status, &u.GoogleSub); err != nil {
		return nil, err
	}
	return u, nil
}

// FindByGoogleSub truy vấn người dùng đã liên kết theo Google sub.
func (r *UserRepo) FindByGoogleSub(ctx context.Context, googleSub string) (*model.User, error) {
	const q = `
		SELECT user_id, email, password_hash, status, COALESCE(google_sub, '') AS google_sub
		FROM users
		WHERE google_sub = $1
		LIMIT 1;
	`
	u := &model.User{}
	if err := r.pool.QueryRow(ctx, q, googleSub).Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Status, &u.GoogleSub); err != nil {
		return nil, err
	}
	return u, nil
}

// LinkGoogleSub liên kết tài khoản local với Google sub.
// Chỉ link nếu user chưa liên kết trước đó hoặc đã liên kết đúng sub này.
func (r *UserRepo) LinkGoogleSub(ctx context.Context, userID uuid.UUID, googleSub string) error {
	const q = `
		UPDATE users
		SET google_sub = $2,
		    google_linked_at = now(),
		    updated_at = now()
		WHERE user_id = $1
		  AND (google_sub IS NULL OR google_sub = '' OR google_sub = $2);
	`
	res, err := r.pool.Exec(ctx, q, userID, googleSub)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrGoogleAlreadyLinkedDifferent
	}
	return nil
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
	const userQuery = `
		SELECT u.user_id, u.email, u.status, 
		       COALESCE(t.full_name, p.full_name, sa.full_name, '') as full_name
		FROM users u
		LEFT JOIN teachers t ON u.user_id = t.user_id
		LEFT JOIN parents p ON u.user_id = p.user_id
		LEFT JOIN school_admins sa ON u.user_id = sa.user_id
		WHERE u.user_id = $1;
	`

	info := &model.UserInfo{}
	if err := r.pool.QueryRow(ctx, userQuery, userID).Scan(&info.UserID, &info.Email, &info.Status, &info.FullName); err != nil {
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

// List lấy danh sách users kèm roles (có thể lọc theo school_id và role).
//
// schoolID == nil: tất cả users (SUPER_ADMIN).
//
// schoolID != nil: chỉ users thuộc trường đó qua teachers/parents (SCHOOL_ADMIN).
func (r *UserRepo) List(ctx context.Context, schoolID *uuid.UUID, roleFilter string, limit, offset int) ([]model.UserInfo, int, error) {
	qAll := `
		-- ARRAY_AGG(): gom nhiều dòng thành một array
		-- COUNT(*) OVER(): tính tổng số dòng mà không bị ảnh hưởng bởi LIMIT/OFFSET
		SELECT u.user_id, u.email, u.status, ARRAY_AGG(r.name ORDER BY r.name) as roles,
		       COUNT(*) OVER() as total_count
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		WHERE 1=1
	`
	if roleFilter != "" {
		qAll += ` AND EXISTS (SELECT 1 FROM user_roles ur2 JOIN roles r2 ON ur2.role_id = r2.role_id WHERE ur2.user_id = u.user_id AND r2.name = '` + roleFilter + `')`
	}
	qAll += `
		GROUP BY u.user_id, u.email, u.status
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2;
	`

	qBySchool := `
		SELECT u.user_id, u.email, u.status, ARRAY_AGG(r.name ORDER BY r.name) as roles,
		       COUNT(*) OVER() as total_count
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		WHERE u.user_id IN (
			SELECT user_id FROM teachers WHERE school_id = $3
			UNION
			SELECT user_id FROM parents WHERE school_id = $3
		)
	`
	if roleFilter != "" {
		qBySchool += ` AND EXISTS (SELECT 1 FROM user_roles ur2 JOIN roles r2 ON ur2.role_id = r2.role_id WHERE ur2.user_id = u.user_id AND r2.name = '` + roleFilter + `')`
	}
	qBySchool += `
		GROUP BY u.user_id, u.email, u.status
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2;
	`

	var rows pgx.Rows
	var err error
	if schoolID != nil {
		rows, err = r.pool.Query(ctx, qBySchool, limit, offset, *schoolID)
	} else {
		rows, err = r.pool.Query(ctx, qAll, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.UserInfo
	var total int
	for rows.Next() {
		var u model.UserInfo
		if err := rows.Scan(&u.UserID, &u.Email, &u.Status, &u.Roles, &total); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, rows.Err()
}

// CreateActive tạo mới user với status='active'
func (r *UserRepo) CreateActive(ctx context.Context, email, passwordHash string) (uuid.UUID, error) {
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

// UpdateEmail cập nhật email của user (admin only)
func (r *UserRepo) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {
	const q = `
		UPDATE users
		SET email = $2, updated_at = now()
		WHERE user_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, userID, email)
	return err
}

// UpdatePassword cập nhật password của user (self-service)
func (r *UserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, email, passwordHash string) error {
	const q = `
		UPDATE users
		SET password_hash = $3, updated_at = now()
		WHERE user_id = $1 AND email = $2;
	`
	_, err := r.pool.Exec(ctx, q, userID, email, passwordHash)
	return err
}

// CreatePending tạo user với status='pending'
func (r *UserRepo) CreatePending(ctx context.Context, email, passwordHash string) (uuid.UUID, error) {
	const q = `
		INSERT INTO users (email, password_hash, status)
		VALUES ($1, $2, 'pending')
		RETURNING user_id;
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, email, passwordHash).Scan(&id)
	return id, err
}

// CreateWithRolesTx tạo user + gán role trong một transaction để tránh partial state.
func (r *UserRepo) CreateWithRolesTx(ctx context.Context, email, passwordHash, status string, roles []string) (uuid.UUID, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const createUserQuery = `
		INSERT INTO users (email, password_hash, status)
		VALUES ($1, $2, $3)
		RETURNING user_id;
	`

	var userID uuid.UUID
	if err := tx.QueryRow(ctx, createUserQuery, email, passwordHash, status).Scan(&userID); err != nil {
		return uuid.Nil, err
	}

	const assignRoleQuery = `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, role_id FROM roles WHERE name = $2
		ON CONFLICT (user_id, role_id) DO NOTHING;
	`

	for _, role := range roles {
		result, err := tx.Exec(ctx, assignRoleQuery, userID, role)
		if err != nil {
			return uuid.Nil, fmt.Errorf("%w: %s", ErrRoleAssignmentFailed, role)
		}

		if result.RowsAffected() == 0 {
			return uuid.Nil, fmt.Errorf("%w: %s", ErrRoleAssignmentFailed, role)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// SetActivationToken lưu activation token
func (r *UserRepo) SetActivationToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	const q = `
		UPDATE users
		SET activation_token = $2, token_expires_at = $3, updated_at = now()
		WHERE user_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, userID, token, expiresAt)
	return err
}

// IsUserInSchool kiểm tra user có thuộc trường hay không (qua teachers/parents)
func (r *UserRepo) IsUserInSchool(ctx context.Context, userID, schoolID uuid.UUID) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1 FROM teachers WHERE user_id = $1 AND school_id = $2
			-- UNION ALL: gộp (nối) kết quả
			UNION ALL
			SELECT 1 FROM parents WHERE user_id = $1 AND school_id = $2
		);
	`

	// const q = `
	// 	SELECT
	// 		EXISTS (SELECT 1 FROM teachers WHERE user_id = $1 AND school_id = $2)
	// 		OR
	// 		EXISTS (SELECT 1 FROM parents  WHERE user_id = $1 AND school_id = $2);
	// `

	var exists bool
	err := r.pool.QueryRow(ctx, q, userID, schoolID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// FindByActivationToken tìm user theo activation token
func (r *UserRepo) FindByActivationToken(ctx context.Context, token string) (*model.UserWithToken, error) {
	const q = `
		SELECT user_id, email, password_hash, status, activation_token, token_expires_at
		FROM users
		WHERE activation_token = $1;
	`
	u := &model.UserWithToken{}
	err := r.pool.QueryRow(ctx, q, token).Scan(
		&u.UserID,
		&u.Email,
		&u.PasswordHash,
		&u.Status,
		&u.ActivationToken,
		&u.TokenExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ActivateWithPassword activate user + set password
func (r *UserRepo) ActivateWithPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	const q = `
		UPDATE users
		SET password_hash = $2,
			status = 'active',
			activation_token = NULL,
			token_expires_at = NULL,
			updated_at = now()
		WHERE user_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, userID, passwordHash)
	return err
}

// CountUsersByRoleAndSchool đếm tổng số user theo role và school
// Nếu schoolID rỗng, đếm tất cả user có role đó trên toàn hệ thống
func (r *UserRepo) CountUsersByRoleAndSchool(ctx context.Context, role string, schoolID *uuid.UUID) (int, error) {
	var q string
	var err error
	var count int

	if schoolID != nil {
		q = `
			SELECT COUNT(DISTINCT u.user_id)
			FROM users u
			JOIN user_roles ur ON ur.user_id = u.user_id
			JOIN roles r ON r.role_id = ur.role_id
			WHERE r.name = $1
			AND u.user_id IN (
				SELECT user_id FROM teachers WHERE school_id = $2
				UNION
				SELECT user_id FROM parents WHERE school_id = $2
			);
		`
		err = r.pool.QueryRow(ctx, q, role, *schoolID).Scan(&count)
	} else {
		q = `
			SELECT COUNT(u.user_id)
			FROM users u
			JOIN user_roles ur ON ur.user_id = u.user_id
			JOIN roles r ON r.role_id = ur.role_id
			WHERE r.name = $1;
		`
		err = r.pool.QueryRow(ctx, q, role).Scan(&count)
	}

	return count, err
}
