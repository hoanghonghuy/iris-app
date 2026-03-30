import { useCallback, useEffect, useMemo, useState } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Class, Pagination, Parent, School, Student } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";

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

  const [assignModal, setAssignModal] = useState<AssignModalState>({ isOpen: false, parentId: null, parentName: null });
  const [unassignAlert, setUnassignAlert] = useState<UnassignAlertState>({
    isOpen: false,
    parentId: null,
    studentId: null,
    studentName: null,
  });

  const fetchParents = useCallback(async () => {
    try {
      setLoading(true);
      setError("");
      const response = await adminApi.getParents({ limit: 20, offset: currentOffset });
      setParents(response.data || []);
      if (response.pagination) {
        setPagination(response.pagination);
      }
    } catch (err: unknown) {
      setError(extractApiErrorMessage(err, "Không thể tải danh sách phụ huynh"));
    } finally {
      setLoading(false);
    }
  }, [currentOffset]);

  useEffect(() => {
    void fetchParents();
  }, [fetchParents]);

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
          setStudents([]);
        }
      } catch {
        setClasses([]);
      }
    };

    void loadClasses();
  }, [selectedSchoolId]);

  useEffect(() => {
    if (!selectedClassId) {
      return;
    }

    const loadStudents = async () => {
      try {
        const response = await adminApi.getStudentsByClass(selectedClassId);
        const studentData = response.data;
        setStudents(studentData || []);

        if (studentData && studentData.length > 0) {
          setSelectedStudentId(studentData[0].student_id);
        } else {
          setSelectedStudentId("");
        }
      } catch {
        setStudents([]);
      }
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
      setAssignModal({ isOpen: false, parentId: null, parentName: null });
      await fetchParents();
    } catch (err: unknown) {
      setError(extractApiErrorMessage(err, "Không thể gán"));
    } finally {
      setActionLoading(false);
    }
  }, [assignModal.parentId, fetchParents, selectedStudentId, students]);

  const confirmUnassign = useCallback(async () => {
    if (!unassignAlert.parentId || !unassignAlert.studentId) {
      return;
    }

    try {
      setActionLoading(true);
      setSuccess("");
      await adminApi.unassignParentFromStudent(unassignAlert.parentId, unassignAlert.studentId);
      setSuccess(`Đã hủy gán học sinh ${unassignAlert.studentName}`);
      setUnassignAlert({ isOpen: false, parentId: null, studentId: null, studentName: null });
      await fetchParents();
    } catch (err: unknown) {
      setError(extractApiErrorMessage(err, "Không thể hủy gán"));
    } finally {
      setActionLoading(false);
    }
  }, [fetchParents, unassignAlert.parentId, unassignAlert.studentId, unassignAlert.studentName]);

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
    handleAssign,
    confirmUnassign,
  };
}