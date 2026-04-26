<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { adminService } from '../../services/adminService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_LOADING_MESSAGE,
  ADMIN_RETRY_BUTTON_TEXT,
  ADMIN_SELECTOR_FETCH_LIMIT,
} from '../../helpers/adminConfig'
import { useAdminCrudList, useAdminUserSearch } from '../../composables/admin'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import PaginationBar from '../../components/common/PaginationBar.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'

const PAGE_SIZE = 10
const USER_SEARCH_MIN_LENGTH = 2
const USER_SEARCH_RESULT_LIMIT = 6

const schools = ref([])

const {
  searchQuery: userSearchQuery,
  searchResults: userSearchResults,
  searchLoading: userSearchLoading,
  selectedUser,
  searchUsers,
  selectUser,
  clearSelectedUser,
  cleanup: cleanupUserSearch,
} = useAdminUserSearch({
  minLength: USER_SEARCH_MIN_LENGTH,
  resultLimit: USER_SEARCH_RESULT_LIMIT,
})

const {
  items: schoolAdmins,
  totalPages,
  currentPage,
  totalItems,
  isLoading,
  errorMessage,
  isModalOpen,
  isSubmitting,
  modalError,
  formData,
  isConfirmOpen,
  itemToDelete,
  fetchItems: fetchSchoolAdmins,
  openAddModal: openCrudAddModal,
  closeDeleteConfirm,
  confirmDelete,
  handleSave: handleCrudSave,
  handleDelete,
} = useAdminCrudList({
  pageSize: PAGE_SIZE,
  fetchPage: ({ page, pageSize }) =>
    adminService.getSchoolAdmins({
      limit: pageSize,
      offset: (page - 1) * pageSize,
    }),
  createEmptyForm: () => ({ school_id: '' }),
  toEditForm: (admin) => ({ school_id: admin.school_id || '' }),
  validateForm: (form) => {
    if (!selectedUser.value || !form.school_id) {
      return 'Vui lòng chọn user và trường học'
    }

    return ''
  },
  createItem: async (form) => {
    const userId = selectedUser.value?.user_id
    if (!userId) {
      return
    }

    const roles = Array.isArray(selectedUser.value.roles) ? selectedUser.value.roles : []
    if (!roles.includes('SCHOOL_ADMIN')) {
      await adminService.assignRole(userId, 'SCHOOL_ADMIN')
    }

    await adminService.createSchoolAdmin({
      user_id: userId,
      school_id: form.school_id,
    })
  },
  deleteItem: (admin) => adminService.deleteSchoolAdmin(admin.admin_id),
  saveErrorMessage: 'Không thể gán school admin',
  deleteErrorPrefix: 'Không thể gỡ school admin',
  onAfterSave: async ({ fetchItems }) => {
    await fetchItems(1)
    return true
  },
  onAfterDelete: async ({ currentPage, currentItemsCount, fetchItems }) => {
    const nextPage = currentItemsCount === 1 && currentPage > 1 ? currentPage - 1 : currentPage

    await fetchItems(nextPage)
    return true
  },
})

const schoolNameById = computed(() => {
  const entries = schools.value
    .filter((school) => school?.school_id && school?.name)
    .map((school) => [school.school_id, school.name])

  return new Map(entries)
})

function getAdminKey(admin) {
  return admin?.admin_id || `${admin?.user_id || 'user'}-${admin?.school_id || 'school'}`
}

function getAdminEmail(admin) {
  return admin?.email || admin?.user?.email || 'N/A'
}

function getAdminSchoolName(admin) {
  return (
    admin?.school_name || admin?.school?.name || schoolNameById.value.get(admin?.school_id) || 'N/A'
  )
}

const deleteMessage = computed(() => {
  if (!itemToDelete.value) {
    return ''
  }

  return `Bạn có chắc muốn gỡ quyền quản trị trường "${getAdminSchoolName(itemToDelete.value)}" của "${getAdminEmail(itemToDelete.value)}"?`
})

