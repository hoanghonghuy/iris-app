<script setup>
import { ref, onMounted, watch } from 'vue'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['google-login'])

const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID
const isScriptLoaded = ref(false)
const isGoogleInitialized = ref(false)

// 1. load script Google GIS
onMounted(() => {
  if (!googleClientId) return

  // kiểm tra nếu script đã load rồi
  if (window.google?.accounts?.id) {
    isScriptLoaded.value = true
    return
  }

  const script = document.createElement('script')
  script.src = 'https://accounts.google.com/gsi/client'
  script.async = true
  script.onload = () => {
    isScriptLoaded.value = true
  }
  document.head.appendChild(script)
})

// 2. khi script load xong → init + render button
// watch cũng xử lý re-render khi component được mount lại
watch(
  [isScriptLoaded, () => props.disabled],
  ([loaded, disabled]) => {
    if (!loaded) return

    if (!isGoogleInitialized.value) {
      window.google.accounts.id.initialize({
        client_id: googleClientId,
        callback: (response) => {
          // gửi ID token lên cho component cha xử lý
          emit('google-login', response.credential)
        },
      })
      isGoogleInitialized.value = true
    }

    const buttonContainer = document.getElementById('google-signin-button')
    if (buttonContainer && !disabled) {
      // Xoá nội dung cũ để render lại
      buttonContainer.innerHTML = ''
      window.google.accounts.id.renderButton(buttonContainer, {
        type: 'standard',
        size: 'large',
        shape: 'rectangular',
        theme: 'outline',
        text: 'signin_with',
        width: Math.min(400, Math.max(200, buttonContainer.clientWidth || 360)),
      })
    } else if (buttonContainer && disabled) {
      buttonContainer.innerHTML = ''
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="google-wrapper">
    <div v-if="googleClientId" class="google-button-wrapper" :class="{ 'is-disabled': disabled }">
      <!-- Cần một wrapper có min-height để không bị giật layout -->
      <div id="google-signin-button" class="google-button-container"></div>
    </div>
    <div v-else class="text-muted text-center text-sm p-4 border rounded">
      Google Sign-In chưa được cấu hình (Thiếu Client ID).
    </div>
  </div>
</template>

<style scoped>
.google-wrapper {
  width: 100%;
  display: flex;
  justify-content: center;
}

.google-button-wrapper {
  width: 100%;
  display: flex;
  justify-content: center;
  min-height: 40px; /* Chiều cao mặc định của nút Google */
}

.is-disabled {
  opacity: 0.5;
  pointer-events: none;
}
</style>
