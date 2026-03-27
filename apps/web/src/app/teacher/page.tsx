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
    <div className="h-full flex flex-col space-y-6">
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
          {/* Stats Overview - Bento Grid */}
          {/* Stats Overview - Bento Grid */}
          <div className="grid gap-2.5 md:gap-3 grid-cols-3 shrink-0">
            <Link href="/teacher/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-primary/30 relative overflow-hidden">
                <CardContent className="p-3 md:p-4 flex flex-col items-center justify-center gap-2 h-full text-center">
                  <div className="p-2.5 bg-primary/10 rounded-2xl text-primary transition-transform group-hover:scale-110 duration-300">
                    <BookOpen className="h-5 w-5" />
                  </div>
                  <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest mt-1">Lớp</p>
                  <p className="text-xl font-bold text-foreground group-hover:text-primary leading-none">{stats?.total_classes || 0}</p>
                </CardContent>
              </Card>
            </Link>

            <Link href="/teacher/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-2/30 relative overflow-hidden">
                <CardContent className="p-3 md:p-4 flex flex-col items-center justify-center gap-2 h-full text-center">
                  <div className="p-2.5 bg-chart-2/10 rounded-2xl text-chart-2 transition-transform group-hover:scale-110 duration-300">
                    <GraduationCap className="h-5 w-5" />
                  </div>
                  <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest mt-1">Trẻ</p>
                  <p className="text-xl font-bold text-foreground group-hover:text-chart-2 leading-none">{stats?.total_students || 0}</p>
                </CardContent>
              </Card>
            </Link>

            <Link href="/teacher/posts" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-3/30 relative overflow-hidden">
                <CardContent className="p-3 md:p-4 flex flex-col items-center justify-center gap-2 h-full text-center">
                  <div className="p-2.5 bg-chart-3/10 rounded-2xl text-chart-3 transition-transform group-hover:scale-110 duration-300">
                    <MessageSquare className="h-5 w-5" />
                  </div>
                  <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest mt-1">Tin</p>
                  <p className="text-xl font-bold text-foreground group-hover:text-chart-3 leading-none">{stats?.total_posts || 0}</p>
                </CardContent>
              </Card>
            </Link>
          </div>

          <div className="grid gap-6 lg:grid-cols-3 flex-1 min-h-0 overflow-hidden">
            {/* Quick Actions (Bento Box style) */}
            <div className="lg:col-span-1 flex flex-col min-h-0 space-y-4">
              <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2 shrink-0">
                Hoạt động Nhanh
              </h2>
              <div className="grid gap-2 md:gap-3 grid-cols-3 shrink-0">
                <Link href="/teacher/attendance" className="h-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-2 flex flex-col items-center justify-center gap-1.5 h-full text-center">
                      <div className="p-2 bg-success/10 text-success rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <ClipboardCheck className="h-5 w-5" />
                      </div>
                      <p className="font-semibold text-[10px] text-foreground group-hover:text-success transition-colors">Điểm danh</p>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/teacher/health" className="h-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-2 flex flex-col items-center justify-center gap-1.5 h-full text-center">
                      <div className="p-2 bg-destructive/10 text-destructive rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <Heart className="h-5 w-5" />
                      </div>
                      <p className="font-semibold text-[10px] text-foreground group-hover:text-destructive transition-colors">Sức khỏe</p>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/teacher/posts" className="h-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-2 flex flex-col items-center justify-center gap-1.5 h-full text-center">
                      <div className="p-2 bg-primary/10 text-primary rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <MessageSquare className="h-5 w-5" />
                      </div>
                      <p className="font-semibold text-[10px] text-foreground group-hover:text-primary transition-colors">Bảng tin</p>
                    </CardContent>
                  </Card>
                </Link>
              </div>

              {/* Extra spacing or other content could go here if needed */}
              <div className="flex-1" />
            </div>

            {/* My Classes List - Spans 2 columns */}
            <div className="lg:col-span-2 flex flex-col min-h-0 space-y-4">
              <div className="flex items-center justify-between shrink-0">
                <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2">
                  Lớp Được Phân Công
                </h2>
                <Link href="/teacher/classes" className="text-xs font-semibold text-primary hover:underline transition-all hover:translate-x-0.5">
                  Xem tất cả
                </Link>
              </div>
              
              <div className="flex-1 overflow-y-auto pr-2 custom-scrollbar">
                {classes.length > 0 ? (
                  <div className="grid gap-4 sm:grid-cols-2">
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
                  <Card className="border-dashed shadow-none">
                    {/* ... empty ... */}
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