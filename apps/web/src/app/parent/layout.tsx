/**
 * Parent Layout
 * Layout chung cho tất cả các trang của phụ huynh.
 * Sử dụng ProtectedRoute + AppShell với parentMenuItems.
 */
"use client";

import React from "react";
import { ProtectedRoute } from "@/components/layout/ProtectedRoute";
import { AppShell } from "@/components/layout/AppShell";
import { parentMenuItems } from "@/components/layout/Sidebar";

export default function ParentLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={["PARENT"]}>
      <AppShell menuItems={parentMenuItems}>
        {children}
      </AppShell>
    </ProtectedRoute>
  );
}