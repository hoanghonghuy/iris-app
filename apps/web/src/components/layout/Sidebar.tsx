/**
 * Sidebar Component
 * Thanh điều hướng bên trái, hiển thị menu theo role.
 * Sẽ được sử dụng trong layout.tsx của từng role (admin, teacher, parent).
 */
"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import {
  School,
  GraduationCap,
  Users,
  UserCog,
  BookUser,
  Heart,
  ShieldCheck,
  ClipboardCheck,
  HeartPulse,
  FileText,
  UserPen,
  Baby,
  Newspaper,
  LayoutDashboard,
} from "lucide-react";

export interface SidebarItem {
  label: string;
  href: string;
  icon: React.ReactNode;
}

// Menu items theo từng role
export const adminMenuItems: SidebarItem[] = [
  { label: "Tổng quan", href: "/admin", icon: <LayoutDashboard className="h-4 w-4" /> },
  { label: "Trường học", href: "/admin/schools", icon: <School className="h-4 w-4" /> },
  { label: "Lớp học", href: "/admin/classes", icon: <GraduationCap className="h-4 w-4" /> },
  { label: "Học sinh", href: "/admin/students", icon: <Users className="h-4 w-4" /> },
  { label: "Người dùng", href: "/admin/users", icon: <UserCog className="h-4 w-4" /> },
  { label: "Giáo viên", href: "/admin/teachers", icon: <BookUser className="h-4 w-4" /> },
  { label: "Phụ huynh", href: "/admin/parents", icon: <Heart className="h-4 w-4" /> },
  { label: "School Admin", href: "/admin/school-admins", icon: <ShieldCheck className="h-4 w-4" /> },
];

export const teacherMenuItems: SidebarItem[] = [
  { label: "Tổng quan", href: "/teacher", icon: <LayoutDashboard className="h-4 w-4" /> },
  { label: "Lớp của tôi", href: "/teacher/classes", icon: <GraduationCap className="h-4 w-4" /> },
  { label: "Điểm danh", href: "/teacher/attendance", icon: <ClipboardCheck className="h-4 w-4" /> },
  { label: "Sức khỏe", href: "/teacher/health", icon: <HeartPulse className="h-4 w-4" /> },
  { label: "Bài đăng", href: "/teacher/posts", icon: <FileText className="h-4 w-4" /> },
  { label: "Hồ sơ", href: "/teacher/profile", icon: <UserPen className="h-4 w-4" /> },
];

export const parentMenuItems: SidebarItem[] = [
  { label: "Tổng quan", href: "/parent", icon: <LayoutDashboard className="h-4 w-4" /> },
  { label: "Con của tôi", href: "/parent/children", icon: <Baby className="h-4 w-4" /> },
  { label: "Bảng tin", href: "/parent/feed", icon: <Newspaper className="h-4 w-4" /> },
];

interface SidebarProps {
  items: SidebarItem[];
  title?: string;
}

export function Sidebar({ items, title = "Iris School" }: SidebarProps) {
  const pathname = usePathname();

  return (
    <aside className="flex h-screen w-64 flex-col border-r bg-white">
      <div className="flex h-14 items-center border-b px-4">
        <h2 className="text-lg font-semibold">{title}</h2>
      </div>
      <nav className="flex-1 space-y-1 p-2">
        {items.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
              pathname === item.href
                ? "bg-zinc-100 text-zinc-900"
                : "text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900"
            )}
          >
            {item.icon}
            {item.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
