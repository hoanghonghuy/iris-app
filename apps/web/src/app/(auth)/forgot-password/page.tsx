/**
 * Forgot Password Page
 * Nhập email → gửi yêu cầu đặt lại mật khẩu.
 * API: POST /api/v1/auth/forgot-password
 */
"use client";

import React, { useState } from "react";
import Link from "next/link";
import { authApi } from "@/lib/api/auth.api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter,
} from "@/components/ui/card";
import { ArrowLeft, Mail, CheckCircle2, Loader2 } from "lucide-react";

function extractErrorMessage(err: unknown): string | undefined {
    return (
        typeof (err as { response?: { data?: { error?: string } } }).response?.data?.error === "string"
            ? (err as { response?: { data?: { error?: string } } }).response?.data?.error
            : undefined
    );
}

export default function ForgotPasswordPage() {
    const [email, setEmail] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [sent, setSent] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setIsSubmitting(true);

        try {
            await authApi.forgotPassword(email);
            setSent(true);
        } catch (err: unknown) {
            setError(extractErrorMessage(err) || "Không thể gửi yêu cầu. Vui lòng thử lại.");
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="flex w-full items-center justify-center w-full max-w-screen-xl flex justify-center">
            <Card className="w-full max-w-md">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold text-center">Quên mật khẩu</CardTitle>
                    <CardDescription className="text-center">
                        Nhập email để nhận link đặt lại mật khẩu
                    </CardDescription>
                </CardHeader>

                {sent ? (
                    <>
                        <CardContent className="flex flex-col items-center gap-4 py-6">
                            <CheckCircle2 className="h-12 w-12 text-success" />
                            <p className="text-center text-sm text-muted-foreground">
                                Nếu email <span className="font-medium text-foreground">{email}</span> tồn tại trong hệ thống,
                                bạn sẽ nhận được link đặt lại mật khẩu trong vài phút.
                            </p>
                        </CardContent>
                        <CardFooter>
                            <Link href="/login" className="w-full">
                                <Button variant="outline" className="w-full">
                                    <ArrowLeft className="mr-2 h-4 w-4" />
                                    Quay lại đăng nhập
                                </Button>
                            </Link>
                        </CardFooter>
                    </>
                ) : (
                    <form onSubmit={handleSubmit}>
                        <CardContent className="space-y-4">
                            {error && (
                                <div className="p-3 text-sm text-destructive-foreground bg-destructive rounded-md">
                                    {error}
                                </div>
                            )}
                            <div className="space-y-2">
                                <label className="text-sm font-medium" htmlFor="email">Email</label>
                                <div className="relative">
                                    <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="email"
                                        type="email"
                                        className="pl-10"
                                        placeholder="name@example.com"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        required
                                    />
                                </div>
                            </div>
                        </CardContent>
                        <CardFooter className="flex flex-col gap-4">
                            <Button className="w-full" type="submit" disabled={isSubmitting}>
                                {isSubmitting ? (
                                    <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Đang gửi...</>
                                ) : (
                                    "Gửi link đặt lại mật khẩu"
                                )}
                            </Button>
                            <Link href="/login" className="text-sm text-muted-foreground hover:text-primary transition-colors">
                                <ArrowLeft className="inline mr-1 h-3 w-3" />
                                Quay lại đăng nhập
                            </Link>
                        </CardFooter>
                    </form>
                )}
            </Card>
        </div>
    );
}
