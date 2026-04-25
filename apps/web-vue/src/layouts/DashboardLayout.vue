<script setup>
import { ref } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const isSidebarOpen = ref(false)

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}

const closeSidebar = () => {
  isSidebarOpen.value = false
}
</script>

<template>
  <div class="dashboard-layout">
    <!-- Sidebar -->
    <AppSidebar :is-open="isSidebarOpen" @close-sidebar="closeSidebar" />

    <!-- Main Content -->
    <div class="dashboard-main">
      <AppHeader @toggle-sidebar="toggleSidebar" />
      
      <main class="dashboard-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style scoped>
.dashboard-layout {
  display: flex;
  min-height: 100vh;
  width: 100%;
  background-color: var(--color-background);
}

.dashboard-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0; /* Cho phép co lại khi cửa sổ nhỏ */
}

@media (min-width: 1024px) {
  .dashboard-main {
    margin-left: var(--sidebar-width);
  }
}

.dashboard-content {
  flex: 1;
  padding: var(--spacing-4) var(--spacing-4) var(--spacing-8);
  overflow-y: auto;
}

@media (min-width: 768px) {
  .dashboard-content {
    padding: var(--spacing-6) var(--spacing-6) var(--spacing-8);
  }
}
</style>
