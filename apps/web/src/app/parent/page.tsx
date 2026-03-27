/**
 * Parent Dashboard
 * Trang tổng quan cho phụ huynh chuẩn Minimalist Pastel.
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Post } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardContent } from "@/components/ui/card";
import { MessageSquare, Loader2, Baby, CalendarClock } from "lucide-react";
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
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const feedData = await parentApi.getMyFeed({ limit: 5 });
        setPosts(feedData?.data || []);
      } catch { /* ignore */ }
      finally { setLoading(false); }
    };
    load();
  }, []);

  return (
    <div className="flex h-full min-h-0 flex-col gap-6">
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
          {/* Quick Actions / Bento Box */}
          <div className="grid gap-2.5 md:gap-3 grid-cols-2 shrink-0">
            <Link href="/parent/children" className="group h-full">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-chart-2/30 relative overflow-hidden bg-gradient-to-br hover:from-card hover:to-chart-2/10 dark:hover:to-chart-2/20">
                <CardContent className="p-2.5 md:p-3 flex flex-col items-center justify-center gap-1.5 md:gap-2 h-full text-center">
                  <div className="p-2 md:p-2.5 bg-chart-2/10 text-chart-2 rounded-xl md:rounded-2xl transition-transform group-hover:scale-110 duration-300">
                    <Baby className="h-5 w-5 md:h-6 md:w-6" />
                  </div>
                  <p className="font-medium md:font-semibold text-xs text-foreground group-hover:text-chart-2 transition-colors">Hồ sơ con</p>
                </CardContent>
              </Card>
            </Link>

            <Link href="/parent/posts" className="group h-full">
              <Card className="h-full transition-all duration-300 hover:shadow-md hover:border-primary/30 relative overflow-hidden bg-gradient-to-br hover:from-card hover:to-primary/10 dark:hover:to-primary/20">
                <CardContent className="p-2.5 md:p-3 flex flex-col items-center justify-center gap-1.5 md:gap-2 h-full text-center">
                  <div className="p-2 md:p-2.5 bg-primary/10 text-primary rounded-xl md:rounded-2xl transition-transform group-hover:scale-110 duration-300">
                    <MessageSquare className="h-5 w-5 md:h-6 md:w-6" />
                  </div>
                  <p className="font-medium md:font-semibold text-xs text-foreground group-hover:text-primary transition-colors">Bảng tin Lớp</p>
                </CardContent>
              </Card>
            </Link>
          </div>

          <div className="flex-1 min-h-0 overflow-hidden">
            {/* Recent Feed - Focused View */}
            <div className="flex h-full min-h-0 flex-col gap-4">
              <div className="flex items-center justify-between shrink-0">
                <h2 className="text-base font-bold tracking-tight text-foreground flex items-center gap-2">
                  Hoạt động mới nhất
                </h2>
                <Link href="/parent/posts" className="text-xs font-semibold text-primary hover:underline transition-all hover:translate-x-0.5">
                  Xem tất cả
                </Link>
              </div>

              <div className="custom-scrollbar flex-1 min-h-0 overflow-y-auto pr-2">
                {posts.length > 0 ? (
                  <div className="mx-auto max-w-4xl space-y-4 pb-4">
                    {posts.map((p) => {
                      const config = postTypeConfig[p.type] || { label: p.type, colorClass: "bg-muted text-muted-foreground" };
                      return (
                        <Card key={p.post_id} className="group hover:shadow-lg transition-all border-border/50 hover:border-primary/20 bg-card/50 backdrop-blur-sm">
                          <CardContent className="p-5 md:p-6">
                            <div className="flex items-center justify-between mb-4 border-b border-border/40 pb-4">
                              <span className={`px-2.5 py-0.5 rounded-full text-[10px] uppercase font-bold tracking-wider ${config.colorClass}`}>
                                {config.label}
                              </span>
                              <span className="text-[11px] font-medium text-muted-foreground/80 flex items-center gap-1.5">
                                <CalendarClock className="h-3.5 w-3.5" />
                                {new Date(p.created_at).toLocaleString("vi-VN")}
                              </span>
                            </div>
                            <p className="text-sm md:text-base text-foreground/90 leading-relaxed whitespace-pre-line font-medium italic">
                              &ldquo;{p.content}&rdquo;
                            </p>
                          </CardContent>
                        </Card>
                      );
                    })}
                  </div>
                ) : (
                  <Card className="border-dashed shadow-none">
                    <CardContent className="flex flex-col items-center justify-center py-20 text-center">
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
