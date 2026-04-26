<script setup>
import { computed } from 'vue'

const MAX_VISIBLE_PAGES = 5

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

const safeTotalPages = computed(() => Math.max(Number(props.totalPages) || 0, 0))

const safeLimit = computed(() => Math.max(Number(props.limit) || 10, 1))

const safeCurrentPage = computed(() => {
  if (safeTotalPages.value === 0) {
    return 1
  }

  const rawPage = Number(props.currentPage) || 1
  return Math.min(Math.max(rawPage, 1), safeTotalPages.value)
})

const shouldShowPagination = computed(() => safeTotalPages.value > 0)

const pages = computed(() => {
  const range = []

  if (!shouldShowPagination.value) {
    return range
  }
  
  let start = Math.max(1, safeCurrentPage.value - Math.floor(MAX_VISIBLE_PAGES / 2))
  let end = start + MAX_VISIBLE_PAGES - 1
  
  if (end > safeTotalPages.value) {
    end = safeTotalPages.value
    start = Math.max(1, end - MAX_VISIBLE_PAGES + 1)
  }
  
  for (let i = start; i <= end; i++) {
    range.push(i)
  }
  
  return range
})

const hasPrev = computed(() => safeCurrentPage.value > 1)
const hasNext = computed(() => safeCurrentPage.value < safeTotalPages.value)

const startItem = computed(() => {
  if (props.totalItems <= 0) {
    return 0
  }

  return (safeCurrentPage.value - 1) * safeLimit.value + 1
})

const endItem = computed(() => {
  if (props.totalItems <= 0) {
    return 0
  }

  return Math.min(safeCurrentPage.value * safeLimit.value, props.totalItems)
})

function goToPage(page) {
  if (page !== safeCurrentPage.value && page >= 1 && page <= safeTotalPages.value) {
    emit('page-change', page)
  }
}
</script>

<template>
  <div v-if="shouldShowPagination" class="pagination-bar">
    <div class="pagination-info text-sm text-muted">
      Hiển thị <span class="font-medium">{{ startItem }}</span> đến <span class="font-medium">{{ endItem }}</span> trong số <span class="font-medium">{{ totalItems }}</span> kết quả
    </div>

    <div class="pagination-controls">
      <button 
        type="button"
        class="pagination-btn" 
        :disabled="!hasPrev" 
        @click="goToPage(safeCurrentPage - 1)"
      >
        Trước
      </button>

      <div class="pagination-pages">
        <button 
          v-for="page in pages" 
          :key="page"
          type="button"
          class="pagination-page-btn"
          :class="{ 'pagination-page-btn--active': page === safeCurrentPage }"
          @click="goToPage(page)"
        >
          {{ page }}
        </button>
      </div>

      <button 
        type="button"
        class="pagination-btn" 
        :disabled="!hasNext" 
        @click="goToPage(safeCurrentPage + 1)"
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
  font-weight: 700;
}

.pagination-info {
  display: none;
}

@media (min-width: 768px) {
  .pagination-info {
    display: block;
  }
}
</style>
