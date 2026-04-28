<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { teacherService } from '../../services/teacherService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'

const authStore = useAuthStore()
const analytics = ref(null)
const classes = ref([])
const isLoading = ref(true)
const errorMessage = ref('')

const displayName = computed(() => {
  const user = authStore.currentUser
  return user?.full_name || user?.email?.split('@')[0] || 'Giáo viên'
})

const analyticsView = computed(() => {
  const data = analytics.value || {}
  return {
    total_classes: data.total_classes ?? classes.value.length ?? 0,
    total_students: data.total_students ?? 0,
    total_posts: data.total_posts ?? data.recent_posts_count ?? 0,
    today_attendance_marked_count:
      data.today_attendance_marked_count ?? data.today_attendance_count ?? 0,
    today_attendance_pending_count: data.today_attendance_pending_count ?? 0,
    pending_appointments: data.pending_appointments ?? 0,
    recent_health_alerts_24h: data.recent_health_alerts_24h ?? 0,
  }
})

const primaryStats = computed(() => [
  { label: 'Lớp', value: analyticsView.value.total_classes, to: '/teacher/classes' },
  { label: 'Trẻ', value: analyticsView.value.total_students, to: '/teacher/classes' },
  { label: 'Tin', value: analyticsView.value.total_posts, to: '/teacher/posts' },
])

const secondaryStats = computed(() => [
  { label: 'Điểm danh hôm nay', value: analyticsView.value.today_attendance_marked_count },
  { label: 'Chưa điểm danh', value: analyticsView.value.today_attendance_pending_count },
  { label: 'Lịch hẹn chờ duyệt', value: analyticsView.value.pending_appointments },
  { label: 'Cảnh báo sức khỏe 24h', value: analyticsView.value.recent_health_alerts_24h },
])

const quickActions = [
  { label: 'Điểm danh', to: '/teacher/attendance' },
  { label: 'Sức khỏe', to: '/teacher/health' },
  { label: 'Bảng tin', to: '/teacher/posts' },
]

async function fetchDashboard() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const [analyticsRes, classesRes] = await Promise.all([
      teacherService.getAnalytics(),
      teacherService.getMyClasses(),
    ])

    analytics.value = analyticsRes?.data ?? analyticsRes ?? {}
    classes.value = normalizeListResponse(classesRes)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  if (!authStore.currentUser) {
    await authStore.fetchCurrentUser()
  }
  fetchDashboard()
})
</script>

<template>
  <div class="teacher-dashboard">
    <div class="dashboard-hero mb-6">
      <div>
        <h2>Xin chào, {{ displayName }}</h2>
        <p class="hero-copy">
          Hôm nay bạn có {{ analyticsView.total_students }} học sinh cần theo dõi.
        </p>
      </div>
      <button class="btn btn--outline" type="button" @click="fetchDashboard" :disabled="isLoading">
        Làm mới
      </button>
    </div>

    <LoadingSpinner v-if="isLoading" message="Đang tải dữ liệu tổng quan..." />

    <div
      v-else-if="errorMessage"
      class="alert alert--error"
    >
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button type="button" class="btn btn--outline mt-2" @click="fetchDashboard">Thử lại</button>
    </div>

    <div v-else class="dashboard-content">
      <div class="stats-grid mb-6">
        <RouterLink
          v-for="card in primaryStats"
          :key="card.to"
          :to="card.to"
          class="stat-card stat-card--link"
        >
          <span class="stat-label">{{ card.label }}</span>
          <strong class="stat-value">{{ card.value }}</strong>
        </RouterLink>
        <div v-for="card in secondaryStats" :key="card.label" class="stat-card">
          <span class="stat-label">{{ card.label }}</span>
          <strong class="stat-value stat-value--small">{{ card.value }}</strong>
        </div>
      </div>

      <div class="dashboard-columns">
        <section class="quick-actions">
          <h3>Hoạt động nhanh</h3>
          <div class="quick-grid">
            <RouterLink
              v-for="action in quickActions"
              :key="action.to"
              :to="action.to"
              class="quick-action"
            >
              {{ action.label }}
            </RouterLink>
          </div>
        </section>

        <section class="class-section">
          <div class="section-header">
            <h3>Lớp được phân công</h3>
            <RouterLink to="/teacher/classes" class="section-link">Xem tất cả</RouterLink>
          </div>

          <div v-if="classes.length === 0" class="empty-card">
            Bạn chưa được phân công lớp nào. Vui lòng liên hệ quản trị viên.
          </div>

          <div v-else class="class-grid">
            <RouterLink
              v-for="cls in classes"
              :key="cls.class_id"
              :to="`/teacher/classes/${cls.class_id}`"
              class="class-card"
            >
              <span class="class-year">{{ cls.school_year || 'Năm học hiện tại' }}</span>
              <h4>{{ cls.name }}</h4>
              <p>Lớp phụ trách chính thức. Nhấn để quản lý học sinh và điểm danh.</p>
            </RouterLink>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  background:
    radial-gradient(circle at top right, color-mix(in srgb, var(--color-primary) 12%, transparent), transparent 60%),
    var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.dashboard-hero h2,
