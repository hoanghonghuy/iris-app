/**
 * Admin Students Page
 * Quản lý học sinh theo lớp: chọn trường → chọn lớp → xem danh sách + tạo mới + tạo mã phụ huynh.
 * API: GET /admin/students/by-class/:class_id, POST /admin/students, POST /admin/students/:id/generate-parent-code
 *
 * TODO: add server-side pagination when student count per class grows
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { School, Class, Student } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Users, Plus, X, Loader2, Calendar, User, KeyRound, Copy, Check, AlertCircle, Search,
} from "lucide-react";

export default function AdminStudentsPage() {
  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [loadingSchools, setLoadingSchools] = useState(true);
  const [loadingClasses, setLoadingClasses] = useState(false);

  const [students, setStudents] = useState<Student[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ full_name: "", dob: "", gender: "male" });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  const [generatingCode, setGeneratingCode] = useState<string | null>(null);
  const [parentCodes, setParentCodes] = useState<Record<string, string>>({});
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [codeError, setCodeError] = useState("");

  const genderLabel: Record<string, string> = { male: "Nam", female: "Nữ", other: "Khác" };

  // ─── Fetch schools ─────────────────────────────────────────────
  useEffect(() => {
    const load = async () => {
      try {
        const data = await adminApi.getSchools();
        setSchools(data || []);
        if (data && data.length > 0) setSelectedSchoolId(data[0].school_id);
      } catch { setError("Không thể tải danh sách trường"); }
      finally { setLoadingSchools(false); }
    };
    load();
  }, []);

  // ─── Fetch classes khi đổi trường ─────────────────────────────
  useEffect(() => {
    if (!selectedSchoolId) return;
    const load = async () => {
      try {
        setLoadingClasses(true); setSelectedClassId(""); setStudents([]); setSearchQuery("");
        const data = await adminApi.getClassesBySchool(selectedSchoolId);
        setClasses(data || []);
        if (data && data.length > 0) setSelectedClassId(data[0].class_id);
      } catch { setClasses([]); }
      finally { setLoadingClasses(false); }
    };
    load();
  }, [selectedSchoolId]);

  // ─── Fetch students khi đổi lớp ──────────────────────────────
  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) return;
    try {
      setLoadingStudents(true); setError("");
      const data = await adminApi.getStudentsByClass(selectedClassId);
      setStudents(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách học sinh");
    } finally { setLoadingStudents(false); }
  }, [selectedClassId]);

  useEffect(() => { fetchStudents(); }, [fetchStudents]);

  // ─── Filter students by search query ──────────────────────────
  const filteredStudents = useMemo(() => {
    if (!searchQuery.trim()) return students;
    const q = searchQuery.toLowerCase();
    return students.filter((s) => s.full_name.toLowerCase().includes(q));
  }, [students, searchQuery]);

  // ─── Create student ───────────────────────────────────────────
  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.full_name.trim()) { setFormError("Họ tên không được để trống"); return; }
    try {
      setSubmitting(true); setFormError("");
      await adminApi.createStudent({
        school_id: selectedSchoolId, class_id: selectedClassId,
        full_name: formData.full_name, dob: formData.dob,
        gender: formData.gender as "male" | "female" | "other",
      });
      setFormData({ full_name: "", dob: "", gender: "male" });
      setShowForm(false); fetchStudents();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo học sinh");
    } finally { setSubmitting(false); }
  };

  // ─── Generate Parent Code ─────────────────────────────────────
  const handleGenerateCode = async (studentId: string) => {
    try {
      setGeneratingCode(studentId); setCodeError("");
      const res = await adminApi.generateParentCode(studentId);
      const code = (res as any)?.data?.parent_code || (res as any)?.parent_code || "";
      setParentCodes((prev) => ({ ...prev, [studentId]: code }));
    } catch (err: any) {
      setCodeError(err.response?.data?.error || "Không thể tạo mã");
    } finally { setGeneratingCode(null); }
  };

  const handleCopy = (studentId: string) => {
    navigator.clipboard.writeText(parentCodes[studentId] || "");
    setCopiedId(studentId);
    setTimeout(() => setCopiedId(null), 2000);
  };

  const selectedClassName = classes.find((c) => c.class_id === selectedClassId)?.name || "";

  if (loadingSchools) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
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
          <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
            <SelectTrigger className="w-[180px]"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
            <SelectContent>
              {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
            </SelectContent>
          </Select>
          {classes.length > 0 && (
            <Select value={selectedClassId} onValueChange={setSelectedClassId}>
              <SelectTrigger className="w-[160px]"><SelectValue placeholder="Chọn lớp" /></SelectTrigger>
              <SelectContent>
                {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          {selectedClassId && (
            <Button size="sm" onClick={() => setShowForm(!showForm)}>
              {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
              {showForm ? "Hủy" : "Thêm HS"}
            </Button>
          )}
        </div>
      </div>

      {/* Toolbar: Search */}
      {!loadingStudents && !error && students.length > 0 && !showForm && (
        <div className="relative max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input 
            type="search" 
            placeholder="Tìm theo tên học sinh..." 
            className="pl-8 bg-white" 
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      )}

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Thêm học sinh — {selectedClassName}</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && (
                <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>
              )}
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="fullName">Họ tên <span className="text-destructive">*</span></Label>
                  <Input id="fullName" placeholder="VD: Bé An" value={formData.full_name}
                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="dob">Ngày sinh <span className="text-destructive">*</span></Label>
                  <Input id="dob" type="date" value={formData.dob}
                    onChange={(e) => setFormData({ ...formData, dob: e.target.value })} required />
                </div>
                <div className="space-y-2">
                  <Label>Giới tính</Label>
                  <Select value={formData.gender} onValueChange={(v) => setFormData({ ...formData, gender: v })}>
                    <SelectTrigger className="w-full"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="male">Nam</SelectItem>
                      <SelectItem value="female">Nữ</SelectItem>
                      <SelectItem value="other">Khác</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Tạo học sinh
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Errors */}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}
      {codeError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{codeError}</AlertDescription></Alert>}

      {/* Loading */}
      {loadingStudents && (
        <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>
      )}

      {/* Empty (No students at all) */}
      {!loadingStudents && !error && students.length === 0 && selectedClassId && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Users className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Chưa có học sinh nào trong {selectedClassName}</p>
            <Button variant="outline" className="mt-4" onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" /> Thêm học sinh đầu tiên
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Empty Search Results */}
      {!loadingStudents && !error && students.length > 0 && filteredStudents.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy học sinh nào phù hợp với &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      {/* Desktop Table */}
      {!loadingStudents && filteredStudents.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Họ tên</th>
                  <th className="px-6 py-3 font-medium">Ngày sinh</th>
                  <th className="px-6 py-3 font-medium">Giới tính</th>
                  <th className="px-6 py-3 font-medium text-right">Mã PH</th>
                </tr>
              </thead>
              <tbody>
                {filteredStudents.map((s) => (
                  <tr key={s.student_id} className="border-b last:border-0 hover:bg-zinc-50">
                    <td className="px-6 py-4 font-medium">{s.full_name}</td>
                    <td className="px-6 py-4 text-muted-foreground">{s.dob}</td>
                    <td className="px-6 py-4"><Badge variant="secondary">{genderLabel[s.gender] || s.gender}</Badge></td>
                    <td className="px-6 py-4 text-right">
                      {parentCodes[s.student_id] ? (
                        <div className="flex items-center justify-end gap-1">
                          <code className="rounded bg-zinc-100 px-2 py-0.5 text-xs font-mono">{parentCodes[s.student_id]}</code>
                          <Button variant="ghost" size="sm" onClick={() => handleCopy(s.student_id)}>
                            {copiedId === s.student_id ? <Check className="h-3 w-3 text-green-600" /> : <Copy className="h-3 w-3" />}
                          </Button>
                        </div>
                      ) : (
                        <Button variant="ghost" size="sm" onClick={() => handleGenerateCode(s.student_id)} disabled={generatingCode === s.student_id}>
                          {generatingCode === s.student_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <KeyRound className="mr-1 h-4 w-4" />} Tạo mã
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
      {!loadingStudents && filteredStudents.length > 0 && (
        <div className="space-y-3 md:hidden">
          {filteredStudents.map((s) => (
            <Card key={s.student_id}>
              <CardContent className="flex items-start gap-3 py-4">
                <User className="mt-0.5 h-5 w-5 shrink-0 text-muted-foreground" />
                <div className="min-w-0 flex-1">
                  <p className="font-medium">{s.full_name}</p>
                  <div className="mt-1 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
                    <span className="flex items-center gap-1"><Calendar className="h-3 w-3" /> {s.dob}</span>
                    <Badge variant="secondary">{genderLabel[s.gender] || s.gender}</Badge>
                  </div>
                  {parentCodes[s.student_id] ? (
                    <div className="mt-2 flex items-center gap-1">
                      <code className="rounded bg-zinc-100 px-2 py-0.5 text-xs font-mono">{parentCodes[s.student_id]}</code>
                      <Button variant="ghost" size="sm" onClick={() => handleCopy(s.student_id)}>
                        {copiedId === s.student_id ? <Check className="h-3 w-3 text-green-600" /> : <Copy className="h-3 w-3" />}
                      </Button>
                    </div>
                  ) : (
                    <Button variant="ghost" size="sm" className="mt-2" onClick={() => handleGenerateCode(s.student_id)} disabled={generatingCode === s.student_id}>
                      {generatingCode === s.student_id ? <Loader2 className="h-3 w-3 animate-spin" /> : <KeyRound className="mr-1 h-3 w-3" />} Tạo mã PH
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
