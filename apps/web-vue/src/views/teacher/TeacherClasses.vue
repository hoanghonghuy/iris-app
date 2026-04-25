<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { GraduationCap, LoaderCircle, User, Calendar, ChevronDown } from 'lucide-vue-next'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateVN } from '../../helpers/dateFormatter'

const classes = ref([])
const selectedClassId = ref('')
const students = ref([])
const loadingClasses = ref(true)
const loadingStudents = ref(false)
const errorMessage = ref('')

const genderLabel = {
  male: 'Nam',
  female: 'Nữ',
  other: 'Khác',
}

const selectedClassName = computed(() => {
  return classes.value.find((classInfo) => classInfo.class_id === selectedClassId.value)?.name || ''
})

async function fetchClasses() {
  loadingClasses.value = true
  errorMessage.value = ''
  try {
    const response = await teacherService.getMyClasses()
    classes.value = response?.data ?? []
    if (classes.value.length > 0) {
      selectedClassId.value = classes.value[0].class_id
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách lớp'
  } finally {
    loadingClasses.value = false
  }
}

async function fetchStudents() {
  if (!selectedClassId.value) return
  loadingStudents.value = true
  errorMessage.value = ''
  try {
    const response = await teacherService.getStudentsInClass(selectedClassId.value)
    students.value = response?.data ?? []
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách học sinh'
  } finally {
    loadingStudents.value = false
  }
}

watch(selectedClassId, fetchStudents)

onMounted(fetchClasses)
</script>

<template>
  <div class="teacher-classes page-stack">
    <div class="toolbar">
      <div v-if="classes.length > 0" class="select-wrap">
        <select v-model="selectedClassId" class="form-input class-select">
          <option v-for="classInfo in classes" :key="classInfo.class_id" :value="classInfo.class_id">
            {{ classInfo.name }} ({{ classInfo.school_year }})
          </option>
        </select>
        <ChevronDown class="select-icon" :size="16" />
      </div>
    </div>

    <div v-if="errorMessage" class="alert alert--error">{{ errorMessage }}</div>

    <div v-if="loadingClasses || loadingStudents" class="loading-block">
      <LoaderCircle class="spin text-muted" :size="32" />
    </div>

    <div v-else-if="classes.length === 0" class="card empty-card">
      <GraduationCap :size="48" class="text-muted" />
      <p>Bạn chưa được phân công lớp nào</p>
    </div>

    <div v-else-if="students.length === 0" class="card empty-card">
      <User :size="48" class="text-muted" />
      <p>Chưa có học sinh nào trong {{ selectedClassName }}</p>
    </div>

    <div v-else class="content-stack">
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
          <div>
            <p class="font-medium m-0">{{ student.full_name }}</p>
            <p class="student-meta">
              <Calendar :size="12" />
              {{ formatDateVN(student.dob) }} · {{ genderLabel[student.gender] || student.gender }}
            </p>
          </div>
        </article>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-stack,
.content-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-6);
}

.toolbar {
  display: flex;
  justify-content: flex-end;
}

.select-wrap {
  position: relative;
  width: 100%;
  max-width: 18rem;
}

.class-select {
  appearance: none;
  padding-right: 2rem;
}

.select-icon {
  pointer-events: none;
  position: absolute;
  top: 50%;
  right: var(--spacing-2);
  transform: translateY(-50%);
  color: var(--color-text-muted);
}

.loading-block,
.empty-card {
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-block {
  padding: 3rem 0;
}

.empty-card {
  flex-direction: column;
  gap: var(--spacing-4);
  padding: 3rem var(--spacing-4);
  color: var(--color-text-muted);
}

.table-card {
  padding: 0;
}

.mobile-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.student-card {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
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
