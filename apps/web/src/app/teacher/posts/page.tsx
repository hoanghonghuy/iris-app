/**
 * Teacher Posts Page
 * Tạo & xem bài đăng (thông báo lớp, nhận xét HS).
 * API: POST /teacher/posts, GET /teacher/classes/:id/posts
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, CreatePostRequest, Pagination, Post, PostType, Student } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { EmptyState } from "@/components/shared/EmptyState";
import { PostCard } from "@/components/shared/PostCard";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { AlertCircle, Loader2, MessageSquare, Plus, X } from "lucide-react";
import { POST_SCOPE_LABELS, POST_TYPE_OPTIONS } from "@/lib/post-config";

type ComposerScope = "class" | "student";

function extractErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as { response?: { data?: { error?: string } } }).response;
    return response?.data?.error || fallback;
  }

  return fallback;
}

function isComposerScope(value: string): value is ComposerScope {
  return value === "class" || value === "student";
}

function isPostType(value: string): value is PostType {
  return POST_TYPE_OPTIONS.some((option) => option.value === value);
}

export default function TeacherPostsPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const [showForm, setShowForm] = useState(false);
  const [scopeType, setScopeType] = useState<ComposerScope>("class");
  const [formStudentId, setFormStudentId] = useState("");
  const [postType, setPostType] = useState<PostType>("announcement");
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const data = await teacherApi.getMyClasses();
        setClasses(data || []);
        if (data && data.length > 0) setSelectedClassId(data[0].class_id);
      } catch { setError("Không thể tải lớp"); }
      finally { setLoading(false); }
    };
    load();
  }, []);

  useEffect(() => {
    if (!selectedClassId) return;
    setCurrentOffset(0);

    const load = async () => {
      try {
        const data = await teacherApi.getStudentsInClass(selectedClassId);
        setStudents(data || []);
        if (data && data.length > 0) setFormStudentId(data[0].student_id);
      } catch { /* ignore */ }
    };
    load();
  }, [selectedClassId]);

  const fetchPosts = useCallback(async () => {
    if (!selectedClassId) return;
    try {
      setLoadingPosts(true); setError("");
      const response = await teacherApi.getClassPosts(selectedClassId, { limit: 20, offset: currentOffset });
      setPosts(response.data || []);
      if (response.pagination) setPagination(response.pagination);
    } catch (err: unknown) {
      setError(extractErrorMessage(err, "Không thể tải bài đăng"));
    } finally { setLoadingPosts(false); }
  }, [selectedClassId, currentOffset]);

  useEffect(() => { fetchPosts(); }, [fetchPosts]);

  const patchPostById = useCallback((postId: string, patch: Partial<Post>) => {
    setPosts((prev) => prev.map((item) => (item.post_id === postId ? { ...item, ...patch } : item)));
  }, []);

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) { setFormError("Nội dung không được trống"); return; }

    const payload: CreatePostRequest = {
      scope_type: scopeType,
      type: postType,
      content,
      class_id: scopeType === "class" ? selectedClassId : undefined,
      student_id: scopeType === "student" ? formStudentId : undefined,
    };

    try {
      setSubmitting(true); setFormError("");
      await teacherApi.createPost(payload);
      setContent("");
      setShowForm(false);
      setCurrentOffset(0);
      fetchPosts();
    } catch (err: unknown) {
      setFormError(extractErrorMessage(err, "Lỗi tạo bài đăng"));
    } finally { setSubmitting(false); }
  };

  if (loading) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="mx-auto w-full max-w-3xl space-y-4">
      <Card>
        <CardContent className="space-y-3 p-4 sm:p-5">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex min-w-0 items-center gap-3">
              <Avatar>
                <AvatarFallback>GV</AvatarFallback>
              </Avatar>
              <Button
                type="button"
                variant="outline"
                className="h-10 w-full justify-start rounded-full px-4 text-muted-foreground sm:w-auto sm:min-w-72"
                onClick={() => setShowForm((prev) => !prev)}
              >
                {showForm ? "Đóng khung soạn bài" : "Bạn muốn chia sẻ điều gì với lớp hôm nay?"}
              </Button>
            </div>

            {classes.length > 0 && (
              <Select value={selectedClassId} onValueChange={setSelectedClassId}>
                <SelectTrigger className="w-full sm:w-[220px]">
                  <SelectValue placeholder="Chọn lớp" />
                </SelectTrigger>
                <SelectContent>
                  {classes.map((classInfo) => (
                    <SelectItem key={classInfo.class_id} value={classInfo.class_id}>
                      {classInfo.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>

          <div className="flex flex-wrap items-center gap-2 border-t pt-3">
            <Badge variant="outline">{pagination.total} bài đăng</Badge>
            <Badge variant="secondary">{POST_SCOPE_LABELS[scopeType]}</Badge>
            <Button size="sm" variant="ghost" onClick={() => setShowForm((prev) => !prev)}>
              {showForm ? <X className="mr-1 h-4 w-4" /> : <Plus className="mr-1 h-4 w-4" />}
              {showForm ? "Đóng" : "Tạo bài"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Tạo bài đăng mới</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreatePost} className="space-y-4">
              {formError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>}
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label>Phạm vi</Label>
                  <Select
                    value={scopeType}
                    onValueChange={(value) => {
                      if (isComposerScope(value)) {
                        setScopeType(value);
                      }
                    }}
                  >
                    <SelectTrigger className="w-full"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="class">Cả lớp</SelectItem>
                      <SelectItem value="student">Từng HS</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                {scopeType === "student" && (
                  <div className="space-y-2">
                    <Label>Học sinh</Label>
                    <Select value={formStudentId} onValueChange={setFormStudentId}>
                      <SelectTrigger className="w-full"><SelectValue placeholder="Chọn HS" /></SelectTrigger>
                      <SelectContent>
                        {students.map((s) => <SelectItem key={s.student_id} value={s.student_id}>{s.full_name}</SelectItem>)}
                      </SelectContent>
                    </Select>
                  </div>
                )}
                <div className="space-y-2">
                  <Label>Loại bài</Label>
                  <Select
                    value={postType}
                    onValueChange={(value) => {
                      if (isPostType(value)) {
                        setPostType(value);
                      }
                    }}
                  >
                    <SelectTrigger className="w-full"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {POST_TYPE_OPTIONS.map((typeOption) => (
                        <SelectItem key={typeOption.value} value={typeOption.value}>
                          {typeOption.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="postContent">Nội dung</Label>
                <Textarea id="postContent" value={content} onChange={(e) => setContent(e.target.value)}
                  placeholder="Nhập nội dung bài đăng..." rows={4} required />
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Đăng
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loadingPosts && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingPosts && posts.length === 0 && (
        <EmptyState
          icon={MessageSquare}
          title="Chưa có bài đăng nào"
          description="Hãy tạo bài đầu tiên để cập nhật thông tin cho lớp học."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Tạo bài đăng
            </Button>
          }
        />
      )}

      {!loadingPosts && posts.length > 0 && (
        <div className="space-y-4">
          {posts.map((post) => (
            <PostCard
              key={post.post_id}
              post={post}
              audience="teacher"
              onPostPatched={patchPostById}
            />
          ))}
        </div>
      )}

      {/* Pagination */}
      {!loadingPosts && posts.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}
    </div>
  );
}
