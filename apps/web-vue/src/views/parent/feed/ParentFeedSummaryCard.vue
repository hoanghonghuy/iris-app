<script setup>
const props = defineProps({
  totalPosts: {
    type: Number,
    default: 0,
  },
})
</script>

<template>
  <div class="card summary-card">
    <h2 class="summary-title">Bảng tin tổng hợp</h2>
    <p class="summary-description">Cập nhật từ giáo viên liên quan đến con của bạn.</p>
    <p class="summary-total">Tổng số bài: {{ props.totalPosts }}</p>
  </div>
</template>

<style scoped>
.summary-card {
  padding: var(--spacing-4);
}

.summary-title,
.summary-description,
.summary-total {
  margin: 0;
}

.summary-title {
  font-size: var(--font-size-lg);
}

.summary-description {
  margin-top: var(--spacing-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.summary-total {
  margin-top: var(--spacing-1);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}
</style>
