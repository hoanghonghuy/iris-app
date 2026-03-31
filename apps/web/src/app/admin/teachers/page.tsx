/**
 * Admin Teachers Page
 * Danh sách giáo viên + gán/hủy gán lớp.
 * API: GET /admin/teachers, POST/DELETE /admin/teachers/:id/classes/:class_id
 */
"use client";

import { PaginationBar } from "@/components/shared/PaginationBar";
import { TableSkeleton } from "@/components/shared/TableSkeleton";
import { CardSkeleton } from "@/components/shared/CardSkeleton";
import { EmptyState } from "@/components/shared/EmptyState";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { ActionModal } from "@/components/shared/ActionModal";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { toast } from "sonner";
import { Phone, Mail, Link2, Search, X, BookUser } from "lucide-react";
import { useAuth } from "@/providers/AuthProvider";
import { Input } from "@/components/ui/input"; // Added missing import for Input
import { extractApiErrorMessage } from "@/lib/api-error";
import { useAdminTeachersPage } from "./useAdminTeachersPage";

export default function AdminTeachersPage() {
  const { role } = useAuth();
  const {
    teachers,
    searchQuery,
    loading,
    error,
    pagination,
    schools,
    classes,
    selectedSchoolId,
    selectedClassId,
    actionLoading,
    assignModal,
    unassignAlert,
    filteredTeachers,
    setSearchQuery,
    setCurrentOffset,
    setSelectedSchoolId,
    setSelectedClassId,
    setActionLoading,
    setAssignModal,
    setUnassignAlert,
    closeAssignModal,
    closeUnassignAlert,
    handleAssign,
    confirmUnassign,
  } = useAdminTeachersPage();

  const handleAssignWithFeedback = async () => {
    if (!selectedClassId || !assignModal.teacherId) {
      return;
    }

    try {
      setActionLoading(true);
      await handleAssign();
      const className = classes.find((classInfo) => classInfo.class_id === selectedClassId)?.name || "";
      toast.success(`Đã gán giáo viên vào lớp ${className}`);
    } catch (err: unknown) {
      toast.error(extractApiErrorMessage(err, "Không thể gán lớp"));
    } finally {
      setActionLoading(false);
    }
  };

  const confirmUnassignWithFeedback = async () => {
    if (!unassignAlert.teacherId || !unassignAlert.classId) {
      return;
    }

    try {
      setActionLoading(true);
      await confirmUnassign();
      toast.success(`Đã hủy gán lớp ${unassignAlert.className}`);
    } catch (err: unknown) {
      toast.error(extractApiErrorMessage(err, "Không thể hủy gán lớp"));
    } finally {
      setActionLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Toolbar: Search */}
      {!loading && !error && teachers.length > 0 && (
        <div className="relative max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Tìm theo tên, email, SĐT..."
            className="pl-8 bg-background"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      )}

      {loading && (
        <>
          <div className="hidden md:block">
            <TableSkeleton columns={4} rows={10} />
          </div>
          <div className="md:hidden">
            <CardSkeleton cards={5} />
          </div>
        </>
      )}

      {!loading && teachers.length === 0 && !error && (
        <EmptyState
          icon={BookUser}
          title="Chưa có giáo viên nào"
          description="Hiện tại hệ thống chưa có dữ liệu giáo viên mới."
        />
      )}

      {!loading && teachers.length > 0 && filteredTeachers.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy giáo viên nào phù hợp với &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      {/* Desktop Table */}
      {!loading && filteredTeachers.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Họ tên</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Lớp Phụ Trách</th>
                  <th className="px-6 py-3 font-medium text-right">Gán lớp</th>
                </tr>
              </thead>
              <tbody>
                {filteredTeachers.map((t) => (
                  <tr key={t.teacher_id} className="border-b last:border-0 hover:bg-muted leading-relaxed">
                    <td className="px-6 py-4">
                      <div className="font-medium text-foreground">{t.full_name}</div>
                      <div className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                        <Phone className="h-3 w-3" /> {t.phone || "—"}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-muted-foreground">{t.email}</td>
                    <td className="px-6 py-4">
                      {t.classes && t.classes.length > 0 ? (
                        <div className="flex flex-wrap gap-1.5">
                          {t.classes.map(c => (
                            <Badge key={c.class_id} variant="secondary" className="pr-1.5 flex items-center gap-1">
                              {c.name}
                              <button
                                onClick={() => setUnassignAlert({ isOpen: true, teacherId: t.teacher_id, classId: c.class_id, className: c.name })}
                                className="ml-0.5 rounded-full p-0.5 hover:bg-destructive/20 hover:text-destructive transition-colors focus:outline-none focus:ring-2 focus:ring-ring"
                                aria-label="Remove"
                              >
                                <X className="h-3 w-3" />
                              </button>
                            </Badge>
                          ))}
                        </div>
                      ) : (
                        <span className="text-sm text-muted-foreground italic">Chưa phân lớp</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-right align-middle">
                        <Button variant="ghost" size="sm" onClick={() => setAssignModal({ isOpen: true, teacherId: t.teacher_id, teacherName: t.full_name })}>
                          <Link2 className="mr-1 h-4 w-4" /> Gán lớp
                        </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent></Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loading && filteredTeachers.length > 0 && (
        <div className="space-y-3 md:hidden">
          {filteredTeachers.map((t) => (
            <Card key={t.teacher_id}>
              <CardContent className="py-4">
                <p className="font-medium">{t.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2"><Mail className="h-3 w-3" /> {t.email}</p>
                  {t.phone && <p className="flex items-center gap-2"><Phone className="h-3 w-3" /> {t.phone}</p>}
                </div>

                <div className="mt-4 border-t border-dashed pt-3">
                  <p className="text-xs text-muted-foreground mb-2 font-medium uppercase tracking-wider">Lớp phụ trách</p>
                  {t.classes && t.classes.length > 0 ? (
                    <div className="flex flex-wrap gap-1.5">
                      {t.classes.map(c => (
                        <Badge key={c.class_id} variant="secondary" className="pl-2 pr-1.5 py-0.5 flex items-center gap-1">
                          {c.name}
                            <button
                              onClick={(e) => { e.preventDefault(); setUnassignAlert({ isOpen: true, teacherId: t.teacher_id, classId: c.class_id, className: c.name }); }}
                              className="ml-1 rounded-full bg-transparent p-0.5 text-muted-foreground hover:bg-destructive/20 hover:text-destructive transition-colors outline-none"
                            aria-label="Remove class"
                          >
                            <X className="h-3 w-3" />
                          </button>
                        </Badge>
                      ))}
                    </div>
                  ) : (
                    <div className="text-sm text-muted-foreground italic">
                      Chưa phân lớp
                    </div>
                  )}
                </div>

                <div className="mt-4 pt-4 border-t border-border/60 flex items-center justify-start">
                    <Button variant="secondary" size="sm" className="w-full hover:bg-primary/20 hover:text-primary" onClick={() => setAssignModal({ isOpen: true, teacherId: t.teacher_id, teacherName: t.full_name })}>
                      <Link2 className="mr-1 h-4 w-4" /> Gán phân lớp
                    </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Pagination */}
      {!loading && teachers.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}

      {/* Assign Modal */}
      <ActionModal
        isOpen={assignModal.isOpen}
        onClose={closeAssignModal}
        onConfirm={handleAssignWithFeedback}
        title="Gán lớp phụ trách"
        description={<>Chọn trường và lớp để gán cho giáo viên <strong>{assignModal.teacherName}</strong>.</>}
        loading={actionLoading}
        disabled={!selectedClassId}
        confirmText="Gán lớp"
      >
        <div className="grid gap-4 py-2">
          {role === 'SUPER_ADMIN' && (
            <div className="space-y-2">
              <div className="flex flex-wrap items-center gap-2">
                <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
                  <SelectTrigger className="w-[180px]"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
                  <SelectContent>
                    {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
            </div>
          )}
          <div className="space-y-2">
            <label className="text-sm font-medium">Lớp học</label>
            <Select value={selectedClassId} onValueChange={setSelectedClassId}>
              <SelectTrigger><SelectValue placeholder="Chọn lớp" /></SelectTrigger>
              <SelectContent>
                {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>
        </div>
      </ActionModal>

      {/* Unassign Alert */}
      <ConfirmAlertDialog
        isOpen={unassignAlert.isOpen}
        onClose={closeUnassignAlert}
        onConfirm={confirmUnassignWithFeedback}
        title="Xác nhận hủy gán"
        description={<>Bạn có chắc chắn muốn hủy gán lớp <strong>{unassignAlert.className}</strong> khỏi giáo viên này? Hành động này sẽ thay đổi quyền truy cập của giáo viên đối với lớp học.</>}
        loading={actionLoading}
        confirmText="Chắc chắn hủy"
      />
    </div>
  );
}
