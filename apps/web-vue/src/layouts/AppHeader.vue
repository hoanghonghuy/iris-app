<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import { PROFILE_ROUTE_BY_ROLE, ROLE_LABELS } from '../helpers/authConfig'
import ThemeToggle from '../components/ThemeToggle.vue'
import { Menu, UserCircle, LogOut, ClipboardCheck } from 'lucide-vue-next'

const emit = defineEmits(['toggle-sidebar'])

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isDropdownOpen = ref(false)

const exactHeaderMeta = {
  '/teacher': { title: 'Tổng quan giáo viên' },
  '/parent': { title: 'Tổng quan phụ huynh' },
  '/admin': { title: 'Tổng quan quản trị' },
}

const prefixHeaderMetaRules = [
  { prefixes: ['/teacher/classes/'], meta: { title: 'Chi tiết lớp học' } },
  { prefixes: ['/teacher/classes'], meta: { title: 'Lớp của tôi' } },
  { prefixes: ['/teacher/health'], meta: { title: 'Sức khỏe học sinh' } },
  { prefixes: ['/teacher/posts'], meta: { title: 'Bài đăng' } },
  { prefixes: ['/teacher/appointments'], meta: { title: 'Lịch hẹn' } },
  { prefixes: ['/teacher/chat'], meta: { title: 'Tin nhắn' } },
  { prefixes: ['/teacher/profile'], meta: { title: 'Hồ sơ cá nhân' } },
  { prefixes: ['/parent/children/'], meta: { title: 'Thông tin con' } },
  { prefixes: ['/parent/children'], meta: { title: 'Con của tôi' } },
  { prefixes: ['/parent/posts', '/parent/feed'], meta: { title: 'Bảng tin' } },
  { prefixes: ['/parent/appointments'], meta: { title: 'Lịch hẹn' } },
  { prefixes: ['/parent/chat'], meta: { title: 'Tin nhắn' } },
  { prefixes: ['/parent/profile'], meta: { title: 'Hồ sơ cá nhân' } },
  { prefixes: ['/admin/schools'], meta: { title: 'Quản lý trường học' } },
  { prefixes: ['/admin/school-admins'], meta: { title: 'Quản lý School Admin' } },
  { prefixes: ['/admin/classes'], meta: { title: 'Quản lý lớp học' } },
  { prefixes: ['/admin/teachers'], meta: { title: 'Quản lý giáo viên' } },
  { prefixes: ['/admin/students'], meta: { title: 'Quản lý học sinh' } },
  { prefixes: ['/admin/parents'], meta: { title: 'Quản lý phụ huynh' } },
  { prefixes: ['/admin/users'], meta: { title: 'Quản lý người dùng' } },
  { prefixes: ['/admin/chat'], meta: { title: 'Tin nhắn' } },
]

function closeDropdown(event) {
  if (!event.target.closest('.user-menu')) {
    isDropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
})

function toggleDropdown() {
  isDropdownOpen.value = !isDropdownOpen.value
}

function handleLogout() {
  authStore.clearAuth()
  router.push('/login')
}

const roleLabel = computed(() => ROLE_LABELS[authStore.currentUserRole] || 'Người dùng')

const profileRoute = computed(() => PROFILE_ROUTE_BY_ROLE[authStore.currentUserRole] || null)

const headerMeta = computed(() => {
  const path = route.path
  const today = new Date().toISOString().slice(0, 10)

  if (path.startsWith('/teacher/attendance')) {
    return { title: 'Điểm danh', subtitle: `Ngày: ${today}`, icon: ClipboardCheck }
  }

  if (exactHeaderMeta[path]) {
    return exactHeaderMeta[path]
  }

  const matchedRule = prefixHeaderMetaRules.find((rule) =>
    rule.prefixes.some((prefix) => path.startsWith(prefix)),
  )

  return matchedRule?.meta || null
})

const userInitials = computed(() => {
  return authStore.currentUser?.email ? authStore.currentUser.email.substring(0, 2).toUpperCase() : 'U'
})
</script>

