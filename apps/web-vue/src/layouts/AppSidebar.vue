<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import { adminMenuItems, teacherMenuItems, parentMenuItems } from '../helpers/authConfig'
import {
  LayoutDashboard,
  School,
  GraduationCap,
  Users,
  UserCog,
  BookUser,
  Heart,
  ShieldCheck,
  ListChecks,
  ClipboardCheck,
  HeartPulse,
  FileText,
  CalendarClock,
  Menu,
  Baby,
  Newspaper,
  Bell,
  MessageSquare,
  X
} from 'lucide-vue-next'

const iconMap = {
  dashboard: LayoutDashboard,
  school: School,
  class: GraduationCap,
  students: Users,
  users: UserCog,
  teacher: BookUser,
  parent: Heart,
  shield: ShieldCheck,
  logs: ListChecks,
  attendance: ClipboardCheck,
  health: HeartPulse,
  post: FileText,
  calendar: CalendarClock,
  menu: Menu,
  child: Baby,
  feed: Newspaper,
  bell: Bell,
  message: MessageSquare
}

defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close-sidebar'])

const route = useRoute()
const authStore = useAuthStore()

const menuItems = computed(() => {
  let items = []
  if (authStore.isAdmin) {
    items = adminMenuItems
  } else if (authStore.isTeacher) {
    items = teacherMenuItems
  } else if (authStore.isParent) {
    items = parentMenuItems
  }

  // Lọc theo role (VD: SUPER_ADMIN mới thấy menu School Admin)
  return items.filter(item => {
    if (!item.roles) return true
    return item.roles.includes(authStore.currentUserRole)
  })
})

const isActive = (path) => {
  // So sánh chính xác cho dashboard tổng quan, so sánh prefix cho các trang con
  if (path === '/admin' || path === '/teacher' || path === '/parent') {
    return route.path === path
  }
  return route.path.startsWith(path)
}

function handleOverlayClick() {
  emit('close-sidebar')
}
</script>

<template>
  <div class="sidebar-wrapper">
    <!-- Overlay cho mobile -->
    <div 
      class="sidebar-overlay" 
      :class="{ 'sidebar-overlay--open': isOpen }" 
      @click="handleOverlayClick"
    ></div>

    <!-- Sidebar cố định -->
    <aside class="sidebar" :class="{ 'sidebar--open': isOpen }">
      <div class="sidebar__header">
        <h2 class="sidebar__brand">🎓 Iris School</h2>
        <button class="sidebar__close-btn lg-hidden" @click="emit('close-sidebar')">
          <X :size="20" />
        </button>
      </div>

      <nav class="sidebar__nav">
        <RouterLink 
          v-for="item in menuItems" 
          :key="item.path" 
          :to="item.path"
          class="sidebar__nav-item"
          :class="{ 'sidebar__nav-item--active': isActive(item.path) }"
          @click="emit('close-sidebar')"
        >
          <span class="sidebar__nav-icon">
            <component :is="iconMap[item.icon] || LayoutDashboard" :size="20" />
          </span>
          <span class="sidebar__nav-label">{{ item.label }}</span>
        </RouterLink>
      </nav>
    </aside>
  </div>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 40;
  width: var(--sidebar-width);
  height: 100vh;
  background-color: var(--color-surface);
  border-right: 1px solid var(--color-border);
  transition: transform 0.3s ease-in-out;
  transform: translateX(-100%);
  display: flex;
  flex-direction: column;
}

.sidebar--open {
  transform: translateX(0);
}

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: var(--color-overlay);
  z-index: 30;
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.3s ease-in-out;
}

.sidebar-overlay--open {
  opacity: 1;
  visibility: visible;
}

@media (min-width: 1024px) {
  .sidebar {
    transform: translateX(0);
  }
  .sidebar-overlay {
    display: none;
  }
  .lg-hidden {
    display: none;
  }
}

.sidebar__header {
  height: var(--header-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-4);
  border-bottom: 1px solid var(--color-border);
}

.sidebar__brand {
  font-size: var(--font-size-xl);
  font-weight: bold;
  color: var(--color-primary);
  margin: 0;
}

.sidebar__close-btn {
  background: none;
  border: none;
  font-size: var(--font-size-lg);
  color: var(--color-text-muted);
  cursor: pointer;
}

.sidebar__nav {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.sidebar__nav-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-md);
  color: var(--color-text);
  transition: all 0.2s;
  font-weight: 500;
}

.sidebar__nav-item:hover {
  background-color: var(--color-background);
  color: var(--color-primary);
}

.sidebar__nav-item--active {
  background-color: var(--color-sidebar-active-bg);
  color: var(--color-primary);
}

.sidebar__nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: inherit;
}
</style>
