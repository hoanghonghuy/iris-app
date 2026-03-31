/**
 * Activate Account Page (Public)
 * User nhận activation token qua email, nhập token + mật khẩu mới để kích hoạt.
 * API: POST /api/v1/users/activate-token
 */
"use client";

import React, { useState } from "react";
import { authApi } from "@/lib/api/auth.api";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ShieldCheck, Loader2 } from "lucide-react";
import Link from "next/link";
import { extractApiErrorRawMessage } from "@/lib/api-error";
import { AUTH_PAGE_CARD_CLASS, AUTH_PAGE_CONTAINER_CLASS } from "@/components/auth/auth-layout";

export default function ActivateAccountPage() {
  const [token, setToken] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token.trim()) { setError("Token không được trống"); return; }
    if (password.length < 6) { setError("Mật khẩu tối thiểu 6 ký tự"); return; }
    if (password !== confirmPassword) { setError("Mật khẩu xác nhận không khớp"); return; }

    try {
      setSubmitting(true); setError("");
      await authApi.activateWithToken({ token: token.trim(), password });
      setSuccess(true);
    } catch (err: unknown) {
      setError(extractApiErrorRawMessage(err) || "Không thể kích hoạt tài khoản");
    } finally { setSubmitting(false); }
  };

  if (success) {
    return (
      <div className={AUTH_PAGE_CONTAINER_CLASS}>
        <Card className={AUTH_PAGE_CARD_CLASS}>
          <CardContent className="flex flex-col items-center py-12">
            <ShieldCheck className="h-16 w-16 text-success" />
            <h2 className="mt-4 text-xl font-semibold">Kích hoạt thành công!</h2>
            <p className="mt-2 text-sm text-muted-foreground">Bạn có thể đăng nhập với mật khẩu mới.</p>
            <Link href="/login">
              <Button className="mt-6">Đăng nhập</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className={AUTH_PAGE_CONTAINER_CLASS}>
      <Card className={AUTH_PAGE_CARD_CLASS}>
        <CardHeader className="text-center">
          <ShieldCheck className="mx-auto h-10 w-10 text-muted-foreground" />
          <CardTitle className="mt-2 text-xl">Kích hoạt tài khoản</CardTitle>
          <CardDescription>Nhập token kích hoạt và đặt mật khẩu mới</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{error}</div>}
            <div className="space-y-2">
              <label htmlFor="token" className="text-sm font-medium">Token kích hoạt</label>
              <Input id="token" placeholder="Nhập token từ email..." value={token} onChange={(e) => setToken(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium">Mật khẩu mới</label>
              <Input id="password" type="password" placeholder="Tối thiểu 6 ký tự" value={password} onChange={(e) => setPassword(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="text-sm font-medium">Xác nhận mật khẩu</label>
              <Input id="confirmPassword" type="password" placeholder="Nhập lại mật khẩu" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required />
            </div>
            <Button type="submit" className="w-full" disabled={submitting}>
              {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Kích hoạt
            </Button>
            <p className="text-center text-sm text-muted-foreground">
              Đã có tài khoản?{" "}
              <Link href="/login" className="font-medium text-foreground hover:underline">Đăng nhập</Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
