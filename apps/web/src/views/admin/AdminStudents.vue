<script setup>
import EmptyState from '../../components/common/EmptyState.vue'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import StudentsTable from './students/StudentsTable.vue'
import StudentFormModal from './students/StudentFormModal.vue'
import { useAdminStudentsPage } from './students/useAdminStudentsPage'
import {
  ADMIN_LOAD_ERROR_TITLE,
  ADMIN_LOADING_MESSAGE,
  ADMIN_RETRY_BUTTON_TEXT,
} from '../../helpers/adminConfig'
import {
  getDeleteStudentConfirmMessage,
  STUDENT_REVOKE_CONFIRM_MESSAGE,
} from './students/studentConfig'

const {
  schools,
  classes,
  selectedSchoolId,
  selectedClassId,
  students,
  searchQuery,
  isBootstrapping,
  isLoadingStudents,
  errorMessage,
  codeError,
  isFormModalOpen,
  isFormSubmitting,
  formError,
  formMode,
  formData,
  generatingCodeStudentId,
  revokingCodeStudentId,
  copiedStudentId,
  isRevokeConfirmOpen,
  isDeleteConfirmOpen,
  deleteTarget,
  isDeleteLoading,
  isSuperAdmin,
  selectedSchoolName,
  selectedClassName,
  filteredStudents,
  openAddModal,
  openEditModal,
  closeFormModal,
  updateFormData,
  submitForm,
  confirmDelete,
  closeDeleteConfirm,
  handleDelete,
  handleGenerateCode,
  confirmRevokeCode,
  closeRevokeConfirm,
  handleRevokeCode,
  handleCopyCode,
  loadStudents,
} = useAdminStudentsPage()
</script>

<template>
  <div class="admin-students page-stack">
    <div class="page-actions">
      <button
        class="btn btn--primary"
        type="button"
        :disabled="!selectedClassId"
        @click="openAddModal"
      >
        + Thêm học sinh
      </button>
    </div>

    <div class="card toolbar-card">
      <div class="toolbar-grid">
        <div v-if="isSuperAdmin" class="form-group mb-0">
          <label class="form-label" for="schoolFilter">Trường học</label>
          <select
            id="schoolFilter"
            v-model="selectedSchoolId"
            class="form-input"
            :disabled="isBootstrapping"
          >
            <option value="" disabled>-- Chọn trường --</option>
            <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
              {{ school.name }}
            </option>
          </select>
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="classFilter">Lớp học</label>
          <select
            id="classFilter"
            v-model="selectedClassId"
            class="form-input"
            :disabled="isBootstrapping || classes.length === 0"
          >
            <option value="" disabled>
              {{ classes.length === 0 ? '-- Không có lớp --' : '-- Chọn lớp --' }}
            </option>
            <option
              v-for="classItem in classes"
              :key="classItem.class_id"
              :value="classItem.class_id"
            >
              {{ classItem.name }}
            </option>
          </select>
        </div>

        <div v-if="students.length > 0" class="form-group mb-0 toolbar-search">
          <label class="form-label" for="studentSearch">Tìm kiếm</label>
          <input
            id="studentSearch"
            v-model="searchQuery"
            type="search"
            class="form-input"
            placeholder="Tìm theo tên học sinh..."
          />
        </div>
      </div>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="loadStudents">
        {{ ADMIN_RETRY_BUTTON_TEXT }}
      </button>
    </div>

    <div v-if="codeError" class="alert alert--error">
      {{ codeError }}
    </div>

    <LoadingSpinner v-if="isBootstrapping || isLoadingStudents" :message="ADMIN_LOADING_MESSAGE" />

    <div v-else-if="schools.length === 0" class="card">
      <EmptyState
        title="Chưa có trường học nào"
        message="Bạn cần tạo trường học trước khi có thể quản lý học sinh."
      />
    </div>

    <div v-else-if="classes.length === 0" class="card">
      <EmptyState
        title="Chưa có lớp học nào"
        :message="`${selectedSchoolName || 'Trường này'} chưa có lớp học nào. Hãy tạo lớp trước khi thêm học sinh.`"
      />
    </div>

    <div v-else-if="students.length === 0" class="card">
      <EmptyState
        title="Chưa có học sinh nào"
        :message="`${selectedClassName || 'Lớp này'} chưa có học sinh nào. Hãy thêm học sinh đầu tiên.`"
        icon="users"
      >
        <template #action>
          <button class="btn btn--primary" type="button" @click="openAddModal">
            Thêm học sinh
          </button>
        </template>
      </EmptyState>
    </div>

    <div v-else-if="filteredStudents.length === 0" class="card empty-search">
      Không tìm thấy học sinh nào phù hợp với "{{ searchQuery }}"
    </div>

    <StudentsTable
      v-else
      :students="filteredStudents"
      :generating-code-student-id="generatingCodeStudentId"
      :revoking-code-student-id="revokingCodeStudentId"
      :copied-student-id="copiedStudentId"
      @edit="openEditModal"
      @delete="confirmDelete"
      @generate-code="handleGenerateCode"
      @copy-code="handleCopyCode"
      @revoke-code="confirmRevokeCode"
    />

    <StudentFormModal
      :is-open="isFormModalOpen"
      :form-mode="formMode"
      :form-data="formData"
      :selected-class-name="selectedClassName"
      :is-submitting="isFormSubmitting"
      :error-message="formError"
      @close="closeFormModal"
      @submit="submitForm"
      @update:form-data="updateFormData"
    />

    <ConfirmDialog
      :is-open="isDeleteConfirmOpen"
      title="Xác nhận xóa"
      :message="getDeleteStudentConfirmMessage(deleteTarget?.full_name)"
      confirm-text="Xóa học sinh"
      is-danger
      :is-loading="isDeleteLoading"
      @confirm="handleDelete"
      @cancel="closeDeleteConfirm"
    />

    <ConfirmDialog
      :is-open="isRevokeConfirmOpen"
      title="Thu hồi mã phụ huynh"
      :message="STUDENT_REVOKE_CONFIRM_MESSAGE"
      confirm-text="Thu hồi"
      is-danger
      :is-loading="Boolean(revokingCodeStudentId)"
      @confirm="handleRevokeCode"
      @cancel="closeRevokeConfirm"
    />
  </div>
</template>

<style scoped>
.page-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-actions {
  display: flex;
  justify-content: flex-end;
}

.toolbar-card {
  padding: var(--spacing-4);
}

.toolbar-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.toolbar-search {
  min-width: 0;
}

.empty-search {
  padding: var(--spacing-6);
  text-align: center;
  color: var(--color-text-muted);
}

.mt-2 {
  margin-top: var(--spacing-2);
}

@media (min-width: 768px) {
  .toolbar-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
    align-items: end;
  }
}
</style>
