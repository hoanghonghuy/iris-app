/**
 * Admin Students Page
 * Quản lý danh sách học sinh theo lớp.
 * API: GET/POST /admin/students, GET /admin/students/by-class/:class_id
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Users } from "lucide-react";

export default function AdminStudentsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Users className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý Học sinh</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách học sinh</CardTitle>
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
