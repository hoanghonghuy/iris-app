import React from "react";
import { Check, History, Loader2 } from "lucide-react";
import { AttendanceChangeLog, AttendanceStatus, Student } from "@/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { formatDateVN } from "@/lib/utils";
import { ATTENDANCE_STATUS_OPTIONS, mapAttendanceStatusLabel } from "../config";

interface StudentAttendanceCardProps {
  student: Student;
  attendanceValue: { status: AttendanceStatus; note: string };
  isDirty: boolean;
  hasSavedToday: boolean;
  isSaving: boolean;
  isCanceling: boolean;
  isHistoryOpen: boolean;
  isHistoryLoading: boolean;
  historyRecords: AttendanceChangeLog[];
  onStatusChange: (studentId: string, status: AttendanceStatus) => void;
  onNoteChange: (studentId: string, note: string) => void;
  onSave: (studentId: string) => void;
  onRevert: (studentId: string) => void;
  onCancelSaved: (studentId: string) => void;
  onToggleHistory: (studentId: string) => void;
}

export function StudentAttendanceCard({
  student,
  attendanceValue,
  isDirty,
  hasSavedToday,
  isSaving,
  isCanceling,
  isHistoryOpen,
  isHistoryLoading,
  historyRecords,
  onStatusChange,
  onNoteChange,
  onSave,
  onRevert,
  onCancelSaved,
  onToggleHistory,
}: StudentAttendanceCardProps) {
  return (
    <Card className={!isDirty && hasSavedToday ? "border-success/30 bg-success/10" : ""}>
      <CardContent className="px-3 py-3 sm:px-4">
        <div className="flex items-start justify-between gap-2 sm:items-center">
          <div className="min-w-0">
            <p className="text-sm font-medium leading-tight truncate">
              {student.full_name}
              <span className="ml-1.5 text-xs font-normal text-muted-foreground">• {formatDateVN(student.dob)}</span>
            </p>
            <p className="mt-0.5 text-xs text-muted-foreground">
              {!hasSavedToday ? "Chưa lưu" : isDirty ? "Đã chỉnh sửa, chưa lưu" : "Đã lưu"}
            </p>
          </div>

          <div className="w-[120px] shrink-0 sm:hidden">
            <Select
              value={attendanceValue.status}
              onValueChange={(value: AttendanceStatus) => onStatusChange(student.student_id, value)}
            >
              <SelectTrigger className="h-8 text-xs">
                <SelectValue placeholder="Trạng thái" />
              </SelectTrigger>
              <SelectContent>
                {ATTENDANCE_STATUS_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="hidden flex-wrap items-center gap-1.5 sm:flex">
            {ATTENDANCE_STATUS_OPTIONS.map((option) => (
              <Badge
                key={option.value}
                variant={attendanceValue.status === option.value ? option.variant : "outline"}
                className={`h-6 cursor-pointer select-none px-2 text-xs transition-all ${attendanceValue.status === option.value ? "ring-2 ring-offset-1 ring-ring" : "opacity-70 hover:opacity-100"}`}
                onClick={() => onStatusChange(student.student_id, option.value)}
              >
                {option.label}
              </Badge>
            ))}
          </div>
        </div>

        <div className="mt-2 flex items-center gap-1.5">
          <Input
            placeholder="Ghi chú..."
            value={attendanceValue.note}
            className="h-8 text-xs"
            onChange={(e) => onNoteChange(student.student_id, e.target.value)}
          />
          <Button
            size="sm"
            className="h-8 px-2.5 text-xs"
            onClick={() => onSave(student.student_id)}
            disabled={isSaving}
            variant={!isDirty && hasSavedToday ? "outline" : "default"}
          >
            {isSaving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : !isDirty && hasSavedToday ? (
              <Check className="h-4 w-4" />
            ) : hasSavedToday ? (
              "Cập nhật"
            ) : (
              "Lưu"
            )}
          </Button>

          {hasSavedToday && isDirty && (
            <Button size="sm" variant="ghost" className="h-8 px-2 text-xs" onClick={() => onRevert(student.student_id)}>
              Hoàn tác
            </Button>
          )}

          {hasSavedToday && !isDirty && (
            <Button
              size="sm"
              variant="ghost"
              className="h-8 px-2 text-xs text-destructive hover:text-destructive"
              disabled={isCanceling}
              onClick={() => onCancelSaved(student.student_id)}
            >
              {isCanceling ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : "Huỷ lưu hôm nay"}
            </Button>
          )}
        </div>

        <div className="mt-2">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-7 px-0 text-xs text-muted-foreground hover:text-foreground"
            onClick={() => onToggleHistory(student.student_id)}
          >
            <History className="mr-1 h-3.5 w-3.5" />
            {isHistoryOpen ? "Ẩn lịch sử" : "Xem lịch sử 30 ngày"}
          </Button>

          {isHistoryOpen && (
            <div className="mt-2 rounded-md border bg-muted/30 p-2">
              {isHistoryLoading ? (
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Loader2 className="h-3.5 w-3.5 animate-spin" />
                  Đang tải lịch sử...
                </div>
              ) : historyRecords.length === 0 ? (
                <p className="text-xs text-muted-foreground">Chưa có lịch sử điểm danh.</p>
              ) : (
                <div className="space-y-1.5">
                  {historyRecords.slice(0, 8).map((record) => (
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
                            ? `Tạo: ${mapAttendanceStatusLabel(record.new_status)}`
                            : record.change_type === "delete"
                              ? `${mapAttendanceStatusLabel(record.old_status)} → Đã huỷ`
                              : `${mapAttendanceStatusLabel(record.old_status)} → ${mapAttendanceStatusLabel(record.new_status)}`}
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
}
