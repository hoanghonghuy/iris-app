import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowRight, BookOpen, HeartPulse, ClipboardCheck, GraduationCap, ShieldCheck, Menu } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";
import { Sheet, SheetContent, SheetTrigger, SheetTitle } from "@/components/ui/sheet";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col bg-zinc-50 dark:bg-zinc-950 font-sans transition-colors duration-300">
      {/* ─── Header Navigation ──────────────────────────────────────────────────────── */}
      <header className="sticky top-0 z-50 w-full border-b border-zinc-200 dark:border-zinc-800 bg-white/80 dark:bg-zinc-950/80 backdrop-blur-md transition-colors duration-300">
        <div className="container mx-auto flex h-16 items-center justify-between px-4 md:px-6 lg:max-w-7xl">
          <div className="flex items-center gap-2">
            <GraduationCap className="h-7 w-7 text-zinc-900 dark:text-zinc-100" />
            <span className="text-xl font-bold tracking-tighter text-zinc-900 dark:text-zinc-100">Iris School</span>
          </div>
          <nav className="hidden md:flex gap-6 text-sm font-medium text-zinc-600 dark:text-zinc-400">
            <a href="#features" className="hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors">Tính năng</a>
            <a href="#about" className="hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors">Về chúng tôi</a>
          </nav>
          <div className="flex items-center gap-2 sm:gap-3">
            <ThemeToggle className="text-zinc-600 dark:text-zinc-400" />
            <div className="hidden sm:flex items-center gap-2 sm:gap-3">
              <Link href="/login">
                <Button variant="ghost" className="text-zinc-900 dark:text-zinc-100 hover:bg-zinc-100 dark:hover:bg-zinc-800">Đăng nhập</Button>
              </Link>
              <Link href="/register">
                <Button>Đăng ký Phụ huynh</Button>
              </Link>
            </div>
            {/* Mobile Navigation */}
            <div className="md:hidden flex items-center ml-2">
              <Sheet>
                <SheetTrigger asChild>
                  <Button variant="outline" size="icon" className="h-9 w-9 bg-transparent border-zinc-200 dark:border-zinc-800 text-zinc-900 dark:text-zinc-100">
                    <Menu className="h-5 w-5" />
                    <span className="sr-only">Toggle navigation menu</span>
                  </Button>
                </SheetTrigger>
                <SheetContent side="right" className="w-[300px] sm:w-[350px] bg-white dark:bg-zinc-950 border-l border-zinc-200 dark:border-zinc-800 flex flex-col p-0">
                  <SheetTitle className="sr-only">Menu Điều Hướng</SheetTitle>
                  <div className="flex flex-col h-full">
                    {/* Sidebar Header with Logo */}
                    <div className="p-6 pr-12 border-b border-zinc-100 dark:border-zinc-900 flex items-center gap-3">
                      <GraduationCap className="h-6 w-6 text-zinc-900 dark:text-zinc-100" />
                      <span className="text-xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">Iris School</span>
                    </div>
                    
                    <nav className="flex-1 flex flex-col py-8 px-6">
                      <div className="flex flex-col gap-2">
                        <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-400 dark:text-zinc-500 mb-2 px-2">Điều hướng</p>
                        <a 
                          href="#features" 
                          className="flex items-center gap-3 px-3 py-3 rounded-xl text-lg font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 hover:bg-zinc-50 dark:hover:bg-zinc-900 transition-all active:scale-95"
                        >
                          Tính năng
                        </a>
                        <a 
                          href="#about" 
                          className="flex items-center gap-3 px-3 py-3 rounded-xl text-lg font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 hover:bg-zinc-50 dark:hover:bg-zinc-900 transition-all active:scale-95"
                        >
                          Về chúng tôi
                        </a>
                      </div>
                      
                      <div className="mt-auto flex flex-col gap-4 pb-8">
                        <div className="h-px bg-zinc-100 dark:bg-zinc-900 w-full mb-4" />
                        <Link href="/login" className="w-full">
                          <Button 
                            variant="outline" 
                            className="w-full justify-center h-12 text-base font-semibold border-zinc-200 dark:border-zinc-800 rounded-xl"
                          >
                            Đăng nhập
                          </Button>
                        </Link>
                        <Link href="/register" className="w-full">
                          <Button 
                            className="w-full justify-center h-12 text-base font-semibold bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-950 hover:bg-zinc-800 dark:hover:bg-zinc-200 rounded-xl shadow-lg shadow-zinc-200/50 dark:shadow-none"
                          >
                            Đăng ký Phụ huynh
                          </Button>
                        </Link>
                      </div>
                    </nav>
                  </div>
                </SheetContent>
              </Sheet>
            </div>
          </div>
        </div>
      </header>

      {/* ─── Hero Section ───────────────────────────────────────────────────────────── */}
      <main className="flex-1">
        {/* Keeps a dark scheme for both modes as a prominent visual focal point (Dark Hero Pattern) */}
        <section className="relative px-4 py-24 md:py-32 lg:py-40 bg-zinc-900 text-white overflow-hidden">
          {/* Background pattern */}
          <div className="absolute inset-0 bg-[linear-gradient(to_right,#4f4f4f2e_1px,transparent_1px),linear-gradient(to_bottom,#4f4f4f2e_1px,transparent_1px)] bg-[size:14px_24px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)]"></div>
          
          <div className="container mx-auto relative z-10 text-center lg:max-w-5xl">
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold tracking-tighter mb-6 bg-clip-text text-transparent bg-gradient-to-r from-zinc-100 to-zinc-400">
              Nền tảng Quản lý Tập trung dành cho Trường học
            </h1>
            <p className="mx-auto max-w-[700px] text-lg md:text-xl text-zinc-400 mb-8 font-light">
              Iris School kết nối Quản trị viên, Giáo viên và Phụ huynh trên một hệ thống duy nhất. Tối ưu hóa quy trình, nâng cao trải nghiệm giáo dục.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/login" className="w-full sm:w-auto">
                <Button size="lg" className="w-full bg-white text-zinc-900 hover:bg-zinc-200">
                  Bắt đầu ngay <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
              <Link href="#features" className="w-full sm:w-auto">
                <Button size="lg" variant="outline" className="w-full border-zinc-700 text-zinc-300 bg-transparent hover:bg-zinc-800 hover:text-white">
                  Khám phá Tính năng
                </Button>
              </Link>
            </div>
          </div>
        </section>

        {/* ─── Features Grid ────────────────────────────────────────────────────────── */}
        <section id="features" className="py-20 md:py-32 bg-white dark:bg-zinc-950 px-4 transition-colors duration-300">
          <div className="container mx-auto lg:max-w-7xl">
            <div className="text-center mb-16">
              <h2 className="text-3xl md:text-4xl font-bold tracking-tighter mb-4 text-zinc-900 dark:text-zinc-100">Mọi công cụ bạn cần</h2>
              <p className="text-lg text-zinc-500 dark:text-zinc-400 max-w-2xl mx-auto">
                Được thiết kế tỉ mỉ phục vụ riêng cho nghiệp vụ giáo dục cấp Mầm non và Tiểu học.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {/* Feature 1 */}
              <div className="flex flex-col p-6 rounded-2xl bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 transition-all hover:shadow-md hover:border-zinc-300 dark:hover:border-zinc-700">
                <div className="p-3 bg-zinc-200/50 dark:bg-zinc-800/50 w-fit rounded-lg mb-4 text-zinc-800 dark:text-zinc-200 transition-colors">
                  <ShieldCheck className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-zinc-900 dark:text-zinc-100">Quản lý Đa điểm trường</h3>
                <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed">
                  Super Admin kiểm soát hệ thống, School Admin quản trị trọn vẹn từng điểm trường với lớp học và học sinh riêng biệt.
                </p>
              </div>

              {/* Feature 2 */}
              <div className="flex flex-col p-6 rounded-2xl bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 transition-all hover:shadow-md hover:border-zinc-300 dark:hover:border-zinc-700">
                <div className="p-3 bg-zinc-200/50 dark:bg-zinc-800/50 w-fit rounded-lg mb-4 text-zinc-800 dark:text-zinc-200 transition-colors">
                  <ClipboardCheck className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-zinc-900 dark:text-zinc-100">Điểm danh Nhanh chóng</h3>
                <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed">
                  Giáo viên dễ dàng điểm danh đầu giờ bằng mã Check in / Check out, dữ liệu đồng bộ ngay lập tức đến phụ huynh.
                </p>
              </div>

              {/* Feature 3 */}
              <div className="flex flex-col p-6 rounded-2xl bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 transition-all hover:shadow-md hover:border-zinc-300 dark:hover:border-zinc-700">
                <div className="p-3 bg-zinc-200/50 dark:bg-zinc-800/50 w-fit rounded-lg mb-4 text-zinc-800 dark:text-zinc-200 transition-colors">
                  <HeartPulse className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-zinc-900 dark:text-zinc-100">Sổ theo dõi Sức khỏe</h3>
                <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed">
                  Ghi nhận chỉ số BMI định kỳ. Cập nhật các thông tin y tế, bữa ăn và theo dõi sát sao tình trạng các bé.
                </p>
              </div>

              {/* Feature 4 */}
              <div className="flex flex-col p-6 rounded-2xl bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 transition-all hover:shadow-md hover:border-zinc-300 dark:hover:border-zinc-700">
                <div className="p-3 bg-zinc-200/50 dark:bg-zinc-800/50 w-fit rounded-lg mb-4 text-zinc-800 dark:text-zinc-200 transition-colors">
                  <BookOpen className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-bold mb-2 text-zinc-900 dark:text-zinc-100">Bảng tin Thông báo</h3>
                <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed">
                  Giáo viên đăng tải hình ảnh hoạt động, thực đơn, dặn dò. Phụ huynh nhận tương tác qua News Feed hiện đại.
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* ─── Footer ─────────────────────────────────────────────────────────────────── */}
      <footer className="bg-zinc-100 dark:bg-zinc-900 py-10 px-4 border-t border-zinc-200 dark:border-zinc-800 transition-colors duration-300">
        <div className="container mx-auto lg:max-w-7xl flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <GraduationCap className="h-5 w-5 text-zinc-600 dark:text-zinc-400" />
            <span className="text-sm font-semibold text-zinc-800 dark:text-zinc-200">Iris School</span>
          </div>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 text-center">
            &copy; {new Date().getFullYear()} Nền tảng Quản lý Trường học.
          </p>
          <div className="flex items-center gap-4 text-sm text-zinc-500 dark:text-zinc-400">
            <a href="#" className="hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors">Bảo mật</a>
            <a href="#" className="hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors">Điều khoản</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
