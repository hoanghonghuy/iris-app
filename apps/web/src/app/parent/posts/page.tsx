/**
 * Parent Posts Page
 * Xem bảng tin (feed) từ giáo viên.
 * API: GET /parent/feed
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Post } from "@/types";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { EmptyState } from "@/components/shared/EmptyState";
import { PostCard } from "@/components/shared/PostCard";
import { AlertCircle, Loader2, MessageSquare } from "lucide-react";

function extractErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as { response?: { data?: { error?: string } } }).response;
    return response?.data?.error || fallback;
  }

  return fallback;
}

export default function ParentPostsPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const data = await parentApi.getMyFeed({ limit: 50 });
        setPosts(data.data || []);
      } catch (err: unknown) {
        setError(extractErrorMessage(err, "Không thể tải bảng tin"));
      } finally { setLoading(false); }
    };
    load();
  }, []);

  const patchPostById = (postId: string, patch: Partial<Post>) => {
    setPosts((prev) => prev.map((item) => (item.post_id === postId ? { ...item, ...patch } : item)));
  };

  return (
    <div className="mx-auto w-full max-w-3xl space-y-4">
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && posts.length === 0 && !error && (
        <EmptyState
          icon={MessageSquare}
          title="Chưa có bài đăng nào"
          description="Bảng tin sẽ hiển thị thông báo và cập nhật từ giáo viên."
        />
      )}

      {!loading && posts.length > 0 && (
        <div className="space-y-4">
          {posts.map((post) => (
            <PostCard
              key={post.post_id}
              post={post}
              authorLabel="Giáo viên"
              audience="parent"
              onPostPatched={patchPostById}
            />
          ))}
        </div>
      )}
    </div>
  );
}
