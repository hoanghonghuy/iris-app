import { ApiError } from "@/types";

export function extractAttendanceErrorMessage(err: unknown): string | undefined {
  return typeof err === "object" &&
    err !== null &&
    "response" in err &&
    typeof (err as { response?: { data?: ApiError } }).response?.data?.error === "string"
    ? (err as { response?: { data?: ApiError } }).response?.data?.error
    : undefined;
}
