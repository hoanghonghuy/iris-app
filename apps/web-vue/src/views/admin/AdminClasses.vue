<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { adminService } from '../../services/adminService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_LOADING_MESSAGE,
  ADMIN_RETRY_BUTTON_TEXT,
  ADMIN_SELECTOR_FETCH_LIMIT,
} from '../../helpers/adminConfig'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import PaginationBar from '../../components/common/PaginationBar.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'
import { useAdminCrudList } from '../../composables/admin/useAdminCrudList'

const authStore = useAuthStore()
const PAGE_SIZE = 10

const schools = ref([])
const selectedSchoolId = ref('')

const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')
const selectedSchoolName = computed(() => {
  return schools.value.find((school) => school.school_id === selectedSchoolId.value)?.name || ''
})

const {
  items: classesList,
  totalPages,
  currentPage,
  totalItems,
  isLoading,
  errorMessage,
  isModalOpen,
  isSubmitting,
  modalError,
  formMode,
  formData,
  isConfirmOpen,
  itemToDelete,
  fetchItems: fetchClasses,
  openAddModal: openCrudAddModal,
  closeModal,
  openEditModal,
  handleSave,
  confirmDelete,
  closeDeleteConfirm,
  handleDelete,
} = useAdminCrudList({
  pageSize: PAGE_SIZE,
  fetchPage: ({ page, pageSize }) => {
    if (!selectedSchoolId.value) {
      return {
        data: [],
        pagination: {
          total: 0,
          limit: pageSize,
          offset: 0,
          has_more: false,
        },
      }
    }

    return adminService.getClassesBySchool(selectedSchoolId.value, {
      limit: pageSize,
      offset: (page - 1) * pageSize,
    })
  },
  createEmptyForm: () => ({
    id: '',
    name: '',
    school_year: '2023-2024',
    school_id: selectedSchoolId.value,
  }),
  toEditForm: (cls) => ({
    ...cls,
    id: cls.class_id,
  }),
  validateForm: (form) => {
    if (!form.name || !form.school_year || !form.school_id) {
      return 'Vui lòng nhập đầy đủ thông tin'
    }

    return ''
  },
  createItem: (form) =>
    adminService.createClass({
      name: form.name,
      school_year: form.school_year,
      school_id: form.school_id,
    }),
  updateItem: (form) =>
    adminService.updateClass(form.id, {
      name: form.name,
      school_year: form.school_year,
    }),
  deleteItem: (cls) => adminService.deleteClass(cls.class_id),
  saveErrorMessage: 'Không thể lưu lớp học',
  deleteErrorPrefix: 'Không thể xóa lớp',
  onAfterSave: async ({ mode, form, currentPage, fetchItems }) => {
    if (mode === 'add' && selectedSchoolId.value !== form.school_id) {
      selectedSchoolId.value = form.school_id
      return true
    }

    await fetchItems(currentPage)
    return true
  },
})