<template>
  <header class="header">
    <div class="header__left">
      <button class="header__menu-btn lg-hidden" @click="emit('toggle-sidebar')">
        <Menu :size="20" />
      </button>

      <span class="role-label">
        {{ roleLabel }}
      </span>
    </div>

    <div v-if="headerMeta" class="header__center">
      <div class="header__title-wrap">
        <component :is="headerMeta.icon" v-if="headerMeta.icon" class="text-muted" :size="16" />
        <p class="header__title">{{ headerMeta.title }}</p>
        <span v-if="headerMeta.subtitle" class="header__subtitle">
          {{ headerMeta.subtitle }}
        </span>
      </div>
    </div>

    <div class="header__right">
      <ThemeToggle class="mr-2" />

      <div class="user-menu relative">
        <button class="user-menu__trigger" @click="toggleDropdown">
          <div class="user-avatar">
            {{ userInitials }}
          </div>
          <span class="user-name">
            {{ authStore.currentUser?.email || 'Đang tải...' }}
          </span>
        </button>

        <div v-if="isDropdownOpen" class="user-menu__dropdown shadow-md">
          <div class="dropdown-header">
            <p class="text-sm font-medium leading-none truncate">{{ authStore.currentUser?.email }}</p>
            <p class="text-xs text-muted mt-1">{{ roleLabel }}</p>
          </div>
          <hr class="my-1 border-border" />

          <RouterLink
            v-if="profileRoute"
            :to="profileRoute"
            class="dropdown-item"
            @click="isDropdownOpen = false"
          >
            <UserCircle :size="16" class="mr-2" />
            Hồ sơ cá nhân
          </RouterLink>
          <hr v-if="profileRoute" class="my-1 border-border" />
          <button class="dropdown-item text-danger w-full text-left" @click="handleLogout">
            <LogOut :size="16" class="mr-2" />
            Đăng xuất
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  height: var(--header-height);
  background-color: color-mix(in srgb, var(--color-background) 80%, transparent);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-4);
  position: sticky;
  top: 0;
  z-index: 30;
  transition: background-color 0.3s;
}

@media (min-width: 1024px) {
  .header {
    padding: 0 var(--spacing-8);
  }
}

.header__left,
.header__right {
  flex: 1;
  display: flex;
  align-items: center;
}

.header__left {
  gap: var(--spacing-2);
}

.header__right {
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.header__center {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.header__title-wrap {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-2);
}

.header__title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin: 0;
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.header__subtitle {
  display: none;
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.role-label,
.user-name {
  display: none;
}

.header__menu-btn {
  background: none;
  border: none;
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-1);
  border-radius: var(--radius-md);
  transition: background-color 0.2s;
}

.header__menu-btn:hover,
.user-menu__trigger:hover,
.dropdown-item:hover {
  background-color: var(--color-background);
}

@media (min-width: 640px) {
  .role-label,
  .user-name,
  .header__subtitle {
    display: inline-block;
  }
}

@media (min-width: 1024px) {
  .lg-hidden {
    display: none;
  }
}

.relative {
  position: relative;
}

.user-menu__trigger {
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-1) var(--spacing-3) var(--spacing-1) var(--spacing-1);
  border-radius: var(--radius-full);
  transition: background-color 0.2s;
}

.user-name {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: var(--color-text);
}

.role-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: var(--radius-full);
  background-color: color-mix(in srgb, var(--color-primary) 12%, transparent);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: var(--font-size-xs);
}

.user-menu__dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: var(--spacing-2);
  width: 224px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--spacing-1) 0;
  box-shadow: var(--shadow-md);
}

.dropdown-header {
  padding: var(--spacing-2) var(--spacing-4);
}

.dropdown-item {
  display: flex;
  align-items: center;
  width: 100%;
  padding: var(--spacing-2) var(--spacing-4);
  font-size: var(--font-size-sm);
  color: var(--color-text);
  background: none;
  border: none;
  transition: background-color 0.2s;
}

.text-danger {
  color: var(--color-danger);
}

.my-1 {
  margin-top: var(--spacing-1);
  margin-bottom: var(--spacing-1);
}

.border-border {
  border-color: var(--color-border);
}

.mr-2 {
  margin-right: var(--spacing-2);
}

.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
