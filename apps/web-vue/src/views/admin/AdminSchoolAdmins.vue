<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'
import PaginationBar from '../../components/PaginationBar.vue'
import ConfirmDialog from '../../components/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'

const PAGE_SIZE = 10
const USER_SEARCH_LIMIT = 100
const USER_SEARCH_MIN_LENGTH = 2
const USER_SEARCH_RESULT_LIMIT = 6

const schoolAdmins = ref([])
const schools = ref([])
const totalPages = ref(0)
const currentPage = ref(1)
const totalItems = ref(0)
const isLoading = ref(true)
const errorMessage = ref('')

const isModalOpen = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const formData = ref({ school_id: '' })

const userSearchQuery = ref('')
const userSearchResults = ref([])
const userSearchLoading = ref(false)
const selectedUser = ref(null)

const isConfirmOpen = ref(false)
const itemToDelete = ref(null)

let userSearchTimerId = null
let userSearchRequestId = 0

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
    admin?.school_name ||
    admin?.school?.name ||
    schoolNameById.value.get(admin?.school_id) ||
    'N/A'
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
  userSearchResults.value = []
  userSearchLoading.value = false
  selectedUser.value = null
}

function openAddModal() {
  resetForm()
  isModalOpen.value = true
}

function closeModal() {
  isModalOpen.value = false
  resetForm()
}

async function fetchSchools() {
  try {
    const data = await adminService.getSchools({ limit: 100, offset: 0 })
    schools.value = Array.isArray(data?.data) ? data.data : []
  } catch (error) {
    console.error('Không thể lấy danh sách trường:', error)
  }
}

