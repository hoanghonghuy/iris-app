/**
 * Admin Schools Page
 * Quản lý danh sách trường học.
 * API: GET/POST /admin/schools
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { School } from "lucide-react";

export default function AdminSchoolsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <School className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý Trường học</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách trường học</CardTitle>
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
