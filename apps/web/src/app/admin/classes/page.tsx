/**
 * Admin Classes Page
 * Quản lý danh sách lớp học theo trường.
 * API: GET/POST /admin/classes, GET /admin/classes/by-school/:school_id
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { GraduationCap } from "lucide-react";

export default function AdminClassesPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <GraduationCap className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý Lớp học</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách lớp học</CardTitle>
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
