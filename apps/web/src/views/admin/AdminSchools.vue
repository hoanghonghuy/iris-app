<script setup>
import { onMounted } from 'vue'
import { adminService } from '../../services/adminService'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import PaginationBar from '../../components/common/PaginationBar.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'
import { useAdminCrudList } from '../../composables/admin/useAdminCrudList'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_LOADING_MESSAGE,
  ADMIN_RETRY_BUTTON_TEXT,
} from '../../helpers/adminConfig'

const PAGE_SIZE = 10

const {
  items: schools,
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
  fetchItems: fetchSchools,
  openAddModal,
  closeModal,
  openEditModal,
  handleSave,
  confirmDelete,
  closeDeleteConfirm,
  handleDelete,
} = useAdminCrudList({
  pageSize: PAGE_SIZE,
  fetchPage: ({ page, pageSize }) =>
    adminService.getSchools({
      limit: pageSize,
      offset: (page - 1) * pageSize,
    }),
  createEmptyForm: () => ({ id: '', name: '', address: '' }),
  toEditForm: (school) => ({ ...school, id: school.school_id }),
  validateForm: (form) => {
    if (!form.name) {
      return 'Tên trường không được để trống'
    }

    return ''
  },
  createItem: (form) =>
    adminService.createSchool({
      name: form.name,
      address: form.address,
    }),
  updateItem: (form) =>
    adminService.updateSchool(form.id, {
      name: form.name,
      address: form.address,
    }),
  deleteItem: (school) => adminService.deleteSchool(school.school_id),
})

onMounted(() => {
  fetchSchools()
})
</script>

<template>
  <div class="admin-schools">
    <div class="flex justify-end items-center mb-6">
      <button class="btn btn--primary" type="button" @click="openAddModal">
        + Thêm trường học
      </button>
    </div>

    <div v-if="errorMessage" class="p-4 mb-6 bg-red-50 text-danger rounded border border-red-200">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchSchools(currentPage)">
        {{ ADMIN_RETRY_BUTTON_TEXT }}
      </button>
    </div>

    <LoadingSpinner v-else-if="isLoading" :message="ADMIN_LOADING_MESSAGE" />

    <div v-else class="card">
      <EmptyState
        v-if="schools.length === 0"
        title="Chưa có trường học nào"
        message="Hãy thêm trường học đầu tiên để bắt đầu quản lý."
      >
        <template #action>
          <button class="btn btn--primary" type="button" @click="openAddModal">
            Thêm trường học
          </button>
        </template>
      </EmptyState>

      <div v-else class="table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th>Tên trường</th>
              <th>Địa chỉ</th>
              <th class="text-right">Thao tác</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="school in schools" :key="school.school_id">
              <td class="font-medium">{{ school.name }}</td>
              <td>{{ school.address }}</td>
              <td class="text-right">
                <div class="flex justify-end gap-2">
                  <button
                    class="btn btn--sm btn--outline"
                    type="button"
                    @click="openEditModal(school)"
                  >
                    Sửa
                  </button>
                  <button
                    class="btn btn--sm btn--danger"
                    type="button"
                    @click="confirmDelete(school)"
                  >
                    Xóa
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="p-4">
          <PaginationBar
            :current-page="currentPage"
            :total-pages="totalPages"
            :total-items="totalItems"
            @page-change="fetchSchools"
          />
        </div>
      </div>
    </div>

    <ActionModal
      :is-open="isModalOpen"
      :title="formMode === 'add' ? 'Thêm trường học mới' : 'Sửa thông tin trường'"
      @close="closeModal"
    >
      <form @submit.prevent="handleSave" class="flex-col gap-4">
        <div
          v-if="modalError"
          class="p-3 bg-red-50 text-danger text-sm rounded border border-red-200"
        >
          {{ modalError }}
        </div>

        <div class="form-group">
          <label class="form-label" for="schoolName"
            >Tên trường <span class="text-danger">*</span></label
          >
          <input
            id="schoolName"
            v-model="formData.name"
            type="text"
            class="form-input"
            placeholder="Nhập tên trường học"
            :disabled="isSubmitting"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="schoolAddress">Địa chỉ</label>
          <textarea
            id="schoolAddress"
            v-model="formData.address"
            class="form-input"
            placeholder="Nhập địa chỉ"
            rows="3"
            :disabled="isSubmitting"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2 mt-4 pt-4 border-t">
          <button
            type="button"
            class="btn btn--outline"
            @click="closeModal"
            :disabled="isSubmitting"
          >
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
      :message="`Bạn có chắc chắn muốn xóa trường '${itemToDelete?.name}' không? Hành động này không thể hoàn tác.`"
      confirm-text="Xóa trường"
      is-danger
      :is-loading="isSubmitting"
      @confirm="handleDelete"
      @cancel="closeDeleteConfirm"
    />
  </div>
</template>

<style scoped>
.mb-6 {
  margin-bottom: var(--spacing-6);
}
.p-3 {
  padding: var(--spacing-3);
}
.p-4 {
  padding: var(--spacing-4);
}
.gap-2 {
  gap: var(--spacing-2);
}
.gap-4 {
  gap: var(--spacing-4);
}
.mt-2 {
  margin-top: var(--spacing-2);
}
.mt-4 {
  margin-top: var(--spacing-4);
}
.pt-4 {
  padding-top: var(--spacing-4);
}
.border-t {
  border-top: 1px solid var(--color-border);
}
.rounded {
  border-radius: var(--radius);
}
.bg-red-50 {
  background-color: var(--color-danger-soft-bg);
}
.border-red-200 {
  border-color: var(--color-danger-soft-border);
}
.text-right {
  text-align: right;
}

.table-responsive {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th,
.table td {
  padding: var(--spacing-3) var(--spacing-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.table th {
  background-color: var(--color-background);
  font-weight: 600;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  text-transform: uppercase;
}

.table tbody tr:hover {
  background-color: var(--color-background);
}
</style>
