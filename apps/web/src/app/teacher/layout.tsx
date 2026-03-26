/**
 * Teacher Layout
 * Layout chung cho tất cả các trang của giáo viên.
 * Sử dụng ProtectedRoute + AppShell với teacherMenuItems.
 */
"use client";

import React from "react";
import { ProtectedRoute } from "@/components/layout/ProtectedRoute";
import { AppShell } from "@/components/layout/AppShell";
import { teacherMenuItems } from "@/components/layout/Sidebar";
import { TEACHER_ALLOWED_ROLES } from "@/lib/auth-config";

export default function TeacherLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={TEACHER_ALLOWED_ROLES}>
      <AppShell menuItems={teacherMenuItems}>
        {children}
      </AppShell>
    </ProtectedRoute>
  );
}