<script setup>
const props = defineProps({
  analytics: {
    type: Object,
    default: null,
  },
  activeAppointmentsCount: {
    type: Number,
    default: 0,
  },
})
</script>

<template>
  <div class="card summary-card">
    <div class="summary-grid">
      <div class="summary-item">
        <p class="summary-label">Tổng số con</p>
        <p class="summary-value">{{ props.analytics?.total_children ?? 0 }}</p>
      </div>
      <div class="summary-item">
        <p class="summary-label">Lịch đang hoạt động</p>
        <p class="summary-value">{{ props.activeAppointmentsCount }}</p>
      </div>
      <div class="summary-item">
        <p class="summary-label">Cảnh báo sức khỏe 24h</p>
        <p class="summary-value">{{ props.analytics?.recent_health_alerts_24h ?? 0 }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.summary-card {
  padding: var(--spacing-4);
}

.summary-grid {
  display: grid;
  gap: var(--spacing-3);
}

.summary-item {
  padding: var(--spacing-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  background: var(--color-background);
}

.summary-label {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.summary-value {
  margin: var(--spacing-2) 0 0;
  font-size: var(--font-size-xl);
  font-weight: 700;
}

@media (min-width: 768px) {
  .summary-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
