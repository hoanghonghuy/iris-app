<script setup>
import { ref, onMounted } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateTime } from '../../helpers/dateFormatter'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import EmptyState from '../../components/common/EmptyState.vue'

const appointments = ref([])
const isLoading = ref(true)
const errorMessage = ref('')

const fetchAppointments = async () => {
  isLoading.value = true
  errorMessage.value = ''
  
  try {
    const data = await teacherService.getAppointments()
    appointments.value = data.data || []
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchAppointments()
})

const getStatusBadge = (status) => {
  switch (status) {
    case 'PENDING': return 'badge badge--outline text-muted'
    case 'CONFIRMED': return 'badge badge--success'
    case 'CANCELLED': return 'badge badge--danger'
    case 'COMPLETED': return 'badge badge--primary'
    default: return 'badge badge--outline'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'PENDING': return 'Chờ xác nhận'
    case 'CONFIRMED': return 'Đã xác nhận'
    case 'CANCELLED': return 'Đã hủy'
    case 'COMPLETED': return 'Hoàn thành'
    default: return status
  }
}
</script>

<template>
  <div class="teacher-schedule">
    <div class="flex justify-end items-center mb-6">
      <button type="button" class="btn btn--primary" disabled title="Tính năng tạo slot đang phát triển">
        + Tạo khung giờ trống
      </button>
    </div>

    <!-- Error State -->
    <div v-if="errorMessage" class="p-4 mb-6 bg-red-50 text-danger rounded border border-red-200">
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button type="button" class="btn btn--outline mt-2" @click="fetchAppointments">Thử lại</button>
    </div>

    <!-- Loading State -->
    <LoadingSpinner v-else-if="isLoading" message="Đang tải lịch hẹn..." />

    <!-- Content -->
    <div v-else class="card">
      <EmptyState 
        v-if="appointments.length === 0" 
        title="Chưa có lịch hẹn nào" 
        message="Không có lịch hẹn nào được lên kế hoạch trong thời gian này."
      />

      <div v-else class="table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th>Thời gian</th>
              <th>Phụ huynh</th>
              <th>Học sinh liên quan</th>
              <th>Nội dung / Ghi chú</th>
              <th>Trạng thái</th>
              <th class="text-right">Thao tác</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="app in appointments" :key="app.appointment_id">
              <td class="whitespace-nowrap font-medium">
                {{ formatDateTime(app.start_time) }}
              </td>
              <td>{{ app.parent_name || 'N/A' }}</td>
              <td>{{ app.student_name || 'N/A' }}</td>
              <td class="text-sm text-muted max-w-xs truncate" :title="app.notes">{{ app.notes || '-' }}</td>
              <td>
                <span :class="getStatusBadge(app.status)">{{ getStatusText(app.status) }}</span>
              </td>
              <td class="text-right">
                <button type="button" class="btn btn--sm btn--outline" disabled>Chi tiết</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mb-6 { margin-bottom: var(--spacing-6); }
.p-4 { padding: var(--spacing-4); }
.mt-2 { margin-top: var(--spacing-2); }
.rounded { border-radius: var(--radius); }
.border { border: 1px solid var(--color-border); }
.bg-red-50 { background-color: var(--color-danger-soft-bg); }
.border-red-200 { border-color: var(--color-danger-soft-border); }
.text-danger { color: var(--color-danger-soft-text); }
.text-sm { font-size: var(--font-size-sm); }
.text-right { text-align: right; }
.whitespace-nowrap { white-space: nowrap; }
.max-w-xs { max-width: 15rem; }
.truncate { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

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
</style>
