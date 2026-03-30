import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { adminApi } from "@/lib/api/admin.api";
import { Class, Pagination, School, Teacher } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";

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

  const [assignModal, setAssignModal] = useState<AssignModalState>({
    isOpen: false,
    teacherId: null,
    teacherName: null,
  });
  const [unassignAlert, setUnassignAlert] = useState<UnassignAlertState>({
    isOpen: false,
    teacherId: null,
    classId: null,
    className: null,
  });

  const fetchTeachers = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const response = await adminApi.getTeachers({ limit: 20, offset: currentOffset });
      setTeachers(response.data || []);
      if (response.pagination) {
        setPagination(response.pagination);
      }
    } catch (err: unknown) {
      const message = extractApiErrorMessage(err, "Không thể tải danh sách giáo viên");
      setError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  }, [currentOffset]);

  useEffect(() => {
    void fetchTeachers();
  }, [fetchTeachers]);

  useEffect(() => {
    const loadSchools = async () => {
      try {
        const response = await adminApi.getSchools();
        const schoolData = response.data;
        setSchools(schoolData || []);
        if (schoolData && schoolData.length > 0) {
          setSelectedSchoolId(schoolData[0].school_id);
        }
      } catch {
      }
    };

    void loadSchools();
  }, []);

  useEffect(() => {
    if (!selectedSchoolId) {
      return;
    }

    const loadClasses = async () => {
      try {
        const response = await adminApi.getClassesBySchool(selectedSchoolId);
        const classData = response.data;
        setClasses(classData || []);

        if (classData && classData.length > 0) {
          setSelectedClassId(classData[0].class_id);
        } else {
          setSelectedClassId("");
        }
      } catch {
        setClasses([]);
      }
    };

    void loadClasses();
  }, [selectedSchoolId]);

  const handleAssign = useCallback(async () => {
    if (!selectedClassId || !assignModal.teacherId) {
      return;
    }

    await adminApi.assignTeacherToClass(assignModal.teacherId, selectedClassId);
    setAssignModal({ isOpen: false, teacherId: null, teacherName: null });
    await fetchTeachers();
  }, [assignModal.teacherId, fetchTeachers, selectedClassId]);

  const confirmUnassign = useCallback(async () => {
    if (!unassignAlert.teacherId || !unassignAlert.classId) {
      return;
    }

    await adminApi.unassignTeacherFromClass(unassignAlert.teacherId, unassignAlert.classId);
    setUnassignAlert({ isOpen: false, teacherId: null, classId: null, className: null });
    await fetchTeachers();
  }, [fetchTeachers, unassignAlert.classId, unassignAlert.teacherId]);

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
    handleAssign,
    confirmUnassign,
  };
}