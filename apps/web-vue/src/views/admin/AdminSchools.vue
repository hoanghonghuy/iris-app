<script setup>
import { ref, onMounted } from 'vue'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'
import PaginationBar from '../../components/PaginationBar.vue'
import ConfirmDialog from '../../components/ConfirmDialog.vue'
import ActionModal from '../../components/ActionModal.vue'

// State danh sách
const schools = ref([])
const totalPages = ref(0)
const currentPage = ref(1)
const totalItems = ref(0)
const isLoading = ref(true)
const errorMessage = ref('')
const PAGE_SIZE = 10

// State Modal Form
const isModalOpen = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const formMode = ref('add') // add | edit
const formData = ref({ id: '', name: '', address: '' })

// State Confirm Dialog
const isConfirmOpen = ref(false)
const itemToDelete = ref(null)

const fetchSchools = async (page = 1) => {
  isLoading.value = true
  errorMessage.value = ''
  currentPage.value = page
  
  try {
    const data = await adminService.getSchools({
      limit: PAGE_SIZE,
      offset: (page - 1) * PAGE_SIZE,
    })
    schools.value = data.data || []
    
    if (data.pagination) {
      totalItems.value = data.pagination.total || 0
      totalPages.value = Math.ceil(data.pagination.total / data.pagination.limit) || 1
    } else {
      totalItems.value = 0
      totalPages.value = 0
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchSchools()
})

// Handlers cho Form
const openAddModal = () => {
  formMode.value = 'add'
  formData.value = { id: '', name: '', address: '' }
  modalError.value = ''
  isModalOpen.value = true
}

const openEditModal = (school) => {
  formMode.value = 'edit'
  formData.value = { ...school, id: school.school_id }
  modalError.value = ''
  isModalOpen.value = true
}

const handleSave = async () => {
  if (!formData.value.name) {
    modalError.value = 'Tên trường không được để trống'
    return
  }

  isSubmitting.value = true
  modalError.value = ''

  try {
    if (formMode.value === 'add') {
      await adminService.createSchool({
        name: formData.value.name,
        address: formData.value.address
      })
    } else {
      await adminService.updateSchool(formData.value.id, {
        name: formData.value.name,
        address: formData.value.address
      })
    }
    
    isModalOpen.value = false
    fetchSchools(currentPage.value) // Tải lại danh sách
  } catch (error) {
    modalError.value = extractErrorMessage(error)
  } finally {
    isSubmitting.value = false
  }
}

// Handlers cho Delete
const confirmDelete = (school) => {
  itemToDelete.value = school
  isConfirmOpen.value = true
}

const handleDelete = async () => {
  if (!itemToDelete.value) return

  isSubmitting.value = true
  
  try {
    await adminService.deleteSchool(itemToDelete.value.school_id)
    isConfirmOpen.value = false
    itemToDelete.value = null
    fetchSchools(currentPage.value)
  } catch (error) {
    errorMessage.value = `Lỗi xóa: ${extractErrorMessage(error)}`
    isConfirmOpen.value = false
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="admin-schools">
    <div class="flex justify-end items-center mb-6">
      <button class="btn btn--primary" @click="openAddModal">
        + Thêm trường học
      </button>
    </div>

    <!-- Error State -->
    <div v-if="errorMessage" class="p-4 mb-6 bg-red-50 text-danger rounded border border-red-200">
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" @click="fetchSchools(currentPage)">Thử lại</button>
    </div>

    <!-- Loading State -->
    <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu..." />

    <!-- Content -->
    <div v-else class="card">
      <EmptyState 
        v-if="schools.length === 0" 
        title="Chưa có trường học nào" 
        message="Hãy thêm trường học đầu tiên để bắt đầu quản lý."
      >
        <template #action>
          <button class="btn btn--primary" @click="openAddModal">Thêm trường học</button>
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
                  <button class="btn btn--sm btn--outline" @click="openEditModal(school)">
                    Sửa
                  </button>
                  <button class="btn btn--sm btn--danger" @click="confirmDelete(school)">
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

    <!-- Modal Form -->
    <ActionModal 
      :is-open="isModalOpen" 
      :title="formMode === 'add' ? 'Thêm trường học mới' : 'Sửa thông tin trường'"
      @close="isModalOpen = false"
    >
      <form @submit.prevent="handleSave" class="flex-col gap-4">
        <div v-if="modalError" class="p-3 bg-red-50 text-danger text-sm rounded border border-red-200">
          {{ modalError }}
        </div>
        
        <div class="form-group">
          <label class="form-label" for="schoolName">Tên trường <span class="text-danger">*</span></label>
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
          <button type="button" class="btn btn--outline" @click="isModalOpen = false" :disabled="isSubmitting">
            Hủy
          </button>
          <button type="submit" class="btn btn--primary" :disabled="isSubmitting">
            {{ isSubmitting ? 'Đang lưu...' : 'Lưu lại' }}
          </button>
        </div>
      </form>
    </ActionModal>

    <!-- Confirm Dialog -->
    <ConfirmDialog 
      :is-open="isConfirmOpen"
      title="Xác nhận xóa"
      :message="`Bạn có chắc chắn muốn xóa trường '${itemToDelete?.name}' không? Hành động này không thể hoàn tác.`"
      confirm-text="Xóa trường"
      is-danger
      :is-loading="isSubmitting"
      @confirm="handleDelete"
      @cancel="isConfirmOpen = false"
    />
  </div>
</template>

<style scoped>
.mb-6 { margin-bottom: var(--spacing-6); }
.p-3 { padding: var(--spacing-3); }
.p-4 { padding: var(--spacing-4); }
.gap-2 { gap: var(--spacing-2); }
.gap-4 { gap: var(--spacing-4); }
.mt-2 { margin-top: var(--spacing-2); }
.mt-4 { margin-top: var(--spacing-4); }
.pt-4 { padding-top: var(--spacing-4); }
.border-t { border-top: 1px solid var(--color-border); }
.rounded { border-radius: var(--radius); }
.bg-red-50 { background-color: var(--color-danger-soft-bg); }
.border-red-200 { border-color: var(--color-danger-soft-border); }
.text-right { text-align: right; }

.table-responsive {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th, .table td {
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
