<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDate } from '../../helpers/dateFormatter'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'

const route = useRoute()
const router = useRouter()
const studentId = route.params.id

const profile = ref(null)
const isLoading = ref(true)
const errorMessage = ref('')

const genderLabel = computed(() => {
  const gender = String(profile.value?.gender || '').toLowerCase()
  if (gender === 'female') {
    return 'Nữ'
  }

  if (gender === 'other') {
    return 'Khác'
  }

  return 'Nam'
})

const ageLabel = computed(() => {
  if (!profile.value?.dob) {
    return ''
  }

  const dob = new Date(profile.value.dob)
  if (Number.isNaN(dob.getTime())) {
    return ''
  }

  const today = new Date()
  let years = today.getFullYear() - dob.getFullYear()
  const monthDiff = today.getMonth() - dob.getMonth()

  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < dob.getDate())) {
    years -= 1
  }

  return `${years} tuổi`
})

async function fetchStudentProfile() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const response = await adminService.getStudentProfile(studentId)
    profile.value = response?.data || null
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải hồ sơ học sinh'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchStudentProfile()
})
</script>

<template>
  <div class="student-detail page-stack">
    <div>
      <button class="btn btn--outline btn--sm" type="button" @click="router.back()">
        ← Quay lại
      </button>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold">Lỗi tải dữ liệu</p>
      <p>{{ errorMessage }}</p>
      <button class="btn btn--outline mt-2" type="button" @click="fetchStudentProfile">
        Thử lại
      </button>
    </div>

    <LoadingSpinner v-else-if="isLoading" message="Đang tải hồ sơ học sinh..." />

    <div v-else-if="!profile" class="card">
      <EmptyState
        title="Không tìm thấy học sinh"
        message="Bản ghi học sinh này không còn tồn tại hoặc bạn không có quyền truy cập."
      />
    </div>

    <div v-else class="detail-grid">
      <div class="card panel">
        <div class="hero">
          <div class="avatar">
            {{ profile.full_name?.charAt(0)?.toUpperCase() }}
          </div>
          <div>
            <h2 class="hero__title">{{ profile.full_name }}</h2>
            <p class="hero__subtitle">
              {{ profile.current_class_name || 'Chưa xếp lớp' }}
            </p>
          </div>
        </div>

        <div class="info-list">
          <div class="info-item">
            <span class="info-label">Ngày sinh</span>
            <span class="info-value">{{ formatDate(profile.dob) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Độ tuổi</span>
            <span class="info-value">{{ ageLabel || 'Không rõ' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Giới tính</span>
            <span class="info-value">{{ genderLabel }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Mã học sinh</span>
            <span class="info-value">{{ profile.student_id }}</span>
          </div>
        </div>
      </div>

      <div class="card panel">
        <div class="section-head">
          <h3 class="section-title">Phụ huynh liên kết</h3>
          <span class="section-badge">{{ profile.parents?.length || 0 }}</span>
        </div>

        <div v-if="profile.parents && profile.parents.length > 0" class="parent-list">
          <article v-for="parent in profile.parents" :key="parent.parent_id" class="parent-card">
            <div class="parent-card__head">
              <p class="parent-card__name">{{ parent.full_name }}</p>
              <span class="badge badge--outline">Phụ huynh</span>
            </div>
            <p class="parent-card__meta">{{ parent.email || 'Chưa có email' }}</p>
            <p class="parent-card__meta">{{ parent.phone || 'Chưa có số điện thoại' }}</p>
          </article>
        </div>

        <EmptyState
          v-else
          title="Chưa có phụ huynh liên kết"
          message="Học sinh này hiện chưa được liên kết với phụ huynh nào."
        />
      </div>

      <div class="card panel detail-grid__wide">
        <div class="section-head">
          <h3 class="section-title">Mã phụ huynh</h3>
        </div>

        <div v-if="profile.active_parent_code" class="parent-code">
          <code class="parent-code__value">{{ profile.active_parent_code }}</code>
          <p class="parent-code__meta">
            Hết hạn: {{ formatDate(profile.code_expires_at) || 'Không rõ' }}
          </p>
        </div>

        <p v-else class="parent-code__meta">
          Hồ sơ chi tiết hiện không có mã phụ huynh đang hoạt động. Bạn có thể tạo hoặc thu hồi mã trực tiếp ở màn danh sách học sinh.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.detail-grid {
  display: grid;
  gap: var(--spacing-4);
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.panel {
  padding: var(--spacing-5);
}

.hero {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
  margin-bottom: var(--spacing-5);
}

.avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 4rem;
  height: 4rem;
  border-radius: 999px;
  background: color-mix(in srgb, var(--color-primary) 14%, var(--color-on-primary));
  color: var(--color-primary);
  font-size: 1.5rem;
  font-weight: 700;
}

.hero__title,
.hero__subtitle,
.section-title,
.parent-card__name,
.parent-card__meta,
.parent-code__meta {
  margin: 0;
}

.hero__title,
.section-title,
.parent-card__name {
  color: var(--color-text);
  font-weight: 700;
}

.hero__subtitle,
.parent-card__meta,
.parent-code__meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  margin-bottom: var(--spacing-4);
}

.section-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2rem;
  height: 2rem;
  padding: 0 0.6rem;
  border-radius: 999px;
  background: var(--color-background);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.info-list,
.parent-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  padding-bottom: var(--spacing-2);
  border-bottom: 1px solid var(--color-border);
}

.info-item:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.info-label {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.info-value {
  color: var(--color-text);
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-align: right;
  word-break: break-word;
}

.parent-card {
  padding: var(--spacing-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
}

.parent-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  margin-bottom: var(--spacing-2);
}

.parent-code__value {
  display: inline-flex;
  align-items: center;
  min-height: 2.25rem;
  padding: 0.35rem 0.75rem;
  border-radius: var(--radius-md);
  background: var(--color-background);
  border: 1px solid var(--color-border);
  font-size: 1rem;
}

.mt-2 {
  margin-top: var(--spacing-2);
}

@media (min-width: 1024px) {
  .detail-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .detail-grid__wide {
    grid-column: 1 / -1;
  }
}
</style>
