/**
 * Teacher Dashboard
 * Trang tổng quan cho giáo viên chuẩn Minimalist Pastel (NurturedLayer style).
 */
"use client";

import React, { useEffect, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, TeacherAnalytics } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { GraduationCap, Users, ClipboardCheck, Heart, Loader2, BookOpen, MessageSquare, ChevronRight } from "lucide-react";
import Link from "next/link";

export default function TeacherDashboard() {
  const { user } = useAuth();
  const [classes, setClasses] = useState<Class[]>([]);
  const [stats, setStats] = useState<TeacherAnalytics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [classesData, statsData] = await Promise.all([
          teacherApi.getMyClasses(),
          teacherApi.getAnalytics()
        ]);
        setClasses(classesData || []);
        setStats(statsData);
      } catch (error) {
        console.error("Lỗi khi tải dữ liệu giáo viên", error);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  return (
    <div className="space-y-8 pb-8">
      {/* Hero Header Area (No card, pure typography & airy spacing) */}
      <div className="flex flex-col gap-1.5 animate-in fade-in slide-in-from-bottom-2 duration-500">
        <div className="flex items-center gap-2 mb-1">
          <span className="bg-primary/15 text-primary px-3 py-1 rounded-full text-xs font-semibold tracking-wide uppercase">
            Teacher Portal
          </span>
        </div>
        <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-foreground">
          Xin chào, {user?.full_name || user?.email?.split('@')[0]}
        </h1>
        <p className="text-muted-foreground text-base max-w-2xl mt-1">
          Hôm nay bạn có {stats?.total_classes || 0} lớp học và {stats?.total_students || 0} học sinh cần theo dõi. Chúc một ngày làm việc hiệu quả!
        </p>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : (
        <>
          {/* Stats Overview - Bento Grid */}
          <div className="grid gap-5 md:grid-cols-3">
            <Link href="/teacher/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-primary/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="space-y-3">
                      <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">Lớp phụ trách</p>
                      <p className="text-4xl font-bold text-foreground group-hover:text-primary transition-colors">{stats?.total_classes || 0}</p>
                    </div>
                    {/* Pill Icon Badge */}
                    <div className="p-3.5 bg-primary/10 rounded-2xl text-primary transition-transform group-hover:scale-110 duration-300">
                      <BookOpen className="h-6 w-6" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            <Link href="/teacher/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-amber-500/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="space-y-3">
                      <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">Học sinh quản lý</p>
                      <p className="text-4xl font-bold text-foreground group-hover:text-amber-500 transition-colors">{stats?.total_students || 0}</p>
                    </div>
                    {/* Amber Pill Icon Badge */}
                    <div className="p-3.5 bg-amber-500/10 rounded-2xl text-amber-500 transition-transform group-hover:scale-110 duration-300">
                      <GraduationCap className="h-6 w-6" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            <Link href="/teacher/posts" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-blue-500/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="space-y-3">
                      <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">Bài đăng đã tạo</p>
                      <p className="text-4xl font-bold text-foreground group-hover:text-blue-500 transition-colors">{stats?.total_posts || 0}</p>
                    </div>
                    {/* Blue Pill Icon Badge */}
                    <div className="p-3.5 bg-blue-500/10 rounded-2xl text-blue-500 transition-transform group-hover:scale-110 duration-300">
                      <MessageSquare className="h-6 w-6" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {/* Quick Actions (Bento Box style) */}
            <div className="lg:col-span-1 space-y-4">
              <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
                Hoạt động Nhanh
              </h2>
              <div className="grid gap-3">
                <Link href="/teacher/attendance">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm">
                    <CardContent className="p-4 flex items-center gap-4">
                      <div className="p-2.5 bg-green-500/10 text-green-600 rounded-xl">
                        <ClipboardCheck className="h-5 w-5" />
                      </div>
                      <div className="flex-1">
                        <p className="font-semibold text-foreground group-hover:text-green-600 transition-colors">Điểm danh</p>
                        <p className="text-xs text-muted-foreground">Chốt sĩ số hàng ngày</p>
                      </div>
                      <ChevronRight className="h-5 w-5 text-muted-foreground opacity-50 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/teacher/health">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm">
                    <CardContent className="p-4 flex items-center gap-4">
                      <div className="p-2.5 bg-rose-500/10 text-rose-600 rounded-xl">
                        <Heart className="h-5 w-5" />
                      </div>
                      <div className="flex-1">
                        <p className="font-semibold text-foreground group-hover:text-rose-600 transition-colors">Sức khỏe</p>
                        <p className="text-xs text-muted-foreground">Cập nhật hồ sơ thể chất</p>
                      </div>
                      <ChevronRight className="h-5 w-5 text-muted-foreground opacity-50 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/teacher/posts">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm">
                    <CardContent className="p-4 flex items-center gap-4">
                      <div className="p-2.5 bg-blue-500/10 text-blue-600 rounded-xl">
                        <Users className="h-5 w-5" />
                      </div>
                      <div className="flex-1">
                        <p className="font-semibold text-foreground group-hover:text-blue-600 transition-colors">Bảng tin</p>
                        <p className="text-xs text-muted-foreground">Gửi thông báo lớp</p>
                      </div>
                      <ChevronRight className="h-5 w-5 text-muted-foreground opacity-50 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                    </CardContent>
                  </Card>
                </Link>
              </div>
            </div>

            {/* My Classes List - Spans 2 columns */}
            <div className="lg:col-span-2 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
                  Lớp Được Phân Công
                </h2>
                <Link href="/teacher/classes" className="text-sm font-medium text-primary hover:underline">
                  Xem tất cả
                </Link>
              </div>
              
              {classes.length > 0 ? (
                <div className="grid gap-4 sm:grid-cols-2">
                  {classes.map((cls) => (
                    <Card key={cls.class_id} className="group hover:shadow-md transition-all duration-300">
                      <CardContent className="p-5">
                        <div className="flex justify-between items-start mb-4">
                          <div className="p-2 bg-primary/10 text-primary rounded-lg">
                            <BookOpen className="h-5 w-5" />
                          </div>
                          <span className="bg-muted text-muted-foreground px-2.5 py-1 rounded-full text-xs font-medium">
                            {cls.school_year}
                          </span>
                        </div>
                        <h3 className="text-xl font-bold group-hover:text-primary transition-colors">{cls.name}</h3>
                        <p className="mt-1 text-sm text-muted-foreground line-clamp-2">
                          Lớp phụ trách chính thức cho năm học hiện tại.
                        </p>
                        <div className="mt-4 pt-4 border-t border-border flex items-center justify-between">
                          <span className="text-xs font-medium text-muted-foreground flex items-center gap-1.5">
                            <Users className="h-4 w-4" /> Đang hoạt động
                          </span>
                          <Link href="/teacher/classes" className="text-sm font-medium text-primary flex items-center group-hover:translate-x-1 transition-transform">
                            Truy cập <ChevronRight className="h-4 w-4 ml-0.5" />
                          </Link>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              ) : (
                <Card className="border-dashed shadow-none">
                  <CardContent className="flex flex-col items-center justify-center py-16 text-center">
                    <div className="p-4 bg-muted rounded-full mb-4">
                      <GraduationCap className="h-8 w-8 text-muted-foreground" />
                    </div>
                    <p className="font-semibold text-lg text-foreground">Chưa có lớp học</p>
                    <p className="mt-1 text-sm text-muted-foreground max-w-sm">Bạn chưa được phân bổ vào danh sách lớp giảng dạy nào. Vui lòng liên hệ Admin.</p>
                  </CardContent>
                </Card>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  );
}