/**
 * Parent Registration Page
 * Phụ huynh tự đăng ký bằng parent code.
 * API: POST /register/parent
 */
"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";

export default function RegisterParentPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [parentCode, setParentCode] = useState("");

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">
            Đăng ký Phụ huynh
          </CardTitle>
          <CardDescription className="text-center">
            Sử dụng mã code được cung cấp bởi nhà trường để đăng ký
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor="parentCode">
              Mã phụ huynh
            </label>
            <Input
              id="parentCode"
              placeholder="VD: ABC12345"
              value={parentCode}
              onChange={(e) => setParentCode(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor="email">
              Email
            </label>
            <Input
              id="email"
              type="email"
              placeholder="name@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor="password">
              Mật khẩu
            </label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Button className="w-full">Đăng ký</Button>
          <p className="text-sm text-muted-foreground text-center">
            Đã có tài khoản?{" "}
            <a href="/login" className="underline">
              Đăng nhập
            </a>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
