/**
 * Teacher Posts Page
 * Tạo & xem bài đăng (thông báo lớp, nhận xét HS).
 * API: POST /teacher/posts, GET /teacher/classes/:id/posts
 */
"use client";

import React from "react";
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
import { isComposerScope, isPostType, useTeacherPostsPage } from "./useTeacherPostsPage";

export default function TeacherPostsPage() {
  const {
    classes,
    selectedClassId,
    students,
    posts,
    loading,
    loadingPosts,
    error,
    pagination,
    showForm,
    scopeType,
    formStudentId,
    postType,
    content,
    submitting,
    formError,
    setSelectedClassId,
    setCurrentOffset,
    setShowForm,
    setScopeType,
    setFormStudentId,
    setPostType,
    setContent,
    patchPostById,
    handleCreatePost,
  } = useTeacherPostsPage();

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
              enableShare={false}
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
