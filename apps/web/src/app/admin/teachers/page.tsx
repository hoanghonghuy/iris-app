/**
 * Admin Teachers Page
 * Danh sách giáo viên + gán/hủy gán lớp.
 * API: GET /admin/teachers, POST/DELETE /admin/teachers/:id/classes/:class_id
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Teacher, School, Class, Pagination } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { TableSkeleton } from "@/components/shared/TableSkeleton";
import { CardSkeleton } from "@/components/shared/CardSkeleton";
import { EmptyState } from "@/components/shared/EmptyState";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { BookUser, Loader2, Phone, Mail, Link2, Unlink, Search, X } from "lucide-react";

export default function AdminTeachersPage() {
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [assigningTeacherId, setAssigningTeacherId] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);

  const fetchTeachers = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const response = await adminApi.getTeachers({ limit: 20, offset: currentOffset });
      setTeachers(response.data || []);
      if (response.pagination) setPagination(response.pagination);
    } catch (err: any) {
      const msg = err.response?.data?.error || "Không thể tải danh sách giáo viên";
      setError(msg);
      toast.error(msg);
    } finally { setLoading(false); }
  }, [currentOffset]);

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
      setActionLoading(true);
      await adminApi.assignTeacherToClass(teacherId, selectedClassId);
      const className = classes.find((c) => c.class_id === selectedClassId)?.name || "";
      toast.success(`Đã gán giáo viên vào lớp ${className}`);
      setAssigningTeacherId(null);
      fetchTeachers();
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Không thể gán lớp");
    } finally { setActionLoading(false); }
  };

  const handleUnassign = async (teacherId: string, classId: string, className: string) => {
    // eslint-disable-next-line no-alert
    if (!window.confirm(`Bạn có chắc chắn muốn hủy gán lớp ${className}?`)) return;
    try {
      setActionLoading(true);
      await adminApi.unassignTeacherFromClass(teacherId, classId);
      toast.success(`Đã hủy gán lớp ${className}`);
      fetchTeachers();
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Không thể hủy gán lớp");
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

      {loading && (
        <>
          <div className="hidden md:block">
            <TableSkeleton columns={4} rows={10} />
          </div>
          <div className="md:hidden">
            <CardSkeleton cards={5} />
          </div>
        </>
      )}

      {!loading && teachers.length === 0 && !error && (
        <EmptyState
          icon={BookUser}
          title="Chưa có giáo viên nào"
          description="Hiện tại hệ thống chưa có dữ liệu giáo viên mới."
        />
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

      {/* Pagination */}
      {!loading && teachers.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}
    </div>
  );
}
