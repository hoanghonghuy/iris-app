<script setup>
import { onMounted, onUnmounted } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: 'md' // sm, md, lg
  }
})

const emit = defineEmits(['close'])

const handleEscape = (e) => {
  if (e.key === 'Escape' && props.isOpen) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="emit('close')">
    <div class="modal-dialog" :class="`modal-dialog--${size}`" role="dialog" aria-modal="true">
      <div class="modal-header">
        <h3 class="font-bold text-lg m-0">{{ title }}</h3>
        <button class="modal-close" @click="emit('close')">✕</button>
      </div>
      
      <div class="modal-body">
        <slot></slot>
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
  overflow-y: auto;
}

.modal-dialog {
  background-color: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  width: 100%;
  animation: modal-pop 0.2s ease-out;
  margin: auto;
  max-height: calc(100vh - 2 * var(--spacing-4));
  display: flex;
  flex-direction: column;
}

.modal-dialog--sm { max-width: 400px; }
.modal-dialog--md { max-width: 600px; }
.modal-dialog--lg { max-width: 800px; }

.modal-header {
  padding: var(--spacing-4) var(--spacing-5);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
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
  overflow-y: auto;
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
