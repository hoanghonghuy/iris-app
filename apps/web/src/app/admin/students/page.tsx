/**
 * Admin Students Page
 * Quản lý học sinh theo lớp: chọn trường → chọn lớp → xem danh sách + tạo mới + tạo mã phụ huynh.
 * API: GET /admin/students/by-class/:class_id, POST /admin/students, POST /admin/students/:id/generate-parent-code
 */
"use client";

import React from "react";
import Link from "next/link";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { ActionModal } from "@/components/shared/ActionModal";
import { ResponsiveSplitView } from "@/components/shared/ResponsiveSplitView";
import { toast } from "sonner";
import {
  Users, Plus, X, Loader2, Calendar, User, KeyRound, Copy, Check, AlertCircle, Search, Pencil, Trash2
} from "lucide-react";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { useAuth } from "@/providers/AuthProvider";
import { genderLabel, useAdminStudentsPage } from "./useAdminStudentsPage";

type StudentFormState = {
  full_name: string;
  dob: string;
  gender: "male" | "female" | "other";
};

export default function AdminStudentsPage() {
  const { role } = useAuth();
  const {
    schools,
    classes,
    selectedSchoolId,
    selectedClassId,
    loadingSchools,
    students,
    searchQuery,
    loadingStudents,
    error,
    showForm,
    formData,
    submitting,
    formError,
    generatingCode,
    revokingCode,
    revokeAlert,
    copiedId,
    codeError,
    filteredStudents,
    selectedClassName,
    editModal,
    editData,
    editLoading,
    deleteAlert,
    deleteLoading,
    setSelectedSchoolId,
    setSelectedClassId,
    setSearchQuery,
    setShowForm,
    setFormData,
    setRevokeAlert,
    closeRevokeAlert,
    setEditModal,
    setEditData,
    setDeleteAlert,
    closeEditModal,
    closeDeleteAlert,
    handleCreate,
    handleEdit,
    handleDelete,
    handleGenerateCode,
    confirmRevokeCode,
    handleCopy,
    getDaysLeft,
  } = useAdminStudentsPage();

  const toDateInputValue = (dob: string) => {
    return dob.includes("T") ? dob.split("T")[0] : dob;
  };

  const toggleCreateForm = () => {
    setShowForm((prev) => !prev);
  };

  const openCreateForm = () => {
    setShowForm(true);
  };

  const handleFormFieldChange = (field: keyof StudentFormState, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleFormGenderChange = (value: string) => {
    if (value === "male" || value === "female" || value === "other") {
      setFormData((prev) => ({ ...prev, gender: value }));
    }
  };

  const handleEditFieldChange = (field: keyof StudentFormState, value: string) => {
    setEditData((prev) => ({ ...prev, [field]: value }));
  };

  const handleEditGenderChange = (value: "male" | "female" | "other") => {
    setEditData((prev) => ({ ...prev, gender: value }));
  };

  const openRevokeAlert = (studentId: string) => {
    setRevokeAlert({ isOpen: true, studentId });
  };

  const openEditModal = (student: (typeof filteredStudents)[number]) => {
    setEditData({
      full_name: student.full_name,
      dob: toDateInputValue(student.dob),
      gender: student.gender as "male" | "female" | "other",
    });
    setEditModal({ isOpen: true, selectedStudent: student });
  };

  const openDeleteAlert = (studentId: string) => {
    setDeleteAlert({ isOpen: true, studentId });
  };

  const handleEditWithFeedback = async () => {
    try {
      await handleEdit();
      toast.success("Cập nhật học sinh thành công");
    } catch {
      toast.error("Không thể cập nhật học sinh");
    }
  };

  const handleDeleteWithFeedback = async () => {
    try {
      await handleDelete();
      toast.success("Xóa học sinh thành công");
    } catch {
      toast.error("Không thể xóa học sinh");
    }
  };

  if (loadingSchools) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-wrap items-center gap-2">
          {role === 'SUPER_ADMIN' && (
            <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
              <SelectTrigger className="w-[180px]"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
              <SelectContent>
                {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          {classes.length > 0 && (
            <Select value={selectedClassId} onValueChange={setSelectedClassId}>
              <SelectTrigger className="w-[160px]"><SelectValue placeholder="Chọn lớp" /></SelectTrigger>
              <SelectContent>
                {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          {selectedClassId && (
            <Button size="sm" onClick={toggleCreateForm}>
              {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
              {showForm ? "Hủy" : "Thêm HS"}
            </Button>
          )}
        </div>
      </div>

      {/* Toolbar: Search */}
      {!loadingStudents && !error && students.length > 0 && !showForm && (
        <div className="relative max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input 
            type="search" 
            placeholder="Tìm theo tên học sinh..." 
            className="pl-8 bg-background " 
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      )}

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Thêm học sinh — {selectedClassName}</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && (
                <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>
              )}
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="fullName">Họ tên <span className="text-destructive">*</span></Label>
                  <Input id="fullName" placeholder="VD: Bé An" value={formData.full_name}
                    onChange={(e) => handleFormFieldChange("full_name", e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="dob">Ngày sinh <span className="text-destructive">*</span></Label>
                  <Input id="dob" type="date" value={formData.dob}
                    onChange={(e) => handleFormFieldChange("dob", e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label>Giới tính</Label>
                  <Select
                    value={formData.gender}
                    onValueChange={handleFormGenderChange}
                  >
                    <SelectTrigger className="w-full"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="male">Nam</SelectItem>
                      <SelectItem value="female">Nữ</SelectItem>
                      <SelectItem value="other">Khác</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Tạo học sinh
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Errors */}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}
      {codeError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{codeError}</AlertDescription></Alert>}

      {/* Loading */}
      {loadingStudents && (
        <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>
      )}

      {/* Empty (No students at all) */}
      {!loadingStudents && !error && students.length === 0 && selectedClassId && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Users className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Chưa có học sinh nào trong {selectedClassName}</p>
            <Button variant="outline" className="mt-4" onClick={openCreateForm}>
              <Plus className="mr-2 h-4 w-4" /> Thêm học sinh đầu tiên
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Empty Search Results */}
      {!loadingStudents && !error && students.length > 0 && filteredStudents.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy học sinh nào phù hợp với &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      <ResponsiveSplitView
        show={!loadingStudents && filteredStudents.length > 0}
        desktop={(
          <Card><CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Họ tên</TableHead>
                  <TableHead>Ngày sinh</TableHead>
                  <TableHead>Giới tính</TableHead>
                  <TableHead className="text-right">Mã PH</TableHead>
                  <TableHead className="text-right">Hành động</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredStudents.map((s) => (
                  <TableRow key={s.student_id}>
                    <TableCell className="font-medium">
                      <Link href={`/admin/students/${s.student_id}`} className="hover:text-primary hover:underline transition-colors">
                        {s.full_name}
                      </Link>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{s.dob}</TableCell>
                    <TableCell><Badge variant="secondary">{genderLabel[s.gender] || s.gender}</Badge></TableCell>
                    <TableCell className="text-right">
                      {s.active_parent_code ? (
                        <div className="flex flex-col items-end gap-1">
                          <div className="flex items-center gap-1">
                            <code className="rounded bg-muted px-2 py-0.5 text-xs font-mono">{s.active_parent_code}</code>
                            <Button variant="ghost" size="icon" className="h-6 w-6" onClick={() => handleCopy(s.active_parent_code as string, s.student_id)}>
                              {copiedId === s.student_id ? <Check className="h-3 w-3 text-success" /> : <Copy className="h-3 w-3" />}
                            </Button>
                            <Button variant="ghost" size="icon" className="h-6 w-6 text-destructive hover:bg-destructive/10 hover:text-destructive" onClick={() => openRevokeAlert(s.student_id)} disabled={revokingCode === s.student_id}>
                              {revokingCode === s.student_id ? <Loader2 className="h-3 w-3 animate-spin" /> : <X className="h-3 w-3" />}
                            </Button>
                          </div>
                          <span className={`text-[10px] ${getDaysLeft(s.code_expires_at) === 'Hết hạn' ? 'text-destructive font-medium' : 'text-muted-foreground'}`}>{getDaysLeft(s.code_expires_at)}</span>
                        </div>
                      ) : (
                        <Button variant="ghost" size="sm" onClick={() => handleGenerateCode(s.student_id)} disabled={generatingCode === s.student_id}>
                          {generatingCode === s.student_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <KeyRound className="mr-1 h-4 w-4" />} Tạo mã
                        </Button>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button variant="ghost" size="sm" onClick={() => openEditModal(s)}>
                        <Pencil className="h-4 w-4 text-muted-foreground hover:text-primary" />
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => openDeleteAlert(s.student_id)}>
                        <Trash2 className="h-4 w-4 text-destructive hover:text-destructive/80" />
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
            {filteredStudents.map((s) => (
              <Card key={s.student_id}>
                <CardContent className="flex items-start gap-3 py-4">
                  <User className="mt-0.5 h-5 w-5 shrink-0 text-muted-foreground" />
                  <div className="min-w-0 flex-1">
                    <p className="font-medium">
                      <Link href={`/admin/students/${s.student_id}`} className="hover:text-primary hover:underline transition-colors">
                        {s.full_name}
                      </Link>
                    </p>
                    <div className="mt-1 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
                      <span className="flex items-center gap-1"><Calendar className="h-3 w-3" /> {s.dob}</span>
                      <Badge variant="secondary">{genderLabel[s.gender] || s.gender}</Badge>
                    </div>
                    {s.active_parent_code ? (
                      <div className="mt-2 flex items-center justify-between gap-1 w-full bg-muted/50 p-2 rounded-md">
                        <div>
                          <code className="rounded bg-background px-2 py-0.5 text-sm font-mono tracking-wider shadow-sm">{s.active_parent_code}</code>
                          <div className={`text-[10px] mt-1 ${getDaysLeft(s.code_expires_at) === 'Hết hạn' ? 'text-destructive font-medium' : 'text-muted-foreground'}`}>{getDaysLeft(s.code_expires_at)}</div>
                        </div>
                        <div className="flex gap-1">
                          <Button variant="secondary" size="sm" className="h-8 shadow-sm" onClick={() => handleCopy(s.active_parent_code as string, s.student_id)}>
                            {copiedId === s.student_id ? <Check className="h-3 w-3 text-success" /> : <Copy className="h-3 w-3" />}
                          </Button>
                          <Button variant="outline" size="sm" className="h-8 shadow-sm text-destructive hover:bg-destructive/10" onClick={() => openRevokeAlert(s.student_id)} disabled={revokingCode === s.student_id}>
                            {revokingCode === s.student_id ? <Loader2 className="h-3 w-3 animate-spin" /> : <X className="h-3 w-3" />}
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <Button variant="secondary" size="sm" className="mt-3 w-full" onClick={() => handleGenerateCode(s.student_id)} disabled={generatingCode === s.student_id}>
                        {generatingCode === s.student_id ? <Loader2 className="h-3 w-3 animate-spin" /> : <KeyRound className="mr-2 h-3 w-3" />} Tạo mã PH
                      </Button>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </>
        )}
      />

      {/* Confirm Revoke Code Alert */}
      <ConfirmAlertDialog
        isOpen={revokeAlert.isOpen}
        onClose={closeRevokeAlert}
        onConfirm={confirmRevokeCode}
        title="Xác nhận thu hồi mã"
        description="Mã phụ huynh hiện tại sẽ bị vô hiệu hóa. Phụ huynh đang sử dụng mã này sẽ bị đăng xuất. Bạn có chắc chắn muốn tiếp tục?"
        loading={revokingCode !== null}
        confirmText="Thu hồi"
        cancelText="Hủy"
      />

      {/* Edit Student Modal */}
      <ActionModal
        isOpen={editModal.isOpen}
        onClose={closeEditModal}
        onConfirm={handleEditWithFeedback}
        title="Sửa thông tin học sinh"
        loading={editLoading}
        confirmText="Lưu"
        cancelText="Hủy"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label>Họ và tên</Label>
            <Input value={editData.full_name} onChange={(e) => handleEditFieldChange("full_name", e.target.value)} />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label>Ngày sinh</Label>
              <Input type="date" value={editData.dob} onChange={(e) => handleEditFieldChange("dob", e.target.value)} />
            </div>
            <div className="space-y-1.5">
              <Label>Giới tính</Label>
              <Select value={editData.gender} onValueChange={handleEditGenderChange}>
                <SelectTrigger><SelectValue placeholder="Chọn giới tính" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="male">Nam</SelectItem>
                  <SelectItem value="female">Nữ</SelectItem>
                  <SelectItem value="other">Khác</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>
      </ActionModal>

      {/* Confirm Delete Student */}
      <ConfirmAlertDialog
        isOpen={deleteAlert.isOpen}
        onClose={closeDeleteAlert}
        onConfirm={handleDeleteWithFeedback}
        title="Xác nhận xóa học sinh"
        description="Bạn có chắc chắn muốn xóa học sinh này? Hành động này không thể hoàn tác."
        loading={deleteLoading}
        confirmText="Xóa"
        cancelText="Hủy"
      />
    </div>
  );
}
