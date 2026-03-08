/**
 * Admin School-Admins Page
 * Danh sách quản trị viên trường (SCHOOL_ADMIN).
 * API: GET /admin/school-admins (nếu có), hoặc dùng /admin/users filter
 */
"use client";

import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { ShieldCheck } from "lucide-react";

export default function AdminSchoolAdminsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <ShieldCheck className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý School Admin</h1>
      </div>

      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <ShieldCheck className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">
            Chức năng quản lý School Admin đang được phát triển
          </p>
          <p className="mt-1 text-xs text-muted-foreground">
            Hiện tại có thể quản lý qua trang Người dùng
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
