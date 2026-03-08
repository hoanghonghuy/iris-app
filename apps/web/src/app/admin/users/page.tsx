/**
 * Admin Users Page
 * Quản lý tài khoản người dùng (tạo, khóa, mở khóa, gán role).
 * API: GET/POST /admin/users, POST /admin/users/:user_id/lock|unlock
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { UserCog } from "lucide-react";

export default function AdminUsersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <UserCog className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý Người dùng</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách tài khoản</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Nội dung sẽ được phát triển.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
