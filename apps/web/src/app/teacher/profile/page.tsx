/**
 * Teacher Profile Page
 * Giáo viên cập nhật hồ sơ cá nhân (chỉ phone).
 * API: PUT /teacher/profile
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { UserPen } from "lucide-react";

export default function TeacherProfilePage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <UserPen className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Hồ sơ cá nhân</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Thông tin giáo viên</CardTitle>
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
