/**
 * Parent Feed Page
 * Bảng tin tổng hợp của tất cả con.
 * API: GET /parent/feed
 *
 * TODO: add server-side pagination when implementing feed content
 */
"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Newspaper } from "lucide-react";

export default function ParentFeedPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Newspaper className="h-8 w-8" />
        <h1 className="text-3xl font-bold tracking-tight">Bảng tin</h1>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Bảng tin tổng hợp</CardTitle>
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
