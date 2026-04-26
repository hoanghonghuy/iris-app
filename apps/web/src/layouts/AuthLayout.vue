<script setup>
import { RouterView, useRouter } from 'vue-router'
import { GraduationCap, ShieldCheck, Zap, Activity } from 'lucide-vue-next'
import ThemeToggle from '../components/ThemeToggle.vue'

const router = useRouter()
const goHome = () => router.push('/')
</script>

<template>
  <div class="auth-layout">
    <!-- Nửa trái: Branding (chỉ hiện trên Desktop) -->
    <div class="auth-layout__brand">
      <!-- Background pattern -->
      <div class="brand-bg">
        <div class="brand-bg-pattern"></div>
        <div class="brand-bg-glow top-glow"></div>
        <div class="brand-bg-glow bottom-glow"></div>
      </div>

      <!-- Header -->
      <div class="brand-header">
        <div class="logo-link" @click="goHome">
          <GraduationCap class="text-primary" :size="28" />
          <span class="logo-text">Iris School</span>
        </div>
      </div>

      <!-- Center content -->
      <div class="brand-content">
        <h1 class="brand-title">
          Kết nối nhà trường và <span class="whitespace-nowrap">phụ huynh</span><br />
          trong một nền tảng toàn diện.
        </h1>
        <p class="brand-subtitle">
          Hệ sinh thái công nghệ giáo dục Iris School mang lại sự tiện lợi, minh bạch và an toàn.
          Tương tác nhanh chóng, theo dõi hành trình học viên hiệu quả mỗi ngày.
        </p>

        <div class="brand-features">
          <div class="feature-item">
            <div class="feature-icon feature-blue">
              <ShieldCheck :size="20" />
            </div>
            <span>Hệ thống bảo mật tối ưu</span>
          </div>
          <div class="feature-item">
            <div class="feature-icon feature-green">
              <Zap :size="20" />
            </div>
            <span>Thông báo tương tác tức thời</span>
          </div>
          <div class="feature-item">
            <div class="feature-icon feature-purple">
              <Activity :size="20" />
            </div>
            <span>Theo dõi sức khỏe & hoạt động</span>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="brand-footer">
        <p>&copy; {{ new Date().getFullYear() }} Iris Education. All rights reserved.</p>
      </div>
    </div>

    <!-- Nửa phải: Form (hiện trên mọi thiết bị) -->
    <div class="auth-layout__form-container">
      <!-- Mobile Header -->
      <div class="mobile-header">
        <div class="logo-link" @click="goHome">
          <GraduationCap class="text-primary" :size="24" />
          <span class="logo-text text-foreground">Iris School</span>
        </div>
        <ThemeToggle />
      </div>

      <!-- Desktop Theme Toggle -->
      <div class="desktop-theme-toggle">
        <ThemeToggle />
      </div>

      <div class="auth-layout__form-center">
        <div class="auth-layout__form-wrapper">
          <RouterView />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-layout {
  display: grid;
  min-height: 100vh;
  width: 100%;
  grid-template-columns: 1fr;
  background-color: var(--color-background);
  transition: background-color 0.3s;
}

@media (min-width: 1024px) {
  .auth-layout {
    grid-template-columns: 1fr 1fr;
  }
}

/* NỬA TRÁI */
.auth-layout__brand {
  display: none;
  flex-direction: column;
  justify-content: space-between;
  background-color: var(--color-auth-brand-bg);
  color: var(--color-on-primary);
  position: relative;
  overflow: hidden;
  border-right: 1px solid var(--color-auth-brand-border);
}

@media (min-width: 1024px) {
  .auth-layout__brand {
    display: flex;
  }
}

.brand-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.brand-bg-pattern {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(to right, var(--color-auth-pattern-line) 1px, transparent 1px),
    linear-gradient(to bottom, var(--color-auth-pattern-line) 1px, transparent 1px);
  background-size: 14px 24px;
  -webkit-mask-image: radial-gradient(ellipse 80% 50% at 50% 0%, black 70%, transparent 100%);
  mask-image: radial-gradient(ellipse 80% 50% at 50% 0%, black 70%, transparent 100%);
}

.brand-bg-glow {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}

.top-glow {
  top: -20%;
  left: -10%;
  width: 70%;
  height: 70%;
  background-color: var(--color-auth-glow-top);
  filter: blur(120px);
}

.bottom-glow {
  bottom: 0%;
  right: -20%;
  width: 60%;
  height: 60%;
  background-color: var(--color-auth-glow-bottom);
  filter: blur(100px);
}

.brand-header,
.brand-content,
.brand-footer {
  position: relative;
  z-index: 10;
  padding: 2rem;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.logo-link:hover {
  opacity: 0.9;
}

.logo-text {
  font-weight: 700;
  font-size: 1.25rem;
  letter-spacing: -0.025em;
  color: inherit;
}

.brand-content {
  width: 100%;
  max-width: 42rem; /* 2xl */
}

.brand-title {
  font-size: 1.875rem;
  font-weight: 600;
  letter-spacing: -0.025em;
  line-height: 1.25;
  margin: 0;
}

@media (min-width: 1280px) {
  .brand-title {
    font-size: 2.25rem;
  }
}

.whitespace-nowrap {
  white-space: nowrap;
}

.brand-subtitle {
  margin-top: 1.5rem;
  color: var(--color-auth-subtitle);
  font-size: 1rem;
  line-height: 1.625;
  max-width: 36rem;
}

.brand-features {
  margin-top: 2.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-auth-feature-text);
}

.feature-icon {
  display: flex;
  height: 2.5rem;
  width: 2.5rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 0.5rem;
  background-color: var(--color-auth-feature-bg);
  border: 1px solid var(--color-auth-feature-border);
  box-shadow: var(--shadow-sm);
}

.feature-blue {
  color: var(--color-auth-feature-blue);
}
.feature-green {
  color: var(--color-auth-feature-green);
}
.feature-purple {
  color: var(--color-auth-feature-purple);
}

.brand-footer p {
  font-size: 0.75rem;
  color: var(--color-auth-footer);
  font-weight: 500;
  margin: 0;
}

/* NỬA PHẢI */
.auth-layout__form-container {
  display: flex;
  flex-direction: column;
  position: relative;
  height: 100dvh;
  overflow-y: auto;
}

@media (min-width: 1024px) {
  .auth-layout__form-container {
    height: auto;
  }
}

.mobile-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem 0.5rem;
}

@media (min-width: 1024px) {
  .mobile-header {
    display: none;
  }
}

.text-primary {
  color: var(--color-primary);
}

.text-foreground {
  color: var(--color-text);
}

.desktop-theme-toggle {
  display: none;
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  z-index: 20;
}

@media (min-width: 1024px) {
  .desktop-theme-toggle {
    display: block;
  }
}

.auth-layout__form-center {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

@media (min-width: 640px) {
  .auth-layout__form-center {
    padding: 2rem;
  }
}

@media (min-width: 768px) {
  .auth-layout__form-center {
    padding: 3rem;
  }
}

.auth-layout__form-wrapper {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
