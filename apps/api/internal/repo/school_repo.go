package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type School struct {
	ID uuid.UUID
	Name string
	Address string
}

// SchoolRepo là một struct dùng để quản lý các thao tác với cơ sở dữ liệu liên quan đến trường học (School).
// Trường pool *pgxpool.Pool là một connection pool tới cơ sở dữ liệu PostgreSQL, giúp tái sử dụng kết nối và tối ưu hiệu suất.
// Các phương thức của SchoolRepo sẽ sử dụng pool này để thực hiện các truy vấn như thêm, sửa, xóa, lấy thông tin trường học.
type SchoolRepo struct {
	pool *pgxpool.Pool
}

//
// NewSchoolRepo là một hàm khởi tạo (constructor) cho struct SchoolRepo.
// Hàm này nhận vào một đối tượng *pgxpool.Pool (connection pool tới PostgreSQL)
// và trả về một con trỏ tới SchoolRepo đã được gán pool này.
// Việc sử dụng connection pool giúp tối ưu hiệu suất khi truy cập cơ sở dữ liệu.
//
func NewSchoolRepo(pool *pgxpool.Pool) *SchoolRepo {
	return &SchoolRepo{
		pool: pool,
	}
}

func (r *SchoolRepo) Create(ctx context.Context, name, address string) (uuid.UUID, error) {
	const q = `
		INSERT INTO schools (name, address)
		VALUES ($1, $2)
		RETURNING schools_id;
	`
	var id uuid.UUID
	// QueryRow được sử dụng để thực hiện một truy vấn SQL mà chỉ trả về một dòng kết quả.
	// ở đây nó thực thi câu lệnh INSERT và trả về giá trị schools_id vừa được tạo.
	// Scan được dùng để đọc giá trị trả về từ truy vấn và gán vào biến id.
	// Nếu truy vấn không trả về dòng nào hoặc có lỗi, Scan sẽ trả về lỗi tương ứng.
	err := r.pool.QueryRow(ctx, q, name, address).Scan(&id)
	return id, err
}

func (r SchoolRepo) List(ctx context.Context) ([]School, error) {
	const q = `
		SELECT schools_id, name, COALESCE(address, '')
		FROM schools
		ORDER BY created_at DESC;
	`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var schools []School
	for rows.Next() {
		var s School
		if err := rows.Scan(&s.ID, &s.Name, &s.Address); err != nil {
			return nil, err
		}
		schools = append(schools, s)
	}
	return schools, rows.Err()
}