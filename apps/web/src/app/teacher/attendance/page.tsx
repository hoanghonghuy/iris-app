/**
 * Teacher Attendance Page
 * Chọn lớp → xem HS → điểm danh từng em.
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/attendance
 * API: GET /teacher/classes, GET /teacher/classes/:id/students, POST /teacher/attendance
 */
"use client";

import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ClipboardCheck, Loader2, AlertCircle } from "lucide-react";
import { TakeControlsCard } from "./components/TakeControlsCard";
import { StudentAttendanceCard } from "./components/StudentAttendanceCard";
import { AttendanceHistoryView } from "./components/AttendanceHistoryView";
import { useTeacherAttendancePage } from "./useTeacherAttendancePage";

export default function TeacherAttendancePage() {
  const {
    classes,
    selectedClassId,
    students,
    loadingClasses,
    loadingStudents,
    error,
    submitting,
    canceling,
    savingAll,
    savingDisplayed,
    historyOpen,
    historyLoading,
    historyByStudent,
    studentSearch,
    listOrderMode,
    takeListFilter,
    showMobileTakeControls,
    viewMode,
    historyFrom,
    historyTo,
    historyStudentId,
    historyStatus,
    historyListLoading,
    historyList,
    historyOffset,
    historyLimit,
    historyTotal,
    historyHasMore,
    dirtyCount,
    displayedStudents,
    displayedDirtyCount,
    displayedSavedCount,
    globalPendingCount,
    attendance,
    hasSavedToday,
    setSelectedClassId,
    setStudentSearch,
    setListOrderMode,
    setTakeListFilter,
    setShowMobileTakeControls,
    setViewMode,
    setHistoryFrom,
    setHistoryTo,
    setHistoryStudentId,
    setHistoryStatus,
    isRowDirty,
    handleMark,
    handleRevertLocal,
    handleCancelSaved,
    handleSaveAll,
    handleSaveDisplayed,
    applyStatusToDisplayed,
    toggleHistory,
    handleHistorySearch,
    handleHistoryPrev,
    handleHistoryNext,
    handleAttendanceStatusChange,
    handleAttendanceNoteChange,
  } = useTeacherAttendancePage();

  if (loadingClasses) {
    return <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        <Button size="sm" variant={viewMode === "take" ? "default" : "outline"} onClick={() => setViewMode("take")}>Điểm danh hôm nay</Button>
        <Button size="sm" variant={viewMode === "history" ? "default" : "outline"} onClick={() => setViewMode("history")}>Lịch sử lớp</Button>
      </div>

      <div className="flex items-center gap-2">
        {classes.length > 0 && (
          <Select value={selectedClassId} onValueChange={setSelectedClassId}>
            <SelectTrigger className="w-full sm:w-[220px]"><SelectValue placeholder="Chọn lớp" /></SelectTrigger>
            <SelectContent>
              {classes.map((c) => <SelectItem key={c.class_id} value={c.class_id}>{c.name}</SelectItem>)}
            </SelectContent>
          </Select>
        )}
      </div>

      {error && <Alert variant="destructive"><AlertCircle className="h-4 w-4" /><AlertDescription>{error}</AlertDescription></Alert>}
      {loadingStudents && <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}

      {!loadingStudents && students.length === 0 && selectedClassId && (
        <Card><CardContent className="flex flex-col items-center justify-center py-12">
          <ClipboardCheck className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-sm text-muted-foreground">Không có học sinh</p>
        </CardContent></Card>
      )}

      {!loadingStudents && students.length > 0 && viewMode === "take" && (
        <div className="space-y-4">
          {/* Summary Ring (P1) */}
          <Card className="bg-card/60 backdrop-blur-sm border-transparent shadow-sm">
            <CardContent className="p-4 flex flex-col sm:flex-row items-center gap-6">
              {(() => {
                const total = students.length;
                const presentCount = students.filter(s => hasSavedToday[s.student_id] && attendance[s.student_id]?.status === "present").length;
                const absentCount = students.filter(s => hasSavedToday[s.student_id] && attendance[s.student_id]?.status === "absent").length;
                const lateCount = students.filter(s => hasSavedToday[s.student_id] && attendance[s.student_id]?.status === "late").length;
                const completedCount = presentCount + absentCount + lateCount;
                const pendingCount = total - completedCount;
                
                const presentPct = total ? (presentCount / total) * 100 : 0;
                const absentPct = total ? (absentCount / total) * 100 : 0;
                const latePct = total ? (lateCount / total) * 100 : 0;

                return (
                  <>
                    <div 
                      className="w-20 h-20 rounded-full flex items-center justify-center shrink-0 shadow-inner"
                      style={{
                        background: `conic-gradient(
                          #22c55e 0% ${presentPct}%, 
                          #ef4444 ${presentPct}% ${presentPct + absentPct}%, 
                          #f59e0b ${presentPct + absentPct}% ${presentPct + absentPct + latePct}%, 
                          hsl(var(--muted)) ${presentPct + absentPct + latePct}% 100%
                        )`
                      }}
                    >
                      <div className="w-[60px] h-[60px] bg-card/95 rounded-full flex flex-col items-center justify-center">
                        <span className="text-sm font-bold leading-none text-foreground">{completedCount}</span>
                        <span className="text-[10px] font-medium text-muted-foreground mt-0.5">/ {total}</span>
                      </div>
                    </div>
                    
                    <div className="flex-1 grid grid-cols-2 md:grid-cols-4 gap-4 w-full">
                      <div className="flex flex-col">
                        <span className="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground uppercase opacity-80">
                          <div className="w-2 h-2 rounded-full bg-success"></div> Có mặt
                        </span>
                        <span className="text-xl font-bold mt-1 text-foreground">{presentCount}</span>
                      </div>

                      <div className="flex flex-col">
                        <span className="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground uppercase opacity-80">
                          <div className="w-2 h-2 rounded-full bg-destructive"></div> Vắng
                        </span>
                        <span className="text-xl font-bold mt-1 text-foreground">{absentCount}</span>
                      </div>

                      <div className="flex flex-col">
                        <span className="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground uppercase opacity-80">
                          <div className="w-2 h-2 rounded-full bg-orange-500"></div> Muộn
                        </span>
                        <span className="text-xl font-bold mt-1 text-foreground">{lateCount}</span>
                      </div>

                      <div className="flex flex-col">
                        <span className="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground uppercase opacity-80">
                          <div className="w-2 h-2 rounded-full bg-muted-foreground/30"></div> Chưa lưu
                        </span>
                        <span className="text-xl font-bold mt-1 text-foreground">{pendingCount}</span>
                      </div>
                    </div>
                  </>
                );
              })()}
            </CardContent>
          </Card>

          <TakeControlsCard
            showMobileTakeControls={showMobileTakeControls}
            studentSearch={studentSearch}
            takeListFilter={takeListFilter}
            listOrderMode={listOrderMode}
            displayedStudentsLength={displayedStudents.length}
            studentsLength={students.length}
            displayedDirtyCount={displayedDirtyCount}
            displayedSavedCount={displayedSavedCount}
            globalPendingCount={globalPendingCount}
            savingDisplayed={savingDisplayed}
            onToggleMobileControls={() => setShowMobileTakeControls((prev) => !prev)}
            onStudentSearchChange={setStudentSearch}
            onTakeListFilterChange={setTakeListFilter}
            onListOrderModeChange={setListOrderMode}
            onApplyStatusToDisplayed={applyStatusToDisplayed}
            onSaveDisplayed={handleSaveDisplayed}
          />

          {displayedStudents.length === 0 && (
            <Card>
              <CardContent className="py-6 text-sm text-muted-foreground">
                Không có học sinh phù hợp với bộ lọc hiện tại. Hãy đổi từ khóa tìm kiếm hoặc chuyển bộ lọc danh sách.
              </CardContent>
            </Card>
          )}

          {displayedStudents.map((student) => {
            const attendanceValue = attendance[student.student_id] || { status: "present", note: "" };
            return (
              <StudentAttendanceCard
                key={student.student_id}
                student={student}
                attendanceValue={attendanceValue}
                isDirty={isRowDirty(student.student_id)}
                hasSavedToday={!!hasSavedToday[student.student_id]}
                isSaving={submitting === student.student_id}
                isCanceling={canceling === student.student_id}
                isHistoryOpen={historyOpen.has(student.student_id)}
                isHistoryLoading={historyLoading.has(student.student_id)}
                historyRecords={historyByStudent[student.student_id] || []}
                onStatusChange={handleAttendanceStatusChange}
                onNoteChange={handleAttendanceNoteChange}
                onSave={handleMark}
                onRevert={handleRevertLocal}
                onCancelSaved={handleCancelSaved}
                onToggleHistory={toggleHistory}
              />
            );
          })}

          {(displayedDirtyCount > 0 || globalPendingCount > 0) && (
            <div className="sticky bottom-3 z-20 rounded-lg border bg-background/95 p-3 shadow-sm backdrop-blur supports-[backdrop-filter]:bg-background/70">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-xs text-muted-foreground sm:text-sm">
                  Còn {displayedDirtyCount} học sinh chưa lưu trong danh sách hiển thị • Toàn lớp còn {globalPendingCount} học sinh chưa lưu.
                </p>

                <div className="flex flex-wrap items-center gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={handleSaveDisplayed}
                    disabled={savingDisplayed || displayedDirtyCount === 0}
                  >
                    {savingDisplayed ? <Loader2 className="h-4 w-4 animate-spin" /> : `Lưu danh sách hiển thị${displayedDirtyCount > 0 ? ` (${displayedDirtyCount})` : ""}`}
                  </Button>
                  <Button size="sm" onClick={handleSaveAll} disabled={savingAll || dirtyCount === 0}>
                    {savingAll ? <Loader2 className="h-4 w-4 animate-spin" /> : `Lưu toàn lớp${dirtyCount > 0 ? ` (${dirtyCount})` : ""}`}
                  </Button>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {!loadingStudents && students.length > 0 && viewMode === "history" && (
        <AttendanceHistoryView
          students={students}
          historyFrom={historyFrom}
          historyTo={historyTo}
          historyStudentId={historyStudentId}
          historyStatus={historyStatus}
          historyListLoading={historyListLoading}
          historyList={historyList}
          historyTotal={historyTotal}
          historyOffset={historyOffset}
          historyLimit={historyLimit}
          historyHasMore={historyHasMore}
          onHistoryFromChange={setHistoryFrom}
          onHistoryToChange={setHistoryTo}
          onHistoryStudentChange={setHistoryStudentId}
          onHistoryStatusChange={setHistoryStatus}
          onHistorySearch={handleHistorySearch}
          onHistoryPrev={handleHistoryPrev}
          onHistoryNext={handleHistoryNext}
        />
      )}
    </div>
  );
}
