import { useMemo, useState } from "react";
import { Student } from "@/types";
import { TakeListFilter } from "./config";

export interface UseAttendanceTakeModeParams {
  students: Student[];
  isRowDirty: (studentId: string) => boolean;
}

export interface UseAttendanceTakeModeResult {
  studentSearch: string;
  listOrderMode: "prioritize" | "original";
  takeListFilter: TakeListFilter;
  showMobileTakeControls: boolean;
  displayedStudents: Student[];
  displayedDirtyCount: number;
  displayedSavedCount: number;
  globalPendingCount: number;
  setStudentSearch: (value: string) => void;
  setListOrderMode: (value: "prioritize" | "original") => void;
  setTakeListFilter: (value: TakeListFilter) => void;
  setShowMobileTakeControls: (value: boolean | ((prev: boolean) => boolean)) => void;
}

export function useAttendanceTakeMode({ students, isRowDirty }: UseAttendanceTakeModeParams): UseAttendanceTakeModeResult {
  const [studentSearch, setStudentSearch] = useState("");
  const [listOrderMode, setListOrderMode] = useState<"prioritize" | "original">("prioritize");
  const [takeListFilter, setTakeListFilter] = useState<TakeListFilter>("all");
  const [showMobileTakeControls, setShowMobileTakeControls] = useState(false);

  const displayedStudentsBase = useMemo(() => {
    const normalizedSearchText = studentSearch.trim().toLowerCase();
    const searchedStudents = normalizedSearchText
      ? students.filter((student) => student.full_name.toLowerCase().includes(normalizedSearchText))
      : students;

    if (takeListFilter === "pending") {
      return searchedStudents.filter((student) => isRowDirty(student.student_id));
    }
    if (takeListFilter === "saved") {
      return searchedStudents.filter((student) => !isRowDirty(student.student_id));
    }
    return searchedStudents;
  }, [students, studentSearch, takeListFilter, isRowDirty]);

  const displayedStudents = useMemo(() => {
    if (listOrderMode === "original") {
      return displayedStudentsBase;
    }
    const unsavedStudents = displayedStudentsBase.filter((student) => isRowDirty(student.student_id));
    const savedStudents = displayedStudentsBase.filter((student) => !isRowDirty(student.student_id));
    return [...unsavedStudents, ...savedStudents];
  }, [displayedStudentsBase, isRowDirty, listOrderMode]);

  const displayedDirtyCount = useMemo(
    () => displayedStudents.filter((student) => isRowDirty(student.student_id)).length,
    [displayedStudents, isRowDirty]
  );
  const displayedSavedCount = displayedStudents.length - displayedDirtyCount;
  const globalPendingCount = useMemo(
    () => students.length - students.filter((student) => !isRowDirty(student.student_id)).length,
    [students, isRowDirty]
  );

  return {
    studentSearch,
    listOrderMode,
    takeListFilter,
    showMobileTakeControls,
    displayedStudents,
    displayedDirtyCount,
    displayedSavedCount,
    globalPendingCount,
    setStudentSearch,
    setListOrderMode,
    setTakeListFilter,
    setShowMobileTakeControls,
  };
}