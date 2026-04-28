<script setup>
import { ref } from 'vue'
import { KeyRound, LoaderCircle } from 'lucide-vue-next'
import { authService } from '../services/authService'
import { extractErrorMessage } from '../helpers/errorHandler'

const password = ref('')
const confirmPassword = ref('')
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

async function handleSubmit() {
  errorMessage.value = ''
  successMessage.value = ''

  if (password.value.length < 6) {
    errorMessage.value = 'Mật khẩu tối thiểu 6 ký tự'
    return
  }

  if (password.value !== confirmPassword.value) {
    errorMessage.value = 'Mật khẩu xác nhận không khớp'
    return
  }

  submitting.value = true
  try {
    await authService.updateMyPassword(password.value)
    successMessage.value = 'Đổi mật khẩu thành công!'
    password.value = ''
    confirmPassword.value = ''
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="card profile-card">
    <div class="card__header">
      <h2 class="card-title">
        <KeyRound :size="20" />
        Đổi mật khẩu
      </h2>
    </div>
    <div class="card__body">
      <form class="flex-col gap-4" @submit.prevent="handleSubmit">
        <div v-if="errorMessage" class="alert alert--error">{{ errorMessage }}</div>
        <div v-if="successMessage" class="alert alert--success">{{ successMessage }}</div>

        <div class="form-group mb-0">
          <label class="form-label" for="newPassword">Mật khẩu mới</label>
          <input
            id="newPassword"
            v-model="password"
            class="form-input"
            type="password"
            placeholder="Tối thiểu 6 ký tự"
            required
          />
        </div>

        <div class="form-group mb-0">
          <label class="form-label" for="confirmNewPassword">Xác nhận mật khẩu</label>
          <input
            id="confirmNewPassword"
            v-model="confirmPassword"
            class="form-input"
            type="password"
            placeholder="Nhập lại mật khẩu mới"
            required
          />
        </div>

        <button type="submit" class="btn btn--primary fit-content" :disabled="submitting">
          <LoaderCircle v-if="submitting" class="spin mr-2" :size="16" />
          Cập nhật mật khẩu
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.profile-card {
  max-width: 32rem;
}

.card-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-lg);
  margin: 0;
}
</style>
