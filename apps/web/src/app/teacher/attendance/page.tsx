/**
 * Teacher Attendance Page
 * Chọn lớp → xem HS → điểm danh từng em.
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/attendance
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/attendance
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, Student, AttendanceStatus, AttendanceChangeLog } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ClipboardCheck, Loader2, Check, AlertCircle, History } from "lucide-react";
import { formatDateVN } from "@/lib/utils";

const statusOptions = [
  { value: "present" as AttendanceStatus, label: "Có mặt", variant: "default" as const },
  { value: "absent" as AttendanceStatus, label: "Vắng", variant: "destructive" as const },
  { value: "late" as AttendanceStatus, label: "Muộn", variant: "secondary" as const },
  { value: "excused" as AttendanceStatus, label: "Có phép", variant: "outline" as const },
];

const statusLabel: Record<AttendanceStatus, string> = {
  present: "Có mặt",
  absent: "Vắng",
  late: "Muộn",
  excused: "Có phép",
};

type AttendanceChangeHistoryRow = AttendanceChangeLog & {
  student_name: string;
};

type TakeListFilter = "all" | "pending" | "saved";

function extractErrorMessage(err: unknown): string | undefined {
  return typeof err === "object" &&
    err !== null &&
    "response" in err &&
    typeof (err as { response?: { data?: { error?: string } } }).response?.data?.error === "string"
    ? (err as { response?: { data?: { error?: string } } }).response?.data?.error
    : undefined;
}

function mapStatusLabel(status?: AttendanceStatus): string {
  if (!status) {
    return "-";
  }
  return statusLabel[status] || status;
}

export default function TeacherAttendancePage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  const [attendance, setAttendance] = useState<Record<string, { status: AttendanceStatus; note: string }>>({});
  const [savedAttendance, setSavedAttendance] = useState<Record<string, { status: AttendanceStatus; note: string }>>({});
  const [hasSavedToday, setHasSavedToday] = useState<Record<string, boolean>>({});

  const [submitting, setSubmitting] = useState<string | null>(null);
  const [canceling, setCanceling] = useState<string | null>(null);
  const [savingAll, setSavingAll] = useState(false);
  const [savingDisplayed, setSavingDisplayed] = useState(false);

  const [historyOpen, setHistoryOpen] = useState<Set<string>>(new Set());
  const [historyLoading, setHistoryLoading] = useState<Set<string>>(new Set());
  const [historyByStudent, setHistoryByStudent] = useState<Record<string, AttendanceChangeLog[]>>({});
  const [studentSearch, setStudentSearch] = useState("");
  const [listOrderMode, setListOrderMode] = useState<"prioritize" | "original">("prioritize");
  const [takeListFilter, setTakeListFilter] = useState<TakeListFilter>("all");
  const [showMobileTakeControls, setShowMobileTakeControls] = useState(false);

  const [today] = useState(() => new Date().toISOString().slice(0, 10));

  const [viewMode, setViewMode] = useState<"take" | "history">("take");
  const [historyFrom, setHistoryFrom] = useState(() => {
    const from = new Date();
    from.setDate(from.getDate() - 7);
    return from.toISOString().slice(0, 10);
  });
  const [historyTo, setHistoryTo] = useState(today);
  const [historyStudentId, setHistoryStudentId] = useState("all");
  const [historyStatus, setHistoryStatus] = useState<AttendanceStatus | "all">("all");
  const [historyListLoading, setHistoryListLoading] = useState(false);
  const [historyList, setHistoryList] = useState<AttendanceChangeHistoryRow[]>([]);
  const [historyOffset, setHistoryOffset] = useState(0);
  const [historyLimit] = useState(20);
  const [historyTotal, setHistoryTotal] = useState(0);
  const [historyHasMore, setHistoryHasMore] = useState(false);

  const isRowDirty = useCallback((studentId: string) => {
    const current = attendance[studentId];
    const saved = savedAttendance[studentId];
    if (!current || !saved) {
      return !hasSavedToday[studentId];
    }
    if (!hasSavedToday[studentId]) {
      return true;
    }
    return current.status !== saved.status || (current.note || "") !== (saved.note || "");
  }, [attendance, savedAttendance, hasSavedToday]);

  const dirtyCount = students.filter((student) => isRowDirty(student.student_id)).length;

  const displayedStudentsBase = useMemo(() => {
    const normalizedSearch = studentSearch.trim().toLowerCase();
    const searched = normalizedSearch
      ? students.filter((student) => student.full_name.toLowerCase().includes(normalizedSearch))
      : students;

    if (takeListFilter === "pending") {
      return searched.filter((student) => isRowDirty(student.student_id));
    }
    if (takeListFilter === "saved") {
      return searched.filter((student) => !isRowDirty(student.student_id));
    }
    return searched;
  }, [students, studentSearch, takeListFilter, isRowDirty]);

  const displayedStudents = useMemo(() => {
    if (listOrderMode === "original") {
      return displayedStudentsBase;
    }

    const unfinished = displayedStudentsBase.filter((student) => isRowDirty(student.student_id));
    const finished = displayedStudentsBase.filter((student) => !isRowDirty(student.student_id));
    return [...unfinished, ...finished];
  }, [displayedStudentsBase, isRowDirty, listOrderMode]);

  const displayedDirtyCount = displayedStudents.filter((student) => isRowDirty(student.student_id)).length;
  const displayedSavedCount = displayedStudents.length - displayedDirtyCount;
  const globalPendingCount = students.length - students.filter((student) => !isRowDirty(student.student_id)).length;

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
      setLoadingStudents(true);
      setError("");
      setHistoryOpen(new Set());
      setHistoryByStudent({});

      const data = await teacherApi.getStudentsInClass(selectedClassId);

      const studentList = data || [];
      setStudents(studentList);

      const init: Record<string, { status: AttendanceStatus; note: string }> = {};
      const savedInit: Record<string, { status: AttendanceStatus; note: string }> = {};
      const hasSavedInit: Record<string, boolean> = {};

      await Promise.all(
        studentList.map(async (student: Student) => {
          try {
            const records = await teacherApi.getStudentAttendance(student.student_id, today, today);
            const todayRecord = records.find((record) => record.date.startsWith(today));

            if (todayRecord) {
              const value = {
                status: todayRecord.status,
                note: todayRecord.note || "",
              };
              init[student.student_id] = value;
              savedInit[student.student_id] = value;
              hasSavedInit[student.student_id] = true;
            } else {
              const defaultValue = { status: "present" as AttendanceStatus, note: "" };
              init[student.student_id] = defaultValue;
              savedInit[student.student_id] = defaultValue;
              hasSavedInit[student.student_id] = false;
            }
          } catch {
            const defaultValue = { status: "present" as AttendanceStatus, note: "" };
            init[student.student_id] = defaultValue;
            savedInit[student.student_id] = defaultValue;
            hasSavedInit[student.student_id] = false;
          }
        })
      );

      setAttendance(init);
      setSavedAttendance(savedInit);
      setHasSavedToday(hasSavedInit);
    } catch (err: unknown) {
      const message = extractErrorMessage(err);
      setError(message || "Không thể tải HS");
    } finally { setLoadingStudents(false); }
  }, [selectedClassId, today]);

  useEffect(() => { fetchStudents(); }, [fetchStudents]);

  useEffect(() => {
    setHistoryStudentId("all");
    setHistoryList([]);
    setHistoryOffset(0);
    setHistoryTotal(0);
    setHistoryHasMore(false);
  }, [selectedClassId]);

  const handleMark = async (studentId: string) => {
    const att = attendance[studentId];
    if (!att) return;
    try {
      setSubmitting(studentId);
      await teacherApi.markAttendance({ student_id: studentId, date: today, status: att.status, note: att.note });
      setSavedAttendance((prev) => ({ ...prev, [studentId]: { status: att.status, note: att.note } }));
      setHasSavedToday((prev) => ({ ...prev, [studentId]: true }));
    } catch (err: unknown) {
      const message = extractErrorMessage(err);
      setError(message || "Lỗi điểm danh");
    } finally { setSubmitting(null); }
  };

  const handleRevertLocal = (studentId: string) => {
    const saved = savedAttendance[studentId];
    if (!saved) {
      return;
    }
    setAttendance((prev) => ({
      ...prev,
      [studentId]: { status: saved.status, note: saved.note },
    }));
  };

  const handleCancelSaved = async (studentId: string) => {
    if (!hasSavedToday[studentId]) {
      return;
    }

    setCanceling(studentId);
    setError("");
    try {
      await teacherApi.cancelAttendance(studentId, today);

      const defaultValue = { status: "present" as AttendanceStatus, note: "" };
      setAttendance((prev) => ({ ...prev, [studentId]: defaultValue }));
      setSavedAttendance((prev) => ({ ...prev, [studentId]: defaultValue }));
      setHasSavedToday((prev) => ({ ...prev, [studentId]: false }));
      setHistoryByStudent((prev) => {
        const next = { ...prev };
        delete next[studentId];
        return next;
      });

      await loadClassHistory(0);
    } catch (err: unknown) {
      const message = extractErrorMessage(err);
      setError(message || "Không thể huỷ điểm danh đã lưu");
    } finally {
      setCanceling(null);
    }
  };

  const handleSaveAll = async () => {
    const dirtyStudents = students.filter((student) => isRowDirty(student.student_id));
    if (dirtyStudents.length === 0) {
      return;
    }

    setSavingAll(true);
    setError("");
    try {
      await Promise.all(
        dirtyStudents.map(async (student) => {
          const att = attendance[student.student_id];
          if (!att) {
            return;
          }
          await teacherApi.markAttendance({
            student_id: student.student_id,
            date: today,
            status: att.status,
            note: att.note,
          });
        })
      );

      setSavedAttendance((prev) => {
        const next = { ...prev };
        dirtyStudents.forEach((student) => {
          const att = attendance[student.student_id];
          if (att) {
            next[student.student_id] = { status: att.status, note: att.note };
          }
        });
        return next;
      });

      setHasSavedToday((prev) => {
        const next = { ...prev };
        dirtyStudents.forEach((student) => {
          next[student.student_id] = true;
        });
        return next;
      });
    } catch (err: unknown) {
      const message = extractErrorMessage(err);
      setError(message || "Lỗi khi lưu hàng loạt");
    } finally {
      setSavingAll(false);
    }
  };

  const handleSaveDisplayed = async () => {
    const dirtyStudents = displayedStudents.filter((student) => isRowDirty(student.student_id));
    if (dirtyStudents.length === 0) {
      return;
    }

    setSavingDisplayed(true);
    setError("");
    try {
      await Promise.all(
        dirtyStudents.map(async (student) => {
          const att = attendance[student.student_id];
          if (!att) {
            return;
          }
          await teacherApi.markAttendance({
            student_id: student.student_id,
            date: today,
            status: att.status,
            note: att.note,
          });
        })
      );

      setSavedAttendance((prev) => {
        const next = { ...prev };
        dirtyStudents.forEach((student) => {
          const att = attendance[student.student_id];
          if (att) {
            next[student.student_id] = { status: att.status, note: att.note };
          }
        });
        return next;
      });

      setHasSavedToday((prev) => {
        const next = { ...prev };
        dirtyStudents.forEach((student) => {
          next[student.student_id] = true;
        });
        return next;
      });
    } catch (err: unknown) {
      const message = extractErrorMessage(err);
      setError(message || "Lỗi khi lưu danh sách đang hiển thị");
    } finally {
      setSavingDisplayed(false);
    }
  };

  const applyStatusToDisplayed = (status: AttendanceStatus) => {
    if (displayedStudents.length === 0) {
      return;
    }

    setAttendance((prev) => {
      const next = { ...prev };
      displayedStudents.forEach((student) => {
        const current = next[student.student_id] || { status: "present" as AttendanceStatus, note: "" };
        next[student.student_id] = { ...current, status };
      });
      return next;
    });
  };

  const toggleHistory = async (studentId: string) => {
    const isOpened = historyOpen.has(studentId);
    const nextOpen = new Set(historyOpen);
    if (isOpened) {
      nextOpen.delete(studentId);
      setHistoryOpen(nextOpen);
      return;
    }
    nextOpen.add(studentId);
    setHistoryOpen(nextOpen);

    if (historyByStudent[studentId]) {
      return;
    }

    const nextLoading = new Set(historyLoading);
    nextLoading.add(studentId);
    setHistoryLoading(nextLoading);
    try {
      const from = new Date();
      from.setDate(from.getDate() - 30);
      const fromDate = from.toISOString().slice(0, 10);

      const records = await teacherApi.getStudentAttendanceChanges(studentId, fromDate, today);
      setHistoryByStudent((prev) => ({ ...prev, [studentId]: records || [] }));
    } catch {
      setHistoryByStudent((prev) => ({ ...prev, [studentId]: [] }));
    } finally {
      const loadingAfter = new Set(historyLoading);
      loadingAfter.delete(studentId);
      setHistoryLoading(loadingAfter);
    }
  };

  const loadClassHistory = async (offset: number = historyOffset) => {
    if (!selectedClassId || students.length === 0) {
      setHistoryList([]);
      setHistoryTotal(0);
      setHistoryHasMore(false);
      return;
    }

    setHistoryListLoading(true);
    setError("");

    try {
      const response = await teacherApi.getClassAttendanceChanges(selectedClassId, {
        from: historyFrom,
        to: historyTo,
        student_id: historyStudentId === "all" ? undefined : historyStudentId,
        status: historyStatus === "all" ? undefined : historyStatus,
        limit: historyLimit,
        offset,
      });

      const list = response.data || [];
      const total = response.pagination?.total || 0;
      const hasMore = response.pagination?.has_more || false;

      setHistoryList(list.map((item) => ({
        ...item,
        student_name: item.student_name || students.find((student) => student.student_id === item.student_id)?.full_name || "Không rõ",
      })));
      setHistoryTotal(total);
      setHistoryHasMore(hasMore);
      setHistoryOffset(offset);
    } catch (err: unknown) {
      const message = extractErrorMessage(err);
      setError(message || "Không thể tải lịch sử điểm danh");
      setHistoryList([]);
      setHistoryTotal(0);
      setHistoryHasMore(false);
    } finally {
      setHistoryListLoading(false);
    }
  };

  const handleHistorySearch = () => {
    loadClassHistory(0);
  };

  const handleHistoryPrev = () => {
    if (historyOffset <= 0) {
      return;
    }
    loadClassHistory(Math.max(0, historyOffset - historyLimit));
  };

  const handleHistoryNext = () => {
    if (!historyHasMore) {
      return;
    }
    loadClassHistory(historyOffset + historyLimit);
  };

  if (loadingClasses) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        <Button size="sm" variant={viewMode === "take" ? "default" : "outline"} onClick={() => setViewMode("take")}>Điểm danh hôm nay</Button>
        <Button size="sm" variant={viewMode === "history" ? "default" : "outline"} onClick={() => setViewMode("history")}>Lịch sử lớp</Button>
      </div>

      <div className="flex items-center gap-2">
        {classes.length > 0 && (
          <Select value={selectedClassId} onValueChange={setSelectedClassId}>
            <SelectTrigger className="w-full sm:w-[220px]"><SelectValue placeholder="Chọn lớp" /></SelectTrigger>
            <SelectContent>
              {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
            </SelectContent>
          </Select>
        )}
      </div>

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}
      {loadingStudents && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingStudents && students.length === 0 && selectedClassId && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <ClipboardCheck className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Không có học sinh</p>
        </CardContent></Card>
      )}

      {!loadingStudents && students.length > 0 && viewMode === "take" && (
        <div className="space-y-2.5">
          <Card>
            <CardContent className="py-3">
              <div className="flex items-center justify-between gap-2 sm:hidden">
                <p className="text-xs text-muted-foreground">
                  Hiển thị {displayedStudents.length}/{students.length} • Chờ lưu {displayedDirtyCount}
                </p>
                <Button
                  size="sm"
                  variant="outline"
                  className="h-8 px-2.5 text-xs"
                  onClick={() => setShowMobileTakeControls((prev) => !prev)}
                  aria-expanded={showMobileTakeControls}
                >
                  {showMobileTakeControls ? "Ẩn bộ lọc" : "Mở bộ lọc"}
                </Button>
              </div>

              <div className={`${showMobileTakeControls ? "mt-2" : "hidden"} sm:mt-0 sm:block`}>
                <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
                  <Input
                    value={studentSearch}
                    onChange={(e) => setStudentSearch(e.target.value)}
                    placeholder="Tìm học sinh theo tên..."
                    className="h-9 text-sm"
                    aria-label="Tìm học sinh theo tên"
                  />

                  <Select value={takeListFilter} onValueChange={(value: TakeListFilter) => setTakeListFilter(value)}>
                    <SelectTrigger className="h-9 text-sm" aria-label="Lọc theo trạng thái lưu">
                      <SelectValue placeholder="Lọc danh sách" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">Tất cả học sinh</SelectItem>
                      <SelectItem value="pending">Chưa lưu / đang sửa</SelectItem>
                      <SelectItem value="saved">Đã lưu</SelectItem>
                    </SelectContent>
                  </Select>

                  <Button
                    size="sm"
                    variant={listOrderMode === "prioritize" ? "default" : "outline"}
                    onClick={() => setListOrderMode("prioritize")}
                    className="h-9"
                  >
                    Ưu tiên chưa lưu
                  </Button>
                  <Button
                    size="sm"
                    variant={listOrderMode === "original" ? "default" : "outline"}
                    onClick={() => setListOrderMode("original")}
                    className="h-9"
                  >
                    Giữ nguyên thứ tự
                  </Button>
                </div>

                <div className="mt-2 flex flex-wrap items-center gap-1.5">
                  <Badge variant="outline" className="text-xs">
                    Toàn lớp chờ lưu: {globalPendingCount}
                  </Badge>
                  <Badge variant="outline" className="text-xs">
                    Đang hiển thị: {displayedStudents.length}/{students.length}
                  </Badge>
                  <Badge variant={displayedDirtyCount > 0 ? "secondary" : "outline"} className="text-xs">
                    Chờ lưu trong danh sách: {displayedDirtyCount}
                  </Badge>
                  <Badge variant={displayedSavedCount > 0 ? "default" : "outline"} className="text-xs">
                    Đã lưu trong danh sách: {displayedSavedCount}
                  </Badge>
                </div>

                <div className="mt-2 flex flex-wrap items-center gap-2">
                  <Button size="sm" variant="outline" onClick={() => applyStatusToDisplayed("present")}>Đặt tất cả hiển thị: Có mặt</Button>
                  <Button size="sm" variant="outline" onClick={() => applyStatusToDisplayed("absent")}>Đặt tất cả hiển thị: Vắng</Button>
                  <Button size="sm" variant="outline" onClick={() => applyStatusToDisplayed("late")}>Đặt tất cả hiển thị: Muộn</Button>
                  <Button size="sm" variant="outline" onClick={() => applyStatusToDisplayed("excused")}>Đặt tất cả hiển thị: Có phép</Button>
                  <Button size="sm" onClick={handleSaveDisplayed} disabled={savingDisplayed || displayedDirtyCount === 0}>
                    {savingDisplayed ? <Loader2 className="h-4 w-4 animate-spin" /> : `Lưu danh sách hiển thị${displayedDirtyCount > 0 ? ` (${displayedDirtyCount})` : ""}`}
                  </Button>
                </div>

                <p className="mt-1 text-xs text-muted-foreground">
                  {listOrderMode === "prioritize" ? "Đang ưu tiên chưa lưu" : "Đang giữ nguyên thứ tự"} • Dùng “Lưu danh sách hiển thị” để chốt nhanh phần đang lọc.
                </p>
              </div>
            </CardContent>
          </Card>

          {displayedStudents.length === 0 && (
            <Card>
              <CardContent className="py-6 text-sm text-muted-foreground">
                Không có học sinh phù hợp với bộ lọc hiện tại. Hãy đổi từ khóa tìm kiếm hoặc chuyển bộ lọc danh sách.
              </CardContent>
            </Card>
          )}

          {displayedStudents.map((s) => {
            const att = attendance[s.student_id] || { status: "present", note: "" };
            const isDirty = isRowDirty(s.student_id);
            const isSavingThisRow = submitting === s.student_id;
            return (
              <Card key={s.student_id} className={!isDirty && hasSavedToday[s.student_id] ? "border-success/30 bg-success/10" : ""}>
                <CardContent className="px-3 py-3 sm:px-4">
                  <div className="flex items-start justify-between gap-2 sm:items-center">
                    <div className="min-w-0">
                      <p className="text-sm font-medium leading-tight truncate">
                        {s.full_name}
                        <span className="ml-1.5 text-xs font-normal text-muted-foreground">• {formatDateVN(s.dob)}</span>
                      </p>
                      <p className="mt-0.5 text-xs text-muted-foreground">
                        {!hasSavedToday[s.student_id] ? "Chưa lưu" : isDirty ? "Đã chỉnh sửa, chưa lưu" : "Đã lưu"}
                      </p>
                    </div>

                    <div className="w-[120px] shrink-0 sm:hidden">
                      <Select
                        value={att.status}
                        onValueChange={(value: AttendanceStatus) =>
                          setAttendance((prev) => ({ ...prev, [s.student_id]: { ...att, status: value } }))
                        }
                      >
                        <SelectTrigger className="h-8 text-xs">
                          <SelectValue placeholder="Trạng thái" />
                        </SelectTrigger>
                        <SelectContent>
                          {statusOptions.map((opt) => (
                            <SelectItem key={opt.value} value={opt.value}>
                              {opt.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div className="hidden flex-wrap items-center gap-1.5 sm:flex">
                      {statusOptions.map((opt) => (
                        <Badge key={opt.value}
                          variant={att.status === opt.value ? opt.variant : "outline"}
                          className={`h-6 cursor-pointer select-none px-2 text-xs transition-all ${att.status === opt.value ? "ring-2 ring-offset-1 ring-zinc-400" : "opacity-70 hover:opacity-100"}`}
                          onClick={() => setAttendance((prev) => ({ ...prev, [s.student_id]: { ...att, status: opt.value } }))}
                        >{opt.label}</Badge>
                      ))}
                    </div>
                  </div>
                  <div className="mt-2 flex items-center gap-1.5">
                    <Input placeholder="Ghi chú..." value={att.note} className="h-8 text-xs"
                      onChange={(e) => setAttendance((prev) => ({ ...prev, [s.student_id]: { ...att, note: e.target.value } }))} />
                    <Button size="sm" className="h-8 px-2.5 text-xs" onClick={() => handleMark(s.student_id)} disabled={isSavingThisRow}
                      variant={!isDirty && hasSavedToday[s.student_id] ? "outline" : "default"}>
                      {isSavingThisRow ? <Loader2 className="h-4 w-4 animate-spin" /> : !isDirty && hasSavedToday[s.student_id] ? <Check className="h-4 w-4" /> : hasSavedToday[s.student_id] ? "Cập nhật" : "Lưu"}
                    </Button>
                    {hasSavedToday[s.student_id] && isDirty && (
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-8 px-2 text-xs"
                        onClick={() => handleRevertLocal(s.student_id)}
                      >
                        Hoàn tác
                      </Button>
                    )}
                    {hasSavedToday[s.student_id] && !isDirty && (
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-8 px-2 text-xs text-destructive hover:text-destructive"
                        disabled={canceling === s.student_id}
                        onClick={() => handleCancelSaved(s.student_id)}
                      >
                        {canceling === s.student_id ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : "Huỷ lưu hôm nay"}
                      </Button>
                    )}
                  </div>

                  <div className="mt-2">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="h-7 px-0 text-xs text-muted-foreground hover:text-foreground"
                      onClick={() => toggleHistory(s.student_id)}
                    >
                      <History className="mr-1 h-3.5 w-3.5" />
                      {historyOpen.has(s.student_id) ? "Ẩn lịch sử" : "Xem lịch sử 30 ngày"}
                    </Button>

                    {historyOpen.has(s.student_id) && (
                      <div className="mt-2 rounded-md border bg-muted/30 p-2">
                        {historyLoading.has(s.student_id) ? (
                          <div className="flex items-center gap-2 text-xs text-muted-foreground">
                            <Loader2 className="h-3.5 w-3.5 animate-spin" />
                            Đang tải lịch sử...
                          </div>
                        ) : (historyByStudent[s.student_id] || []).length === 0 ? (
                          <p className="text-xs text-muted-foreground">Chưa có lịch sử điểm danh.</p>
                        ) : (
                          <div className="space-y-1.5">
                            {(historyByStudent[s.student_id] || []).slice(0, 8).map((record) => (
                              <div key={record.change_id} className="space-y-0.5 text-xs">
                                <div className="flex items-start justify-between gap-2">
                                  <span className="text-muted-foreground">{new Date(record.changed_at).toLocaleString("vi-VN")}</span>
                                  <span className="font-medium">
                                    {record.change_type === "create" ? "Tạo mới" : record.change_type === "delete" ? "Huỷ lưu" : "Cập nhật"}
                                  </span>
                                </div>
                                <div className="flex items-start justify-between gap-2 text-muted-foreground">
                                  <span>
                                    {record.change_type === "create"
                                      ? `Tạo: ${mapStatusLabel(record.new_status)}`
                                      : record.change_type === "delete"
                                        ? `${mapStatusLabel(record.old_status)} → Đã huỷ`
                                        : `${mapStatusLabel(record.old_status)} → ${mapStatusLabel(record.new_status)}`}
                                  </span>
                                  <span className="line-clamp-1 max-w-[45%] text-right">
                                    {record.change_type === "delete"
                                      ? `${record.old_note || "-"} → Đã xoá`
                                      : `${record.old_note || "-"} → ${record.new_note || "-"}`}
                                  </span>
                                </div>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            );
          })}

          {(displayedDirtyCount > 0 || globalPendingCount > 0) && (
            <div className="sticky bottom-3 z-20 rounded-lg border bg-background/95 p-3 shadow-sm backdrop-blur supports-[backdrop-filter]:bg-background/70">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-xs text-muted-foreground sm:text-sm">
                  Còn {displayedDirtyCount} học sinh chưa lưu trong danh sách hiển thị • Toàn lớp còn {globalPendingCount} học sinh chưa lưu.
                </p>

                <div className="flex flex-wrap items-center gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={handleSaveDisplayed}
                    disabled={savingDisplayed || displayedDirtyCount === 0}
                  >
                    {savingDisplayed ? <Loader2 className="h-4 w-4 animate-spin" /> : `Lưu danh sách hiển thị${displayedDirtyCount > 0 ? ` (${displayedDirtyCount})` : ""}`}
                  </Button>
                  <Button size="sm" onClick={handleSaveAll} disabled={savingAll || dirtyCount === 0}>
                    {savingAll ? <Loader2 className="h-4 w-4 animate-spin" /> : `Lưu toàn lớp${dirtyCount > 0 ? ` (${dirtyCount})` : ""}`}
                  </Button>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {!loadingStudents && students.length > 0 && viewMode === "history" && (
        <div className="space-y-3">
          <Card>
            <CardContent className="space-y-3 py-4">
              <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
                <Input type="date" value={historyFrom} onChange={(e) => setHistoryFrom(e.target.value)} className="h-9 text-sm" />
                <Input type="date" value={historyTo} onChange={(e) => setHistoryTo(e.target.value)} className="h-9 text-sm" />

                <Select value={historyStudentId} onValueChange={setHistoryStudentId}>
                  <SelectTrigger className="h-9 text-sm"><SelectValue placeholder="Học sinh" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Tất cả học sinh</SelectItem>
                    {students.map((student) => (
                      <SelectItem key={student.student_id} value={student.student_id}>{student.full_name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <Select value={historyStatus} onValueChange={(value: AttendanceStatus | "all") => setHistoryStatus(value)}>
                  <SelectTrigger className="h-9 text-sm"><SelectValue placeholder="Trạng thái" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Tất cả trạng thái</SelectItem>
                    {statusOptions.map((option) => (
                      <SelectItem key={option.value} value={option.value}>{option.label}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="flex items-center justify-between">
                <Button size="sm" onClick={handleHistorySearch} disabled={historyListLoading}>
                  {historyListLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : "Xem lịch sử"}
                </Button>
                <p className="text-xs text-muted-foreground">Tổng bản ghi: {historyTotal}</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="py-4">
              {historyListLoading ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Loader2 className="h-4 w-4 animate-spin" /> Đang tải lịch sử...
                </div>
              ) : historyList.length === 0 ? (
                <p className="text-sm text-muted-foreground">Không có dữ liệu lịch sử phù hợp.</p>
              ) : (
                <div className="space-y-2">
                  {historyList.map((record) => (
                    <div key={record.change_id} className="flex items-start justify-between gap-3 rounded-md border px-3 py-2 text-sm">
                      <div className="min-w-0">
                        <p className="font-medium leading-tight">{record.student_name}</p>
                        <p className="text-xs text-muted-foreground">{new Date(record.changed_at).toLocaleString("vi-VN")}</p>
                      </div>

                      <div className="shrink-0 text-right">
                        <p className="text-sm font-medium">
                          {record.change_type === "create"
                            ? `Tạo: ${mapStatusLabel(record.new_status)}`
                            : record.change_type === "delete"
                              ? `${mapStatusLabel(record.old_status)} → Đã huỷ`
                              : `${mapStatusLabel(record.old_status)} → ${mapStatusLabel(record.new_status)}`}
                        </p>
                        <p className="line-clamp-1 text-xs text-muted-foreground">
                          {record.change_type === "delete"
                            ? `${record.old_note || "-"} → Đã xoá`
                            : `${record.old_note || "-"} → ${record.new_note || "-"}`}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          {!historyListLoading && (
            <div className="flex items-center justify-between">
              <Button size="sm" variant="outline" onClick={handleHistoryPrev} disabled={historyOffset === 0}>
                Trang trước
              </Button>
              <p className="text-xs text-muted-foreground">
                {historyTotal === 0 ? "0-0" : `${historyOffset + 1}-${Math.min(historyOffset + historyLimit, historyTotal)}`} / {historyTotal}
              </p>
              <Button size="sm" variant="outline" onClick={handleHistoryNext} disabled={!historyHasMore}>
                Trang sau
              </Button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
