/**
 * Teacher Attendance Page
 * Chọn lớp → xem HS → điểm danh từng em.
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/attendance
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, Student } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  ClipboardCheck,
  Loader2,
  ChevronDown,
  Check,
} from "lucide-react";

const statusOptions = [
  { value: "present", label: "Có mặt", color: "bg-green-100 text-green-700" },
  { value: "absent", label: "Vắng", color: "bg-red-100 text-red-700" },
  { value: "late", label: "Muộn", color: "bg-yellow-100 text-yellow-700" },
  { value: "excused", label: "Có phép", color: "bg-blue-100 text-blue-700" },
];

export default function TeacherAttendancePage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  // Attendance state per student: { studentId: { status, note } }
  const [attendance, setAttendance] = useState<Record<string, { status: string; note: string }>>({});
  const [submitting, setSubmitting] = useState<string | null>(null);
  const [submitted, setSubmitted] = useState<Set<string>>(new Set());
  const [today] = useState(() => new Date().toISOString().slice(0, 10));

  useEffect(() => {
    const load = async () => {
      try {
        const data = await teacherApi.getMyClasses();
        setClasses(data || []);
        if (data && data.length > 0) setSelectedClassId(data[0].class_id);
      } catch {
        setError("Không thể tải lớp");
      } finally {
        setLoadingClasses(false);
      }
    };
    load();
  }, []);

  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) return;
    try {
      setLoadingStudents(true);
      setError("");
      setSubmitted(new Set());
      const data = await teacherApi.getStudentsInClass(selectedClassId);
      setStudents(data || []);
      const init: Record<string, { status: string; note: string }> = {};
      (data || []).forEach((s: Student) => {
        init[s.student_id] = { status: "present", note: "" };
      });
      setAttendance(init);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải HS");
    } finally {
      setLoadingStudents(false);
    }
  }, [selectedClassId]);

  useEffect(() => { fetchStudents(); }, [fetchStudents]);

  const handleMark = async (studentId: string) => {
    const att = attendance[studentId];
    if (!att) return;
    try {
      setSubmitting(studentId);
      await teacherApi.markAttendance({
        student_id: studentId,
        date: today,
        status: att.status as any,
        note: att.note,
      });
      setSubmitted((prev) => new Set(prev).add(studentId));
    } catch (err: any) {
      setError(err.response?.data?.error || "Lỗi điểm danh");
    } finally {
      setSubmitting(null);
    }
  };

  if (loadingClasses) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <ClipboardCheck className="h-7 w-7" />
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Điểm danh</h1>
            <p className="text-sm text-muted-foreground">Ngày: {today}</p>
          </div>
        </div>
        {classes.length > 0 && (
          <div className="relative">
            <select value={selectedClassId} onChange={(e) => setSelectedClassId(e.target.value)}
              className="h-9 appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
              {classes.map((c) => (
                <option key={c.class_id} value={c.class_id}>{c.name}</option>
              ))}
            </select>
            <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          </div>
        )}
      </div>

      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {loadingStudents && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingStudents && students.length === 0 && selectedClassId && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <ClipboardCheck className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Không có học sinh</p>
        </CardContent></Card>
      )}

      {!loadingStudents && students.length > 0 && (
        <div className="space-y-3">
          {students.map((s) => {
            const att = attendance[s.student_id] || { status: "present", note: "" };
            const isDone = submitted.has(s.student_id);
            return (
              <Card key={s.student_id} className={isDone ? "border-green-200 bg-green-50/50" : ""}>
                <CardContent className="py-4">
                  <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                    <div className="min-w-0">
                      <p className="font-medium">{s.full_name}</p>
                      <p className="text-sm text-muted-foreground">{s.dob}</p>
                    </div>
                    <div className="flex flex-wrap items-center gap-2">
                      {statusOptions.map((opt) => (
                        <button
                          key={opt.value}
                          onClick={() => setAttendance((prev) => ({ ...prev, [s.student_id]: { ...att, status: opt.value } }))}
                          className={`rounded-full px-3 py-1 text-xs font-medium transition-colors ${
                            att.status === opt.value ? opt.color + " ring-2 ring-offset-1 ring-zinc-400" : "bg-zinc-100 text-zinc-500 hover:bg-zinc-200"
                          }`}
                        >
                          {opt.label}
                        </button>
                      ))}
                    </div>
                  </div>
                  <div className="mt-3 flex items-center gap-2">
                    <Input
                      placeholder="Ghi chú..."
                      value={att.note}
                      onChange={(e) => setAttendance((prev) => ({ ...prev, [s.student_id]: { ...att, note: e.target.value } }))}
                      className="text-sm"
                    />
                    <Button
                      size="sm"
                      onClick={() => handleMark(s.student_id)}
                      disabled={submitting === s.student_id || isDone}
                      variant={isDone ? "outline" : "default"}
                    >
                      {submitting === s.student_id ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : isDone ? (
                        <Check className="h-4 w-4" />
                      ) : (
                        "Lưu"
                      )}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
