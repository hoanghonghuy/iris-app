/**
 * Parent Layout
 * Layout chung cho tất cả các trang của phụ huynh.
 * Sử dụng ProtectedRoute để đảm bảo chỉ phụ huynh mới truy cập được.
 */
import React from 'react';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';

export default function ParentLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={['PARENT']}>
      <div className="min-h-screen bg-zinc-50">
        {/* Sidebar & Header sẽ được thêm vào đây sau */}
        <main className="p-4 md:p-8">
          {children}
        </main>
      </div>
    </ProtectedRoute>
  );
}