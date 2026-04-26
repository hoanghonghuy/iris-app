<script setup>
import { Link2, Pencil } from 'lucide-vue-next'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import AdminPeopleList from '../../components/admin/AdminPeopleList.vue'
import AdminPeopleEditModal from '../../components/admin/AdminPeopleEditModal.vue'
import AdminPeopleAssignModal from '../../components/admin/AdminPeopleAssignModal.vue'
import { adminService } from '../../services/adminService'
import { createAdminPersonEditFormConfig } from '../../helpers/adminPeopleFormConfig'
import { useAdminPeopleManagement } from '../../composables/admin'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_LOADING_MESSAGE,
  ADMIN_RETRY_BUTTON_TEXT,
} from '../../helpers/adminConfig'

const PAGE_SIZE = 20
const parentEditFormConfig = createAdminPersonEditFormConfig({ idField: 'parent_id' })

const {
  items: parents,
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
  students,
  selectedSchoolId,
  selectedClassId,
  selectedStudentId,
  filteredItems: filteredParents,
  updateSearchQuery,
  updateSelectedSchoolId,
  updateSelectedClassId,
  updateSelectedStudentId,
  updateEditForm,
  fetchItems: fetchParents,
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
  fetchList: adminService.getParents,
  ...parentEditFormConfig,
  updateItem: (form) =>
    adminService.updateParent(form.parent_id, {
      full_name: form.full_name.trim(),
      phone: form.phone?.trim() || '',
      school_id: form.school_id,
    }),
  updateErrorMessage: 'Không thể cập nhật phụ huynh',
  assignItem: ({ target, selectedStudentId: studentId }) =>
    adminService.assignParentToStudent(target.parent_id, studentId),
  enableStudentSelector: true,
  toUnassignTarget: (parent, student) => ({
    parent_id: parent.parent_id,
    student_id: student.student_id,
    student_name: student.full_name,
  }),
  unassignItem: (target) =>
    adminService.unassignParentFromStudent(target.parent_id, target.student_id),
})
</script>

<template>
  <div class="admin-parents">
    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchParents(currentPage)">
        {{ ADMIN_RETRY_BUTTON_TEXT }}
      </button>
    </div>

    <LoadingSpinner v-else-if="isLoading" :message="ADMIN_LOADING_MESSAGE" />

    <div v-else>
      <AdminPeopleList
        :items="parents"
        :filtered-items="filteredParents"
        :search-query="searchQuery"
        search-placeholder="Tìm theo tên, email, SĐT..."
        empty-title="Chưa có phụ huynh nào"
        empty-message="Chưa có phụ huynh nào đăng ký tài khoản trên hệ thống."
        empty-search-label="phụ huynh"
        item-key-field="parent_id"
        relation-field="children"
        relation-key-field="student_id"
        relation-name-field="full_name"
        relation-column-title="Học sinh quản lý"
        action-column-title="Gán học sinh"
        :action-column-width="220"
        no-relation-text="Chưa ghép học sinh"
        remove-relation-title="Hủy gán học sinh"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :limit="PAGE_SIZE"
        @update:search-query="updateSearchQuery"
        @remove-relation="({ item, relation }) => openUnassignDialog(item, relation)"
        @page-change="fetchParents"
      >
        <template #desktop-actions="{ item }">
          <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(item)">
            <Pencil :size="14" />
            <span>Sửa</span>
          </button>
          <button class="btn btn--sm btn--outline" type="button" @click="openAssignModal(item)">
            <Link2 :size="14" />
            <span>Gán học sinh</span>
          </button>
        </template>

        <template #mobile-head-extra="{ item }">
          <button class="btn btn--sm btn--outline" type="button" @click="openAssignModal(item)">
            <Link2 :size="14" />
            <span>Gán học sinh</span>
          </button>
        </template>

        <template #mobile-actions="{ item }">
          <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(item)">
            Sửa
          </button>
        </template>
      </AdminPeopleList>
    </div>

    <AdminPeopleEditModal
      :is-open="isEditModalOpen"
      :title="`Sửa thông tin phụ huynh - ${editForm.full_name || ''}`"
      :form-data="editForm"
      :error-message="editError"
      :is-loading="editLoading"
      :is-super-admin="isSuperAdmin"
      :schools="schools"
      full-name-input-id="parentFullName"
      phone-input-id="parentPhone"
      school-input-id="parentSchool"
      full-name-placeholder="Ví dụ: Trần Thị B"
      submit-text="Lưu lại"
      submitting-text="Đang lưu..."
      @close="closeEditModal"
      @submit="handleEdit"
      @update:form-data="updateEditForm"
    />

    <AdminPeopleAssignModal
      :is-open="isAssignModalOpen"
      :title="`Gán học sinh - ${assignTarget?.full_name || ''}`"
      :error-message="assignError"
      :is-loading="assignLoading"
      :is-super-admin="isSuperAdmin"
      :schools="schools"
      :classes="classes"
      :students="students"
      :selected-school-id="selectedSchoolId"
      :selected-class-id="selectedClassId"
      :selected-student-id="selectedStudentId"
      include-student-selector
      class-label="Chọn lớp"
      class-empty-text="Không có lớp"
      student-label="Chọn học sinh"
      student-empty-text="Không có học sinh"
      submit-text="Gán học sinh"
      submitting-text="Đang gán..."
      @close="closeAssignModal"
      @submit="handleAssign"
      @update:selected-school-id="updateSelectedSchoolId"
      @update:selected-class-id="updateSelectedClassId"
      @update:selected-student-id="updateSelectedStudentId"
    />

    <ConfirmDialog
      :is-open="isUnassignOpen"
      title="Xác nhận hủy gán"
      :message="`Bạn có chắc muốn hủy gán học sinh '${unassignTarget?.student_name || ''}' khỏi phụ huynh này?`"
      confirm-text="Hủy gán"
      is-danger
      :is-loading="unassignLoading"
      @confirm="handleUnassign"
      @cancel="closeUnassignDialog"
    />
  </div>
</template>

<style scoped>
.admin-parents {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.mt-2 {
  margin-top: var(--spacing-2);
}

</style>
