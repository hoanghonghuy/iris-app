/**
 * ChangePasswordForm
 * Component dùng chung cho tất cả role để đổi mật khẩu.
 * API: PUT /api/v1/me/password
 */
"use client";

import React, { useState } from "react";
import { authApi } from "@/lib/api/auth.api";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { KeyRound, Loader2 } from "lucide-react";

export function ChangePasswordForm() {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password.length < 6) { setError("Mật khẩu tối thiểu 6 ký tự"); return; }
    if (password !== confirmPassword) { setError("Mật khẩu xác nhận không khớp"); return; }

    try {
      setSubmitting(true); setError(""); setSuccess("");
      await authApi.updateMyPassword(password);
      setSuccess("Đổi mật khẩu thành công!");
      setPassword(""); setConfirmPassword("");
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể đổi mật khẩu");
    } finally { setSubmitting(false); }
  };

  return (
    <Card className="max-w-lg">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-lg">
          <KeyRound className="h-5 w-5" /> Đổi mật khẩu
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{error}</div>}
          {success && <div className="rounded-md bg-success/10 p-3 text-sm text-success">{success}</div>}
          <div className="space-y-2">
            <label htmlFor="newPassword" className="text-sm font-medium">Mật khẩu mới</label>
            <Input id="newPassword" type="password" placeholder="Tối thiểu 6 ký tự" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          <div className="space-y-2">
            <label htmlFor="confirmNewPassword" className="text-sm font-medium">Xác nhận mật khẩu</label>
            <Input id="confirmNewPassword" type="password" placeholder="Nhập lại mật khẩu mới" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required />
          </div>
          <Button type="submit" disabled={submitting}>
            {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Cập nhật mật khẩu
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
