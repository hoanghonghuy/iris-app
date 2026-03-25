/**
 * Parent Child Detail Page
 * Xem thông tin chi tiết của một con cụ thể.
 * API: GET /parent/children
 */
"use client";

import React, { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { parentApi } from "@/lib/api/parent.api";
import { Student } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Loader2, Calendar, ArrowLeft } from "lucide-react";
import { formatDateVN } from "@/lib/utils";

const genderLabel: Record<string, string> = { male: "Nam", female: "Nữ", other: "Khác" };

function getAgeFromDOB(dob: string): number | null {
  const match = dob.match(/^(\d{4})-(\d{2})-(\d{2})/);
  if (!match) {
    return null;
  }

  const birthDate = new Date(Number(match[1]), Number(match[2]) - 1, Number(match[3]));
  if (Number.isNaN(birthDate.getTime())) {
    return null;
  }

  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();
  const dayDiff = today.getDate() - birthDate.getDate();

  if (monthDiff < 0 || (monthDiff === 0 && dayDiff < 0)) {
    age -= 1;
  }

  return age;
}

export default function ParentChildDetailPage() {
  const params = useParams<{ studentId: string }>();
  const studentId = params?.studentId;

  const [children, setChildren] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const data = await parentApi.getMyChildren();
        setChildren(data || []);
      } catch (err: unknown) {
        const message =
          typeof err === "object" &&
          err !== null &&
          "response" in err &&
          typeof (err as { response?: { data?: { error?: string } } }).response?.data?.error === "string"
            ? (err as { response?: { data?: { error?: string } } }).response?.data?.error
            : undefined;
        setError(message || "Không thể tải thông tin con");
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  const child = useMemo(
    () => children.find((item) => item.student_id === studentId),
    [children, studentId]
  );
  const childAge = useMemo(() => (child ? getAgeFromDOB(child.dob) : null), [child]);

  return (
    <div className="space-y-6">
      <Link href="/parent/children" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="h-4 w-4" />
        Quay lại danh sách con
      </Link>

      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {!loading && !error && !child && (
        <Card>
          <CardContent className="py-8 text-sm text-muted-foreground">
            Không tìm thấy thông tin của bé này.
          </CardContent>
        </Card>
      )}

      {!loading && !error && child && (
        <>
          <Card>
            <CardHeader>
              <CardTitle>{child.full_name}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm text-muted-foreground">
              <p className="flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Ngày sinh: {formatDateVN(child.dob)}
              </p>
              <p>Giới tính: {genderLabel[child.gender] || child.gender}</p>
              {childAge !== null && <p>{childAge} tuổi</p>}
              <p>Lớp hiện tại: {child.current_class_name || "Chưa có thông tin"}</p>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}
