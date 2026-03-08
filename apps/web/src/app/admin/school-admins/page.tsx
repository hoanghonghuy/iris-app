/**
 * Admin School Admins Page
 * Quản lý School Admin: listing + tạo mới + xóa.
 * API: GET/POST/DELETE /admin/school-admins
 */
"use client";

import React, { useEffect, useState, useCallback } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { UserInfo, School } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  ShieldCheck, Loader2, Plus, X, Trash2, ChevronDown, Mail,
} from "lucide-react";

export default function AdminSchoolAdminsPage() {
  const [admins, setAdmins] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Create form
  const [showForm, setShowForm] = useState(false);
  const [schools, setSchools] = useState<School[]>([]);
  const [users, setUsers] = useState<UserInfo[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedUserId, setSelectedUserId] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchAdmins = useCallback(async () => {
    try {
      setLoading(true); setError("");
      const data = await adminApi.getSchoolAdmins();
      setAdmins(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể tải danh sách");
    } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchAdmins(); }, [fetchAdmins]);

  // Load schools + users when form opens
  useEffect(() => {
    if (!showForm) return;
    const load = async () => {
      try {
        const [schoolData, userData] = await Promise.all([
          adminApi.getSchools(),
          adminApi.getUsers({ limit: 100 }),
        ]);
        setSchools(schoolData || []);
        const userList = (userData as any).data || [];
        setUsers(Array.isArray(userList) ? userList : []);
        if (schoolData && schoolData.length > 0) setSelectedSchoolId(schoolData[0].school_id);
        if (userList && userList.length > 0) setSelectedUserId(userList[0].user_id);
      } catch { /* ignore */ }
    };
    load();
  }, [showForm]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedUserId || !selectedSchoolId) { setFormError("Chọn user và trường"); return; }
    try {
      setSubmitting(true); setFormError(""); setSuccess("");
      await adminApi.createSchoolAdmin({ user_id: selectedUserId, school_id: selectedSchoolId });
      setSuccess("Đã tạo School Admin thành công!");
      setShowForm(false);
      fetchAdmins();
    } catch (err: any) {
      setFormError(err.response?.data?.error || "Không thể tạo");
    } finally { setSubmitting(false); }
  };

  const handleDelete = async (adminId: string) => {
    try {
      setDeletingId(adminId);
      await adminApi.deleteSchoolAdmin(adminId);
      fetchAdmins();
    } catch (err: any) {
      setError(err.response?.data?.error || "Không thể xóa");
    } finally { setDeletingId(null); }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <ShieldCheck className="h-7 w-7" />
          <h1 className="text-2xl font-bold tracking-tight">Quản lý School Admin</h1>
        </div>
        <Button size="sm" onClick={() => { setShowForm(!showForm); setSuccess(""); }}>
          {showForm ? <X className="mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
          {showForm ? "Hủy" : "Thêm School Admin"}
        </Button>
      </div>

      {success && <div className="rounded-md bg-green-100 p-4 text-sm text-green-700">{success}</div>}
      {error && <div className="rounded-md bg-destructive/10 p-4 text-sm text-destructive">{error}</div>}

      {showForm && (
        <Card>
          <CardHeader><CardTitle className="text-lg">Gán School Admin</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              {formError && <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{formError}</div>}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium">User</label>
                  <div className="relative">
                    <select value={selectedUserId} onChange={(e) => setSelectedUserId(e.target.value)}
                      className="h-9 w-full appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                      {users.map((u) => <option key={u.user_id} value={u.user_id}>{u.email}</option>)}
                    </select>
                    <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  </div>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Trường</label>
                  <div className="relative">
                    <select value={selectedSchoolId} onChange={(e) => setSelectedSchoolId(e.target.value)}
                      className="h-9 w-full appearance-none rounded-md border bg-white py-1 pl-3 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                      {schools.map((s) => <option key={s.school_id} value={s.school_id}>{s.name}</option>)}
                    </select>
                    <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  </div>
                </div>
              </div>
              <div className="flex justify-end">
                <Button type="submit" disabled={submitting}>
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

      {/* Desktop Table */}
      {!loading && admins.length > 0 && (
        <div className="hidden md:block">
          <Card><CardContent className="p-0">
            <table className="w-full">
              <thead>
                <tr className="border-b text-left text-sm text-muted-foreground">
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Trường</th>
                  <th className="px-6 py-3 font-medium text-right">Hành động</th>
                </tr>
              </thead>
              <tbody>
                {admins.map((a) => (
                  <tr key={a.school_admin_id || a.user_id} className="border-b last:border-0 hover:bg-zinc-50">
                    <td className="px-6 py-4 font-medium">{a.email || a.user_id}</td>
                    <td className="px-6 py-4 text-muted-foreground">{a.school_name || a.school_id}</td>
                    <td className="px-6 py-4 text-right">
                      <Button variant="ghost" size="sm" onClick={() => handleDelete(a.school_admin_id || a.user_id)} disabled={deletingId === (a.school_admin_id || a.user_id)}>
                        {deletingId === (a.school_admin_id || a.user_id) ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="mr-1 h-4 w-4 text-destructive" />}
                        Xóa
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
      {!loading && admins.length > 0 && (
        <div className="space-y-3 md:hidden">
          {admins.map((a) => (
            <Card key={a.school_admin_id || a.user_id}>
              <CardContent className="py-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="flex items-center gap-2 font-medium"><Mail className="h-4 w-4 text-muted-foreground" /> {a.email || a.user_id}</p>
                    <p className="mt-1 text-sm text-muted-foreground">{a.school_name || a.school_id}</p>
                  </div>
                  <Button variant="ghost" size="sm" onClick={() => handleDelete(a.school_admin_id || a.user_id)}>
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
