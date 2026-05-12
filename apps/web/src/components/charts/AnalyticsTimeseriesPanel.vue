<script setup lang="ts">
import { computed } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  Filler,
} from 'chart.js'
import { Line, Bar, Doughnut } from 'vue-chartjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  Filler,
)

const props = defineProps({
  payload: {
    type: Object,
    default: null,
  },
  title: {
    type: String,
    default: 'Xu hướng',
  },
})

const COLORS = ['#6366f1', '#8b5cf6', '#94a3b8', '#f59e0b', '#10b981', '#ef4444']
const SERIES_LABEL_MAP = {
  present_rate: 'Tỷ lệ có mặt',
  attendance_status: 'Điểm danh theo trạng thái',
  health_alerts: 'Cảnh báo sức khỏe',
  appointments_by_status: 'Lịch hẹn theo trạng thái',
  attendance_marked: 'Lượt điểm danh đã ghi nhận',
  attendance_marked_vs_pending: 'Điểm danh đã ghi nhận và chưa ghi nhận',
  population_by_role: 'Phân bố theo vai trò',
  present: 'Có mặt',
  absent: 'Vắng mặt',
  late: 'Đi muộn',
  excused: 'Có phép',
  pending: 'Chờ xác nhận',
  confirmed: 'Đã xác nhận',
  completed: 'Hoàn thành',
  cancelled: 'Đã hủy',
  no_show: 'Không đến',
  marked: 'Đã điểm danh',
  not_marked: 'Chưa điểm danh',
}

function prettifyKey(key) {
  if (!key) return ''
  // Fallback: chỉ capitalize chữ cái đầu tiên của cả chuỗi, không capitalize mỗi từ
  const normalized = key.replaceAll('_', ' ')
  return normalized.charAt(0).toUpperCase() + normalized.slice(1)
}

function normalizeSeriesLabel(labelOrId) {
  return SERIES_LABEL_MAP[labelOrId] || prettifyKey(labelOrId)
}

function formatMetricValue(series) {
  if (!series?.points?.length) return '--'
  const lastPoint = series.points[series.points.length - 1]
  
  // Nếu series có components (stacked data), tính tổng
  if (lastPoint?.components && typeof lastPoint.components === 'object') {
    const total = Object.values(lastPoint.components).reduce((sum, val) => sum + (Number(val) || 0), 0)
    if (series.unit === 'percent') return `${total.toFixed(1)}%`
    return `${Math.round(total)}`
  }
  
  // Series có value đơn giản
  const raw = Number(lastPoint?.value ?? 0)
  if (!Number.isFinite(raw)) return '--'
  if (series.unit === 'percent') return `${raw.toFixed(1)}%`
  return `${Math.round(raw)}`
}

function formatBucket(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return `${d.getUTCDate().toString().padStart(2, '0')}/${(d.getUTCMonth() + 1).toString().padStart(2, '0')}`
}

function localizedRoleLabel(roleKey) {
  if (roleKey === 'teacher') return 'Giáo viên'
  if (roleKey === 'parent') return 'Phụ huynh'
  if (roleKey === 'student') return 'Học sinh'
  return roleKey
}

const gridItems = computed(() => {
  const p = props.payload
  if (!p?.series?.length) return []

  return p.series.map((series) => {
    const pts = Array.isArray(series.points) ? series.points : []

    if (series.id === 'population_by_role' && pts.length === 1 && pts[0].components) {
      const comp = pts[0].components
      const labels = Object.keys(comp).map((k) => localizedRoleLabel(k))
      const data = Object.values(comp)
      return {
        key: series.id,
        kind: 'doughnut',
        title: normalizeSeriesLabel(series.label || series.id),
        component: 'Doughnut',
        chartData: {
          labels,
          datasets: [
            {
              data,
              backgroundColor: COLORS.slice(0, labels.length),
              borderWidth: 0,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: { position: 'bottom' },
          },
        },
      }
    }

    const first = pts[0]
    const hasComponents = Boolean(
      first?.components &&
      typeof first.components === 'object' &&
      Object.keys(first.components).length > 0,
    )

    if (hasComponents) {
      const labels = pts.map((x) => formatBucket(x.bucket_start))
      const keys = Object.keys(first.components)
      const datasets = keys.map((key, idx) => ({
        label: normalizeSeriesLabel(key),
        data: pts.map((pt) => pt.components[key] ?? 0),
        backgroundColor: COLORS[idx % COLORS.length],
        stack: 'a',
      }))
      return {
        key: series.id,
        kind: 'bar',
        title: normalizeSeriesLabel(series.label || series.id),
        component: 'Bar',
        chartData: { labels, datasets },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: { position: 'bottom' },
          },
          scales: {
            x: { stacked: true },
            y: { stacked: true, beginAtZero: true },
          },
        },
      }
    }

    const labels = pts.map((x) => formatBucket(x.bucket_start))
    const data = pts.map((x) => (x.value == null ? null : Number(x.value)))
    return {
      key: series.id,
      kind: 'line',
      title: normalizeSeriesLabel(series.label || series.id),
      component: 'Line',
      chartData: {
        labels,
        datasets: [
          {
            label: normalizeSeriesLabel(series.label || series.id),
            data,
            borderColor: '#6366f1',
            backgroundColor: 'rgba(99, 102, 241, 0.15)',
            fill: true,
            tension: 0.25,
            spanGaps: true,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
        },
        scales: {
          y: {
            beginAtZero: true,
            ...(series.unit === 'percent' ? { max: 100 } : {}),
          },
        },
      },
    }
  })
})

