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
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Heart, Loader2, Plus, X, AlertCircle, CheckCircle2 } from "lucide-react";

const severityOptions = [
  { value: "normal", label: "Bình thường", variant: "secondary" as const },
  { value: "watch", label: "Theo dõi", variant: "outline" as const },
  { value: "urgent", label: "Khẩn cấp", variant: "destructive" as const },
];

export default function TeacherHealthPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

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
      setLoadingStudents(true); setError("");
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
      setSubmitting(true); setFormError(""); setSuccess("");
      await teacherApi.createHealthLog({
        student_id: formStudentId,
        temperature: temperature ? parseFloat(temperature) : undefined,
        symptoms: symptoms || undefined,
        severity: severity as any,
        note: note || undefined,
      });
      setSuccess("Đã ghi nhận sức khỏe thành công!");
      setTemperature(""); setSymptoms(""); setSeverity("normal"); setNote("");
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
            <Select value={selectedClassId} onValueChange={setSelectedClassId}>
              <SelectTrigger className="w-[180px]"><SelectValue placeholder="Chọn lớp" /></SelectTrigger>
              <SelectContent>
                {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          <Button size="sm" onClick={() => setShowForm(!showForm)}>
            {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
            {showForm ? "Đóng" : "Ghi nhận"}
          </Button>
        </div>
      </div>

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}
      {success && <Alert><CheckCircle2 className="h-4 w-4 text-green-600" /><AlertDescription>{success}</AlertDescription></Alert>}

      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Ghi nhận sức khỏe</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {formError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>Học sinh</Label>
                  <Select value={formStudentId} onValueChange={setFormStudentId}>
                    <SelectTrigger className="w-full"><SelectValue placeholder="Chọn HS" /></SelectTrigger>
                    <SelectContent>
                      {students.map((s) => <SelectItem key={s.student_id} value={s.student_id}>{s.full_name}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="temperature">Nhiệt độ (°C)</Label>
                  <Input id="temperature" type="number" step="0.1" placeholder="36.5" value={temperature} onChange={(e) => setTemperature(e.target.value)} />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="symptoms">Triệu chứng</Label>
                <Input id="symptoms" placeholder="VD: Ho nhẹ, sổ mũi..." value={symptoms} onChange={(e) => setSymptoms(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Mức độ</Label>
                <div className="flex gap-2">
                  {severityOptions.map((opt) => (
                    <Badge key={opt.value} variant={severity === opt.value ? opt.variant : "outline"}
                      className={`cursor-pointer select-none transition-all ${severity === opt.value ? "ring-2 ring-offset-1 ring-zinc-400" : "opacity-60 hover:opacity-100"}`}
                      onClick={() => setSeverity(opt.value)}>{opt.label}</Badge>
                  ))}
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="healthNote">Ghi chú</Label>
                <Textarea id="healthNote" placeholder="Ghi chú thêm..." value={note} onChange={(e) => setNote(e.target.value)} />
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Lưu
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loadingStudents && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingStudents && !showForm && students.length > 0 && (
        <Card><CardContent className="py-4">
          <p className="text-sm text-muted-foreground">{students.length} học sinh trong lớp. Nhấn &ldquo;Ghi nhận&rdquo; để thêm nhật ký sức khỏe.</p>
        </CardContent></Card>
      )}
    </div>
  );
}
