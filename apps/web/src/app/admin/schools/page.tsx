/**
 * Admin Schools Page
 * Quản lý danh sách trường học: xem, tạo mới.
 * API: GET /admin/schools, POST /admin/schools
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { School, CreateSchoolRequest } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { School as SchoolIcon, Plus, X, Loader2, MapPin } from "lucide-react";

export default function AdminSchoolsPage() {
  // ─── State ────────────────────────────────────────────────────────

  const [schools, setSchools] = useState<School[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Form state
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState<CreateSchoolRequest>({ name: "", address: "" });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  // ─── Data fetching ────────────────────────────────────────────────

  const fetchSchools = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const data = await adminApi.getSchools();
      setSchools(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách trường học");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSchools();
  }, [fetchSchools]);

  // ─── Create school ────────────────────────────────────────────────

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      setFormError("Tên trường không được để trống");
      return;
    }

    try {
      setSubmitting(true);
      setFormError("");
      await adminApi.createSchool(formData);
      setFormData({ name: "", address: "" });
      setShowForm(false);
      fetchSchools(); // refresh list
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo trường học");
    } finally {
      setSubmitting(false);
    }
  };

  // ─── Render ───────────────────────────────────────────────────────

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <SchoolIcon className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Quản lý Trường học</h1>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
          {showForm ? "Hủy" : "Thêm trường"}
        </Button>
      </div>

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Thêm trường học mới</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && (
                <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                  {formError}
                </div>
              )}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label htmlFor="name" className="text-sm font-medium">
                    Tên trường <span className="text-destructive">*</span>
                  </label>
                  <Input
                    id="name"
                    placeholder="VD: Trường Mầm Non Hoa Sen"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="address" className="text-sm font-medium">
                    Địa chỉ
                  </label>
                  <Input
                    id="address"
                    placeholder="VD: 123 Nguyễn Văn A, Q.1, TP.HCM"
                    value={formData.address}
                    onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  />
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Tạo trường
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Error State */}
      {error && (
        <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">
          {error}
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {/* Empty State */}
      {!loading && !error && schools.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <SchoolIcon className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Chưa có trường học nào</p>
            <Button variant="outline" className="mt-4" onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Thêm trường đầu tiên
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Desktop Table (md+) */}
      {!loading && schools.length > 0 && (
        <div className="hidden md:block">
          <Card>
            <CardContent className="p-0">
              <table className="w-full">
                <thead>
                  <tr className="border-b text-left text-sm text-muted-foreground">
                    <th className="px-6 py-3 font-medium">Tên trường</th>
                    <th className="px-6 py-3 font-medium">Địa chỉ</th>
                  </tr>
                </thead>
                <tbody>
                  {schools.map((school) => (
                    <tr key={school.school_id} className="border-b last:border-0 hover:bg-zinc-50">
                      <td className="px-6 py-4 font-medium">{school.name}</td>
                      <td className="px-6 py-4 text-muted-foreground">
                        {school.address || "—"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Mobile Card List (<md) */}
      {!loading && schools.length > 0 && (
        <div className="space-y-3 md:hidden">
          {schools.map((school) => (
            <Card key={school.school_id}>
              <CardContent className="flex items-start gap-3 py-4">
                <SchoolIcon className="mt-0.5 h-5 w-5 shrink-0 text-muted-foreground" />
                <div className="min-w-0 flex-1">
                  <p className="font-medium">{school.name}</p>
                  {school.address && (
                    <p className="mt-1 flex items-center gap-1 text-sm text-muted-foreground">
                      <MapPin className="h-3 w-3" />
                      {school.address}
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
