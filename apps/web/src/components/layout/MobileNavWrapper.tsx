/**
 * MobileNavWrapper Component
 * Wrapper client-only để cho phép dùng dynamic({ ssr: false }) từ Server Component.
 * next/dynamic với ssr:false chỉ hoạt động bên trong Client Component.
 */
"use client";

import dynamic from "next/dynamic";

// Sheet của Radix UI sinh aria-controls ID ngẫu nhiên → phải bỏ SSR hoàn toàn
const MobileNav = dynamic(
  () => import("@/components/layout/MobileNav").then((m) => m.MobileNav),
  { ssr: false }
);

export function MobileNavWrapper() {
  return <MobileNav />;
}
