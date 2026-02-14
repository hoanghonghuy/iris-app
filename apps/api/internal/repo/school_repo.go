package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SchoolRepo là một struct dùng để quản lý các thao tác với cơ sở dữ liệu liên quan đến trường học (School).
// Trường pool *pgxpool.Pool là một connection pool tới cơ sở dữ liệu PostgreSQL, giúp tái sử dụng kết nối và tối ưu hiệu suất.
// Các phương thức của SchoolRepo sẽ sử dụng pool này để thực hiện các truy vấn như thêm, sửa, xóa, lấy thông tin trường học.
type SchoolRepo struct {
	pool *pgxpool.Pool
}

// NewSchoolRepo là một hàm khởi tạo (constructor) cho struct SchoolRepo.
// Hàm này nhận vào một đối tượng *pgxpool.Pool (connection pool tới PostgreSQL)
// và trả về một con trỏ tới SchoolRepo đã được gán pool này.
// Việc sử dụng connection pool giúp tối ưu hiệu suất khi truy cập cơ sở dữ liệu.
func NewSchoolRepo(pool *pgxpool.Pool) *SchoolRepo {
	return &SchoolRepo{
		pool: pool,
	}
}

func (r *SchoolRepo) Create(ctx context.Context, name, address string) (uuid.UUID, error) {
	const q = `
		INSERT INTO schools (name, address)
		VALUES ($1, $2)
		RETURNING school_id;
	`
	var id uuid.UUID
	// QueryRow được sử dụng để thực hiện một truy vấn SQL mà chỉ trả về một dòng kết quả.
	// ở đây nó thực thi câu lệnh INSERT và trả về giá trị schools_id vừa được tạo.
	// Scan được dùng để đọc giá trị trả về từ truy vấn và gán vào biến id.
	// Nếu truy vấn không trả về dòng nào hoặc có lỗi, Scan sẽ trả về lỗi tương ứng.
	err := r.pool.QueryRow(ctx, q, name, address).Scan(&id)
	return id, err
}

func (r *SchoolRepo) List(ctx context.Context, limit, offset int) ([]model.School, int, error) {
	const q = `
		SELECT school_id, name, COALESCE(address, ''),
		       COUNT(*) OVER() AS total_count
		FROM schools
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`
	rows, err := r.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var schools []model.School
	var total int
	for rows.Next() {
		var s model.School
		if err := rows.Scan(&s.ID, &s.Name, &s.Address, &total); err != nil {
			return nil, 0, err
		}
		schools = append(schools, s)
	}
	return schools, total, rows.Err()
}

// GetByID lấy thông tin trường học theo school_id
func (r *SchoolRepo) GetByID(ctx context.Context, schoolID uuid.UUID) (*model.School, error) {
	const q = `
		SELECT school_id, name, COALESCE(address, '')
		FROM schools
		WHERE school_id = $1;
	`
	var s model.School
	err := r.pool.QueryRow(ctx, q, schoolID).Scan(&s.ID, &s.Name, &s.Address)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
