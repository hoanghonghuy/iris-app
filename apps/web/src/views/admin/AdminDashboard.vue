<script setup>
import { computed, ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { ADMIN_LOAD_ERROR_TITLE, ADMIN_RETRY_BUTTON_TEXT } from '../../helpers/adminConfig'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import AnalyticsTimeseriesPanel from '../../components/charts/AnalyticsTimeseriesPanel.vue'

const authStore = useAuthStore()
const analytics = ref(null)
const isLoading = ref(true)
const errorMessage = ref('')
const timeseries = ref(null)
const chartsLoading = ref(false)
const chartsError = ref('')
const selectedRange = ref('14d')
const activeTab = ref('overview')
const RANGE_OPTIONS = [
  { value: '7d', label: '7 ngày' },
  { value: '14d', label: '14 ngày' },
  { value: '30d', label: '30 ngày' },
]

const displayName = computed(() => {
  const user = authStore.currentUser
  return user?.full_name || user?.email?.split('@')[0] || 'Admin'
})

const statCards = computed(() => [
  { label: 'Trường học', value: analyticsView.value.total_schools, to: '/admin/schools' },
  { label: 'Lớp học', value: analyticsView.value.total_classes, to: '/admin/classes' },
  { label: 'Giáo viên', value: analyticsView.value.total_teachers, to: '/admin/teachers' },
  { label: 'Học sinh', value: analyticsView.value.total_students, to: '/admin/students' },
  { label: 'Phụ huynh', value: analyticsView.value.total_parents, to: '/admin/parents' },
])

const quickActions = [
  { label: 'Quản lý trường', to: '/admin/schools' },
  { label: 'Quản lý người dùng', to: '/admin/users' },
  { label: 'Quản lý lớp', to: '/admin/classes' },
]

const analyticsView = computed(() => {
  const data = analytics.value || {}
  const attendance = data.today_attendance || {}
  const rawAttendanceRate = attendance.attendance_rate ?? data.today_attendance_rate ?? 0
  const attendanceRate = Number(rawAttendanceRate)
  const recentHealthIssues = Array.isArray(data.recent_health_issues)
    ? data.recent_health_issues
    : []

  return {
    ...data,
    total_schools: data.total_schools ?? 0,
    total_classes: data.total_classes ?? 0,
    total_students: data.total_students ?? 0,
    total_teachers: data.total_teachers ?? 0,
    total_parents: data.total_parents ?? 0,
    total_users: data.total_users ?? 0,
    today_attendance: {
      not_recorded: attendance.not_recorded ?? 0,
      present: attendance.present ?? 0,
      absent: attendance.absent ?? 0,
      attendance_rate: Number.isFinite(attendanceRate) ? attendanceRate : 0,
    },
    recent_health_issues: recentHealthIssues,
    recent_health_count: data.recent_health_alerts_24h ?? recentHealthIssues.length,
  }
})

const fetchAnalytics = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const res = await adminService.getAnalytics()
    analytics.value = res?.data ?? res ?? {}
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

const fetchTimeseries = async () => {
  chartsLoading.value = true
  chartsError.value = ''
  try {
    const res = await adminService.getAnalyticsTimeseries({
      range: selectedRange.value,
      interval: 'day',
    })
    timeseries.value = res?.data ?? res ?? null
  } catch (error) {
    chartsError.value = extractErrorMessage(error)
    timeseries.value = null
  } finally {
    chartsLoading.value = false
  }
}

const refreshAll = async () => {
  await Promise.all([fetchAnalytics(), fetchTimeseries()])
}

const changeRange = async (rangeValue) => {
  if (selectedRange.value === rangeValue || chartsLoading.value) return
  selectedRange.value = rangeValue
  await fetchTimeseries()
}

onMounted(async () => {
  if (!authStore.currentUser) {
    await authStore.fetchCurrentUser()
  }
  await refreshAll()
})
</script>

<template>
  <div class="admin-dashboard">
    <div class="dashboard-hero mb-6">
      <div>
        <h2>Xin chào, {{ displayName }}</h2>
      </div>
      <button class="btn btn--outline" type="button" @click="refreshAll" :disabled="isLoading || chartsLoading">
        Làm mới
      </button>
    </div>

    <div class="view-switch mb-6">
      <button
        type="button"
        class="view-switch__btn"
        :class="{ 'view-switch__btn--active': activeTab === 'overview' }"
        @click="activeTab = 'overview'"
      >
        Tổng quan
      </button>
      <button
        type="button"
        class="view-switch__btn"
        :class="{ 'view-switch__btn--active': activeTab === 'charts' }"
        @click="activeTab = 'charts'"
      >
        Biểu đồ
      </button>
    </div>

    <LoadingSpinner v-if="isLoading" message="Đang tải dữ liệu thống kê..." />

    <div v-else-if="errorMessage" class="alert alert--error">
      <p class="font-bold">{{ ADMIN_LOAD_ERROR_TITLE }}</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="refreshAll">
        {{ ADMIN_RETRY_BUTTON_TEXT }}
      </button>
    </div>

    <div v-else-if="analytics" class="dashboard-content">
      <!-- Tab Tổng quan -->
      <div v-show="activeTab === 'overview'">
        <!-- Cards thống kê chính -->
        <div class="grid-cards mb-8">
          <RouterLink
            v-for="card in statCards"
            :key="card.to"
            :to="card.to"
            class="card stat-card p-4"
          >
            <p class="text-sm font-medium text-muted">{{ card.label }}</p>
            <p class="text-2xl font-bold mt-1">{{ card.value }}</p>
          </RouterLink>
        </div>

        <div class="grid-2-cols gap-6">
          <div class="card p-4">
            <p class="text-sm font-medium text-muted uppercase">Tỷ lệ điểm danh hôm nay</p>
            <p class="text-3xl font-bold mt-2">
              {{ analyticsView.today_attendance.attendance_rate.toFixed(1) }}%
            </p>
          </div>

          <div class="card p-4">
            <p class="text-sm font-medium text-muted uppercase">Cảnh báo sức khỏe 24h</p>
            <p class="text-3xl font-bold mt-2">{{ analyticsView.recent_health_count }}</p>
          </div>
        </div>

        <section class="quick-actions mt-8">
          <h3>Quản lý nhanh</h3>
          <div class="quick-grid">
            <RouterLink
              v-for="action in quickActions"
              :key="action.to"
              :to="action.to"
              class="card quick-action"
            >
              {{ action.label }}
            </RouterLink>
          </div>
        </section>
      </div>

      <!-- Tab Biểu đồ -->
      <div v-show="activeTab === 'charts'">
        <section class="card range-filter mb-6">
          <div class="range-filter__heading">
            <h3>Khoảng thời gian biểu đồ</h3>
            <span class="text-muted text-sm">Xem nhanh theo 7, 14 hoặc 30 ngày gần nhất</span>
          </div>
          <div class="range-filter__actions">
            <button
              v-for="option in RANGE_OPTIONS"
              :key="option.value"
              type="button"
              class="btn"
              :class="selectedRange === option.value ? 'btn--primary' : 'btn--outline'"
              :disabled="chartsLoading"
              @click="changeRange(option.value)"
            >
              {{ option.label }}
            </button>
          </div>
        </section>

        <div v-if="chartsLoading" class="text-muted text-sm">Đang tải biểu đồ...</div>
        <div v-else-if="chartsError" class="alert alert--error">{{ chartsError }}</div>
        <AnalyticsTimeseriesPanel
          v-else-if="timeseries"
          :payload="timeseries"
          :title="`Biểu đồ xu hướng (${selectedRange.replace('d', ' ngày')})`"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Layout cục bộ */
.dashboard-hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  background:
    radial-gradient(circle at top right, color-mix(in srgb, var(--color-primary) 14%, transparent), transparent 60%),
    var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.dashboard-hero h2,
.quick-actions h3 {
  margin: 0;
  color: var(--color-text);
}

.dashboard-hero h2 {
  margin-top: 0;
  font-size: clamp(1.5rem, 3.2vw, var(--font-size-3xl));
  font-weight: 800;
}

.grid-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-4);
}

