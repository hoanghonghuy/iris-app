/**
 * chatHelpers.ts
 * Các hàm tiện ích cho tính năng Chat (pure functions, không có side effect).
 */

/**
 * parseJwtPayload giải mã phần payload của JWT token (không xác thực chữ ký).
 * Tránh gọi API thêm round-trip chỉ để lấy user_id từ token đã có sẵn.
 */
export function parseJwtPayload(token: string): { user_id: string; email: string } | null {
  try {
    const base64 = token.split(".")[1];
    const json = atob(base64.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(json);
  } catch {
    return null;
  }
}

/** formatTime hiển thị thời gian ngắn gọn (HH:mm) theo locale tiếng Việt */
export function formatTime(iso: string): string {
  return new Date(iso).toLocaleTimeString("vi-VN", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

/**
 * getInitials trả về 2 ký tự đầu viết hoa từ tên đầy đủ.
 * Ví dụ: "Nguyễn Văn A" → "NA"
 */
export function getInitials(name: string): string {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .substring(0, 2)
    .toUpperCase();
}