const summaryItems = computed(() => {
  const p = props.payload
  if (!p?.series?.length) return []
  return p.series
    .map((series) => ({
      key: series.id,
      label: normalizeSeriesLabel(series.label || series.id),
      value: formatMetricValue(series),
    }))
    .slice(0, 4)
})

function shouldSpanFull(index) {
  const total = gridItems.value.length
  // Nếu có số lẻ biểu đồ (3, 5, 7...) và đây là item cuối cùng, span full width
  return total % 2 === 1 && index === total - 1
}
</script>

<template>
  <section v-if="payload?.series?.length" class="analytics-ts-panel">
    <h3 class="analytics-ts-panel__title">{{ title }}</h3>
    <p class="analytics-ts-panel__meta text-muted text-sm">
      Khoảng: {{ payload.meta?.range || '14d' }} · Bucket: {{ payload.meta?.interval || 'day' }} (UTC)
    </p>
    <div v-if="summaryItems.length" class="analytics-ts-panel__summary">
      <div
        v-for="item in summaryItems"
        :key="item.key"
        class="analytics-ts-panel__summary-item"
      >
        <span class="analytics-ts-panel__summary-label">{{ item.label }}</span>
        <strong class="analytics-ts-panel__summary-value">{{ item.value }}</strong>
      </div>
    </div>
    <div class="analytics-ts-panel__grid">
      <div
        v-for="(item, index) in gridItems"
        :key="item.key"
        class="card analytics-ts-panel__card"
        :class="{ 'analytics-ts-panel__card--span-full': shouldSpanFull(index) }"
      >
        <p class="analytics-ts-panel__chart-title">{{ item.title }}</p>
        <div class="analytics-ts-panel__chart-wrap">
          <Line
            v-if="item.component === 'Line'"
            :data="item.chartData"
            :options="item.options"
          />
          <Bar
            v-else-if="item.component === 'Bar'"
            :data="item.chartData"
            :options="item.options"
          />
          <Doughnut
            v-else-if="item.component === 'Doughnut'"
            :data="item.chartData"
            :options="item.options"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.analytics-ts-panel__title {
  margin: 0 0 var(--spacing-1);
  font-size: var(--font-size-lg);
  font-weight: 800;
  color: var(--color-text);
}

.analytics-ts-panel__meta {
  margin: 0 0 var(--spacing-4);
}

.analytics-ts-panel__summary {
  display: grid;
  gap: var(--spacing-2);
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin-bottom: var(--spacing-4);
}

.analytics-ts-panel__summary-item {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  padding: var(--spacing-2) var(--spacing-3);
}

.analytics-ts-panel__summary-label {
  display: block;
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
  margin-bottom: var(--spacing-1);
}

.analytics-ts-panel__summary-value {
  color: var(--color-text);
  font-size: var(--font-size-base);
}

.analytics-ts-panel__grid {
  display: grid;
  gap: var(--spacing-4);
  grid-template-columns: 1fr;
}

@media (min-width: 900px) {
  .analytics-ts-panel__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .analytics-ts-panel__card--span-full {
    grid-column: 1 / -1;
  }
}

.analytics-ts-panel__card {
  padding: var(--spacing-4);
}

.analytics-ts-panel__chart-title {
  margin: 0 0 var(--spacing-2);
  font-size: var(--font-size-sm);
  font-weight: 700;
  color: var(--color-text-muted);
}

.analytics-ts-panel__chart-wrap {
  position: relative;
  min-height: 260px;
  max-height: 320px;
}

@media (max-width: 767px) {
  .analytics-ts-panel__summary {
    grid-template-columns: 1fr;
  }

  .analytics-ts-panel__card {
    padding: var(--spacing-3);
  }

  .analytics-ts-panel__chart-wrap {
    min-height: 220px;
    max-height: 260px;
  }
}
</style>
