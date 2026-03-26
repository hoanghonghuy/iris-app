import React from "react";
import { Loader2 } from "lucide-react";
import { AttendanceStatus, Student } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ATTENDANCE_STATUS_OPTIONS, mapAttendanceStatusLabel } from "../config";

type AttendanceChangeHistoryRow = {
  change_id: string;
  student_id: string;
  student_name: string;
  old_status?: AttendanceStatus;
  new_status?: AttendanceStatus;
  old_note?: string;
  new_note?: string;
  change_type: "create" | "update" | "delete";
  changed_at: string;
};

interface AttendanceHistoryViewProps {
  students: Student[];
  historyFrom: string;
  historyTo: string;
  historyStudentId: string;
  historyStatus: AttendanceStatus | "all";
  historyListLoading: boolean;
  historyList: AttendanceChangeHistoryRow[];
  historyTotal: number;
  historyOffset: number;
  historyLimit: number;
  historyHasMore: boolean;
  onHistoryFromChange: (value: string) => void;
  onHistoryToChange: (value: string) => void;
  onHistoryStudentChange: (value: string) => void;
  onHistoryStatusChange: (value: AttendanceStatus | "all") => void;
  onHistorySearch: () => void;
  onHistoryPrev: () => void;
  onHistoryNext: () => void;
}

export function AttendanceHistoryView({
  students,
  historyFrom,
  historyTo,
  historyStudentId,
  historyStatus,
  historyListLoading,
  historyList,
  historyTotal,
  historyOffset,
  historyLimit,
  historyHasMore,
  onHistoryFromChange,
  onHistoryToChange,
  onHistoryStudentChange,
  onHistoryStatusChange,
  onHistorySearch,
  onHistoryPrev,
  onHistoryNext,
}: AttendanceHistoryViewProps) {
  return (
    <div className="space-y-3">
      <Card>
        <CardContent className="space-y-3 py-4">
          <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
            <Input type="date" value={historyFrom} onChange={(e) => onHistoryFromChange(e.target.value)} className="h-9 text-sm" />
            <Input type="date" value={historyTo} onChange={(e) => onHistoryToChange(e.target.value)} className="h-9 text-sm" />

            <Select value={historyStudentId} onValueChange={onHistoryStudentChange}>
              <SelectTrigger className="h-9 text-sm"><SelectValue placeholder="Học sinh" /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Tất cả học sinh</SelectItem>
                {students.map((student) => (
                  <SelectItem key={student.student_id} value={student.student_id}>{student.full_name}</SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select value={historyStatus} onValueChange={(value: AttendanceStatus | "all") => onHistoryStatusChange(value)}>
              <SelectTrigger className="h-9 text-sm"><SelectValue placeholder="Trạng thái" /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Tất cả trạng thái</SelectItem>
                {ATTENDANCE_STATUS_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>{option.label}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center justify-between">
            <Button size="sm" onClick={onHistorySearch} disabled={historyListLoading}>
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
                        ? `Tạo: ${mapAttendanceStatusLabel(record.new_status)}`
                        : record.change_type === "delete"
                          ? `${mapAttendanceStatusLabel(record.old_status)} → Đã huỷ`
                          : `${mapAttendanceStatusLabel(record.old_status)} → ${mapAttendanceStatusLabel(record.new_status)}`}
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
          <Button size="sm" variant="outline" onClick={onHistoryPrev} disabled={historyOffset === 0}>
            Trang trước
          </Button>
          <p className="text-xs text-muted-foreground">
            {historyTotal === 0 ? "0-0" : `${historyOffset + 1}-${Math.min(historyOffset + historyLimit, historyTotal)}`} / {historyTotal}
          </p>
          <Button size="sm" variant="outline" onClick={onHistoryNext} disabled={!historyHasMore}>
            Trang sau
          </Button>
        </div>
      )}
    </div>
  );
}
