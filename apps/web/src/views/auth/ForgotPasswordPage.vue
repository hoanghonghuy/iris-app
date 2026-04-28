<script setup>
import { ref } from 'vue'
import { Mail, CheckCircle2, ArrowLeft, Loader2 } from 'lucide-vue-next'
import { authService } from '../../services/authService'
import { extractErrorMessage } from '../../helpers/errorHandler'

const email = ref('')
const isLoading = ref(false)
const errorMessage = ref('')
const isSuccess = ref(false)

async function handleForgotPassword() {
  if (!email.value || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    errorMessage.value = 'Vui lòng nhập Email hợp lệ'
    return
  }

  isLoading.value = true
  errorMessage.value = ''
  isSuccess.value = false

  try {
    await authService.forgotPassword(email.value)
    isSuccess.value = true
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="auth-card">
    <div class="text-center mb-6">
      <h2 class="text-2xl font-bold mb-1">Quên mật khẩu</h2>
      <p class="text-muted text-sm">Nhập email để nhận mã đặt lại mật khẩu</p>
    </div>

    <!-- Trạng thái thành công -->
    <template v-if="isSuccess">
      <div class="flex-col items-center gap-4 py-6 text-center">
        <CheckCircle2 class="text-success mx-auto" :size="48" />
        <p class="text-sm text-muted mt-4">
          Nếu email <span class="font-medium text-foreground">{{ email }}</span> tồn tại trong hệ
          thống, bạn sẽ nhận được mã đặt lại mật khẩu trong vài phút.
        </p>
      </div>
      <div class="mt-4">
        <RouterLink
          to="/login"
          class="btn btn--outline w-full inline-flex items-center justify-center gap-2"
        >
          <ArrowLeft :size="16" /> Quay lại đăng nhập
        </RouterLink>
      </div>
    </template>

    <!-- Form nhập email -->
    <template v-else>
      <div
        v-if="errorMessage"
        class="alert alert--error"
      >
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleForgotPassword" class="flex-col gap-4">
        <div class="form-group">
          <label class="form-label" for="email">Email</label>
          <div class="input-icon-wrapper">
            <Mail class="input-icon" :size="16" />
            <input
              id="email"
              v-model="email"
              type="email"
              class="form-input form-input--icon"
              placeholder="name@example.com"
              :disabled="isLoading"
              required
            />
          </div>
        </div>

        <button type="submit" class="btn btn--primary w-full mt-2" :disabled="isLoading">
          <template v-if="isLoading">
            <Loader2 class="spin mr-2" :size="16" /> Đang gửi...
          </template>
          <template v-else> Gửi mã đặt lại mật khẩu </template>
        </button>
      </form>

      <div class="text-center text-sm mt-6">
        <RouterLink
          to="/login"
          class="text-muted hover-text-primary inline-flex items-center gap-1"
        >
          <ArrowLeft :size="12" /> Quay lại đăng nhập
        </RouterLink>
      </div>
    </template>
  </div>
</template>

<style scoped>
.input-icon-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-muted);
  pointer-events: none;
}

.form-input--icon {
  padding-left: 2.5rem;
}
</style>
