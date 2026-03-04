/**
 * Admin Layout
 * Layout chung cho tất cả các trang của Admin (SUPER_ADMIN & SCHOOL_ADMIN).
 * Sử dụng ProtectedRoute để đảm bảo chỉ Admin mới truy cập được.
 */
import React from 'react';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={['SUPER_ADMIN', 'SCHOOL_ADMIN']}>
      <div className="min-h-screen bg-zinc-50">
        {/* Sidebar & Header sẽ được thêm vào đây sau */}
        <main className="p-4 md:p-8">
          {children}
        </main>
      </div>
    </ProtectedRoute>
  );
}