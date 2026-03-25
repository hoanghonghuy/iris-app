/**
 * Admin Users Page
 * Quản lý tài khoản: listing, tạo user, khóa/mở khóa.
 * API: GET/POST /admin/users, POST /admin/users/:id/lock|unlock
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { UserInfo, Pagination } from "@/types";
import { useAuth } from "@/providers/AuthProvider";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { TableSkeleton } from "@/components/shared/TableSkeleton";
import { CardSkeleton } from "@/components/shared/CardSkeleton";
import { EmptyState } from "@/components/shared/EmptyState";
import { ConfirmAlertDialog } from "@/components/shared/ConfirmAlertDialog";
import { toast } from "sonner";
import {
  Loader2, Lock, Unlock, Shield, Mail, Plus, X, AlertCircle, Search, UserCog
} from "lucide-react";

const roleLabels: Record<string, string> = {
  SUPER_ADMIN: "Super Admin", SCHOOL_ADMIN: "School Admin",
  TEACHER: "Giáo viên", PARENT: "Phụ huynh",
};
const allRoles = ["TEACHER", "PARENT", "SCHOOL_ADMIN"];

const statusVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  active: "default", pending: "secondary", locked: "destructive",
};
const statusLabel: Record<string, string> = {
  active: "Hoạt động", pending: "Chờ kích hoạt", locked: "Đã khóa",
};

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

  const [authActionAlert, setAuthActionAlert] = useState<{isOpen: boolean, userId: string | null, action: "lock" | "unlock" | null}>({isOpen: false, userId: null, action: null});

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const params: any = { limit: 20, offset: currentOffset };
      if (roleFilter !== "ALL") {
        params.role = roleFilter;
      }
      const response = await adminApi.getUsers(params);
      const data = (response as any).data || response || [];
      setUsers(Array.isArray(data) ? data : []);
      if (response.pagination) setPagination(response.pagination);
    } catch (err: any) {
      const msg = err.response?.data?.error || "Không thể tải danh sách người dùng";
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
    } catch (err: any) {
      toast.error(err.response?.data?.error || `Không thể ${authActionAlert.action === "lock" ? "khóa" : "mở khóa"}`);
    } finally {
      setActionLoading(null);
      setAuthActionAlert({ isOpen: false, userId: null, action: null });
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
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo user");
      toast.error(err.response?.data?.error || "Không thể tạo user");
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
        <Card>
          <CardHeader><CardTitle className="text-lg">Tạo tài khoản mới</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{formError}</AlertDescription></Alert>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="userEmail">Email <span className="text-destructive">*</span></Label>
                  <Input id="userEmail" type="email" placeholder="user@example.com" value={formEmail} onChange={(e) => setFormEmail(e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label>Vai trò <span className="text-destructive">*</span></Label>
                  <div className="flex flex-wrap gap-2">
                    {allRoles.map((role) => (
                      <Badge key={role} variant={formRoles.includes(role) ? "default" : "outline"}
                        className="cursor-pointer select-none" onClick={() => toggleRole(role)}>
                        {roleLabels[role]}
                      </Badge>
                    ))}
                  </div>
                </div>
              </div>
              <p className="text-xs text-muted-foreground">User sẽ ở trạng thái &ldquo;Chờ kích hoạt&rdquo;. Họ cần dùng activation token để đặt mật khẩu.</p>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
                  {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Tạo user
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}



      {/* Toolbar: Search box & Role Filter */}
      {!loading && !error && (users.length > 0 || roleFilter !== "ALL") && !showForm && (
        <div className="flex items-center gap-3">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Tìm theo email..."
              className="pl-8 bg-background min-w-0"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          
          <Select 
            value={roleFilter} 
            onValueChange={(val) => { 
              setRoleFilter(val); 
              setCurrentOffset(0); 
            }}
          >
            <SelectTrigger className="w-[140px] shrink-0">
              <SelectValue placeholder="Tất cả vai trò" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">Tất cả vai trò</SelectItem>
              <SelectItem value="TEACHER">Giáo viên</SelectItem>
              <SelectItem value="PARENT">Phụ huynh</SelectItem>
              <SelectItem value="SCHOOL_ADMIN">School Admin</SelectItem>
              <SelectItem value="SUPER_ADMIN">Super Admin</SelectItem>
            </SelectContent>
          </Select>
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
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Vai trò</th>
                  <th className="px-6 py-3 font-medium">Trạng thái</th>
                  <th className="px-6 py-3 font-medium text-right">Hành động</th>
                </tr>
              </thead>
              <tbody>
                {filteredUsers.map((user) => (
                  <tr key={user.user_id} className="border-b last:border-0 hover:bg-muted">
                    <td className="px-6 py-4 font-medium">{user.email}</td>
                    <td className="px-6 py-4">
                      <div className="flex flex-wrap gap-1">
                        {user.roles?.map((r) => <Badge key={r} variant="secondary"><Shield className="h-3 w-3" /> {roleLabels[r] || r}</Badge>)}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <Badge variant={statusVariant[user.status] || "secondary"}>{statusLabel[user.status] || user.status}</Badge>
                    </td>
                    <td className="px-6 py-4 text-right">
                      {user.status === "active" ? (
                        <Button variant="ghost" size="sm" onClick={() => setAuthActionAlert({ isOpen: true, userId: user.user_id, action: "lock" })} disabled={actionLoading === user.user_id || currentUser?.user_id === user.user_id} title={currentUser?.user_id === user.user_id ? "Bạn không thể tự khóa chính mình" : ""}>
                          {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Lock className="mr-1 h-4 w-4 text-destructive" />} <span className="text-destructive">Khóa</span>
                        </Button>
                      ) : user.status === "locked" ? (
                        <Button variant="ghost" size="sm" onClick={() => setAuthActionAlert({ isOpen: true, userId: user.user_id, action: "unlock" })} disabled={actionLoading === user.user_id}>
                          {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Unlock className="mr-1 h-4 w-4 text-success" />} <span className="text-success">Mở khóa</span>
                        </Button>
                      ) : null}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent></Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loading && filteredUsers.length > 0 && (
        <div className="space-y-3 md:hidden">
          {filteredUsers.map((user) => (
            <Card key={user.user_id}>
              <CardContent className="py-4">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0 flex-1">
                    <p className="flex items-center gap-2 font-medium"><Mail className="h-4 w-4 shrink-0 text-muted-foreground" /><span className="truncate">{user.email}</span></p>
                    <div className="mt-2 flex flex-wrap gap-1">
                      {user.roles?.map((r) => <Badge key={r} variant="secondary"><Shield className="h-3 w-3" /> {roleLabels[r] || r}</Badge>)}
                      <Badge variant={statusVariant[user.status] || "secondary"}>{statusLabel[user.status] || user.status}</Badge>
                    </div>
                  </div>
                  <div>
                    {user.status === "active" ? (
                      <Button variant="ghost" size="sm" onClick={() => setAuthActionAlert({ isOpen: true, userId: user.user_id, action: "lock" })} disabled={actionLoading === user.user_id || currentUser?.user_id === user.user_id} title={currentUser?.user_id === user.user_id ? "Bạn không thể tự khóa chính mình" : ""}>
                        {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Lock className="h-4 w-4 text-destructive" />}
                      </Button>
                    ) : user.status === "locked" ? (
                      <Button variant="ghost" size="sm" onClick={() => setAuthActionAlert({ isOpen: true, userId: user.user_id, action: "unlock" })} disabled={actionLoading === user.user_id}>
                        {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Unlock className="h-4 w-4 text-success" />}
                      </Button>
                    ) : null}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Pagination */}
      {!loading && users.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}

      {/* Lock/Unlock Confirmation */}
      <ConfirmAlertDialog
        isOpen={authActionAlert.isOpen}
        onClose={() => setAuthActionAlert({ isOpen: false, userId: null, action: null })}
        onConfirm={confirmAuthAction}
        title={authActionAlert.action === "lock" ? "Xác nhận khóa tài khoản" : "Xác nhận mở khóa tài khoản"}
        description={authActionAlert.action === "lock" ? "Việc khóa tài khoản sẽ ngay lập tức vô hiệu hóa các phiên đăng nhập của người dùng này và ngăn họ truy cập vào hệ thống. Bạn có chắc chắn muốn khóa?" : "Tài khoản sẽ có thể đăng nhập lại bình thường sau khi được mở khóa. Bạn có tiếp tục?"}
        loading={!!actionLoading}
        confirmText={authActionAlert.action === "lock" ? "Khóa tài khoản" : "Mở khóa tài khoản"}
      />
    </div>
  );
}
