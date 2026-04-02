import React from "react";
import { Loader2, Lock, Shield, Unlock } from "lucide-react";
import { UserInfo, UserRole, UserStatus } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";

interface UsersDesktopTableProps {
  users: UserInfo[];
  actionLoading: string | null;
  currentUserId?: string;
  roleLabels: Record<UserRole, string>;
  statusLabels: Record<UserStatus, string>;
  statusVariants: Record<UserStatus, "default" | "secondary" | "destructive" | "outline">;
  onRequestLock: (userId: string) => void;
  onRequestUnlock: (userId: string) => void;
}

export function UsersDesktopTable({
  users,
  actionLoading,
  currentUserId,
  roleLabels,
  statusLabels,
  statusVariants,
  onRequestLock,
  onRequestUnlock,
}: UsersDesktopTableProps) {
  return (
    <div className="hidden md:block">
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Email</TableHead>
                <TableHead>Vai trò</TableHead>
                <TableHead>Trạng thái</TableHead>
                <TableHead className="text-right">Hành động</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.user_id}>
                  <TableCell className="font-medium">{user.email}</TableCell>
                  <TableCell>
                    <div className="flex flex-wrap gap-1">
                      {user.roles?.map((role) => (
                        <Badge key={role} variant="secondary">
                          <Shield className="h-3 w-3" /> {roleLabels[role] || role}
                        </Badge>
                      ))}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={statusVariants[user.status] || "secondary"}>
                      {statusLabels[user.status] || user.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    {user.status === "active" ? (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => onRequestLock(user.user_id)}
                        disabled={actionLoading === user.user_id || currentUserId === user.user_id}
                        title={currentUserId === user.user_id ? "Bạn không thể tự khóa chính mình" : ""}
                      >
                        {actionLoading === user.user_id ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Lock className="mr-1 h-4 w-4 text-destructive" />
                        )}
                        <span className="text-destructive">Khóa</span>
                      </Button>
                    ) : user.status === "locked" ? (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => onRequestUnlock(user.user_id)}
                        disabled={actionLoading === user.user_id}
                      >
                        {actionLoading === user.user_id ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Unlock className="mr-1 h-4 w-4 text-success" />
                        )}
                        <span className="text-success">Mở khóa</span>
                      </Button>
                    ) : null}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
