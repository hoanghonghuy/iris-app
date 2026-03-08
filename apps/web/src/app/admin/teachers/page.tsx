/**
 * Admin Teachers Page
 * Danh sách giáo viên + gán/hủy gán lớp.
 * API: GET /admin/teachers, POST/DELETE /admin/teachers/:id/classes/:class_id
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Teacher, School, Class } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { BookUser, Loader2, Phone, Mail, Link2, Unlink, AlertCircle, CheckCircle2, Search, X } from "lucide-react";

export default function AdminTeachersPage() {
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
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
      fetchTeachers();
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể gán lớp");
    } finally { setActionLoading(false); }
  };

  const handleUnassign = async (teacherId: string, classId: string, className: string) => {
    // eslint-disable-next-line no-alert
    if (!window.confirm(`Bạn có chắc chắn muốn hủy gán lớp ${className}?`)) return;
    try {
      setActionLoading(true); setSuccess("");
      await adminApi.unassignTeacherFromClass(teacherId, classId);
      setSuccess(`Đã hủy gán lớp ${className}`);
      fetchTeachers();
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể hủy gán lớp");
    } finally { setActionLoading(false); }
  };

  const filteredTeachers = useMemo(() => {
    if (!searchQuery.trim()) return teachers;
    const q = searchQuery.toLowerCase();
    return teachers.filter((t) => 
      t.full_name?.toLowerCase().includes(q) || 
      t.email?.toLowerCase().includes(q) ||
      t.phone?.includes(q)
    );
  }, [teachers, searchQuery]);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <BookUser className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý Giáo viên</h1>
      </div>

      {success && <Alert><CheckCircle2 className="h-4 w-4 text-green-600" /><AlertDescription>{success}</AlertDescription></Alert>}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {/* Toolbar: Search */}
      {!loading && !error && teachers.length > 0 && (
        <div className="relative max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input 
            type="search" 
            placeholder="Tìm theo tên, email, SĐT..." 
            className="pl-8 bg-white" 
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      )}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && teachers.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <BookUser className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có giáo viên nào</p>
        </CardContent></Card>
      )}

      {!loading && teachers.length > 0 && filteredTeachers.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy giáo viên nào phù hợp với &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      {/* Desktop Table */}
      {!loading && filteredTeachers.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Họ tên</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Lớp Phụ Trách</th>
                  <th className="px-6 py-3 font-medium text-right">Gán lớp</th>
                </tr>
              </thead>
              <tbody>
                {filteredTeachers.map((t) => (
                  <tr key={t.teacher_id} className="border-b last:border-0 hover:bg-zinc-50 leading-relaxed">
                    <td className="px-6 py-4">
                      <div className="font-medium text-slate-900">{t.full_name}</div>
                      <div className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                        <Phone className="h-3 w-3" /> {t.phone || "—"}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-muted-foreground">{t.email}</td>
                    <td className="px-6 py-4">
                      {t.classes && t.classes.length > 0 ? (
                        <div className="flex flex-wrap gap-1.5">
                          {t.classes.map(c => (
                            <Badge key={c.class_id} variant="secondary" className="pr-1.5 flex items-center gap-1">
                              {c.name}
                              <button 
                                onClick={() => handleUnassign(t.teacher_id, c.class_id, c.name)}
                                className="ml-0.5 rounded-full p-0.5 hover:bg-destructive/20 hover:text-destructive transition-colors focus:outline-none focus:ring-2 focus:ring-ring"
                                aria-label="Remove"
                              >
                                <X className="h-3 w-3" />
                              </button>
                            </Badge>
                          ))}
                        </div>
                      ) : (
                        <span className="text-sm text-muted-foreground italic">Chưa phân lớp</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-right align-middle">
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
      {!loading && filteredTeachers.length > 0 && (
        <div className="space-y-3 md:hidden">
          {filteredTeachers.map((t) => (
            <Card key={t.teacher_id}>
              <CardContent className="py-4">
                <p className="font-medium">{t.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2"><Mail className="h-3 w-3" /> {t.email}</p>
                  {t.phone && <p className="flex items-center gap-2"><Phone className="h-3 w-3" /> {t.phone}</p>}
                </div>
                
                <div className="mt-4 border-t border-dashed pt-3">
                  <p className="text-xs text-muted-foreground mb-2 font-medium uppercase tracking-wider">Lớp phụ trách</p>
                  {t.classes && t.classes.length > 0 ? (
                    <div className="flex flex-wrap gap-1.5">
                      {t.classes.map(c => (
                        <Badge key={c.class_id} variant="secondary" className="pl-2 pr-1.5 py-0.5 flex items-center gap-1">
                          {c.name}
                          <button 
                            onClick={(e) => { e.preventDefault(); handleUnassign(t.teacher_id, c.class_id, c.name); }}
                            className="ml-1 rounded-full bg-transparent p-0.5 text-muted-foreground hover:bg-destructive/20 hover:text-destructive transition-colors outline-none"
                            aria-label="Remove class"
                          >
                            <X className="h-3 w-3" />
                          </button>
                        </Badge>
                      ))}
                    </div>
                  ) : (
                    <div className="text-sm text-muted-foreground italic">
                      Chưa phân lớp
                    </div>
                  )}
                </div>

                <div className="mt-4 pt-4 border-t border-zinc-100 flex items-center justify-start">
                  {assigningTeacherId === t.teacher_id ? (
                    <div className="flex flex-wrap items-center gap-2 w-full">
                      <div className="flex w-full gap-2">
                         <Select value={selectedClassId} onValueChange={setSelectedClassId}>
                           <SelectTrigger className="flex-1" size="sm"><SelectValue placeholder="Chọn Lớp" /></SelectTrigger>
                           <SelectContent>
                             {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
                           </SelectContent>
                         </Select>
                      </div>
                      <div className="flex w-full gap-2 mt-1">
                        <Button size="sm" className="flex-1" onClick={() => handleAssign(t.teacher_id)} disabled={actionLoading}>
                          {actionLoading ? <Loader2 className="h-3 w-3 animate-spin mr-1" /> : <Link2 className="h-4 w-4 mr-1" />} Gán
                        </Button>
                        <Button size="sm" variant="outline" className="flex-1" onClick={() => setAssigningTeacherId(null)}>Hủy</Button>
                      </div>
                    </div>
                  ) : (
                    <Button variant="outline" size="sm" className="w-full text-blue-600 bg-blue-50 hover:bg-blue-100 border-blue-200" onClick={() => setAssigningTeacherId(t.teacher_id)}>
                      <Link2 className="mr-1 h-4 w-4" /> Gán phân lớp
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
