export function extractErrorMessage(error) {
  // Lỗi mạng hoặc lỗi hệ thống không có response data
  if (!error.data) {
    return error.message || 'Có lỗi xảy ra, vui lòng thử lại sau.'
  }

  // API trả lỗi dạng { error: "message" }
  if (error.data.error) {
    return error.data.error
  }

  // Fallback
  return 'Lỗi máy chủ không xác định.'
}
