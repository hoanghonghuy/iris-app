"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Appointment, AppointmentStatus, Class } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
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
  Clock3,
  Download,
  Loader2,
  RefreshCcw,
  UserRound,
} from "lucide-react";
import { toast } from "sonner";

const APPOINTMENT_STATUS_CONFIG: Record<
  AppointmentStatus,
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

const APPOINTMENT_STATUS_OPTIONS: Array<{ value: AppointmentStatus; label: string }> = [
  { value: "pending", label: "Chờ xác nhận" },
  { value: "confirmed", label: "Đã xác nhận" },
  { value: "cancelled", label: "Đã hủy" },
  { value: "completed", label: "Hoàn tất" },
  { value: "no_show", label: "Vắng mặt" },
];

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

export default function TeacherAppointmentsPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [loading, setLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [updatingAppointmentId, setUpdatingAppointmentId] = useState<string | null>(null);
  const [createModalOpen, setCreateModalOpen] = useState(false);

  const [classId, setClassId] = useState("");
  const [startTime, setStartTime] = useState("");
  const [durationMinutes, setDurationMinutes] = useState(30);
  const [bufferMinutes, setBufferMinutes] = useState(10);
  const [maxBookingsPerDay, setMaxBookingsPerDay] = useState(12);
  const [note, setNote] = useState("");

  const [statusFilter, setStatusFilter] = useState<"" | AppointmentStatus>("");
  const [filterFromDate, setFilterFromDate] = useState(() => {
    const date = new Date();
    date.setDate(date.getDate() - 6);
    return toDateInputValue(date);
  });
  const [filterToDate, setFilterToDate] = useState(() => toDateInputValue(new Date()));

  const availableClassOptions = useMemo(() => classes ?? [], [classes]);
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
  const minStartTime = new Date(Date.now() - new Date().getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16);

  const stats = useMemo(
    () => ({
      totalClasses: availableClassOptions.length,
      totalAppointments: appointments.length,
      pendingCount: appointments.filter((a) => a.status === "pending").length,
      confirmedCount: appointments.filter((a) => a.status === "confirmed").length,
    }),
    [appointments, availableClassOptions.length],
  );

  const groupedAppointments = useMemo(() => {
    const sorted = [...appointments].sort(
      (a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime(),
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
  }, [appointments]);

  const loadData = useCallback(async () => {
    setLoading(true);
    setErrorMessage(null);

    try {
      const from = filterFromDate ? new Date(`${filterFromDate}T00:00:00`).toISOString() : undefined;
      const to = filterToDate ? new Date(`${filterToDate}T23:59:59.999`).toISOString() : undefined;

      if (from && to && new Date(from).getTime() > new Date(to).getTime()) {
        setErrorMessage("Khoảng ngày lọc không hợp lệ: Từ ngày phải nhỏ hơn hoặc bằng Đến ngày.");
        setAppointments([]);
        return;
      }

      const [classData, appointmentRes] = await Promise.all([
        teacherApi.getMyClasses(),
        teacherApi.getAppointments({
          limit: 200,
          offset: 0,
          status: statusFilter || undefined,
          from,
          to,
        }),
      ]);

      setClasses(classData || []);
      setAppointments(appointmentRes.data || []);

      if (!classId && classData?.length) {
        setClassId(classData[0].class_id);
      }
    } catch (error) {
      setErrorMessage(extractApiErrorMessage(error, "Không thể tải dữ liệu lịch hẹn."));
    } finally {
      setLoading(false);
    }
  }, [classId, filterFromDate, filterToDate, statusFilter]);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  const createSlot = async () => {
    if (!classId || !startTime) return;

    setSubmitting(true);
    try {
      const startDate = new Date(startTime);
      if (Number.isNaN(startDate.getTime())) {
        toast.error("Thời gian bắt đầu không hợp lệ.");
        return;
      }

      const proposedStartMs = startDate.getTime();
      const proposedEndMs = proposedStartMs + durationMinutes * 60000;
      const bufferMs = Math.max(0, bufferMinutes) * 60000;

      const dayStart = new Date(startDate);
      dayStart.setHours(0, 0, 0, 0);
      const dayEnd = new Date(startDate);
      dayEnd.setHours(23, 59, 59, 999);

      const dayAppointmentsRes = await teacherApi.getAppointments({
        limit: 200,
        offset: 0,
        from: dayStart.toISOString(),
        to: dayEnd.toISOString(),
      });

      const dayAppointments = (dayAppointmentsRes.data || []).filter((a) => a.status !== "cancelled");

      if (dayAppointments.length >= maxBookingsPerDay) {
        toast.error(`Đã đạt giới hạn ${maxBookingsPerDay} lịch trong ngày này.`);
        return;
      }

      const conflicting = dayAppointments.find((a) => {
        const existingStart = new Date(a.start_time).getTime();
        const existingEnd = new Date(a.end_time).getTime();

        const hasSpacing =
          proposedEndMs + bufferMs <= existingStart || proposedStartMs >= existingEnd + bufferMs;

        return !hasSpacing;
      });

      if (conflicting) {
        toast.error(
          `Khung giờ mới chưa đảm bảo khoảng nghỉ ${bufferMinutes} phút với lịch ${formatDateTime(conflicting.start_time)}.`,
        );
        return;
      }

      const startISO = startDate.toISOString();

      await teacherApi.createAppointmentSlot({
        class_id: classId,
        start_time: startISO,
        duration_minutes: durationMinutes,
        buffer_minutes: bufferMinutes,
        max_bookings_per_day: maxBookingsPerDay,
        note: note.trim() || undefined,
      });
      toast.success("Tạo khung giờ thành công.");
      setStartTime("");
      setNote("");
      setCreateModalOpen(false);
      await loadData();
    } catch (error) {
      toast.error(extractApiErrorMessage(error, "Không thể tạo khung giờ."));
    } finally {
      setSubmitting(false);
    }
  };

  const updateStatus = async (appointmentId: string, status: AppointmentStatus) => {
    setUpdatingAppointmentId(appointmentId);

    try {
      await teacherApi.updateAppointmentStatus(
        appointmentId,
        status,
        status === "cancelled" ? "teacher_cancelled" : undefined,
      );
      toast.success("Cập nhật trạng thái lịch hẹn thành công.");
      await loadData();
    } catch (error) {
      toast.error(extractApiErrorMessage(error, "Không thể cập nhật trạng thái lịch hẹn."));
    } finally {
      setUpdatingAppointmentId(null);
    }
  };

  const exportAppointmentsCsv = () => {
    if (appointments.length === 0) {
      toast.error("Không có dữ liệu để xuất CSV.");
      return;
    }

    const rows = [...appointments]
      .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime())
      .map((a) => [
        formatDayHeading(getLocalDateKey(a.start_time)),
        a.student_name || a.student_id,
        a.class_name || a.class_id,
        a.parent_name || a.parent_id,
        APPOINTMENT_STATUS_CONFIG[a.status].label,
        formatDateTime(a.start_time),
        formatDateTime(a.end_time),
        timezoneDisplay,
        a.note || "",
      ]);

    downloadCsv(
      `teacher-appointments-${toDateInputValue(new Date())}.csv`,
      [
        "Ngay",
        "HocSinh",
        "Lop",
        "PhuHuynh",
        "TrangThai",
        "BatDau",
        "KetThuc",
        "MuiGio",
        "GhiChu",
      ],
      rows,
    );

    toast.success("Đã xuất CSV theo bộ lọc hiện tại.");
  };

  return (
    <div className="space-y-6 pb-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-bold">Lịch hẹn phụ huynh</h1>
        <p className="text-sm text-muted-foreground">
          Tạo khung giờ, theo dõi yêu cầu đặt lịch và cập nhật trạng thái trong một luồng làm việc.
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
              <Button type="button" variant="outline" size="sm" onClick={() => void loadData()}>
                Thử lại
              </Button>
            </div>
          </AlertDescription>
        </Alert>
      )}

      <Card>
        <CardContent className="grid gap-2 p-4 text-sm md:grid-cols-4">
          <div><b>Số lớp đang phụ trách:</b> {stats.totalClasses}</div>
          <div><b>Tổng lịch hẹn:</b> {stats.totalAppointments}</div>
          <div><b>Đang chờ xác nhận:</b> {stats.pendingCount}</div>
          <div><b>Đã xác nhận:</b> {stats.confirmedCount}</div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="font-semibold">Tạo khung giờ mới</h2>
            <p className="text-sm text-muted-foreground">
              Mở biểu mẫu trong modal để thêm khung giờ mà không làm rối màn danh sách.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Dialog open={createModalOpen} onOpenChange={setCreateModalOpen}>
              <DialogTrigger asChild>
                <Button type="button">Tạo khung giờ mới</Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                  <DialogTitle>Tạo khung giờ mới</DialogTitle>
                  <DialogDescription>
                    Điền thời gian, khoảng nghỉ giữa hai lịch và số lịch tối đa/ngày để hệ thống kiểm tra trước khi tạo.
                  </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                  <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-4">
                    <div className="space-y-1.5">
                      <Label htmlFor="teacher-appointment-class">Lớp học</Label>
                      <Select value={classId || undefined} onValueChange={setClassId}>
                        <SelectTrigger id="teacher-appointment-class" className="w-full">
                          <SelectValue placeholder={availableClassOptions.length > 0 ? "Chọn lớp" : "Chưa có lớp"} />
                        </SelectTrigger>
                        <SelectContent>
                          {availableClassOptions.map((c) => (
                            <SelectItem key={c.class_id} value={c.class_id}>
                              {c.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div className="space-y-1.5">
                      <Label htmlFor="teacher-appointment-start-time">Thời gian bắt đầu</Label>
                      <Input
                        id="teacher-appointment-start-time"
                        type="datetime-local"
                        value={startTime}
                        min={minStartTime}
                        onChange={(e) => setStartTime(e.target.value)}
                      />
                    </div>

                    <div className="space-y-1.5">
                      <Label htmlFor="teacher-appointment-duration">Thời lượng</Label>
                      <Select
                        value={String(durationMinutes)}
                        onValueChange={(value) => setDurationMinutes(Number(value))}
                      >
                        <SelectTrigger id="teacher-appointment-duration" className="w-full">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="15">15 phút</SelectItem>
                          <SelectItem value="20">20 phút</SelectItem>
                          <SelectItem value="30">30 phút</SelectItem>
                          <SelectItem value="45">45 phút</SelectItem>
                          <SelectItem value="60">60 phút</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div className="space-y-1.5">
                      <Label htmlFor="teacher-appointment-end-time">Kết thúc dự kiến</Label>
                      <div className="flex h-9 items-center rounded-md border bg-muted/40 px-3 text-sm text-muted-foreground">
                        <Clock3 className="mr-2 h-4 w-4" />
                        {startTime
                          ? formatDateTime(new Date(new Date(startTime).getTime() + durationMinutes * 60000).toISOString())
                          : "Chưa xác định"}
                      </div>
                    </div>
                  </div>

                  <div className="grid gap-3 md:grid-cols-2">
                    <div className="space-y-1.5">
                      <Label htmlFor="teacher-buffer-minutes">Khoảng nghỉ giữa 2 lịch (phút)</Label>
                      <Input
                        id="teacher-buffer-minutes"
                        type="number"
                        min={0}
                        step={5}
                        value={bufferMinutes}
                        onChange={(e) => setBufferMinutes(Math.max(0, Number(e.target.value || 0)))}
                      />
                      <p className="text-xs text-muted-foreground">
                        Hệ thống sẽ giữ khoảng nghỉ này trước hoặc sau mỗi lịch để tránh đặt quá sát nhau.
                      </p>
                    </div>

                    <div className="space-y-1.5">
                      <Label htmlFor="teacher-max-bookings-per-day">Số lịch tối đa trong ngày</Label>
                      <Input
                        id="teacher-max-bookings-per-day"
                        type="number"
                        min={1}
                        step={1}
                        value={maxBookingsPerDay}
                        onChange={(e) => setMaxBookingsPerDay(Math.max(1, Number(e.target.value || 1)))}
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="teacher-appointment-note">Ghi chú cho phụ huynh (tuỳ chọn)</Label>
                    <Textarea
                      id="teacher-appointment-note"
                      placeholder="Ví dụ: Vui lòng chuẩn bị các câu hỏi liên quan tới tiến độ học tập của bé..."
                      value={note}
                      onChange={(e) => setNote(e.target.value)}
                      maxLength={500}
                      rows={3}
                    />
                    <p className="text-xs text-muted-foreground">Tối đa 500 ký tự.</p>
                  </div>
                </div>

                <DialogFooter>
                  <Button type="button" variant="outline" onClick={() => setCreateModalOpen(false)}>
                    Đóng
                  </Button>
                  <Button
                    type="button"
                    onClick={createSlot}
                    disabled={submitting || !classId || !startTime}
                  >
                    {submitting ? "Đang tạo..." : "Tạo khung giờ"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
            <Button type="button" variant="outline" onClick={() => void loadData()} disabled={loading}>
              <RefreshCcw className="h-4 w-4" />
              Làm mới
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-3 p-4 md:p-5">
          <div className="flex flex-col gap-2">
            <h2 className="font-semibold">Danh sách lịch hẹn</h2>
            <div className="grid gap-2 md:grid-cols-4">
              <Select
                value={statusFilter || "all"}
                onValueChange={(value) => setStatusFilter(value === "all" ? "" : (value as AppointmentStatus))}
              >
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Tất cả trạng thái</SelectItem>
                  {APPOINTMENT_STATUS_OPTIONS.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              <div className="flex items-center gap-2">
                <Input
                  type="date"
                  value={filterFromDate}
                  onChange={(e) => setFilterFromDate(e.target.value)}
                />
                <span className="text-xs text-muted-foreground">đến</span>
                <Input
                  type="date"
                  value={filterToDate}
                  onChange={(e) => setFilterToDate(e.target.value)}
                />
              </div>

              <div className="flex items-center gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    const to = toDateInputValue(new Date());
                    const fromDate = new Date();
                    fromDate.setDate(fromDate.getDate() - 6);
                    setFilterFromDate(toDateInputValue(fromDate));
                    setFilterToDate(to);
                  }}
                >
                  7 ngày gần nhất
                </Button>
                <Button type="button" variant="outline" size="sm" onClick={exportAppointmentsCsv}>
                  <Download className="h-3.5 w-3.5" />
                  Xuất CSV
                </Button>
              </div>
            </div>
          </div>

          {loading ? (
            <div className="py-8 flex justify-center"><Loader2 className="h-6 w-6 animate-spin" /></div>
          ) : appointments.length === 0 ? (
            <p className="text-sm text-muted-foreground">Chưa có lịch hẹn.</p>
          ) : (
            <div className="max-h-[560px] space-y-4 overflow-y-auto pr-1">
              {groupedAppointments.map((group) => (
                <div key={group.dateKey} className="space-y-2">
                  <div className="sticky top-0 z-10 flex items-center gap-2 rounded-md bg-background/95 px-1 py-1 text-xs font-semibold text-muted-foreground backdrop-blur">
                    <CalendarDays className="h-3.5 w-3.5" />
                    {formatDayHeading(group.dateKey)}
                  </div>

                  {group.items.map((a) => (
                    <div key={a.appointment_id} className="rounded-lg border p-3">
                      <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
                        <div className="space-y-1 text-sm">
                          <div className="font-medium">
                            {a.student_name || a.student_id} - {a.class_name || a.class_id}
                          </div>
                          <div className="text-muted-foreground">
                            {formatDateTime(a.start_time)} - {formatDateTime(a.end_time)}
                          </div>
                          <div className="text-xs text-muted-foreground">Múi giờ: {timezoneDisplay}</div>
                          <div className="flex items-center gap-1.5 text-muted-foreground">
                            <UserRound className="h-3.5 w-3.5" />
                            Phụ huynh: {a.parent_name || a.parent_id}
                          </div>
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
                            size="sm"
                            onClick={() => void updateStatus(a.appointment_id, "confirmed")}
                            disabled={a.status !== "pending" || updatingAppointmentId === a.appointment_id}
                          >
                            Xác nhận
                          </Button>

                          <Button
                            type="button"
                            size="sm"
                            variant="outline"
                            onClick={() => void updateStatus(a.appointment_id, "completed")}
                            disabled={a.status !== "confirmed" || updatingAppointmentId === a.appointment_id}
                          >
                            Hoàn tất
                          </Button>

                          <Button
                            type="button"
                            size="sm"
                            variant="outline"
                            onClick={() => void updateStatus(a.appointment_id, "no_show")}
                            disabled={a.status !== "confirmed" || updatingAppointmentId === a.appointment_id}
                          >
                            Vắng mặt
                          </Button>

                          <Button
                            type="button"
                            size="sm"
                            variant="destructive"
                            onClick={() => void updateStatus(a.appointment_id, "cancelled")}
                            disabled={(a.status === "cancelled" || a.status === "completed") || updatingAppointmentId === a.appointment_id}
                          >
                            {updatingAppointmentId === a.appointment_id ? "Đang xử lý..." : "Hủy"}
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
