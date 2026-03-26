import { UserRole } from "@/types";

export const DASHBOARD_ROUTE_BY_ROLE: Record<UserRole, string> = {
  SUPER_ADMIN: "/admin",
  SCHOOL_ADMIN: "/admin",
  TEACHER: "/teacher",
  PARENT: "/parent",
};

export const PROFILE_ROUTE_BY_ROLE: Record<UserRole, string | null> = {
  SUPER_ADMIN: null,
  SCHOOL_ADMIN: null,
  TEACHER: "/teacher/profile",
  PARENT: "/parent/profile",
};

export const ROLE_LABELS: Record<UserRole, string> = {
  SUPER_ADMIN: "Quản trị viên cấp cao",
  SCHOOL_ADMIN: "Quản trị viên trường",
  TEACHER: "Giáo viên",
  PARENT: "Phụ huynh",
};

export const ADMIN_ALLOWED_ROLES: UserRole[] = ["SUPER_ADMIN", "SCHOOL_ADMIN"];
export const TEACHER_ALLOWED_ROLES: UserRole[] = ["TEACHER"];
export const PARENT_ALLOWED_ROLES: UserRole[] = ["PARENT"];

export function getDashboardRouteByRole(role: UserRole | null | undefined): string {
  if (!role) {
    return "/";
  }
  return DASHBOARD_ROUTE_BY_ROLE[role] || "/";
}
