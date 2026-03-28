/**
 * Teacher Profile Page
 * Cập nhật số điện thoại + đổi mật khẩu.
 * API: PUT /teacher/profile, PUT /me/password
 */
"use client";

import React, { useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Loader2 } from "lucide-react";
import { ChangePasswordForm } from "@/components/shared/ChangePasswordForm";

function extractApiError(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as { response?: { data?: { error?: string } } }).response;
    return response?.data?.error || fallback;
  }
  return fallback;
}

export default function TeacherProfilePage() {
  const { user } = useAuth();
  const [phone, setPhone] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setSubmitting(true); setError(""); setSuccess("");
      await teacherApi.updateMyProfile(phone);
      setSuccess("Cập nhật thành công!");
    } catch (error: unknown) {
      setError(extractApiError(error, "Không thể cập nhật"));
    } finally { setSubmitting(false); }
  };

  return (
    <div className="space-y-6">
      <Card className="max-w-lg">
        <CardHeader><CardTitle className="text-lg">Thông tin tài khoản</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium text-muted-foreground">Email</label>
            <p className="text-sm">{user?.email || "—"}</p>
          </div>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{error}</div>}
            {success && <div className="rounded-md bg-success/10 p-3 text-sm text-success">{success}</div>}
            <div className="space-y-2">
              <label htmlFor="phone" className="text-sm font-medium">Số điện thoại</label>
              <Input id="phone" type="tel" placeholder="0900 000 000" value={phone} onChange={(e) => setPhone(e.target.value)} />
            </div>
            <Button type="submit" disabled={submitting}>
              {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Cập nhật
            </Button>
          </form>
        </CardContent>
      </Card>

      <ChangePasswordForm />
    </div>
  );
}
