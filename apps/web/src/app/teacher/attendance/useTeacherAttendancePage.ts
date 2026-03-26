import { useCallback, useEffect, useMemo, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { AttendanceChangeLog, AttendanceStatus, Class, Student } from "@/types";
import { extractAttendanceErrorMessage } from "./utils";
import { useAttendanceTakeMode } from "./useAttendanceTakeMode";
import { useAttendanceHistoryMode } from "./useAttendanceHistoryMode";

type AttendanceValue = { status: AttendanceStatus; note: string };

function getDefaultAttendanceValue(): AttendanceValue {
  return { status: "present", note: "" };
}

export function useTeacherAttendancePage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  const [attendance, setAttendance] = useState<Record<string, AttendanceValue>>({});
  const [savedAttendance, setSavedAttendance] = useState<Record<string, AttendanceValue>>({});
  const [hasSavedToday, setHasSavedToday] = useState<Record<string, boolean>>({});

  const [submitting, setSubmitting] = useState<string | null>(null);
  const [canceling, setCanceling] = useState<string | null>(null);
  const [savingAll, setSavingAll] = useState(false);
  const [savingDisplayed, setSavingDisplayed] = useState(false);

  const [historyOpen, setHistoryOpen] = useState<Set<string>>(new Set());
  const [historyLoading, setHistoryLoading] = useState<Set<string>>(new Set());
  const [historyByStudent, setHistoryByStudent] = useState<Record<string, AttendanceChangeLog[]>>({});

  const [today] = useState(() => new Date().toISOString().slice(0, 10));
  const [viewMode, setViewMode] = useState<"take" | "history">("take");

  const isRowDirty = useCallback((studentId: string) => {
    const currentValue = attendance[studentId];
    const savedValue = savedAttendance[studentId];
    if (!currentValue || !savedValue) {
      return !hasSavedToday[studentId];
    }
    if (!hasSavedToday[studentId]) {
      return true;
    }
    return currentValue.status !== savedValue.status || (currentValue.note || "") !== (savedValue.note || "");
  }, [attendance, savedAttendance, hasSavedToday]);

  const dirtyCount = useMemo(
    () => students.filter((student) => isRowDirty(student.student_id)).length,
    [students, isRowDirty]
  );

  const {
    studentSearch,
    listOrderMode,
    takeListFilter,
    showMobileTakeControls,
    displayedStudents,
    displayedDirtyCount,
    displayedSavedCount,
    globalPendingCount,
    setStudentSearch,
    setListOrderMode,
    setTakeListFilter,
    setShowMobileTakeControls,
  } = useAttendanceTakeMode({ students, isRowDirty });

  const {
    historyFrom,
    historyTo,
    historyStudentId,
    historyStatus,
    historyListLoading,
    historyList,
    historyOffset,
    historyLimit,
    historyTotal,
    historyHasMore,
    setHistoryFrom,
    setHistoryTo,
    setHistoryStudentId,
    setHistoryStatus,
    loadClassHistory,
    handleHistorySearch,
    handleHistoryPrev,
    handleHistoryNext,
  } = useAttendanceHistoryMode({
    selectedClassId,
    students,
    onError: setError,
  });

  useEffect(() => {
    const loadClasses = async () => {
      try {
        const classList = await teacherApi.getMyClasses();
        setClasses(classList || []);
        if (classList && classList.length > 0) {
          setSelectedClassId(classList[0].class_id);
        }
      } catch {
        setError("Không thể tải lớp");
      } finally {
        setLoadingClasses(false);
      }
    };

    loadClasses();
  }, []);

  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) {
      return;
    }

    try {
      setLoadingStudents(true);
      setError("");
      setHistoryOpen(new Set());
      setHistoryByStudent({});

      const studentList = (await teacherApi.getStudentsInClass(selectedClassId)) || [];
      setStudents(studentList);

      const initialAttendance: Record<string, AttendanceValue> = {};
      const initialSavedAttendance: Record<string, AttendanceValue> = {};
      const initialHasSavedToday: Record<string, boolean> = {};

      await Promise.all(
        studentList.map(async (student: Student) => {
          try {
            const records = await teacherApi.getStudentAttendance(student.student_id, today, today);
            const todayRecord = records.find((record) => record.date.startsWith(today));

            if (todayRecord) {
              const savedValue = { status: todayRecord.status, note: todayRecord.note || "" };
              initialAttendance[student.student_id] = savedValue;
              initialSavedAttendance[student.student_id] = savedValue;
              initialHasSavedToday[student.student_id] = true;
              return;
            }
          } catch {
          }

          const defaultValue = getDefaultAttendanceValue();
          initialAttendance[student.student_id] = defaultValue;
          initialSavedAttendance[student.student_id] = defaultValue;
          initialHasSavedToday[student.student_id] = false;
        })
      );

      setAttendance(initialAttendance);
      setSavedAttendance(initialSavedAttendance);
      setHasSavedToday(initialHasSavedToday);
    } catch (errorValue: unknown) {
      const message = extractAttendanceErrorMessage(errorValue);
      setError(message || "Không thể tải HS");
    } finally {
      setLoadingStudents(false);
    }
  }, [selectedClassId, today]);

  useEffect(() => {
    fetchStudents();
  }, [fetchStudents]);

  const handleMark = useCallback(async (studentId: string) => {
    const currentAttendance = attendance[studentId];
    if (!currentAttendance) {
      return;
    }

    try {
      setSubmitting(studentId);
      await teacherApi.markAttendance({
        student_id: studentId,
        date: today,
        status: currentAttendance.status,
        note: currentAttendance.note,
      });
      setSavedAttendance((previous) => ({
        ...previous,
        [studentId]: { status: currentAttendance.status, note: currentAttendance.note },
      }));
      setHasSavedToday((previous) => ({ ...previous, [studentId]: true }));
    } catch (errorValue: unknown) {
      const message = extractAttendanceErrorMessage(errorValue);
      setError(message || "Lỗi điểm danh");
    } finally {
      setSubmitting(null);
    }
  }, [attendance, today]);

  const handleRevertLocal = useCallback((studentId: string) => {
    const savedValue = savedAttendance[studentId];
    if (!savedValue) {
      return;
    }
    setAttendance((previous) => ({
      ...previous,
      [studentId]: { status: savedValue.status, note: savedValue.note },
    }));
  }, [savedAttendance]);

  const handleCancelSaved = useCallback(async (studentId: string) => {
    if (!hasSavedToday[studentId]) {
      return;
    }

    setCanceling(studentId);
    setError("");

    try {
      await teacherApi.cancelAttendance(studentId, today);

      const defaultValue = getDefaultAttendanceValue();
      setAttendance((previous) => ({ ...previous, [studentId]: defaultValue }));
      setSavedAttendance((previous) => ({ ...previous, [studentId]: defaultValue }));
      setHasSavedToday((previous) => ({ ...previous, [studentId]: false }));
      setHistoryByStudent((previous) => {
        const nextValue = { ...previous };
        delete nextValue[studentId];
        return nextValue;
      });

      await loadClassHistory(0);
    } catch (errorValue: unknown) {
      const message = extractAttendanceErrorMessage(errorValue);
      setError(message || "Không thể huỷ điểm danh đã lưu");
    } finally {
      setCanceling(null);
    }
  }, [hasSavedToday, today, loadClassHistory]);

  const handleSaveAll = useCallback(async () => {
    const dirtyStudents = students.filter((student) => isRowDirty(student.student_id));
    if (dirtyStudents.length === 0) {
      return;
    }

    setSavingAll(true);
    setError("");

    try {
      await Promise.all(
        dirtyStudents.map(async (student) => {
          const currentAttendance = attendance[student.student_id];
          if (!currentAttendance) {
            return;
          }
          await teacherApi.markAttendance({
            student_id: student.student_id,
            date: today,
            status: currentAttendance.status,
            note: currentAttendance.note,
          });
        })
      );

      setSavedAttendance((previous) => {
        const nextValue = { ...previous };
        dirtyStudents.forEach((student) => {
          const currentAttendance = attendance[student.student_id];
          if (currentAttendance) {
            nextValue[student.student_id] = {
              status: currentAttendance.status,
              note: currentAttendance.note,
            };
          }
        });
        return nextValue;
      });

      setHasSavedToday((previous) => {
        const nextValue = { ...previous };
        dirtyStudents.forEach((student) => {
          nextValue[student.student_id] = true;
        });
        return nextValue;
      });
    } catch (errorValue: unknown) {
      const message = extractAttendanceErrorMessage(errorValue);
      setError(message || "Lỗi khi lưu hàng loạt");
    } finally {
      setSavingAll(false);
    }
  }, [students, isRowDirty, attendance, today]);

  const handleSaveDisplayed = useCallback(async () => {
    const dirtyStudents = displayedStudents.filter((student) => isRowDirty(student.student_id));
    if (dirtyStudents.length === 0) {
      return;
    }

    setSavingDisplayed(true);
    setError("");

    try {
      await Promise.all(
        dirtyStudents.map(async (student) => {
          const currentAttendance = attendance[student.student_id];
          if (!currentAttendance) {
            return;
          }
          await teacherApi.markAttendance({
            student_id: student.student_id,
            date: today,
            status: currentAttendance.status,
            note: currentAttendance.note,
          });
        })
      );

      setSavedAttendance((previous) => {
        const nextValue = { ...previous };
        dirtyStudents.forEach((student) => {
          const currentAttendance = attendance[student.student_id];
          if (currentAttendance) {
            nextValue[student.student_id] = {
              status: currentAttendance.status,
              note: currentAttendance.note,
            };
          }
        });
        return nextValue;
      });

      setHasSavedToday((previous) => {
        const nextValue = { ...previous };
        dirtyStudents.forEach((student) => {
          nextValue[student.student_id] = true;
        });
        return nextValue;
      });
    } catch (errorValue: unknown) {
      const message = extractAttendanceErrorMessage(errorValue);
      setError(message || "Lỗi khi lưu danh sách đang hiển thị");
    } finally {
      setSavingDisplayed(false);
    }
  }, [displayedStudents, isRowDirty, attendance, today]);

  const applyStatusToDisplayed = useCallback((status: AttendanceStatus) => {
    if (displayedStudents.length === 0) {
      return;
    }

    setAttendance((previous) => {
      const nextValue = { ...previous };
      displayedStudents.forEach((student) => {
        const currentValue = nextValue[student.student_id] || getDefaultAttendanceValue();
        nextValue[student.student_id] = { ...currentValue, status };
      });
      return nextValue;
    });
  }, [displayedStudents]);

  const handleAttendanceStatusChange = useCallback((studentId: string, status: AttendanceStatus) => {
    setAttendance((previous) => {
      const currentValue = previous[studentId] || getDefaultAttendanceValue();
      return {
        ...previous,
        [studentId]: { ...currentValue, status },
      };
    });
  }, []);

  const handleAttendanceNoteChange = useCallback((studentId: string, note: string) => {
    setAttendance((previous) => {
      const currentValue = previous[studentId] || getDefaultAttendanceValue();
      return {
        ...previous,
        [studentId]: { ...currentValue, note },
      };
    });
  }, []);

  const toggleHistory = useCallback(async (studentId: string) => {
    let shouldFetchHistory = false;
    setHistoryOpen((previous) => {
      const nextValue = new Set(previous);
      if (nextValue.has(studentId)) {
        nextValue.delete(studentId);
      } else {
        nextValue.add(studentId);
        shouldFetchHistory = true;
      }
      return nextValue;
    });

    if (!shouldFetchHistory || historyByStudent[studentId]) {
      return;
    }

    setHistoryLoading((previous) => {
      const nextValue = new Set(previous);
      nextValue.add(studentId);
      return nextValue;
    });

    try {
      const fromDate = new Date();
      fromDate.setDate(fromDate.getDate() - 30);
      const formattedFromDate = fromDate.toISOString().slice(0, 10);
      const records = await teacherApi.getStudentAttendanceChanges(studentId, formattedFromDate, today);
      setHistoryByStudent((previous) => ({ ...previous, [studentId]: records || [] }));
    } catch {
      setHistoryByStudent((previous) => ({ ...previous, [studentId]: [] }));
    } finally {
      setHistoryLoading((previous) => {
        const nextValue = new Set(previous);
        nextValue.delete(studentId);
        return nextValue;
      });
    }
  }, [historyByStudent, today]);

  return {
    classes,
    selectedClassId,
    students,
    loadingClasses,
    loadingStudents,
    error,
    submitting,
    canceling,
    savingAll,
    savingDisplayed,
    historyOpen,
    historyLoading,
    historyByStudent,
    studentSearch,
    listOrderMode,
    takeListFilter,
    showMobileTakeControls,
    viewMode,
    historyFrom,
    historyTo,
    historyStudentId,
    historyStatus,
    historyListLoading,
    historyList,
    historyOffset,
    historyLimit,
    historyTotal,
    historyHasMore,
    dirtyCount,
    displayedStudents,
    displayedDirtyCount,
    displayedSavedCount,
    globalPendingCount,
    attendance,
    hasSavedToday,
    setSelectedClassId,
    setStudentSearch,
    setListOrderMode,
    setTakeListFilter,
    setShowMobileTakeControls,
    setViewMode,
    setHistoryFrom,
    setHistoryTo,
    setHistoryStudentId,
    setHistoryStatus,
    isRowDirty,
    handleMark,
    handleRevertLocal,
    handleCancelSaved,
    handleSaveAll,
    handleSaveDisplayed,
    applyStatusToDisplayed,
    toggleHistory,
    handleHistorySearch,
    handleHistoryPrev,
    handleHistoryNext,
    handleAttendanceStatusChange,
    handleAttendanceNoteChange,
  };
}