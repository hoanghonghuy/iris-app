/**
 * Admin Teachers Page
 * Quản lý giáo viên: xem, sửa, gán/gỡ lớp.
 * API: GET/PUT /admin/teachers, POST/DELETE /admin/teachers/:id/classes/:id
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { BookUser } from "lucide-react";

export default function AdminTeachersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <BookUser className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý Giáo viên</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách giáo viên</CardTitle>
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
