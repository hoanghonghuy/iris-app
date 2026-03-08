/**
 * Teacher Health Page
 * Ghi nhận sức khỏe HS: chọn lớp → HS → ghi nhận nhiệt độ, triệu chứng, mức độ.
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/health
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, Student } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Heart, Loader2, ChevronDown, Plus, X } from "lucide-react";

const severityOptions = [
  { value: "normal", label: "Bình thường", color: "bg-green-100 text-green-700" },
  { value: "watch", label: "Theo dõi", color: "bg-yellow-100 text-yellow-700" },
  { value: "urgent", label: "Khẩn cấp", color: "bg-red-100 text-red-700" },
];

export default function TeacherHealthPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  // Form
  const [showForm, setShowForm] = useState(false);
  const [formStudentId, setFormStudentId] = useState("");
  const [temperature, setTemperature] = useState("");
  const [symptoms, setSymptoms] = useState("");
  const [severity, setSeverity] = useState("normal");
  const [note, setNote] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const data = await teacherApi.getMyClasses();
        setClasses(data || []);
        if (data && data.length > 0) setSelectedClassId(data[0].class_id);
      } catch { setError("Không thể tải lớp"); }
      finally { setLoadingClasses(false); }
    };
    load();
  }, []);

  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) return;
    try {
      setLoadingStudents(true);
      setError("");
      const data = await teacherApi.getStudentsInClass(selectedClassId);
      setStudents(data || []);
      if (data && data.length > 0) setFormStudentId(data[0].student_id);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải HS");
    } finally { setLoadingStudents(false); }
  }, [selectedClassId]);

  useEffect(() => { fetchStudents(); }, [fetchStudents]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formStudentId) { setFormError("Chọn học sinh"); return; }
    try {
      setSubmitting(true);
      setFormError("");
      setSuccess("");
      await teacherApi.createHealthLog({
        student_id: formStudentId,
        temperature: temperature ? parseFloat(temperature) : undefined,
        symptoms: symptoms || undefined,
        severity: severity as any,
        note: note || undefined,
      });
      setSuccess("Đã ghi nhận sức khỏe thành công!");
      setTemperature("");
      setSymptoms("");
      setSeverity("normal");
      setNote("");
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Lỗi ghi nhận");
    } finally { setSubmitting(false); }
  };

  if (loadingClasses) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <Heart className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Sức khỏe Học sinh</h1>
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
            {showForm ? "Đóng" : "Ghi nhận"}
          </Button>
        </div>
      </div>

      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}
      {success && <div className="rounded-md bg-green-100 p-4 text-sm text-green-700">{success}</div>}

      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Ghi nhận sức khỏe</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {formError && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{formError}</div>}
              <div className="grid gap-4 sm:grid-cols-2">
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
                <div className="space-y-2">
                  <label className="text-sm font-medium">Nhiệt độ (°C)</label>
                  <Input type="number" step="0.1" placeholder="36.5" value={temperature} onChange={(e) => setTemperature(e.target.value)} />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Triệu chứng</label>
                <Input placeholder="VD: Ho nhẹ, sổ mũi..." value={symptoms} onChange={(e) => setSymptoms(e.target.value)} />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Mức độ</label>
                <div className="flex gap-2">
                  {severityOptions.map((opt) => (
                    <button key={opt.value} type="button" onClick={() => setSeverity(opt.value)}
                      className={`rounded-full px-3 py-1 text-xs font-medium transition-colors ${severity === opt.value ? opt.color + " ring-2 ring-offset-1 ring-zinc-400" : "bg-zinc-100 text-zinc-500"}`}>
                      {opt.label}
                    </button>
                  ))}
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Ghi chú</label>
                <Input placeholder="Ghi chú thêm..." value={note} onChange={(e) => setNote(e.target.value)} />
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Lưu
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loadingStudents && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingStudents && !showForm && students.length > 0 && (
        <Card>
          <CardContent className="py-4">
            <p className="text-sm text-muted-foreground">
              {students.length} học sinh trong lớp. Nhấn &ldquo;Ghi nhận&rdquo; để thêm nhật ký sức khỏe.
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
