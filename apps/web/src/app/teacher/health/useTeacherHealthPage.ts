import { FormEvent, useCallback, useEffect, useState } from "react";
import { teacherApi } from "@/lib/api/teacher.api";
import { Class, HealthLog, Student } from "@/types";
import { loadListWithDefaultSelection } from "@/lib/list-loaders";
import { extractApiErrorRawMessage } from "@/lib/api-error";

export type Severity = "normal" | "watch" | "urgent";

export function useTeacherHealthPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [selectedClassId, setSelectedClassId] = useState("");
  const [students, setStudents] = useState<Student[]>([]);
  const [loadingClasses, setLoadingClasses] = useState(true);
  const [loadingStudents, setLoadingStudents] = useState(false);
  const [error, setError] = useState("");

  const [showForm, setShowForm] = useState(false);
  const [formStudentId, setFormStudentId] = useState("");
  const [temperature, setTemperature] = useState("");
  const [symptoms, setSymptoms] = useState("");
  const [severity, setSeverity] = useState<Severity>("normal");
  const [note, setNote] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");

  const [historyStudentId, setHistoryStudentId] = useState("");
  const [historyFrom, setHistoryFrom] = useState("");
  const [historyTo, setHistoryTo] = useState("");
  const [historyLogs, setHistoryLogs] = useState<HealthLog[]>([]);
  const [loadingHistory, setLoadingHistory] = useState(false);
  const [historyError, setHistoryError] = useState("");

  useEffect(() => {
    const loadClasses = async () => {
      await loadListWithDefaultSelection({
        fetchList: () => teacherApi.getMyClasses(),
        setList: setClasses,
        setSelectedId: setSelectedClassId,
        getId: (classItem) => classItem.class_id,
        onError: () => setError("Không thể tải lớp"),
        onFinally: () => setLoadingClasses(false),
      });
    };

    void loadClasses();
  }, []);

  const fetchStudents = useCallback(async () => {
    if (!selectedClassId) {
      return;
    }

    try {
      setLoadingStudents(true);
      setError("");
      const studentData = await teacherApi.getStudentsInClass(selectedClassId);
      setStudents(studentData || []);

      if (studentData && studentData.length > 0) {
        setFormStudentId(studentData[0].student_id);
        setHistoryStudentId(studentData[0].student_id);
      } else {
        setHistoryLogs([]);
      }
    } catch (err: unknown) {
      setError(extractApiErrorRawMessage(err) || "Không thể tải HS");
    } finally {
      setLoadingStudents(false);
    }
  }, [selectedClassId]);

  useEffect(() => {
    void fetchStudents();
  }, [fetchStudents]);

  const fetchHistory = useCallback(async () => {
    if (!historyStudentId) {
      setHistoryLogs([]);
      return;
    }

    try {
      setLoadingHistory(true);
      setHistoryError("");
      const logs = await teacherApi.getStudentHealth(
        historyStudentId,
        historyFrom || undefined,
        historyTo || undefined
      );
      setHistoryLogs(logs || []);
    } catch (err: unknown) {
      setHistoryError(extractApiErrorRawMessage(err) || "Không thể tải lịch sử sức khỏe");
    } finally {
      setLoadingHistory(false);
    }
  }, [historyFrom, historyStudentId, historyTo]);

  useEffect(() => {
    void fetchHistory();
  }, [fetchHistory]);

  const handleSubmit = useCallback(async (event: FormEvent) => {
    event.preventDefault();
    if (!formStudentId) {
      setFormError("Chọn học sinh");
      return;
    }

    try {
      setSubmitting(true);
      setFormError("");
      setSuccess("");
      await teacherApi.createHealthLog({
        student_id: formStudentId,
        temperature: temperature ? parseFloat(temperature) : undefined,
        symptoms: symptoms || undefined,
        severity,
        note: note || undefined,
      });
      setSuccess("Đã ghi nhận sức khỏe thành công!");
      setTemperature("");
      setSymptoms("");
      setSeverity("normal");
      setNote("");
      await fetchHistory();
    } catch (err: unknown) {
      setFormError(extractApiErrorRawMessage(err) || "Lỗi ghi nhận");
    } finally {
      setSubmitting(false);
    }
  }, [fetchHistory, formStudentId, note, severity, symptoms, temperature]);

  return {
    classes,
    selectedClassId,
    students,
    loadingClasses,
    loadingStudents,
    error,
    showForm,
    formStudentId,
    temperature,
    symptoms,
    severity,
    note,
    submitting,
    formError,
    success,
    historyStudentId,
    historyFrom,
    historyTo,
    historyLogs,
    loadingHistory,
    historyError,
    setSelectedClassId,
    setShowForm,
    setFormStudentId,
    setTemperature,
    setSymptoms,
    setSeverity,
    setNote,
    setHistoryStudentId,
    setHistoryFrom,
    setHistoryTo,
    fetchHistory,
    handleSubmit,
  };
}