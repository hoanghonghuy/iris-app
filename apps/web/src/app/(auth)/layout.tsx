import { ThemeToggle } from "@/components/ThemeToggle";
import { GraduationCap, ShieldCheck, Zap, Activity } from "lucide-react";
import Link from "next/link";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="grid min-h-screen w-full lg:grid-cols-2 bg-background relative transition-colors duration-300">
      
      {/* Left side: Branding/Graphic panel */}
      <div className="hidden lg:flex flex-col justify-between overflow-hidden bg-zinc-950 border-r border-zinc-800/50 text-white relative">
        {/* Subtle mesh/pattern background */}
        <div className="absolute inset-0 z-0">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,#4f4f4f2e_1px,transparent_1px),linear-gradient(to_bottom,#4f4f4f2e_1px,transparent_1px)] bg-[size:14px_24px] [mask-image:radial-gradient(ellipse_80%_50%_at_50%_0%,#000_70%,transparent_100%)]"></div>
          <div className="absolute -top-[20%] -left-[10%] w-[70%] h-[70%] rounded-full bg-primary/20 blur-[120px] pointer-events-none" />
          <div className="absolute bottom-[0%] -right-[20%] w-[60%] h-[60%] rounded-full bg-blue-600/10 blur-[100px] pointer-events-none" />
        </div>

        {/* Header inside left panel */}
        <div className="relative z-10 p-8 flex items-center gap-2">
          <Link href="/" className="flex items-center gap-2 hover:opacity-90 transition-opacity text-white">
            <GraduationCap className="h-7 w-7 text-primary" />
            <span className="font-bold text-xl tracking-tight">Iris School</span>
          </Link>
        </div>

        {/* Center content */}
        <div className="relative z-10 px-8 py-8 w-full max-w-2xl">
          <h1 className="text-3xl xl:text-4xl font-semibold tracking-tight leading-tight">
             Kết nối nhà trường và <span className="whitespace-nowrap">phụ huynh</span><br/> trong một nền tảng toàn diện.
          </h1>
          <p className="mt-6 text-zinc-400 text-base leading-relaxed max-w-xl">
             Hệ sinh thái công nghệ giáo dục Iris School mang lại sự tiện lợi, minh bạch và an toàn. Tương tác nhanh chóng, theo dõi hành trình học viên hiệu quả mỗi ngày.
          </p>

          <div className="mt-10 space-y-5">
             <div className="flex items-center gap-4 text-sm font-medium text-zinc-300">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-zinc-900 border border-zinc-800 text-blue-400 shadow-sm">
                  <ShieldCheck className="h-5 w-5" />
                </div>
                <span>Hệ thống bảo mật tối ưu</span>
             </div>
             <div className="flex items-center gap-4 text-sm font-medium text-zinc-300">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-zinc-900 border border-zinc-800 text-green-400 shadow-sm">
                  <Zap className="h-5 w-5" />
                </div>
                <span>Thông báo tương tác tức thời</span>
             </div>
             <div className="flex items-center gap-4 text-sm font-medium text-zinc-300">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-zinc-900 border border-zinc-800 text-purple-400 shadow-sm">
                  <Activity className="h-5 w-5" />
                </div>
                <span>Theo dõi sức khỏe & hoạt động</span>
             </div>
          </div>
        </div>

        {/* Footer inside left panel */}
        <div className="relative z-10 p-8">
          <p className="text-xs text-zinc-500 font-medium">&copy; {new Date().getFullYear()} Iris Education. All rights reserved.</p>
        </div>
      </div>

      {/* Right side: Auth Form Pane */}
      <div className="flex flex-col relative w-full h-[100dvh] lg:h-auto overflow-y-auto">
        {/* Mobile Header (Hidden on Desktop) */}
        <div className="flex lg:hidden items-center justify-between p-4 sm:p-6 pb-2">
          <Link href="/" className="flex items-center gap-2 hover:opacity-80 transition-opacity">
            <GraduationCap className="h-6 w-6 text-primary" />
            <span className="font-bold text-foreground">Iris School</span>
          </Link>
          <ThemeToggle className="text-muted-foreground" />
        </div>

        {/* Desktop Theme Toggle */}
        <div className="hidden lg:block absolute top-6 right-6 z-20">
          <ThemeToggle className="text-muted-foreground" />
        </div>

        {/* Form Container */}
        <div className="flex flex-1 items-center justify-center p-4 sm:p-8 md:p-12">
          {children}
        </div>
      </div>
    </div>
  );
}
