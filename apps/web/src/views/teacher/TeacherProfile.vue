<script setup>
import { ref } from 'vue'
import { LoaderCircle } from 'lucide-vue-next'
import { useAuthStore } from '../../stores/authStore'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import ChangePasswordForm from '../../components/ChangePasswordForm.vue'

const authStore = useAuthStore()
const phone = ref('')
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

async function handleSubmit() {
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await teacherService.updateMyProfile(phone.value)
    successMessage.value = 'Cập nhật thành công!'
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể cập nhật'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="profile-stack">
    <div class="card profile-card">
      <div class="card__body profile-body">
        <div class="field-readonly">
          <label>Email</label>
          <p>{{ authStore.currentUser?.email || '—' }}</p>
        </div>

        <form class="form-stack" @submit.prevent="handleSubmit">
          <div v-if="errorMessage" class="alert alert--error">{{ errorMessage }}</div>
          <div v-if="successMessage" class="alert alert--success">{{ successMessage }}</div>

          <div class="form-group mb-0">
            <label class="form-label" for="phone">Số điện thoại</label>
            <input
              id="phone"
              v-model="phone"
              class="form-input"
              type="tel"
              placeholder="0900 000 000"
            />
          </div>

          <button class="btn btn--primary fit-content" type="submit" :disabled="submitting">
            <LoaderCircle v-if="submitting" class="spin mr-2" :size="16" />
            Cập nhật
          </button>
        </form>
      </div>
    </div>

    <ChangePasswordForm />
  </div>
</template>

<style scoped>
.profile-stack,
.profile-body,
.form-stack {
  display: flex;
  flex-direction: column;
}

.profile-stack {
  gap: var(--spacing-6);
}

.profile-body,
.form-stack {
  gap: var(--spacing-4);
}

.profile-card {
  max-width: 32rem;
}

.title {
  margin: 0;
  font-size: var(--font-size-lg);
}

.field-readonly label {
  display: block;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  font-weight: 500;
  margin-bottom: var(--spacing-2);
}

.field-readonly p {
  margin: 0;
  font-size: var(--font-size-sm);
}

.fit-content {
  width: fit-content;
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
