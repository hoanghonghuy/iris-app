/**
 * Reset Password Page
 * Nhập email + mã reset + mật khẩu mới → gọi API → redirect về login.
 * API: POST /api/v1/auth/reset-password
 */
"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { authApi } from "@/lib/api/auth.api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter,
} from "@/components/ui/card";
import { ArrowLeft, CheckCircle2, Loader2 } from "lucide-react";

function extractErrorMessage(err: unknown): string | undefined {
    return (
        typeof (err as { response?: { data?: { error?: string } } }).response?.data?.error === "string"
            ? (err as { response?: { data?: { error?: string } } }).response?.data?.error
            : undefined
    );
}

export default function ResetPasswordPage() {
    const router = useRouter();
    const [email, setEmail] = useState("");
    const [token, setToken] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (!email) {
            setError("Email không được để trống");
            return;
        }
        if (!token) {
            setError("Mã đặt lại mật khẩu không được để trống");
            return;
        }
        if (password.length < 6) {
            setError("Mật khẩu phải có ít nhất 6 ký tự");
            return;
        }
        if (password !== confirmPassword) {
            setError("Mật khẩu xác nhận không khớp");
            return;
        }

        setIsSubmitting(true);
        try {
            await authApi.resetPassword(email, token, password);
            setSuccess(true);
            setTimeout(() => router.push("/login"), 3000);
        } catch (err: unknown) {
            setError(extractErrorMessage(err) || "Không thể đặt lại mật khẩu. Mã có thể đã hết hạn.");
        } finally {
            setIsSubmitting(false);
        }
    };

    if (success) {
        return (
            <div className="flex w-full max-w-screen-xl justify-center">
                <Card className="w-full max-w-md">
                    <CardContent className="flex flex-col items-center gap-4 py-8">
                        <CheckCircle2 className="h-12 w-12 text-success" />
                        <p className="text-center text-sm font-medium">Đặt lại mật khẩu thành công!</p>
                        <p className="text-center text-sm text-muted-foreground">
                            Đang chuyển về trang đăng nhập...
                        </p>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="flex w-full max-w-screen-xl justify-center">
            <Card className="w-full max-w-md">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold text-center">Đặt lại mật khẩu</CardTitle>
                    <CardDescription className="text-center">
                        Nhập email, mã đặt lại mật khẩu và mật khẩu mới
                    </CardDescription>
                </CardHeader>
                <form onSubmit={handleSubmit}>
                    <CardContent className="space-y-4">
                        {error && (
                            <div className="p-3 text-sm text-destructive-foreground bg-destructive rounded-md">
                                {error}
                            </div>
                        )}
                        <div className="space-y-2">
                            <label className="text-sm font-medium" htmlFor="email">Email</label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="name@example.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-medium" htmlFor="token">Mã đặt lại mật khẩu</label>
                            <Input
                                id="token"
                                type="text"
                                placeholder="Nhập mã trong email"
                                value={token}
                                onChange={(e) => setToken(e.target.value)}
                                required
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-medium" htmlFor="password">Mật khẩu mới</label>
                            <Input
                                id="password"
                                type="password"
                                placeholder="Tối thiểu 6 ký tự"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                minLength={6}
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-medium" htmlFor="confirmPassword">Xác nhận mật khẩu</label>
                            <Input
                                id="confirmPassword"
                                type="password"
                                placeholder="Nhập lại mật khẩu"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                required
                                minLength={6}
                            />
                        </div>
                    </CardContent>
                    <CardFooter className="flex flex-col gap-4">
                        <Button className="w-full" type="submit" disabled={isSubmitting}>
                            {isSubmitting ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Đang xử lý...</>
                            ) : (
                                "Đặt lại mật khẩu"
                            )}
                        </Button>
                        <Link href="/login" className="text-sm text-muted-foreground hover:text-primary transition-colors">
                            <ArrowLeft className="inline mr-1 h-3 w-3" />
                            Quay lại đăng nhập
                        </Link>
                    </CardFooter>
                </form>
            </Card>
        </div>
    );
}
