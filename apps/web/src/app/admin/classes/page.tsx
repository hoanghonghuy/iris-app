/**
 * Admin Classes Page
 * Quản lý lớp học theo trường: chọn trường → xem danh sách lớp + tạo lớp mới.
 * API: GET /admin/classes/by-school/:school_id, POST /admin/classes
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { School, Class } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ResponsiveSplitView } from "@/components/shared/ResponsiveSplitView";
import { toast } from "sonner";
import { GraduationCap, Plus, X, Loader2, Calendar, AlertCircle, Pencil, Trash2 } from "lucide-react";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { ActionModal } from "@/components/shared/ActionModal";
import { useAuth } from "@/providers/AuthProvider";
import { extractApiErrorMessage } from "@/lib/api-error";

type ClassEditModalState = {
  isOpen: boolean;
  classItem: Class | null;
};

type ClassDeleteAlertState = {
  isOpen: boolean;
  classId: string | null;
};

type ClassFormState = {
  name: string;
  school_year: string;
};

const INITIAL_CLASS_EDIT_MODAL_STATE: ClassEditModalState = {
  isOpen: false,
  classItem: null,
};

const INITIAL_CLASS_DELETE_ALERT_STATE: ClassDeleteAlertState = {
  isOpen: false,
  classId: null,
};

export default function AdminClassesPage() {
  const { role } = useAuth();
  const [schools, setSchools] = useState<School[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState<string>("");
  const [loadingSchools, setLoadingSchools] = useState(true);

  const [classes, setClasses] = useState<Class[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(false);
  const [error, setError] = useState("");

  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState<ClassFormState>({ name: "", school_year: "" });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  // Edit/Delete state
  const [editModal, setEditModal] = useState<ClassEditModalState>(INITIAL_CLASS_EDIT_MODAL_STATE);
  const [editData, setEditData] = useState<ClassFormState>({ name: "", school_year: "" });
  const [editLoading, setEditLoading] = useState(false);
  const [deleteAlert, setDeleteAlert] = useState<ClassDeleteAlertState>(INITIAL_CLASS_DELETE_ALERT_STATE);
  const [deleteLoading, setDeleteLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        const response = await adminApi.getSchools();
        const data = response.data;
        setSchools(data || []);
        if (data && data.length > 0) setSelectedSchoolId(data[0].school_id);
      } catch { setError("Không thể tải danh sách trường"); }
      finally { setLoadingSchools(false); }
    };
    load();
  }, []);

  const fetchClasses = useCallback(async () => {
    if (!selectedSchoolId) return;
    try {
      setLoadingClasses(true); setError("");
      const response = await adminApi.getClassesBySchool(selectedSchoolId);
      const data = response.data;
      setClasses(data || []);
    } catch (error: unknown) {
      setError(extractApiErrorMessage(error, "Không thể tải danh sách lớp"));
    } finally { setLoadingClasses(false); }
  }, [selectedSchoolId]);

  useEffect(() => { fetchClasses(); }, [fetchClasses]);

  const toggleCreateForm = () => {
    setShowForm((prev) => !prev);
  };

  const openCreateForm = () => {
    setShowForm(true);
  };

  const handleFormFieldChange = (field: keyof ClassFormState, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleEditFieldChange = (field: keyof ClassFormState, value: string) => {
    setEditData((prev) => ({ ...prev, [field]: value }));
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) { setFormError("Tên lớp không được để trống"); return; }
    if (!formData.school_year.trim()) { setFormError("Năm học không được để trống"); return; }
    try {
      setSubmitting(true); setFormError("");
      await adminApi.createClass({ school_id: selectedSchoolId, name: formData.name, school_year: formData.school_year });
      setFormData({ name: "", school_year: "" }); setShowForm(false); fetchClasses();
    } catch (error: unknown) {
      setFormError(extractApiErrorMessage(error, "Không thể tạo lớp"));
    } finally { setSubmitting(false); }
  };

  const openEditModal = (classItem: Class) => {
    setEditData({ name: classItem.name, school_year: classItem.school_year });
    setEditModal({ isOpen: true, classItem });
  };

  const closeEditModal = () => {
    setEditModal(INITIAL_CLASS_EDIT_MODAL_STATE);
  };

  const openDeleteAlert = (classId: string) => {
    setDeleteAlert({ isOpen: true, classId });
  };

  const closeDeleteAlert = () => {
    setDeleteAlert(INITIAL_CLASS_DELETE_ALERT_STATE);
  };

  const handleEdit = async () => {
    if (!editModal.classItem) return;
    try {
      setEditLoading(true);
      await adminApi.updateClass(editModal.classItem.class_id, editData);
      toast.success("Cập nhật lớp thành công");
      closeEditModal();
      await fetchClasses();
    } catch (error: unknown) {
      toast.error(extractApiErrorMessage(error, "Không thể cập nhật lớp"));
    } finally {
      setEditLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!deleteAlert.classId) return;
    try {
      setDeleteLoading(true);
      await adminApi.deleteClass(deleteAlert.classId);
      toast.success("Xóa lớp thành công");
      closeDeleteAlert();
      await fetchClasses();
    } catch (error: unknown) {
      toast.error(extractApiErrorMessage(error, "Không thể xóa lớp"));
    } finally {
      setDeleteLoading(false);
    }
  };

  const selectedSchoolName = schools.find((s) => s.school_id === selectedSchoolId)?.name || "";

  if (loadingSchools) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2">
          {role === 'SUPER_ADMIN' && (
            <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
              <SelectTrigger className="w-[200px]"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
              <SelectContent>
                {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          {selectedSchoolId && (
            <Button size="sm" onClick={toggleCreateForm}>
              {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
              {showForm ? "Hủy" : "Thêm lớp"}
            </Button>
          )}
        </div>
      </div>

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Thêm lớp mới — {selectedSchoolName}</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="className">Tên lớp <span className="text-destructive">*</span></Label>
                  <Input id="className" placeholder="VD: Lá Non" value={formData.name}
                    onChange={(e) => handleFormFieldChange("name", e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="schoolYear">Năm học <span className="text-destructive">*</span></Label>
                  <Input id="schoolYear" placeholder="VD: 2025-2026" value={formData.school_year}
                    onChange={(e) => handleFormFieldChange("school_year", e.target.value)} required />
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Tạo lớp
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {loadingClasses && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingClasses && !error && classes.length === 0 && selectedSchoolId && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <GraduationCap className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có lớp nào trong {selectedSchoolName}</p>
          <Button variant="outline" className="mt-4" onClick={openCreateForm}>
            <Plus className="mr-2 h-4 w-4" /> Thêm lớp đầu tiên
          </Button>
        </CardContent></Card>
      )}

      <ResponsiveSplitView
        show={!loadingClasses && classes.length > 0}
        desktop={(
          <Card><CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Tên lớp</TableHead>
                  <TableHead>Năm học</TableHead>
                  <TableHead className="text-right">Hành động</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {classes.map((classItem) => (
                  <TableRow key={classItem.class_id}>
                    <TableCell className="font-medium">{classItem.name}</TableCell>
                    <TableCell><Badge variant="secondary">{classItem.school_year}</Badge></TableCell>
                    <TableCell className="text-right">
                      <Button variant="ghost" size="sm" onClick={() => openEditModal(classItem)}>
                        <Pencil className="mr-1 h-3.5 w-3.5" /> Sửa
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => openDeleteAlert(classItem.class_id)}>
                        <Trash2 className="mr-1 h-3.5 w-3.5 text-destructive" /> Xóa
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
            {classes.map((classItem) => (
              <Card key={classItem.class_id}>
                <CardContent className="flex items-start gap-3 py-4">
                  <GraduationCap className="mt-0.5 h-5 w-5 shrink-0 text-muted-foreground" />
                  <div className="min-w-0 flex-1">
                    <p className="font-medium">{classItem.name}</p>
                    <div className="mt-1 flex items-center gap-1 text-sm text-muted-foreground">
                      <Calendar className="h-3 w-3" />
                      <Badge variant="secondary">{classItem.school_year}</Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </>
        )}
      />

      <ActionModal
        isOpen={editModal.isOpen}
        onClose={closeEditModal}
        onConfirm={handleEdit}
        title="Sửa lớp học"
        loading={editLoading}
        confirmText="Lưu"
        cancelText="Hủy"
      >
        <div className="space-y-3">
          <div className="space-y-1.5">
            <Label>Tên lớp</Label>
            <Input value={editData.name} onChange={(e) => handleEditFieldChange("name", e.target.value)} />
          </div>
          <div className="space-y-1.5">
            <Label>Năm học</Label>
            <Input value={editData.school_year} onChange={(e) => handleEditFieldChange("school_year", e.target.value)} />
          </div>
        </div>
      </ActionModal>

      <ConfirmAlertDialog
        isOpen={deleteAlert.isOpen}
        onClose={closeDeleteAlert}
        onConfirm={handleDelete}
        title="Xác nhận xóa lớp"
        description="Bạn có chắc chắn muốn xóa lớp này? Hành động này có thể ảnh hưởng tới học sinh đang thuộc lớp."
        loading={deleteLoading}
        confirmText="Xóa"
        cancelText="Hủy"
      />
    </div>
  );
}
