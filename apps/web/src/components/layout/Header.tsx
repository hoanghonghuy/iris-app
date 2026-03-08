/**
 * Header Component
 * Thanh header hiển thị thông tin user và nút đăng xuất.
 * Sẽ được sử dụng trong layout.tsx của từng role.
 */
"use client";

import React from "react";
import { useAuth } from "@/providers/AuthProvider";
import { Button } from "@/components/ui/button";
import { LogOut } from "lucide-react";

const roleLabels: Record<string, string> = {
  SUPER_ADMIN: "Quản trị viên cấp cao",
  SCHOOL_ADMIN: "Quản trị viên trường",
  TEACHER: "Giáo viên",
  PARENT: "Phụ huynh",
};

export function Header() {
  const { user, role, logout } = useAuth();

  return (
    <header className="flex h-14 items-center justify-between border-b bg-white px-6">
      <div className="text-sm text-muted-foreground">
        {role && roleLabels[role]}
      </div>
      <div className="flex items-center gap-4">
        <span className="text-sm font-medium">{user?.email}</span>
        <Button variant="ghost" size="sm" onClick={logout}>
          <LogOut className="mr-2 h-4 w-4" />
          Đăng xuất
        </Button>
      </div>
    </header>
  );
}
