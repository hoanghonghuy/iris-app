/**
 * Parent Dashboard
 * Trang tổng quan cho phụ huynh chuẩn Minimalist Pastel.
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { ParentAnalytics, Post } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardContent } from "@/components/ui/card";
import { MessageSquare, Loader2, Baby, CalendarClock, Heart, ClipboardCheck, ChevronRight } from "lucide-react";
import Link from "next/link";

const postTypeConfig: Record<string, { label: string, colorClass: string }> = {
  announcement: { label: "Thông báo", colorClass: "bg-primary/10 text-primary" },
  activity: { label: "Hoạt động", colorClass: "bg-chart-2/10 text-chart-2" },
  daily_note: { label: "Nhận xét", colorClass: "bg-chart-1/10 text-chart-1" },
  health_note: { label: "Sức khỏe", colorClass: "bg-destructive/10 text-destructive" },
};

export default function ParentDashboard() {
  const { user } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [stats, setStats] = useState<ParentAnalytics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [feedData, analytics] = await Promise.all([
          parentApi.getMyFeed({ limit: 5 }),
          parentApi.getAnalytics(),
        ]);
        setPosts(feedData?.data || []);
        setStats(analytics);
      } catch { /* ignore */ }
      finally { setLoading(false); }
    };
    load();
  }, []);

  return (
    <div className="flex h-full min-h-0 flex-col gap-6 pb-6">
      <div className="flex flex-col gap-1.5 animate-in fade-in slide-in-from-bottom-2 duration-500 shrink-0">
        <div className="hidden md:flex items-center gap-2 mb-1">
          <span className="bg-primary/15 text-primary px-3 py-1 rounded-full text-xs font-semibold tracking-wide uppercase">
            Parent Portal
          </span>
        </div>
        <h1 className="text-2xl md:text-3xl font-extrabold tracking-tight text-foreground">
          Xin chào, {user?.full_name || user?.email?.split('@')[0]}
        </h1>
        <p className="text-muted-foreground text-sm max-w-2xl mt-1 hidden md:block">
          Hôm nay con bạn có hoạt động gì mới?
        </p>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : (
        <>
          {/* Stat Cards Grid P0 Fix */}
          <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-4 shrink-0">
            <Card className="shadow-sm border-transparent transition-all hover:border-chart-2/30">
              <CardContent className="p-4 flex items-center gap-4">
                <div className="p-3 bg-chart-2/10 text-chart-2 rounded-2xl">
                  <Baby className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-[11px] font-semibold text-muted-foreground uppercase opacity-80 whitespace-nowrap">Tổng số con</p>
                  <p className="text-2xl font-bold">{stats?.total_children ?? 0}</p>
                </div>
              </CardContent>
            </Card>

            <Card className="shadow-sm border-transparent transition-all hover:border-primary/30">
              <CardContent className="p-4 flex items-center gap-4">
                <div className="p-3 bg-primary/10 text-primary rounded-2xl">
                  <CalendarClock className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-[11px] font-semibold text-muted-foreground uppercase opacity-80 whitespace-nowrap">Lịch sắp tới</p>
                  <p className="text-2xl font-bold">{stats?.upcoming_appointments ?? 0}</p>
                </div>
              </CardContent>
            </Card>

            <Card className="shadow-sm border-transparent transition-all hover:border-chart-1/30">
              <CardContent className="p-4 flex items-center gap-4">
                <div className="p-3 bg-chart-1/10 text-chart-1 rounded-2xl">
                  <MessageSquare className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-[11px] font-semibold text-muted-foreground uppercase opacity-80 whitespace-nowrap">Bài đăng (7 ngày)</p>
                  <p className="text-2xl font-bold">{stats?.recent_posts_7d ?? 0}</p>
                </div>
              </CardContent>
            </Card>

            <Card className="shadow-sm border-orange-500/20 bg-orange-500/5 transition-all hover:border-orange-500/50">
              <CardContent className="p-4 flex items-center gap-4">
                <div className="p-3 bg-orange-500/20 text-orange-600 dark:text-orange-400 rounded-2xl">
                  <Heart className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-[11px] font-semibold text-orange-600 dark:text-orange-400 uppercase opacity-80 whitespace-nowrap">Cảnh báo sức khỏe</p>
                  <p className="text-2xl font-bold text-orange-600 dark:text-orange-400">{stats?.recent_health_alerts_7d ?? 0}</p>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Today's Attendance Strip (P0 Mockup) */}
          <div className="shrink-0 gap-4 grid lg:grid-cols-3">
             <div className="lg:col-span-2">
                <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2 mb-3">
                  Tình trạng hôm nay
                </h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <Card className="bg-card/60 backdrop-blur-sm border-transparent shadow-sm hover:shadow-md transition-shadow cursor-default">
                     <CardContent className="p-3.5 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="h-10 w-10 rounded-full bg-success/20 flex items-center justify-center">
                             <Baby className="h-5 w-5 text-success" />
                          </div>
                          <div>
                             <p className="font-semibold text-sm">Gia An</p>
                             <p className="text-xs text-muted-foreground mt-0.5">Lớp Lá Non</p>
                          </div>
                        </div>
                        <div className="flex flex-col items-end">
                           <span className="inline-flex items-center rounded-full bg-success/10 px-2.5 py-0.5 text-[11px] font-semibold text-success">
                              <ClipboardCheck className="mr-1 h-3 w-3" /> Có mặt
                           </span>
                           <span className="text-[10px] text-muted-foreground mt-1">Lúc 08:30</span>
                        </div>
                     </CardContent>
                  </Card>

                  <Card className="bg-card/60 backdrop-blur-sm border-transparent shadow-sm hover:shadow-md transition-shadow cursor-default relative overflow-hidden">
                     {/* Red accent border top for absent */}
                     <div className="absolute top-0 left-0 right-0 h-1 bg-destructive/60"></div>
                     <CardContent className="p-3.5 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center">
                             <Baby className="h-5 w-5 text-muted-foreground" />
                          </div>
                          <div>
                             <p className="font-semibold text-sm text-foreground">Khánh Băng</p>
                             <p className="text-xs text-muted-foreground mt-0.5">Lớp Chồi 1</p>
                          </div>
                        </div>
                        <div className="flex flex-col items-end">
                           <span className="inline-flex items-center rounded-full bg-destructive/10 px-2.5 py-0.5 text-[11px] font-semibold text-destructive">
                              Vắng mặt
                           </span>
                        </div>
                     </CardContent>
                  </Card>
                </div>
             </div>

             {/* Quick Actions Restyled */}
             <div className="lg:col-span-1">
               <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2 mb-3">
                  Tác vụ nhanh
                </h2>
               <div className="grid gap-3 grid-cols-1 shrink-0 h-full">
                  <Link href="/parent/children" className="block w-full">
                    <Card className="group hover:shadow-md transition-all cursor-pointer border-transparent hover:border-chart-2/30 shadow-sm bg-card/60">
                      <CardContent className="p-3.5 flex items-center justify-between gap-3">
                        <div className="flex items-center gap-3">
                          <div className="p-2.5 bg-chart-2/10 text-chart-2 rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                            <Baby className="h-5 w-5" />
                          </div>
                          <div className="text-left">
                            <p className="font-semibold text-sm text-foreground group-hover:text-chart-2 transition-colors">Hồ sơ con</p>
                            <p className="text-xs font-medium text-muted-foreground mt-0.5">Xem chi tiết</p>
                          </div>
                        </div>
                        <ChevronRight className="h-4 w-4 text-muted-foreground opacity-20 group-hover:opacity-100 transition-opacity group-hover:text-chart-2 group-hover:translate-x-1" />
                      </CardContent>
                    </Card>
                  </Link>

                  <Link href="/parent/posts" className="block w-full">
                    <Card className="group hover:shadow-md transition-all cursor-pointer border-transparent hover:border-primary/30 shadow-sm bg-card/60">
                      <CardContent className="p-3.5 flex items-center justify-between gap-3">
                        <div className="flex items-center gap-3">
                          <div className="p-2.5 bg-primary/10 text-primary rounded-xl shrink-0 transition-transform group-hover:scale-110 duration-300">
                            <MessageSquare className="h-5 w-5" />
                          </div>
                          <div className="text-left">
                            <p className="font-semibold text-sm text-foreground group-hover:text-primary transition-colors">Bảng tin Lớp</p>
                            <p className="text-xs font-medium text-muted-foreground mt-0.5">Hoạt động chung</p>
                          </div>
                        </div>
                        <ChevronRight className="h-4 w-4 text-muted-foreground opacity-20 group-hover:opacity-100 transition-opacity group-hover:text-primary group-hover:translate-x-1" />
                      </CardContent>
                    </Card>
                  </Link>
                </div>
             </div>
          </div>

          <div className="flex-1 min-h-0 overflow-hidden mt-4">
            {/* Recent Feed - Focused View */}
            <div className="flex h-full min-h-0 flex-col gap-4">
              <div className="flex items-center justify-between shrink-0">
                <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2">
                  Bảng tin mới nhất
                </h2>
                <Link href="/parent/posts" className="text-xs font-semibold text-primary hover:underline transition-all hover:translate-x-0.5">
                  Xem tất cả
                </Link>
              </div>

              <div className="custom-scrollbar flex-1 min-h-0 overflow-y-auto pr-2">
                {posts.length > 0 ? (
                  <div className="grid gap-4 lg:grid-cols-2 pb-4">
                    {posts.map((p, index) => {
                      const config = postTypeConfig[p.type] || { label: p.type, colorClass: "bg-muted text-muted-foreground" };
                      const mockImages = [
                        "https://images.unsplash.com/photo-1502086223501-7ea6ecd79368?q=80&w=500&auto=format&fit=crop",
                        "https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?q=80&w=500&auto=format&fit=crop",
                        "https://images.unsplash.com/photo-1516627145497-ae6968895b74?q=80&w=500&auto=format&fit=crop"
                      ];

                      return (
                        <Card key={p.post_id} className="group hover:shadow-lg transition-all border-border/50 hover:border-primary/20 bg-card/50 backdrop-blur-sm h-fit overflow-hidden">
                          {((p.type === 'activity' || p.type === 'announcement') && index % 2 === 0) && (
                            <div className="w-full h-40 overflow-hidden bg-muted relative">
                              <img 
                                src={mockImages[index % mockImages.length]} 
                                alt="Activity Photo" 
                                className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                              />
                              <div className="absolute top-3 left-3">
                                <span className={`px-2.5 py-0.5 rounded-full text-[10px] uppercase font-bold tracking-wider bg-black/60 text-white backdrop-blur-md`}>
                                  {config.label}
                                </span>
                              </div>
                            </div>
                          )}
                          <CardContent className={`p-5 md:p-6 ${((p.type === 'activity' || p.type === 'announcement') && index % 2 === 0) ? 'pt-4' : ''}`}>
                            {!((p.type === 'activity' || p.type === 'announcement') && index % 2 === 0) && (
                              <div className="flex items-center justify-between mb-4 border-b border-border/40 pb-4">
                                <span className={`px-2.5 py-0.5 rounded-full text-[10px] uppercase font-bold tracking-wider ${config.colorClass}`}>
                                  {config.label}
                                </span>
                                <span className="text-[11px] font-medium text-muted-foreground/80 flex items-center gap-1.5">
                                  <CalendarClock className="h-3.5 w-3.5" />
                                  {new Date(p.created_at).toLocaleString("vi-VN", { hour: '2-digit', minute: '2-digit', day: '2-digit', month: '2-digit' })}
                                </span>
                              </div>
                            )}
                            {((p.type === 'activity' || p.type === 'announcement') && index % 2 === 0) && (
                               <div className="flex items-center justify-between mb-2">
                                <span className="text-[11px] font-medium text-muted-foreground/80 flex items-center gap-1.5">
                                  <CalendarClock className="h-3.5 w-3.5" />
                                  {new Date(p.created_at).toLocaleString("vi-VN", { hour: '2-digit', minute: '2-digit', day: '2-digit', month: '2-digit' })}
                                </span>
                               </div>
                            )}
                            <p className="text-sm md:text-base text-foreground/90 leading-relaxed whitespace-pre-line font-medium italic line-clamp-3">
                              &ldquo;{p.content}&rdquo;
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
          </div>
        </>
      )}
    </div>
  );
}
