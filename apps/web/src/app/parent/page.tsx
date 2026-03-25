/**
 * Parent Dashboard
 * Trang tổng quan cho phụ huynh chuẩn Minimalist Pastel.
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Student, Post } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardContent } from "@/components/ui/card";
import { Users, MessageSquare, Loader2, Baby, CalendarClock, ChevronRight, Activity, HeartPulse } from "lucide-react";
import Link from "next/link";

const postTypeConfig: Record<string, { label: string, colorClass: string }> = {
  announcement: { label: "Thông báo", colorClass: "bg-primary/10 text-primary" },
  activity: { label: "Hoạt động", colorClass: "bg-chart-2/10 text-chart-2" },
  daily_note: { label: "Nhận xét", colorClass: "bg-chart-1/10 text-chart-1" },
  health_note: { label: "Sức khỏe", colorClass: "bg-destructive/10 text-destructive" },
};

export default function ParentDashboard() {
  const { user } = useAuth();
  const [children, setChildren] = useState<Student[]>([]);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [childData, feedData] = await Promise.all([
          parentApi.getMyChildren(),
          parentApi.getMyFeed({ limit: 5 }),
        ]);
        setChildren(childData || []);
        setPosts((feedData as any)?.data || []);
      } catch { /* ignore */ }
      finally { setLoading(false); }
    };
    load();
  }, []);

  return (
    <div className="space-y-8 pb-8">
      {/* Hero Header Area */}
      <div className="flex flex-col gap-1.5 animate-in fade-in slide-in-from-bottom-2 duration-500">
        <div className="flex items-center gap-2 mb-1">
          <span className="bg-primary/15 text-primary px-3 py-1 rounded-full text-xs font-semibold tracking-wide uppercase">
            Parent Portal
          </span>
        </div>
        <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-foreground">
          Xin chào, {user?.full_name || user?.email?.split('@')[0]}
        </h1>
        <p className="text-muted-foreground text-base max-w-2xl mt-1">
          Hôm nay con bạn có hoạt động gì mới? Bạn hiện đang theo dõi hồ sơ của {children.length} bé.
        </p>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : (
        <>
          {/* Quick Actions / Bento Box */}
          <div className="grid gap-5 md:grid-cols-2 lg:grid-cols-2">
            <Link href="/parent/children" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-amber-500/30 relative overflow-hidden bg-gradient-to-br hover:from-card hover:to-amber-50/50 dark:hover:to-amber-950/20">
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="p-4 bg-amber-500/10 text-amber-600 rounded-2xl transition-transform group-hover:scale-110 duration-300">
                        <Baby className="h-7 w-7" />
                      </div>
                      <div>
                        <p className="text-xl font-bold text-foreground group-hover:text-amber-600 transition-colors">Hồ sơ Con em</p>
                        <p className="text-sm text-muted-foreground mt-0.5">Quản lý và xem bảng điểm, sức khỏe</p>
                      </div>
                    </div>
                    <ChevronRight className="h-5 w-5 text-muted-foreground opacity-30 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                  </div>
                </CardContent>
              </Card>
            </Link>

            <Link href="/parent/posts" className="group">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-primary/30 relative overflow-hidden bg-gradient-to-br hover:from-card hover:to-primary/10 dark:hover:to-primary/20">
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="p-4 bg-primary/10 text-primary rounded-2xl transition-transform group-hover:scale-110 duration-300">
                        <MessageSquare className="h-7 w-7" />
                      </div>
                      <div>
                        <p className="text-xl font-bold text-foreground group-hover:text-primary transition-colors">Bảng tin Lớp</p>
                        <p className="text-sm text-muted-foreground mt-0.5">Cập nhật thông báo từ giáo viên</p>
                      </div>
                    </div>
                    <ChevronRight className="h-5 w-5 text-muted-foreground opacity-30 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                  </div>
                </CardContent>
              </Card>
            </Link>
          </div>

          <div className="grid gap-8 lg:grid-cols-3">
            {/* Children List */}
            <div className="lg:col-span-1 space-y-4">
              <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
                Con em của bạn
              </h2>
              {children.length > 0 ? (
                <div className="grid gap-3">
                  {children.map((child) => (
                    <Card key={child.student_id} className="group hover:border-primary/50 transition-colors">
                      <CardContent className="p-4 flex items-center gap-4">
                        <div className="h-12 w-12 shrink-0 rounded-2xl bg-primary/10 text-primary flex items-center justify-center font-bold text-lg">
                          {child.full_name?.charAt(0)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-semibold text-foreground truncate">{child.full_name}</p>
                          <p className="text-xs text-muted-foreground mt-0.5 flex items-center gap-1.5">
                            <CalendarClock className="h-3.5 w-3.5" /> Ngày sinh: {child.dob}
                          </p>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              ) : (
                <Card className="border-dashed shadow-none">
                  <CardContent className="flex flex-col items-center justify-center py-10 text-center">
                    <div className="p-3 bg-muted rounded-full mb-3">
                      <Baby className="h-6 w-6 text-muted-foreground" />
                    </div>
                    <p className="font-medium text-foreground">Chưa có dữ liệu học sinh</p>
                    <p className="mt-1 text-xs text-muted-foreground">Vui lòng liên hệ nhà trường để liên kết tài khoản.</p>
                  </CardContent>
                </Card>
              )}
            </div>

            {/* Recent Feed - Timeline Style */}
            <div className="lg:col-span-2 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-bold tracking-tight text-foreground flex items-center gap-2">
                  Hoạt động mới nhất
                </h2>
                <Link href="/parent/posts" className="text-sm font-medium text-primary hover:underline">
                  Xem tất cả
                </Link>
              </div>

              {posts.length > 0 ? (
                <div className="space-y-4">
                  {posts.map((p) => {
                    const config = postTypeConfig[p.type] || { label: p.type, colorClass: "bg-muted text-muted-foreground" };
                    return (
                      <Card key={p.post_id} className="group hover:shadow-md transition-shadow">
                        <CardContent className="p-5">
                          <div className="flex items-center justify-between mb-3 border-b border-border/50 pb-3">
                            <span className={`px-2.5 py-1 rounded-md text-xs font-semibold tracking-wide ${config.colorClass}`}>
                              {config.label}
                            </span>
                            <span className="text-xs font-medium text-muted-foreground flex items-center gap-1.5">
                              <CalendarClock className="h-3.5 w-3.5" />
                              {new Date(p.created_at).toLocaleString("vi-VN")}
                            </span>
                          </div>
                          <p className="text-sm text-foreground leading-relaxed whitespace-pre-line">
                            {p.content}
                          </p>
                        </CardContent>
                      </Card>
                    );
                  })}
                </div>
              ) : (
                <Card className="border-dashed shadow-none">
                  <CardContent className="flex flex-col items-center justify-center py-16 text-center">
                    <div className="p-4 bg-muted rounded-full mb-4">
                      <MessageSquare className="h-8 w-8 text-muted-foreground" />
                    </div>
                    <p className="font-semibold text-lg text-foreground">Không có bài viết mới</p>
                    <p className="mt-1 text-sm text-muted-foreground max-w-sm">Chưa có hoạt động hay thông báo nào được đăng tải lúc này.</p>
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