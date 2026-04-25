<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'
import PaginationBar from '../../components/PaginationBar.vue'
import ConfirmDialog from '../../components/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'

const authStore = useAuthStore()

const PAGE_SIZE = 20
const ROLE_OPTIONS = ['TEACHER', 'PARENT', 'SCHOOL_ADMIN', 'SUPER_ADMIN']
const CREATABLE_USER_ROLES = ['TEACHER', 'PARENT', 'SCHOOL_ADMIN']

const ROLE_LABELS = {
  SUPER_ADMIN: 'Super Admin',
  SCHOOL_ADMIN: 'School Admin',
  TEACHER: 'Giáo viên',
  PARENT: 'Phụ huynh',
}

const STATUS_LABELS = {
  active: 'Hoạt động',
  pending: 'Chờ kích hoạt',
  locked: 'Đã khóa',
}

const STATUS_BADGES = {
  active: 'badge badge--success',
  pending: 'badge badge--outline',
  locked: 'badge badge--danger',
}

const users = ref([])
const totalPages = ref(0)
const currentPage = ref(1)
const totalItems = ref(0)
const isLoading = ref(true)
const errorMessage = ref('')

const searchQuery = ref('')
const selectedRoleFilter = ref('ALL')
const actionLoadingUserId = ref('')

const isAddModalOpen = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const newUserData = ref({
  email: '',
  roles: ['TEACHER'],
})

const isLockConfirmOpen = ref(false)
const userToToggle = ref(null)

function normalizeUserStatus(user) {
  if (user?.status) {
    return user.status
  }

  if (user?.is_active === false) {
    return 'locked'
  }

  if (user?.is_active === true) {
    return 'active'
  }

  return 'pending'
}

function formatRoles(roles) {
  if (!Array.isArray(roles) || roles.length === 0) {
    return ['USER']
  }

  return roles
}

const filteredUsers = computed(() => {
  const search = searchQuery.value.trim().toLowerCase()
  if (!search) {
    return users.value
  }

  return users.value.filter((user) => user.email?.toLowerCase().includes(search))
})

const currentUserId = computed(() => authStore.currentUser?.user_id || '')

const emptyState = computed(() => {
  if (users.value.length === 0 && selectedRoleFilter.value === 'ALL') {
    return {
      title: 'Chưa có người dùng nào',
      message: 'Hiện tại hệ thống chưa có dữ liệu người dùng mới.',
      showCreateAction: true,
    }
  }

  if (users.value.length === 0 && selectedRoleFilter.value !== 'ALL') {
    return {
      title: 'Không tìm thấy người dùng',
      message: 'Không có dữ liệu phù hợp với vai trò đang lọc.',
      showCreateAction: false,
    }
  }

  if (users.value.length > 0 && filteredUsers.value.length === 0) {
    return {
      title: 'Không có kết quả tìm kiếm',
      message: `Không tìm thấy người dùng nào mang email "${searchQuery.value}".`,
      showCreateAction: false,
    }
  }

  return null
})

