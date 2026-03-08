/**
 * Admin Parents Page
 * Quản lý phụ huynh: xem, gán/gỡ học sinh.
 * API: GET /admin/parents, POST/DELETE /admin/parents/:id/students/:id
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Heart } from "lucide-react";

export default function AdminParentsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Heart className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Quản lý Phụ huynh</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách phụ huynh</CardTitle>
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
