/**
 * Admin Classes Page
 * Quản lý lớp học theo trường: chọn trường → xem danh sách lớp + tạo lớp mới.
 * API: GET /admin/classes/by-school/:school_id, POST /admin/classes
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { School, Class } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { GraduationCap, Plus, X, Loader2, Calendar, AlertCircle } from "lucide-react";
import { useAuth } from "@/providers/AuthProvider";
export default function AdminClassesPage() {
  const { role } = useAuth();
  const [schools, setSchools] = useState<School[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState<string>("");
  const [loadingSchools, setLoadingSchools] = useState(true);

  const [classes, setClasses] = useState<Class[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(false);
  const [error, setError] = useState("");

  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ name: "", school_year: "" });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const response = await adminApi.getSchools();
        const data = response.data;
        setSchools(data || []);
        if (data && data.length > 0) setSelectedSchoolId(data[0].school_id);
      } catch { setError("Không thể tải danh sách trường"); }
      finally { setLoadingSchools(false); }
    };
    load();
  }, []);

  const fetchClasses = useCallback(async () => {
    if (!selectedSchoolId) return;
    try {
      setLoadingClasses(true); setError("");
      const response = await adminApi.getClassesBySchool(selectedSchoolId);
      const data = response.data;
      setClasses(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách lớp");
    } finally { setLoadingClasses(false); }
  }, [selectedSchoolId]);

  useEffect(() => { fetchClasses(); }, [fetchClasses]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) { setFormError("Tên lớp không được để trống"); return; }
    if (!formData.school_year.trim()) { setFormError("Năm học không được để trống"); return; }
    try {
      setSubmitting(true); setFormError("");
      await adminApi.createClass({ school_id: selectedSchoolId, name: formData.name, school_year: formData.school_year });
      setFormData({ name: "", school_year: "" }); setShowForm(false); fetchClasses();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo lớp");
    } finally { setSubmitting(false); }
  };

  const selectedSchoolName = schools.find((s) => s.school_id === selectedSchoolId)?.name || "";

  if (loadingSchools) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2">
          {role === 'SUPER_ADMIN' && (
            <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
              <SelectTrigger className="w-[200px]"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
              <SelectContent>
                {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          {selectedSchoolId && (
            <Button size="sm" onClick={() => setShowForm(!showForm)}>
              {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
              {showForm ? "Hủy" : "Thêm lớp"}
            </Button>
          )}
        </div>
      </div>

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Thêm lớp mới — {selectedSchoolName}</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="className">Tên lớp <span className="text-destructive">*</span></Label>
                  <Input id="className" placeholder="VD: Lá Non" value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="schoolYear">Năm học <span className="text-destructive">*</span></Label>
                  <Input id="schoolYear" placeholder="VD: 2025-2026" value={formData.school_year}
                    onChange={(e) => setFormData({ ...formData, school_year: e.target.value })} required />
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Tạo lớp
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {loadingClasses && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingClasses && !error && classes.length === 0 && selectedSchoolId && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <GraduationCap className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có lớp nào trong {selectedSchoolName}</p>
          <Button variant="outline" className="mt-4" onClick={() => setShowForm(true)}>
            <Plus className="mr-2 h-4 w-4" /> Thêm lớp đầu tiên
          </Button>
        </CardContent></Card>
      )}

      {/* Desktop Table */}
      {!loadingClasses && classes.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Tên lớp</th>
                  <th className="px-6 py-3 font-medium">Năm học</th>
                </tr>
              </thead>
              <tbody>
                {classes.map((cls) => (
                  <tr key={cls.class_id} className="border-b last:border-0 hover:bg-muted">
                    <td className="px-6 py-4 font-medium">{cls.name}</td>
                    <td className="px-6 py-4"><Badge variant="secondary">{cls.school_year}</Badge></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent></Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loadingClasses && classes.length > 0 && (
        <div className="space-y-3 md:hidden">
          {classes.map((cls) => (
            <Card key={cls.class_id}>
              <CardContent className="flex items-start gap-3 py-4">
                <GraduationCap className="mt-0.5 h-5 w-5 shrink-0 text-muted-foreground" />
                <div className="min-w-0 flex-1">
                  <p className="font-medium">{cls.name}</p>
                  <div className="mt-1 flex items-center gap-1 text-sm text-muted-foreground">
                    <Calendar className="h-3 w-3" />
                    <Badge variant="secondary">{cls.school_year}</Badge>
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
