<script setup>
import { onBeforeUnmount, ref } from 'vue'
import { Link2, Pencil, Trash2, X } from 'lucide-vue-next'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import {
  createAdminPersonEditFormConfig,
  createAdminPersonRelationConfig,
} from '../../helpers/adminPeopleFormConfig'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'
import AdminPeopleList from '../../components/admin/AdminPeopleList.vue'
import AdminPeopleEditModal from '../../components/admin/AdminPeopleEditModal.vue'
import AdminPeopleAssignModal from '../../components/admin/AdminPeopleAssignModal.vue'
import { useAdminPeopleManagement, useAdminUserSearch } from '../../composables/admin'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_LOADING_MESSAGE,
  ADMIN_RETRY_BUTTON_TEXT,
} from '../../helpers/adminConfig'

const PAGE_SIZE = 20
const USER_SEARCH_MIN_LENGTH = 2
const USER_SEARCH_RESULT_LIMIT = 6

const teacherEditFormConfig = createAdminPersonEditFormConfig({ idField: 'teacher_id' })
const teacherRelationConfig = createAdminPersonRelationConfig({
  ownerIdField: 'teacher_id',
  relationIdField: 'class_id',
  relationNameField: 'name',
  unassignNameField: 'class_name',
  assignSelectionField: 'selectedClassId',
  assignService: adminService.assignTeacherToClass,
  unassignService: adminService.unassignTeacherFromClass,
})

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
  items: teachers,
  totalPages,
  currentPage,
  totalItems,
  isLoading,
  errorMessage,
  searchQuery,
  isSuperAdmin,
  isEditModalOpen,
  editLoading,
  editError,
  editForm,
  isAssignModalOpen,
  assignTarget,
  assignLoading,
  assignError,
  isUnassignOpen,
  unassignTarget,
  unassignLoading,
  schools,
  classes,
  selectedSchoolId,
  selectedClassId,
  filteredItems: filteredTeachers,
  updateSearchQuery,
  updateSelectedSchoolId,
  updateSelectedClassId,
  updateEditForm,
  fetchItems: fetchTeachers,
  openAssignModal,
  closeAssignModal,
  openEditModal,
  closeEditModal,
  handleEdit,
  handleAssign,
  openUnassignDialog,
  closeUnassignDialog,
  handleUnassign,
} = useAdminPeopleManagement({
  pageSize: PAGE_SIZE,
  searchFields: ['full_name', 'email', 'phone'],
  fetchList: adminService.getTeachers,
  ...teacherEditFormConfig,
  updateItem: (form) =>
    adminService.updateTeacher(form.teacher_id, {
      full_name: form.full_name.trim(),
      phone: form.phone?.trim() || '',
      school_id: form.school_id,
    }),
  updateErrorMessage: 'Không thể cập nhật giáo viên',
  ...teacherRelationConfig,
})

const isDeleteOpen = ref(false)
const deleteTarget = ref(null)
const deleteLoading = ref(false)

const isCreateModalOpen = ref(false)
const createLoading = ref(false)
const createError = ref('')
const createForm = ref({ school_id: '' })

function openDeleteDialog(teacher) {
  deleteTarget.value = teacher
  isDeleteOpen.value = true
}

function closeDeleteDialog() {
  isDeleteOpen.value = false
  deleteTarget.value = null
}

async function handleDelete() {
  if (!deleteTarget.value) {
    return
  }

  deleteLoading.value = true
  try {
    await adminService.deleteTeacher(deleteTarget.value.teacher_id)
    closeDeleteDialog()
    await fetchTeachers(currentPage.value)
  } catch (error) {
    errorMessage.value = `Lỗi xóa: ${extractErrorMessage(error)}`
    closeDeleteDialog()
  } finally {
    deleteLoading.value = false
  }
}

function openCreateModal() {
  createForm.value = { school_id: selectedSchoolId.value || '' }
  createError.value = ''
  userSearchQuery.value = ''
  clearSelectedUser()
  isCreateModalOpen.value = true
}

function closeCreateModal() {
  isCreateModalOpen.value = false
  createForm.value = { school_id: '' }
  createError.value = ''
  userSearchQuery.value = ''
  clearSelectedUser()
}

