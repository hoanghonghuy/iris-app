<script setup>
import EmptyState from '../../components/common/EmptyState.vue'
</script>

<template>
  <div class="teacher-menu">
    <div class="card p-8 text-center mt-8">
      <EmptyState 
        title="Tính năng đang phát triển" 
        message="Phân hệ quản lý thực đơn đang được xây dựng và sẽ sớm ra mắt trong các phiên bản tiếp theo."
        icon="box"
      />
    </div>
  </div>
</template>

<style scoped>
.mb-6 { margin-bottom: var(--spacing-6); }
.mt-8 { margin-top: var(--spacing-8); }
.p-8 { padding: var(--spacing-8); }
.text-center { text-align: center; }
</style>
