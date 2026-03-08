/**
 * Admin Parents Page
 * Danh sách phụ huynh + gán/hủy gán học sinh.
 * API: GET /admin/parents, POST/DELETE /admin/parents/:id/students/:student_id
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Parent, School, Class, Student } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Heart, Loader2, Phone, Mail, Link2, Unlink, AlertCircle, CheckCircle2, Search } from "lucide-react";

export default function AdminParentsPage() {
  const [parents, setParents] = useState<Parent[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [students, setStudents] = useState<Student[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [selectedStudentId, setSelectedStudentId] = useState("");
  const [assigningParentId, setAssigningParentId] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [success, setSuccess] = useState("");

  const fetchParents = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const data = await adminApi.getParents();
      setParents(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách phụ huynh");
    } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchParents(); }, [fetchParents]);

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
        else { setSelectedClassId(""); setStudents([]); }
      } catch { setClasses([]); }
    };
    load();
  }, [selectedSchoolId]);

  useEffect(() => {
    if (!selectedClassId) return;
    const load = async () => {
      try {
        const data = await adminApi.getStudentsByClass(selectedClassId);
        setStudents(data || []);
        if (data && data.length > 0) setSelectedStudentId(data[0].student_id);
        else setSelectedStudentId("");
      } catch { setStudents([]); }
    };
    load();
  }, [selectedClassId]);

  const handleAssign = async (parentId: string) => {
    if (!selectedStudentId) return;
    try {
      setActionLoading(true); setSuccess("");
      await adminApi.assignParentToStudent(parentId, selectedStudentId);
      const studentName = students.find((s) => s.student_id === selectedStudentId)?.full_name || "";
      setSuccess(`Đã gán phụ huynh cho ${studentName}`);
      setAssigningParentId(null);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể gán");
    } finally { setActionLoading(false); }
  };

  const filteredParents = useMemo(() => {
    if (!searchQuery.trim()) return parents;
    const q = searchQuery.toLowerCase();
    return parents.filter((p) => 
      p.full_name?.toLowerCase().includes(q) || 
      p.email?.toLowerCase().includes(q) ||
      p.phone?.includes(q)
    );
  }, [parents, searchQuery]);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Heart className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý Phụ huynh</h1>
      </div>

      {success && <Alert><CheckCircle2 className="h-4 w-4 text-green-600" /><AlertDescription>{success}</AlertDescription></Alert>}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {/* Toolbar: Search box */}
      {!loading && !error && parents.length > 0 && (
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

      {!loading && parents.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <Heart className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có phụ huynh nào</p>
        </CardContent></Card>
      )}

      {!loading && parents.length > 0 && filteredParents.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy phụ huynh nào phù hợp với &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      {/* Desktop Table */}
      {!loading && filteredParents.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Họ tên</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Điện thoại</th>
                  <th className="px-6 py-3 font-medium text-right">Gán học sinh</th>
                </tr>
              </thead>
              <tbody>
                {filteredParents.map((p) => (
                  <tr key={p.parent_id} className="border-b last:border-0 hover:bg-zinc-50">
                    <td className="px-6 py-4 font-medium">{p.full_name}</td>
                    <td className="px-6 py-4 text-muted-foreground">{p.email}</td>
                    <td className="px-6 py-4 text-muted-foreground">{p.phone || "—"}</td>
                    <td className="px-6 py-4 text-right">
                      {assigningParentId === p.parent_id ? (
                        <div className="flex items-center justify-end gap-1">
                          <Select value={selectedClassId} onValueChange={setSelectedClassId}>
                            <SelectTrigger className="w-[120px]" size="sm"><SelectValue placeholder="Lớp" /></SelectTrigger>
                            <SelectContent>
                              {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
                            </SelectContent>
                          </Select>
                          <Select value={selectedStudentId} onValueChange={setSelectedStudentId}>
                            <SelectTrigger className="w-[130px]" size="sm"><SelectValue placeholder="HS" /></SelectTrigger>
                            <SelectContent>
                              {students.map((s) => <SelectItem key={s.student_id} value={s.student_id}>{s.full_name}</SelectItem>)}
                            </SelectContent>
                          </Select>
                          <Button size="sm" onClick={() => handleAssign(p.parent_id)} disabled={actionLoading || !selectedStudentId}>
                            {actionLoading ? <Loader2 className="h-3 w-3 animate-spin" /> : <Link2 className="h-3 w-3" />}
                          </Button>
                          <Button size="sm" variant="ghost" onClick={() => setAssigningParentId(null)}><Unlink className="h-3 w-3" /></Button>
                        </div>
                      ) : (
                        <Button variant="ghost" size="sm" onClick={() => setAssigningParentId(p.parent_id)}>
                          <Link2 className="mr-1 h-4 w-4" /> Gán HS
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
      {!loading && filteredParents.length > 0 && (
        <div className="space-y-3 md:hidden">
          {filteredParents.map((p) => (
            <Card key={p.parent_id}>
              <CardContent className="py-4">
                <p className="font-medium">{p.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2"><Mail className="h-3 w-3" /> {p.email}</p>
                  {p.phone && <p className="flex items-center gap-2"><Phone className="h-3 w-3" /> {p.phone}</p>}
                </div>
                <div className="mt-3">
                  <Button variant="outline" size="sm" onClick={() => setAssigningParentId(p.parent_id)}>
                    <Link2 className="mr-1 h-3 w-3" /> Gán học sinh
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
