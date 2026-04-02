/**
 * Admin School Admins Page
 * Quản lý School Admin: listing + tạo mới + xóa.
 * API: GET/POST/DELETE /admin/school-admins
 */
"use client";

import React, { useRef, useEffect } from "react";
import { SchoolAdmin } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { ResponsiveSplitView } from "@/components/shared/ResponsiveSplitView";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { ShieldCheck, Loader2, Plus, X, Trash2, Mail, AlertCircle, CheckCircle2, Search } from "lucide-react";
import { useAdminSchoolAdminsPage } from "./useAdminSchoolAdminsPage";

export default function AdminSchoolAdminsPage() {
  const {
    admins,
    loading,
    error,
    pagination,
    showForm,
    schools,
    selectedSchoolId,
    userSearchQuery,
    userSearchResults,
    userSearchLoading,
    selectedUser,
    showUserDropdown,
    submitting,
    formError,
    success,
    deletingId,
    deleteAlert,
    setCurrentOffset,
    setShowForm,
    setSelectedSchoolId,
    setUserSearchQuery,
    selectUser,
    clearSelectedUser,
    setSuccess,
    setDeleteAlert,
    closeDeleteAlert,
    handleCreate,
    confirmDelete,
  } = useAdminSchoolAdminsPage();

  const getAdminId = (admin: SchoolAdmin): string => admin.admin_id;

  // Ref để đóng dropdown khi click bên ngoài
  const dropdownRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        // Không cần set state trực tiếp — chỉ blur input
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

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
                {/* User Search */}
                <div className="space-y-2">
                  <Label>User</Label>
                  <div className="relative" ref={dropdownRef}>
                    {selectedUser ? (
                      <div className="flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm">
                        <Mail className="h-4 w-4 shrink-0 text-muted-foreground" />
                        <span className="flex-1 truncate">{selectedUser.email}</span>
                        <button
                          type="button"
                          onClick={clearSelectedUser}
                          className="shrink-0 rounded-full p-0.5 hover:bg-muted transition-colors"
                        >
                          <X className="h-3.5 w-3.5 text-muted-foreground" />
                        </button>
                      </div>
                    ) : (
                      <>
                        <div className="relative">
                          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                          <Input
                            type="text"
                            placeholder="Tìm theo email..."
                            className="pl-8"
                            value={userSearchQuery}
                            onChange={(e) => setUserSearchQuery(e.target.value)}
                            autoComplete="off"
                          />
                          {userSearchLoading && (
                            <Loader2 className="absolute right-2.5 top-2.5 h-4 w-4 animate-spin text-muted-foreground" />
                          )}
                        </div>
                        {showUserDropdown && userSearchResults.length > 0 && (
                          <div className="absolute z-50 mt-1 w-full rounded-lg border bg-popover shadow-lg overflow-hidden animate-in fade-in slide-in-from-top-1 duration-150">
                            {userSearchResults.map((user) => (
                              <button
                                key={user.user_id}
                                type="button"
                                className="w-full text-left px-3 py-2.5 text-sm hover:bg-muted transition-colors flex items-center gap-2 border-b last:border-0"
                                onClick={() => selectUser(user)}
                              >
                                <Mail className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                                <span className="truncate">{user.email}</span>
                              </button>
                            ))}
                          </div>
                        )}
                        {userSearchQuery.trim().length >= 2 && !userSearchLoading && userSearchResults.length === 0 && (
                          <div className="absolute z-50 mt-1 w-full rounded-lg border bg-popover p-3 text-sm text-muted-foreground shadow-lg">
                            Không tìm thấy user nào
                          </div>
                        )}
                      </>
                    )}
                  </div>
                </div>
                {/* School Select */}
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
                <Button type="submit" disabled={submitting || !selectedUser}>
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
