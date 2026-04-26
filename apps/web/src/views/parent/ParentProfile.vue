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
.profile-stack,
.profile-body {
  display: flex;
  flex-direction: column;
}

.profile-stack {
  gap: var(--spacing-6);
}

.profile-body {
  gap: var(--spacing-4);
}

.profile-card {
  max-width: 32rem;
}

.title {
  margin: 0;
  font-size: var(--font-size-lg);
}

.field-readonly label {
  display: block;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  font-weight: 500;
  margin-bottom: var(--spacing-2);
}

.field-readonly p {
  margin: 0;
  font-size: var(--font-size-sm);
}
</style>
