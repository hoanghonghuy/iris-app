/**
 * Parent Profile Page
 * Xem thông tin tài khoản + đổi mật khẩu.
 * API: PUT /me/password
 */
"use client";

import React from "react";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { User } from "lucide-react";
import { ChangePasswordForm } from "@/components/shared/ChangePasswordForm";

export default function ParentProfilePage() {
  const { user } = useAuth();

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <User className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Hồ sơ cá nhân</h1>
      </div>

      <Card className="max-w-lg">
        <CardHeader><CardTitle className="text-lg">Thông tin tài khoản</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium text-muted-foreground">Email</label>
            <p className="text-sm">{user?.email || "—"}</p>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-muted-foreground">Vai trò</label>
            <p className="text-sm">Phụ huynh</p>
          </div>
        </CardContent>
      </Card>

      <ChangePasswordForm />
    </div>
  );
}
