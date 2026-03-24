/**
 * AppShell Component
 * Layout wrapper chung cho các trang sau khi đăng nhập.
 * Kết hợp Sidebar + Header + Main content với responsive behavior.
 * Dùng trong layout.tsx của từng role.
 */
"use client";

import React, { useState } from "react";
import { Sidebar, SidebarItem } from "@/components/layout/Sidebar";
import { Header } from "@/components/layout/Header";

interface AppShellProps {
  children: React.ReactNode;
  menuItems: SidebarItem[];
}

export function AppShell({ children, menuItems }: AppShellProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="flex h-screen overflow-hidden bg-background transition-colors duration-300">
      {/* Sidebar */}
      <Sidebar
        items={menuItems}
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
      />

      {/* Main area (Header + Content) */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <Header onMenuToggle={() => setSidebarOpen(!sidebarOpen)} />
        <main className="flex-1 overflow-y-auto p-4 md:p-6 lg:p-8">
          {children}
        </main>
      </div>
    </div>
  );
}