async function fetchUsers(page = 1) {
  isLoading.value = true
  errorMessage.value = ''
  currentPage.value = page

  try {
    const params = {
      limit: PAGE_SIZE,
      offset: (page - 1) * PAGE_SIZE,
    }

    if (selectedRoleFilter.value !== 'ALL') {
      params.role = selectedRoleFilter.value
    }

    const data = await adminService.getUsers(params)
    users.value = Array.isArray(data?.data) ? data.data : []

    if (data?.pagination) {
      totalItems.value = data.pagination.total || 0
      totalPages.value = Math.ceil(data.pagination.total / data.pagination.limit) || 1
    } else {
      totalItems.value = users.value.length
      totalPages.value = users.value.length > 0 ? 1 : 0
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách người dùng'
  } finally {
    isLoading.value = false
  }
}

function handleRoleFilterChange() {
  fetchUsers(1)
}

function openAddModal() {
  newUserData.value = {
    email: '',
    roles: ['TEACHER'],
  }
  modalError.value = ''
  isAddModalOpen.value = true
}

function toggleCreateRole(role) {
  const currentRoles = [...newUserData.value.roles]
  if (currentRoles.includes(role)) {
    newUserData.value.roles = currentRoles.filter((item) => item !== role)
    return
  }

  newUserData.value.roles = [...currentRoles, role]
}

async function handleAddUser() {
  if (!newUserData.value.email.trim()) {
    modalError.value = 'Vui lòng nhập email'
    return
  }

  if (newUserData.value.roles.length === 0) {
    modalError.value = 'Chọn ít nhất 1 vai trò'
    return
  }

  isSubmitting.value = true
  modalError.value = ''

  try {
    await adminService.createUser({
      email: newUserData.value.email.trim(),
      roles: newUserData.value.roles,
    })

    isAddModalOpen.value = false
    fetchUsers(1)
  } catch (error) {
    modalError.value = extractErrorMessage(error) || 'Không thể tạo user'
  } finally {
    isSubmitting.value = false
  }
}

function confirmToggleLock(user) {
  userToToggle.value = user
  isLockConfirmOpen.value = true
}

async function handleToggleLock() {
  if (!userToToggle.value) return

  const targetStatus = normalizeUserStatus(userToToggle.value)
  actionLoadingUserId.value = userToToggle.value.user_id
  isSubmitting.value = true

  try {
    if (targetStatus === 'locked') {
      await adminService.unlockUser(userToToggle.value.user_id)
    } else {
      await adminService.lockUser(userToToggle.value.user_id)
    }

    isLockConfirmOpen.value = false
    userToToggle.value = null
    fetchUsers(currentPage.value)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || `Không thể ${targetStatus === 'locked' ? 'mở khóa' : 'khóa'} tài khoản`
    isLockConfirmOpen.value = false
  } finally {
    actionLoadingUserId.value = ''
    isSubmitting.value = false
  }
}

onMounted(async () => {
  if (!authStore.currentUser && authStore.isAuthenticated) {
    await authStore.fetchCurrentUser()
  }
  fetchUsers()
})
</script>

<template>
  <div class="admin-users page-stack">
    <div class="page-actions">
      <button class="btn btn--primary" type="button" @click="openAddModal">
        + Tạo user
      </button>
    </div>

    <div
      v-if="!isLoading && !errorMessage && (users.length > 0 || selectedRoleFilter !== 'ALL') && !isAddModalOpen"
      class="card toolbar-card"
    >
      <div class="toolbar-grid">
        <input
          v-model="searchQuery"
          class="form-input"
          type="search"
          placeholder="Tìm theo email..."
        />

        <select v-model="selectedRoleFilter" class="form-input toolbar-select" @change="handleRoleFilterChange">
          <option value="ALL">Tất cả vai trò</option>
          <option v-for="role in ROLE_OPTIONS" :key="role" :value="role">
            {{ ROLE_LABELS[role] || role }}
          </option>
        </select>
      </div>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchUsers(currentPage)">Thử lại</button>
    </div>

    <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu..." />

    <div v-else-if="emptyState" class="card">
      <EmptyState :title="emptyState.title" :message="emptyState.message" icon="users">
        <template v-if="emptyState.showCreateAction" #action>
          <button class="btn btn--primary" type="button" @click="openAddModal">
            Tạo user đầu tiên
          </button>
        </template>
      </EmptyState>
    </div>

    <template v-else>
      <div class="card desktop-table">
        <div class="table-responsive">
          <table class="table">
            <thead>
              <tr>
                <th>Email</th>
                <th>Vai trò</th>
                <th>Trạng thái</th>
                <th class="text-right">Thao tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.user_id">
                <td class="font-medium">{{ user.email }}</td>
                <td>
                  <div class="role-list">
                    <span
                      v-for="role in formatRoles(user.roles)"
                      :key="role"
                      class="badge badge--outline"
                    >
                      {{ ROLE_LABELS[role] || role }}
                    </span>
                  </div>
                </td>
                <td>
                  <span :class="STATUS_BADGES[normalizeUserStatus(user)] || 'badge badge--outline'">
                    {{ STATUS_LABELS[normalizeUserStatus(user)] || normalizeUserStatus(user) }}
                  </span>
                </td>
                <td class="text-right">
                  <button
                    v-if="normalizeUserStatus(user) === 'active'"
                    class="btn btn--sm btn--danger"
                    type="button"
                    :disabled="actionLoadingUserId === user.user_id || currentUserId === user.user_id"
                    :title="currentUserId === user.user_id ? 'Bạn không thể tự khóa chính mình' : ''"
                    @click="confirmToggleLock(user)"
                  >
                    {{ actionLoadingUserId === user.user_id ? 'Đang xử lý...' : 'Khóa' }}
                  </button>

                  <button
                    v-else-if="normalizeUserStatus(user) === 'locked'"
                    class="btn btn--sm btn--success"
                    type="button"
                    :disabled="actionLoadingUserId === user.user_id"
                    @click="confirmToggleLock(user)"
                  >
                    {{ actionLoadingUserId === user.user_id ? 'Đang xử lý...' : 'Mở khóa' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="mobile-list">
        <article v-for="user in filteredUsers" :key="user.user_id" class="card user-card">
          <div class="user-card__head">
            <p class="user-email">{{ user.email }}</p>
            <span :class="STATUS_BADGES[normalizeUserStatus(user)] || 'badge badge--outline'">
              {{ STATUS_LABELS[normalizeUserStatus(user)] || normalizeUserStatus(user) }}
            </span>
          </div>

          <div class="role-list">
            <span
              v-for="role in formatRoles(user.roles)"
              :key="role"
              class="badge badge--outline"
            >
              {{ ROLE_LABELS[role] || role }}
            </span>
          </div>

          <div class="user-card__actions">
            <button
              v-if="normalizeUserStatus(user) === 'active'"
              class="btn btn--sm btn--danger"
              type="button"
              :disabled="actionLoadingUserId === user.user_id || currentUserId === user.user_id"
              :title="currentUserId === user.user_id ? 'Bạn không thể tự khóa chính mình' : ''"
              @click="confirmToggleLock(user)"
            >
              {{ actionLoadingUserId === user.user_id ? 'Đang xử lý...' : 'Khóa' }}
            </button>

            <button
              v-else-if="normalizeUserStatus(user) === 'locked'"
              class="btn btn--sm btn--success"
              type="button"
              :disabled="actionLoadingUserId === user.user_id"
              @click="confirmToggleLock(user)"
            >
              {{ actionLoadingUserId === user.user_id ? 'Đang xử lý...' : 'Mở khóa' }}
            </button>
          </div>
        </article>
      </div>

      <PaginationBar
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :limit="PAGE_SIZE"
        @page-change="fetchUsers"
      />
    </template>

    <ActionModal
      :is-open="isAddModalOpen"
      title="Tạo user mới"
      @close="isAddModalOpen = false"
    >
      <form class="modal-form" @submit.prevent="handleAddUser">
        <div v-if="modalError" class="alert alert--error">
          {{ modalError }}
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="userEmail">Email</label>
          <input
            id="userEmail"
            v-model="newUserData.email"
            type="email"
            class="form-input"
            placeholder="Nhập email"
            :disabled="isSubmitting"
            required
          />
        </div>

        <div class="form-group mb-0">
          <label class="form-label">Vai trò ban đầu</label>
          <div class="role-picker">
            <label
              v-for="role in CREATABLE_USER_ROLES"
              :key="role"
              class="role-option"
              :class="{ 'role-option--active': newUserData.roles.includes(role) }"
            >
              <input
                type="checkbox"
                class="role-option__checkbox"
                :checked="newUserData.roles.includes(role)"
                :disabled="isSubmitting"
                @change="toggleCreateRole(role)"
              />
              <span>{{ ROLE_LABELS[role] || role }}</span>
            </label>
          </div>
        </div>

        <div class="modal-actions">
          <button type="button" class="btn btn--outline" :disabled="isSubmitting" @click="isAddModalOpen = false">
            Hủy
          </button>
          <button type="submit" class="btn btn--primary" :disabled="isSubmitting">
            {{ isSubmitting ? 'Đang tạo...' : 'Tạo tài khoản' }}
          </button>
        </div>
      </form>
    </ActionModal>

    <ConfirmDialog
      :is-open="isLockConfirmOpen"
      :title="normalizeUserStatus(userToToggle) === 'locked' ? 'Xác nhận mở khóa tài khoản' : 'Xác nhận khóa tài khoản'"
      :message="normalizeUserStatus(userToToggle) === 'locked'
        ? `Tài khoản ${userToToggle?.email || ''} sẽ có thể đăng nhập lại bình thường.`
        : `Tài khoản ${userToToggle?.email || ''} sẽ không thể đăng nhập được nữa. Bạn có chắc chắn?`"
      :confirm-text="normalizeUserStatus(userToToggle) === 'locked' ? 'Mở khóa tài khoản' : 'Khóa tài khoản'"
      :is-danger="normalizeUserStatus(userToToggle) !== 'locked'"
      :is-loading="isSubmitting"
      @confirm="handleToggleLock"
      @cancel="isLockConfirmOpen = false"
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
.desktop-table,
.user-card {
  padding: var(--spacing-4);
}

.toolbar-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: minmax(0, 1fr);
}

.toolbar-select {
  max-width: 14rem;
}

.role-list,
.user-card__actions,
.modal-actions,
.role-picker {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.user-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.user-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.user-email,
.user-meta {
  margin: 0;
}

.user-email {
  font-weight: 700;
  color: var(--color-text);
}

.user-meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.role-picker {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
}

.role-option {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  min-height: 2.75rem;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  cursor: pointer;
}

.role-option--active {
  border-color: color-mix(in srgb, var(--color-primary) 30%, var(--color-border));
  background: color-mix(in srgb, var(--color-primary) 8%, transparent);
}

.role-option__checkbox {
  margin: 0;
}

.modal-actions {
  justify-content: flex-end;
}

.mobile-list {
  display: none;
}

@media (min-width: 768px) {
  .toolbar-grid {
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: center;
  }
}

@media (max-width: 767px) {
  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: flex;
  }

  .toolbar-select {
    max-width: none;
  }
}
</style>