function resetForm() {
  formData.value = { school_id: '' }
  modalError.value = ''
  userSearchQuery.value = ''
  clearSelectedUser()
}

function openAddModal() {
  resetForm()
  openCrudAddModal()
}

function closeModal() {
  isModalOpen.value = false
  resetForm()
}

async function fetchSchools() {
  try {
    const data = await adminService.getSchools({ limit: ADMIN_SELECTOR_FETCH_LIMIT, offset: 0 })
    schools.value = normalizeListResponse(data)
  } catch (error) {
    console.error('Không thể lấy danh sách trường:', error)
  }
}

watch(userSearchQuery, (value) => {
  if (
    selectedUser.value &&
    value ===
      (selectedUser.value.email || selectedUser.value.full_name || selectedUser.value.user_id)
  ) {
    return
  }

  const query = value.trim()
  if (query.length >= USER_SEARCH_MIN_LENGTH) {
    searchUsers(query)
  }
})

async function handleSave() {
  await handleCrudSave()

  if (!isModalOpen.value) {
    resetForm()
  }
}

onMounted(() => {
  fetchSchoolAdmins()
  fetchSchools()
})

onBeforeUnmount(() => {
  cleanupUserSearch()
})
</script>

<template>
  <div class="page-stack">
    <div class="page-actions">
      <button class="btn btn--primary" type="button" @click="openAddModal">
        + Gán Quản trị viên
      </button>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchSchoolAdmins(currentPage)">
        {{ ADMIN_RETRY_BUTTON_TEXT }}
      </button>
    </div>

    <LoadingSpinner v-else-if="isLoading" :message="ADMIN_LOADING_MESSAGE" />

    <div v-else-if="schoolAdmins.length === 0" class="card">
      <EmptyState
        title="Chưa có quản trị viên trường"
        message="Hãy gán người dùng làm quản trị viên cho một trường."
        icon="users"
      />
    </div>

    <template v-else>
      <div class="card desktop-table">
        <div class="table-responsive">
          <table class="table">
            <thead>
              <tr>
                <th>Người dùng</th>
                <th>Quản lý trường học</th>
                <th class="text-right">Thao tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="admin in schoolAdmins" :key="getAdminKey(admin)">
                <td class="font-medium">
                  <div>{{ getAdminEmail(admin) }}</div>
                </td>
                <td>
                  <span class="badge badge--primary">{{ getAdminSchoolName(admin) }}</span>
                </td>
                <td class="text-right">
                  <button
                    class="btn btn--sm btn--danger"
                    type="button"
                    @click="confirmDelete(admin)"
                  >
                    Gỡ bỏ
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="mobile-list">
        <article
          v-for="admin in schoolAdmins"
          :key="getAdminKey(admin)"
          class="card school-admin-card"
        >
          <div class="school-admin-card__row">
            <div>
              <p class="school-admin-card__label">Người dùng</p>
              <p class="school-admin-card__value">{{ getAdminEmail(admin) }}</p>
            </div>
            <button class="btn btn--sm btn--danger" type="button" @click="confirmDelete(admin)">
              Gỡ bỏ
            </button>
          </div>

          <div>
            <p class="school-admin-card__label">Trường học</p>
            <span class="badge badge--primary">{{ getAdminSchoolName(admin) }}</span>
          </div>
        </article>
      </div>

      <PaginationBar
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :limit="PAGE_SIZE"
        @page-change="fetchSchoolAdmins"
      />
    </template>

    <ActionModal :is-open="isModalOpen" title="Gán quản trị viên cho trường" @close="closeModal">
      <form class="modal-form" @submit.prevent="handleSave">
        <div v-if="modalError" class="alert alert--error">
          {{ modalError }}
        </div>

        <div class="alert alert--info">
          Hệ thống sẽ tự động cấp quyền <strong>SCHOOL_ADMIN</strong> nếu user được chọn chưa có vai
          trò này.
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="schoolAdminUserSearch">Người dùng</label>

          <div v-if="selectedUser" class="selected-user">
            <div class="selected-user__row">
              <div>
                <p class="selected-user__email">{{ selectedUser.email || selectedUser.user_id }}</p>
                <p class="selected-user__meta">{{ selectedUser.user_id }}</p>
              </div>
              <button
                class="btn btn--outline btn--sm"
                type="button"
                :disabled="isSubmitting"
                @click="clearSelectedUser"
              >
                Đổi user
              </button>
            </div>
          </div>

          <div v-else class="user-search">
            <input
              id="schoolAdminUserSearch"
              v-model="userSearchQuery"
              type="search"
              class="form-input"
              placeholder="Tìm theo email..."
              autocomplete="off"
              :disabled="isSubmitting"
            />

            <p class="search-helper">Nhập ít nhất 2 ký tự để tìm user.</p>

            <div v-if="userSearchLoading" class="search-feedback">Đang tìm user...</div>

            <div
              v-else-if="
                userSearchQuery.trim().length >= USER_SEARCH_MIN_LENGTH &&
                userSearchResults.length > 0
              "
              class="search-results"
            >
              <button
                v-for="user in userSearchResults"
                :key="user.user_id"
                class="search-result"
                type="button"
                @click="selectUser(user)"
              >
                <span class="search-result__email">{{ user.email || user.user_id }}</span>
                <span class="search-result__meta">{{ user.user_id }}</span>
              </button>
            </div>

            <div
              v-else-if="
                userSearchQuery.trim().length >= USER_SEARCH_MIN_LENGTH && !userSearchLoading
              "
              class="search-feedback"
            >
              Không tìm thấy user nào.
            </div>
          </div>
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="schoolSelect">Trường học</label>
          <select
            id="schoolSelect"
            v-model="formData.school_id"
            class="form-input"
            :disabled="isSubmitting || schools.length === 0"
            required
          >
            <option value="" disabled>
              {{ schools.length === 0 ? 'Đang tải danh sách trường...' : 'Chọn trường học' }}
            </option>
            <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
              {{ school.name }}
            </option>
          </select>
        </div>

        <div class="modal-actions">
          <button
            type="button"
            class="btn btn--outline"
            :disabled="isSubmitting"
            @click="closeModal"
          >
            Hủy
          </button>
          <button type="submit" class="btn btn--primary" :disabled="isSubmitting">
            {{ isSubmitting ? 'Đang xử lý...' : 'Gán quyền' }}
          </button>
        </div>
      </form>
    </ActionModal>

    <ConfirmDialog
      :is-open="isConfirmOpen"
      title="Gỡ quyền quản trị"
      :message="deleteMessage"
      confirm-text="Gỡ bỏ"
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