async function fetchSchoolAdmins(page = 1) {
  isLoading.value = true
  errorMessage.value = ''
  currentPage.value = page

  try {
    const data = await adminService.getSchoolAdmins({
      limit: PAGE_SIZE,
      offset: (page - 1) * PAGE_SIZE,
    })

    schoolAdmins.value = Array.isArray(data?.data) ? data.data : []

    if (data?.pagination) {
      const limit = data.pagination.limit || PAGE_SIZE
      totalItems.value = data.pagination.total || 0
      totalPages.value = Math.ceil(totalItems.value / limit) || 1
    } else {
      totalItems.value = schoolAdmins.value.length
      totalPages.value = schoolAdmins.value.length > 0 ? 1 : 0
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách school admin'
  } finally {
    isLoading.value = false
  }
}

async function searchUsers(query) {
  const normalizedQuery = query.trim().toLowerCase()
  const requestId = ++userSearchRequestId

  if (normalizedQuery.length < USER_SEARCH_MIN_LENGTH) {
    userSearchResults.value = []
    userSearchLoading.value = false
    return
  }

  userSearchLoading.value = true

  try {
    const matches = []
    let offset = 0
    let hasMore = true

    while (hasMore && matches.length < USER_SEARCH_RESULT_LIMIT) {
      const response = await adminService.getUsers({
        limit: USER_SEARCH_LIMIT,
        offset,
      })
      const items = Array.isArray(response?.data) ? response.data : []

      items.forEach((user) => {
        if (
          user?.email?.toLowerCase().includes(normalizedQuery) &&
          !matches.some((item) => item.user_id === user.user_id)
        ) {
          matches.push(user)
        }
      })

      const pagination = response?.pagination
      hasMore = Boolean(pagination?.has_more) && items.length > 0
      offset += pagination?.limit || USER_SEARCH_LIMIT
    }

    if (requestId !== userSearchRequestId) {
      return
    }

    userSearchResults.value = matches.slice(0, USER_SEARCH_RESULT_LIMIT)
  } catch (error) {
    if (requestId !== userSearchRequestId) {
      return
    }

    modalError.value = extractErrorMessage(error) || 'Không thể tìm user'
    userSearchResults.value = []
  } finally {
    if (requestId === userSearchRequestId) {
      userSearchLoading.value = false
    }
  }
}

function selectUser(user) {
  selectedUser.value = user
  userSearchQuery.value = user.email || ''
  userSearchResults.value = []
}

function clearSelectedUser() {
  selectedUser.value = null
  userSearchQuery.value = ''
  userSearchResults.value = []
  modalError.value = ''
}

watch(userSearchQuery, (value) => {
  if (selectedUser.value && value === selectedUser.value.email) {
    userSearchResults.value = []
    userSearchLoading.value = false
    return
  }

  selectedUser.value = null
  modalError.value = ''

  if (userSearchTimerId) {
    clearTimeout(userSearchTimerId)
  }

  const query = value.trim()
  if (query.length < USER_SEARCH_MIN_LENGTH) {
    userSearchResults.value = []
    userSearchLoading.value = false
    return
  }

  userSearchTimerId = setTimeout(() => {
    searchUsers(query)
  }, 250)
})

async function handleSave() {
  if (!selectedUser.value || !formData.value.school_id) {
    modalError.value = 'Vui lòng chọn user và trường học'
    return
  }

  isSubmitting.value = true
  modalError.value = ''

  try {
    const roles = Array.isArray(selectedUser.value.roles) ? selectedUser.value.roles : []
    if (!roles.includes('SCHOOL_ADMIN')) {
      await adminService.assignRole(selectedUser.value.user_id, 'SCHOOL_ADMIN')
    }

    await adminService.createSchoolAdmin({
      user_id: selectedUser.value.user_id,
      school_id: formData.value.school_id,
    })

    closeModal()
    fetchSchoolAdmins(1)
  } catch (error) {
    modalError.value = extractErrorMessage(error) || 'Không thể gán school admin'
  } finally {
    isSubmitting.value = false
  }
}

function confirmDelete(admin) {
  itemToDelete.value = admin
  isConfirmOpen.value = true
}

async function handleDelete() {
  if (!itemToDelete.value) {
    return
  }

  isSubmitting.value = true

  try {
    await adminService.deleteSchoolAdmin(itemToDelete.value.admin_id)
    isConfirmOpen.value = false
    itemToDelete.value = null

    const nextPage = schoolAdmins.value.length === 1 && currentPage.value > 1
      ? currentPage.value - 1
      : currentPage.value

    fetchSchoolAdmins(nextPage)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể gỡ school admin'
    isConfirmOpen.value = false
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  fetchSchoolAdmins()
  fetchSchools()
})

onBeforeUnmount(() => {
  if (userSearchTimerId) {
    clearTimeout(userSearchTimerId)
  }
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
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchSchoolAdmins(currentPage)">
        Thử lại
      </button>
    </div>

    <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu..." />

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
                  <button class="btn btn--sm btn--danger" type="button" @click="confirmDelete(admin)">
                    Gỡ bỏ
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="mobile-list">
        <article v-for="admin in schoolAdmins" :key="getAdminKey(admin)" class="card school-admin-card">
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

    <ActionModal
      :is-open="isModalOpen"
      title="Gán quản trị viên cho trường"
      @close="closeModal"
    >
      <form class="modal-form" @submit.prevent="handleSave">
        <div v-if="modalError" class="alert alert--error">
          {{ modalError }}
        </div>

        <div class="alert alert--info">
          Hệ thống sẽ tự động cấp quyền <strong>SCHOOL_ADMIN</strong> nếu user được chọn chưa có vai trò này.
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="schoolAdminUserSearch">Người dùng</label>

          <div v-if="selectedUser" class="selected-user">
            <div class="selected-user__row">
              <div>
                <p class="selected-user__email">{{ selectedUser.email || selectedUser.user_id }}</p>
                <p class="selected-user__meta">{{ selectedUser.user_id }}</p>
              </div>
              <button class="btn btn--outline btn--sm" type="button" :disabled="isSubmitting" @click="clearSelectedUser">
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

            <div v-if="userSearchLoading" class="search-feedback">
              Đang tìm user...
            </div>

            <div
              v-else-if="userSearchQuery.trim().length >= USER_SEARCH_MIN_LENGTH && userSearchResults.length > 0"
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
              v-else-if="userSearchQuery.trim().length >= USER_SEARCH_MIN_LENGTH && !userSearchLoading"
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
          <button type="button" class="btn btn--outline" :disabled="isSubmitting" @click="closeModal">
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
      @cancel="isConfirmOpen = false"
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
