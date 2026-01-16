package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AdminSchoolHandler struct {
	Schools *repo.SchoolRepo
}

type CreateSchoolReq struct {
	Name string `json:"name" binding:"required,min=2"`
	Address string `json:"address"`
}


// Hàm Create là một handler cho endpoint tạo mới trường học (school).
// Khi nhận một HTTP request, hàm này sẽ thực hiện các bước sau:
//
// 	1. Đọc và kiểm tra dữ liệu đầu vào từ body của request (dạng JSON) vào biến req.
//    Nếu dữ liệu không hợp lệ (ví dụ thiếu trường "name" hoặc "name" quá ngắn), trả về lỗi 400.
//
// 	2. Tạo một context với timeout 3 giây để đảm bảo các thao tác với database không bị treo quá lâu.
//
// 	3. Gọi phương thức Create của SchoolRepo để lưu thông tin trường học mới vào database.
//    Nếu có lỗi khi lưu (ví dụ lỗi kết nối database), trả về lỗi 500.
//
// 	4. Nếu tạo thành công, trả về thông tin school_id của trường học vừa được tạo với mã thành công.
func (h *AdminSchoolHandler) Create(c *gin.Context) {
	var req CreateSchoolReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid payload")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	
	id, err := h.Schools.Create(ctx, req.Name, req.Address)
	if err != nil {
		fail(c, http.StatusInternalServerError, "db error")
		return
	}
	created(c, gin.H{
		"school_id": id,
	})
}

func (h *AdminSchoolHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	
	schools, err := h.Schools.List(ctx)
	if err != nil {
		fail(c, http.StatusInternalServerError, "db error")
		return
	}
	ok(c, schools)
}