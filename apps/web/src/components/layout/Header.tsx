/**
 * Header Component
 * Thanh header hiển thị hamburger toggle, role, user info, đăng xuất.
 */
"use client";

import React from "react";
import { useAuth } from "@/providers/AuthProvider";
import { Button } from "@/components/ui/button";
import { LogOut, Menu } from "lucide-react";

const roleLabels: Record<string, string> = {
  SUPER_ADMIN: "Quản trị viên cấp cao",
  SCHOOL_ADMIN: "Quản trị viên trường",
  TEACHER: "Giáo viên",
  PARENT: "Phụ huynh",
};

interface HeaderProps {
  onMenuToggle: () => void;
}

export function Header({ onMenuToggle }: HeaderProps) {
  const { user, role, logout } = useAuth();

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

      {/* Right: email + logout */}
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium hidden sm:inline">{user?.email}</span>
        <Button variant="ghost" size="sm" onClick={logout}>
          <LogOut className="h-4 w-4" />
          <span className="hidden sm:inline ml-1">Đăng xuất</span>
        </Button>
      </div>
    </header>
  );
}
