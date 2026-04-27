<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Laptop, Moon, Sun } from 'lucide-vue-next'

const isOpen = ref(false)
const theme = ref('system')
const isDark = ref(false)
let mediaQuery = null

const currentThemeIcon = computed(() => (isDark.value ? Moon : Sun))
const themeButtonLabel = computed(() => `Giao diện hiện tại: ${theme.value}`)

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const closeMenu = (event) => {
  if (!event.target.closest('.theme-toggle-container')) {
    isOpen.value = false
  }
}

const syncDocumentTheme = (darkMode) => {
  isDark.value = darkMode
  document.documentElement.classList.toggle('dark', darkMode)
}

const applyTheme = (newTheme) => {
  theme.value = newTheme
  localStorage.setItem('theme', newTheme)

  if (newTheme === 'dark') {
    syncDocumentTheme(true)
  } else if (newTheme === 'light') {
    syncDocumentTheme(false)
  } else {
    syncDocumentTheme(
      mediaQuery?.matches ?? window.matchMedia('(prefers-color-scheme: dark)').matches,
    )
  }

  isOpen.value = false
}

const handleSystemThemeChange = (event) => {
  if (theme.value === 'system') {
    syncDocumentTheme(event.matches)
  }
}

onMounted(() => {
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  applyTheme(localStorage.getItem('theme') || 'system')
  mediaQuery.addEventListener('change', handleSystemThemeChange)
  document.addEventListener('click', closeMenu)
})

onUnmounted(() => {
  mediaQuery?.removeEventListener('change', handleSystemThemeChange)
  document.removeEventListener('click', closeMenu)
})
</script>

<template>
  <div class="theme-toggle-container relative">
    <button
      type="button"
      class="theme-toggle-btn"
      title="Chọn giao diện"
      :aria-label="themeButtonLabel"
      :aria-expanded="isOpen"
      aria-haspopup="menu"
      @click="toggleMenu"
    >
      <component :is="currentThemeIcon" :size="20" />
      <span class="sr-only">Toggle theme</span>
    </button>

    <div v-if="isOpen" class="theme-menu">
      <button
        type="button"
        class="theme-menu-item"
        :class="{ active: theme === 'light' }"
        @click="applyTheme('light')"
      >
        <Sun :size="16" class="mr-2" /> Sáng
      </button>
      <button
        type="button"
        class="theme-menu-item"
        :class="{ active: theme === 'dark' }"
        @click="applyTheme('dark')"
      >
        <Moon :size="16" class="mr-2" /> Tối
      </button>
      <button
        type="button"
        class="theme-menu-item"
        :class="{ active: theme === 'system' }"
        @click="applyTheme('system')"
      >
        <Laptop :size="16" class="mr-2" /> Hệ thống
      </button>
    </div>
  </div>
</template>

<style scoped>
.relative {
  position: relative;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

.mr-2 {
  margin-right: 0.5rem;
}

.theme-toggle-btn {
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 44px;
  min-width: 44px;
  padding: var(--spacing-2);
  border-radius: var(--radius);
  transition: all var(--transition-fast);
}

.theme-toggle-btn:hover {
  background: var(--color-surface-muted);
  color: var(--color-text);
}

.theme-toggle-btn:focus-visible,
.theme-menu-item:focus-visible {
  box-shadow: 0 0 0 3px var(--color-primary-focus-ring);
}

.theme-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.5rem;
  width: 150px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  padding: 0.25rem;
  z-index: 50;
  display: flex;
  flex-direction: column;
}

.theme-menu-item {
  display: flex;
  align-items: center;
  min-height: 40px;
  width: 100%;
  padding: 0.5rem;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--color-text);
  font-size: 0.875rem;
  cursor: pointer;
  text-align: left;
  transition: background-color var(--transition-fast);
}

.theme-menu-item:hover {
  background-color: var(--color-background);
}

.theme-menu-item.active {
  background-color: var(--color-background);
  color: var(--color-primary);
  font-weight: 500;
}
</style>
