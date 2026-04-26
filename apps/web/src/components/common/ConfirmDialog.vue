<script setup>
import { watchEffect } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    required: true
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: 'Xác nhận'
  },
  cancelText: {
    type: String,
    default: 'Hủy'
  },
  isDanger: {
    type: Boolean,
    default: false
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['confirm', 'cancel'])

function handleCancel() {
  if (!props.isLoading) {
    emit('cancel')
  }
}

const handleEscape = (e) => {
  if (e.key === 'Escape' && props.isOpen && !props.isLoading) {
    handleCancel()
  }
}

watchEffect((onCleanup) => {
  if (!props.isOpen) {
    return
  }

  document.addEventListener('keydown', handleEscape)
  onCleanup(() => {
    document.removeEventListener('keydown', handleEscape)
  })
})
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop">
    <div class="modal-dialog" role="dialog" aria-modal="true">
      <div class="modal-header">
        <h3 class="font-bold text-lg m-0">{{ title }}</h3>
        <button 
          class="modal-close" 
          type="button"
          aria-label="Đóng"
          @click="handleCancel"
          :disabled="isLoading"
        >
          ✕
        </button>
      </div>
      
      <div class="modal-body">
        <p>{{ message }}</p>
      </div>
      
      <div class="modal-footer">
        <button 
          class="btn btn--outline" 
          type="button"
          @click="handleCancel"
          :disabled="isLoading"
        >
          {{ cancelText }}
        </button>
        <button 
          class="btn" 
          type="button"
          :class="isDanger ? 'btn--danger' : 'btn--primary'"
          @click="emit('confirm')"
          :disabled="isLoading"
        >
          {{ isLoading ? 'Đang xử lý...' : confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: var(--color-overlay);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: var(--spacing-4);
}

.modal-dialog {
  background-color: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  width: 100%;
  max-width: 450px;
  overflow: hidden;
  animation: modal-pop 0.2s ease-out;
}

.modal-header {
  padding: var(--spacing-4) var(--spacing-5);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-close {
  background: none;
  border: none;
  color: var(--color-text-muted);
  font-size: var(--font-size-lg);
  cursor: pointer;
  padding: 0;
  display: flex;
}

.modal-close:hover {
  color: var(--color-text);
}

.modal-body {
  padding: var(--spacing-5);
  color: var(--color-text);
  line-height: 1.5;
}

.modal-footer {
  padding: var(--spacing-4) var(--spacing-5);
  border-top: 1px solid var(--color-border);
  background-color: var(--color-background);
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-3);
}

@keyframes modal-pop {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>
