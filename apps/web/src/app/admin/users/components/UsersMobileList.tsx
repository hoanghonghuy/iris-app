import React from "react";
import { Loader2, Lock, Mail, Shield, Unlock } from "lucide-react";
import { UserInfo, UserRole, UserStatus } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

interface UsersMobileListProps {
  users: UserInfo[];
  actionLoading: string | null;
  currentUserId?: string;
  roleLabels: Record<UserRole, string>;
  statusLabels: Record<UserStatus, string>;
  statusVariants: Record<UserStatus, "default" | "secondary" | "destructive" | "outline">;
  onRequestLock: (userId: string) => void;
  onRequestUnlock: (userId: string) => void;
}

export function UsersMobileList({
  users,
  actionLoading,
  currentUserId,
  roleLabels,
  statusLabels,
  statusVariants,
  onRequestLock,
  onRequestUnlock,
}: UsersMobileListProps) {
  return (
    <div className="space-y-3 md:hidden">
      {users.map((user) => (
        <Card key={user.user_id}>
          <CardContent className="py-4">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0 flex-1">
                <p className="flex items-center gap-2 font-medium">
                  <Mail className="h-4 w-4 shrink-0 text-muted-foreground" />
                  <span className="truncate">{user.email}</span>
                </p>
                <div className="mt-2 flex flex-wrap gap-1">
                  {user.roles?.map((role) => (
                    <Badge key={role} variant="secondary">
                      <Shield className="h-3 w-3" /> {roleLabels[role] || role}
                    </Badge>
                  ))}
                  <Badge variant={statusVariants[user.status] || "secondary"}>
                    {statusLabels[user.status] || user.status}
                  </Badge>
                </div>
              </div>

              <div>
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
                      <Lock className="h-4 w-4 text-destructive" />
                    )}
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
                      <Unlock className="h-4 w-4 text-success" />
                    )}
                  </Button>
                ) : null}
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
