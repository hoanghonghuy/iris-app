import { FormEvent, useCallback, useEffect, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, CreatePostRequest, Pagination, Post, PostType, Student } from "@/types";
import { POST_TYPE_OPTIONS } from "@/lib/post-config";

export type ComposerScope = "class" | "student";

function extractErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as { response?: { data?: { error?: string } } }).response;
    return response?.data?.error || fallback;
  }

  return fallback;
}

export function isComposerScope(value: string): value is ComposerScope {
  return value === "class" || value === "student";
}

export function isPostType(value: string): value is PostType {
  return POST_TYPE_OPTIONS.some((option) => option.value === value);
}

export function useTeacherPostsPage() {
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
    const loadClasses = async () => {
      try {
        const classData = await teacherApi.getMyClasses();
        setClasses(classData || []);
        if (classData && classData.length > 0) {
          setSelectedClassId(classData[0].class_id);
        }
      } catch {
        setError("Không thể tải lớp");
      } finally {
        setLoading(false);
      }
    };

    void loadClasses();
  }, []);

  useEffect(() => {
    if (!selectedClassId) {
      return;
    }

    setCurrentOffset(0);

    const loadStudents = async () => {
      try {
        const studentData = await teacherApi.getStudentsInClass(selectedClassId);
        setStudents(studentData || []);
        if (studentData && studentData.length > 0) {
          setFormStudentId(studentData[0].student_id);
        }
      } catch {
      }
    };

    void loadStudents();
  }, [selectedClassId]);

  const fetchPosts = useCallback(async () => {
    if (!selectedClassId) {
      return;
    }

    try {
      setLoadingPosts(true);
      setError("");
      const response = await teacherApi.getClassPosts(selectedClassId, { limit: 20, offset: currentOffset });
      setPosts(response.data || []);
      if (response.pagination) {
        setPagination(response.pagination);
      }
    } catch (errorValue: unknown) {
      setError(extractErrorMessage(errorValue, "Không thể tải bài đăng"));
    } finally {
      setLoadingPosts(false);
    }
  }, [currentOffset, selectedClassId]);

  useEffect(() => {
    void fetchPosts();
  }, [fetchPosts]);

  const patchPostById = useCallback((postId: string, patch: Partial<Post>) => {
    setPosts((prev) => prev.map((item) => (item.post_id === postId ? { ...item, ...patch } : item)));
  }, []);

  const handleCreatePost = useCallback(async (event: FormEvent) => {
    event.preventDefault();
    if (!content.trim()) {
      setFormError("Nội dung không được trống");
      return;
    }

    const payload: CreatePostRequest = {
      scope_type: scopeType,
      type: postType,
      content,
      class_id: scopeType === "class" ? selectedClassId : undefined,
      student_id: scopeType === "student" ? formStudentId : undefined,
    };

    try {
      setSubmitting(true);
      setFormError("");
      await teacherApi.createPost(payload);
      setContent("");
      setShowForm(false);
      setCurrentOffset(0);
      await fetchPosts();
    } catch (errorValue: unknown) {
      setFormError(extractErrorMessage(errorValue, "Lỗi tạo bài đăng"));
    } finally {
      setSubmitting(false);
    }
  }, [content, fetchPosts, formStudentId, postType, scopeType, selectedClassId]);

  return {
    classes,
    selectedClassId,
    students,
    posts,
    loading,
    loadingPosts,
    error,
    pagination,
    currentOffset,
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
  };
}