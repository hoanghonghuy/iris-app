<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeft, Calendar, LoaderCircle } from 'lucide-vue-next'
import { parentService } from '../../services/parentService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateVN } from '../../helpers/dateFormatter'

const route = useRoute()
const children = ref([])
const loading = ref(true)
const errorMessage = ref('')

const genderLabel = {
  male: 'Nam',
  female: 'Nữ',
  other: 'Khác',
}

const child = computed(() => {
  return children.value.find((item) => item.student_id === route.params.studentId) || null
})

const childAge = computed(() => {
  if (!child.value?.dob) return null
  const match = child.value.dob.match(/^(\d{4})-(\d{2})-(\d{2})/)
  if (!match) return null

  const birthDate = new Date(Number(match[1]), Number(match[2]) - 1, Number(match[3]))
  if (Number.isNaN(birthDate.getTime())) return null

  const today = new Date()
  let age = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()
  const dayDiff = today.getDate() - birthDate.getDate()
  if (monthDiff < 0 || (monthDiff === 0 && dayDiff < 0)) age -= 1
  return age
})

async function loadChild() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await parentService.getMyChildren()
    children.value = response?.data ?? []
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải thông tin con'
  } finally {
    loading.value = false
  }
}

onMounted(loadChild)
</script>

<template>
  <div class="child-detail page-stack">
    <RouterLink to="/parent/children" class="back-link">
      <ArrowLeft :size="16" />
      Quay lại danh sách con
    </RouterLink>

    <div v-if="errorMessage" class="alert alert--error">{{ errorMessage }}</div>

    <div v-if="loading" class="loading-block">
      <LoaderCircle class="spin text-muted" :size="32" />
    </div>

    <div v-else-if="!child" class="card">
      <div class="card__body text-sm text-muted">Không tìm thấy thông tin của bé này.</div>
    </div>

    <div v-else class="card">
      <div class="card__header">
        <h2 class="title">{{ child.full_name }}</h2>
      </div>
      <div class="card__body info-stack">
        <p>
          <Calendar :size="16" />
          Ngày sinh: {{ formatDateVN(child.dob) }}
        </p>
        <p>Giới tính: {{ genderLabel[child.gender] || child.gender }}</p>
        <p v-if="childAge !== null">{{ childAge }} tuổi</p>
        <p>Lớp hiện tại: {{ child.current_class_name || 'Chưa có thông tin' }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-stack,
.info-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  width: fit-content;
}

.back-link:hover {
  color: var(--color-text);
}

.loading-block {
  display: flex;
  justify-content: center;
  padding: 3rem 0;
}

.title {
  margin: 0;
  font-size: var(--font-size-lg);
}

.info-stack {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.info-stack p {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  margin: 0;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
