"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Appointment, AppointmentStatus, Class } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Loader2 } from "lucide-react";

export default function TeacherAppointmentsPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [classId, setClassId] = useState("");
  const [startTime, setStartTime] = useState("");
  const [durationMinutes, setDurationMinutes] = useState(30);
  const [note, setNote] = useState("");

  const [statusFilter, setStatusFilter] = useState<"" | AppointmentStatus>("");

  const availableClassOptions = useMemo(() => classes ?? [], [classes]);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [classData, appointmentRes] = await Promise.all([
        teacherApi.getMyClasses(),
        teacherApi.getAppointments({ limit: 50, offset: 0, status: statusFilter || undefined }),
      ]);
      setClasses(classData || []);
      setAppointments(appointmentRes.data || []);
      if (!classId && classData?.length) {
        setClassId(classData[0].class_id);
      }
    } finally {
      setLoading(false);
    }
  }, [classId, statusFilter]);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  const createSlot = async () => {
    if (!classId || !startTime) return;
    setSubmitting(true);
    try {
      const startISO = new Date(startTime).toISOString();
      await teacherApi.createAppointmentSlot({
        class_id: classId,
        start_time: startISO,
        duration_minutes: durationMinutes,
        note,
      });
      setNote("");
      await loadData();
    } finally {
      setSubmitting(false);
    }
  };

  const updateStatus = async (appointmentId: string, status: AppointmentStatus) => {
    await teacherApi.updateAppointmentStatus(appointmentId, status, status === "cancelled" ? "teacher_cancelled" : undefined);
    await loadData();
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Lịch hẹn phụ huynh</h1>

      <Card>
        <CardContent className="p-4 space-y-3">
          <h2 className="font-semibold">Tạo khung giờ mới</h2>
          <div className="grid gap-3 md:grid-cols-4">
            <select className="border rounded px-3 py-2" value={classId} onChange={(e) => setClassId(e.target.value)}>
              <option value="">Chọn lớp</option>
              {availableClassOptions.map((c) => (
                <option key={c.class_id} value={c.class_id}>{c.name}</option>
              ))}
            </select>
            <input className="border rounded px-3 py-2" type="datetime-local" value={startTime} onChange={(e) => setStartTime(e.target.value)} />
            <input className="border rounded px-3 py-2" type="number" min={15} step={5} value={durationMinutes} onChange={(e) => setDurationMinutes(Number(e.target.value || 30))} />
            <input className="border rounded px-3 py-2" placeholder="Ghi chú" value={note} onChange={(e) => setNote(e.target.value)} />
          </div>
          <button className="px-4 py-2 rounded bg-primary text-primary-foreground disabled:opacity-50" onClick={createSlot} disabled={submitting || !classId || !startTime}>
            {submitting ? "Đang tạo..." : "Tạo slot"}
          </button>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4 space-y-3">
          <div className="flex items-center justify-between">
            <h2 className="font-semibold">Danh sách lịch hẹn</h2>
            <select className="border rounded px-3 py-2" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value as "" | AppointmentStatus)}>
              <option value="">Tất cả trạng thái</option>
              <option value="pending">pending</option>
              <option value="confirmed">confirmed</option>
              <option value="cancelled">cancelled</option>
              <option value="completed">completed</option>
              <option value="no_show">no_show</option>
            </select>
          </div>

          {loading ? (
            <div className="py-8 flex justify-center"><Loader2 className="h-6 w-6 animate-spin" /></div>
          ) : appointments.length === 0 ? (
            <p className="text-sm text-muted-foreground">Chưa có lịch hẹn.</p>
          ) : (
            <div className="space-y-2">
              {appointments.map((a) => (
                <div key={a.appointment_id} className="border rounded p-3 flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                  <div className="text-sm">
                    <div><b>{a.student_name || a.student_id}</b> - {a.class_name || a.class_id}</div>
                    <div className="text-muted-foreground">{new Date(a.start_time).toLocaleString("vi-VN")} - {new Date(a.end_time).toLocaleString("vi-VN")}</div>
                    <div>Phụ huynh: {a.parent_name || a.parent_id}</div>
                    <div>Trạng thái: <b>{a.status}</b></div>
                  </div>
                  <div className="flex gap-2">
                    <button className="border rounded px-3 py-1 text-xs" onClick={() => updateStatus(a.appointment_id, "confirmed")} disabled={a.status !== "pending"}>Xác nhận</button>
                    <button className="border rounded px-3 py-1 text-xs" onClick={() => updateStatus(a.appointment_id, "completed")} disabled={a.status !== "confirmed"}>Hoàn tất</button>
                    <button className="border rounded px-3 py-1 text-xs" onClick={() => updateStatus(a.appointment_id, "no_show")} disabled={a.status !== "confirmed"}>No show</button>
                    <button className="border rounded px-3 py-1 text-xs" onClick={() => updateStatus(a.appointment_id, "cancelled")} disabled={a.status === "cancelled" || a.status === "completed"}>Hủy</button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
