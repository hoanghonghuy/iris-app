<script setup>
import { ref } from 'vue'
import { ShieldCheck, Loader2 } from 'lucide-vue-next'
import { authService } from '../../services/authService'
import { extractErrorMessage } from '../../helpers/errorHandler'

const token = ref('')
const password = ref('')
const confirmPassword = ref('')

const isLoading = ref(false)
const errorMessage = ref('')
const isSuccess = ref(false)

async function handleActivate() {
  errorMessage.value = ''

  if (!token.value.trim()) {
    errorMessage.value = 'Token không được để trống'
    return
  }
  if (password.value.length < 6) {
    errorMessage.value = 'Mật khẩu tối thiểu 6 ký tự'
    return
  }
  if (password.value !== confirmPassword.value) {
    errorMessage.value = 'Mật khẩu xác nhận không khớp'
    return
  }

  isLoading.value = true

  try {
    await authService.activateWithToken(token.value.trim(), password.value)
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
    <!-- Trạng thái thành công -->
    <template v-if="isSuccess">
      <div class="flex-col items-center py-12 text-center">
        <ShieldCheck class="success-icon mx-auto" :size="64" />
        <h2 class="text-xl font-semibold mt-4">Kích hoạt thành công!</h2>
        <p class="text-sm text-muted mt-2">Bạn có thể đăng nhập với mật khẩu mới.</p>
        <RouterLink to="/login" class="btn btn--primary mt-6 inline-flex"> Đăng nhập </RouterLink>
      </div>
    </template>

    <!-- Form kích hoạt -->
    <template v-else>
      <div class="text-center mb-6">
        <div class="flex justify-center mb-2">
          <ShieldCheck class="text-muted" :size="40" />
        </div>
        <h2 class="text-xl font-bold mb-1">Kích hoạt tài khoản</h2>
        <p class="text-muted text-sm">Nhập token kích hoạt và đặt mật khẩu mới</p>
      </div>

      <div
        v-if="errorMessage"
        class="mb-4 p-3 bg-red-50 text-danger text-sm rounded border border-red-200"
      >
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleActivate" class="flex-col gap-4">
        <div class="form-group">
          <label class="form-label" for="token">Token kích hoạt</label>
          <input
            id="token"
            v-model="token"
            type="text"
            class="form-input"
            placeholder="Nhập token từ email..."
            :disabled="isLoading"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="password">Mật khẩu mới</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-input"
            placeholder="Tối thiểu 6 ký tự"
            :disabled="isLoading"
            required
            minlength="6"
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="confirmPassword">Xác nhận mật khẩu</label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            class="form-input"
            placeholder="Nhập lại mật khẩu"
            :disabled="isLoading"
            required
            minlength="6"
          />
        </div>

        <button type="submit" class="btn btn--primary w-full" :disabled="isLoading">
          <template v-if="isLoading">
            <Loader2 class="spin mr-2" :size="16" />
          </template>
          Kích hoạt
        </button>

        <div class="text-center text-sm mt-4 text-muted">
          Đã có tài khoản?
          <RouterLink to="/login" class="font-medium text-foreground hover-underline"
            >Đăng nhập</RouterLink
          >
        </div>
      </form>
    </template>
  </div>
</template>

<style scoped>
.success-icon {
  color: var(--color-success);
}

.text-foreground {
  color: var(--color-text);
}

.hover-underline:hover {
  text-decoration: underline;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
