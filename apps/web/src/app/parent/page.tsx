/**
 * Parent Dashboard
 * Trang tổng quan cho phụ huynh: xem con em + feed mới nhất.
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Student, Post } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Users, MessageSquare, Loader2, User } from "lucide-react";
import Link from "next/link";

const postTypeLabels: Record<string, string> = {
  announcement: "Thông báo",
  activity: "Hoạt động",
  daily_note: "Nhận xét ngày",
  health_note: "Sức khỏe",
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

  if (loading) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold tracking-tight">Xin chào, Phụ huynh</h1>

      <Card>
        <CardContent className="py-4">
          <p className="font-medium">{user?.email}</p>
          <p className="mt-1 text-sm text-muted-foreground">
            Bạn có {children.length} con em đang theo dõi
          </p>
        </CardContent>
      </Card>

      {/* Quick links */}
      <div className="grid gap-4 sm:grid-cols-2">
        <Link href="/parent/children">
          <Card className="transition-colors hover:bg-zinc-50">
            <CardHeader className="pb-2">
              <Users className="h-6 w-6 text-muted-foreground" />
              <CardTitle className="text-lg">Con em</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Xem thông tin con em của bạn</CardDescription>
            </CardContent>
          </Card>
        </Link>

        <Link href="/parent/posts">
          <Card className="transition-colors hover:bg-zinc-50">
            <CardHeader className="pb-2">
              <MessageSquare className="h-6 w-6 text-muted-foreground" />
              <CardTitle className="text-lg">Bảng tin</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Thông báo và nhận xét từ giáo viên</CardDescription>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Children cards */}
      {children.length > 0 && (
        <div>
          <h2 className="mb-3 text-lg font-semibold">Con em của bạn</h2>
          <div className="grid gap-3 sm:grid-cols-2">
            {children.map((child) => (
              <Card key={child.student_id}>
                <CardContent className="flex items-center gap-3 py-4">
                  <User className="h-8 w-8 shrink-0 rounded-full bg-zinc-100 p-1.5 text-muted-foreground" />
                  <div>
                    <p className="font-medium">{child.full_name}</p>
                    <p className="text-sm text-muted-foreground">{child.dob}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Recent feed */}
      {posts.length > 0 && (
        <div>
          <h2 className="mb-3 text-lg font-semibold">Bài đăng gần đây</h2>
          <div className="space-y-3">
            {posts.map((p) => (
              <Card key={p.post_id}>
                <CardContent className="py-4">
                  <div className="flex items-center gap-2">
                    <span className="rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs font-medium">
                      {postTypeLabels[p.type] || p.type}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {new Date(p.created_at).toLocaleString("vi-VN")}
                    </span>
                  </div>
                  <p className="mt-2 text-sm whitespace-pre-line">{p.content}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}