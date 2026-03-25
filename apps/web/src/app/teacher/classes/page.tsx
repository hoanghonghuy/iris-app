/**
 * Teacher Classes Page
 * Giáo viên xem lớp được phân công → chọn lớp → xem danh sách học sinh.
 * API: GET /teacher/classes, GET /teacher/classes/:class_id/students
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, Student } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { formatDateVN } from "@/lib/utils";
import {
  GraduationCap,
  Loader2,
  ChevronDown,
  User,
  Calendar,
} from "lucide-react";

export default function TeacherClassesPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  const genderLabel: Record<string, string> = { male: "Nam", female: "Nữ", other: "Khác" };

  useEffect(() => {
    const load = async () => {
      try {
        const data = await teacherApi.getMyClasses();
        setClasses(data || []);
        if (data && data.length > 0) {
          setSelectedClassId(data[0].class_id);
        }
      } catch {
        setError("Không thể tải danh sách lớp");
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
      const data = await teacherApi.getStudentsInClass(selectedClassId);
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

  const selectedClassName = classes.find((c) => c.class_id === selectedClassId)?.name || "";

  if (loadingClasses) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <GraduationCap className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Lớp của tôi</h1>
        </div>

        {classes.length > 0 && (
          <div className="relative">
            <select
              value={selectedClassId}
              onChange={(e) => setSelectedClassId(e.target.value)}
              className="h-9 appearance-none rounded-md border bg-background py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            >
              {classes.map((c) => (
                <option key={c.class_id} value={c.class_id}>
                  {c.name} ({c.school_year})
                </option>
              ))}
            </select>
            <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          </div>
        )}
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
      )}

      {loadingStudents && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {!loadingStudents && students.length === 0 && selectedClassId && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <User className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">
              Chưa có học sinh nào trong {selectedClassName}
            </p>
          </CardContent>
        </Card>
      )}

      {classes.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <GraduationCap className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Bạn chưa được phân công lớp nào</p>
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
                    <tr key={s.student_id} className="border-b last:border-0 hover:bg-muted">
                      <td className="px-6 py-4 font-medium">{s.full_name}</td>
                      <td className="px-6 py-4 text-muted-foreground">{formatDateVN(s.dob)}</td>
                      <td className="px-6 py-4 text-muted-foreground">{genderLabel[s.gender] || s.gender}</td>
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
                <div>
                  <p className="font-medium">{s.full_name}</p>
                  <p className="mt-1 flex items-center gap-1 text-sm text-muted-foreground">
                    <Calendar className="h-3 w-3" /> {formatDateVN(s.dob)} · {genderLabel[s.gender] || s.gender}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
