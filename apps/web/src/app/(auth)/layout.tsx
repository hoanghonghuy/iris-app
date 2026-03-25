import { ThemeToggle } from "@/components/ThemeToggle";
import { GraduationCap } from "lucide-react";
import Link from "next/link";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col bg-background transition-colors duration-300 relative">
      <div className="absolute top-4 right-4 sm:top-6 sm:right-6">
        <ThemeToggle className="text-muted-foreground" />
      </div>
      <Link href="/" className="absolute top-4 left-4 sm:top-6 sm:left-6 flex items-center gap-2 hover:opacity-80 transition-opacity">
        <GraduationCap className="h-6 w-6 text-foreground" />
        <span className="font-bold text-foreground hidden sm:inline-block">Iris School</span>
      </Link>
      <div className="flex flex-1 items-center justify-center p-4">
        {children}
      </div>
    </div>
  );
}
