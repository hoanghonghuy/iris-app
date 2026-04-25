<script setup>
import { useRouter } from 'vue-router'
import { GraduationCap, ArrowRight, ShieldCheck, ClipboardCheck, HeartPulse, BookOpen, Menu } from 'lucide-vue-next'
import ThemeToggle from '../components/ThemeToggle.vue'

const router = useRouter()

const navigateToLogin = () => {
  router.push('/login')
}

const navigateToRegister = () => {
  router.push('/register')
}
</script>

<template>
  <div class="landing-page flex-col">
    <!-- Header Navigation -->
    <header class="header sticky">
      <div class="container flex justify-between items-center">
        <div class="flex items-center gap-2">
          <div class="logo-icon flex-center text-foreground"><GraduationCap :size="28" /></div>
          <span class="text-xl font-bold tracking-tight">Iris School</span>
        </div>
        
        <nav class="hidden md:flex gap-6 text-sm font-medium text-muted">
          <a href="#features" class="hover-text-primary">Tính năng</a>
          <a href="#about" class="hover-text-primary">Về chúng tôi</a>
        </nav>
        
        <div class="flex items-center gap-3">
          <ThemeToggle />
          <div class="hidden sm:flex items-center gap-3">
            <button class="btn btn--outline border-none" @click="navigateToLogin">Đăng nhập</button>
            <button class="btn btn--primary" @click="navigateToRegister">Đăng ký Phụ huynh</button>
          </div>
          <!-- Mobile Menu Toggle (Tạm ẩn) -->
          <button class="md:hidden btn btn--sm btn--outline border-none"><Menu :size="20" /></button>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <main class="flex-1">
      <section class="hero relative flex-center text-center">
        <!-- Glow background -->
        <div class="hero-glow"></div>
        <div class="hero-pattern"></div>
        
        <div class="container relative z-10">
          <h1 class="hero-title mb-6">
            Nền tảng Quản lý Tập trung dành cho Trường học
          </h1>
          <p class="hero-subtitle mb-8 mx-auto">
            Iris School kết nối Quản trị viên, Giáo viên và Phụ huynh trên một hệ thống duy nhất. Tối ưu hóa quy trình, nâng cao trải nghiệm giáo dục.
          </p>
          <div class="flex justify-center gap-4 hero-actions">
            <button class="btn btn--primary btn--lg w-full sm:w-auto" @click="navigateToLogin">
              Bắt đầu ngay <ArrowRight class="ml-2" :size="20" />
            </button>
            <a href="#features" class="btn btn--outline btn--lg w-full sm:w-auto hero-btn-outline text-center">
              Khám phá Tính năng
            </a>
          </div>
        </div>
      </section>

      <!-- Features Grid -->
      <section id="features" class="features py-20 bg-background">
        <div class="container">
          <div class="text-center mb-16">
            <h2 class="text-3xl font-bold mb-4">Mọi công cụ bạn cần</h2>
            <p class="text-lg text-muted max-w-2xl mx-auto">
              Được thiết kế tỉ mỉ phục vụ riêng cho nghiệp vụ giáo dục cấp Mầm non và Tiểu học.
            </p>
          </div>

          <div class="features-grid">
            <!-- Feature 1 -->
            <div class="feature-card card p-6">
              <div class="feature-icon mb-4 text-foreground"><ShieldCheck :size="24" /></div>
              <h3 class="text-xl font-bold mb-2">Quản lý Đa điểm trường</h3>
              <p class="text-muted leading-relaxed">
                Super Admin kiểm soát hệ thống, School Admin quản trị trọn vẹn từng điểm trường với lớp học và học sinh riêng biệt.
              </p>
            </div>

            <!-- Feature 2 -->
            <div class="feature-card card p-6">
              <div class="feature-icon mb-4 text-foreground"><ClipboardCheck :size="24" /></div>
              <h3 class="text-xl font-bold mb-2">Điểm danh Nhanh chóng</h3>
              <p class="text-muted leading-relaxed">
                Giáo viên dễ dàng điểm danh đầu giờ bằng mã Check in / Check out, dữ liệu đồng bộ ngay lập tức đến phụ huynh.
              </p>
            </div>

            <!-- Feature 3 -->
            <div class="feature-card card p-6">
              <div class="feature-icon mb-4 text-foreground"><HeartPulse :size="24" /></div>
              <h3 class="text-xl font-bold mb-2">Sổ theo dõi Sức khỏe</h3>
              <p class="text-muted leading-relaxed">
                Ghi nhận chỉ số BMI định kỳ. Cập nhật các thông tin y tế, bữa ăn và theo dõi sát sao tình trạng các bé.
              </p>
            </div>

            <!-- Feature 4 -->
            <div class="feature-card card p-6">
              <div class="feature-icon mb-4 text-foreground"><BookOpen :size="24" /></div>
              <h3 class="text-xl font-bold mb-2">Bảng tin Thông báo</h3>
              <p class="text-muted leading-relaxed">
                Giáo viên đăng tải hình ảnh hoạt động, thực đơn, dặn dò. Phụ huynh nhận tương tác qua News Feed hiện đại.
              </p>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="footer bg-muted border-t">
      <div class="container flex justify-between items-center flex-col md:flex-row gap-4">
        <div class="flex items-center gap-2">
          <div class="text-muted"><GraduationCap :size="20" /></div>
          <span class="text-sm font-bold text-foreground">Iris School</span>
        </div>
        <p class="text-sm text-muted text-center m-0">
          &copy; {{ new Date().getFullYear() }} Nền tảng Quản lý Trường học.
        </p>
        <div class="flex items-center gap-4 text-sm text-muted">
          <a href="#" class="hover-text-primary">Bảo mật</a>
          <a href="#" class="hover-text-primary">Điều khoản</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.landing-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: system-ui, -apple-system, sans-serif;
}

