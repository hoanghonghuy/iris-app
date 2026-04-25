<script setup>
import { computed } from 'vue'

const props = defineProps({
  currentPage: {
    type: Number,
    required: true
  },
  totalPages: {
    type: Number,
    required: true
  },
  totalItems: {
    type: Number,
    default: 0
  },
  limit: {
    type: Number,
    default: 10
  }
})

const emit = defineEmits(['page-change'])

// Lấy danh sách trang hiển thị (giới hạn tối đa 5 trang xung quanh trang hiện tại)
const pages = computed(() => {
  const range = []
  const maxPagesToShow = 5
  
  let start = Math.max(1, props.currentPage - Math.floor(maxPagesToShow / 2))
  let end = start + maxPagesToShow - 1
  
  if (end > props.totalPages) {
    end = props.totalPages
    start = Math.max(1, end - maxPagesToShow + 1)
  }
  
  for (let i = start; i <= end; i++) {
    range.push(i)
  }
  
  return range
})

const hasPrev = computed(() => props.currentPage > 1)
const hasNext = computed(() => props.currentPage < props.totalPages)

const startItem = computed(() => {
  if (props.totalItems === 0) return 0
  return (props.currentPage - 1) * props.limit + 1
})

const endItem = computed(() => {
  return Math.min(props.currentPage * props.limit, props.totalItems)
})

function goToPage(page) {
  if (page !== props.currentPage && page >= 1 && page <= props.totalPages) {
    emit('page-change', page)
  }
}
</script>

<template>
  <div class="pagination-bar" v-if="totalPages > 0">
    <div class="pagination-info text-sm text-muted hidden md-block">
      Hiển thị <span class="font-medium">{{ startItem }}</span> đến <span class="font-medium">{{ endItem }}</span> trong số <span class="font-medium">{{ totalItems }}</span> kết quả
    </div>

    <div class="pagination-controls">
      <button 
        class="pagination-btn" 
        :disabled="!hasPrev" 
        @click="goToPage(currentPage - 1)"
      >
        Trước
      </button>

      <div class="pagination-pages">
        <button 
          v-for="page in pages" 
          :key="page"
          class="pagination-page-btn"
          :class="{ 'pagination-page-btn--active': page === currentPage }"
          @click="goToPage(page)"
        >
          {{ page }}
        </button>
      </div>

      <button 
        class="pagination-btn" 
        :disabled="!hasNext" 
        @click="goToPage(currentPage + 1)"
      >
        Sau
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-4) 0;
  border-top: 1px solid var(--color-border);
  margin-top: var(--spacing-4);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.pagination-pages {
  display: flex;
  gap: var(--spacing-1);
}

.pagination-btn, .pagination-page-btn {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-sm);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.pagination-btn {
  padding: var(--spacing-1) var(--spacing-3);
}

.pagination-page-btn {
  min-width: 32px;
  height: 32px;
  padding: 0 var(--spacing-1);
}

.pagination-btn:hover:not(:disabled), 
.pagination-page-btn:hover:not(.pagination-page-btn--active) {
  background-color: var(--color-background);
  border-color: var(--color-text-muted);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-page-btn--active {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  color: var(--color-on-primary);
  font-weight: bold;
}

.hidden { display: none; }
@media (min-width: 768px) {
  .md-block { display: block; }
}
</style>
