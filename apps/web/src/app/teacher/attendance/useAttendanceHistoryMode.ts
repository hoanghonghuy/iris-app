import { useCallback, useEffect, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { AttendanceChangeLog, AttendanceStatus, Student } from "@/types";
import { extractAttendanceErrorMessage } from "./utils";

export type AttendanceChangeHistoryRow = AttendanceChangeLog & {
  student_name: string;
};

export interface UseAttendanceHistoryModeParams {
  selectedClassId: string;
  students: Student[];
  onError: (message: string) => void;
}

export interface UseAttendanceHistoryModeResult {
  historyFrom: string;
  historyTo: string;
  historyStudentId: string;
  historyStatus: AttendanceStatus | "all";
  historyListLoading: boolean;
  historyList: AttendanceChangeHistoryRow[];
  historyOffset: number;
  historyLimit: number;
  historyTotal: number;
  historyHasMore: boolean;
  setHistoryFrom: (value: string) => void;
  setHistoryTo: (value: string) => void;
  setHistoryStudentId: (value: string) => void;
  setHistoryStatus: (value: AttendanceStatus | "all") => void;
  loadClassHistory: (offset?: number) => Promise<void>;
  handleHistorySearch: () => void;
  handleHistoryPrev: () => void;
  handleHistoryNext: () => void;
}

export function useAttendanceHistoryMode({ selectedClassId, students, onError }: UseAttendanceHistoryModeParams): UseAttendanceHistoryModeResult {
  const [historyFrom, setHistoryFrom] = useState(() => {
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - 7);
    return startDate.toISOString().slice(0, 10);
  });
  const [historyTo, setHistoryTo] = useState(() => new Date().toISOString().slice(0, 10));
  const [historyStudentId, setHistoryStudentId] = useState("all");
  const [historyStatus, setHistoryStatus] = useState<AttendanceStatus | "all">("all");
  const [historyListLoading, setHistoryListLoading] = useState(false);
  const [historyList, setHistoryList] = useState<AttendanceChangeHistoryRow[]>([]);
  const [historyOffset, setHistoryOffset] = useState(0);
  const [historyLimit] = useState(20);
  const [historyTotal, setHistoryTotal] = useState(0);
  const [historyHasMore, setHistoryHasMore] = useState(false);

  useEffect(() => {
    setHistoryStudentId("all");
    setHistoryList([]);
    setHistoryOffset(0);
    setHistoryTotal(0);
    setHistoryHasMore(false);
  }, [selectedClassId]);

  const loadClassHistory = useCallback(async (offset: number = historyOffset) => {
    if (!selectedClassId || students.length === 0) {
      setHistoryList([]);
      setHistoryTotal(0);
      setHistoryHasMore(false);
      return;
    }

    setHistoryListLoading(true);
    onError("");

    try {
      const response = await teacherApi.getClassAttendanceChanges(selectedClassId, {
        from: historyFrom,
        to: historyTo,
        student_id: historyStudentId === "all" ? undefined : historyStudentId,
        status: historyStatus === "all" ? undefined : historyStatus,
        limit: historyLimit,
        offset,
      });

      const responseList = response.data || [];
      const total = response.pagination?.total || 0;
      const hasMore = response.pagination?.has_more || false;

      setHistoryList(
        responseList.map((item) => ({
          ...item,
          student_name:
            item.student_name ||
            students.find((student) => student.student_id === item.student_id)?.full_name ||
            "Không rõ",
        }))
      );
      setHistoryTotal(total);
      setHistoryHasMore(hasMore);
      setHistoryOffset(offset);
    } catch (errorValue: unknown) {
      const message = extractAttendanceErrorMessage(errorValue);
      onError(message || "Không thể tải lịch sử điểm danh");
      setHistoryList([]);
      setHistoryTotal(0);
      setHistoryHasMore(false);
    } finally {
      setHistoryListLoading(false);
    }
  }, [selectedClassId, students, historyFrom, historyTo, historyStudentId, historyStatus, historyLimit, historyOffset, onError]);

  const handleHistorySearch = useCallback(() => {
    loadClassHistory(0);
  }, [loadClassHistory]);

  const handleHistoryPrev = useCallback(() => {
    if (historyOffset <= 0) {
      return;
    }
    loadClassHistory(Math.max(0, historyOffset - historyLimit));
  }, [historyOffset, historyLimit, loadClassHistory]);

  const handleHistoryNext = useCallback(() => {
    if (!historyHasMore) {
      return;
    }
    loadClassHistory(historyOffset + historyLimit);
  }, [historyHasMore, historyOffset, historyLimit, loadClassHistory]);

  return {
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
  };
}