async function fetchSchools() {
  try {
    const data = await adminService.getSchools({ limit: ADMIN_SELECTOR_FETCH_LIMIT, offset: 0 })
    schools.value = normalizeListResponse(data)

    if (schools.value.length > 0) {
      selectedSchoolId.value = schools.value[0].school_id
    } else {
      isLoading.value = false
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách trường'
    isLoading.value = false
  }
}

function openAddModal() {
  openCrudAddModal()
}

watch(selectedSchoolId, () => {
  fetchClasses(1)
})

onMounted(async () => {
  if (!authStore.currentUser && authStore.isAuthenticated) {
    await authStore.fetchCurrentUser()
  }

  await fetchSchools()
})
</script>

<template>
  <div class="admin-classes page-stack">
    <div class="page-actions">
      <button class="btn btn--primary" type="button" :disabled="schools.length === 0" @click="openAddModal">
        + Thêm lớp học
      </button>
    </div>

    <div class="card toolbar-card">
      <div v-if="isSuperAdmin" class="form-group mb-0">
        <label class="form-label" for="schoolFilter">Chọn trường học</label>
        <select id="schoolFilter" v-model="selectedSchoolId" class="form-input" :disabled="isLoading">
          <option value="" disabled>-- Chọn trường --</option>
          <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
            {{ school.name }}
          </option>
        </select>
      </div>

      <div v-else class="school-context">
        <p class="school-context__label">Trường hiện tại</p>
        <p class="school-context__name">{{ selectedSchoolName || 'Không xác định' }}</p>
      </div>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchClasses(currentPage)">{{ ADMIN_RETRY_BUTTON_TEXT }}</button>
    </div>

    <LoadingSpinner v-else-if="isLoading" :message="ADMIN_LOADING_MESSAGE" />

    <div v-else class="card classes-shell">
      <EmptyState
        v-if="schools.length === 0"
        title="Chưa có trường học nào"
        message="Bạn cần tạo trường học trước khi có thể thêm lớp."
      />

      <EmptyState
        v-else-if="classesList.length === 0"
        title="Chưa có lớp học nào"
        :message="`${selectedSchoolName || 'Trường này'} chưa có lớp học nào. Hãy thêm lớp đầu tiên.`"
      >
        <template #action>
          <button class="btn btn--primary" type="button" @click="openAddModal">
            Thêm lớp học
          </button>
        </template>
      </EmptyState>

      <template v-else>
        <div class="desktop-table table-responsive">
          <table class="table">
            <thead>
              <tr>
                <th>Tên lớp</th>
                <th>Năm học</th>
                <th class="text-right">Thao tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="cls in classesList" :key="cls.class_id">
                <td class="font-medium">{{ cls.name }}</td>
                <td><span class="badge badge--outline">{{ cls.school_year }}</span></td>
                <td class="text-right">
                  <div class="table-actions">
                    <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(cls)">
                      Sửa
                    </button>
                    <button class="btn btn--sm btn--danger" type="button" @click="confirmDelete(cls)">
                      Xóa
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="mobile-list">
          <article v-for="cls in classesList" :key="cls.class_id" class="mobile-card">
            <div class="mobile-card__head">
              <p class="mobile-card__title">{{ cls.name }}</p>
              <span class="badge badge--outline">{{ cls.school_year }}</span>
            </div>

            <div class="mobile-card__actions">
              <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(cls)">
                Sửa
              </button>
              <button class="btn btn--sm btn--danger" type="button" @click="confirmDelete(cls)">
                Xóa
              </button>
            </div>
          </article>
        </div>

        <PaginationBar
          :current-page="currentPage"
          :total-pages="totalPages"
          :total-items="totalItems"
          :limit="PAGE_SIZE"
          @page-change="fetchClasses"
        />
      </template>
    </div>

    <ActionModal
      :is-open="isModalOpen"
      :title="formMode === 'add' ? 'Thêm lớp học mới' : 'Sửa thông tin lớp'"
      @close="closeModal"
    >
      <form class="modal-form" @submit.prevent="handleSave">
        <div v-if="modalError" class="alert alert--error">
          {{ modalError }}
        </div>

        <div v-if="formMode === 'add' && isSuperAdmin" class="form-group mb-0">
          <label class="form-label" for="schoolSelect">Thuộc trường</label>
          <select id="schoolSelect" v-model="formData.school_id" class="form-input" :disabled="isSubmitting" required>
            <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
              {{ school.name }}
            </option>
          </select>
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="className">Tên lớp</label>
          <input
            id="className"
            v-model="formData.name"
            type="text"
            class="form-input"
            placeholder="Ví dụ: Lá Non"
            :disabled="isSubmitting"
            required
          />
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="schoolYear">Năm học</label>
          <input
            id="schoolYear"
            v-model="formData.school_year"
            type="text"
            class="form-input"
            placeholder="Ví dụ: 2025-2026"
            :disabled="isSubmitting"
            required
          />
        </div>

        <div class="modal-actions">
          <button type="button" class="btn btn--outline" :disabled="isSubmitting" @click="closeModal">
            Hủy
          </button>
          <button type="submit" class="btn btn--primary" :disabled="isSubmitting">
            {{ isSubmitting ? 'Đang lưu...' : 'Lưu lại' }}
          </button>
        </div>
      </form>
    </ActionModal>

    <ConfirmDialog
      :is-open="isConfirmOpen"
      title="Xác nhận xóa"
      :message="`Bạn có chắc chắn muốn xóa lớp '${itemToDelete?.name || ''}' không? Hành động này không thể hoàn tác.`"
      confirm-text="Xóa lớp"
      is-danger
      :is-loading="isSubmitting"
      @confirm="handleDelete"
      @cancel="closeDeleteConfirm"
    />
  </div>
</template>

<style scoped>
.page-stack,
.mobile-list,
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-actions {
  display: flex;
  justify-content: flex-end;
}

.toolbar-card,
.classes-shell,
.mobile-card {
  padding: var(--spacing-4);
}

.school-context__label,
.school-context__name,
.mobile-card__title,
.mobile-card__meta {
  margin: 0;
}

.school-context__label,
.mobile-card__meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.school-context__name,
.mobile-card__title {
  color: var(--color-text);
  font-weight: 700;
}

.school-context__name {
  margin-top: var(--spacing-1);
}

.desktop-table {
  display: block;
}

.mobile-list {
  display: none;
}

.table-actions,
.mobile-card__actions,
.modal-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.table-actions {
  justify-content: flex-end;
}

.mobile-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
}

.mobile-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.mobile-card__actions {
  margin-top: var(--spacing-3);
}

.modal-actions {
  justify-content: flex-end;
}

.mt-2 {
  margin-top: var(--spacing-2);
}

.text-right {
  text-align: right;
}

@media (max-width: 767px) {
  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: flex;
  }
}
</style>
