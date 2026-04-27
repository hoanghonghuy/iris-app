<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/authStore'
import { authService } from '../../services/authService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import GoogleSignInButton from '../../components/GoogleSignInButton.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// State form chính
const email = ref('')
const password = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

// State cho flow Link Password
const isLinkMode = ref(false)
const tempGoogleIdToken = ref(null)
const linkPassword = ref('')

async function handleLogin() {
  if (!email.value || !password.value) {
    errorMessage.value = 'Vui lòng nhập đầy đủ email và mật khẩu'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    const data = await authService.login(email.value, password.value)
    handleLoginSuccess(data)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

async function handleGoogleLogin(idToken) {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const data = await authService.loginWithGoogle(idToken)
    handleLoginSuccess(data)
  } catch (error) {
    const errorMsg = extractErrorMessage(error)
    // Nếu API trả về GOOGLE_LINK_PASSWORD_REQUIRED -> Bật form nhập mật khẩu
    if (
      error?.data?.error_code === 'GOOGLE_LINK_PASSWORD_REQUIRED' ||
      errorMsg === 'GOOGLE_LINK_PASSWORD_REQUIRED'
    ) {
      isLinkMode.value = true
      tempGoogleIdToken.value = idToken
      errorMessage.value =
        'Tài khoản này đã tồn tại. Vui lòng nhập mật khẩu để liên kết với Google.'
    } else {
      errorMessage.value = errorMsg
    }
  } finally {
    isLoading.value = false
  }
}

async function handleLinkPassword() {
  if (!linkPassword.value) {
    errorMessage.value = 'Vui lòng nhập mật khẩu'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    const data = await authService.linkGooglePassword(tempGoogleIdToken.value, linkPassword.value)
    handleLoginSuccess(data)
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
  } finally {
    isLoading.value = false
  }
}

function handleLoginSuccess(data) {
  // Backend trả về data.data.access_token
  const token = data?.data?.access_token || data?.access_token
  if (!token) {
    errorMessage.value = 'Không nhận được token từ server'
    return
  }
  authStore.setToken(token)

  // Gọi API lấy thông tin user để lưu role
  authStore.fetchCurrentUser().then((user) => {
    if (!user) return
    // Chuyển hướng
    const redirect = route.query.redirect
    if (redirect) {
      router.push(redirect)
      return
    }
    // BE trả roles là array
    const primaryRole = Array.isArray(user.roles) ? user.roles[0] : user.role
    if (primaryRole === 'SUPER_ADMIN' || primaryRole === 'SCHOOL_ADMIN') {
      router.push('/admin')
    } else if (primaryRole === 'TEACHER') {
      router.push('/teacher')
    } else if (primaryRole === 'PARENT') {
      router.push('/parent')
    } else {
      router.push('/')
    }
  })
}

function cancelLinkMode() {
  isLinkMode.value = false
  tempGoogleIdToken.value = null
  linkPassword.value = ''
  errorMessage.value = ''
}
</script>

<template>
  <div class="auth-card">
    <div class="text-center mb-6">
      <h2 class="text-2xl font-bold mb-2">Iris School</h2>
      <p class="text-muted text-sm">Đăng nhập để quản lý thông tin trường học</p>
    </div>

    <div
      v-if="errorMessage"
      class="mb-4 p-3 bg-red-50 text-danger text-sm rounded border border-red-200"
    >
      {{ errorMessage }}
    </div>

    <!-- Form Liên kết Google (Link Mode) -->
    <form v-if="isLinkMode" @submit.prevent="handleLinkPassword" class="flex-col gap-4">
      <div class="form-group">
        <label class="form-label" for="link-password">Mật khẩu hiện tại</label>
        <input
          id="link-password"
          v-model="linkPassword"
          type="password"
          class="form-input"
          placeholder="Nhập mật khẩu của bạn"
          :disabled="isLoading"
          required
        />
      </div>

      <div class="flex gap-2 mt-4">
        <button
          type="button"
          class="btn btn--outline w-full"
          @click="cancelLinkMode"
          :disabled="isLoading"
        >
          Hủy
        </button>
        <button type="submit" class="btn btn--primary w-full" :disabled="isLoading">
          {{ isLoading ? 'Đang xử lý...' : 'Liên kết tài khoản' }}
        </button>
      </div>
    </form>

    <!-- Form Đăng nhập thường -->
    <form v-else @submit.prevent="handleLogin" class="flex-col gap-4">
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
        <div class="flex justify-between items-center mb-1">
          <label class="form-label m-0" for="password">Mật khẩu</label>
          <RouterLink to="/forgot-password" class="text-xs text-muted hover-text-primary"
            >Quên mật khẩu?</RouterLink
          >
        </div>
        <input
          id="password"
          v-model="password"
          type="password"
          class="form-input"
          placeholder="••••••••"
          :disabled="isLoading"
          required
        />
      </div>

      <button type="submit" class="btn btn--primary w-full mt-2" :disabled="isLoading">
        {{ isLoading ? 'Đang đăng nhập...' : 'Đăng nhập' }}
      </button>

      <div class="divider my-4">
        <span class="divider-text">phương thức khác</span>
      </div>

      <GoogleSignInButton :disabled="isLoading" @google-login="handleGoogleLogin" />

      <div class="text-center text-sm mt-6 text-muted">
        Phụ huynh chưa có tài khoản?
        <RouterLink to="/register" class="font-medium text-primary">Đăng ký tại đây</RouterLink>
      </div>
    </form>
  </div>
</template>

<style scoped>
.mb-6 {
  margin-bottom: var(--spacing-6);
}
.bg-red-50 {
  background-color: var(--color-danger-soft-bg);
}
.border-red-200 {
  border-color: var(--color-danger-soft-border);
}
.rounded {
  border-radius: var(--radius);
}
.my-4 {
  margin-top: var(--spacing-4);
  margin-bottom: var(--spacing-4);
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--color-border);
}

.divider:not(:empty)::before {
  margin-right: 0.25em;
}

.divider:not(:empty)::after {
  margin-left: 0.25em;
}

.divider-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  padding: 0 var(--spacing-2);
  text-transform: uppercase;
}

.hover-text-primary {
  transition: color var(--transition-fast);
}

.hover-text-primary:hover {
  color: var(--color-primary);
}

.text-primary:focus-visible,
.hover-text-primary:focus-visible {
  border-radius: var(--radius-sm);
  box-shadow: 0 0 0 3px var(--color-primary-focus-ring);
}
</style>
