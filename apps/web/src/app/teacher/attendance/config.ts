import { AttendanceStatus } from "@/types";

export type TakeListFilter = "all" | "pending" | "saved";

export const ATTENDANCE_STATUS_OPTIONS = [
  { value: "present" as AttendanceStatus, label: "Có mặt", variant: "default" as const },
  { value: "absent" as AttendanceStatus, label: "Vắng", variant: "destructive" as const },
  { value: "late" as AttendanceStatus, label: "Muộn", variant: "secondary" as const },
  { value: "excused" as AttendanceStatus, label: "Có phép", variant: "outline" as const },
];

export const ATTENDANCE_STATUS_LABEL: Record<AttendanceStatus, string> = {
  present: "Có mặt",
  absent: "Vắng",
  late: "Muộn",
  excused: "Có phép",
};

export function mapAttendanceStatusLabel(status?: AttendanceStatus): string {
  if (!status) {
    return "-";
  }
  return ATTENDANCE_STATUS_LABEL[status] || status;
}
