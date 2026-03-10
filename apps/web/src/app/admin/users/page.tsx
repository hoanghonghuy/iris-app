/**
 * Admin Users Page
 * Quản lý tài khoản: listing, tạo user, khóa/mở khóa.
 * API: GET/POST /admin/users, POST /admin/users/:id/lock|unlock
 */
"use client";

import React, { useEffect, useState, useCallback, useMemo } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { UserInfo, Pagination } from "@/types";
import { PaginationBar } from "@/components/shared/PaginationBar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  UserCog, Loader2, Lock, Unlock, Shield, Mail, Plus, X, AlertCircle, CheckCircle2, Search
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
  const [users, setUsers] = useState<UserInfo[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const [showForm, setShowForm] = useState(false);
  const [formEmail, setFormEmail] = useState("");
  const [formRoles, setFormRoles] = useState<string[]>(["TEACHER"]);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const response = await adminApi.getUsers({ limit: 20, offset: currentOffset });
      const data = (response as any).data || response || [];
      setUsers(Array.isArray(data) ? data : []);
      if (response.pagination) setPagination(response.pagination);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách người dùng");
    } finally { setLoading(false); }
  }, [currentOffset]);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

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

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formEmail.trim()) { setFormError("Email không được trống"); return; }
    if (formRoles.length === 0) { setFormError("Chọn ít nhất 1 vai trò"); return; }
    try {
      setSubmitting(true); setFormError(""); setSuccess("");
      await adminApi.createUser({ email: formEmail, roles: formRoles });
      setSuccess(`Đã tạo user ${formEmail}. User cần kích hoạt tài khoản.`);
      setFormEmail(""); setFormRoles(["TEACHER"]); setShowForm(false); fetchUsers();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo user");
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
        <div className="flex items-center gap-3">
          <UserCog className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Quản lý Người dùng</h1>
        </div>
        <Button size="sm" onClick={() => { setShowForm(!showForm); setSuccess(""); }}>
          {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
          {showForm ? "Hủy" : "Tạo user"}
        </Button>
      </div>

      {success && (
        <Alert><CheckCircle2 className="h-4 w-4 text-green-600" /><AlertDescription>{success}</AlertDescription></Alert>
      )}

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

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}
      
      {/* Toolbar: Search box */}
      {!loading && !error && users.length > 0 && !showForm && (
        <div className="relative max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input 
            type="search" 
            placeholder="Tìm theo email..." 
            className="pl-8 bg-white" 
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      )}

      {loading && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loading && users.length === 0 && !error && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <UserCog className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Chưa có người dùng nào</p>
        </CardContent></Card>
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
                  <tr key={user.user_id} className="border-b last:border-0 hover:bg-zinc-50">
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
          ))}
        </div>
      )}

      {/* Pagination */}
      {!loading && users.length > 0 && (
        <PaginationBar pagination={pagination} onPageChange={setCurrentOffset} />
      )}
    </div>
  );
}
