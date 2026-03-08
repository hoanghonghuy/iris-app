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

export default function TeacherLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={["TEACHER"]}>
      <AppShell menuItems={teacherMenuItems}>
        {children}
      </AppShell>
    </ProtectedRoute>
  );
}