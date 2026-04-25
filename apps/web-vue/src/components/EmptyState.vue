<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  message: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: 'box' // Tạm thời dùng string để map với SVG inline bên dưới
  }
})
</script>

<template>
  <div class="empty-state">
    <div class="empty-state__icon text-muted mb-4">
      <!-- Tạm thời dùng SVG cứng dựa trên prop icon -->
      <svg v-if="icon === 'heart'" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
      </svg>
      <svg v-else-if="icon === 'users'" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
      </svg>
      <svg v-else width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line>
      </svg>
    </div>
    <h3 class="empty-state__title font-bold text-lg mb-2">{{ title }}</h3>
    <p v-if="message" class="empty-state__message text-muted">{{ message }}</p>
    
    <div v-if="$slots.action" class="empty-state__action mt-6">
      <slot name="action"></slot>
    </div>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: var(--spacing-8) var(--spacing-4);
  width: 100%;
}

.mb-2 { margin-bottom: var(--spacing-2); }
.mb-4 { margin-bottom: var(--spacing-4); }
.mt-6 { margin-top: var(--spacing-6); }

.empty-state__icon {
  opacity: 0.5;
}

.empty-state__title {
  color: var(--color-text);
}

.empty-state__message {
  max-width: 400px;
}
</style>