/* Utilities cho layout */
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding-left: var(--spacing-4);
  padding-right: var(--spacing-4);
}

.sticky {
  position: sticky;
  top: 0;
  z-index: 50;
}

.flex-1 { flex: 1; }
.flex-center { display: flex; align-items: center; justify-content: center; }
.text-center { text-align: center; }
.mx-auto { margin-left: auto; margin-right: auto; }
.w-full { width: 100%; }
.m-0 { margin: 0; }
.mb-2 { margin-bottom: var(--spacing-2); }
.mb-4 { margin-bottom: var(--spacing-4); }
.mb-6 { margin-bottom: var(--spacing-6); }
.mb-8 { margin-bottom: var(--spacing-8); }
.mb-16 { margin-bottom: 4rem; }
.ml-2 { margin-left: var(--spacing-2); }
.gap-2 { gap: var(--spacing-2); }
.gap-3 { gap: var(--spacing-3); }
.gap-4 { gap: var(--spacing-4); }
.gap-6 { gap: var(--spacing-6); }
.py-20 { padding-top: 5rem; padding-bottom: 5rem; }
.p-6 { padding: var(--spacing-6); }
.border-none { border: none; }
.border-t { border-top: 1px solid var(--color-border); }
.text-sm { font-size: var(--font-size-sm); }
.text-lg { font-size: var(--font-size-lg); }
.text-xl { font-size: var(--font-size-xl, 1.25rem); }
.text-3xl { font-size: 1.875rem; }
.font-bold { font-weight: 700; }
.font-medium { font-weight: 500; }
.tracking-tight { letter-spacing: -0.025em; }
.leading-relaxed { line-height: 1.625; }
.text-muted { color: var(--color-text-muted); }
.text-foreground { color: var(--color-text); }
.bg-muted { background-color: var(--color-background); } /* fallback */
.bg-background { background-color: var(--color-background); }
.hover-text-primary:hover { color: var(--color-primary); transition: color 0.2s; }

/* Header */
.header {
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-header-backdrop);
  backdrop-filter: blur(12px);
  padding: 1rem 0;
}
.dark .header {
  background-color: var(--color-header-backdrop-dark);
}
.logo-icon {
  width: 28px; height: 28px;
  font-size: 20px;
}

/* Hero Section */
.hero {
  padding: 6rem 1rem 8rem;
  overflow: hidden;
}
@media (min-width: 768px) {
  .hero { padding: 8rem 1rem 10rem; }
}

.hero-glow {
  position: absolute;
  top: 50%; left: 50%;
  transform: translate(-50%, -50%);
  width: 800px; height: 600px;
  background-color: var(--color-hero-glow);
  border-radius: 50%;
  filter: blur(120px);
  pointer-events: none;
}

.hero-pattern {
  position: absolute;
  inset: 0;
  background-image: 
    linear-gradient(to right, var(--color-border) 1px, transparent 1px),
    linear-gradient(to bottom, var(--color-border) 1px, transparent 1px);
  background-size: 14px 24px;
  opacity: 0.4;
}

.hero-title {
  font-size: 2.5rem;
  font-weight: 800;
  letter-spacing: -0.05em;
  background: linear-gradient(to right, var(--color-text), var(--color-text-muted));
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
@media (min-width: 768px) { .hero-title { font-size: 3.75rem; } }
@media (min-width: 1024px) { .hero-title { font-size: 4.5rem; } }

.hero-subtitle {
  max-width: 700px;
  font-size: 1.125rem;
  font-weight: 300;
}
@media (min-width: 768px) { .hero-subtitle { font-size: 1.25rem; } }

.hero-actions {
  flex-direction: column;
}
@media (min-width: 640px) {
  .hero-actions { flex-direction: row; }
}

.btn--lg {
  padding: 0.75rem 2rem;
  font-size: 1.125rem;
}
.hero-btn-outline {
  background-color: var(--color-hero-button-overlay);
  backdrop-filter: blur(4px);
  text-decoration: none;
}
.dark .hero-btn-outline {
  background-color: var(--color-hero-button-overlay-dark);
}

/* Features */
.features-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}
@media (min-width: 768px) { .features-grid { grid-template-columns: repeat(2, 1fr); } }
@media (min-width: 1024px) { .features-grid { grid-template-columns: repeat(3, 1fr); } }

.feature-card {
  background-color: color-mix(in srgb, var(--color-background) 50%, var(--color-surface));
  transition: all 0.3s ease;
  border-radius: 1rem;
}
.feature-card:hover {
  box-shadow: var(--shadow-md);
  border-color: color-mix(in srgb, var(--color-primary) 30%, transparent);
}

.feature-icon {
  width: 48px;
  height: 48px;
  background-color: var(--color-background);
  border-radius: var(--radius);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

/* Footer */
.footer {
  padding: 2.5rem 1rem;
}

/* Responsive utils */
.hidden { display: none; }
@media (min-width: 640px) { .sm\:flex { display: flex; } .sm\:w-auto { width: auto; } }
@media (min-width: 768px) { .md\:flex { display: flex; } .md\:hidden { display: none; } .md\:flex-row { flex-direction: row; } }
</style>
