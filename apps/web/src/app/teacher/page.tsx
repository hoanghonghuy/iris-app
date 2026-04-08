/**
 * Teacher Dashboard
 * Trang tổng quan cho giáo viên chuẩn Minimalist Pastel (NurturedLayer style).
 */
"use client";

import React, { useEffect, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, TeacherAnalytics } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardContent } from "@/components/ui/card";
import { GraduationCap, Users, ClipboardCheck, Heart, Loader2, BookOpen, MessageSquare, ChevronRight, Calendar } from "lucide-react";
import Link from "next/link";
import { cn } from "@/lib/utils";

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
    <div className="flex h-full min-h-0 flex-col gap-6">
      {/* Hero Header Area (No card, pure typography & airy spacing) */}
      <div className="flex flex-col gap-1.5 animate-in fade-in slide-in-from-bottom-2 duration-500 shrink-0">
        <div className="hidden md:flex items-center gap-2 mb-1">
          <span className="bg-primary/15 text-primary px-3 py-1 rounded-full text-xs font-semibold tracking-wide uppercase">
            Teacher Portal
          </span>
        </div>
        <h1 className="text-2xl md:text-3xl font-extrabold tracking-tight text-foreground">
          Xin chào, {user?.full_name || user?.email?.split('@')[0]}
        </h1>
        <p className="text-muted-foreground text-sm max-w-2xl mt-1 hidden md:block">
          Hôm nay bạn có {stats?.total_students || 0} học sinh cần theo dõi.
        </p>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : (
        <>
          {/* Today's Status Card (P1) */}
          <Card className="shrink-0 bg-card border border-border shadow-sm relative overflow-hidden">
            <CardContent className="p-5 md:p-6">
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
                <div className="space-y-1.5 flex-1">
                  <h2 className="text-xl font-bold text-foreground flex items-center gap-2">
                    <Calendar className="h-5 w-5 text-primary" />
                    Hôm nay: {new Date().toLocaleDateString('vi-VN', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })}
                  </h2>
                  <p className="text-sm text-muted-foreground">
                    Lớp phụ trách chính: <span className="font-semibold text-foreground">{classes.length > 0 ? classes[0].name : "Chưa có"}</span>
                  </p>
                </div>
                
                <div className="flex flex-col sm:flex-row sm:items-center gap-4 md:gap-8">
                  {/* Attendance Progress Line */}
                  <div className="space-y-2 min-w-[140px]">
                    <div className="flex items-center justify-between text-sm">
                      <span className="font-semibold flex items-center gap-1.5 text-success">
                        <ClipboardCheck className="h-4 w-4" /> Điểm danh
                      </span>
                      <span className="font-bold">{stats?.total_students ? `${stats?.total_students}/${stats?.total_students}` : "0/0"}</span>
                    </div>
                    <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                      <div className="h-full bg-success rounded-full" style={{ width: '100%' }}></div>
                    </div>
                  </div>
                  
                  <div className="h-10 w-px bg-border/50 hidden sm:block"></div>
                  
                  {/* Health Alerts */}
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-muted text-muted-foreground rounded-lg">
                      <Heart className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Sức khỏe</p>
                      <p className="text-sm font-bold text-foreground">0 cảnh báo</p>
                    </div>
                  </div>
                  
                  <div className="h-10 w-px bg-border/50 hidden sm:block"></div>

                  {/* Posts Activity */}
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-primary/10 text-primary rounded-lg">
                      <MessageSquare className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Bài viết (7 ngày)</p>
                      <p className="text-sm font-bold text-foreground">{stats?.total_posts || 0} bài mới</p>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <div className="grid flex-1 min-h-0 gap-6 overflow-hidden lg:grid-cols-3">
            {/* Quick Actions (Bento Box style) */}
            <div className="flex min-h-0 flex-col gap-4 lg:col-span-1">
              <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2 shrink-0">
                Hoạt động Nhanh
              </h2>
              <div className="grid gap-3 grid-cols-1 shrink-0">
                <Link href="/teacher/attendance" className="block w-full">
                  <Card className="group hover:shadow-md transition-all cursor-pointer border-transparent hover:border-success/30 shadow-sm bg-card/60">
                    <CardContent className="p-3.5 flex items-center justify-between gap-3">
                      <div className="flex items-center gap-3">
                        <div className="p-2.5 bg-success/10 text-success rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                          <ClipboardCheck className="h-5 w-5" />
                        </div>
                        <div className="text-left">
                          <p className="font-semibold text-sm text-foreground group-hover:text-success transition-colors">Điểm danh</p>
                          <p className="text-xs font-medium text-muted-foreground mt-0.5">3/3 đã điểm danh</p>
                        </div>
                      </div>
                      <ChevronRight className="h-4 w-4 text-muted-foreground opacity-20 group-hover:opacity-100 transition-opacity group-hover:text-success group-hover:translate-x-1" />
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/teacher/health" className="block w-full">
                  <Card className="group hover:shadow-md transition-all cursor-pointer border-transparent hover:border-orange-500/30 shadow-sm bg-card/60">
                    <CardContent className="p-3.5 flex items-center justify-between gap-3">
                      <div className="flex items-center gap-3">
                        <div className="p-2.5 bg-orange-500/10 text-orange-500 rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                          <Heart className="h-5 w-5" />
                        </div>
                        <div className="text-left">
                          <p className="font-semibold text-sm text-foreground group-hover:text-orange-500 transition-colors">Sức khỏe</p>
                          <p className="text-xs font-medium text-orange-600 dark:text-orange-400 mt-0.5">1 cảnh báo mới</p>
                        </div>
                      </div>
                      <ChevronRight className="h-4 w-4 text-muted-foreground opacity-20 group-hover:opacity-100 transition-opacity group-hover:text-orange-500 group-hover:translate-x-1" />
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/teacher/posts" className="block w-full">
                  <Card className="group hover:shadow-md transition-all cursor-pointer border-transparent hover:border-primary/30 shadow-sm bg-card/60">
                    <CardContent className="p-3.5 flex items-center justify-between gap-3">
                      <div className="flex items-center gap-3">
                        <div className="p-2.5 bg-primary/10 text-primary rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                          <MessageSquare className="h-5 w-5" />
                        </div>
                        <div className="text-left">
                          <p className="font-semibold text-sm text-foreground group-hover:text-primary transition-colors">Bảng tin</p>
                          <p className="text-xs font-medium text-muted-foreground mt-0.5">2 bài tuần này</p>
                        </div>
                      </div>
                      <ChevronRight className="h-4 w-4 text-muted-foreground opacity-20 group-hover:opacity-100 transition-opacity group-hover:text-primary group-hover:translate-x-1" />
                    </CardContent>
                  </Card>
                </Link>
              </div>

              {/* Extra spacing or other content could go here if needed */}
              <div className="flex-1" />
            </div>

            {/* My Classes List - Spans 2 columns */}
            <div className="flex min-h-0 flex-col gap-4 lg:col-span-2">
              <div className="flex items-center justify-between shrink-0">
                <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2">
                  Lớp Được Phân Công
                </h2>
                <Link href="/teacher/classes" className="text-xs font-semibold text-primary hover:underline transition-all hover:translate-x-0.5">
                  Xem tất cả
                </Link>
              </div>
              
              <div className="custom-scrollbar flex-1 min-h-0 overflow-y-auto pr-2">
                {classes.length > 0 ? (
                  <div className={cn("grid gap-4 pb-4", classes.length === 1 ? "grid-cols-1" : "sm:grid-cols-2")}>
                    {classes.map((cls) => (
                      <Card key={cls.class_id} className="group hover:shadow-lg transition-all border-border/50 hover:border-primary/20 bg-card/50 backdrop-blur-sm">
                        <CardContent className="p-4 md:p-5">
                          <div className="flex justify-between items-start mb-4">
                            <div className="p-2 bg-primary/10 text-primary rounded-xl transition-transform group-hover:scale-110">
                              <BookOpen className="h-5 w-5" />
                            </div>
                            <span className="bg-primary/5 text-primary px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider">
                              {cls.school_year}
                            </span>
                          </div>
                          <h3 className="text-xl font-bold group-hover:text-primary transition-colors">{cls.name}</h3>
                          <p className="mt-2 text-sm text-muted-foreground/80 line-clamp-2 leading-relaxed">
                            Lớp phụ trách chính thức cho năm học hiện tại. Nhấn để quản lý học sinh và điểm danh.
                          </p>
                          <div className="mt-5 pt-4 border-t border-border/40 flex items-center justify-between">
                            <span className="text-[11px] font-bold text-success/80 flex items-center gap-1.5 uppercase tracking-wide">
                              <Users className="h-3.5 w-3.5" /> Đang hoạt động
                            </span>
                            <Link href="/teacher/classes" className="text-sm font-bold text-primary flex items-center group-hover:translate-x-1 transition-transform">
                              Quản lý <ChevronRight className="h-4 w-4 ml-0.5" />
                            </Link>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                ) : (
                  <Card className="border-dashed bg-card/40 shadow-none">
                    <CardContent className="flex min-h-[280px] flex-col items-center justify-center py-16 text-center">
                      <div className="mb-4 rounded-full bg-muted p-4">
                        <GraduationCap className="h-8 w-8 text-muted-foreground" />
                      </div>
                      <p className="text-lg font-semibold text-foreground">Chưa có lớp học</p>
                      <p className="mt-2 max-w-sm text-sm text-muted-foreground">
                        Bạn chưa được phân bổ vào danh sách lớp giảng dạy nào. Vui lòng liên hệ Admin.
                      </p>
                    </CardContent>
                  </Card>
                )}
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
