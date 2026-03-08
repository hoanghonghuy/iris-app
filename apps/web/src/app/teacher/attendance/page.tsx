/**
 * Teacher Attendance Page
 * Điểm danh học sinh trong lớp.
 * API: POST /teacher/attendance, GET /teacher/students/:student_id/attendance
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ClipboardCheck } from "lucide-react";

export default function TeacherAttendancePage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <ClipboardCheck className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Điểm danh</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Điểm danh học sinh</CardTitle>
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
