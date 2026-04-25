<script setup>
defineProps({
  message: {
    type: String,
    default: 'Đang tải...'
  },
  fullScreen: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <div class="loading-wrapper" :class="{ 'loading-wrapper--fullscreen': fullScreen }">
    <div class="spinner"></div>
    <p v-if="message" class="loading-message">{{ message }}</p>
  </div>
</template>

<style scoped>
.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-8);
  width: 100%;
}

.loading-wrapper--fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: var(--color-loading-overlay);
  z-index: 50;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-message {
  margin-top: var(--spacing-4);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
