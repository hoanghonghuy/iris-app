/**
 * PaginationBar Component
 * Shared pagination controls for list pages.
 * Displays: « Prev | Page X / Y | Next » + total count.
 */
"use client";

import React from "react";
import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { Pagination } from "@/types";
import { cn } from "@/lib/utils";

const DEFAULT_PAGE_SIZE = 20;

interface PaginationBarProps {
  pagination: Pagination;
  onPageChange: (offset: number) => void;
  className?: string;
}

export function PaginationBar({ pagination, onPageChange, className }: PaginationBarProps) {
  const { total, limit, offset, has_more } = pagination;

  const pageSize = limit || DEFAULT_PAGE_SIZE;
  const currentPage = Math.floor(offset / pageSize) + 1;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const canGoPrev = offset > 0;
  const canGoNext = has_more;

  if (total <= pageSize) return null; // No pagination needed

  return (
    <div className={cn("flex flex-col sm:flex-row items-center justify-between gap-3 pt-4", className)}>
      <p className="text-sm text-muted-foreground">
        Tổng <span className="font-medium text-foreground">{total}</span> bản ghi
      </p>
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          disabled={!canGoPrev}
          onClick={() => onPageChange(Math.max(0, offset - pageSize))}
        >
          <ChevronLeft className="h-4 w-4 mr-1" />
          Trước
        </Button>
        <span className="text-sm text-muted-foreground min-w-[80px] text-center">
          Trang {currentPage} / {totalPages}
        </span>
        <Button
          variant="outline"
          size="sm"
          disabled={!canGoNext}
          onClick={() => onPageChange(offset + pageSize)}
        >
          Sau
          <ChevronRight className="h-4 w-4 ml-1" />
        </Button>
      </div>
    </div>
  );
}
