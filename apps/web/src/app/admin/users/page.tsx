/**
 * Admin Users Page
 * Quản lý tài khoản: listing, tạo user, khóa/mở khóa.
 * API: GET/POST /admin/users, POST /admin/users/:id/lock|unlock
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { ApiResponse, Pagination, UserInfo, UserRole } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Button } from "@/components/ui/button";
import { TableSkeleton } from "@/components/shared/TableSkeleton";
import { CardSkeleton } from "@/components/shared/CardSkeleton";
import { EmptyState } from "@/components/shared/EmptyState";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { toast } from "sonner";
import { Plus, UserCog, X } from "lucide-react";
import {
  CREATABLE_USER_ROLES,
  USER_ROLE_LABELS,
  USER_STATUS_LABEL,
  USER_STATUS_VARIANT,
} from "./config";
import { extractApiErrorMessage } from "./utils";
import { UserCreateForm } from "./components/UserCreateForm";
import { UsersToolbar } from "./components/UsersToolbar";
import { UsersDesktopTable } from "./components/UsersDesktopTable";
import { UsersMobileList } from "./components/UsersMobileList";

const INITIAL_AUTH_ACTION_ALERT_STATE = {
  isOpen: false,
  userId: null,
  action: null,
} as const;

export default function AdminUsersPage() {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState<UserInfo[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);
  const [roleFilter, setRoleFilter] = useState("ALL");

  const [showForm, setShowForm] = useState(false);
  const [formEmail, setFormEmail] = useState("");
  const [formRoles, setFormRoles] = useState<string[]>(["TEACHER"]);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  const [authActionAlert, setAuthActionAlert] = useState<{isOpen: boolean, userId: string | null, action: "lock" | "unlock" | null}>({
    ...INITIAL_AUTH_ACTION_ALERT_STATE,
  });

  const closeAuthActionAlert = useCallback(() => {
    setAuthActionAlert({ ...INITIAL_AUTH_ACTION_ALERT_STATE });
  }, []);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const params: { limit: number; offset: number; role?: UserRole } = { limit: 20, offset: currentOffset };
      if (roleFilter !== "ALL") {
        params.role = roleFilter as UserRole;
      }
      const response = await adminApi.getUsers(params);
      const data = (response as ApiResponse<UserInfo[]>).data || [];
      setUsers(Array.isArray(data) ? data : []);
      if (response.pagination) setPagination(response.pagination);
    } catch (err: unknown) {
      const msg = extractApiErrorMessage(err) || "Không thể tải danh sách người dùng";
      setError(msg);
      toast.error(msg);
    } finally { setLoading(false); }
  }, [currentOffset, roleFilter]);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

  const confirmAuthAction = async () => {
    if (!authActionAlert.userId || !authActionAlert.action) return;
    try {
      setActionLoading(authActionAlert.userId);
      if (authActionAlert.action === "lock") {
        await adminApi.lockUser(authActionAlert.userId);
        toast.success("Đã khóa người dùng");
      } else {
        await adminApi.unlockUser(authActionAlert.userId);
        toast.success("Đã mở khóa người dùng");
      }
      fetchUsers();
    } catch (err: unknown) {
      toast.error(extractApiErrorMessage(err) || `Không thể ${authActionAlert.action === "lock" ? "khóa" : "mở khóa"}`);
    } finally {
      setActionLoading(null);
      closeAuthActionAlert();
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formEmail.trim()) { setFormError("Email không được trống"); return; }
    if (formRoles.length === 0) { setFormError("Chọn ít nhất 1 vai trò"); return; }
    try {
      setSubmitting(true); setFormError("");
      await adminApi.createUser({ email: formEmail, roles: formRoles });
      toast.success(`Đã tạo user ${formEmail}. User cần kích hoạt tài khoản.`);
      setFormEmail(""); setFormRoles(["TEACHER"]); setShowForm(false); fetchUsers();
    } catch (err: unknown) {
      const message = extractApiErrorMessage(err) || "Không thể tạo user";
      setFormError(message);
      toast.error(message);
    } finally { setSubmitting(false); }
  };

  const toggleRole = (role: string) => {
    setFormRoles((prev) => prev.includes(role) ? prev.filter((r) => r !== role) : [...prev, role]);
  };

  const filteredUsers = useMemo(() => {
    if (!searchQuery.trim()) return users;
    const q = searchQuery.toLowerCase();
    return users.filter((u) => u.email?.toLowerCase().includes(q));
  }, [users, searchQuery]);

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Button size="sm" onClick={() => { setShowForm(!showForm); }}>
          {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
          {showForm ? "Hủy" : "Tạo user"}
        </Button>
      </div>

      {showForm && (
        <UserCreateForm
          formError={formError}
          formEmail={formEmail}
          formRoles={formRoles}
          submitting={submitting}
          creatableRoles={CREATABLE_USER_ROLES}
          roleLabels={USER_ROLE_LABELS}
          onEmailChange={setFormEmail}
          onToggleRole={toggleRole}
          onSubmit={handleCreate}
        />
      )}



      {/* Toolbar: Search box & Role Filter */}
      {!loading && !error && (users.length > 0 || roleFilter !== "ALL") && !showForm && (
        <UsersToolbar
          searchQuery={searchQuery}
          roleFilter={roleFilter}
          onSearchChange={setSearchQuery}
          onRoleFilterChange={(value) => {
            setRoleFilter(value);
            setCurrentOffset(0);
          }}
        />
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

      {!loading && users.length === 0 && roleFilter === "ALL" && !error && (
        <EmptyState
          icon={UserCog}
          title="Chưa có người dùng nào"
          description="Hiện tại hệ thống chưa có dữ liệu người dùng mới."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Tạo user đầu tiên
            </Button>
          }
        />
      )}

      {!loading && users.length === 0 && roleFilter !== "ALL" && !error && (
        <div className="rounded-lg border border-dashed p-8 text-center mt-4">
          <p className="text-sm text-muted-foreground">Không tìm thấy người dùng nào với vai trò này.</p>
        </div>
      )}

      {!loading && users.length > 0 && filteredUsers.length === 0 && (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <p className="text-sm text-muted-foreground">Không tìm thấy người dùng nào mang email &ldquo;{searchQuery}&rdquo;</p>
        </div>
      )}

      {/* Desktop Table */}
      {!loading && filteredUsers.length > 0 && (
        <UsersDesktopTable
          users={filteredUsers}
          actionLoading={actionLoading}
          currentUserId={currentUser?.user_id}
          roleLabels={USER_ROLE_LABELS}
          statusLabels={USER_STATUS_LABEL}
          statusVariants={USER_STATUS_VARIANT}
          onRequestLock={(userId) => setAuthActionAlert({ isOpen: true, userId, action: "lock" })}
          onRequestUnlock={(userId) => setAuthActionAlert({ isOpen: true, userId, action: "unlock" })}
        />
      )}

      {/* Mobile Cards */}
      {!loading && filteredUsers.length > 0 && (
        <UsersMobileList
          users={filteredUsers}
          actionLoading={actionLoading}
          currentUserId={currentUser?.user_id}
          roleLabels={USER_ROLE_LABELS}
          statusLabels={USER_STATUS_LABEL}
          statusVariants={USER_STATUS_VARIANT}
          onRequestLock={(userId) => setAuthActionAlert({ isOpen: true, userId, action: "lock" })}
          onRequestUnlock={(userId) => setAuthActionAlert({ isOpen: true, userId, action: "unlock" })}
        />
      )}

      {/* Pagination */}
      {!loading && users.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}

      {/* Lock/Unlock Confirmation */}
      <ConfirmAlertDialog
        isOpen={authActionAlert.isOpen}
        onClose={closeAuthActionAlert}
        onConfirm={confirmAuthAction}
        title={authActionAlert.action === "lock" ? "Xác nhận khóa tài khoản" : "Xác nhận mở khóa tài khoản"}
        description={authActionAlert.action === "lock" ? "Việc khóa tài khoản sẽ ngay lập tức vô hiệu hóa các phiên đăng nhập của người dùng này và ngăn họ truy cập vào hệ thống. Bạn có chắc chắn muốn khóa?" : "Tài khoản sẽ có thể đăng nhập lại bình thường sau khi được mở khóa. Bạn có tiếp tục?"}
        loading={!!actionLoading}
        confirmText={authActionAlert.action === "lock" ? "Khóa tài khoản" : "Mở khóa tài khoản"}
      />
    </div>
  );
}
