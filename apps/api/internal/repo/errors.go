package repo

import "errors"

// Sentinel errors kỹ thuật của tầng repo.
// Service layer chịu trách nhiệm map những lỗi này sang business errors tương ứng.

// ErrNoRowsUpdated được trả về khi câu UPDATE không ảnh hưởng đến hàng nào.
// Dùng để báo hiệu điều kiện WHERE không thỏa mãn (không phải lỗi DB),
// ví dụ: atomic check-before-update (parent code exhausted, token already used...).
var ErrNoRowsUpdated = errors.New("no rows updated")
