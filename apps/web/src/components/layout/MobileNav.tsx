/**
 * MobileNav Component
 * Menu điều hướng dạng Sheet trên mobile cho Landing Page.
 * Bắt buộc là Client Component vì Radix Sheet dùng state nội bộ.
 */
"use client";

import Link from "next/link";
import { GraduationCap, Menu } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger, SheetTitle } from "@/components/ui/sheet";

export function MobileNav() {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          className="h-9 w-9 bg-transparent border-zinc-200 dark:border-zinc-800 text-zinc-900 dark:text-zinc-100"
        >
          <Menu className="h-5 w-5" />
          <span className="sr-only">Toggle navigation menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent
        side="right"
        className="w-[300px] sm:w-[350px] bg-white dark:bg-zinc-950 border-l border-zinc-200 dark:border-zinc-800 flex flex-col p-0"
      >
        <SheetTitle className="sr-only">Menu Điều Hướng</SheetTitle>
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="p-6 pr-12 border-b border-zinc-100 dark:border-zinc-900 flex items-center gap-3">
            <GraduationCap className="h-6 w-6 text-zinc-900 dark:text-zinc-100" />
            <span className="text-xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">
              Iris School
            </span>
          </div>

          <nav className="flex-1 flex flex-col py-8 px-6">
            <div className="flex flex-col gap-2">
              <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-400 dark:text-zinc-500 mb-2 px-2">
                Điều hướng
              </p>
              <a
                href="#features"
                className="flex items-center gap-3 px-3 py-3 rounded-xl text-lg font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 hover:bg-zinc-50 dark:hover:bg-zinc-900 transition-all active:scale-95"
              >
                Tính năng
              </a>
              <a
                href="#about"
                className="flex items-center gap-3 px-3 py-3 rounded-xl text-lg font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 hover:bg-zinc-50 dark:hover:bg-zinc-900 transition-all active:scale-95"
              >
                Về chúng tôi
              </a>
            </div>

            <div className="mt-auto flex flex-col gap-4 pb-8">
              <div className="h-px bg-zinc-100 dark:bg-zinc-900 w-full mb-4" />
              <Link href="/login" className="w-full">
                <Button
                  variant="outline"
                  className="w-full justify-center h-12 text-base font-semibold border-zinc-200 dark:border-zinc-800 rounded-xl"
                >
                  Đăng nhập
                </Button>
              </Link>
              <Link href="/register" className="w-full">
                <Button className="w-full justify-center h-12 text-base font-semibold bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-950 hover:bg-zinc-800 dark:hover:bg-zinc-200 rounded-xl shadow-lg shadow-zinc-200/50 dark:shadow-none">
                  Đăng ký Phụ huynh
                </Button>
              </Link>
            </div>
          </nav>
        </div>
      </SheetContent>
    </Sheet>
  );
}
