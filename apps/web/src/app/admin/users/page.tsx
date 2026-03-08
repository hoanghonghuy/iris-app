/**
 * Admin Users Page
 * Quản lý tài khoản: xem danh sách, khóa/mở khóa user.
 * API: GET /admin/users, POST /admin/users/:id/lock, POST /admin/users/:id/unlock
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { UserInfo, UserRole } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  UserCog,
  Loader2,
  Lock,
  Unlock,
  Shield,
  Mail,
  User,
} from "lucide-react";

// Role labels
const roleLabels: Record<string, string> = {
  SUPER_ADMIN: "Super Admin",
  SCHOOL_ADMIN: "School Admin",
  TEACHER: "Giáo viên",
  PARENT: "Phụ huynh",
};

// Status labels & styles
const statusConfig: Record<string, { label: string; className: string }> = {
  active: { label: "Hoạt động", className: "bg-green-100 text-green-700" },
  pending: { label: "Chờ kích hoạt", className: "bg-yellow-100 text-yellow-700" },
  locked: { label: "Đã khóa", className: "bg-red-100 text-red-700" },
};

export default function AdminUsersPage() {
  // ─── State ────────────────────────────────────────────────────────

  const [users, setUsers] = useState<UserInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  // ─── Fetch users ──────────────────────────────────────────────────

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const response = await adminApi.getUsers({ limit: 100 });
      // response wrapped: { data: [...], pagination: {...} }
      const data = (response as any).data || response || [];
      setUsers(Array.isArray(data) ? data : []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách người dùng");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  // ─── Lock / Unlock ────────────────────────────────────────────────

  const handleLock = async (userId: string) => {
    try {
      setActionLoading(userId);
      await adminApi.lockUser(userId);
      fetchUsers();
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể khóa tài khoản");
    } finally {
      setActionLoading(null);
    }
  };

  const handleUnlock = async (userId: string) => {
    try {
      setActionLoading(userId);
      await adminApi.unlockUser(userId);
      fetchUsers();
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể mở khóa tài khoản");
    } finally {
      setActionLoading(null);
    }
  };

  // ─── Render ───────────────────────────────────────────────────────

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <UserCog className="h-7 w-7" />
        <h1 className="text-2xl font-bold tracking-tight">Quản lý Người dùng</h1>
      </div>

      {/* Error */}
      {error && (
        <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
      )}

      {/* Loading */}
      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {/* Empty */}
      {!loading && users.length === 0 && !error && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <UserCog className="h-12 w-12 text-muted-foreground/50" />
            <p className="mt-4 text-sm text-muted-foreground">Chưa có người dùng nào</p>
          </CardContent>
        </Card>
      )}

      {/* Desktop Table (md+) */}
      {!loading && users.length > 0 && (
        <div className="hidden md:block">
          <Card>
            <CardContent className="p-0">
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
                              <span
                                key={r}
                                className="inline-flex items-center gap-1 rounded-full bg-zinc-100 px-2.5 py-0.5 text-xs font-medium"
                              >
                                <Shield className="h-3 w-3" />
                                {roleLabels[r] || r}
                              </span>
                            ))}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${status.className}`}>
                            {status.label}
                          </span>
                        </td>
                        <td className="px-6 py-4 text-right">
                          {user.status === "active" ? (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleLock(user.user_id)}
                              disabled={actionLoading === user.user_id}
                            >
                              {actionLoading === user.user_id ? (
                                <Loader2 className="h-4 w-4 animate-spin" />
                              ) : (
                                <Lock className="mr-1 h-4 w-4" />
                              )}
                              Khóa
                            </Button>
                          ) : user.status === "locked" ? (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleUnlock(user.user_id)}
                              disabled={actionLoading === user.user_id}
                            >
                              {actionLoading === user.user_id ? (
                                <Loader2 className="h-4 w-4 animate-spin" />
                              ) : (
                                <Unlock className="mr-1 h-4 w-4" />
                              )}
                              Mở khóa
                            </Button>
                          ) : null}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Mobile Cards (<md) */}
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
                          <span
                            key={r}
                            className="inline-flex items-center gap-1 rounded-full bg-zinc-100 px-2 py-0.5 text-xs"
                          >
                            <Shield className="h-3 w-3" />
                            {roleLabels[r] || r}
                          </span>
                        ))}
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${status.className}`}>
                          {status.label}
                        </span>
                      </div>
                    </div>
                    <div>
                      {user.status === "active" ? (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleLock(user.user_id)}
                          disabled={actionLoading === user.user_id}
                        >
                          {actionLoading === user.user_id ? (
                            <Loader2 className="h-4 w-4 animate-spin" />
                          ) : (
                            <Lock className="h-4 w-4" />
                          )}
                        </Button>
                      ) : user.status === "locked" ? (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleUnlock(user.user_id)}
                          disabled={actionLoading === user.user_id}
                        >
                          {actionLoading === user.user_id ? (
                            <Loader2 className="h-4 w-4 animate-spin" />
                          ) : (
                            <Unlock className="h-4 w-4" />
                          )}
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
