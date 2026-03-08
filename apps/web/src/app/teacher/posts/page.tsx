/**
 * Teacher Posts Page
 * Quản lý bài đăng (cho lớp hoặc học sinh).
 * API: POST /teacher/posts, GET /teacher/classes/:id/posts, GET /teacher/students/:id/posts
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { FileText } from "lucide-react";

export default function TeacherPostsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <FileText className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Bài đăng</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Quản lý bài đăng</CardTitle>
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
