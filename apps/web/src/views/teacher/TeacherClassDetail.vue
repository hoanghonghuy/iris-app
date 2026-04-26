<script setup>
import { computed, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeft, Calendar, GraduationCap, LoaderCircle, User } from 'lucide-vue-next'
import { useTeacherClassSelection } from '../../composables/teacher'
import { formatDateVN } from '../../helpers/dateFormatter'

const route = useRoute()

const {
  classes,
  selectedClassId,
  students,
  isLoadingClasses,
  isLoadingStudents,
  errorMessage,
  fetchClasses,
} = useTeacherClassSelection()

const genderLabel = {
  male: 'Nam',
  female: 'Nữ',
  other: 'Khác',
}

const classId = computed(() => String(route.params.classId || ''))
const selectedClass = computed(() => {
  return classes.value.find((classInfo) => String(classInfo.class_id) === classId.value) || null
})

const loading = computed(() => isLoadingClasses.value || isLoadingStudents.value)

// Sync route param with selectedClassId
watch(
  classId,
  (newClassId) => {
    if (newClassId && newClassId !== selectedClassId.value) {
      selectedClassId.value = newClassId
    }
  },
  { immediate: true },
)

// Load classes on mount
watch(
  isLoadingClasses,
  (isLoading) => {
    if (!isLoading && classId.value && !selectedClass.value) {
      errorMessage.value = 'Bạn không được phân công lớp này hoặc lớp không còn tồn tại.'
    }
  },
  { immediate: true },
)

fetchClasses()
</script>

<template>
  <div class="teacher-class-detail page-stack">
    <div class="detail-head">
      <div class="head-copy">
        <RouterLink to="/teacher/classes" class="back-link">
          <ArrowLeft :size="16" />
          Quay lại danh sách lớp
        </RouterLink>

        <h2>{{ selectedClass?.name || 'Chi tiết lớp học' }}</h2>
        <p class="text-muted">
          {{ selectedClass?.school_year || 'Năm học hiện tại' }}
        </p>
      </div>

      <RouterLink v-if="selectedClass" to="/teacher/attendance" class="btn btn--outline">
        Điểm danh lớp này
      </RouterLink>
    </div>

    <div v-if="errorMessage" class="alert alert--error">
      {{ errorMessage }}
    </div>

    <div v-if="loading" class="loading-block">
      <LoaderCircle class="spin text-muted" :size="32" />
    </div>

    <template v-else>
      <section v-if="selectedClass" class="summary-grid">
        <div class="card summary-card">
          <span class="summary-label">Lớp</span>
          <strong class="summary-value">{{ selectedClass.name }}</strong>
        </div>

        <div class="card summary-card">
          <span class="summary-label">Năm học</span>
          <strong class="summary-value">{{ selectedClass.school_year || 'Đang cập nhật' }}</strong>
        </div>

        <div class="card summary-card">
          <span class="summary-label">Số học sinh</span>
          <strong class="summary-value">{{ students.length }}</strong>
        </div>
      </section>

      <div v-if="selectedClass && students.length === 0" class="card empty-card">
        <GraduationCap :size="44" class="text-muted" />
        <p>Chưa có học sinh nào trong lớp này.</p>
      </div>

      <div v-else-if="selectedClass" class="content-stack">
        <div class="card table-card desktop-only">
          <div class="table-responsive">
            <table class="table">
              <thead>
                <tr>
                  <th>Họ tên</th>
                  <th>Ngày sinh</th>
                  <th>Giới tính</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="student in students" :key="student.student_id">
                  <td class="font-medium">{{ student.full_name }}</td>
                  <td class="text-muted">{{ formatDateVN(student.dob) }}</td>
                  <td class="text-muted">{{ genderLabel[student.gender] || student.gender }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="mobile-list">
          <article v-for="student in students" :key="student.student_id" class="card student-card">
            <User class="text-muted shrink-0" :size="20" />
            <div class="student-copy">
              <p class="font-medium m-0">{{ student.full_name }}</p>
              <p class="student-meta">
                <Calendar :size="12" />
                {{ formatDateVN(student.dob) }} ·
                {{ genderLabel[student.gender] || student.gender }}
              </p>
            </div>
          </article>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.page-stack,
.content-stack,
.mobile-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.detail-head {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.head-copy h2,
.head-copy p,
.empty-card p,
.student-copy p {
  margin: 0;
}

.back-link {
  width: fit-content;
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  color: var(--color-primary);
  text-decoration: none;
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.head-copy h2 {
  margin-top: var(--spacing-1);
  font-size: clamp(1.4rem, 4vw, 1.8rem);
  color: var(--color-text);
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--spacing-3);
}

.summary-card {
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.summary-label {
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
  font-weight: 700;
  text-transform: uppercase;
}

.summary-value {
  color: var(--color-text);
  font-size: var(--font-size-xl);
  line-height: 1.1;
}

.loading-block,
.empty-card {
  display: flex;
  justify-content: center;
}

.loading-block {
  padding: 3rem 0;
}

.empty-card {
  align-items: center;
  flex-direction: column;
  gap: var(--spacing-3);
  padding: 3rem var(--spacing-4);
  color: var(--color-text-muted);
  text-align: center;
}

.table-card {
  padding: 0;
}

.mobile-list {
  gap: var(--spacing-3);
}

.student-card {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
}

.student-copy {
  min-width: 0;
}

.student-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  margin-top: var(--spacing-1);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.desktop-only {
  display: none;
}

.spin {
  animation: spin 1s linear infinite;
}

@media (min-width: 768px) {
  .detail-head {
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-end;
  }

  .summary-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .desktop-only {
    display: block;
  }

  .mobile-list {
    display: none;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
