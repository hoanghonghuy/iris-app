import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowRight, BookOpen, HeartPulse, ClipboardCheck, GraduationCap, ShieldCheck } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";
import { MobileNavWrapper } from "@/components/layout/MobileNavWrapper";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col bg-background font-sans transition-colors duration-300">
      {/* ─── Header Navigation ──────────────────────────────────────────────────────── */}
      <header className="sticky top-0 z-50 w-full border-b border-border bg-background/80 backdrop-blur-md transition-colors duration-300">
        <div className="container mx-auto flex h-16 items-center justify-between px-4 md:px-6 lg:max-w-7xl">
          <div className="flex items-center gap-2">
            <GraduationCap className="h-7 w-7 text-foreground" />
            <span className="text-xl font-bold tracking-tighter text-foreground">Iris School</span>
          </div>
          <nav className="hidden md:flex gap-6 text-sm font-medium text-muted-foreground">
            <a href="#features" className="hover:text-foreground transition-colors">Tính năng</a>
            <a href="#about" className="hover:text-foreground transition-colors">Về chúng tôi</a>
          </nav>
          <div className="flex items-center gap-2 sm:gap-3">
            <ThemeToggle className="text-muted-foreground" />
            <div className="hidden sm:flex items-center gap-2 sm:gap-3">
              <Link href="/login">
                <Button variant="ghost" className="text-foreground hover:bg-muted">Đăng nhập</Button>
              </Link>
              <Link href="/register">
                <Button>Đăng ký Phụ huynh</Button>
              </Link>
            </div>
            {/* Mobile Navigation — client-only wrapper để tránh SSR hydration mismatch với Radix */}
            <div className="md:hidden flex items-center ml-2">
              <MobileNavWrapper />
            </div>
          </div>
        </div>
      </header>

      {/* ─── Hero Section ───────────────────────────────────────────────────────────── */}
      <main className="flex-1">
        {/* Keeps a dark scheme for both modes as a prominent visual focal point (Dark Hero Pattern) */}
        <section className="relative px-4 py-24 md:py-32 lg:py-40 bg-background text-foreground overflow-hidden">
          {/* Subtle Primary Glow Background */}
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[600px] bg-primary/10 rounded-full blur-[120px] pointer-events-none"></div>
          
          {/* Background pattern */}
          <div className="absolute inset-0 bg-[linear-gradient(to_right,var(--border)_1px,transparent_1px),linear-gradient(to_bottom,var(--border)_1px,transparent_1px)] bg-[size:14px_24px] opacity-40"></div>
          
          <div className="container mx-auto relative z-10 text-center lg:max-w-5xl">
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold tracking-tighter mb-6 bg-clip-text text-transparent bg-gradient-to-r from-foreground to-foreground/70">
              Nền tảng Quản lý Tập trung dành cho Trường học
            </h1>
            <p className="mx-auto max-w-[700px] text-lg md:text-xl text-muted-foreground mb-8 font-light">
              Iris School kết nối Quản trị viên, Giáo viên và Phụ huynh trên một hệ thống duy nhất. Tối ưu hóa quy trình, nâng cao trải nghiệm giáo dục.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/login" className="w-full sm:w-auto">
                <Button size="lg" className="w-full">
                  Bắt đầu ngay <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
              <Link href="#features" className="w-full sm:w-auto">
                <Button size="lg" variant="outline" className="w-full bg-background/50 backdrop-blur-sm">
                  Khám phá Tính năng
                </Button>
              </Link>
            </div>
          </div>
        </section>

        {/* ─── Features Grid ────────────────────────────────────────────────────────── */}
        <section id="features" className="py-20 md:py-32 bg-background px-4 transition-colors duration-300">
          <div className="container mx-auto lg:max-w-7xl">
            <div className="text-center mb-16">
              <h2 className="text-3xl md:text-4xl font-bold tracking-tighter mb-4 text-foreground">Mọi công cụ bạn cần</h2>
              <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                Được thiết kế tỉ mỉ phục vụ riêng cho nghiệp vụ giáo dục cấp Mầm non và Tiểu học.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {/* Feature 1 */}
              <div className="flex flex-col p-6 rounded-2xl bg-muted/50 border border-border transition-all hover:shadow-md hover:border-primary/30">
                <div className="p-3 bg-muted w-fit rounded-lg mb-4 text-foreground transition-colors">
                  <ShieldCheck className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-foreground">Quản lý Đa điểm trường</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Super Admin kiểm soát hệ thống, School Admin quản trị trọn vẹn từng điểm trường với lớp học và học sinh riêng biệt.
                </p>
              </div>

              {/* Feature 2 */}
              <div className="flex flex-col p-6 rounded-2xl bg-muted/50 border border-border transition-all hover:shadow-md hover:border-primary/30">
                <div className="p-3 bg-muted w-fit rounded-lg mb-4 text-foreground transition-colors">
                  <ClipboardCheck className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-foreground">Điểm danh Nhanh chóng</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Giáo viên dễ dàng điểm danh đầu giờ bằng mã Check in / Check out, dữ liệu đồng bộ ngay lập tức đến phụ huynh.
                </p>
              </div>

              {/* Feature 3 */}
              <div className="flex flex-col p-6 rounded-2xl bg-muted/50 border border-border transition-all hover:shadow-md hover:border-primary/30">
                <div className="p-3 bg-muted w-fit rounded-lg mb-4 text-foreground transition-colors">
                  <HeartPulse className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-foreground">Sổ theo dõi Sức khỏe</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Ghi nhận chỉ số BMI định kỳ. Cập nhật các thông tin y tế, bữa ăn và theo dõi sát sao tình trạng các bé.
                </p>
              </div>

              {/* Feature 4 */}
              <div className="flex flex-col p-6 rounded-2xl bg-muted/50 border border-border transition-all hover:shadow-md hover:border-primary/30">
                <div className="p-3 bg-muted w-fit rounded-lg mb-4 text-foreground transition-colors">
                  <BookOpen className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-foreground">Bảng tin Thông báo</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Giáo viên đăng tải hình ảnh hoạt động, thực đơn, dặn dò. Phụ huynh nhận tương tác qua News Feed hiện đại.
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* ─── Footer ─────────────────────────────────────────────────────────────────── */}
      <footer className="bg-muted py-10 px-4 border-t border-border transition-colors duration-300">
        <div className="container mx-auto lg:max-w-7xl flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <GraduationCap className="h-5 w-5 text-muted-foreground" />
            <span className="text-sm font-semibold text-foreground">Iris School</span>
          </div>
          <p className="text-sm text-muted-foreground text-center">
            &copy; {new Date().getFullYear()} Nền tảng Quản lý Trường học.
          </p>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <a href="#" className="hover:text-foreground transition-colors">Bảo mật</a>
            <a href="#" className="hover:text-foreground transition-colors">Điều khoản</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
