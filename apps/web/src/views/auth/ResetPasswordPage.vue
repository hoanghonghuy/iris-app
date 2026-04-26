<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { CheckCircle2, ArrowLeft, Loader2 } from 'lucide-vue-next'
import { authService } from '../../services/authService'
import { extractErrorMessage } from '../../helpers/errorHandler'

const router = useRouter()

const email = ref('')
const token = ref('')
const password = ref('')
const confirmPassword = ref('')

const isLoading = ref(false)
const errorMessage = ref('')
const isSuccess = ref(false)

async function handleResetPassword() {
  errorMessage.value = ''

  if (!email.value || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    errorMessage.value = 'Vui lòng nhập Email hợp lệ'
    return
  }
  if (!token.value.trim()) {
    errorMessage.value = 'Mã đặt lại mật khẩu không được để trống'
    return
  }
  if (password.value.length < 6) {
    errorMessage.value = 'Mật khẩu phải có ít nhất 6 ký tự'
    return
  }
  if (password.value !== confirmPassword.value) {
    errorMessage.value = 'Mật khẩu xác nhận không khớp'
    return
  }

  isLoading.value = true
  
  try {
    await authService.resetPassword(email.value, token.value, password.value)
    isSuccess.value = true
    setTimeout(() => router.push('/login'), 3000)
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
      <div class="flex-col items-center gap-4 py-8 text-center">
        <CheckCircle2 class="success-icon mx-auto" :size="48" />
        <p class="text-sm font-medium mt-4">Đặt lại mật khẩu thành công!</p>
        <p class="text-sm text-muted mt-1">Đang chuyển về trang đăng nhập...</p>
      </div>
    </template>

    <!-- Form đặt lại mật khẩu -->
    <template v-else>
      <div class="text-center mb-6">
        <h2 class="text-2xl font-bold mb-1">Đặt lại mật khẩu</h2>
        <p class="text-muted text-sm">Nhập email, mã đặt lại mật khẩu và mật khẩu mới</p>
      </div>

      <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 text-danger text-sm rounded border border-red-200">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleResetPassword" class="flex-col gap-4">
        <div class="form-group">
          <label class="form-label" for="email">Email</label>
          <input 
            id="email" 
            v-model="email" 
            type="email" 
            class="form-input" 
            placeholder="name@example.com"
            :disabled="isLoading"
            required 
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="token">Mã đặt lại mật khẩu</label>
          <input 
            id="token" 
            v-model="token" 
            type="text" 
            class="form-input" 
            placeholder="Nhập mã trong email"
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

        <button type="submit" class="btn btn--primary w-full mt-2" :disabled="isLoading">
          <template v-if="isLoading">
            <Loader2 class="spin mr-2" :size="16" /> Đang xử lý...
          </template>
          <template v-else>
            Đặt lại mật khẩu
          </template>
        </button>
      </form>

      <div class="text-center text-sm mt-6">
        <RouterLink to="/login" class="text-muted hover-text-primary inline-flex items-center gap-1">
          <ArrowLeft :size="12" /> Quay lại đăng nhập
        </RouterLink>
      </div>
    </template>
  </div>
</template>

<style scoped>
.success-icon {
  color: var(--color-success);
}

.hover-text-primary:hover {
  color: var(--color-primary);
  transition: color 0.2s;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
