/**
 * Register Parent Page (Public)
 * Phụ huynh tự đăng ký bằng parent code.
 * API: POST /api/v1/register/parent
 */
"use client";

import React, { useState } from "react";
import { authApi } from "@/lib/api/auth.api";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Heart, Loader2 } from "lucide-react";
import Link from "next/link";

export default function RegisterParentPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [parentCode, setParentCode] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) { setError("Email không được trống"); return; }
    if (password.length < 6) { setError("Mật khẩu tối thiểu 6 ký tự"); return; }
    if (password !== confirmPassword) { setError("Mật khẩu xác nhận không khớp"); return; }
    if (!parentCode.trim()) { setError("Mã phụ huynh không được trống"); return; }

    try {
      setSubmitting(true); setError("");
      await authApi.registerParent({ email, password, parent_code: parentCode });
      setSuccess(true);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể đăng ký");
    } finally { setSubmitting(false); }
  };

  if (success) {
    return (
      <div className="flex w-full items-center justify-center w-full max-w-screen-xl flex justify-center">
        <Card className="w-full max-w-md">
          <CardContent className="flex flex-col items-center py-12">
            <Heart className="h-16 w-16 text-success" />
            <h2 className="mt-4 text-xl font-semibold">Đăng ký thành công!</h2>
            <p className="mt-2 text-sm text-muted-foreground">Bạn có thể đăng nhập ngay bây giờ.</p>
            <Link href="/login">
              <Button className="mt-6">Đăng nhập</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex w-full items-center justify-center w-full max-w-screen-xl flex justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <Heart className="mx-auto h-10 w-10 text-muted-foreground" />
          <CardTitle className="mt-2 text-xl">Đăng ký Phụ huynh</CardTitle>
          <CardDescription>Sử dụng mã phụ huynh từ nhà trường để đăng ký</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{error}</div>}
            <div className="space-y-2">
              <label htmlFor="parentCode" className="text-sm font-medium">Mã phụ huynh <span className="text-destructive">*</span></label>
              <Input id="parentCode" placeholder="Nhập mã từ nhà trường..." value={parentCode} onChange={(e) => setParentCode(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <label htmlFor="email" className="text-sm font-medium">Email <span className="text-destructive">*</span></label>
              <Input id="email" type="email" placeholder="parent@example.com" value={email} onChange={(e) => setEmail(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium">Mật khẩu</label>
              <Input id="password" type="password" placeholder="Tối thiểu 6 ký tự" value={password} onChange={(e) => setPassword(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="text-sm font-medium">Xác nhận mật khẩu</label>
              <Input id="confirmPassword" type="password" placeholder="Nhập lại mật khẩu" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required />
            </div>
            <Button type="submit" className="w-full" disabled={submitting}>
              {submitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Đang xử lý...
                </>
              ) : (
                "Đăng ký"
              )}
            </Button>
            <p className="text-center text-sm text-muted-foreground">
              Đã có tài khoản?{" "}
              <Link href="/login" className="font-medium text-zinc-900 dark:text-zinc-100 hover:underline">Đăng nhập</Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
