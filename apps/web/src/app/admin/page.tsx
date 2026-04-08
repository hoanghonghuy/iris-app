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
import { School, BookOpen, Users, GraduationCap, UsersRound, Loader2, Settings, TrendingUp, PieChart as PieChartIcon } from 'lucide-react';
import Link from 'next/link';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';

const attendanceData = [
  { name: 'T2', present: 95, absent: 5 },
  { name: 'T3', present: 98, absent: 2 },
  { name: 'T4', present: 92, absent: 8 },
  { name: 'T5', present: 96, absent: 4 },
  { name: 'T6', present: 99, absent: 1 },
  { name: 'T7', present: 85, absent: 15 },
];

const distributionData = [
  { name: 'Mầm Non', value: 35 },
  { name: 'Chồi Non', value: 45 },
  { name: 'Lá Non', value: 20 },
];
const COLORS = ['#10b981', '#f59e0b', '#3b82f6'];

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
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {/* School Stats - Teal Pill */}
            <Link href="/admin/schools" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md border-transparent hover:border-chart-1/30 relative overflow-hidden bg-card/60 backdrop-blur-sm">
                <CardContent className="p-5">
                  <div className="flex flex-col h-full gap-4">
                    <div className="flex items-center justify-between">
                      <div className="p-2.5 bg-chart-1/10 text-chart-1 rounded-xl transition-transform group-hover:scale-110 duration-300">
                        <School className="h-5 w-5" />
                      </div>
                      <span className="inline-flex items-center rounded-full bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium text-emerald-500">
                        +2 tuần này
                      </span>
                    </div>
                    <div>
                      <p className="text-3xl font-bold text-foreground group-hover:text-chart-1 transition-colors">{stats.total_schools}</p>
                      <p className="text-xs font-medium text-muted-foreground mt-1 uppercase tracking-wider">Trường học</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Classes Stats - Indigo Pill */}
            <Link href="/admin/classes" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md border-transparent hover:border-chart-3/30 relative overflow-hidden bg-card/60 backdrop-blur-sm">
                <CardContent className="p-5">
                  <div className="flex flex-col h-full gap-4">
                    <div className="flex items-center justify-between">
                      <div className="p-2.5 bg-chart-3/10 text-chart-3 rounded-xl transition-transform group-hover:scale-110 duration-300">
                        <BookOpen className="h-5 w-5" />
                      </div>
                      <span className="inline-flex items-center rounded-full bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium text-emerald-500">
                        +5% tháng này
                      </span>
                    </div>
                    <div>
                      <p className="text-3xl font-bold text-foreground group-hover:text-chart-3 transition-colors">{stats.total_classes}</p>
                      <p className="text-xs font-medium text-muted-foreground mt-1 uppercase tracking-wider">Lớp học</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Students Stats - Amber Pill */}
            <Link href="/admin/students" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md border-transparent hover:border-chart-2/30 relative overflow-hidden bg-card/60 backdrop-blur-sm">
                <CardContent className="p-5">
                  <div className="flex flex-col h-full gap-4">
                    <div className="flex items-center justify-between">
                      <div className="p-2.5 bg-chart-2/10 text-chart-2 rounded-xl transition-transform group-hover:scale-110 duration-300">
                        <GraduationCap className="h-5 w-5" />
                      </div>
                      <span className="inline-flex items-center rounded-full bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium text-emerald-500">
                        +12% học kỳ
                      </span>
                    </div>
                    <div>
                      <p className="text-3xl font-bold text-foreground group-hover:text-chart-2 transition-colors">{stats.total_students}</p>
                      <p className="text-xs font-medium text-muted-foreground mt-1 uppercase tracking-wider">Học sinh</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>

            {/* Users (Teachers + Parents) - Merging stats visually to save space or just another grid item */}
            <div className="grid grid-rows-2 gap-4">
              <Link href="/admin/teachers" className="group h-full">
                <Card className="h-full transition-all duration-300 hover:shadow-md border-transparent hover:border-chart-4/30 relative overflow-hidden bg-card/60 backdrop-blur-sm flex items-center">
                  <CardContent className="p-4 flex items-center justify-between w-full h-full pb-4">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-chart-4/10 text-chart-4 rounded-lg transition-transform group-hover:scale-110 duration-300">
                        <Users className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">Giáo viên</p>
                        <p className="text-xl font-bold text-foreground">{stats.total_teachers}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
              
              <Link href="/admin/parents" className="group h-full">
                <Card className="h-full transition-all duration-300 hover:shadow-md border-transparent hover:border-chart-5/30 relative overflow-hidden bg-card/60 backdrop-blur-sm flex items-center">
                  <CardContent className="p-4 flex items-center justify-between w-full h-full pb-4">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-chart-5/10 text-chart-5 rounded-lg transition-transform group-hover:scale-110 duration-300">
                        <UsersRound className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">Phụ huynh</p>
                        <p className="text-xl font-bold text-foreground">{stats.total_parents}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            </div>
          </div>

          {/* Analytics Overview Section (P2) */}
          <div className="grid gap-6 lg:grid-cols-3 mt-8">
            <Card className="lg:col-span-2 shadow-sm border-transparent hover:border-border transition-colors">
              <CardContent className="p-6">
                <div className="flex items-center justify-between mb-6">
                  <div>
                    <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
                      <TrendingUp className="h-5 w-5 text-primary" />
                      Điểm danh 7 ngày qua
                    </h2>
                    <p className="text-sm text-muted-foreground mt-1">Tỉ lệ tham gia trung bình (mock data)</p>
                  </div>
                </div>
                <div className="h-[280px] w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={attendanceData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="hsl(var(--border))" />
                      <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }} dy={10} />
                      <YAxis axisLine={false} tickLine={false} tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }} unit="%" />
                      <RechartsTooltip 
                        contentStyle={{ backgroundColor: 'hsl(var(--card))', borderColor: 'hsl(var(--border))', borderRadius: '8px' }}
                        itemStyle={{ color: 'hsl(var(--foreground))' }}
                      />
                      <Bar dataKey="present" name="Có mặt" fill="#10b981" radius={[4, 4, 0, 0]} stackId="a" />
                      <Bar dataKey="absent" name="Vắng mặt" fill="#ef4444" radius={[4, 4, 0, 0]} stackId="a" />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>

            <Card className="shadow-sm border-transparent hover:border-border transition-colors">
              <CardContent className="p-6 flex flex-col h-full">
                <div className="mb-4">
                  <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
                    <PieChartIcon className="h-5 w-5 text-chart-2" />
                    Phân bổ Học sinh
                  </h2>
                  <p className="text-sm text-muted-foreground mt-1">Theo độ tuổi (mock data)</p>
                </div>
                <div className="flex-1 flex items-center justify-center min-h-[220px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={distributionData}
                        cx="50%"
                        cy="50%"
                        innerRadius={60}
                        outerRadius={80}
                        paddingAngle={5}
                        dataKey="value"
                      >
                        {distributionData.map((entry, index) => (
                          <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <RechartsTooltip 
                         contentStyle={{ backgroundColor: 'hsl(var(--card))', borderColor: 'hsl(var(--border))', borderRadius: '8px' }}
                         itemStyle={{ color: 'hsl(var(--foreground))' }}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
                <div className="flex justify-center gap-4 mt-2">
                  {distributionData.map((entry, index) => (
                    <div key={entry.name} className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
                      <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[index % COLORS.length] }}></div>
                      {entry.name}
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="grid gap-6 lg:grid-cols-3 mt-8">
            <div className="lg:col-span-2 flex flex-col">
              <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2 mb-4">
                Lịch sử hoạt động
              </h2>
              <Card className="border-transparent shadow-sm bg-card/40 flex-1">
                <CardContent className="p-0">
                  <div className="divide-y divide-border/50">
                    <div className="flex items-center justify-between p-4 px-6 hover:bg-muted/30 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="h-8 w-8 rounded-full bg-chart-1/20 flex items-center justify-center text-chart-1">
                          <School className="h-4 w-4" />
                        </div>
                        <div>
                          <p className="text-sm font-medium">Trường Mầm non Hoa Mai Quận 7 vừa cập nhật thông tin.</p>
                          <p className="text-xs text-muted-foreground">Bởi: school-admin@iris.local</p>
                        </div>
                      </div>
                      <span className="text-xs text-muted-foreground truncate ml-4">10 phút trước</span>
                    </div>
                    <div className="flex items-center justify-between p-4 px-6 hover:bg-muted/30 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="h-8 w-8 rounded-full bg-chart-2/20 flex items-center justify-center text-chart-2">
                          <GraduationCap className="h-4 w-4" />
                        </div>
                        <div>
                          <p className="text-sm font-medium">Đã thêm 5 học sinh mới vào Lớp Lá Non.</p>
                          <p className="text-xs text-muted-foreground">Bởi: admin@iris.local</p>
                        </div>
                      </div>
                      <span className="text-xs text-muted-foreground truncate ml-4">2 giờ trước</span>
                    </div>
                    <div className="flex items-center justify-between p-4 px-6 hover:bg-muted/30 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="h-8 w-8 rounded-full bg-chart-4/20 flex items-center justify-center text-chart-4">
                          <Users className="h-4 w-4" />
                        </div>
                        <div>
                          <p className="text-sm font-medium">Phân công Giáo viên Nguyễn Thị Lan Anh vào Lớp Lá Non.</p>
                          <p className="text-xs text-muted-foreground">Bởi: school-admin@iris.local</p>
                        </div>
                      </div>
                      <span className="text-xs text-muted-foreground truncate ml-4">Hôm qua</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="lg:col-span-1 flex flex-col">
              <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2 mb-4">
                Quản lý Nhanh
              </h2>
              <div className="grid gap-3 grid-cols-2 lg:grid-cols-1 shrink-0 h-full">
                <Link href="/admin/schools" className="block w-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-chart-1/30 shadow-sm h-full">
                    <CardContent className="p-3.5 flex items-center justify-start gap-4">
                      <div className="p-3 bg-chart-1/10 text-chart-1 rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <School className="h-5 w-5" />
                      </div>
                      <div className="text-left font-semibold text-sm text-foreground group-hover:text-chart-1 transition-colors">
                        QL Trường học
                      </div>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/admin/users" className="block w-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-chart-3/30 shadow-sm h-full">
                    <CardContent className="p-3.5 flex items-center justify-start gap-4">
                      <div className="p-3 bg-chart-3/10 text-chart-3 rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <Users className="h-5 w-5" />
                      </div>
                      <div className="text-left font-semibold text-sm text-foreground group-hover:text-chart-3 transition-colors">
                        QL Người dùng
                      </div>
                    </CardContent>
                  </Card>
                </Link>

                <Link href="/admin/classes" className="block w-full">
                  <Card className="group hover:bg-muted/50 transition-colors cursor-pointer border-transparent hover:border-chart-2/30 shadow-sm h-full">
                    <CardContent className="p-3.5 flex items-center justify-start gap-4">
                      <div className="p-3 bg-chart-2/10 text-chart-2 rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                        <BookOpen className="h-5 w-5" />
                      </div>
                      <div className="text-left font-semibold text-sm text-foreground group-hover:text-chart-2 transition-colors">
                        QL Lớp học
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