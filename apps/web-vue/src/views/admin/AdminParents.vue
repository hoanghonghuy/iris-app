<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { Link2, Pencil, Phone, X } from 'lucide-vue-next'
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

const parents = ref([])
const totalPages = ref(0)
const currentPage = ref(1)
const totalItems = ref(0)
const isLoading = ref(true)
const errorMessage = ref('')

const searchQuery = ref('')
const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')

const isEditModalOpen = ref(false)
const editLoading = ref(false)
const editError = ref('')
const editForm = ref({
  parent_id: '',
  full_name: '',
  phone: '',
  school_id: '',
})

const isAssignModalOpen = ref(false)
const assignTarget = ref(null)
const assignLoading = ref(false)
const assignError = ref('')

const isUnassignOpen = ref(false)
const unassignTarget = ref(null)
const unassignLoading = ref(false)

const schools = ref([])
const classes = ref([])
const students = ref([])
const selectedSchoolId = ref('')
const selectedClassId = ref('')
const selectedStudentId = ref('')

const filteredParents = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) {
    return parents.value
  }

  return parents.value.filter((parent) =>
    parent.full_name?.toLowerCase().includes(query) ||
    parent.email?.toLowerCase().includes(query) ||
    parent.phone?.includes(query)
  )
})

async function fetchParents(page = 1) {
  isLoading.value = true
  errorMessage.value = ''
  currentPage.value = page

  try {
    const data = await adminService.getParents({
      limit: PAGE_SIZE,
      offset: (page - 1) * PAGE_SIZE,
    })

    parents.value = Array.isArray(data?.data) ? data.data : []

    if (data?.pagination) {
      totalItems.value = data.pagination.total || 0
      totalPages.value = Math.ceil(totalItems.value / Math.max(data.pagination.limit || PAGE_SIZE, 1)) || 1
    } else {
      totalItems.value = parents.value.length
      totalPages.value = parents.value.length > 0 ? 1 : 0
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

async function fetchSchoolsForSelector() {
  try {
    const data = await adminService.getSchools({ limit: 100, offset: 0 })
    schools.value = Array.isArray(data?.data) ? data.data : []
    selectedSchoolId.value = schools.value[0]?.school_id || ''
  } catch {
    schools.value = []
    selectedSchoolId.value = ''
  }
}

async function fetchClassesForSelector() {
  if (!selectedSchoolId.value) {
    classes.value = []
    students.value = []
    selectedClassId.value = ''
    selectedStudentId.value = ''
    return
  }

  classes.value = []
  students.value = []
  selectedClassId.value = ''
  selectedStudentId.value = ''

  try {
    const data = await adminService.getClassesBySchool(selectedSchoolId.value, { limit: 100, offset: 0 })
    classes.value = Array.isArray(data?.data) ? data.data : []
    selectedClassId.value = classes.value[0]?.class_id || ''
  } catch {
    classes.value = []
  }
}

async function fetchStudentsForSelector() {
  if (!selectedClassId.value) {
    students.value = []
    selectedStudentId.value = ''
    return
  }

  students.value = []
  selectedStudentId.value = ''

  try {
    const data = await adminService.getStudentsByClass(selectedClassId.value, { limit: 100, offset: 0 })
    students.value = Array.isArray(data?.data) ? data.data : []
    selectedStudentId.value = students.value[0]?.student_id || ''
  } catch {
    students.value = []
  }
}

watch(selectedSchoolId, () => {
  fetchClassesForSelector()
})

watch(selectedClassId, () => {
  fetchStudentsForSelector()
})

onMounted(() => {
  fetchParents()
  fetchSchoolsForSelector()
})

function openAssignModal(parent) {
  assignTarget.value = parent
  assignError.value = ''
  isAssignModalOpen.value = true
}

function openEditModal(parent) {
  editError.value = ''
  editForm.value = {
    parent_id: parent.parent_id,
    full_name: parent.full_name || '',
    phone: parent.phone || '',
    school_id: parent.school_id || selectedSchoolId.value || '',
  }
  isEditModalOpen.value = true
}

async function handleEdit() {
  if (!editForm.value.parent_id || !editForm.value.full_name.trim() || !editForm.value.school_id) {
    editError.value = 'Vui lòng nhập đầy đủ thông tin bắt buộc'
    return
  }

  editLoading.value = true
  editError.value = ''
  try {
    await adminService.updateParent(editForm.value.parent_id, {
      full_name: editForm.value.full_name.trim(),
      phone: editForm.value.phone?.trim() || '',
      school_id: editForm.value.school_id,
    })

    isEditModalOpen.value = false
    await fetchParents(currentPage.value)
  } catch (error) {
    editError.value = extractErrorMessage(error) || 'Không thể cập nhật phụ huynh'
  } finally {
    editLoading.value = false
  }
}

async function handleAssign() {
  if (!selectedStudentId.value || !assignTarget.value) {
    return
  }

  assignLoading.value = true
  assignError.value = ''
  try {
    await adminService.assignParentToStudent(assignTarget.value.parent_id, selectedStudentId.value)
    isAssignModalOpen.value = false
    await fetchParents(currentPage.value)
  } catch (error) {
    assignError.value = extractErrorMessage(error)
  } finally {
    assignLoading.value = false
  }
}

function openUnassignDialog(parent, student) {
  unassignTarget.value = {
    parent_id: parent.parent_id,
    student_id: student.student_id,
    student_name: student.full_name,
  }
  isUnassignOpen.value = true
}

async function handleUnassign() {
  if (!unassignTarget.value) {
    return
  }

  unassignLoading.value = true
  try {
    await adminService.unassignParentFromStudent(unassignTarget.value.parent_id, unassignTarget.value.student_id)
    isUnassignOpen.value = false
    unassignTarget.value = null
    await fetchParents(currentPage.value)
  } catch (error) {
    errorMessage.value = `Lỗi hủy gán: ${extractErrorMessage(error)}`
    isUnassignOpen.value = false
  } finally {
    unassignLoading.value = false
  }
}
</script>

<template>
  <div class="admin-parents page-stack">
    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchParents(currentPage)">Thử lại</button>
    </div>

    <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu..." />

    <div v-else class="page-stack">
      <div v-if="parents.length > 0" class="card toolbar-card">
        <div class="toolbar-grid">
          <input
            v-model="searchQuery"
            type="search"
            class="form-input"
            placeholder="Tìm theo tên, email, SĐT..."
          />
        </div>
      </div>

      <div v-if="parents.length === 0" class="card">
        <EmptyState
          title="Chưa có phụ huynh nào"
          message="Chưa có phụ huynh nào đăng ký tài khoản trên hệ thống."
        />
      </div>

      <div v-else-if="filteredParents.length === 0" class="card empty-search">
        Không tìm thấy phụ huynh nào phù hợp với "{{ searchQuery }}"
      </div>

      <template v-else>
        <div class="card desktop-table">
          <div class="table-responsive">
            <table class="table">
              <thead>
                <tr>
                  <th>Họ tên</th>
                  <th>Email</th>
                  <th>Học sinh quản lý</th>
                  <th class="action-column">Gán học sinh</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="parent in filteredParents" :key="parent.parent_id">
                  <td>
                    <div class="font-medium">{{ parent.full_name || 'Chưa cập nhật' }}</div>
                    <div class="text-xs text-muted mt-1 phone-line">
                      <Phone :size="12" />
                      <span>{{ parent.phone || '-' }}</span>
                    </div>
                  </td>
                  <td class="text-muted">{{ parent.email || '-' }}</td>
                  <td>
                    <template v-if="parent.children && parent.children.length > 0">
                      <div class="flex flex-wrap gap-1">
                        <span
                          v-for="child in parent.children"
                          :key="child.student_id"
                          class="badge badge--outline badge--sm flex items-center gap-1"
                        >
                          {{ child.full_name }}
                          <button
                            class="badge-remove-btn"
                            type="button"
                            title="Hủy gán học sinh"
                            @click="openUnassignDialog(parent, child)"
                          >
                            <X :size="11" />
                          </button>
                        </span>
                      </div>
                    </template>
                    <span v-else class="text-muted text-sm italic">Chưa ghép học sinh</span>
                  </td>
                  <td class="action-column">
                    <div class="table-action-buttons">
                      <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(parent)">
                        <Pencil :size="14" />
                        <span>Sửa</span>
                      </button>
                      <button class="btn btn--sm btn--outline" type="button" @click="openAssignModal(parent)">
                        <Link2 :size="14" />
                        <span>Gán HS</span>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="mobile-list">
          <article v-for="parent in filteredParents" :key="parent.parent_id" class="card mobile-card">
            <div class="mobile-card__head">
              <p class="mobile-card__title">{{ parent.full_name || 'Chưa cập nhật' }}</p>
              <button class="btn btn--sm btn--outline" type="button" @click="openAssignModal(parent)">
                <Link2 :size="14" />
                <span>Gán HS</span>
              </button>
            </div>

            <p class="mobile-card__meta">{{ parent.email || '-' }}</p>
            <p class="mobile-card__meta mobile-card__phone">
              <Phone :size="12" />
              <span>{{ parent.phone || '-' }}</span>
            </p>

            <div v-if="parent.children && parent.children.length > 0" class="mobile-card__chips">
              <span
                v-for="child in parent.children"
                :key="child.student_id"
                class="badge badge--outline badge--sm flex items-center gap-1"
              >
                {{ child.full_name }}
                <button
                  class="badge-remove-btn"
                  type="button"
                  title="Hủy gán học sinh"
                  @click="openUnassignDialog(parent, child)"
                >
                  <X :size="11" />
                </button>
              </span>
            </div>
            <p v-else class="mobile-card__meta italic">Chưa ghép học sinh</p>
            <div class="mobile-card__actions">
              <button class="btn btn--sm btn--outline" type="button" @click="openEditModal(parent)">
                Sửa
              </button>
            </div>
          </article>
        </div>

        <PaginationBar
          :current-page="currentPage"
          :total-pages="totalPages"
          :total-items="totalItems"
          :limit="PAGE_SIZE"
          @page-change="fetchParents"
        />
      </template>
    </div>

    <ActionModal
      :is-open="isEditModalOpen"
      :title="`Sửa thông tin phụ huynh - ${editForm.full_name || ''}`"
      @close="isEditModalOpen = false"
    >
      <form class="modal-form" @submit.prevent="handleEdit">
        <div v-if="editError" class="alert alert--error">
          {{ editError }}
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="parentFullName">Họ và tên <span class="text-danger">*</span></label>
          <input
            id="parentFullName"
            v-model="editForm.full_name"
            type="text"
            class="form-input"
            placeholder="Ví dụ: Trần Thị B"
            :disabled="editLoading"
            required
          />
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="parentPhone">Số điện thoại</label>
          <input
            id="parentPhone"
            v-model="editForm.phone"
            type="text"
            class="form-input"
            placeholder="Nhập số điện thoại"
            :disabled="editLoading"
          />
        </div>

        <div v-if="isSuperAdmin" class="form-group mb-0">
          <label class="form-label" for="parentSchool">Trường học <span class="text-danger">*</span></label>
          <select id="parentSchool" v-model="editForm.school_id" class="form-input" :disabled="editLoading" required>
            <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
              {{ school.name }}
            </option>
          </select>
        </div>

        <div class="modal-actions">
          <button class="btn btn--outline" type="button" :disabled="editLoading" @click="isEditModalOpen = false">
            Hủy
          </button>
          <button class="btn btn--primary" type="submit" :disabled="editLoading">
            {{ editLoading ? 'Đang lưu...' : 'Lưu lại' }}
          </button>
        </div>
      </form>
    </ActionModal>

    <ActionModal
      :is-open="isAssignModalOpen"
      :title="`Gán học sinh - ${assignTarget?.full_name || ''}`"
      @close="isAssignModalOpen = false"
    >
      <div class="modal-form">
        <div v-if="assignError" class="alert alert--error">
          {{ assignError }}
        </div>

        <div v-if="authStore.currentUserRole === 'SUPER_ADMIN'" class="form-group mb-0">
          <label class="form-label">Chọn trường</label>
          <select v-model="selectedSchoolId" class="form-input">
            <option v-for="school in schools" :key="school.school_id" :value="school.school_id">
              {{ school.name }}
            </option>
          </select>
        </div>

        <div class="form-group mb-0">
          <label class="form-label">Chọn lớp <span class="text-danger">*</span></label>
          <select v-model="selectedClassId" class="form-input" :disabled="classes.length === 0">
            <option v-if="classes.length === 0" value="" disabled>Không có lớp</option>
            <option v-for="cls in classes" :key="cls.class_id" :value="cls.class_id">
              {{ cls.name }}
            </option>
          </select>
        </div>

        <div class="form-group mb-0">
          <label class="form-label">Chọn học sinh <span class="text-danger">*</span></label>
          <select v-model="selectedStudentId" class="form-input" :disabled="students.length === 0">
            <option v-if="students.length === 0" value="" disabled>Không có học sinh</option>
            <option v-for="student in students" :key="student.student_id" :value="student.student_id">
              {{ student.full_name }}
            </option>
          </select>
        </div>

        <div class="modal-actions">
          <button class="btn btn--outline" type="button" :disabled="assignLoading" @click="isAssignModalOpen = false">
            Hủy
          </button>
          <button
            class="btn btn--primary"
            type="button"
            :disabled="assignLoading || !selectedStudentId"
            @click="handleAssign"
          >
            {{ assignLoading ? 'Đang gán...' : 'Gán học sinh' }}
          </button>
        </div>
      </div>
    </ActionModal>

    <ConfirmDialog
      :is-open="isUnassignOpen"
      title="Xác nhận hủy gán"
      :message="`Bạn có chắc muốn hủy gán học sinh '${unassignTarget?.student_name || ''}' khỏi phụ huynh này?`"
      confirm-text="Hủy gán"
      is-danger
      :is-loading="unassignLoading"
      @confirm="handleUnassign"
      @cancel="isUnassignOpen = false"
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

.toolbar-card,
.desktop-table,
.mobile-card {
  padding: var(--spacing-4);
}

.toolbar-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: minmax(0, 1fr);
}

.empty-search {
  padding: var(--spacing-6);
  text-align: center;
  color: var(--color-text-muted);
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
  vertical-align: middle;
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

.action-column {
  width: 220px;
  min-width: 220px;
  text-align: right !important;
  white-space: nowrap;
}

.table-action-buttons {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.table-action-buttons .btn,
.mobile-card__head .btn,
.mobile-card__actions .btn {
  gap: 0.35rem;
}

.phone-line,
.mobile-card__phone {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
}

.badge--sm {
  font-size: 0.7rem;
  padding: 2px 6px;
}

.badge-remove-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-muted);
  padding: 0;
  line-height: 1;
}

.badge-remove-btn:hover {
  color: var(--color-danger);
}

.mobile-list {
  display: none;
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

.mobile-card__title,
.mobile-card__meta {
  margin: 0;
}

.mobile-card__title {
  color: var(--color-text);
  font-weight: 700;
}

.mobile-card__meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.mobile-card__chips {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.mobile-card__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.mt-1 {
  margin-top: var(--spacing-1);
}

.mt-2 {
  margin-top: var(--spacing-2);
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
