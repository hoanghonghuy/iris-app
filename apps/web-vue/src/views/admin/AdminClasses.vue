<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'
import PaginationBar from '../../components/PaginationBar.vue'
import ConfirmDialog from '../../components/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'

const authStore = useAuthStore()

const PAGE_SIZE = 10

const schools = ref([])
const selectedSchoolId = ref('')

const classesList = ref([])
const totalPages = ref(0)
const currentPage = ref(1)
const totalItems = ref(0)
const isLoading = ref(true)
const errorMessage = ref('')

const isModalOpen = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const formMode = ref('add')
const formData = ref({ id: '', name: '', school_year: '', school_id: '' })

const isConfirmOpen = ref(false)
const itemToDelete = ref(null)

const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')
const selectedSchoolName = computed(() => {
  return schools.value.find((school) => school.school_id === selectedSchoolId.value)?.name || ''
})

async function fetchSchools() {
  try {
    const data = await adminService.getSchools({ limit: 100, offset: 0 })
    schools.value = Array.isArray(data?.data) ? data.data : []

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

async function fetchClasses(page = 1) {
  if (!selectedSchoolId.value) return

  isLoading.value = true
  errorMessage.value = ''
  currentPage.value = page

  try {
    const data = await adminService.getClassesBySchool(selectedSchoolId.value, {
      limit: PAGE_SIZE,
      offset: (page - 1) * PAGE_SIZE,
    })

    classesList.value = Array.isArray(data?.data) ? data.data : []

    if (data?.pagination) {
      totalItems.value = data.pagination.total || 0
      totalPages.value = Math.ceil(data.pagination.total / data.pagination.limit) || 1
    } else {
      totalItems.value = classesList.value.length
      totalPages.value = classesList.value.length > 0 ? 1 : 0
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách lớp'
  } finally {
    isLoading.value = false
  }
}

watch(selectedSchoolId, () => {
  fetchClasses(1)
})

onMounted(async () => {
  if (!authStore.currentUser && authStore.isAuthenticated) {
    await authStore.fetchCurrentUser()
  }

  fetchSchools()
})

function openAddModal() {
  formMode.value = 'add'
  formData.value = {
    id: '',
    name: '',
    school_year: '2023-2024',
    school_id: selectedSchoolId.value,
  }
  modalError.value = ''
  isModalOpen.value = true
}

function openEditModal(cls) {
  formMode.value = 'edit'
  formData.value = {
    ...cls,
    id: cls.class_id,
  }
  modalError.value = ''
  isModalOpen.value = true
}

async function handleSave() {
  if (!formData.value.name || !formData.value.school_year || !formData.value.school_id) {
    modalError.value = 'Vui lòng nhập đầy đủ thông tin'
    return
  }

  isSubmitting.value = true
  modalError.value = ''

  try {
    if (formMode.value === 'add') {
      await adminService.createClass({
        name: formData.value.name,
        school_year: formData.value.school_year,
        school_id: formData.value.school_id,
      })
    } else {
      await adminService.updateClass(formData.value.id, {
        name: formData.value.name,
        school_year: formData.value.school_year,
      })
    }

    isModalOpen.value = false

    if (formMode.value === 'add' && selectedSchoolId.value !== formData.value.school_id) {
      selectedSchoolId.value = formData.value.school_id
    } else {
      fetchClasses(currentPage.value)
    }
  } catch (error) {
    modalError.value = extractErrorMessage(error) || 'Không thể lưu lớp học'
  } finally {
    isSubmitting.value = false
  }
}

function confirmDelete(cls) {
  itemToDelete.value = cls
  isConfirmOpen.value = true
}

async function handleDelete() {
  if (!itemToDelete.value) return

  isSubmitting.value = true

  try {
    await adminService.deleteClass(itemToDelete.value.class_id)
    isConfirmOpen.value = false
    itemToDelete.value = null
    fetchClasses(currentPage.value)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể xóa lớp'
    isConfirmOpen.value = false
  } finally {
    isSubmitting.value = false
  }
}
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
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchClasses(currentPage)">Thử lại</button>
    </div>

    <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu..." />

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
      @close="isModalOpen = false"
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
          <button type="button" class="btn btn--outline" :disabled="isSubmitting" @click="isModalOpen = false">
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
