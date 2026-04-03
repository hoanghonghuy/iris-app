"use client";

import { useCallback, useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Appointment, AppointmentSlot, ParentAnalytics, Student } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Loader2 } from "lucide-react";

export default function ParentAppointmentsPage() {
  const [children, setChildren] = useState<Student[]>([]);
  const [selectedStudentId, setSelectedStudentId] = useState("");
  const [slots, setSlots] = useState<AppointmentSlot[]>([]);
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [stats, setStats] = useState<ParentAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [bookingNote, setBookingNote] = useState("");

  const resolveSyncedStudentId = useCallback((
    nextChildren: Student[],
    nextAppointments: Appointment[],
    explicitStudentId?: string,
  ) => {
    const hasStudent = (id: string) => nextChildren.some((s) => s.student_id === id);

    if (explicitStudentId && hasStudent(explicitStudentId)) {
      return explicitStudentId;
    }

    if (selectedStudentId && hasStudent(selectedStudentId)) {
      return selectedStudentId;
    }

    const activeAppointment = nextAppointments.find(
      (a) => a.status === "pending" || a.status === "confirmed",
    );

    if (activeAppointment?.student_id && hasStudent(activeAppointment.student_id)) {
      return activeAppointment.student_id;
    }

    const anyAppointmentStudentId = nextAppointments[0]?.student_id;
    if (anyAppointmentStudentId && hasStudent(anyAppointmentStudentId)) {
      return anyAppointmentStudentId;
    }

    return nextChildren[0]?.student_id || "";
  }, [selectedStudentId]);

  const load = useCallback(async (studentId?: string) => {
    setLoading(true);
    try {
      const [childrenData, appointmentsRes, analytics] = await Promise.all([
        parentApi.getMyChildren(),
        parentApi.getAppointments({ limit: 50, offset: 0 }),
        parentApi.getAnalytics(),
      ]);
      setChildren(childrenData || []);
      setAppointments(appointmentsRes.data || []);
      setStats(analytics);

      const effectiveStudentId = resolveSyncedStudentId(
        childrenData || [],
        appointmentsRes.data || [],
        studentId,
      );

      if (effectiveStudentId) {
        setSelectedStudentId(effectiveStudentId);
        const slotRes = await parentApi.getAvailableAppointmentSlots({ student_id: effectiveStudentId, limit: 50, offset: 0 });
        setSlots(slotRes.data || []);
      } else {
        setSlots([]);
      }
    } finally {
      setLoading(false);
    }
  }, [resolveSyncedStudentId]);

  useEffect(() => {
    void load();
  }, [load]);

  const onSelectStudent = async (studentId: string) => {
    setSelectedStudentId(studentId);
    const slotRes = await parentApi.getAvailableAppointmentSlots({ student_id: studentId, limit: 50, offset: 0 });
    setSlots(slotRes.data || []);
  };

  const syncStudentFromAppointment = async (studentId: string) => {
    await onSelectStudent(studentId);
  };

  const bookSlot = async (slotId: string) => {
    if (!selectedStudentId) return;
    await parentApi.createAppointment({ student_id: selectedStudentId, slot_id: slotId, note: bookingNote || undefined });
    setBookingNote("");
    await load(selectedStudentId);
  };

  const cancelAppointment = async (appointmentId: string) => {
    await parentApi.cancelAppointment(appointmentId, "parent_cancelled");
    await load(selectedStudentId);
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Lịch hẹn với giáo viên</h1>

      <Card>
        <CardContent className="p-4 grid gap-2 md:grid-cols-4 text-sm">
          <div><b>Tổng số con:</b> {stats?.total_children ?? 0}</div>
          <div><b>Lịch sắp tới:</b> {stats?.upcoming_appointments ?? 0}</div>
          <div><b>Bài đăng 7 ngày:</b> {stats?.recent_posts_7d ?? 0}</div>
          <div><b>Cảnh báo sức khỏe:</b> {stats?.recent_health_alerts_7d ?? 0}</div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4 space-y-3">
          <h2 className="font-semibold">Đặt lịch</h2>
          <div className="grid gap-3 md:grid-cols-3">
            <select className="border rounded px-3 py-2" value={selectedStudentId} onChange={(e) => void onSelectStudent(e.target.value)}>
              <option value="">Chọn học sinh</option>
              {children.map((s) => (
                <option value={s.student_id} key={s.student_id}>{s.full_name}</option>
              ))}
            </select>
            <input className="border rounded px-3 py-2 md:col-span-2" placeholder="Ghi chú khi đặt lịch" value={bookingNote} onChange={(e) => setBookingNote(e.target.value)} />
          </div>

          {loading ? (
            <div className="py-8 flex justify-center"><Loader2 className="h-6 w-6 animate-spin" /></div>
          ) : slots.length === 0 ? (
            <p className="text-sm text-muted-foreground">Không có slot phù hợp ở thời điểm hiện tại.</p>
          ) : (
            <div className="space-y-2">
              {slots.map((slot) => (
                <div key={slot.slot_id} className="border rounded p-3 flex items-center justify-between">
                  <div className="text-sm">
                    <div><b>{slot.teacher_name || slot.teacher_id}</b> - {slot.class_name || slot.class_id}</div>
                    <div className="text-muted-foreground">{new Date(slot.start_time).toLocaleString("vi-VN")} - {new Date(slot.end_time).toLocaleString("vi-VN")}</div>
                  </div>
                  <button className="px-3 py-1 rounded bg-primary text-primary-foreground text-xs" onClick={() => void bookSlot(slot.slot_id)}>
                    Đặt lịch
                  </button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4 space-y-3">
          <h2 className="font-semibold">Lịch đã đặt</h2>
          {appointments.length === 0 ? (
            <p className="text-sm text-muted-foreground">Bạn chưa có lịch hẹn nào.</p>
          ) : (
            <div className="space-y-2">
              {appointments.map((a) => (
                <div key={a.appointment_id} className="border rounded p-3 flex items-center justify-between gap-2">
                  <div className="text-sm">
                    <div><b>{a.student_name || a.student_id}</b> - {a.teacher_name || a.teacher_id}</div>
                    <div className="text-muted-foreground">{new Date(a.start_time).toLocaleString("vi-VN")} - {new Date(a.end_time).toLocaleString("vi-VN")}</div>
                    <div>Trạng thái: <b>{a.status}</b></div>
                  </div>
                  <div className="flex items-center gap-2">
                    <button className="border rounded px-3 py-1 text-xs" onClick={() => void syncStudentFromAppointment(a.student_id)}>
                      Chọn học sinh này
                    </button>
                    <button className="border rounded px-3 py-1 text-xs" onClick={() => void cancelAppointment(a.appointment_id)} disabled={a.status !== "pending" && a.status !== "confirmed"}>
                      Hủy lịch
                    </button>
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
