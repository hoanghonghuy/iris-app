/**
 * Parent Children Page
 * Xem danh sách con em.
 * API: GET /parent/children
 */
"use client";

import React, { useEffect, useState } from "react";
import { parentApi } from "@/lib/api/parent.api";
import { Student } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Users, Loader2, User, Calendar } from "lucide-react";

const genderLabel: Record<string, string> = { male: "Nam", female: "Nữ", other: "Khác" };

export default function ParentChildrenPage() {
  const [children, setChildren] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const data = await parentApi.getMyChildren();
        setChildren(data || []);
      } catch (err: any) {
        setError(err.response?.data?.error || "Không thể tải danh sách con em");
      } finally { setLoading(false); }
    };
    load();
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Users className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Con em của tôi</h1>
      </div>

      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && children.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <Users className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có con em nào được liên kết</p>
        </CardContent></Card>
      )}

      {!loading && children.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {children.map((child) => (
            <Card key={child.student_id}>
              <CardContent className="flex items-start gap-4 py-5">
                <User className="h-10 w-10 shrink-0 rounded-full bg-zinc-100 p-2 text-muted-foreground" />
                <div>
                  <p className="text-lg font-medium">{child.full_name}</p>
                  <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                    <p className="flex items-center gap-2">
                      <Calendar className="h-3 w-3" />
                      Ngày sinh: {child.dob}
                    </p>
                    <p>Giới tính: {genderLabel[child.gender] || child.gender}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
