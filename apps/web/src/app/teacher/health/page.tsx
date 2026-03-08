/**
 * Teacher Health Log Page
 * Nhật ký sức khỏe học sinh.
 * API: POST /teacher/health, GET /teacher/students/:student_id/health
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { HeartPulse } from "lucide-react";

export default function TeacherHealthPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <HeartPulse className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Nhật ký sức khỏe</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Sức khỏe học sinh</CardTitle>
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
