/**
 * Parent Child Detail Page
 * Xem thông tin và bài đăng liên quan đến một đứa con cụ thể.
 * API: GET /parent/children/:student_id/posts
 *      GET /parent/children/:student_id/class-posts
 *      GET /parent/children/:student_id/student-posts
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { User } from "lucide-react";

export default function ParentChildDetailPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <User className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Thông tin con</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Bài đăng liên quan</CardTitle>
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
