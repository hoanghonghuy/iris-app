package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherScopeRepo struct {
	pool *pgxpool.Pool
}

func NewTeacherScopeRepo(pool *pgxpool.Pool) *TeacherScopeRepo {
	return &TeacherScopeRepo{
		pool: pool,
	}
}

// ListMyClass liệt kê các lớp học mà giáo viên (theo user_id) được phân công giảng dạy.
func (r *TeacherScopeRepo) ListMyClasses(ctx context.Context, teacherUserID uuid.UUID) ([]model.Class, error) {
	const q = `
		SELECT c.class_id, c.school_id, c.name, c.school_year
		FROM classes c
		JOIN teacher_classes tc ON tc.class_id = c.class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
		ORDER BY c.school_year DESC, c.name;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []model.Class
	for rows.Next() {
		var c model.Class
		if err := rows.Scan(&c.ClassID, &c.SchoolID, &c.Name, &c.SchoolYear); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, rows.Err()
}

// ListMyStudentsInClass liệt kê học sinh trong một lớp nếu giáo viên được phân công dạy lớp đó
func (r *TeacherScopeRepo) ListMyStudentsInClass(ctx context.Context, teacherUserID, classID uuid.UUID) ([]model.Student, error) {
	const q = `
		SELECT s.student_id, s.school_id, s.current_class_id, s.full_name, s.dob, s.gender
		FROM students s
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND s.current_class_id = $2
		ORDER BY s.full_name;
	`
	rows, err := r.pool.Query(ctx, q, teacherUserID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.StudentID, &s.SchoolID, &s.CurrentClassID, &s.FullName, &s.DOB, &s.Gender); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

// UpsertAttendance: Giáo viên chỉ có thể điểm danh cho học sinh trong lớp của mình.
func (r *TeacherScopeRepo) UpsertAttendance(ctx context.Context, teacherUserID, studentID uuid.UUID,
	date time.Time, status string, checkInAt, checkOutAt *time.Time, note string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const qExisting = `
		SELECT ar.attendance_id, ar.status, COALESCE(ar.note, '')
		FROM attendance_records ar
		JOIN students s ON s.student_id = ar.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
		  AND ar.student_id = $2
		  AND ar.date = $3
		FOR UPDATE;
	`

	var attendanceID uuid.UUID
	var oldStatus string
	var oldNote string
	err = tx.QueryRow(ctx, qExisting, teacherUserID, studentID, date).Scan(&attendanceID, &oldStatus, &oldNote)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		const qInsert = `
			INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
			SELECT s.student_id, $3, $4, $5, $6, $7, $1
			FROM students s
			JOIN teacher_classes tc ON tc.class_id = s.current_class_id
			JOIN teachers t ON t.teacher_id = tc.teacher_id
			WHERE t.user_id = $1 AND s.student_id = $2
			RETURNING attendance_id;
		`

		if err := tx.QueryRow(ctx, qInsert, teacherUserID, studentID, date, status, checkInAt, checkOutAt, note).Scan(&attendanceID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNoRowsUpdated
			}
			return err
		}

		const qInsertLogCreate = `
			INSERT INTO attendance_change_logs (
				attendance_id, student_id, date, change_type,
				new_status, new_note, changed_by
			)
			VALUES ($1, $2, $3, 'create', $4, $5, $6);
		`
		if _, err := tx.Exec(ctx, qInsertLogCreate, attendanceID, studentID, date, status, note, teacherUserID); err != nil {
			return err
		}

		return tx.Commit(ctx)
	}

	const qUpdate = `
		UPDATE attendance_records
		SET status = $2,
			check_in_at = $3,
			check_out_at = $4,
			note = $5,
			recorded_by = $6,
			updated_at = now()
		WHERE attendance_id = $1;
	`

	tag, err := tx.Exec(ctx, qUpdate, attendanceID, status, checkInAt, checkOutAt, note, teacherUserID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsUpdated
	}

	if oldStatus != status || oldNote != note {
		const qInsertLogUpdate = `
			INSERT INTO attendance_change_logs (
				attendance_id, student_id, date, change_type,
				old_status, new_status, old_note, new_note, changed_by
			)
			VALUES ($1, $2, $3, 'update', $4, $5, $6, $7, $8);
		`
		if _, err := tx.Exec(ctx, qInsertLogUpdate, attendanceID, studentID, date, oldStatus, status, oldNote, note, teacherUserID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// DeleteAttendanceForDate huỷ bản ghi điểm danh của một học sinh trong ngày nếu giáo viên có quyền.
func (r *TeacherScopeRepo) DeleteAttendanceForDate(ctx context.Context, teacherUserID, studentID uuid.UUID, date time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const qExisting = `
		SELECT ar.attendance_id, ar.status, COALESCE(ar.note, '')
		FROM attendance_records ar
		JOIN students s ON s.student_id = ar.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
		  AND ar.student_id = $2
		  AND ar.date = $3
		FOR UPDATE;
	`

	var attendanceID uuid.UUID
	var oldStatus string
	var oldNote string
	if err := tx.QueryRow(ctx, qExisting, teacherUserID, studentID, date).Scan(&attendanceID, &oldStatus, &oldNote); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRowsUpdated
		}
		return err
	}

	const qDelete = `
		DELETE FROM attendance_records
		WHERE attendance_id = $1;
	`
	if tag, err := tx.Exec(ctx, qDelete, attendanceID); err != nil {
		return err
	} else if tag.RowsAffected() == 0 {
		return ErrNoRowsUpdated
	}

	const qInsertDeleteLog = `
		INSERT INTO attendance_change_logs (
			attendance_id, student_id, date, change_type,
			old_status, old_note, changed_by
		)
		VALUES ($1, $2, $3, 'delete', $4, $5, $6);
	`
	if _, err := tx.Exec(ctx, qInsertDeleteLog, attendanceID, studentID, date, oldStatus, oldNote, teacherUserID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// CreateClassPost tạo bài đăng cho một lớp học nếu giáo viên được phân công dạy lớp đó.
func (r *TeacherScopeRepo) CreateClassPost(ctx context.Context, teacherUserID, classID uuid.UUID,
	postType, content string) (uuid.UUID, error) {
	const q = `
		INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
		SELECT $1, 'class', $2, $3, $4
		FROM teacher_classes tc
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND tc.class_id = $2
		RETURNING post_id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, teacherUserID, classID, postType, content).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrNoRowsUpdated
		}
		return uuid.Nil, err
	}

	return id, nil
}

