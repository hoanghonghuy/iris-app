/**
 * Admin Layout
 * Layout chung cho tất cả các trang của Admin (SUPER_ADMIN & SCHOOL_ADMIN).
 * Sử dụng ProtectedRoute + AppShell với adminMenuItems.
 */
"use client";

import React from "react";
import { ProtectedRoute } from "@/components/layout/ProtectedRoute";
import { AppShell } from "@/components/layout/AppShell";
import { adminMenuItems } from "@/components/layout/Sidebar";

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={["SUPER_ADMIN", "SCHOOL_ADMIN"]}>
      <AppShell menuItems={adminMenuItems}>
        {children}
      </AppShell>
    </ProtectedRoute>
  );
}