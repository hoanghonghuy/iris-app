/**
 * Admin Parents Page
 * Danh sách phụ huynh + gán/hủy gán học sinh.
 * API: GET /admin/parents, POST/DELETE /admin/parents/:id/students/:student_id
 */
"use client";

import React from "react";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { ActionModal } from "@/components/shared/ActionModal";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { Heart, Loader2, Phone, Mail, Link2, AlertCircle, CheckCircle2, Search, X } from "lucide-react";
import { useAuth } from "@/providers/AuthProvider";
import { useAdminParentsPage } from "./useAdminParentsPage";

export default function AdminParentsPage() {
  const { role } = useAuth();
  const {
    parents,
    searchQuery,
    loading,
    error,
    pagination,
    schools,
    classes,
    students,
    selectedSchoolId,
    selectedClassId,
    selectedStudentId,
    actionLoading,
    success,
    assignModal,
    unassignAlert,
    filteredParents,
    setSearchQuery,
    setCurrentOffset,
    setSelectedSchoolId,
    setSelectedClassId,
    setSelectedStudentId,
    setAssignModal,
    setUnassignAlert,
    handleAssign,
    confirmUnassign,
  } = useAdminParentsPage();

  return (
    <div className="space-y-6">
      {success && <Alert><CheckCircle2 className="h-4 w-4 text-success" /><AlertDescription>{success}</AlertDescription></Alert>}
      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}

      {/* Toolbar: Search box */}
      {!loading && !error && parents.length > 0 && (
        <div className="relative max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input 
            type="search" 
            placeholder="Tìm theo tên, email, SĐT..." 
            className="pl-8 bg-background " 
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      )}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && parents.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <Heart className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có phụ huynh nào</p>
        </CardContent></Card>
      )}

      {!loading && parents.length > 0 && filteredParents.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy phụ huynh nào phù hợp với &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      {/* Desktop Table */}
      {!loading && filteredParents.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Họ tên</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Học Sinh Quản Lý</th>
                  <th className="px-6 py-3 font-medium text-right">Gán học sinh</th>
                </tr>
              </thead>
              <tbody>
                {filteredParents.map((p) => (
                  <tr key={p.parent_id} className="border-b last:border-0 hover:bg-muted leading-relaxed">
                    <td className="px-6 py-4">
                      <div className="font-medium text-foreground">{p.full_name}</div>
                      <div className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                        <Phone className="h-3 w-3" /> {p.phone || "—"}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-muted-foreground">{p.email}</td>
                    <td className="px-6 py-4">
                      {p.children && p.children.length > 0 ? (
                        <div className="flex flex-wrap gap-1.5">
                          {p.children.map(c => (
                            <Badge key={c.student_id} variant="secondary" className="pr-1.5 flex items-center gap-1">
                              {c.full_name}
                              <button 
                                onClick={() => setUnassignAlert({ isOpen: true, parentId: p.parent_id, studentId: c.student_id, studentName: c.full_name })}
                                className="ml-0.5 rounded-full p-0.5 hover:bg-destructive/20 hover:text-destructive transition-colors focus:outline-none focus:ring-2 focus:ring-ring"
                                aria-label="Remove"
                              >
                                <X className="h-3 w-3" />
                              </button>
                            </Badge>
                          ))}
                        </div>
                      ) : (
                        <span className="text-sm text-muted-foreground italic">Chưa ghép học sinh</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-right align-middle">
                        <Button variant="ghost" size="sm" onClick={() => setAssignModal({ isOpen: true, parentId: p.parent_id, parentName: p.full_name })}>
                          <Link2 className="mr-1 h-4 w-4" /> Gán HS
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
      {!loading && filteredParents.length > 0 && (
        <div className="space-y-3 md:hidden">
          {filteredParents.map((p) => (
            <Card key={p.parent_id}>
              <CardContent className="py-4">
                <p className="font-medium">{p.full_name}</p>
                <div className="mt-2 space-y-1 text-sm text-muted-foreground">
                  <p className="flex items-center gap-2"><Mail className="h-3 w-3" /> {p.email}</p>
                  {p.phone && <p className="flex items-center gap-2"><Phone className="h-3 w-3" /> {p.phone}</p>}
                </div>

                <div className="mt-4 border-t border-dashed pt-3">
                  <p className="text-xs text-muted-foreground mb-2 font-medium uppercase tracking-wider">Học sinh thuộc quản lý</p>
                  {p.children && p.children.length > 0 ? (
                    <div className="flex flex-wrap gap-1.5">
                      {p.children.map(c => (
                        <Badge key={c.student_id} variant="secondary" className="pl-2 pr-1.5 py-0.5 flex items-center gap-1">
                          {c.full_name}
                          <button 
                            onClick={(e: React.MouseEvent<HTMLButtonElement>) => { e.preventDefault(); setUnassignAlert({ isOpen: true, parentId: p.parent_id, studentId: c.student_id, studentName: c.full_name }); }}
                            className="ml-1 rounded-full bg-transparent p-0.5 text-muted-foreground hover:bg-destructive/20 hover:text-destructive transition-colors outline-none"
                            aria-label="Remove child"
                          >
                            <X className="h-3 w-3" />
                          </button>
                        </Badge>
                      ))}
                    </div>
                  ) : (
                    <div className="text-sm text-muted-foreground italic">
                      Chưa ghép học sinh
                    </div>
                  )}
                </div>

                <div className="mt-4 pt-4 border-t border-border/60 flex items-center justify-start">
                    <Button variant="secondary" size="sm" className="w-full hover:bg-primary/20 hover:text-primary" onClick={() => setAssignModal({ isOpen: true, parentId: p.parent_id, parentName: p.full_name })}>
                      <Link2 className="mr-1 h-4 w-4" /> Gán học sinh
                    </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Pagination */}
      {!loading && parents.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}

      {/* Assign Modal */}
      <ActionModal
        isOpen={assignModal.isOpen}
        onClose={() => setAssignModal({ isOpen: false, parentId: null, parentName: null })}
        onConfirm={handleAssign}
        title="Gán học sinh"
        description={<>Chọn trường, lớp và học sinh để gán cho phụ huynh <strong>{assignModal.parentName}</strong>.</>}
        loading={actionLoading}
        disabled={!selectedStudentId}
        confirmText="Gán học sinh"
      >
        <div className="grid gap-4 py-2">
          {role === 'SUPER_ADMIN' && (
            <div className="flex flex-wrap items-center gap-2">
            <Select value={selectedSchoolId} onValueChange={setSelectedSchoolId}>
              <SelectTrigger className="w-[180px]"><SelectValue placeholder="Chọn trường" /></SelectTrigger>
              <SelectContent>
                {schools.map((s) => <SelectItem key={s.school_id} value={s.school_id}>{s.name}</SelectItem>)}
              </SelectContent>
            </Select>
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
          <div className="space-y-2">
            <label className="text-sm font-medium">Học sinh</label>
            <Select value={selectedStudentId} onValueChange={setSelectedStudentId}>
              <SelectTrigger><SelectValue placeholder="Chọn học sinh" /></SelectTrigger>
              <SelectContent>
                {students.map((s) => <SelectItem key={s.student_id} value={s.student_id}>{s.full_name}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>
        </div>
      </ActionModal>

      {/* Unassign Alert */}
      <ConfirmAlertDialog
        isOpen={unassignAlert.isOpen}
        onClose={() => setUnassignAlert({ isOpen: false, parentId: null, studentId: null, studentName: null })}
        onConfirm={confirmUnassign}
        title="Xác nhận hủy gán"
        description={<>Bạn có chắc chắn muốn hủy gán học sinh <strong>{unassignAlert.studentName}</strong> khỏi phụ huynh này?</>}
        loading={actionLoading}
        confirmText="Chắc chắn hủy"
      />
    </div>
  );
}
