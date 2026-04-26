<script setup>
import { useAuthStore } from '../../stores/authStore'

const authStore = useAuthStore()
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">Tổng quan Dashboard</h2>
    <div class="card p-6">
      <p>Chào mừng <strong>{{ authStore.currentUser?.email }}</strong>!</p>
      <p class="text-muted mt-2">Tính năng đang được phát triển trong Phase tiếp theo.</p>
    </div>
  </div>
</template>
