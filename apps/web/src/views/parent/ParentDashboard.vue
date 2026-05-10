<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { parentService } from '../../services/parentService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { POST_TYPE_META } from '../../helpers/postConfig'
import { formatDateTime } from '@/helpers/dateFormatter'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import AnalyticsTimeseriesPanel from '../../components/charts/AnalyticsTimeseriesPanel.vue'

const authStore = useAuthStore()
const analytics = ref(null)
const posts = ref([])
const children = ref([])
const chartStudentId = ref('')
const isLoading = ref(true)
const errorMessage = ref('')
const timeseries = ref(null)
const chartsLoading = ref(false)
const chartsError = ref('')
const selectedRange = ref('14d')
const activeTab = ref('overview')
const chartChildSelectId = 'parent-chart-child-select'
const RANGE_OPTIONS = [
  { value: '7d', label: '7 ngày' },
  { value: '14d', label: '14 ngày' },
  { value: '30d', label: '30 ngày' },
]
const PARENT_SERIES_ORDER = [
  'child_attendance',
  'health_alerts',
  'appointments_by_status',
]

const displayName = computed(() => {
  const user = authStore.currentUser
  return user?.full_name || user?.email?.split('@')[0] || 'Phụ huynh'
})

const analyticsView = computed(() => {
  const data = analytics.value || {}
  return {
    total_children: data.total_children ?? 0,
    recent_health_alerts_24h: data.recent_health_alerts_24h ?? 0,
    today_attendance_present_count: data.today_attendance_present_count ?? 0,
    today_attendance_pending_count: data.today_attendance_pending_count ?? 0,
  }
})

const quickActions = [
  { label: 'Hồ sơ con', to: '/parent/children' },
  { label: 'Bảng tin lớp', to: '/parent/posts' },
]

function getPostMeta(type) {
  return POST_TYPE_META[type] || { label: type || 'Bài đăng', badgeClass: 'badge--outline' }
}

const parentTimeseriesPayload = computed(() => {
  const payload = timeseries.value
  if (!payload?.series?.length) return null

  const indexMap = new Map(payload.series.map((item) => [item.id, item]))
  const series = PARENT_SERIES_ORDER
    .map((id) => indexMap.get(id))
    .filter(Boolean)
    .slice(0, 3)

  if (!series.length) return null
  return { ...payload, series }
})

async function fetchTimeseries() {
  if (!chartStudentId.value) {
    timeseries.value = null
    return
  }
  chartsLoading.value = true
  chartsError.value = ''
  try {
    const res = await parentService.getAnalyticsTimeseries({
      range: selectedRange.value,
      interval: 'day',
      student_id: chartStudentId.value,
    })
    timeseries.value = res?.data ?? res ?? null
  } catch (error) {
    chartsError.value = extractErrorMessage(error)
    timeseries.value = null
  } finally {
    chartsLoading.value = false
  }
}

async function fetchDashboard() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const [analyticsRes, feedRes, childrenRes] = await Promise.all([
      parentService.getAnalytics(),
      parentService.getMyFeed({ limit: 5 }),
      parentService.getMyChildren(),
    ])

    analytics.value = analyticsRes?.data ?? analyticsRes ?? {}
    posts.value = normalizeListResponse(feedRes)
    children.value = normalizeListResponse(childrenRes)
    if (!chartStudentId.value && children.value.length) {
      chartStudentId.value = children.value[0].student_id
    }
    await fetchTimeseries()
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

function changeRange(rangeValue) {
  if (selectedRange.value === rangeValue || chartsLoading.value) return
  selectedRange.value = rangeValue
  fetchTimeseries()
}
</script>

