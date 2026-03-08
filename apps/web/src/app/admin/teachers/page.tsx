/**
 * Admin Teachers Page
 * Danh sách giáo viên.
 * API: GET /admin/teachers
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Teacher } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { BookUser, Loader2, Phone, Mail } from "lucide-react";

export default function AdminTeachersPage() {
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchTeachers = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const data = await adminApi.getTeachers();
      setTeachers(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách giáo viên");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTeachers();
  }, [fetchTeachers]);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <BookUser className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý Giáo viên</h1>
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
      )}

      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {!loading && teachers.length === 0 && !error && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <BookUser className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Chưa có giáo viên nào</p>
          </CardContent>
        </Card>
      )}

      {/* Desktop Table */}
      {!loading && teachers.length > 0 && (
        <div className="hidden md:block">
          <Card>
            <CardContent className="p-0">
              <table className="w-full">
                <thead>
                  <tr className="border-b text-left text-sm text-muted-foreground">
                    <th className="px-6 py-3 font-medium">Họ tên</th>
                    <th className="px-6 py-3 font-medium">Email</th>
                    <th className="px-6 py-3 font-medium">Điện thoại</th>
                  </tr>
                </thead>
                <tbody>
                  {teachers.map((t) => (
                    <tr key={t.teacher_id} className="border-b last:border-0 hover:bg-zinc-50">
                      <td className="px-6 py-4 font-medium">{t.full_name}</td>
                      <td className="px-6 py-4 text-muted-foreground">{t.email}</td>
                      <td className="px-6 py-4 text-muted-foreground">{t.phone || "—"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loading && teachers.length > 0 && (
        <div className="space-y-3 md:hidden">
          {teachers.map((t) => (
            <Card key={t.teacher_id}>
              <CardContent className="py-4">
                <p className="font-medium">{t.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2">
                    <Mail className="h-3 w-3" /> {t.email}
                  </p>
                  {t.phone && (
                    <p className="flex items-center gap-2">
                      <Phone className="h-3 w-3" /> {t.phone}
                    </p>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
