/**
 * Teacher Classes Page
 * Danh sách lớp học được phân công cho giáo viên.
 * API: GET /teacher/classes
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { GraduationCap } from "lucide-react";

export default function TeacherClassesPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <GraduationCap className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Lớp học của tôi</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách lớp được phân công</CardTitle>
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