// CreateStudentPost tạo bài đăng cho một học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) CreateStudentPost(ctx context.Context, teacherUserID, studentID uuid.UUID,
	postType, content string) (uuid.UUID, error) {
	const q = `
		INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
		SELECT $1, 'student', $2, $3, $4
		FROM students s
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND s.student_id = $2
		RETURNING post_id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, teacherUserID, studentID, postType, content).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrNoRowsUpdated
		}
		return uuid.Nil, err
	}

	return id, nil
}

// UpdatePost cập nhật nội dung bài đăng nếu người dùng hiện tại là tác giả.
// Đồng thời ghi lịch sử trước/sau chỉnh sửa vào post_edit_history.
func (r *TeacherScopeRepo) UpdatePost(ctx context.Context, authorUserID, postID uuid.UUID, newContent string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const qGetCurrent = `
		SELECT content
		FROM posts
		WHERE post_id = $1 AND author_user_id = $2
		FOR UPDATE;
	`

	var oldContent string
	err = tx.QueryRow(ctx, qGetCurrent, postID, authorUserID).Scan(&oldContent)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRowsUpdated
		}
		return err
	}

	const qUpdate = `
		UPDATE posts
		SET content = $3,
			updated_at = now()
		WHERE post_id = $1 AND author_user_id = $2;
	`

	tag, err := tx.Exec(ctx, qUpdate, postID, authorUserID, newContent)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsUpdated
	}

	const qInsertHistory = `
		INSERT INTO post_edit_history (post_id, old_content, new_content, edited_by)
		VALUES ($1, $2, $3, $4);
	`

	if _, err := tx.Exec(ctx, qInsertHistory, postID, oldContent, newContent, authorUserID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeletePost xóa bài đăng nếu người dùng hiện tại là tác giả.
func (r *TeacherScopeRepo) DeletePost(ctx context.Context, authorUserID, postID uuid.UUID) error {
	const q = `
		DELETE FROM posts
		WHERE post_id = $1
		  AND author_user_id = $2;
	`

	tag, err := r.pool.Exec(ctx, q, postID, authorUserID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}

// ListClassPosts liệt kê bài đăng của một lớp nếu giáo viên được phân công dạy lớp đó.
func (r *TeacherScopeRepo) ListClassPosts(ctx context.Context, teacherUserID, classID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	// Truy vấn join bổ sung dữ liệu quan hệ cho kết quả trả về
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content,
			COALESCE(lc.like_count, 0) AS like_count,
			COALESCE(cc.comment_count, 0) AS comment_count,
			COALESCE(sc.share_count, 0) AS share_count,
			EXISTS(
				SELECT 1
				FROM post_interactions self_like
				WHERE self_like.post_id = p.post_id AND self_like.user_id = $1 AND self_like.action_type = 'like'
			) AS liked_by_me,
			p.created_at, p.updated_at,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN teacher_classes tc ON tc.class_id = p.class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS like_count
			FROM post_interactions
			WHERE action_type = 'like'
			GROUP BY post_id
		) lc ON lc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM post_comments
			GROUP BY post_id
		) cc ON cc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS share_count
			FROM post_interactions
			WHERE action_type = 'share'
			GROUP BY post_id
		) sc ON sc.post_id = p.post_id
		WHERE t.user_id = $1 AND p.class_id = $2 AND p.scope_type = 'class'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID, classID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []model.Post
	var total int
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content,
			&p.LikeCount, &p.CommentCount, &p.ShareCount, &p.LikedByMe,
			&p.CreatedAt, &p.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

// ListAttendanceByStudent liệt kê lịch sử điểm danh của một học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) ListAttendanceByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.AttendanceRecord, error) {
	const q = `
		SELECT ar.attendance_id, ar.student_id, ar.date, ar.status,
			ar.check_in_at, ar.check_out_at, COALESCE(ar.note,''), ar.recorded_by
		FROM attendance_records ar
		JOIN students s ON s.student_id = ar.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
			AND ar.student_id = $2
			AND ar.date BETWEEN $3 AND $4
		ORDER BY ar.date DESC;
	`
	rows, err := r.pool.Query(ctx, q, teacherUserID, studentID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.AttendanceRecord
	for rows.Next() {
		var x model.AttendanceRecord
		if err := rows.Scan(
			&x.AttendanceID, &x.StudentID, &x.Date, &x.Status,
			&x.CheckInAt, &x.CheckOutAt, &x.Note, &x.RecordedBy,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// ListAttendanceChangeLogsByStudent liệt kê lịch sử chỉnh sửa điểm danh của một học sinh.
func (r *TeacherScopeRepo) ListAttendanceChangeLogsByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.AttendanceChangeLog, error) {
	const q = `
		SELECT acl.change_id, acl.attendance_id, acl.student_id, acl.date, acl.change_type,
			acl.old_status, acl.new_status, acl.old_note, acl.new_note,
			acl.changed_by, acl.changed_at
		FROM attendance_change_logs acl
		JOIN students s ON s.student_id = acl.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
			AND acl.student_id = $2
			AND acl.date BETWEEN $3 AND $4
		ORDER BY acl.changed_at DESC;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID, studentID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.AttendanceChangeLog
	for rows.Next() {
		var x model.AttendanceChangeLog
		if err := rows.Scan(
			&x.ChangeID, &x.AttendanceID, &x.StudentID, &x.Date, &x.ChangeType,
			&x.OldStatus, &x.NewStatus, &x.OldNote, &x.NewNote,
			&x.ChangedBy, &x.ChangedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}

	return out, rows.Err()
}

// ListAttendanceChangeLogsByClass liệt kê lịch sử chỉnh sửa điểm danh theo lớp có phân trang.
func (r *TeacherScopeRepo) ListAttendanceChangeLogsByClass(ctx context.Context, teacherUserID, classID uuid.UUID,
	studentID *uuid.UUID, status *string, from, to time.Time, limit, offset int) ([]model.AttendanceChangeLog, int, error) {
	const q = `
		SELECT acl.change_id, acl.attendance_id, acl.student_id, s.full_name, acl.date, acl.change_type,
			acl.old_status, acl.new_status, acl.old_note, acl.new_note,
			acl.changed_by, acl.changed_at,
			COUNT(*) OVER() AS total_count
		FROM attendance_change_logs acl
		JOIN students s ON s.student_id = acl.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
			AND s.current_class_id = $2
			AND ($3::uuid IS NULL OR acl.student_id = $3)
			AND ($4::varchar IS NULL OR COALESCE(acl.new_status, acl.old_status) = $4)
			AND acl.date BETWEEN $5 AND $6
		ORDER BY acl.changed_at DESC
		LIMIT $7 OFFSET $8;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID, classID, studentID, status, from, to, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []model.AttendanceChangeLog
	var total int
	for rows.Next() {
		var x model.AttendanceChangeLog
		if err := rows.Scan(
			&x.ChangeID, &x.AttendanceID, &x.StudentID, &x.StudentName, &x.Date, &x.ChangeType,
			&x.OldStatus, &x.NewStatus, &x.OldNote, &x.NewNote,
			&x.ChangedBy, &x.ChangedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, x)
	}

	return out, total, rows.Err()
}

// ListStudentPosts liệt kê bài đăng của một học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) ListStudentPosts(ctx context.Context, teacherUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	// Truy vấn join bổ sung dữ liệu quan hệ cho kết quả trả về
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content,
			COALESCE(lc.like_count, 0) AS like_count,
			COALESCE(cc.comment_count, 0) AS comment_count,
			COALESCE(sc.share_count, 0) AS share_count,
			EXISTS(
				SELECT 1
				FROM post_interactions self_like
				WHERE self_like.post_id = p.post_id AND self_like.user_id = $1 AND self_like.action_type = 'like'
			) AS liked_by_me,
			p.created_at, p.updated_at,
			COUNT(*) OVER() AS total_count
		FROM posts p
		JOIN students s ON s.student_id = p.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS like_count
			FROM post_interactions
			WHERE action_type = 'like'
			GROUP BY post_id
		) lc ON lc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM post_comments
			GROUP BY post_id
		) cc ON cc.post_id = p.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS share_count
			FROM post_interactions
			WHERE action_type = 'share'
			GROUP BY post_id
		) sc ON sc.post_id = p.post_id
		WHERE t.user_id = $1 AND p.student_id = $2 AND p.scope_type = 'student'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []model.Post
	var total int
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content,
			&p.LikeCount, &p.CommentCount, &p.ShareCount, &p.LikedByMe,
			&p.CreatedAt, &p.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

// CountMyStudents đếm tổng số học sinh đang được phân công giảng dạy bởi giáo viên này
func (r *TeacherScopeRepo) CountMyStudents(ctx context.Context, teacherUserID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(DISTINCT s.student_id)
		FROM students s
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1;
	`
	var count int
	err := r.pool.QueryRow(ctx, q, teacherUserID).Scan(&count)
	return count, err
}

// CountMyPosts đếm tổng số bài đăng được tạo bởi giáo viên này
func (r *TeacherScopeRepo) CountMyPosts(ctx context.Context, teacherUserID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(p.post_id)
		FROM posts p
		WHERE p.author_user_id = $1;
	`
	var count int
	err := r.pool.QueryRow(ctx, q, teacherUserID).Scan(&count)
	return count, err
}
