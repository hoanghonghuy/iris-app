import { FormEvent, useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { adminApi } from "@/lib/api/admin.api";
import { useAuth } from "@/providers/AuthProvider";
import { Pagination, School, SchoolAdmin, UserInfo } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";

type DeleteAlertState = {
  isOpen: boolean;
  adminId: string | null;
};

export function useAdminSchoolAdminsPage() {
  const { role } = useAuth();
  const router = useRouter();

  const [admins, setAdmins] = useState<SchoolAdmin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 20, offset: 0, has_more: false });
  const [currentOffset, setCurrentOffset] = useState(0);

  const [showForm, setShowForm] = useState(false);
  const [schools, setSchools] = useState<School[]>([]);
  const [users, setUsers] = useState<UserInfo[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState("");
  const [selectedUserId, setSelectedUserId] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [deleteAlert, setDeleteAlert] = useState<DeleteAlertState>({ isOpen: false, adminId: null });

  const fetchAdmins = useCallback(async () => {
    if (role !== "SUPER_ADMIN") {
      return;
    }

    try {
      setLoading(true);
      setError("");
      const response = await adminApi.getSchoolAdmins({ limit: 20, offset: currentOffset });
      setAdmins(response.data || []);
      if (response.pagination) {
        setPagination(response.pagination);
      }
    } catch (errorValue: unknown) {
      setError(extractApiErrorMessage(errorValue, "Không thể tải danh sách"));
    } finally {
      setLoading(false);
    }
  }, [currentOffset, role]);

  useEffect(() => {
    if (role && role !== "SUPER_ADMIN") {
      router.replace("/admin");
    } else if (role === "SUPER_ADMIN") {
      void fetchAdmins();
    }
  }, [fetchAdmins, role, router]);

  useEffect(() => {
    if (!showForm) {
      return;
    }

    const loadFormData = async () => {
      try {
        const [schoolResponse, userData] = await Promise.all([
          adminApi.getSchools(),
          adminApi.getUsers({ limit: 100 }),
        ]);
        const schoolData = schoolResponse.data;
        setSchools(schoolData || []);

        const userList = userData.data || [];
        setUsers(Array.isArray(userList) ? userList : []);

        if (schoolData && schoolData.length > 0) {
          setSelectedSchoolId(schoolData[0].school_id);
        }
        if (userList && userList.length > 0) {
          setSelectedUserId(userList[0].user_id);
        }
      } catch {
      }
    };

    void loadFormData();
  }, [showForm]);

  const handleCreate = useCallback(async (event: FormEvent) => {
    event.preventDefault();
    if (!selectedUserId || !selectedSchoolId) {
      setFormError("Chọn user và trường");
      return;
    }

    try {
      setSubmitting(true);
      setFormError("");
      setSuccess("");
      await adminApi.createSchoolAdmin({ user_id: selectedUserId, school_id: selectedSchoolId });
      setSuccess("Đã tạo School Admin thành công!");
      setShowForm(false);
      await fetchAdmins();
    } catch (errorValue: unknown) {
      setFormError(extractApiErrorMessage(errorValue, "Không thể tạo"));
    } finally {
      setSubmitting(false);
    }
  }, [fetchAdmins, selectedSchoolId, selectedUserId]);

  const confirmDelete = useCallback(async () => {
    if (!deleteAlert.adminId) {
      return;
    }

    try {
      setDeletingId(deleteAlert.adminId);
      await adminApi.deleteSchoolAdmin(deleteAlert.adminId);
      await fetchAdmins();
    } catch (errorValue: unknown) {
      setError(extractApiErrorMessage(errorValue, "Không thể xóa"));
    } finally {
      setDeletingId(null);
      setDeleteAlert({ isOpen: false, adminId: null });
    }
  }, [deleteAlert.adminId, fetchAdmins]);

  return {
    role,
    admins,
    loading,
    error,
    pagination,
    showForm,
    schools,
    users,
    selectedSchoolId,
    selectedUserId,
    submitting,
    formError,
    success,
    deletingId,
    deleteAlert,
    setCurrentOffset,
    setShowForm,
    setSelectedSchoolId,
    setSelectedUserId,
    setSuccess,
    setDeleteAlert,
    handleCreate,
    confirmDelete,
  };
}