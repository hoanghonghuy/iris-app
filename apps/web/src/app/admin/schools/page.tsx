/**
 * Admin Schools Page
 * Quản lý danh sách trường học: xem, tạo mới.
 * API: GET /admin/schools, POST /admin/schools
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { School, CreateSchoolRequest, Pagination } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { TableSkeleton } from "@/components/shared/TableSkeleton";
import { CardSkeleton } from "@/components/shared/CardSkeleton";
import { EmptyState } from "@/components/shared/EmptyState";
import { ResponsiveSplitView } from "@/components/shared/ResponsiveSplitView";
import { toast } from "sonner";
import { School as SchoolIcon, Plus, X, Loader2, MapPin } from "lucide-react";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { extractApiErrorMessage } from "@/lib/api-error";

export default function AdminSchoolsPage() {
  // ─── State ────────────────────────────────────────────────────────

  const [schools, setSchools] = useState<School[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

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
      const response = await adminApi.getSchools({ limit: 20, offset: currentOffset });
      setSchools(response.data || []);
      if (response.pagination) setPagination(response.pagination);
    } catch (error: unknown) {
      const msg = extractApiErrorMessage(error, "Không thể tải danh sách trường học");
      setError(msg);
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }, [currentOffset]);

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
      toast.success("Tạo trường học thành công");
      fetchSchools(); // refresh list
    } catch (error: unknown) {
      setFormError(extractApiErrorMessage(error, "Không thể tạo trường học"));
    } finally {
      setSubmitting(false);
    }
  };

  // ─── Render ───────────────────────────────────────────────────────

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
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

      <ResponsiveSplitView
        show={loading}
        desktop={<TableSkeleton columns={2} rows={5} />}
        mobile={<CardSkeleton cards={3} />}
      />

      {/* Empty State */}
      {!loading && !error && schools.length === 0 && (
        <EmptyState
          icon={SchoolIcon}
          title="Chưa có trường học nào"
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Thêm trường đầu tiên
            </Button>
          }
        />
      )}

      <ResponsiveSplitView
        show={!loading && schools.length > 0}
        desktop={(
          <Card>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Tên trường</TableHead>
                    <TableHead>Địa chỉ</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {schools.map((school) => (
                    <TableRow key={school.school_id}>
                      <TableCell className="font-medium">{school.name}</TableCell>
                      <TableCell className="text-muted-foreground">
                        {school.address || "—"}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        )}
        mobileClassName="space-y-3 md:hidden"
        mobile={(
          <>
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
          </>
        )}
      />

      {/* Pagination */}
      {!loading && schools.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}
    </div>
  );
}