.stat-card,
.quick-action {
  color: inherit;
  text-decoration: none;
  transition:
    border-color var(--transition-fast),
    box-shadow var(--transition-fast),
    transform var(--transition-fast);
}

.stat-card:hover,
.quick-action:hover {
  border-color: color-mix(in srgb, var(--color-primary) 35%, var(--color-border));
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.stat-card:focus-visible,
.quick-action:focus-visible {
  box-shadow: 0 0 0 3px var(--color-primary-focus-ring);
}

.quick-actions h3 {
  font-size: var(--font-size-lg);
  font-weight: 800;
  margin-bottom: var(--spacing-3);
}

.range-filter {
  padding: var(--spacing-3);
}

.range-filter__heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-3);
}

.range-filter__heading h3 {
  margin: 0;
  font-size: var(--font-size-base);
  color: var(--color-text);
}

.range-filter__actions {
  display: flex;
  gap: var(--spacing-2);
  flex-wrap: wrap;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-3);
}

.quick-action {
  min-height: 5.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-4);
  text-align: center;
  font-weight: 700;
  border-radius: var(--radius-lg);
}

@media (min-width: 768px) {
  .grid-cards {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1024px) {
  .grid-cards {
    grid-template-columns: repeat(5, 1fr);
  }
  .grid-2-cols {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .dashboard-hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .range-filter__heading {
    flex-direction: column;
    align-items: flex-start;
  }

  .dashboard-hero .btn {
    width: 100%;
  }

  .quick-grid {
    grid-template-columns: 1fr;
  }

  .range-filter__actions .btn {
    flex: 1 1 100%;
  }
}
</style>
