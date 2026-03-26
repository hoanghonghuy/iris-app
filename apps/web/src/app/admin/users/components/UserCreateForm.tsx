import React from "react";
import { AlertCircle, Loader2 } from "lucide-react";
import { UserRole } from "@/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface UserCreateFormProps {
  formError: string;
  formEmail: string;
  formRoles: string[];
  submitting: boolean;
  creatableRoles: UserRole[];
  roleLabels: Record<UserRole, string>;
  onEmailChange: (value: string) => void;
  onToggleRole: (role: string) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export function UserCreateForm({
  formError,
  formEmail,
  formRoles,
  submitting,
  creatableRoles,
  roleLabels,
  onEmailChange,
  onToggleRole,
  onSubmit,
}: UserCreateFormProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Tạo tài khoản mới</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="space-y-4">
          {formError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{formError}</AlertDescription>
            </Alert>
          )}

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="userEmail">
                Email <span className="text-destructive">*</span>
              </Label>
              <Input
                id="userEmail"
                type="email"
                placeholder="user@example.com"
                value={formEmail}
                onChange={(e) => onEmailChange(e.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label>
                Vai trò <span className="text-destructive">*</span>
              </Label>
              <div className="flex flex-wrap gap-2">
                {creatableRoles.map((role) => (
                  <Badge
                    key={role}
                    variant={formRoles.includes(role) ? "default" : "outline"}
                    className="cursor-pointer select-none"
                    onClick={() => onToggleRole(role)}
                  >
                    {roleLabels[role]}
                  </Badge>
                ))}
              </div>
            </div>
          </div>

          <p className="text-xs text-muted-foreground">
            User sẽ ở trạng thái &ldquo;Chờ kích hoạt&rdquo;. Họ cần dùng activation token để đặt mật khẩu.
          </p>

          <div className="flex justify-end">
            <Button type="submit" disabled={submitting}>
              {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Tạo user
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