async function handleCreateTeacher() {
  if (!selectedUser.value || !createForm.value.school_id) {
    createError.value = 'Vui lòng chọn user và trường học'
    return
  }

  createLoading.value = true
  createError.value = ''

  try {
    const roles = Array.isArray(selectedUser.value.roles) ? selectedUser.value.roles : []
    if (!roles.includes('TEACHER')) {
      await adminService.assignRole(selectedUser.value.user_id, 'TEACHER')
    }

    await adminService.createTeacher({
      user_id: selectedUser.value.user_id,
      school_id: createForm.value.school_id,
    })

    closeCreateModal()
    await fetchTeachers(1)
  } catch (error) {
    createError.value = extractErrorMessage(error) || 'Không thể gán giáo viên'
  } finally {
    createLoading.value = false
  }
}

onBeforeUnmount(() => {
  cleanupUserSearch()
})
</script>

<template>
  <div class="admin-teachers">
    <div class="page-actions">
      <button class="btn btn--primary" type="button" @click="openCreateModal">
        + Gán giáo viên
      </button>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchTeachers(currentPage)">
        {{ ADMIN_RETRY_BUTTON_TEXT }}
      </button>
    </div>

    <LoadingSpinner v-else-if="isLoading" :message="ADMIN_LOADING_MESSAGE" />

    <div v-else>
      <AdminPeopleList
        :items="teachers"
        :filtered-items="filteredTeachers"
        :search-query="searchQuery"
        search-placeholder="Tìm theo tên, email, SĐT..."
        empty-title="Chưa có giáo viên nào"
        empty-message="Hãy tạo tài khoản user và cấp quyền TEACHER."
        empty-search-label="giáo viên"
        item-key-field="teacher_id"
        relation-field="classes"
        relation-key-field="class_id"
        relation-name-field="name"
        relation-column-title="Lớp phụ trách"
        action-column-title="Hành động"
        :action-column-width="260"
        no-relation-text="Chưa phân lớp"
        remove-relation-title="Hủy gán lớp"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :limit="PAGE_SIZE"
        @update:search-query="updateSearchQuery"
        @remove-relation="({ item, relation }) => openUnassignDialog(item, relation)"
        @page-change="fetchTeachers"
      >
        <template #desktop-actions="{ item }">
          <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(item)">
            <Pencil :size="14" />
            <span>Sửa</span>
          </button>
          <button class="btn btn--sm btn--outline" type="button" @click="openAssignModal(item)">
            <Link2 :size="14" />
            <span>Gán lớp</span>
          </button>
          <button class="btn btn--sm btn--danger" type="button" @click="openDeleteDialog(item)">
            <Trash2 :size="14" />
            <span>Xóa</span>
          </button>
        </template>

        <template #mobile-actions="{ item }">
          <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(item)">
            Sửa
          </button>
          <button class="btn btn--sm btn--outline" type="button" @click="openAssignModal(item)">
            <Link2 :size="14" />
            <span>Gán lớp</span>
          </button>
          <button class="btn btn--sm btn--danger" type="button" @click="openDeleteDialog(item)">
            Xóa
          </button>
        </template>
      </AdminPeopleList>
    </div>

    <AdminPeopleEditModal
      :is-open="isEditModalOpen"
      :title="`Sửa thông tin giáo viên - ${editForm.full_name || ''}`"
      :form-data="editForm"
      :error-message="editError"
      :is-loading="editLoading"
      :is-super-admin="isSuperAdmin"
      :schools="schools"
      full-name-input-id="teacherFullName"
      phone-input-id="teacherPhone"
      school-input-id="teacherSchool"
      full-name-placeholder="Ví dụ: Nguyễn Văn A"
      submit-text="Lưu lại"
      submitting-text="Đang lưu..."
      @close="closeEditModal"
      @submit="handleEdit"
      @update:form-data="updateEditForm"
    />

    <AdminPeopleAssignModal
      :is-open="isAssignModalOpen"
      :title="`Gán lớp phụ trách - ${assignTarget?.full_name || ''}`"
      :error-message="assignError"
      :is-loading="assignLoading"
      :is-super-admin="isSuperAdmin"
      :schools="schools"
      :classes="classes"
      :selected-school-id="selectedSchoolId"
      :selected-class-id="selectedClassId"
      class-label="Chọn lớp"
      class-empty-text="Không có lớp"
      submit-text="Gán lớp"
      submitting-text="Đang gán..."
      @close="closeAssignModal"
      @submit="handleAssign"
      @update:selected-school-id="updateSelectedSchoolId"
      @update:selected-class-id="updateSelectedClassId"
    />

    <ActionModal :is-open="isCreateModalOpen" title="Gán giáo viên" @close="closeCreateModal">
      <div class="modal-form">
        <div v-if="createError" class="alert alert--error">{{ createError }}</div>

        <div class="form-group">
          <label class="form-label">Tìm user</label>
          <input
            v-model="userSearchQuery"
            type="text"
            class="form-input"
            placeholder="Nhập email để tìm user..."
            :disabled="createLoading"
            @input="searchUsers(userSearchQuery)"
          />
          <div v-if="userSearchLoading" class="text-muted text-sm mt-1">Đang tìm...</div>
          <div v-if="userSearchResults.length > 0" class="user-search-results">
            <button
              v-for="user in userSearchResults"
              :key="user.user_id"
              type="button"
              class="user-search-item"
              @click="selectUser(user)"
            >
              <span class="font-medium">{{ user.email }}</span>
              <span class="text-xs text-muted">{{ user.roles?.join(', ') || 'No roles' }}</span>
            </button>
          </div>
          <div v-if="selectedUser" class="selected-user-badge">
            <span>{{ selectedUser.email }}</span>
            <button type="button" @click="clearSelectedUser">
              <X :size="14" />
            </button>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label" for="createSchoolId">Trường học</label>
          <select
            id="createSchoolId"
            v-model="createForm.school_id"
            class="form-input"
            :disabled="createLoading"
            required
          >
            <option value="">-- Chọn trường --</option>
            <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
              {{ school.name }}
            </option>
          </select>
        </div>

        <div class="modal-actions">
          <button
            class="btn btn--outline"
            type="button"
            :disabled="createLoading"
            @click="closeCreateModal"
          >
            Hủy
          </button>
          <button
            class="btn btn--primary"
            type="button"
            :disabled="createLoading || !selectedUser || !createForm.school_id"
            @click="handleCreateTeacher"
          >
            {{ createLoading ? 'Đang gán...' : 'Gán giáo viên' }}
          </button>
        </div>
      </div>
    </ActionModal>

    <ConfirmDialog
      :is-open="isUnassignOpen"
      title="Xác nhận hủy gán"
      :message="`Bạn có chắc muốn hủy gán lớp '${unassignTarget?.class_name || ''}' khỏi giáo viên này?`"
      confirm-text="Hủy gán"
      is-danger
      :is-loading="unassignLoading"
      @confirm="handleUnassign"
      @cancel="closeUnassignDialog"
    />

    <ConfirmDialog
      :is-open="isDeleteOpen"
      title="Xác nhận xóa giáo viên"
      :message="`Bạn có chắc chắn muốn xóa giáo viên '${deleteTarget?.full_name || ''}' không? Hành động này không thể hoàn tác.`"
      confirm-text="Xóa giáo viên"
      is-danger
      :is-loading="deleteLoading"
      @confirm="handleDelete"
      @cancel="closeDeleteDialog"
    />
  </div>
</template>

<style scoped>
.admin-teachers,
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-actions {
  display: flex;
  justify-content: flex-end;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.user-search-results {
  margin-top: var(--spacing-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  max-height: 200px;
  overflow-y: auto;
}

.user-search-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  padding: var(--spacing-3);
  width: 100%;
  text-align: left;
  background: none;
  border: none;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: background-color 0.2s;
}

.user-search-item:last-child {
  border-bottom: none;
}

.user-search-item:hover {
  background-color: var(--color-background);
}

.selected-user-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  margin-top: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  background-color: var(--color-primary-light);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
}

.selected-user-badge button {
  display: inline-flex;
  align-items: center;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-muted);
  padding: 0;
}

.selected-user-badge button:hover {
  color: var(--color-danger);
}

.mt-1 {
  margin-top: var(--spacing-1);
}

.mt-2 {
  margin-top: var(--spacing-2);
}
</style>
