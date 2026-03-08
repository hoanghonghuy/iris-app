/**
 * Admin Teachers Page
 * Danh sách giáo viên + gán/hủy gán lớp.
 * API: GET /admin/teachers, POST/DELETE /admin/teachers/:id/classes/:class_id
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Teacher, School, Class } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { BookUser, Loader2, Phone, Mail, Link2, Unlink, AlertCircle, CheckCircle2 } from "lucide-react";

export default function AdminTeachersPage() {
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [assigningTeacherId, setAssigningTeacherId] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [success, setSuccess] = useState("");

  const fetchTeachers = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const data = await adminApi.getTeachers();
      setTeachers(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách giáo viên");
    } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchTeachers(); }, [fetchTeachers]);

  useEffect(() => {
    const load = async () => {
      try {
        const data = await adminApi.getSchools();
        setSchools(data || []);
        if (data && data.length > 0) setSelectedSchoolId(data[0].school_id);
      } catch { /* ignore */ }
    };
    load();
  }, []);

  useEffect(() => {
    if (!selectedSchoolId) return;
    const load = async () => {
      try {
        const data = await adminApi.getClassesBySchool(selectedSchoolId);
        setClasses(data || []);
        if (data && data.length > 0) setSelectedClassId(data[0].class_id);
      } catch { setClasses([]); }
    };
    load();
  }, [selectedSchoolId]);

  const handleAssign = async (teacherId: string) => {
    if (!selectedClassId) return;
    try {
      setActionLoading(true); setSuccess("");
      await adminApi.assignTeacherToClass(teacherId, selectedClassId);
      const className = classes.find((c) => c.class_id === selectedClassId)?.name || "";
      setSuccess(`Đã gán giáo viên vào lớp ${className}`);
      setAssigningTeacherId(null);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể gán lớp");
    } finally { setActionLoading(false); }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <BookUser className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý Giáo viên</h1>
      </div>

      {success && <Alert><CheckCircle2 className="h-4 w-4 text-green-600" /><AlertDescription>{success}</AlertDescription></Alert>}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && teachers.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <BookUser className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có giáo viên nào</p>
        </CardContent></Card>
      )}

      {/* Desktop Table */}
      {!loading && teachers.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Họ tên</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Điện thoại</th>
                  <th className="px-6 py-3 font-medium text-right">Gán lớp</th>
                </tr>
              </thead>
              <tbody>
                {teachers.map((t) => (
                  <tr key={t.teacher_id} className="border-b last:border-0 hover:bg-zinc-50">
                    <td className="px-6 py-4 font-medium">{t.full_name}</td>
                    <td className="px-6 py-4 text-muted-foreground">{t.email}</td>
                    <td className="px-6 py-4 text-muted-foreground">{t.phone || "—"}</td>
                    <td className="px-6 py-4 text-right">
                      {assigningTeacherId === t.teacher_id ? (
                        <div className="flex items-center justify-end gap-2">
                          <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
                            <SelectTrigger className="w-[140px]" size="sm"><SelectValue placeholder="Trường" /></SelectTrigger>
                            <SelectContent>
                              {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
                            </SelectContent>
                          </Select>
                          <Select value={selectedClassId} onValueChange={setSelectedClassId}>
                            <SelectTrigger className="w-[120px]" size="sm"><SelectValue placeholder="Lớp" /></SelectTrigger>
                            <SelectContent>
                              {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
                            </SelectContent>
                          </Select>
                          <Button size="sm" onClick={() => handleAssign(t.teacher_id)} disabled={actionLoading}>
                            {actionLoading ? <Loader2 className="h-3 w-3 animate-spin" /> : <Link2 className="h-3 w-3" />}
                          </Button>
                          <Button size="sm" variant="ghost" onClick={() => setAssigningTeacherId(null)}>
                            <Unlink className="h-3 w-3" />
                          </Button>
                        </div>
                      ) : (
                        <Button variant="ghost" size="sm" onClick={() => setAssigningTeacherId(t.teacher_id)}>
                          <Link2 className="mr-1 h-4 w-4" /> Gán lớp
                        </Button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent></Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loading && teachers.length > 0 && (
        <div className="space-y-3 md:hidden">
          {teachers.map((t) => (
            <Card key={t.teacher_id}>
              <CardContent className="py-4">
                <p className="font-medium">{t.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2"><Mail className="h-3 w-3" /> {t.email}</p>
                  {t.phone && <p className="flex items-center gap-2"><Phone className="h-3 w-3" /> {t.phone}</p>}
                </div>
                <div className="mt-3">
                  {assigningTeacherId === t.teacher_id ? (
                    <div className="flex flex-wrap items-center gap-2">
                      <Select value={selectedClassId} onValueChange={setSelectedClassId}>
                        <SelectTrigger className="w-[140px]" size="sm"><SelectValue placeholder="Lớp" /></SelectTrigger>
                        <SelectContent>
                          {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
                        </SelectContent>
                      </Select>
                      <Button size="sm" onClick={() => handleAssign(t.teacher_id)} disabled={actionLoading}>
                        {actionLoading ? <Loader2 className="h-3 w-3 animate-spin" /> : "Gán"}
                      </Button>
                      <Button size="sm" variant="ghost" onClick={() => setAssigningTeacherId(null)}>Hủy</Button>
                    </div>
                  ) : (
                    <Button variant="outline" size="sm" onClick={() => setAssigningTeacherId(t.teacher_id)}>
                      <Link2 className="mr-1 h-3 w-3" /> Gán lớp
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
