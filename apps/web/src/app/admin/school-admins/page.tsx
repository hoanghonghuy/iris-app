/**
 * Admin School Admins Page (SUPER_ADMIN only)
 * Quản lý School Admin: tạo, xem, xóa.
 * API: GET/POST/DELETE /admin/school-admins
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ShieldCheck } from "lucide-react";

export default function AdminSchoolAdminsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <ShieldCheck className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý School Admin</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách School Admin</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Chỉ SUPER_ADMIN mới truy cập được trang này. Nội dung sẽ được phát triển.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
