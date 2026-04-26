import { ref } from 'vue'

/**
 * Shared composable for authentication form handling
 * Used across Login, Register, ForgotPassword, ResetPassword, Activate pages
 */
export function useAuthForm() {
  const isLoading = ref(false)
  const errorMessage = ref('')
  const isSuccess = ref(false)

  function clearError() {
    errorMessage.value = ''
  }

  function setError(error) {
    errorMessage.value = error
  }

  function setSuccess() {
    isSuccess.value = true
    errorMessage.value = ''
  }

  async function handleSubmit(submitFn) {
    isLoading.value = true
    errorMessage.value = ''
    isSuccess.value = false

    try {
      await submitFn()
      isSuccess.value = true
    } catch (error) {
      errorMessage.value = error.message || 'Đã xảy ra lỗi'
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    errorMessage,
    isSuccess,
    clearError,
    setError,
    setSuccess,
    handleSubmit,
  }
}
