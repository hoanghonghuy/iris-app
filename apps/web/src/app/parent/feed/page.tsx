/**
 * Parent Feed Page
 * Bảng tin tổng hợp của tất cả con.
 * API: GET /parent/feed
 */
"use client";

import React, { useCallback, useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Pagination, Post } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { PostCard } from "@/components/shared/PostCard";
import { EmptyState } from "@/components/shared/EmptyState";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AlertCircle, Loader2, MessageSquare } from "lucide-react";

function extractErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as { response?: { data?: { error?: string } } }).response;
    return response?.data?.error || fallback;
  }

  return fallback;
}

export default function ParentFeedPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const fetchFeed = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const response = await parentApi.getMyFeed({ limit: 20, offset: currentOffset });
      setPosts(response.data || []);
      if (response.pagination) {
        setPagination(response.pagination);
      }
    } catch (err: unknown) {
      setError(extractErrorMessage(err, "Không thể tải bảng tin"));
    } finally {
      setLoading(false);
    }
  }, [currentOffset]);

  useEffect(() => {
    fetchFeed();
  }, [fetchFeed]);

  const patchPostById = (postId: string, patch: Partial<Post>) => {
    setPosts((prev) => prev.map((item) => (item.post_id === postId ? { ...item, ...patch } : item)));
  };

  return (
    <div className="mx-auto w-full max-w-3xl space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Bảng tin tổng hợp</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">Cập nhật từ giáo viên liên quan đến con của bạn.</p>
          <p className="mt-1 text-xs text-muted-foreground/80">Tổng số bài: {pagination.total}</p>
        </CardContent>
      </Card>

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

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
              enableShare={false}
              onPostPatched={patchPostById}
            />
          ))}
        </div>
      )}

      {!loading && posts.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}
    </div>
  );
}
