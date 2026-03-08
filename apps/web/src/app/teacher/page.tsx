/**
 * Teacher Dashboard
 * Trang tổng quan cho giáo viên: xem danh sách lớp được phân công.
 */
"use client";

import React, { useEffect, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { GraduationCap, Users, ClipboardCheck, Heart, Loader2 } from "lucide-react";
import Link from "next/link";

export default function TeacherDashboard() {
  const { user } = useAuth();
  const [classes, setClasses] = useState<Class[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const data = await teacherApi.getMyClasses();
        setClasses(data || []);
      } catch {
        // Ignore
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold tracking-tight">Bảng điều khiển Giáo viên</h1>

      <Card>
        <CardContent className="py-4">
          <p className="font-medium">Xin chào, {user?.email}</p>
          <p className="mt-1 text-sm text-muted-foreground">
            Bạn đang được phân công dạy {loading ? "..." : classes.length} lớp
          </p>
        </CardContent>
      </Card>

      {/* Quick links */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Link href="/teacher/classes">
          <Card className="transition-colors hover:bg-zinc-50">
            <CardHeader className="pb-2">
              <GraduationCap className="h-6 w-6 text-muted-foreground" />
              <CardTitle className="text-lg">Lớp của tôi</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Xem danh sách lớp và học sinh</CardDescription>
            </CardContent>
          </Card>
        </Link>

        <Link href="/teacher/attendance">
          <Card className="transition-colors hover:bg-zinc-50">
            <CardHeader className="pb-2">
              <ClipboardCheck className="h-6 w-6 text-muted-foreground" />
              <CardTitle className="text-lg">Điểm danh</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Điểm danh hàng ngày cho học sinh</CardDescription>
            </CardContent>
          </Card>
        </Link>

        <Link href="/teacher/health">
          <Card className="transition-colors hover:bg-zinc-50">
            <CardHeader className="pb-2">
              <Heart className="h-6 w-6 text-muted-foreground" />
              <CardTitle className="text-lg">Sức khỏe</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Ghi nhận sức khỏe học sinh</CardDescription>
            </CardContent>
          </Card>
        </Link>

        <Link href="/teacher/posts">
          <Card className="transition-colors hover:bg-zinc-50">
            <CardHeader className="pb-2">
              <Users className="h-6 w-6 text-muted-foreground" />
              <CardTitle className="text-lg">Bài đăng</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Thông báo và nhận xét học sinh</CardDescription>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* My classes list */}
      {loading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : classes.length > 0 ? (
        <div>
          <h2 className="mb-3 text-lg font-semibold">Lớp được phân công</h2>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {classes.map((cls) => (
              <Card key={cls.class_id}>
                <CardContent className="py-4">
                  <p className="font-medium">{cls.name}</p>
                  <p className="mt-1 text-sm text-muted-foreground">{cls.school_year}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      ) : (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <GraduationCap className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Bạn chưa được phân công lớp nào</p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}