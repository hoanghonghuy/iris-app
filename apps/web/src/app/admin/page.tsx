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
        <div className="flex items-center gap-2 mb-1">
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
        <p className="text-muted-foreground text-base max-w-2xl mt-1">
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
          <div className="grid gap-5 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
            {/* School Stats - Teal Pill */}
            <Link href="/admin/schools" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-teal-500/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex flex-col h-full justify-between gap-4">
                    <div className="p-3 bg-teal-500/10 text-teal-600 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <School className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground mb-1 uppercase tracking-wider">Trường học</p>
                      <p className="text-3xl font-bold text-foreground group-hover:text-teal-600 transition-colors">{stats.total_schools}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Classes Stats - Indigo Pill */}
            <Link href="/admin/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-indigo-500/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex flex-col h-full justify-between gap-4">
                    <div className="p-3 bg-indigo-500/10 text-indigo-600 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <BookOpen className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground mb-1 uppercase tracking-wider">Lớp học</p>
                      <p className="text-3xl font-bold text-foreground group-hover:text-indigo-600 transition-colors">{stats.total_classes}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Teachers Stats - Rose Pill */}
            <Link href="/admin/teachers" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-rose-500/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex flex-col h-full justify-between gap-4">
                    <div className="p-3 bg-rose-500/10 text-rose-600 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <Users className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground mb-1 uppercase tracking-wider">Giáo viên</p>
                      <p className="text-3xl font-bold text-foreground group-hover:text-rose-600 transition-colors">{stats.total_teachers}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Students Stats - Amber Pill */}
            <Link href="/admin/students" className="group md:col-span-2 xl:col-span-1">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-amber-500/30 relative overflow-hidden bg-gradient-to-br hover:from-card hover:to-amber-50/50 dark:hover:to-amber-950/20">
                <CardContent className="p-6">
                  <div className="flex flex-col h-full justify-between gap-4">
                    <div className="p-3 bg-amber-500/10 text-amber-600 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <GraduationCap className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground mb-1 uppercase tracking-wider">Học sinh</p>
                      <p className="text-3xl font-bold text-foreground group-hover:text-amber-600 transition-colors">{stats.total_students}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Parents Stats - Blue Pill */}
            <Link href="/admin/parents" className="group md:col-span-2 xl:col-span-1">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-blue-500/30 relative overflow-hidden">
                <CardContent className="p-6">
                  <div className="flex flex-col h-full justify-between gap-4">
                    <div className="p-3 bg-blue-500/10 text-blue-600 rounded-2xl w-fit transition-transform group-hover:scale-110 duration-300">
                      <UsersRound className="h-6 w-6" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-muted-foreground mb-1 uppercase tracking-wider">Phụ huynh</p>
                      <p className="text-3xl font-bold text-foreground group-hover:text-blue-600 transition-colors">{stats.total_parents}</p>
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
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                
                <Link href="/admin/schools">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-5 flex items-start gap-4 h-full flex-col sm:flex-row">
                      <div className="p-3 bg-primary/10 text-primary rounded-xl shrink-0">
                        <Settings className="h-6 w-6" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center justify-between">
                          <p className="font-semibold text-foreground group-hover:text-primary transition-colors">Quản lý Trường học</p>
                          <ChevronRight className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                        </div>
                        <p className="text-sm text-muted-foreground line-clamp-2">
                          Thêm, sửa, thiết lập cấu hình cơ bản cho các trường học trong hệ thống.
                        </p>
                      </div>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/admin/users">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-5 flex items-start gap-4 h-full flex-col sm:flex-row">
                      <div className="p-3 bg-indigo-500/10 text-indigo-600 rounded-xl shrink-0">
                        <Users className="h-6 w-6" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center justify-between">
                          <p className="font-semibold text-foreground group-hover:text-indigo-600 transition-colors">Quản lý Người dùng</p>
                          <ChevronRight className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                        </div>
                        <p className="text-sm text-muted-foreground line-clamp-2">
                          Kiểm soát tài khoản truy cập hệ thống của Quản trị viên, Giáo viên, Phụ huynh.
                        </p>
                      </div>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/admin/classes">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-border shadow-sm h-full">
                    <CardContent className="p-5 flex items-start gap-4 h-full flex-col sm:flex-row">
                      <div className="p-3 bg-amber-500/10 text-amber-600 rounded-xl shrink-0">
                        <BookOpen className="h-6 w-6" />
                      </div>
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center justify-between">
                          <p className="font-semibold text-foreground group-hover:text-amber-600 transition-colors">Thiết lập Lớp học</p>
                          <ChevronRight className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                        </div>
                        <p className="text-sm text-muted-foreground line-clamp-2">
                          Tạo mới lớp học, phân công giáo viên chủ nhiệm và quản lý sĩ số.
                        </p>
                      </div>
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