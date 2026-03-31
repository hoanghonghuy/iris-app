import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { adminApi } from "@/lib/api/admin.api";
import { Class, School, Student } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";
import { fetchCollectionWithState, loadListWithDefaultSelection } from "@/lib/list-loaders";

type StudentFormData = {
  full_name: string;
  dob: string;
  gender: "male" | "female" | "other";
};

type RevokeAlertState = {
  isOpen: boolean;
  studentId: string | null;
};

const INITIAL_REVOKE_ALERT_STATE: RevokeAlertState = {
  isOpen: false,
  studentId: null,
};

export const genderLabel: Record<string, string> = {
  male: "Nam",
  female: "Nữ",
  other: "Khác",
};

export function useAdminStudentsPage() {
  const [schools, setSchools] = useState<School[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedClassId, setSelectedClassId] = useState("");
  const [loadingSchools, setLoadingSchools] = useState(true);

  const [students, setStudents] = useState<Student[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState<StudentFormData>({ full_name: "", dob: "", gender: "male" });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");

  const [generatingCode, setGeneratingCode] = useState<string | null>(null);
  const [revokingCode, setRevokingCode] = useState<string | null>(null);
  const [revokeAlert, setRevokeAlert] = useState<RevokeAlertState>(INITIAL_REVOKE_ALERT_STATE);
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [codeError, setCodeError] = useState("");

  const closeRevokeAlert = useCallback(() => {
    setRevokeAlert(INITIAL_REVOKE_ALERT_STATE);
  }, []);

  useEffect(() => {
    const loadSchools = async () => {
      await loadListWithDefaultSelection({
        fetchList: async () => (await adminApi.getSchools()).data,
        setList: setSchools,
        setSelectedId: setSelectedSchoolId,
        getId: (school) => school.school_id,
        onError: () => setError("Không thể tải danh sách trường"),
        onFinally: () => setLoadingSchools(false),
      });
    };

    void loadSchools();
  }, []);

  useEffect(() => {
    if (!selectedSchoolId) {
      return;
    }

    const loadClasses = async () => {
      setSelectedClassId("");
      setStudents([]);
      setSearchQuery("");

      await loadListWithDefaultSelection({
        fetchList: async () => (await adminApi.getClassesBySchool(selectedSchoolId)).data,
        setList: setClasses,
        setSelectedId: setSelectedClassId,
        getId: (classItem) => classItem.class_id,
        onError: () => setClasses([]),
      });
    };

    void loadClasses();
  }, [selectedSchoolId]);

  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) {
      return;
    }

    await fetchCollectionWithState({
      fetcher: () => adminApi.getStudentsByClass(selectedClassId),
      setItems: setStudents,
      fallbackError: "Không thể tải danh sách học sinh",
      setLoading: setLoadingStudents,
      setError,
    });
  }, [selectedClassId]);

  useEffect(() => {
    void fetchStudents();
  }, [fetchStudents]);

  const filteredStudents = useMemo(() => {
    if (!searchQuery.trim()) {
      return students;
    }

    const normalizedQuery = searchQuery.toLowerCase();
    return students.filter((student) => student.full_name.toLowerCase().includes(normalizedQuery));
  }, [students, searchQuery]);

  const handleCreate = useCallback(async (event: FormEvent) => {
    event.preventDefault();

    if (!formData.full_name.trim()) {
      setFormError("Họ tên không được để trống");
      return;
    }

    try {
      setSubmitting(true);
      setFormError("");
      await adminApi.createStudent({
        school_id: selectedSchoolId,
        class_id: selectedClassId,
        full_name: formData.full_name,
        dob: formData.dob,
        gender: formData.gender,
      });

      setFormData({ full_name: "", dob: "", gender: "male" });
      setShowForm(false);
      await fetchStudents();
    } catch (errorValue: unknown) {
      setFormError(extractApiErrorMessage(errorValue, "Không thể tạo học sinh"));
    } finally {
      setSubmitting(false);
    }
  }, [fetchStudents, formData, selectedClassId, selectedSchoolId]);

  const handleGenerateCode = useCallback(async (studentId: string) => {
    try {
      setGeneratingCode(studentId);
      setCodeError("");
      const response = await adminApi.generateParentCode(studentId);
      const parentCode = response.data?.parent_code || "";
      const expiresAt = response.data?.expires_at || "";

      setStudents((prev) => prev.map((student) => (
        student.student_id === studentId
          ? { ...student, active_parent_code: parentCode, code_expires_at: expiresAt }
          : student
      )));
    } catch (errorValue: unknown) {
      setCodeError(extractApiErrorMessage(errorValue, "Không thể tạo mã"));
    } finally {
      setGeneratingCode(null);
    }
  }, []);

  const confirmRevokeCode = useCallback(async () => {
    if (!revokeAlert.studentId) {
      return;
    }

    try {
      setRevokingCode(revokeAlert.studentId);
      setCodeError("");
      await adminApi.revokeParentCode(revokeAlert.studentId);

      setStudents((prev) => prev.map((student) => (
        student.student_id === revokeAlert.studentId
          ? { ...student, active_parent_code: undefined, code_expires_at: undefined }
          : student
      )));

      closeRevokeAlert();
    } catch (errorValue: unknown) {
      setCodeError(extractApiErrorMessage(errorValue, "Không thể thu hồi mã"));
    } finally {
      setRevokingCode(null);
    }
  }, [closeRevokeAlert, revokeAlert.studentId]);

  const handleCopy = useCallback((code: string, studentId: string) => {
    void navigator.clipboard.writeText(code);
    setCopiedId(studentId);
    setTimeout(() => setCopiedId(null), 2000);
  }, []);

  const getDaysLeft = useCallback((dateString?: string) => {
    if (!dateString) {
      return null;
    }

    const diff = new Date(dateString).getTime() - new Date().getTime();
    if (diff <= 0) {
      return "Hết hạn";
    }

    return `Còn ${Math.ceil(diff / (1000 * 3600 * 24))} ngày`;
  }, []);

  const selectedClassName = useMemo(
    () => classes.find((classInfo) => classInfo.class_id === selectedClassId)?.name || "",
    [classes, selectedClassId]
  );

  return {
    schools,
    classes,
    selectedSchoolId,
    selectedClassId,
    loadingSchools,
    students,
    searchQuery,
    loadingStudents,
    error,
    showForm,
    formData,
    submitting,
    formError,
    generatingCode,
    revokingCode,
    revokeAlert,
    copiedId,
    codeError,
    filteredStudents,
    selectedClassName,
    setSelectedSchoolId,
    setSelectedClassId,
    setSearchQuery,
    setShowForm,
    setFormData,
    setRevokeAlert,
    closeRevokeAlert,
    handleCreate,
    handleGenerateCode,
    confirmRevokeCode,
    handleCopy,
    getDaysLeft,
  };
}