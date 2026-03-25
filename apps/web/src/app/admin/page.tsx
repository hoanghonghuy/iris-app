/**
 * Admin Dashboard
 * Trang chủ dành cho Admin sau khi đăng nhập, chuẩn Minimalist Pastel.
 */
"use client";

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/providers/AuthProvider';
import { Card, CardContent } from '@/components/ui/card';
import { adminApi } from '@/lib/api/admin.api';
import { AdminAnalytics } from '@/types';
import { School, BookOpen, Users, GraduationCap, UsersRound, Loader2, ChevronRight, Settings } from 'lucide-react';
import Link from 'next/link';

export default function AdminDashboard() {
  const { user, role } = useAuth();
  const [stats, setStats] = useState<AdminAnalytics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const data = await adminApi.getAnalytics();
        setStats(data);
      } catch (error) {
        console.error("Lỗi khi tải thống kê", error);
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, []);

  return (
    <div className="space-y-8 pb-8">
      {/* Hero Header Area */}
      <div className="flex flex-col gap-1.5 animate-in fade-in slide-in-from-bottom-2 duration-500">
        <div className="hidden md:flex items-center gap-2 mb-1">
          <span className="bg-primary/15 text-primary px-3 py-1 rounded-full text-xs font-semibold tracking-wide uppercase">
            Admin Portal
          </span>
          <span className="bg-muted text-muted-foreground px-3 py-1 rounded-full text-xs font-medium tracking-wide">
            {role === 'SUPER_ADMIN' ? 'Hệ thống' : 'Cơ sở'}
          </span>
        </div>
        <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-foreground">
          Xin chào, {user?.full_name || user?.email?.split('@')[0]}
        </h1>
        <p className="text-muted-foreground text-base max-w-2xl mt-1 hidden md:block">
          Trung tâm điều khiển hệ thống. Hiện tại có {stats?.total_schools || 0} trường và {stats?.total_students || 0} học sinh đang hoạt động.
        </p>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : stats ? (
        <>
          {/* Stats Overview - Premium Bento Grid */}
          <div className="grid gap-2.5 md:gap-3 grid-cols-2 md:grid-cols-3 xl:grid-cols-5">
            {/* School Stats - Teal Pill */}
            <Link href="/admin/schools" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-1/30 relative overflow-hidden">
                <CardContent className="p-3.5 md:p-5">
                  <div className="flex flex-col h-full justify-between gap-2.5 md:gap-4">
                    <div className="p-3 bg-chart-1/10 text-chart-1 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <School className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-[11px] md:text-xs font-semibold text-muted-foreground mb-0.5 uppercase tracking-wider">Trường học</p>
                      <p className="text-2xl md:text-3xl font-bold text-foreground group-hover:text-chart-1 transition-colors">{stats.total_schools}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Classes Stats - Indigo Pill */}
            <Link href="/admin/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-3/30 relative overflow-hidden">
                <CardContent className="p-3.5 md:p-5">
                  <div className="flex flex-col h-full justify-between gap-2.5 md:gap-4">
                    <div className="p-3 bg-chart-3/10 text-chart-3 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <BookOpen className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-[11px] md:text-xs font-semibold text-muted-foreground mb-0.5 uppercase tracking-wider">Lớp học</p>
                      <p className="text-2xl md:text-3xl font-bold text-foreground group-hover:text-chart-3 transition-colors">{stats.total_classes}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Teachers Stats - Rose Pill */}
            <Link href="/admin/teachers" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-4/30 relative overflow-hidden">
                <CardContent className="p-3.5 md:p-5">
                  <div className="flex flex-col h-full justify-between gap-2.5 md:gap-4">
                    <div className="p-3 bg-chart-4/10 text-chart-4 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <Users className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-[11px] md:text-xs font-semibold text-muted-foreground mb-0.5 uppercase tracking-wider">Giáo viên</p>
                      <p className="text-2xl md:text-3xl font-bold text-foreground group-hover:text-chart-4 transition-colors">{stats.total_teachers}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Students Stats - Amber Pill */}
            <Link href="/admin/students" className="group md:col-span-2 xl:col-span-1">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-2/30 relative overflow-hidden bg-gradient-to-br hover:from-card hover:to-chart-2/10 dark:hover:to-chart-2/20">
                <CardContent className="p-3.5 md:p-5">
                  <div className="flex flex-col h-full justify-between gap-2.5 md:gap-4">
                    <div className="p-3 bg-chart-2/10 text-chart-2 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <GraduationCap className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-[11px] md:text-xs font-semibold text-muted-foreground mb-0.5 uppercase tracking-wider">Học sinh</p>
                      <p className="text-2xl md:text-3xl font-bold text-foreground group-hover:text-chart-2 transition-colors">{stats.total_students}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Parents Stats - Blue Pill */}
            <Link href="/admin/parents" className="group col-span-2 sm:col-span-1 md:col-span-2 xl:col-span-1">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-5/30 relative overflow-hidden">
                <CardContent className="p-3.5 md:p-5">
                  <div className="flex flex-col h-full justify-between gap-2.5 md:gap-4">
                    <div className="p-3 bg-chart-5/10 text-chart-5 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <UsersRound className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-[11px] md:text-xs font-semibold text-muted-foreground mb-0.5 uppercase tracking-wider">Phụ huynh</p>
                      <p className="text-2xl md:text-3xl font-bold text-foreground group-hover:text-chart-5 transition-colors">{stats.total_parents}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <div className="lg:col-span-3 space-y-4">
              <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2 mt-2">
                Quản lý Nhanh
              </h2>
              <div className="grid gap-2 md:gap-3 grid-cols-3">
                
                <Link href="/admin/schools" className="h-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-2 md:p-3 flex flex-col items-center justify-center gap-1.5 md:gap-2 h-full text-center">
                      <div className="p-2.5 md:p-3 bg-primary/10 text-primary rounded-xl md:rounded-2xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <Settings className="h-5 w-5 md:h-6 md:w-6" />
                      </div>
                      <p className="font-medium md:font-semibold text-[11px] md:text-xs text-foreground group-hover:text-primary transition-colors">QL Trường</p>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/admin/users" className="h-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-2 md:p-3 flex flex-col items-center justify-center gap-1.5 md:gap-2 h-full text-center">
                      <div className="p-2.5 md:p-3 bg-chart-3/10 text-chart-3 rounded-xl md:rounded-2xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <Users className="h-5 w-5 md:h-6 md:w-6" />
                      </div>
                      <p className="font-medium md:font-semibold text-[11px] md:text-xs text-foreground group-hover:text-chart-3 transition-colors">QL Users</p>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/admin/classes" className="h-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-2 md:p-3 flex flex-col items-center justify-center gap-1.5 md:gap-2 h-full text-center">
                      <div className="p-2.5 md:p-3 bg-chart-2/10 text-chart-2 rounded-xl md:rounded-2xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <BookOpen className="h-5 w-5 md:h-6 md:w-6" />
                      </div>
                      <p className="font-medium md:font-semibold text-[11px] md:text-xs text-foreground group-hover:text-chart-2 transition-colors">QL Lớp</p>
                    </CardContent>
                  </Card>
                </Link>
                
              </div>
            </div>
          </div>
        </>
      ) : null}
    </div>
  );
}