<template>
  <div class="parent-dashboard">
    <div class="dashboard-hero mb-6">
      <div>
        <h2>Xin chào, {{ displayName }}</h2>
        <p class="hero-copy">Hôm nay con bạn có hoạt động gì mới?</p>
      </div>
      <button
        class="btn btn--outline"
        type="button"
        @click="fetchDashboard"
        :disabled="isLoading || chartsLoading"
      >
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

    <LoadingSpinner v-if="isLoading" message="Đang tải thông tin..." />

    <div
      v-else-if="errorMessage"
      class="alert alert--error"
    >
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button type="button" class="btn btn--outline mt-2" @click="fetchDashboard">Thử lại</button>
    </div>

    <div v-else class="dashboard-content">
      <!-- Tab Tổng quan -->
      <div v-show="activeTab === 'overview'">
        <div class="stats-grid mb-6">
          <div class="stat-card">
            <span class="stat-label">Tổng số con</span>
            <strong class="stat-value">{{ analyticsView.total_children }}</strong>
          </div>
          <div class="stat-card">
            <span class="stat-label">Cảnh báo sức khỏe 24h</span>
            <strong class="stat-value">{{ analyticsView.recent_health_alerts_24h }}</strong>
          </div>
          <div class="stat-card">
            <span class="stat-label">Con có mặt hôm nay</span>
            <strong class="stat-value stat-value--small">{{
              analyticsView.today_attendance_present_count
            }}</strong>
          </div>
          <div class="stat-card">
            <span class="stat-label">Con chưa điểm danh</span>
            <strong class="stat-value stat-value--small">{{
              analyticsView.today_attendance_pending_count
            }}</strong>
          </div>
        </div>

        <section class="quick-actions mb-8">
          <h3>Truy cập nhanh</h3>
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

        <section class="feed-section">
          <div class="section-header">
            <h3>Hoạt động mới nhất</h3>
            <RouterLink to="/parent/posts" class="section-link">Xem tất cả</RouterLink>
          </div>

          <div v-if="posts.length === 0" class="empty-card">
            Chưa có bài viết hay thông báo mới nào được đăng tải.
          </div>

          <div v-else class="feed-list">
            <article v-for="post in posts" :key="post.post_id" class="feed-card">
              <div class="feed-meta">
                <span class="badge" :class="getPostMeta(post.type).badgeClass">
                  {{ getPostMeta(post.type).label }}
                </span>
                <span class="feed-date">{{ formatDateTime(post.created_at) }}</span>
              </div>
              <p class="feed-content">{{ post.content }}</p>
            </article>
          </div>
        </section>
      </div>

      <!-- Tab Biểu đồ -->
      <div v-show="activeTab === 'charts'">
        <div v-if="children.length > 1" class="card chart-child-filter mb-6">
          <label :for="chartChildSelectId" class="form-label mb-1">Biểu đồ theo con</label>
          <select :id="chartChildSelectId" v-model="chartStudentId" class="form-input" @change="fetchTimeseries">
            <option v-for="c in children" :key="c.student_id" :value="c.student_id">
              {{ c.full_name }}
            </option>
          </select>
        </div>
        <div v-else-if="children.length === 0" class="empty-card mb-6">
          Chưa có học sinh liên kết với tài khoản này nên chưa hiển thị biểu đồ.
        </div>

        <section v-if="children.length > 0" class="card range-filter mb-6">
          <div class="range-filter__heading">
            <h3>Khoảng thời gian biểu đồ</h3>
            <span class="text-muted text-sm">Theo dõi nhanh diễn biến gần đây</span>
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
          v-else-if="parentTimeseriesPayload && chartStudentId"
          :payload="parentTimeseriesPayload"
          :title="`Xu hướng của con trong ${selectedRange.replace('d', ' ngày')} gần nhất`"
        />
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
    radial-gradient(circle at top right, color-mix(in srgb, var(--color-primary) 12%, transparent), transparent 58%),
    var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.dashboard-hero h2,
.quick-actions h3,
.feed-section h3 {
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

.stat-card,
.quick-action,
.feed-card,
.empty-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.stat-card {
  min-height: 6rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: var(--spacing-2);
  padding: var(--spacing-3);
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
.feed-section h3 {
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
  color: inherit;
  text-align: center;
  text-decoration: none;
  font-weight: 700;
  transition:
    border-color 0.2s,
    box-shadow 0.2s,
    transform 0.2s;
}

.quick-action:hover {
  border-color: color-mix(in srgb, var(--color-primary) 35%, var(--color-border));
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.quick-action:focus-visible,
.section-link:focus-visible {
  box-shadow: 0 0 0 3px var(--color-primary-focus-ring);
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

.feed-list {
  display: grid;
  gap: var(--spacing-4);
}

.feed-card {
  padding: var(--spacing-3);
}

.feed-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  margin-bottom: var(--spacing-3);
}

.feed-date {
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.feed-content,
.empty-card {
  margin: 0;
  color: var(--color-text);
  white-space: pre-line;
}

.empty-card {
  padding: var(--spacing-4);
  text-align: center;
  color: var(--color-text-muted);
}

.chart-child-filter {
  padding: var(--spacing-3);
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

@media (min-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
  }

  .quick-actions h3,
  .feed-section h3 {
    font-size: var(--font-size-lg);
  }
}

@media (max-width: 767px) {
  .dashboard-hero,
  .section-header {
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
}

@media (max-width: 479px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .quick-grid {
    grid-template-columns: 1fr;
  }

  .range-filter__actions .btn {
    flex: 1 1 100%;
  }
}
</style>
