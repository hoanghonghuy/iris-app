/**
 * Admin School Admins Page
 * Quản lý School Admin: listing + tạo mới + xóa.
 * API: GET/POST/DELETE /admin/school-admins
 */
"use client";

import React from "react";
import { SchoolAdmin } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { ResponsiveSplitView } from "@/components/shared/ResponsiveSplitView";
import { ShieldCheck, Loader2, Plus, X, Trash2, Mail, AlertCircle, CheckCircle2 } from "lucide-react";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { useAdminSchoolAdminsPage } from "./useAdminSchoolAdminsPage";

export default function AdminSchoolAdminsPage() {
  const {
    admins,
    loading,
    error,
    pagination,
    showForm,
    schools,
    users,
    selectedSchoolId,
    selectedUserId,
    submitting,
    formError,
    success,
    deletingId,
    deleteAlert,
    setCurrentOffset,
    setShowForm,
    setSelectedSchoolId,
    setSelectedUserId,
    setSuccess,
    setDeleteAlert,
    closeDeleteAlert,
    handleCreate,
    confirmDelete,
  } = useAdminSchoolAdminsPage();

  const getAdminId = (admin: SchoolAdmin): string => admin.admin_id;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Button size="sm" onClick={() => { setShowForm(!showForm); setSuccess(""); }}>
          {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
          {showForm ? "Hủy" : "Thêm School Admin"}
        </Button>
      </div>

      {success && <Alert><CheckCircle2 className="h-4 w-4 text-success" /><AlertDescription>{success}</AlertDescription></Alert>}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Gán School Admin</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>User</Label>
                  <Select value={selectedUserId} onValueChange={setSelectedUserId}>
                    <SelectTrigger className="w-full"><SelectValue placeholder="Chọn user" /></SelectTrigger>
                    <SelectContent>
                      {users.map((u) => <SelectItem key={u.user_id} value={u.user_id}>{u.email}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Trường</Label>
                  <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
                    <SelectTrigger className="w-full"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
                    <SelectContent>
                      {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Gán
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && admins.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <ShieldCheck className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có School Admin nào</p>
        </CardContent></Card>
      )}

      <ResponsiveSplitView
        show={!loading && admins.length > 0}
        desktop={(
          <Card><CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Email</TableHead>
                  <TableHead>Trường</TableHead>
                  <TableHead className="text-right">Hành động</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {admins.map((a) => (
                  <TableRow key={getAdminId(a)}>
                    <TableCell className="font-medium">{a.email || a.user_id}</TableCell>
                    <TableCell className="text-muted-foreground">{a.school_name || a.school_id}</TableCell>
                    <TableCell className="text-right">
                      <Button variant="ghost" size="sm" onClick={() => setDeleteAlert({ isOpen: true, adminId: getAdminId(a) })} disabled={deletingId === getAdminId(a)}>
                        {deletingId === getAdminId(a) ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="mr-1 h-4 w-4 text-destructive" />} Xóa
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent></Card>
        )}
        mobileClassName="space-y-3 md:hidden"
        mobile={(
          <>
            {admins.map((a) => (
              <Card key={getAdminId(a)}>
                <CardContent className="py-4">
                  <div className="flex items-start justify-between">
                    <div>
                      <p className="flex items-center gap-2 font-medium"><Mail className="h-4 w-4 text-muted-foreground" /> {a.email || a.user_id}</p>
                      <p className="mt-1 text-sm text-muted-foreground">{a.school_name || a.school_id}</p>
                    </div>
                    <Button variant="ghost" size="sm" onClick={() => setDeleteAlert({ isOpen: true, adminId: getAdminId(a) })}>
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </>
        )}
      />

      {/* Pagination */}
      {!loading && admins.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}

      {/* Delete Confirmation */}
      <ConfirmAlertDialog
        isOpen={deleteAlert.isOpen}
        onClose={closeDeleteAlert}
        onConfirm={confirmDelete}
        title="Xác nhận xóa"
        description="Bạn có chắc chắn muốn xóa School Admin này? Hành động này không thể hoàn tác."
        loading={!!deletingId}
        confirmText="Xác nhận xóa"
      />
    </div>
  );
}
