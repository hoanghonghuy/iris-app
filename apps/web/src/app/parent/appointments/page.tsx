"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Appointment, AppointmentSlot, ParentAnalytics, Student } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import {
  AlertTriangle,
  CalendarDays,
  CalendarClock,
  Download,
  Loader2,
  NotebookText,
  RefreshCcw,
  UserRound,
} from "lucide-react";
import { toast } from "sonner";

const APPOINTMENT_STATUS_CONFIG: Record<
  Appointment["status"],
  {
    label: string;
    variant: "default" | "secondary" | "outline" | "destructive";
  }
> = {
  pending: { label: "Chờ xác nhận", variant: "secondary" },
  confirmed: { label: "Đã xác nhận", variant: "default" },
  cancelled: { label: "Đã hủy", variant: "destructive" },
  completed: { label: "Hoàn tất", variant: "outline" },
  no_show: { label: "Vắng mặt", variant: "outline" },
};

const csvEscape = (value: unknown) => {
  const text = String(value ?? "");
  if (/[",\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`;
  }
  return text;
};

const downloadCsv = (filename: string, headers: string[], rows: Array<Array<unknown>>) => {
  const lines = [
    headers.map(csvEscape).join(","),
    ...rows.map((row) => row.map(csvEscape).join(",")),
  ];

  const blob = new Blob([`\uFEFF${lines.join("\n")}`], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};

const toDateInputValue = (date: Date) => {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000);
  return local.toISOString().slice(0, 10);
};

const getLocalDateKey = (value: string | Date) => {
  const date = value instanceof Date ? value : new Date(value);
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
};

const formatDayHeading = (dateKey: string) =>
  new Date(`${dateKey}T00:00:00`).toLocaleDateString("vi-VN", {
    weekday: "long",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });

const formatDateTime = (value: string) =>
  new Date(value).toLocaleString("vi-VN", {
    hour: "2-digit",
    minute: "2-digit",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    timeZoneName: "short",
  });

export default function ParentAppointmentsPage() {
  const [children, setChildren] = useState<Student[]>([]);
  const [selectedStudentId, setSelectedStudentId] = useState("");
  const [slots, setSlots] = useState<AppointmentSlot[]>([]);
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [stats, setStats] = useState<ParentAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingSlots, setLoadingSlots] = useState(false);
  const [submittingBooking, setSubmittingBooking] = useState(false);
  const [cancellingId, setCancellingId] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [bookingNote, setBookingNote] = useState("");
  const [historyFromDate, setHistoryFromDate] = useState(() => {
    const date = new Date();
    date.setDate(date.getDate() - 6);
    return toDateInputValue(date);
  });
  const [historyToDate, setHistoryToDate] = useState(() => toDateInputValue(new Date()));
  const [historyTab, setHistoryTab] = useState<"active" | "cancelled">("active");

  const timeZone = useMemo(() => Intl.DateTimeFormat().resolvedOptions().timeZone || "Local", []);
  const utcOffsetLabel = useMemo(() => {
    const totalMinutes = -new Date().getTimezoneOffset();
    const sign = totalMinutes >= 0 ? "+" : "-";
    const abs = Math.abs(totalMinutes);
    const hours = String(Math.floor(abs / 60)).padStart(2, "0");
    const mins = String(abs % 60).padStart(2, "0");
    return `UTC${sign}${hours}:${mins}`;
  }, []);
  const timezoneDisplay = `${timeZone} (${utcOffsetLabel})`;

  const resolveSyncedStudentId = useCallback(
    (
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
    },
    [selectedStudentId],
  );

  const load = useCallback(
    async (studentId?: string) => {
      setLoading(true);
      setErrorMessage(null);

      try {
        const from = historyFromDate ? new Date(`${historyFromDate}T00:00:00`).toISOString() : undefined;
        const to = historyToDate ? new Date(`${historyToDate}T23:59:59.999`).toISOString() : undefined;

        if (from && to && new Date(from).getTime() > new Date(to).getTime()) {
          setErrorMessage("Khoảng ngày lọc không hợp lệ: Từ ngày phải nhỏ hơn hoặc bằng Đến ngày.");
          setAppointments([]);
          return;
        }

        const [childrenData, appointmentsRes, analytics] = await Promise.all([
          parentApi.getMyChildren(),
          parentApi.getAppointments({ limit: 200, offset: 0, from, to }),
          parentApi.getAnalytics(),
        ]);

        const nextChildren = childrenData || [];
        const nextAppointments = appointmentsRes.data || [];

        setChildren(nextChildren);
        setAppointments(nextAppointments);
        setStats(analytics);

        const effectiveStudentId = resolveSyncedStudentId(
          nextChildren,
          nextAppointments,
          studentId,
        );

        if (!effectiveStudentId) {
          setSelectedStudentId("");
          setSlots([]);
          return;
        }

        setSelectedStudentId(effectiveStudentId);
        setLoadingSlots(true);

        const slotRes = await parentApi.getAvailableAppointmentSlots({
          student_id: effectiveStudentId,
          limit: 50,
          offset: 0,
        });
        setSlots(slotRes.data || []);
      } catch (error) {
        setErrorMessage(extractApiErrorMessage(error, "Không thể tải dữ liệu lịch hẹn."));
      } finally {
        setLoading(false);
        setLoadingSlots(false);
      }
    },
    [historyFromDate, historyToDate, resolveSyncedStudentId],
  );

  useEffect(() => {
    void load();
  }, [load]);

  const onSelectStudent = async (studentId: string) => {
    setSelectedStudentId(studentId);
    setSlots([]);

    if (!studentId) {
      return;
    }

    setLoadingSlots(true);
    try {
      const slotRes = await parentApi.getAvailableAppointmentSlots({
        student_id: studentId,
        limit: 50,
        offset: 0,
      });
      setSlots(slotRes.data || []);
    } catch (error) {
      toast.error(extractApiErrorMessage(error, "Không thể tải khung giờ khả dụng."));
    } finally {
      setLoadingSlots(false);
    }
  };

  const syncStudentFromAppointment = async (studentId: string) => {
    await onSelectStudent(studentId);
  };

  const bookSlot = async (slotId: string) => {
    if (!selectedStudentId) return;

    setSubmittingBooking(true);
    try {
      await parentApi.createAppointment({
        student_id: selectedStudentId,
        slot_id: slotId,
        note: bookingNote.trim() || undefined,
      });
      setBookingNote("");
      toast.success("Đặt lịch thành công.");
      await load(selectedStudentId);
    } catch (error) {
      toast.error(extractApiErrorMessage(error, "Không thể đặt lịch. Vui lòng thử lại."));
    } finally {
      setSubmittingBooking(false);
    }
  };

  const cancelAppointment = async (appointmentId: string) => {
    setCancellingId(appointmentId);
    try {
      await parentApi.cancelAppointment(appointmentId, "parent_cancelled");
      toast.success("Đã hủy lịch hẹn.");
      await load(selectedStudentId);
    } catch (error) {
      toast.error(extractApiErrorMessage(error, "Không thể hủy lịch hẹn."));
    } finally {
      setCancellingId(null);
    }
  };

  const activeAppointmentsCount = appointments.filter(
    (a) => a.status === "pending" || a.status === "confirmed",
  ).length;

  const historyAppointments = useMemo(
    () => appointments.filter((a) => (historyTab === "cancelled" ? a.status === "cancelled" : a.status !== "cancelled")),
    [appointments, historyTab],
  );

  const groupedHistoryAppointments = useMemo(() => {
    const sorted = [...historyAppointments].sort(
      (a, b) => new Date(b.start_time).getTime() - new Date(a.start_time).getTime(),
    );

    const groups = new Map<string, Appointment[]>();
    for (const appt of sorted) {
      const key = getLocalDateKey(appt.start_time);
      const existing = groups.get(key);
      if (existing) {
        existing.push(appt);
      } else {
        groups.set(key, [appt]);
      }
    }

    return Array.from(groups.entries()).map(([dateKey, items]) => ({ dateKey, items }));
  }, [historyAppointments]);

  const exportHistoryCsv = () => {
    if (historyAppointments.length === 0) {
      toast.error("Không có dữ liệu để xuất CSV.");
      return;
    }

    const rows = [...historyAppointments]
      .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime())
      .map((a) => [
        formatDayHeading(getLocalDateKey(a.start_time)),
        a.student_name || a.student_id,
        a.teacher_name || a.teacher_id,
        APPOINTMENT_STATUS_CONFIG[a.status].label,
        formatDateTime(a.start_time),
        formatDateTime(a.end_time),
        timezoneDisplay,
        a.note || "",
      ]);

    downloadCsv(
      `parent-appointments-${historyTab}-${toDateInputValue(new Date())}.csv`,
      ["Ngay", "HocSinh", "GiaoVien", "TrangThai", "BatDau", "KetThuc", "MuiGio", "GhiChu"],
      rows,
    );

    toast.success("Đã xuất CSV theo bộ lọc hiện tại.");
  };

  return (
    <div className="space-y-6 pb-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-bold">Lịch hẹn với giáo viên</h1>
        <p className="text-sm text-muted-foreground">
          Chọn học sinh, xem khung giờ còn trống và đặt lịch ngay trong một màn hình.
        </p>
        <p className="text-xs text-muted-foreground">Múi giờ hiển thị: {timezoneDisplay}</p>
      </div>

      {errorMessage && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertTitle>Tải dữ liệu thất bại</AlertTitle>
          <AlertDescription className="w-full">
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
              <span>{errorMessage}</span>
              <Button variant="outline" size="sm" onClick={() => void load(selectedStudentId)}>
                Thử lại
              </Button>
            </div>
          </AlertDescription>
        </Alert>
      )}

      <Card>
        <CardContent className="grid gap-2 p-4 text-sm md:grid-cols-3">
          <div><b>Tổng số con:</b> {stats?.total_children ?? 0}</div>
          <div><b>Lịch đang hoạt động:</b> {activeAppointmentsCount}</div>
          <div><b>Cảnh báo sức khỏe 24h:</b> {stats?.recent_health_alerts_24h ?? 0}</div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-4 p-4 md:p-5">
          <div className="space-y-1">
            <h2 className="font-semibold">Đặt lịch mới</h2>
            <p className="text-sm text-muted-foreground">
              Các khung giờ bên phải được lọc theo học sinh bạn đang chọn.
            </p>
          </div>

          <div className="grid gap-4 lg:grid-cols-[320px_minmax(0,1fr)]">
            <div className="space-y-3">
              <div className="space-y-1.5">
                <Label htmlFor="parent-appointment-student">Chọn học sinh</Label>
                <Select
                  value={selectedStudentId || undefined}
                  onValueChange={(value) => void onSelectStudent(value)}
                >
                  <SelectTrigger id="parent-appointment-student" className="w-full">
                    <SelectValue placeholder={children.length > 0 ? "Chọn học sinh" : "Chưa có học sinh"} />
                  </SelectTrigger>
                  <SelectContent>
                    {children.map((student) => (
                      <SelectItem key={student.student_id} value={student.student_id}>
                        {student.full_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="parent-appointment-note">Ghi chú khi đặt lịch</Label>
                <Textarea
                  id="parent-appointment-note"
                  placeholder="Ví dụ: Mong muốn trao đổi về tình hình học tập tuần này..."
                  value={bookingNote}
                  onChange={(e) => setBookingNote(e.target.value)}
                  maxLength={500}
                  rows={4}
                />
                <p className="text-xs text-muted-foreground">Tối đa 500 ký tự.</p>
              </div>
            </div>

            <div className="space-y-3">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <h3 className="text-sm font-semibold">Khung giờ khả dụng</h3>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  onClick={() => void load(selectedStudentId)}
                  disabled={loading || loadingSlots || !selectedStudentId}
                >
                  <RefreshCcw className="h-3.5 w-3.5" />
                  Làm mới
                </Button>
              </div>

              {!selectedStudentId ? (
                <Alert>
                  <CalendarClock className="h-4 w-4" />
                  <AlertTitle>Chưa chọn học sinh</AlertTitle>
                  <AlertDescription>Vui lòng chọn học sinh để xem danh sách khung giờ có thể đặt.</AlertDescription>
                </Alert>
              ) : loading || loadingSlots ? (
                <div className="flex min-h-[180px] items-center justify-center rounded-lg border border-dashed">
                  <Loader2 className="h-6 w-6 animate-spin text-primary" />
                </div>
              ) : slots.length === 0 ? (
                <div className="rounded-lg border border-dashed p-4 text-sm text-muted-foreground">
                  Hiện chưa có khung giờ phù hợp cho học sinh này.
                </div>
              ) : (
                <div className="max-h-[360px] space-y-2 overflow-y-auto pr-1">
                  {slots.map((slot) => (
                    <div key={slot.slot_id} className="rounded-lg border p-3">
                      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                        <div className="space-y-1 text-sm">
                          <div className="flex items-center gap-1.5 font-medium">
                            <UserRound className="h-4 w-4 text-muted-foreground" />
                            {slot.teacher_name || slot.teacher_id}
                          </div>
                          <div className="text-muted-foreground">
                            {slot.class_name || slot.class_id}
                          </div>
                          <div className="text-muted-foreground">
                            {formatDateTime(slot.start_time)} - {formatDateTime(slot.end_time)}
                          </div>
                          <div className="text-xs text-muted-foreground">Múi giờ: {timezoneDisplay}</div>
                          {slot.note && (
                            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                              <NotebookText className="h-3.5 w-3.5" />
                              <span>{slot.note}</span>
                            </div>
                          )}
                        </div>
                        <Button
                          type="button"
                          size="sm"
                          onClick={() => void bookSlot(slot.slot_id)}
                          disabled={submittingBooking}
                        >
                          {submittingBooking ? "Đang đặt..." : "Đặt lịch"}
                        </Button>
                      </div>
                    </div>
                  ))}
                  </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-3 p-4 md:p-5">
          <div className="space-y-2">
            <h2 className="font-semibold">Lịch đã đặt</h2>
            <div className="grid gap-2 md:grid-cols-[auto_auto_1fr]">
              <Button
                type="button"
                size="sm"
                variant={historyTab === "active" ? "default" : "outline"}
                onClick={() => setHistoryTab("active")}
              >
                Đang hoạt động
              </Button>
              <Button
                type="button"
                size="sm"
                variant={historyTab === "cancelled" ? "default" : "outline"}
                onClick={() => setHistoryTab("cancelled")}
              >
                Đã hủy
              </Button>
            </div>

            <div className="grid gap-2 md:grid-cols-[1fr_1fr_auto_auto]">
              <div className="flex items-center gap-2">
                <Input
                  type="date"
                  value={historyFromDate}
                  onChange={(e) => setHistoryFromDate(e.target.value)}
                />
              </div>
              <div className="flex items-center gap-2">
                <Input
                  type="date"
                  value={historyToDate}
                  onChange={(e) => setHistoryToDate(e.target.value)}
                />
              </div>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => {
                  const to = toDateInputValue(new Date());
                  const fromDate = new Date();
                  fromDate.setDate(fromDate.getDate() - 6);
                  setHistoryFromDate(toDateInputValue(fromDate));
                  setHistoryToDate(to);
                }}
              >
                7 ngày gần nhất
              </Button>
              <Button type="button" variant="outline" size="sm" onClick={exportHistoryCsv}>
                <Download className="h-3.5 w-3.5" />
                Xuất CSV
              </Button>
            </div>
          </div>

          {loading ? (
            <div className="py-8 flex justify-center"><Loader2 className="h-6 w-6 animate-spin" /></div>
          ) : groupedHistoryAppointments.length === 0 ? (
            <p className="text-sm text-muted-foreground">Bạn chưa có lịch hẹn nào.</p>
          ) : (
            <div className="max-h-[460px] space-y-4 overflow-y-auto pr-1">
              {groupedHistoryAppointments.map((group) => (
                <div key={group.dateKey} className="space-y-2">
                  <div className="sticky top-0 z-10 flex items-center gap-2 rounded-md bg-background/95 px-1 py-1 text-xs font-semibold text-muted-foreground backdrop-blur">
                    <CalendarDays className="h-3.5 w-3.5" />
                    {formatDayHeading(group.dateKey)}
                  </div>

                  {group.items.map((a) => (
                    <div key={a.appointment_id} className="rounded-lg border p-3">
                      <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                        <div className="space-y-1 text-sm">
                          <div className="font-medium">{a.student_name || a.student_id} - {a.teacher_name || a.teacher_id}</div>
                          <div className="text-muted-foreground">{formatDateTime(a.start_time)} - {formatDateTime(a.end_time)}</div>
                          <div className="text-xs text-muted-foreground">Múi giờ: {timezoneDisplay}</div>
                          {a.note && <div className="text-muted-foreground">Ghi chú: {a.note}</div>}
                          <div className="pt-1">
                            <Badge variant={APPOINTMENT_STATUS_CONFIG[a.status].variant}>
                              {APPOINTMENT_STATUS_CONFIG[a.status].label}
                            </Badge>
                          </div>
                        </div>

                        <div className="flex flex-wrap items-center gap-2">
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => void syncStudentFromAppointment(a.student_id)}
                          >
                            Chọn học sinh
                          </Button>

                          <Button
                            type="button"
                            variant="destructive"
                            size="sm"
                            onClick={() => void cancelAppointment(a.appointment_id)}
                            disabled={(a.status !== "pending" && a.status !== "confirmed") || cancellingId === a.appointment_id}
                          >
                            {cancellingId === a.appointment_id ? "Đang hủy..." : "Hủy lịch"}
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