.desktop-table,
.school-admin-card {
  padding: var(--spacing-4);
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
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
}

.table tbody tr:hover {
  background-color: var(--color-background);
}

.text-right {
  text-align: right;
}

.school-admin-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.school-admin-card__row,
.selected-user__row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.school-admin-card__label,
.selected-user__meta,
.search-helper,
.search-feedback,
.search-result__meta {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.school-admin-card__label {
  margin-bottom: var(--spacing-1);
}

.school-admin-card__value,
.selected-user__email,
.search-result__email {
  margin: 0;
  font-weight: 700;
  color: var(--color-text);
}

.selected-user,
.search-results,
.search-feedback {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
}

.selected-user,
.search-feedback {
  padding: 0.875rem;
}

.search-helper {
  margin-top: var(--spacing-2);
}

.search-results {
  margin-top: var(--spacing-2);
  overflow: hidden;
}

.search-result {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  padding: 0.875rem;
  border: 0;
  border-bottom: 1px solid var(--color-border);
  background: transparent;
  text-align: left;
  cursor: pointer;
}

.search-result:last-child {
  border-bottom: 0;
}

.search-result:hover {
  background: var(--color-background);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.mobile-list {
  display: none;
}

@media (max-width: 767px) {
  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: flex;
  }

  .school-admin-card__row,
  .selected-user__row {
    flex-direction: column;
  }
}
</style>
