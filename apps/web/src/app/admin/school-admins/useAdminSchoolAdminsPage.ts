import { FormEvent, useCallback, useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { adminApi } from "@/lib/api/admin.api";
import { useAuth } from "@/providers/AuthProvider";
import { Pagination, School, SchoolAdmin, UserInfo } from "@/types";
import { extractApiErrorMessage } from "@/lib/api-error";
import { fetchCollectionWithState } from "@/lib/list-loaders";

type DeleteAlertState = {
  isOpen: boolean;
  adminId: string | null;
};

const INITIAL_DELETE_ALERT_STATE: DeleteAlertState = {
  isOpen: false,
  adminId: null,
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
  const [selectedSchoolId, setSelectedSchoolId] = useState("");

  // User search state
  const [userSearchQuery, setUserSearchQuery] = useState("");
  const [userSearchResults, setUserSearchResults] = useState<UserInfo[]>([]);
  const [userSearchLoading, setUserSearchLoading] = useState(false);
  const [selectedUser, setSelectedUser] = useState<UserInfo | null>(null);
  const [showUserDropdown, setShowUserDropdown] = useState(false);
  const searchTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [success, setSuccess] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [deleteAlert, setDeleteAlert] = useState<DeleteAlertState>(INITIAL_DELETE_ALERT_STATE);

  const closeDeleteAlert = useCallback(() => {
    setDeleteAlert(INITIAL_DELETE_ALERT_STATE);
  }, []);

  const fetchAdmins = useCallback(async () => {
    if (role !== "SUPER_ADMIN") {
      return;
    }

    await fetchCollectionWithState({
      fetcher: () => adminApi.getSchoolAdmins({ limit: 20, offset: currentOffset }),
      setItems: setAdmins,
      fallbackError: "Không thể tải danh sách",
      setLoading,
      setError,
      setPagination: (value) => setPagination(value as Pagination),
    });
  }, [currentOffset, role]);

  useEffect(() => {
    if (role && role !== "SUPER_ADMIN") {
      router.replace("/admin");
    } else if (role === "SUPER_ADMIN") {
      void fetchAdmins();
    }
  }, [fetchAdmins, role, router]);

  // Load schools khi mở form
  useEffect(() => {
    if (!showForm) {
      return;
    }

    const loadSchools = async () => {
      try {
        const schoolResponse = await adminApi.getSchools();
        const schoolData = schoolResponse.data;
        setSchools(schoolData || []);
        if (schoolData && schoolData.length > 0) {
          setSelectedSchoolId(schoolData[0].school_id);
        }
      } catch {}
    };

    void loadSchools();
    // Reset user search khi mở form
    setUserSearchQuery("");
    setUserSearchResults([]);
    setSelectedUser(null);
    setShowUserDropdown(false);
  }, [showForm]);

  // Debounce search users
  useEffect(() => {
    if (searchTimerRef.current) {
      clearTimeout(searchTimerRef.current);
    }

    const query = userSearchQuery.trim();
    if (query.length < 2) {
      setUserSearchResults([]);
      setShowUserDropdown(false);
      return;
    }

    searchTimerRef.current = setTimeout(async () => {
      try {
        setUserSearchLoading(true);
        const response = await adminApi.getUsers({ limit: 10 });
        const allUsers = response.data || [];
        // Lọc phía client theo email (API chưa có search param)
        const filtered = allUsers.filter(
          (u) => u.email.toLowerCase().includes(query.toLowerCase())
        );
        setUserSearchResults(filtered);
        setShowUserDropdown(filtered.length > 0);
      } catch {
        setUserSearchResults([]);
      } finally {
        setUserSearchLoading(false);
      }
    }, 300);

    return () => {
      if (searchTimerRef.current) {
        clearTimeout(searchTimerRef.current);
      }
    };
  }, [userSearchQuery]);

  const selectUser = useCallback((user: UserInfo) => {
    setSelectedUser(user);
    setUserSearchQuery(user.email);
    setShowUserDropdown(false);
  }, []);

  const clearSelectedUser = useCallback(() => {
    setSelectedUser(null);
    setUserSearchQuery("");
    setUserSearchResults([]);
    setShowUserDropdown(false);
  }, []);

  const handleCreate = useCallback(async (event: FormEvent) => {
    event.preventDefault();
    if (!selectedUser || !selectedSchoolId) {
      setFormError("Chọn user và trường");
      return;
    }

    try {
      setSubmitting(true);
      setFormError("");
      setSuccess("");
      await adminApi.createSchoolAdmin({ user_id: selectedUser.user_id, school_id: selectedSchoolId });
      setSuccess("Đã tạo School Admin thành công!");
      setShowForm(false);
      clearSelectedUser();
      await fetchAdmins();
    } catch (errorValue: unknown) {
      setFormError(extractApiErrorMessage(errorValue, "Không thể tạo"));
    } finally {
      setSubmitting(false);
    }
  }, [clearSelectedUser, fetchAdmins, selectedSchoolId, selectedUser]);

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
      closeDeleteAlert();
    }
  }, [closeDeleteAlert, deleteAlert.adminId, fetchAdmins]);

  return {
    role,
    admins,
    loading,
    error,
    pagination,
    showForm,
    schools,
    selectedSchoolId,
    // User search
    userSearchQuery,
    userSearchResults,
    userSearchLoading,
    selectedUser,
    showUserDropdown,
    submitting,
    formError,
    success,
    deletingId,
    deleteAlert,
    setCurrentOffset,
    setShowForm,
    setSelectedSchoolId,
    setUserSearchQuery,
    selectUser,
    clearSelectedUser,
    setSuccess,
    setDeleteAlert,
    closeDeleteAlert,
    handleCreate,
    confirmDelete,
  };
}