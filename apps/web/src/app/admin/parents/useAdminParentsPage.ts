import { useCallback, useEffect, useMemo, useState } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Class, Pagination, Parent, School, Student } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";
import { fetchCollectionWithState, loadListWithDefaultSelection } from "@/lib/list-loaders";

type AssignModalState = {
  isOpen: boolean;
  parentId: string | null;
  parentName: string | null;
};

type UnassignAlertState = {
  isOpen: boolean;
  parentId: string | null;
  studentId: string | null;
  studentName: string | null;
};

const INITIAL_ASSIGN_MODAL_STATE: AssignModalState = {
  isOpen: false,
  parentId: null,
  parentName: null,
};

const INITIAL_UNASSIGN_ALERT_STATE: UnassignAlertState = {
  isOpen: false,
  parentId: null,
  studentId: null,
  studentName: null,
};

export function useAdminParentsPage() {
  const [parents, setParents] = useState<Parent[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [students, setStudents] = useState<Student[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [selectedStudentId, setSelectedStudentId] = useState("");
  const [actionLoading, setActionLoading] = useState(false);
  const [success, setSuccess] = useState("");

  const [assignModal, setAssignModal] = useState<AssignModalState>(INITIAL_ASSIGN_MODAL_STATE);
  const [unassignAlert, setUnassignAlert] = useState<UnassignAlertState>(INITIAL_UNASSIGN_ALERT_STATE);

  const closeAssignModal = useCallback(() => {
    setAssignModal(INITIAL_ASSIGN_MODAL_STATE);
  }, []);

  const closeUnassignAlert = useCallback(() => {
    setUnassignAlert(INITIAL_UNASSIGN_ALERT_STATE);
  }, []);

  const fetchParents = useCallback(async () => {
    await fetchCollectionWithState({
      fetcher: () => adminApi.getParents({ limit: 20, offset: currentOffset }),
      setItems: setParents,
      fallbackError: "Không thể tải danh sách phụ huynh",
      setLoading,
      setError,
      setPagination: (value) => setPagination(value as Pagination),
    });
  }, [currentOffset]);

  useEffect(() => {
    void fetchParents();
  }, [fetchParents]);

  useEffect(() => {
    const loadSchools = async () => {
      await loadListWithDefaultSelection({
        fetchList: async () => (await adminApi.getSchools()).data,
        setList: setSchools,
        setSelectedId: setSelectedSchoolId,
        getId: (school) => school.school_id,
      });
    };

    void loadSchools();
  }, []);

  useEffect(() => {
    if (!selectedSchoolId) {
      return;
    }

    const loadClasses = async () => {
      await loadListWithDefaultSelection({
        fetchList: async () => (await adminApi.getClassesBySchool(selectedSchoolId)).data,
        setList: setClasses,
        setSelectedId: setSelectedClassId,
        getId: (classItem) => classItem.class_id,
        onEmpty: () => {
          setSelectedClassId("");
          setStudents([]);
        },
        onError: () => setClasses([]),
      });
    };

    void loadClasses();
  }, [selectedSchoolId]);

  useEffect(() => {
    if (!selectedClassId) {
      return;
    }

    const loadStudents = async () => {
      await loadListWithDefaultSelection({
        fetchList: async () => (await adminApi.getStudentsByClass(selectedClassId)).data,
        setList: setStudents,
        setSelectedId: setSelectedStudentId,
        getId: (student) => student.student_id,
        onEmpty: () => setSelectedStudentId(""),
        onError: () => setStudents([]),
      });
    };

    void loadStudents();
  }, [selectedClassId]);

  const handleAssign = useCallback(async () => {
    if (!selectedStudentId || !assignModal.parentId) {
      return;
    }

    try {
      setActionLoading(true);
      setSuccess("");
      await adminApi.assignParentToStudent(assignModal.parentId, selectedStudentId);
      const studentName = students.find((student) => student.student_id === selectedStudentId)?.full_name || "";
      setSuccess(`Đã gán phụ huynh cho ${studentName}`);
      closeAssignModal();
      await fetchParents();
    } catch (err: unknown) {
      setError(extractApiErrorMessage(err, "Không thể gán"));
    } finally {
      setActionLoading(false);
    }
  }, [assignModal.parentId, closeAssignModal, fetchParents, selectedStudentId, students]);

  const confirmUnassign = useCallback(async () => {
    if (!unassignAlert.parentId || !unassignAlert.studentId) {
      return;
    }

    try {
      setActionLoading(true);
      setSuccess("");
      await adminApi.unassignParentFromStudent(unassignAlert.parentId, unassignAlert.studentId);
      setSuccess(`Đã hủy gán học sinh ${unassignAlert.studentName}`);
      closeUnassignAlert();
      await fetchParents();
    } catch (err: unknown) {
      setError(extractApiErrorMessage(err, "Không thể hủy gán"));
    } finally {
      setActionLoading(false);
    }
  }, [closeUnassignAlert, fetchParents, unassignAlert.parentId, unassignAlert.studentId, unassignAlert.studentName]);

  const filteredParents = useMemo(() => {
    if (!searchQuery.trim()) {
      return parents;
    }

    const normalizedQuery = searchQuery.toLowerCase();
    return parents.filter(
      (parent) =>
        parent.full_name?.toLowerCase().includes(normalizedQuery) ||
        parent.email?.toLowerCase().includes(normalizedQuery) ||
        parent.phone?.includes(normalizedQuery)
    );
  }, [parents, searchQuery]);

  return {
    parents,
    searchQuery,
    loading,
    error,
    pagination,
    schools,
    classes,
    students,
    selectedSchoolId,
    selectedClassId,
    selectedStudentId,
    actionLoading,
    success,
    assignModal,
    unassignAlert,
    filteredParents,
    setSearchQuery,
    setCurrentOffset,
    setSelectedSchoolId,
    setSelectedClassId,
    setSelectedStudentId,
    setAssignModal,
    setUnassignAlert,
    closeAssignModal,
    closeUnassignAlert,
    handleAssign,
    confirmUnassign,
  };
}