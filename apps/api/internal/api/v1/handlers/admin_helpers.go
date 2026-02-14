package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
)

// extractAdminSchoolID đọc school_id từ context (InjectAdminScope đã set).
// SUPER_ADMIN → nil (không giới hạn trường)
// SCHOOL_ADMIN → *uuid.UUID (chỉ trường mình)
func extractAdminSchoolID(c *gin.Context) *uuid.UUID {
	// Lấy giá trị school_id từ context (middleware InjectAdminScope đã set)
	// - c.Get(middleware.CtxAdminSchoolID) trích xuất giá trị được lưu dưới key CtxAdminSchoolID trong context của Gin.
	// - raw: là giá trị được trả về (kiểu any), có thể là string (với SCHOOL_ADMIN), "" (với SUPER_ADMIN), hoặc nil nếu chưa được set.
	// - exists: là boolean cho biết key CtxAdminSchoolID có tồn tại trong context hay không.
	//   + Nếu exists == false: nghĩa là middleware InjectAdminScope chưa đặt giá trị nào vào context cho key này,
	//     điều này thường xảy ra khi:
	//       • Request không đi qua middleware InjectAdminScope (lỗi cấu hình route),
	//       • Hoặc người dùng không có quyền admin (không được xác thực là admin nên middleware không set gì).
	//     Trong cả hai trường hợp, ta không thể xác định được scope admin → an toàn nhất là xử lý như SUPER_ADMIN (không giới hạn trường),
	//     nên trả về nil.
	//   + Nếu exists == true: tiếp tục xử lý giá trị raw để xác định chính xác loại admin (SUPER_ADMIN hay SCHOOL_ADMIN).
	raw, exists := c.Get(middleware.CtxAdminSchoolID)
	if !exists {
		return nil
	}

	// Thực hiện type assertion để chuyển đổi giá trị raw (kiểu any) sang string.
	// - Nếu raw không phải là string (ok == false), tức là có thể là nil hoặc kiểu khác,
	//   thì không thể trích xuất school_id hợp lệ → trả về nil (xử lý như SUPER_ADMIN).
	// - Nếu raw là string nhưng rỗng (str == ""), điều này tương ứng với trường hợp SUPER_ADMIN
	//   (middleware InjectAdminScope đã set school_id = "" cho SUPER_ADMIN),
	//   nên cũng trả về nil để biểu thị không giới hạn trường.
	str, ok := raw.(string)
	if !ok || str == "" {
		return nil
	}

	// Parse string → uuid.UUID. SCHOOL_ADMIN luôn có school_id hợp lệ
	id, err := uuid.Parse(str)
	if err != nil {
		return nil
	}
	return &id
}
