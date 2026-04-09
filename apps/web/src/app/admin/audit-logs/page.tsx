"use client";

import { useCallback, useEffect, useState } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { AuditLog } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Loader2 } from "lucide-react";
import { useAuth } from "@/providers/AuthProvider";

export default function AdminAuditLogsPage() {
  const { role } = useAuth();
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [q, setQ] = useState("");
  const [entityType, setEntityType] = useState("");
  const [action, setAction] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");

  const load = useCallback(async () => {
    if (role !== "SUPER_ADMIN") {
      setLogs([]);
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
        limit: 100,
        offset: 0,
      });
      setLogs(res.data || []);
    } catch {
      setLogs([]);
      setError("Khong the tai du lieu Audit Logs.");
    } finally {
      setLoading(false);
    }
  }, [action, entityType, from, q, role, to]);

  useEffect(() => {
    void load();
  }, [load]);

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
              <input className="border rounded px-3 py-2" placeholder="Tim kiem chi tiet" value={q} onChange={(e) => setQ(e.target.value)} />
              <input className="border rounded px-3 py-2" placeholder="Entity type" value={entityType} onChange={(e) => setEntityType(e.target.value)} />
              <input className="border rounded px-3 py-2" placeholder="Action" value={action} onChange={(e) => setAction(e.target.value)} />
              <input className="border rounded px-3 py-2" type="datetime-local" value={from} onChange={(e) => setFrom(e.target.value)} />
              <input className="border rounded px-3 py-2" type="datetime-local" value={to} onChange={(e) => setTo(e.target.value)} />
              <button className="px-4 py-2 rounded bg-primary text-primary-foreground md:col-span-5 w-fit" onClick={() => void load()}>Loc logs</button>
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
        </>
      )}
    </div>
  );
}
