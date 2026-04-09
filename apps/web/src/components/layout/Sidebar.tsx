/**
 * Sidebar Component
 * Thanh điều hướng bên trái, hiển thị menu theo role.
 * Responsive: full (lg) → icon-only (md) → overlay (mobile).
 */
"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { X } from "lucide-react";
import { useAuth } from "@/providers/AuthProvider";
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
  Baby,
  Newspaper,
  LayoutDashboard,
  MessageSquare,
  CalendarClock,
  ListChecks,
} from "lucide-react";

// ─── Menu items config ──────────────────────────────────────────────

export interface SidebarItem {
  label: string;
  href: string;
  icon: React.ReactNode;
  roles?: string[];
}

export const adminMenuItems: SidebarItem[] = [
  { label: "Tổng quan", href: "/admin", icon: <LayoutDashboard className="h-5 w-5" /> },
  { label: "Trường học", href: "/admin/schools", icon: <School className="h-5 w-5" /> },
  { label: "Lớp học", href: "/admin/classes", icon: <GraduationCap className="h-5 w-5" /> },
  { label: "Học sinh", href: "/admin/students", icon: <Users className="h-5 w-5" /> },
  { label: "Người dùng", href: "/admin/users", icon: <UserCog className="h-5 w-5" /> },
  { label: "Giáo viên", href: "/admin/teachers", icon: <BookUser className="h-5 w-5" /> },
  { label: "Phụ huynh", href: "/admin/parents", icon: <Heart className="h-5 w-5" /> },
  { label: "School Admin", href: "/admin/school-admins", icon: <ShieldCheck className="h-5 w-5" />, roles: ["SUPER_ADMIN"] },
  { label: "Audit Logs", href: "/admin/audit-logs", icon: <ListChecks className="h-5 w-5" />, roles: ["SUPER_ADMIN"] },
  { label: "Tin nhắn", href: "/admin/chat", icon: <MessageSquare className="h-5 w-5" /> },
];

export const teacherMenuItems: SidebarItem[] = [
  { label: "Tổng quan", href: "/teacher", icon: <LayoutDashboard className="h-5 w-5" /> },
  { label: "Lớp của tôi", href: "/teacher/classes", icon: <GraduationCap className="h-5 w-5" /> },
  { label: "Điểm danh", href: "/teacher/attendance", icon: <ClipboardCheck className="h-5 w-5" /> },
  { label: "Sức khỏe", href: "/teacher/health", icon: <HeartPulse className="h-5 w-5" /> },
  { label: "Bài đăng", href: "/teacher/posts", icon: <FileText className="h-5 w-5" /> },
  { label: "Lịch hẹn", href: "/teacher/appointments", icon: <CalendarClock className="h-5 w-5" /> },
  { label: "Tin nhắn", href: "/teacher/chat", icon: <MessageSquare className="h-5 w-5" /> },
];

export const parentMenuItems: SidebarItem[] = [
  { label: "Tổng quan", href: "/parent", icon: <LayoutDashboard className="h-5 w-5" /> },
  { label: "Con của tôi", href: "/parent/children", icon: <Baby className="h-5 w-5" /> },
  { label: "Bảng tin", href: "/parent/feed", icon: <Newspaper className="h-5 w-5" /> },
  { label: "Lịch hẹn", href: "/parent/appointments", icon: <CalendarClock className="h-5 w-5" /> },
  { label: "Tin nhắn", href: "/parent/chat", icon: <MessageSquare className="h-5 w-5" /> },
];

// ─── Sidebar Component ──────────────────────────────────────────────

interface SidebarProps {
  items: SidebarItem[];
  isOpen: boolean;
  onClose: () => void;
}

export function Sidebar({ items, isOpen, onClose }: SidebarProps) {
  const pathname = usePathname();
  const { role } = useAuth();

  // Lọc items theo role
  const visibleItems = items.filter(item => !item.roles || (role && item.roles.includes(role)));

  // Kiểm tra active: exact match cho root paths, startsWith cho sub-paths
  const isActive = (href: string) => {
    if (href === "/admin" || href === "/teacher" || href === "/parent") {
      return pathname === href;
    }
    return pathname.startsWith(href);
  };

  return (
    <>
      {/* ── Overlay (mobile only) ── */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* ── Sidebar ── */}
      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-50 flex flex-col border-r border-sidebar-border bg-sidebar text-sidebar-foreground transition-transform duration-200 ease-in-out shadow-lg lg:shadow-none",
          // Mobile: ẩn mặc định, hiện khi isOpen
          isOpen ? "translate-x-0" : "-translate-x-full",
          // Desktop (lg): luôn hiện, width 256px
          "lg:translate-x-0 lg:static lg:w-64",
          // Mobile width khi mở
          "w-64"
        )}
      >
        {/* ── Logo + Close (mobile) ── */}
        <div className="flex h-14 items-center justify-between border-b border-sidebar-border px-4">
          <h2 className="text-lg font-semibold text-sidebar-foreground">🌸 Iris School</h2>
          <button
            className="rounded-md p-1 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground lg:hidden transition-colors"
            onClick={onClose}
          >
            <X className="h-5 w-5 text-sidebar-foreground" />
          </button>
        </div>

        {/* ── Navigation ── */}
        <nav className="flex-1 space-y-1 overflow-y-auto p-3">
          {visibleItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              onClick={onClose}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-colors",
                isActive(item.href)
                  ? "bg-sidebar-primary/10 text-sidebar-primary font-semibold"
                  : "text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground font-medium"
              )}
            >
              {item.icon}
              <span>{item.label}</span>
            </Link>
          ))}
        </nav>
      </aside>
    </>
  );
}
