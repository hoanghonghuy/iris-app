/**
 * Admin Students Page
 * Quản lý học sinh theo lớp: chọn trường → chọn lớp → xem danh sách + tạo mới.
 * API: GET /admin/students/by-class/:class_id, POST /admin/students
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { School, Class, Student, CreateStudentRequest } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Users,
  Plus,
  X,
  Loader2,
  ChevronDown,
  Calendar,
  User,
} from "lucide-react";

export default function AdminStudentsPage() {
  // ─── State ────────────────────────────────────────────────────────

  // Schools + Classes (cascading selectors)
  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [loadingSchools, setLoadingSchools] = useState(true);
  const [loadingClasses, setLoadingClasses] = useState(false);

  // Students
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  // Form
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ full_name: "", dob: "", gender: "male" });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  // ─── Fetch schools ────────────────────────────────────────────────

  useEffect(() => {
    const load = async () => {
      try {
        const data = await adminApi.getSchools();
        setSchools(data || []);
        if (data && data.length > 0) {
          setSelectedSchoolId(data[0].school_id);
        }
      } catch {
        setError("Không thể tải danh sách trường");
      } finally {
        setLoadingSchools(false);
      }
    };
    load();
  }, []);

  // ─── Fetch classes khi đổi trường ────────────────────────────────

  useEffect(() => {
    if (!selectedSchoolId) return;
    const load = async () => {
      try {
        setLoadingClasses(true);
        setSelectedClassId("");
        setStudents([]);
        const data = await adminApi.getClassesBySchool(selectedSchoolId);
        setClasses(data || []);
        if (data && data.length > 0) {
          setSelectedClassId(data[0].class_id);
        }
      } catch {
        setClasses([]);
      } finally {
        setLoadingClasses(false);
      }
    };
    load();
  }, [selectedSchoolId]);

  // ─── Fetch students khi đổi lớp ──────────────────────────────────

  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) return;
    try {
      setLoadingStudents(true);
      setError("");
      const data = await adminApi.getStudentsByClass(selectedClassId);
      setStudents(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách học sinh");
    } finally {
      setLoadingStudents(false);
    }
  }, [selectedClassId]);

  useEffect(() => {
    fetchStudents();
  }, [fetchStudents]);

  // ─── Create student ───────────────────────────────────────────────

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.full_name.trim()) {
      setFormError("Họ tên không được để trống");
      return;
    }

    try {
      setSubmitting(true);
      setFormError("");
      await adminApi.createStudent({
        school_id: selectedSchoolId,
        class_id: selectedClassId,
        full_name: formData.full_name,
        dob: formData.dob,
        gender: formData.gender as "male" | "female" | "other",
      });
      setFormData({ full_name: "", dob: "", gender: "male" });
      setShowForm(false);
      fetchStudents();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo học sinh");
    } finally {
      setSubmitting(false);
    }
  };

  // ─── Helpers ──────────────────────────────────────────────────────

  const genderLabel: Record<string, string> = {
    male: "Nam",
    female: "Nữ",
    other: "Khác",
  };

  const selectedClassName = classes.find((c) => c.class_id === selectedClassId)?.name || "";

  // ─── Render ───────────────────────────────────────────────────────

  if (loadingSchools) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <Users className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Quản lý Học sinh</h1>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          {/* School selector */}
          <div className="relative">
            <select
              value={selectedSchoolId}
              onChange={(e) => setSelectedSchoolId(e.target.value)}
              className="h-9 appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            >
              {schools.map((s) => (
                <option key={s.school_id} value={s.school_id}>{s.name}</option>
              ))}
            </select>
            <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          </div>

          {/* Class selector */}
          {classes.length > 0 && (
            <div className="relative">
              <select
                value={selectedClassId}
                onChange={(e) => setSelectedClassId(e.target.value)}
                className="h-9 appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              >
                {classes.map((c) => (
                  <option key={c.class_id} value={c.class_id}>{c.name}</option>
                ))}
              </select>
              <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            </div>
          )}

          {/* Add button */}
          {selectedClassId && (
            <Button size="sm" onClick={() => setShowForm(!showForm)}>
              {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
              {showForm ? "Hủy" : "Thêm HS"}
            </Button>
          )}
        </div>
      </div>

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">
              Thêm học sinh — {selectedClassName}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && (
                <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                  {formError}
                </div>
              )}
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <label htmlFor="fullName" className="text-sm font-medium">
                    Họ tên <span className="text-destructive">*</span>
                  </label>
                  <Input
                    id="fullName"
                    placeholder="VD: Bé An"
                    value={formData.full_name}
                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="dob" className="text-sm font-medium">
                    Ngày sinh <span className="text-destructive">*</span>
                  </label>
                  <Input
                    id="dob"
                    type="date"
                    value={formData.dob}
                    onChange={(e) => setFormData({ ...formData, dob: e.target.value })}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="gender" className="text-sm font-medium">
                    Giới tính
                  </label>
                  <div className="relative">
                    <select
                      id="gender"
                      value={formData.gender}
                      onChange={(e) => setFormData({ ...formData, gender: e.target.value })}
                      className="h-9 w-full appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                    >
                      <option value="male">Nam</option>
                      <option value="female">Nữ</option>
                      <option value="other">Khác</option>
                    </select>
                    <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  </div>
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Tạo học sinh
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Error */}
      {error && (
        <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
      )}

      {/* Loading */}
      {(loadingClasses || loadingStudents) && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {/* Empty */}
      {!loadingStudents && !error && students.length === 0 && selectedClassId && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Users className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">
              Chưa có học sinh nào trong {selectedClassName}
            </p>
            <Button variant="outline" className="mt-4" onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Thêm học sinh đầu tiên
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Desktop Table */}
      {!loadingStudents && students.length > 0 && (
        <div className="hidden md:block">
          <Card>
            <CardContent className="p-0">
              <table className="w-full">
                <thead>
                  <tr className="border-b text-left text-sm text-muted-foreground">
                    <th className="px-6 py-3 font-medium">Họ tên</th>
                    <th className="px-6 py-3 font-medium">Ngày sinh</th>
                    <th className="px-6 py-3 font-medium">Giới tính</th>
                  </tr>
                </thead>
                <tbody>
                  {students.map((s) => (
                    <tr key={s.student_id} className="border-b last:border-0 hover:bg-zinc-50">
                      <td className="px-6 py-4 font-medium">{s.full_name}</td>
                      <td className="px-6 py-4 text-muted-foreground">{s.dob}</td>
                      <td className="px-6 py-4 text-muted-foreground">
                        {genderLabel[s.gender] || s.gender}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loadingStudents && students.length > 0 && (
        <div className="space-y-3 md:hidden">
          {students.map((s) => (
            <Card key={s.student_id}>
              <CardContent className="flex items-start gap-3 py-4">
                <User className="mt-0.5 h-5 w-5 shrink-0 text-muted-foreground" />
                <div className="min-w-0 flex-1">
                  <p className="font-medium">{s.full_name}</p>
                  <div className="mt-1 flex flex-wrap gap-3 text-sm text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {s.dob}
                    </span>
                    <span>{genderLabel[s.gender] || s.gender}</span>
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