.quick-actions h3,
.class-section h3,
.class-card h4 {
  margin: 0;
  color: var(--color-text);
}

.dashboard-hero h2 {
  font-size: clamp(1.6rem, 4vw, 2.2rem);
  font-weight: 800;
}

.hero-copy {
  margin: var(--spacing-1) 0 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-base);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-3);
}

.dashboard-columns {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-5);
}

.stat-card,
.quick-action,
.class-card,
.empty-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.stat-card,
.quick-action,
.class-card {
  color: inherit;
  text-decoration: none;
  transition:
    border-color var(--transition-fast),
    box-shadow var(--transition-fast),
    transform var(--transition-fast);
}

.stat-card:hover,
.quick-action:hover,
.class-card:hover {
  border-color: color-mix(in srgb, var(--color-primary) 35%, var(--color-border));
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.stat-card {
  min-height: 6.25rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: var(--spacing-2);
  padding: var(--spacing-3);
}

.stat-card--link {
  cursor: pointer;
}

.stat-label {
  color: var(--color-text-muted);
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
}

.stat-value {
  color: var(--color-text);
  font-size: 2rem;
  line-height: 1;
}

.stat-value--small {
  font-size: 1.6rem;
}

.quick-actions h3,
.class-section h3 {
  font-size: var(--font-size-base);
  font-weight: 800;
  margin-bottom: var(--spacing-2);
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-2);
}

.quick-action {
  min-height: 4.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3);
  text-align: center;
  font-weight: 700;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  margin-bottom: var(--spacing-2);
}

.section-link {
  color: var(--color-primary);
  font-size: var(--font-size-xs);
  font-weight: 700;
  text-decoration: none;
}

.stat-card:focus-visible,
.quick-action:focus-visible,
.class-card:focus-visible,
.section-link:focus-visible {
  box-shadow: 0 0 0 3px var(--color-primary-focus-ring);
}

.class-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-3);
}

.class-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
  padding: var(--spacing-3);
}

.class-card p,
.empty-card {
  margin: 0;
  color: var(--color-text-muted);
}

.class-year {
  width: fit-content;
  border-radius: var(--radius-full);
  background: color-mix(in srgb, var(--color-primary) 8%, transparent);
  color: var(--color-primary);
  padding: var(--spacing-1) var(--spacing-2);
  font-size: 0.68rem;
  font-weight: 800;
  text-transform: uppercase;
}

.empty-card {
  padding: var(--spacing-4);
  text-align: center;
}

@media (min-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(3, 1fr);
  }

  .quick-grid {
    grid-template-columns: repeat(3, 1fr);
  }

  .class-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .dashboard-hero {
    align-items: flex-end;
  }

  .quick-actions h3,
  .class-section h3 {
    font-size: var(--font-size-lg);
  }
}

@media (min-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
  }

  .dashboard-columns {
    grid-template-columns: minmax(15rem, 18rem) minmax(0, 1fr);
  }

  .class-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .dashboard-hero {
    flex-direction: column;
  }
}

@media (max-width: 479px) {
  .quick-grid {
    grid-template-columns: 1fr;
  }
}
</style>
