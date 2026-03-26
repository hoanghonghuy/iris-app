import { UserRole, UserStatus } from "@/types";

export const USER_ROLE_LABELS: Record<UserRole, string> = {
  SUPER_ADMIN: "Super Admin",
  SCHOOL_ADMIN: "School Admin",
  TEACHER: "Giáo viên",
  PARENT: "Phụ huynh",
};

export const CREATABLE_USER_ROLES: UserRole[] = ["TEACHER", "PARENT", "SCHOOL_ADMIN"];

export const USER_STATUS_VARIANT: Record<UserStatus, "default" | "secondary" | "destructive" | "outline"> = {
  active: "default",
  pending: "secondary",
  locked: "destructive",
};

export const USER_STATUS_LABEL: Record<UserStatus, string> = {
  active: "Hoạt động",
  pending: "Chờ kích hoạt",
  locked: "Đã khóa",
};
