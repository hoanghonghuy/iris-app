/**
 * Header Component
 * Thanh header hiển thị hamburger toggle, role, user info, dropdown menu.
 */
"use client";

import React from "react";
import Link from "next/link";
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
import { LogOut, Menu, UserCircle } from "lucide-react";

const roleLabels: Record<string, string> = {
  SUPER_ADMIN: "Quản trị viên cấp cao",
  SCHOOL_ADMIN: "Quản trị viên trường",
  TEACHER: "Giáo viên",
  PARENT: "Phụ huynh",
};

// Map role to profile route
const profileRoutes: Record<string, string | null> = {
  SUPER_ADMIN: null, // Admin hasn't got a profile page yet
  SCHOOL_ADMIN: null,
  TEACHER: "/teacher/profile",
  PARENT: "/parent/profile",
};

interface HeaderProps {
  onMenuToggle: () => void;
}

export function Header({ onMenuToggle }: HeaderProps) {
  const { user, role, logout } = useAuth();
  const profileRoute = role ? profileRoutes[role] : null;

  // Generate initials for avatar
  const initials = user?.email
    ? user.email.substring(0, 2).toUpperCase()
    : "U";

  return (
    <header className="flex h-14 items-center justify-between border-b bg-white px-4 lg:px-6">
      {/* Left: hamburger (mobile/tablet) + role label */}
      <div className="flex items-center gap-3">
        <button
          className="rounded-md p-1.5 hover:bg-zinc-100 lg:hidden"
          onClick={onMenuToggle}
        >
          <Menu className="h-5 w-5" />
        </button>
        <span className="text-sm text-muted-foreground hidden sm:inline">
          {role && roleLabels[role]}
        </span>
      </div>

      {/* Right: User Dropdown Menu */}
      <div className="flex items-center gap-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="relative h-9 rounded-full pl-2 pr-4 focus-visible:ring-0">
              <div className="flex items-center gap-2">
                <Avatar className="h-7 w-7">
                  <AvatarFallback className="bg-primary/10 text-primary text-xs font-medium">
                    {initials}
                  </AvatarFallback>
                </Avatar>
                <span className="text-sm font-medium hidden sm:inline-block max-w-[150px] truncate">
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
                  {role && roleLabels[role]}
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
