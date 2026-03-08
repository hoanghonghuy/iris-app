/**
 * Admin Parents Page
 * Danh sách phụ huynh.
 * API: GET /admin/parents
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Parent } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Heart, Loader2, Phone, Mail } from "lucide-react";

export default function AdminParentsPage() {
  const [parents, setParents] = useState<Parent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchParents = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const data = await adminApi.getParents();
      setParents(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách phụ huynh");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchParents();
  }, [fetchParents]);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Heart className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý Phụ huynh</h1>
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
      )}

      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {!loading && parents.length === 0 && !error && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Heart className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Chưa có phụ huynh nào</p>
          </CardContent>
        </Card>
      )}

      {/* Desktop Table */}
      {!loading && parents.length > 0 && (
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
                  {parents.map((p) => (
                    <tr key={p.parent_id} className="border-b last:border-0 hover:bg-zinc-50">
                      <td className="px-6 py-4 font-medium">{p.full_name}</td>
                      <td className="px-6 py-4 text-muted-foreground">{p.email}</td>
                      <td className="px-6 py-4 text-muted-foreground">{p.phone || "—"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loading && parents.length > 0 && (
        <div className="space-y-3 md:hidden">
          {parents.map((p) => (
            <Card key={p.parent_id}>
              <CardContent className="py-4">
                <p className="font-medium">{p.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2">
                    <Mail className="h-3 w-3" /> {p.email}
                  </p>
                  {p.phone && (
                    <p className="flex items-center gap-2">
                      <Phone className="h-3 w-3" /> {p.phone}
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
