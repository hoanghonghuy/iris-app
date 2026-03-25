/**
 * Parent Posts Page
 * Xem bảng tin (feed) từ giáo viên.
 * API: GET /parent/feed
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Post } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { MessageSquare, Loader2 } from "lucide-react";

const postTypeLabels: Record<string, string> = {
  announcement: "Thông báo",
  activity: "Hoạt động",
  daily_note: "Nhận xét ngày",
  health_note: "Sức khỏe",
};

export default function ParentPostsPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const data = await parentApi.getMyFeed({ limit: 50 });
        setPosts((data as any)?.data || []);
      } catch (err: any) {
        setError(err.response?.data?.error || "Không thể tải bảng tin");
      } finally { setLoading(false); }
    };
    load();
  }, []);

  return (
    <div className="space-y-6">
      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && posts.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <MessageSquare className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có bài đăng nào</p>
        </CardContent></Card>
      )}

      {!loading && posts.length > 0 && (
        <div className="space-y-3">
          {posts.map((p) => (
            <Card key={p.post_id}>
              <CardContent className="py-4">
                <div className="flex items-center gap-2">
                  <span className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-medium">
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
      )}
    </div>
  );
}
