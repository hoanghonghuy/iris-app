/**
 * Teacher Health Page
 * Ghi nhận sức khỏe HS: chọn lớp → HS → ghi nhận nhiệt độ, triệu chứng, mức độ.
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/health, GET /teacher/students/:student_id/health
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, HealthLog, Student } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { EmptyState } from "@/components/shared/EmptyState";
import { Loader2, Plus, X, AlertCircle, CheckCircle2, HeartPulse, RefreshCw } from "lucide-react";

type Severity = "normal" | "watch" | "urgent";

const severityOptions: Array<{
  value: Severity;
  label: string;
  variant: "secondary" | "outline" | "destructive";
}> = [
  { value: "normal", label: "Bình thường", variant: "secondary" as const },
  { value: "watch", label: "Theo dõi", variant: "outline" as const },
  { value: "urgent", label: "Khẩn cấp", variant: "destructive" as const },
];

function extractErrorMessage(err: unknown): string | undefined {
  return (
    typeof (err as { response?: { data?: { error?: string } } }).response?.data?.error === "string"
      ? (err as { response?: { data?: { error?: string } } }).response?.data?.error
      : undefined
  );
}

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
  const [severity, setSeverity] = useState<Severity>("normal");
  const [note, setNote] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");

  const [historyStudentId, setHistoryStudentId] = useState("");
  const [historyFrom, setHistoryFrom] = useState("");
  const [historyTo, setHistoryTo] = useState("");
  const [historyLogs, setHistoryLogs] = useState<HealthLog[]>([]);
  const [loadingHistory, setLoadingHistory] = useState(false);
  const [historyError, setHistoryError] = useState("");

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
      if (data && data.length > 0) {
        setFormStudentId(data[0].student_id);
        setHistoryStudentId(data[0].student_id);
      } else {
        setHistoryLogs([]);
      }
    } catch (err: unknown) {
      setError(extractErrorMessage(err) || "Không thể tải HS");
    } finally { setLoadingStudents(false); }
  }, [selectedClassId]);

  useEffect(() => { fetchStudents(); }, [fetchStudents]);

  const fetchHistory = useCallback(async () => {
    if (!historyStudentId) {
      setHistoryLogs([]);
      return;
    }

    try {
      setLoadingHistory(true);
      setHistoryError("");
      const logs = await teacherApi.getStudentHealth(
        historyStudentId,
        historyFrom || undefined,
        historyTo || undefined,
      );
      setHistoryLogs(logs || []);
    } catch (err: unknown) {
      setHistoryError(extractErrorMessage(err) || "Không thể tải lịch sử sức khỏe");
    } finally {
      setLoadingHistory(false);
    }
  }, [historyStudentId, historyFrom, historyTo]);

  useEffect(() => {
    fetchHistory();
  }, [fetchHistory]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formStudentId) { setFormError("Chọn học sinh"); return; }
    try {
      setSubmitting(true); setFormError(""); setSuccess("");
      await teacherApi.createHealthLog({
        student_id: formStudentId,
        temperature: temperature ? parseFloat(temperature) : undefined,
        symptoms: symptoms || undefined,
        severity,
        note: note || undefined,
      });
      setSuccess("Đã ghi nhận sức khỏe thành công!");
      setTemperature(""); setSymptoms(""); setSeverity("normal"); setNote("");
      fetchHistory();
    } catch (err: unknown) {
      setFormError(extractErrorMessage(err) || "Lỗi ghi nhận");
    } finally { setSubmitting(false); }
  };

  const severityLabel: Record<Severity, string> = {
    normal: "Bình thường",
    watch: "Theo dõi",
    urgent: "Khẩn cấp",
  };

  const severityVariant: Record<Severity, "secondary" | "outline" | "destructive"> = {
    normal: "secondary",
    watch: "outline",
    urgent: "destructive",
  };

  if (loadingClasses) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-end">
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
      {success && <Alert><CheckCircle2 className="h-4 w-4 text-success" /><AlertDescription>{success}</AlertDescription></Alert>}

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
                        className={`cursor-pointer select-none transition-all ${severity === opt.value ? "ring-2 ring-offset-1 ring-ring" : "opacity-60 hover:opacity-100"}`}
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

      {!loadingStudents && students.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Lịch sử sức khỏe</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-3 sm:grid-cols-4">
              <div className="space-y-2 sm:col-span-2">
                <Label>Học sinh</Label>
                <Select value={historyStudentId} onValueChange={setHistoryStudentId}>
                  <SelectTrigger className="w-full"><SelectValue placeholder="Chọn HS" /></SelectTrigger>
                  <SelectContent>
                    {students.map((s) => (
                      <SelectItem key={s.student_id} value={s.student_id}>{s.full_name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="historyFrom">Từ ngày</Label>
                <Input id="historyFrom" type="date" value={historyFrom} onChange={(e) => setHistoryFrom(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="historyTo">Đến ngày</Label>
                <Input id="historyTo" type="date" value={historyTo} onChange={(e) => setHistoryTo(e.target.value)} />
              </div>
            </div>

            <div className="flex justify-end">
              <Button variant="outline" onClick={fetchHistory} disabled={loadingHistory || !historyStudentId}>
                {loadingHistory ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <RefreshCw className="mr-2 h-4 w-4" />
                )}
                Làm mới
              </Button>
            </div>

            {historyError && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{historyError}</AlertDescription>
              </Alert>
            )}

            {!loadingHistory && historyLogs.length === 0 && !historyError && (
              <EmptyState
                icon={HeartPulse}
                title="Chưa có nhật ký sức khỏe"
                description="Nhật ký mới sẽ hiển thị sau khi giáo viên ghi nhận cho học sinh."
              />
            )}

            {!loadingHistory && historyLogs.length > 0 && (
              <div className="space-y-3">
                {historyLogs.map((log) => (
                  <div key={log.health_log_id} className="rounded-lg border p-3">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <p className="text-sm text-muted-foreground">{new Date(log.recorded_at).toLocaleString("vi-VN")}</p>
                      <Badge variant={severityVariant[log.severity as Severity] || "outline"}>
                        {severityLabel[log.severity as Severity] || log.severity}
                      </Badge>
                    </div>
                    <div className="mt-2 space-y-1 text-sm">
                      <p><span className="font-medium">Nhiệt độ:</span> {typeof log.temperature === "number" ? `${log.temperature}°C` : "Không ghi"}</p>
                      <p><span className="font-medium">Triệu chứng:</span> {log.symptoms || "Không ghi"}</p>
                      <p><span className="font-medium">Ghi chú:</span> {log.note || "Không ghi"}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
