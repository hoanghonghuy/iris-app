<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Calendar, LoaderCircle, User, Users } from 'lucide-vue-next'
import { parentService } from '../../services/parentService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateVN } from '../../helpers/dateFormatter'

const children = ref([])
const loading = ref(true)
const errorMessage = ref('')

const genderLabel = {
  male: 'Nam',
  female: 'Nữ',
  other: 'Khác',
}

async function fetchChildren() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await parentService.getMyChildren()
    children.value = response?.data ?? []
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách con'
  } finally {
    loading.value = false
  }
}

onMounted(fetchChildren)
</script>

<template>
  <div class="parent-children page-stack">
    <div v-if="errorMessage" class="alert alert--error">{{ errorMessage }}</div>

    <div v-if="loading" class="loading-block">
      <LoaderCircle class="spin text-muted" :size="32" />
    </div>

    <div v-else-if="children.length === 0" class="card empty-card">
      <Users :size="48" class="text-muted" />
      <p>Chưa có con nào được liên kết</p>
    </div>

    <div v-else class="children-grid">
      <RouterLink
        v-for="child in children"
        :key="child.student_id"
        :to="`/parent/children/${child.student_id}`"
        class="card child-card"
      >
        <User class="child-icon" :size="40" />
        <div>
          <p class="child-name">{{ child.full_name }}</p>
          <div class="child-meta">
            <p>
              <Calendar :size="12" />
              Ngày sinh: {{ formatDateVN(child.dob) }}
            </p>
            <p>Giới tính: {{ genderLabel[child.gender] || child.gender }}</p>
          </div>
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.page-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-6);
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

.empty-card p {
  margin: 0;
  font-size: var(--font-size-sm);
}

.children-grid {
  display: grid;
  gap: var(--spacing-4);
}

.child-card {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-4);
  padding: var(--spacing-5);
  color: var(--color-text);
  transition: border-color 0.2s;
}

.child-card:hover {
  border-color: color-mix(in srgb, var(--color-primary) 50%, var(--color-border));
}

.child-icon {
  flex-shrink: 0;
  border-radius: var(--radius-full);
  background: var(--color-background);
  color: var(--color-text-muted);
  padding: var(--spacing-2);
}

.child-name {
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
}

.child-meta {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  margin-top: var(--spacing-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.child-meta p {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  margin: 0;
}

.spin {
  animation: spin 1s linear infinite;
}

@media (min-width: 640px) {
  .children-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (min-width: 1024px) {
  .children-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
