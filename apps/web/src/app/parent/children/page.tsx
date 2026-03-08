/**
 * Parent Children Page
 * Danh sách con của phụ huynh.
 * API: GET /parent/children
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Baby } from "lucide-react";

export default function ParentChildrenPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Baby className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Con của tôi</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Danh sách con</CardTitle>
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
