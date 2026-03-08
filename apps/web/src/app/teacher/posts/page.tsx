/**
 * Teacher Posts Page
 * Tạo & xem bài đăng (thông báo lớp, nhận xét HS).
 * API: POST /teacher/posts, GET /teacher/classes/:id/posts
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, Post, Student } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { MessageSquare, Loader2, ChevronDown, Plus, X } from "lucide-react";

const postTypeLabels: Record<string, string> = {
  announcement: "Thông báo",
  activity: "Hoạt động",
  daily_note: "Nhận xét ngày",
  health_note: "Sức khỏe",
};

export default function TeacherPostsPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [error, setError] = useState("");

  // Form
  const [showForm, setShowForm] = useState(false);
  const [scopeType, setScopeType] = useState<"class" | "student">("class");
  const [formStudentId, setFormStudentId] = useState("");
  const [postType, setPostType] = useState("announcement");
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

  // Fetch students khi đổi lớp (cho scope student)
  useEffect(() => {
    if (!selectedClassId) return;
    const load = async () => {
      try {
        const data = await teacherApi.getStudentsInClass(selectedClassId);
        setStudents(data || []);
        if (data && data.length > 0) setFormStudentId(data[0].student_id);
      } catch { /* ignore */ }
    };
    load();
  }, [selectedClassId]);

  // Fetch posts khi đổi lớp
  const fetchPosts = useCallback(async () => {
    if (!selectedClassId) return;
    try {
      setLoadingPosts(true);
      setError("");
      const data = await teacherApi.getClassPosts(selectedClassId);
      setPosts(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải bài đăng");
    } finally { setLoadingPosts(false); }
  }, [selectedClassId]);

  useEffect(() => { fetchPosts(); }, [fetchPosts]);

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) { setFormError("Nội dung không được trống"); return; }
    try {
      setSubmitting(true);
      setFormError("");
      await teacherApi.createPost({
        scope_type: scopeType,
        class_id: scopeType === "class" ? selectedClassId : undefined,
        student_id: scopeType === "student" ? formStudentId : undefined,
        type: postType as any,
        content,
      });
      setContent("");
      setShowForm(false);
      fetchPosts();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Lỗi tạo bài đăng");
    } finally { setSubmitting(false); }
  };

  if (loading) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <MessageSquare className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Bài đăng</h1>
        </div>
        <div className="flex items-center gap-2">
          {classes.length > 0 && (
            <div className="relative">
              <select value={selectedClassId} onChange={(e) => setSelectedClassId(e.target.value)}
                className="h-9 appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                {classes.map((c) => (<option key={c.class_id} value={c.class_id}>{c.name}</option>))}
              </select>
              <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            </div>
          )}
          <Button size="sm" onClick={() => setShowForm(!showForm)}>
            {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
            {showForm ? "Hủy" : "Tạo bài"}
          </Button>
        </div>
      </div>

      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Tạo bài đăng mới</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreatePost} className="space-y-4">
              {formError && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{formError}</div>}
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Phạm vi</label>
                  <div className="relative">
                    <select value={scopeType} onChange={(e) => setScopeType(e.target.value as "class" | "student")}
                      className="h-9 w-full appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                      <option value="class">Cả lớp</option>
                      <option value="student">Từng HS</option>
                    </select>
                    <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  </div>
                </div>
                {scopeType === "student" && (
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Học sinh</label>
                    <div className="relative">
                      <select value={formStudentId} onChange={(e) => setFormStudentId(e.target.value)}
                        className="h-9 w-full appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                        {students.map((s) => (<option key={s.student_id} value={s.student_id}>{s.full_name}</option>))}
                      </select>
                      <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    </div>
                  </div>
                )}
                <div className="space-y-2">
                  <label className="text-sm font-medium">Loại bài</label>
                  <div className="relative">
                    <select value={postType} onChange={(e) => setPostType(e.target.value)}
                      className="h-9 w-full appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                      <option value="announcement">Thông báo</option>
                      <option value="activity">Hoạt động</option>
                      <option value="daily_note">Nhận xét ngày</option>
                      <option value="health_note">Sức khỏe</option>
                    </select>
                    <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  </div>
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Nội dung</label>
                <textarea
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  placeholder="Nhập nội dung bài đăng..."
                  rows={4}
                  className="w-full rounded-md border p-3 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                  required
                />
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Đăng
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loadingPosts && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingPosts && posts.length === 0 && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <MessageSquare className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có bài đăng nào cho lớp này</p>
        </CardContent></Card>
      )}

      {!loadingPosts && posts.length > 0 && (
        <div className="space-y-3">
          {posts.map((p) => (
            <Card key={p.post_id}>
              <CardContent className="py-4">
                <div className="flex items-start justify-between gap-2">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs font-medium">
                        {postTypeLabels[p.type] || p.type}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {new Date(p.created_at).toLocaleString("vi-VN")}
                      </span>
                    </div>
                    <p className="mt-2 text-sm whitespace-pre-line">{p.content}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
