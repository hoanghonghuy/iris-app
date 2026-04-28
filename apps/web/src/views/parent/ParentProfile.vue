<script setup>
import { useAuthStore } from '../../stores/authStore'
import ChangePasswordForm from '../../components/ChangePasswordForm.vue'

const authStore = useAuthStore()
</script>

<template>
  <div class="profile-stack">
    <div class="card profile-card">
      <div class="card__body profile-body">
        <div class="field-readonly">
          <label>Email</label>
          <p>{{ authStore.currentUser?.email || '—' }}</p>
        </div>
        <div class="field-readonly">
          <label>Vai trò</label>
          <p>Phụ huynh</p>
        </div>
      </div>
    </div>

    <ChangePasswordForm />
  </div>
</template>

<style scoped>
.title {
  margin: 0;
  font-size: var(--font-size-lg);
}
</style>
