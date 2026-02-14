package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherClassRepo struct {
	pool *pgxpool.Pool
}

func NewTeacherClassRepo(pool *pgxpool.Pool) *TeacherClassRepo {
	return &TeacherClassRepo{
		pool: pool,
	}
}

// Assign: gán một giáo viên cho một lớp
//   - Thêm mới mối quan hệ teacher↔class nếu chưa tồn tại
//   - Không làm gì (không lỗi) nếu mối quan hệ đã có
//
// Input:
//   - ctx: Context cho timeout/cancellation
//   - teacherID: UUID của giáo viên cần gán
//   - classID: UUID của lớp cần gán
func (r *TeacherClassRepo) Assign(ctx context.Context, teacherID, classID uuid.UUID) error {
	const q = `
		INSERT INTO teacher_classes (teacher_id, class_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`
	_, err := r.pool.Exec(ctx, q, teacherID, classID)
	return err
}

// ListTeachersOfClass: Lấy danh sách giáo viên đang dạy một lớp nào đó
//
// Input:
//   - ctx: Context cho timeout/cancellation
//   - classID: UUID của lớp cần tìm giáo viên
//
// Output:
//   - []uuid.UUID: Danh sách giáo viên giảng dạy lớp
//   - error: Nếu query có lỗi
func (r *TeacherClassRepo) ListTeachersOfClass(ctx context.Context, classID uuid.UUID) ([]uuid.UUID, error) {
	const q = `
		SELECT teacher_id
		FROM teacher_classes
		WHERE class_id = $1
		ORDER BY teacher_id;
	`
	rows, err := r.pool.Query(ctx, q, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

// Unassign: Hủy gán giáo viên khỏi lớp (xóa relationship)
//
// Input:
//   - ctx: Context cho timeout/cancellation
//   - teacherID: UUID của giáo viên cần xóa khỏi lớp
//   - classID: UUID của lớp cần xóa khỏi giáo viên
func (r *TeacherClassRepo) Unassign(ctx context.Context, teacherID, classID uuid.UUID) error {
	const q = `DELETE FROM teacher_classes WHERE teacher_id=$1 AND class_id=$2;`
	_, err := r.pool.Exec(ctx, q, teacherID, classID)
	return err
}

// kiểm tra xem một giáo viên có được gán cho một lớp hay không
func (r *TeacherClassRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID uuid.UUID) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM teacher_classes WHERE teacher_id=$1 AND class_id=$2);`
	var exists bool
	err := r.pool.QueryRow(ctx, q, teacherID, classID).Scan(&exists)
	return exists, err
}

// ListTeacherDetailsOfClass: Lấy danh sách giáo viên (kèm thông tin chi tiết) đang dạy một lớp.
// Thay thế cho pattern N+1: ListTeachersOfClass (lấy IDs) + loop GetByTeacherID
//
// Input:
//   - ctx: Context cho timeout/cancellation
//   - classID: UUID của lớp cần tìm giáo viên
//
// Output:
//   - []model.Teacher: Danh sách giáo viên với đầy đủ thông tin
//   - error: Nếu query có lỗi
func (r *TeacherClassRepo) ListTeacherDetailsOfClass(ctx context.Context, classID uuid.UUID) ([]model.Teacher, error) {
	const q = `
		SELECT t.teacher_id, t.user_id, u.email, t.full_name, COALESCE(t.phone,''), t.school_id
		FROM teacher_classes tc
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		JOIN users u ON u.user_id = t.user_id
		WHERE tc.class_id = $1
		ORDER BY t.full_name;
	`
	rows, err := r.pool.Query(ctx, q, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []model.Teacher
	for rows.Next() {
		var t model.Teacher
		if err := rows.Scan(&t.TeacherID, &t.UserID, &t.Email, &t.FullName, &t.Phone, &t.SchoolID); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}
