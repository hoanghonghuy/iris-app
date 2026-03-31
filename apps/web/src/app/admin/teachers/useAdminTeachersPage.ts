import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { adminApi } from "@/lib/api/admin.api";
import { Class, Pagination, School, Teacher } from "@/types";
import { fetchCollectionWithState, loadListWithDefaultSelection } from "@/lib/list-loaders";

type AssignModalState = {
  isOpen: boolean;
  teacherId: string | null;
  teacherName: string | null;
};

type UnassignAlertState = {
  isOpen: boolean;
  teacherId: string | null;
  classId: string | null;
  className: string | null;
};

const INITIAL_ASSIGN_MODAL_STATE: AssignModalState = {
  isOpen: false,
  teacherId: null,
  teacherName: null,
};

const INITIAL_UNASSIGN_ALERT_STATE: UnassignAlertState = {
  isOpen: false,
  teacherId: null,
  classId: null,
  className: null,
};

export function useAdminTeachersPage() {
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [actionLoading, setActionLoading] = useState(false);

  const [assignModal, setAssignModal] = useState<AssignModalState>(INITIAL_ASSIGN_MODAL_STATE);
  const [unassignAlert, setUnassignAlert] = useState<UnassignAlertState>(INITIAL_UNASSIGN_ALERT_STATE);

  const closeAssignModal = useCallback(() => {
    setAssignModal(INITIAL_ASSIGN_MODAL_STATE);
  }, []);

  const closeUnassignAlert = useCallback(() => {
    setUnassignAlert(INITIAL_UNASSIGN_ALERT_STATE);
  }, []);

  const fetchTeachers = useCallback(async () => {
    await fetchCollectionWithState({
      fetcher: () => adminApi.getTeachers({ limit: 20, offset: currentOffset }),
      setItems: setTeachers,
      fallbackError: "Không thể tải danh sách giáo viên",
      setLoading,
      setError,
      setPagination: (value) => setPagination(value as Pagination),
      onErrorMessage: (message) => toast.error(message),
    });
  }, [currentOffset]);

  useEffect(() => {
    void fetchTeachers();
  }, [fetchTeachers]);

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
        onEmpty: () => setSelectedClassId(""),
        onError: () => setClasses([]),
      });
    };

    void loadClasses();
  }, [selectedSchoolId]);

  const handleAssign = useCallback(async () => {
    if (!selectedClassId || !assignModal.teacherId) {
      return;
    }

    await adminApi.assignTeacherToClass(assignModal.teacherId, selectedClassId);
    closeAssignModal();
    await fetchTeachers();
  }, [assignModal.teacherId, closeAssignModal, fetchTeachers, selectedClassId]);

  const confirmUnassign = useCallback(async () => {
    if (!unassignAlert.teacherId || !unassignAlert.classId) {
      return;
    }

    await adminApi.unassignTeacherFromClass(unassignAlert.teacherId, unassignAlert.classId);
    closeUnassignAlert();
    await fetchTeachers();
  }, [closeUnassignAlert, fetchTeachers, unassignAlert.classId, unassignAlert.teacherId]);

  const filteredTeachers = useMemo(() => {
    if (!searchQuery.trim()) {
      return teachers;
    }

    const normalizedQuery = searchQuery.toLowerCase();
    return teachers.filter(
      (teacher) =>
        teacher.full_name?.toLowerCase().includes(normalizedQuery) ||
        teacher.email?.toLowerCase().includes(normalizedQuery) ||
        teacher.phone?.includes(normalizedQuery)
    );
  }, [teachers, searchQuery]);

  return {
    teachers,
    searchQuery,
    loading,
    error,
    pagination,
    schools,
    classes,
    selectedSchoolId,
    selectedClassId,
    actionLoading,
    assignModal,
    unassignAlert,
    filteredTeachers,
    setSearchQuery,
    setCurrentOffset,
    setSelectedSchoolId,
    setSelectedClassId,
    setActionLoading,
    setAssignModal,
    setUnassignAlert,
    closeAssignModal,
    closeUnassignAlert,
    handleAssign,
    confirmUnassign,
  };
}