<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Heart } from 'lucide-vue-next'
import { authService } from '../../services/authService'
import { useAuthStore } from '../../stores/authStore'
import { extractErrorMessage } from '../../helpers/errorHandler'
import GoogleSignInButton from '../../components/GoogleSignInButton.vue'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const parentCode = ref('')

const isLoading = ref(false)
const errorMessage = ref('')

// Flow 1: Đăng ký thường
async function handleRegister() {
  if (!email.value || !password.value || !parentCode.value) {
    errorMessage.value = 'Vui lòng nhập đầy đủ các trường bắt buộc'
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
  errorMessage.value = ''
  
  try {
    await authService.registerParent(
      email.value, 
      password.value, 
      parentCode.value
    )
    
    // Đăng ký thành công -> Tự động đăng nhập
    const data = await authService.login(email.value, password.value)
    handleLoginSuccess(data)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

// Flow 2: Đăng ký Google
async function handleGoogleRegister(idToken) {
  if (!parentCode.value) {
    errorMessage.value = 'Vui lòng nhập Mã Phụ Huynh trước khi tiếp tục với Google'
    return
  }

  isLoading.value = true
  errorMessage.value = ''
  
  try {
    const data = await authService.registerParentWithGoogle(
      idToken,
      parentCode.value
    )
    handleLoginSuccess(data)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

function handleLoginSuccess(data) {
  const token = data?.data?.access_token || data?.access_token
  if (!token) return
  authStore.setToken(token)
  authStore.fetchCurrentUser().then((user) => {
    if (!user) return
    router.push('/parent')
  })
}
</script>

<template>
  <div class="auth-card">
    <div class="text-center mb-6">
      <div class="flex justify-center mb-2">
        <Heart class="text-muted" :size="40" />
      </div>
      <h2 class="text-xl font-bold mb-1">Đăng ký Phụ huynh</h2>
      <p class="text-muted text-sm">Sử dụng mã phụ huynh từ nhà trường để đăng ký</p>
    </div>

    <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 text-danger text-sm rounded border border-red-200">
      {{ errorMessage }}
    </div>

    <form @submit.prevent="handleRegister" class="flex-col gap-4">
      <div class="form-group">
        <label class="form-label" for="parentCode">Mã phụ huynh <span class="text-danger">*</span></label>
        <input 
          id="parentCode" 
          v-model="parentCode" 
          type="text" 
          class="form-input" 
          placeholder="Nhập mã từ nhà trường..."
          :disabled="isLoading"
          required 
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="email">Email <span class="text-danger">*</span></label>
        <input 
          id="email" 
          v-model="email" 
          type="email" 
          class="form-input" 
          placeholder="parent@example.com"
          :disabled="isLoading"
          required 
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="password">Mật khẩu <span class="text-danger">*</span></label>
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
        <label class="form-label" for="confirmPassword">Xác nhận mật khẩu <span class="text-danger">*</span></label>
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
        {{ isLoading ? 'Đang xử lý...' : 'Đăng ký bằng Email' }}
      </button>

      <div class="divider my-4">
        <span class="divider-text">hoặc đăng nhập bằng</span>
      </div>

      <GoogleSignInButton 
        :disabled="isLoading" 
        @google-login="handleGoogleRegister" 
      />

      <div class="text-center text-sm mt-6 text-muted">
        Đã có tài khoản? 
        <RouterLink to="/login" class="font-medium text-foreground hover-underline">Đăng nhập</RouterLink>
      </div>
    </form>
  </div>
</template>

<style scoped>
.hover-underline:hover { text-decoration: underline; }

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  margin-top: var(--spacing-4);
  margin-bottom: var(--spacing-4);
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--color-border);
}

.divider:not(:empty)::before {
  margin-right: .25em;
}

.divider:not(:empty)::after {
  margin-left: .25em;
}

.divider-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  padding: 0 var(--spacing-2);
  text-transform: uppercase;
}

.text-foreground {
  color: var(--color-text);
}
</style>
