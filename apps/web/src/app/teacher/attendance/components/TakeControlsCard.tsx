import React from "react";
import { Loader2 } from "lucide-react";
import { AttendanceStatus } from "@/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { TakeListFilter } from "../config";

interface TakeControlsCardProps {
  showMobileTakeControls: boolean;
  studentSearch: string;
  takeListFilter: TakeListFilter;
  listOrderMode: "prioritize" | "original";
  displayedStudentsLength: number;
  studentsLength: number;
  displayedDirtyCount: number;
  displayedSavedCount: number;
  globalPendingCount: number;
  savingDisplayed: boolean;
  onToggleMobileControls: () => void;
  onStudentSearchChange: (value: string) => void;
  onTakeListFilterChange: (value: TakeListFilter) => void;
  onListOrderModeChange: (value: "prioritize" | "original") => void;
  onApplyStatusToDisplayed: (status: AttendanceStatus) => void;
  onSaveDisplayed: () => void;
}

export function TakeControlsCard({
  showMobileTakeControls,
  studentSearch,
  takeListFilter,
  listOrderMode,
  displayedStudentsLength,
  studentsLength,
  displayedDirtyCount,
  displayedSavedCount,
  globalPendingCount,
  savingDisplayed,
  onToggleMobileControls,
  onStudentSearchChange,
  onTakeListFilterChange,
  onListOrderModeChange,
  onApplyStatusToDisplayed,
  onSaveDisplayed,
}: TakeControlsCardProps) {
  return (
    <Card>
      <CardContent className="py-3">
        <div className="flex items-center justify-between gap-2 sm:hidden">
          <p className="text-xs text-muted-foreground">
            Hiển thị {displayedStudentsLength}/{studentsLength} • Chờ lưu {displayedDirtyCount}
          </p>
          <Button
            size="sm"
            variant="outline"
            className="h-8 px-2.5 text-xs"
            onClick={onToggleMobileControls}
            aria-expanded={showMobileTakeControls}
          >
            {showMobileTakeControls ? "Ẩn bộ lọc" : "Mở bộ lọc"}
          </Button>
        </div>

        <div className={`${showMobileTakeControls ? "mt-2" : "hidden"} sm:mt-0 sm:block`}>
          <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
            <Input
              value={studentSearch}
              onChange={(e) => onStudentSearchChange(e.target.value)}
              placeholder="Tìm học sinh theo tên..."
              className="h-9 text-sm"
              aria-label="Tìm học sinh theo tên"
            />

            <Select value={takeListFilter} onValueChange={onTakeListFilterChange}>
              <SelectTrigger className="h-9 text-sm" aria-label="Lọc theo trạng thái lưu">
                <SelectValue placeholder="Lọc danh sách" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Tất cả học sinh</SelectItem>
                <SelectItem value="pending">Chưa lưu / đang sửa</SelectItem>
                <SelectItem value="saved">Đã lưu</SelectItem>
              </SelectContent>
            </Select>

            <Button
              size="sm"
              variant={listOrderMode === "prioritize" ? "default" : "outline"}
              onClick={() => onListOrderModeChange("prioritize")}
              className="h-9"
            >
              Ưu tiên chưa lưu
            </Button>
            <Button
              size="sm"
              variant={listOrderMode === "original" ? "default" : "outline"}
              onClick={() => onListOrderModeChange("original")}
              className="h-9"
            >
              Giữ nguyên thứ tự
            </Button>
          </div>

          <div className="mt-2 flex flex-wrap items-center gap-1.5">
            <Badge variant="outline" className="text-xs">
              Toàn lớp chờ lưu: {globalPendingCount}
            </Badge>
            <Badge variant="outline" className="text-xs">
              Đang hiển thị: {displayedStudentsLength}/{studentsLength}
            </Badge>
            <Badge variant={displayedDirtyCount > 0 ? "secondary" : "outline"} className="text-xs">
              Chờ lưu trong danh sách: {displayedDirtyCount}
            </Badge>
            <Badge variant={displayedSavedCount > 0 ? "default" : "outline"} className="text-xs">
              Đã lưu trong danh sách: {displayedSavedCount}
            </Badge>
          </div>

          <div className="mt-2 flex flex-wrap items-center gap-2">
            <Button size="sm" variant="outline" onClick={() => onApplyStatusToDisplayed("present")}>Đặt tất cả hiển thị: Có mặt</Button>
            <Button size="sm" variant="outline" onClick={() => onApplyStatusToDisplayed("absent")}>Đặt tất cả hiển thị: Vắng</Button>
            <Button size="sm" variant="outline" onClick={() => onApplyStatusToDisplayed("late")}>Đặt tất cả hiển thị: Muộn</Button>
            <Button size="sm" variant="outline" onClick={() => onApplyStatusToDisplayed("excused")}>Đặt tất cả hiển thị: Có phép</Button>
            <Button size="sm" onClick={onSaveDisplayed} disabled={savingDisplayed || displayedDirtyCount === 0}>
              {savingDisplayed ? <Loader2 className="h-4 w-4 animate-spin" /> : `Lưu danh sách hiển thị${displayedDirtyCount > 0 ? ` (${displayedDirtyCount})` : ""}`}
            </Button>
          </div>

          <p className="mt-1 text-xs text-muted-foreground">
            {listOrderMode === "prioritize" ? "Đang ưu tiên chưa lưu" : "Đang giữ nguyên thứ tự"} • Dùng “Lưu danh sách hiển thị” để chốt nhanh phần đang lọc.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
