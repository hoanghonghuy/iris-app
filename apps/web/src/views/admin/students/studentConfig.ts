export const STUDENT_INITIAL_FORM = {
  id: '',
  full_name: '',
  dob: '',
  gender: 'male',
}

export const STUDENT_COPY_FEEDBACK_TIMEOUT_MS = 2000

export const STUDENT_ERROR_MESSAGES = {
  LOAD_STUDENTS_LIST: 'Không thể tải danh sách học sinh',
  LOAD_STUDENTS_DATA: 'Không thể tải dữ liệu học sinh',
  SAVE_STUDENT: 'Không thể lưu học sinh',
  DELETE_PREFIX: 'Lỗi xóa',
  GENERATE_PARENT_CODE: 'Không thể tạo mã phụ huynh',
  REVOKE_PARENT_CODE: 'Không thể thu hồi mã phụ huynh',
  COPY_PARENT_CODE: 'Không thể sao chép mã phụ huynh',
  REQUIRED_FULL_NAME: 'Họ tên không được để trống',
  REQUIRED_DOB: 'Ngày sinh không được để trống',
}

export const STUDENT_REVOKE_CONFIRM_MESSAGE =
  'Mã phụ huynh hiện tại sẽ bị vô hiệu hóa. Phụ huynh đang sử dụng mã này sẽ bị đăng xuất. Bạn có chắc chắn muốn tiếp tục?'

export function getDeleteStudentConfirmMessage(studentName) {
  return `Bạn có chắc muốn xóa học sinh '${studentName || ''}' không? Hành động này không thể hoàn tác.`
}
