/**
 * Header Component
 * Thanh header hiển thị hamburger toggle, role, user info, dropdown menu, theme toggle.
 */
"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/providers/AuthProvider";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ClipboardCheck, LogOut, Menu, UserCircle } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";
import { PROFILE_ROUTE_BY_ROLE, ROLE_LABELS } from "@/lib/auth-config";

interface HeaderProps {
  onMenuToggle: () => void;
}

type HeaderMeta = {
  title: string;
  subtitle?: string;
  icon?: React.ReactNode;
};

function resolveHeaderMeta(pathname: string | null): HeaderMeta | null {
  if (!pathname) {
    return null;
  }

  const today = new Date().toISOString().slice(0, 10);

  if (pathname.startsWith("/teacher/attendance")) {
    return {
      title: "Điểm danh",
      subtitle: `Ngày: ${today}`,
      icon: <ClipboardCheck className="h-4 w-4 shrink-0 text-muted-foreground" />,
    };
  }

  if (pathname.startsWith("/teacher/classes/")) {
    return { title: "Chi tiết lớp học" };
  }

  if (pathname.startsWith("/teacher/classes")) {
    return { title: "Lớp của tôi" };
  }

  if (pathname.startsWith("/teacher/health")) {
    return { title: "Sức khỏe học sinh" };
  }

  if (pathname.startsWith("/teacher/posts")) {
    return { title: "Bài đăng" };
  }

  if (pathname.startsWith("/teacher/profile")) {
    return { title: "Hồ sơ cá nhân" };
  }

  if (pathname === "/teacher") {
    return { title: "Tổng quan giáo viên" };
  }

  if (pathname.startsWith("/parent/children/")) {
    return { title: "Thông tin con" };
  }

  if (pathname.startsWith("/parent/children")) {
    return { title: "Con của tôi" };
  }

  if (pathname.startsWith("/parent/posts") || pathname.startsWith("/parent/feed")) {
    return { title: "Bảng tin" };
  }

  if (pathname.startsWith("/parent/profile")) {
    return { title: "Hồ sơ cá nhân" };
  }

  if (pathname === "/parent") {
    return { title: "Tổng quan phụ huynh" };
  }

  if (pathname.startsWith("/admin/schools")) {
    return { title: "Quản lý trường học" };
  }

  if (pathname.startsWith("/admin/school-admins")) {
    return { title: "Quản lý School Admin" };
  }

  if (pathname.startsWith("/admin/classes")) {
    return { title: "Quản lý lớp học" };
  }

  if (pathname.startsWith("/admin/teachers")) {
    return { title: "Quản lý giáo viên" };
  }

  if (pathname.startsWith("/admin/students")) {
    return { title: "Quản lý học sinh" };
  }

  if (pathname.startsWith("/admin/parents")) {
    return { title: "Quản lý phụ huynh" };
  }

  if (pathname.startsWith("/admin/users")) {
    return { title: "Quản lý người dùng" };
  }

  if (pathname === "/admin") {
    return { title: "Tổng quan quản trị" };
  }

  return null;
}

export function Header({ onMenuToggle }: HeaderProps) {
  const { user, role, logout } = useAuth();
  const pathname = usePathname();
  const profileRoute = role ? PROFILE_ROUTE_BY_ROLE[role] : null;
  const headerMeta = resolveHeaderMeta(pathname);

  // Generate initials for avatar
  const initials = user?.email
    ? user.email.substring(0, 2).toUpperCase()
    : "U";

  return (
    <header className="flex h-14 items-center justify-between border-b border-border bg-background px-4 lg:px-6 transition-colors duration-300">
      {/* Left: hamburger (mobile/tablet) + role label */}
      <div className="flex items-center gap-3">
        <button
          className="rounded-md p-1.5 hover:bg-muted lg:hidden transition-colors"
          onClick={onMenuToggle}
        >
          <Menu className="h-5 w-5 text-foreground" />
        </button>
        <span className="text-sm text-muted-foreground hidden sm:inline">
          {role && ROLE_LABELS[role]}
        </span>
      </div>

      {headerMeta && (
        <div className="mx-2 min-w-0 flex-1 items-center justify-center gap-1.5 overflow-hidden px-1 text-center sm:gap-2">
          <div className="flex min-w-0 items-center justify-center gap-1.5 sm:gap-2">
            {headerMeta.icon}
            <p className="truncate text-xs font-semibold text-foreground sm:text-sm">{headerMeta.title}</p>
            {headerMeta.subtitle && (
              <span className="hidden text-xs text-muted-foreground sm:inline">{headerMeta.subtitle}</span>
            )}
          </div>
        </div>
      )}

      {/* Right: Theme Toggle + User Dropdown Menu */}
      <div className="flex items-center gap-2 sm:gap-4">
        <ThemeToggle className="text-muted-foreground" />
        
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="relative h-9 rounded-full pl-2 pr-4 focus-visible:ring-0 hover:bg-muted transition-colors">
              <div className="flex items-center gap-2">
                <Avatar className="h-7 w-7">
                  <AvatarFallback className="bg-primary/10 text-primary text-xs font-medium">
                    {initials}
                  </AvatarFallback>
                </Avatar>
                <span className="text-sm font-medium hidden sm:inline-block max-w-[150px] truncate text-foreground">
                  {user?.email}
                </span>
              </div>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56" align="end" forceMount>
            <DropdownMenuLabel className="font-normal">
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium leading-none truncate">{user?.email}</p>
                <p className="text-xs leading-none text-muted-foreground">
                  {role && ROLE_LABELS[role]}
                </p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {profileRoute && (
              <>
                <DropdownMenuItem asChild>
                  <Link href={profileRoute} className="cursor-pointer w-full flex items-center">
                    <UserCircle className="mr-2 h-4 w-4" />
                    <span>Hồ sơ cá nhân</span>
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
              </>
            )}
            <DropdownMenuItem onClick={logout} className="cursor-pointer text-destructive focus:text-destructive">
              <LogOut className="mr-2 h-4 w-4" />
              <span>Đăng xuất</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
