/**
 * Admin Users Page
 * Quản lý tài khoản: listing, tạo user, khóa/mở khóa, gán role.
 * API: GET/POST /admin/users, POST /admin/users/:id/lock|unlock|roles
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { UserInfo } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  UserCog, Loader2, Lock, Unlock, Shield, Mail, Plus, X, ChevronDown,
} from "lucide-react";

const roleLabels: Record<string, string> = {
  SUPER_ADMIN: "Super Admin", SCHOOL_ADMIN: "School Admin",
  TEACHER: "Giáo viên", PARENT: "Phụ huynh",
};
const allRoles = ["TEACHER", "PARENT", "SCHOOL_ADMIN"];

const statusConfig: Record<string, { label: string; className: string }> = {
  active: { label: "Hoạt động", className: "bg-green-100 text-green-700" },
  pending: { label: "Chờ kích hoạt", className: "bg-yellow-100 text-yellow-700" },
  locked: { label: "Đã khóa", className: "bg-red-100 text-red-700" },
};

export default function AdminUsersPage() {
  const [users, setUsers] = useState<UserInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  // Create form
  const [showForm, setShowForm] = useState(false);
  const [formEmail, setFormEmail] = useState("");
  const [formRoles, setFormRoles] = useState<string[]>(["TEACHER"]);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const response = await adminApi.getUsers({ limit: 100 });
      const data = (response as any).data || response || [];
      setUsers(Array.isArray(data) ? data : []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách người dùng");
    } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

  // Lock / Unlock
  const handleLock = async (userId: string) => {
    try { setActionLoading(userId); await adminApi.lockUser(userId); fetchUsers(); }
    catch (err: any) { setError(err.response?.data?.error || "Không thể khóa"); }
    finally { setActionLoading(null); }
  };
  const handleUnlock = async (userId: string) => {
    try { setActionLoading(userId); await adminApi.unlockUser(userId); fetchUsers(); }
    catch (err: any) { setError(err.response?.data?.error || "Không thể mở khóa"); }
    finally { setActionLoading(null); }
  };

  // Create User
  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formEmail.trim()) { setFormError("Email không được trống"); return; }
    if (formRoles.length === 0) { setFormError("Chọn ít nhất 1 vai trò"); return; }
    try {
      setSubmitting(true); setFormError(""); setSuccess("");
      await adminApi.createUser({ email: formEmail, roles: formRoles });
      setSuccess(`Đã tạo user ${formEmail}. User cần kích hoạt tài khoản.`);
      setFormEmail(""); setFormRoles(["TEACHER"]); setShowForm(false);
      fetchUsers();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo user");
    } finally { setSubmitting(false); }
  };

  const toggleRole = (role: string) => {
    setFormRoles((prev) =>
      prev.includes(role) ? prev.filter((r) => r !== role) : [...prev, role]
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <UserCog className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Quản lý Người dùng</h1>
        </div>
        <Button size="sm" onClick={() => { setShowForm(!showForm); setSuccess(""); }}>
          {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
          {showForm ? "Hủy" : "Tạo user"}
        </Button>
      </div>

      {/* Success */}
      {success && <div className="rounded-md bg-green-100 p-4 text-sm text-green-700">{success}</div>}

      {/* Create Form */}
      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Tạo tài khoản mới</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{formError}</div>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Email <span className="text-destructive">*</span></label>
                  <Input type="email" placeholder="user@example.com" value={formEmail} onChange={(e) => setFormEmail(e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Vai trò <span className="text-destructive">*</span></label>
                  <div className="flex flex-wrap gap-2">
                    {allRoles.map((role) => (
                      <button key={role} type="button" onClick={() => toggleRole(role)}
                        className={`rounded-full px-3 py-1 text-xs font-medium transition-colors ${
                          formRoles.includes(role) ? "bg-zinc-800 text-white" : "bg-zinc-100 text-zinc-500 hover:bg-zinc-200"
                        }`}>
                        {roleLabels[role]}
                      </button>
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

      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && users.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <UserCog className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có người dùng nào</p>
        </CardContent></Card>
      )}

      {/* Desktop Table */}
      {!loading && users.length > 0 && (
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
                {users.map((user) => {
                  const status = statusConfig[user.status] || statusConfig.pending;
                  return (
                    <tr key={user.user_id} className="border-b last:border-0 hover:bg-zinc-50">
                      <td className="px-6 py-4 font-medium">{user.email}</td>
                      <td className="px-6 py-4">
                        <div className="flex flex-wrap gap-1">
                          {user.roles?.map((r) => (
                            <span key={r} className="inline-flex items-center gap-1 rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs font-medium">
                              <Shield className="h-3 w-3" /> {roleLabels[r] || r}
                            </span>
                          ))}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <span className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${status.className}`}>{status.label}</span>
                      </td>
                      <td className="px-6 py-4 text-right">
                        {user.status === "active" ? (
                          <Button variant="ghost" size="sm" onClick={() => handleLock(user.user_id)} disabled={actionLoading === user.user_id}>
                            {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Lock className="mr-1 h-4 w-4" />} Khóa
                          </Button>
                        ) : user.status === "locked" ? (
                          <Button variant="ghost" size="sm" onClick={() => handleUnlock(user.user_id)} disabled={actionLoading === user.user_id}>
                            {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Unlock className="mr-1 h-4 w-4" />} Mở khóa
                          </Button>
                        ) : null}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </CardContent></Card>
        </div>
      )}

      {/* Mobile Cards */}
      {!loading && users.length > 0 && (
        <div className="space-y-3 md:hidden">
          {users.map((user) => {
            const status = statusConfig[user.status] || statusConfig.pending;
            return (
              <Card key={user.user_id}>
                <CardContent className="py-4">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0 flex-1">
                      <p className="flex items-center gap-2 font-medium">
                        <Mail className="h-4 w-4 shrink-0 text-muted-foreground" />
                        <span className="truncate">{user.email}</span>
                      </p>
                      <div className="mt-2 flex flex-wrap gap-1">
                        {user.roles?.map((r) => (
                          <span key={r} className="inline-flex items-center gap-1 rounded-full bg-zinc-100 px-2 py-0.5 text-xs">
                            <Shield className="h-3 w-3" /> {roleLabels[r] || r}
                          </span>
                        ))}
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${status.className}`}>{status.label}</span>
                      </div>
                    </div>
                    <div>
                      {user.status === "active" ? (
                        <Button variant="ghost" size="sm" onClick={() => handleLock(user.user_id)} disabled={actionLoading === user.user_id}>
                          {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Lock className="h-4 w-4" />}
                        </Button>
                      ) : user.status === "locked" ? (
                        <Button variant="ghost" size="sm" onClick={() => handleUnlock(user.user_id)} disabled={actionLoading === user.user_id}>
                          {actionLoading === user.user_id ? <Loader2 className="h-4 w-4 animate-spin" /> : <Unlock className="h-4 w-4" />}
                        </Button>
                      ) : null}
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
