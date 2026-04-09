"use client";

import { useCallback, useEffect, useState } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { AuditLog, Pagination } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Loader2 } from "lucide-react";
import { useAuth } from "@/providers/AuthProvider";

const DEFAULT_LIMIT = 20;

export default function AdminAuditLogsPage() {
  const { role } = useAuth();
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [pagination, setPagination] = useState<Pagination>({
    total: 0,
    limit: DEFAULT_LIMIT,
    offset: 0,
    has_more: false,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [q, setQ] = useState("");
  const [entityType, setEntityType] = useState("");
  const [action, setAction] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");

  const [pendingQ, setPendingQ] = useState("");
  const [pendingEntityType, setPendingEntityType] = useState("");
  const [pendingAction, setPendingAction] = useState("");
  const [pendingFrom, setPendingFrom] = useState("");
  const [pendingTo, setPendingTo] = useState("");

  const load = useCallback(async () => {
    if (role !== "SUPER_ADMIN") {
      setLogs([]);
      setPagination({ total: 0, limit: DEFAULT_LIMIT, offset: 0, has_more: false });
      setError("Ban khong co quyen truy cap Audit Logs.");
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const res = await adminApi.getAuditLogs({
        q: q || undefined,
        entity_type: entityType || undefined,
        action: action || undefined,
        from: from ? new Date(from).toISOString() : undefined,
        to: to ? new Date(to).toISOString() : undefined,
        limit: pagination.limit,
        offset: pagination.offset,
      });
      setLogs(res.data || []);
      setPagination({
        total: res.pagination?.total ?? 0,
        limit: res.pagination?.limit ?? pagination.limit,
        offset: res.pagination?.offset ?? pagination.offset,
        has_more: res.pagination?.has_more ?? false,
      });
    } catch {
      setLogs([]);
      setPagination((prev) => ({ ...prev, total: 0, has_more: false }));
      setError("Khong the tai du lieu Audit Logs.");
    } finally {
      setLoading(false);
    }
  }, [action, entityType, from, pagination.limit, pagination.offset, q, role, to]);

  const applyFilters = () => {
    setQ(pendingQ);
    setEntityType(pendingEntityType);
    setAction(pendingAction);
    setFrom(pendingFrom);
    setTo(pendingTo);
    setPagination((prev) => ({ ...prev, offset: 0 }));
  };

  const previousPage = () => {
    setPagination((prev) => ({
      ...prev,
      offset: Math.max(prev.offset - prev.limit, 0),
    }));
  };

  const nextPage = () => {
    setPagination((prev) => ({
      ...prev,
      offset: prev.offset + prev.limit,
    }));
  };

  const updatePageSize = (nextLimit: number) => {
    setPagination((prev) => ({
      ...prev,
      limit: nextLimit,
      offset: 0,
    }));
  };

  useEffect(() => {
    void load();
  }, [load]);

  const fromRow = logs.length === 0 ? 0 : pagination.offset + 1;
  const toRow = pagination.offset + logs.length;
  const currentPage = Math.floor(pagination.offset / Math.max(pagination.limit, 1)) + 1;
  const totalPages = pagination.total > 0 ? Math.ceil(pagination.total / Math.max(pagination.limit, 1)) : 1;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Audit Logs</h1>

      {role !== "SUPER_ADMIN" && (
        <Card>
          <CardContent className="p-4 text-sm text-muted-foreground">
            Ban khong co quyen truy cap trang nay.
          </CardContent>
        </Card>
      )}

      {role === "SUPER_ADMIN" && (
        <>
          {error && (
            <Card>
              <CardContent className="p-4 text-sm text-destructive">{error}</CardContent>
            </Card>
          )}

          <Card>
            <CardContent className="p-4 grid gap-3 md:grid-cols-5">
              <input className="border rounded px-3 py-2" placeholder="Tim kiem chi tiet" value={pendingQ} onChange={(e) => setPendingQ(e.target.value)} />
              <input className="border rounded px-3 py-2" placeholder="Entity type" value={pendingEntityType} onChange={(e) => setPendingEntityType(e.target.value)} />
              <input className="border rounded px-3 py-2" placeholder="Action" value={pendingAction} onChange={(e) => setPendingAction(e.target.value)} />
              <input className="border rounded px-3 py-2" type="datetime-local" value={pendingFrom} onChange={(e) => setPendingFrom(e.target.value)} />
              <input className="border rounded px-3 py-2" type="datetime-local" value={pendingTo} onChange={(e) => setPendingTo(e.target.value)} />
              <div className="md:col-span-5 flex flex-wrap items-center gap-3">
                <button className="px-4 py-2 rounded bg-primary text-primary-foreground" onClick={applyFilters}>Loc logs</button>
                <label className="text-sm text-muted-foreground flex items-center gap-2">
                  Page size
                  <select
                    className="border rounded px-2 py-1 text-sm"
                    value={pagination.limit}
                    onChange={(e) => updatePageSize(Number(e.target.value))}
                  >
                    <option value={20}>20</option>
                    <option value={50}>50</option>
                    <option value={100}>100</option>
                  </select>
                </label>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-0 overflow-auto">
              {loading ? (
                <div className="py-8 flex justify-center"><Loader2 className="h-6 w-6 animate-spin" /></div>
              ) : logs.length === 0 ? (
                <div className="p-4 text-sm text-muted-foreground">Khong co log.</div>
              ) : (
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b bg-muted/40 text-left">
                      <th className="p-3">Thoi gian</th>
                      <th className="p-3">Actor</th>
                      <th className="p-3">Role</th>
                      <th className="p-3">Action</th>
                      <th className="p-3">Entity</th>
                      <th className="p-3">Details</th>
                    </tr>
                  </thead>
                  <tbody>
                    {logs.map((log) => (
                      <tr key={log.audit_log_id} className="border-b align-top">
                        <td className="p-3 whitespace-nowrap">{new Date(log.created_at).toLocaleString("vi-VN")}</td>
                        <td className="p-3">{log.actor_user_id}</td>
                        <td className="p-3">{log.actor_role || "-"}</td>
                        <td className="p-3">{log.action}</td>
                        <td className="p-3">{log.entity_type}{log.entity_id ? ` / ${log.entity_id}` : ""}</td>
                        <td className="p-3 max-w-[420px]">
                          <pre className="whitespace-pre-wrap break-words text-xs">{JSON.stringify(log.details ?? {}, null, 2)}</pre>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </CardContent>
          </Card>

          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-3 text-sm">
            <p className="text-muted-foreground">
              Hiển thị {fromRow}-{toRow} trên tổng {pagination.total} bản ghi.
            </p>
            <div className="flex items-center gap-2">
              <span className="text-muted-foreground">Trang {currentPage}/{totalPages}</span>
              <button
                className="px-3 py-1.5 border rounded disabled:opacity-50"
                onClick={previousPage}
                disabled={loading || pagination.offset === 0}
              >
                Trước
              </button>
              <button
                className="px-3 py-1.5 border rounded disabled:opacity-50"
                onClick={nextPage}
                disabled={loading || !pagination.has_more}
              >
                Sau
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
