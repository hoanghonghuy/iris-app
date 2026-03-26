/**
 * Admin Student Detail Page
 * Trang chi tiết hồ sơ học sinh dành cho Admin.
 * Hiển thị thông tin cá nhân, lớp học, điểm danh, sức khỏe và bài đăng (Timeline).
 * Thiết kế Modern Mix SaaS (Bento Grid).
 */
"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api/admin.api";
import { StudentProfile } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { 
  ArrowLeft, Loader2, User, Calendar, MapPin, 
  HeartPulse, ClipboardList, Activity, Phone, Mail, GraduationCap
} from "lucide-react";
import { formatDateVN } from "@/lib/utils";

const genderLabel: Record<string, string> = { male: "Nam", female: "Nữ", other: "Khác" };

export default function AdminStudentDetailPage() {
  const params = useParams<{ id: string }>();
  const studentId = params?.id;

  const [profile, setProfile] = useState<StudentProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!studentId) return;
    const fetchProfile = async () => {
      try {
        setLoading(true);
        const data = await adminApi.getStudentProfile(studentId);
        setProfile(data);
      } catch (err: any) {
        setError(err.response?.data?.error || "Không thể tải hồ sơ học sinh");
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, [studentId]);

  if (loading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="h-10 w-10 animate-spin text-primary/60" />
          <p className="text-sm text-muted-foreground animate-pulse">Đang tải hồ sơ học sinh...</p>
        </div>
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="space-y-4">
        <Link href="/admin/students">
          <Button variant="ghost" size="sm" className="gap-2">
            <ArrowLeft className="h-4 w-4" /> Quay lại danh sách
          </Button>
        </Link>
        <Card className="border-destructive/20 bg-destructive/5">
          <CardContent className="flex flex-col items-center py-10">
            <User className="h-12 w-12 text-destructive mb-4 opacity-50" />
            <p className="text-destructive font-medium">{error || "Không tìm thấy học sinh"}</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const ageMatch = profile.dob.match(/^(\d{4})-(\d{2})-(\d{2})/);
  let ageString = "";
  if (ageMatch) {
    const y = Number(ageMatch[1]);
    const m = Number(ageMatch[2]);
    const d = Number(ageMatch[3]);
    const birthDate = new Date(y, m - 1, d);
    const today = new Date();
    let yDiff = today.getFullYear() - birthDate.getFullYear();
    const mDiff = today.getMonth() - birthDate.getMonth();
    if (mDiff < 0 || (mDiff === 0 && today.getDate() < birthDate.getDate())) {
      yDiff--;
    }
    ageString = `(${yDiff} tuổi)`;
  }

  return (
    <div className="space-y-6 max-w-7xl mx-auto pb-10">
      {/* Top Header */}
      <div className="flex items-center justify-between">
        <Link href="/admin/students">
          <Button variant="outline" size="sm" className="gap-2 group shadow-sm bg-card hover:bg-muted/50 border-muted">
            <ArrowLeft className="h-4 w-4 transition-transform group-hover:-translate-x-1" />
            <span>Danh sách</span>
          </Button>
        </Link>
        <div className="flex gap-2 text-sm text-muted-foreground bg-primary/5 px-3 py-1.5 rounded-full border border-primary/10">
          <span className="font-mono text-xs text-primary font-medium tracking-tight">ID: {profile.student_id.split("-")[0].toUpperCase()}</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* LỚP BÊN TRÁI: IDENTIFICATION COLUMN */}
        <div className="lg:col-span-4 space-y-6">
          {/* Main ID Card */}
          <Card className="border-border/60 shadow-sm overflow-hidden relative">
            {/* Color banner */}
            <div className="h-20 bg-gradient-to-r from-primary/30 to-blue-500/20 w-full absolute top-0 left-0" />
            <CardContent className="pt-10 px-6 pb-6 relative z-10 flex flex-col items-center text-center">
              <div className="h-24 w-24 rounded-full bg-card border-4 border-card shadow-sm flex items-center justify-center mb-4">
                 {/* Fake Avatar */}
                 <div className="h-full w-full rounded-full bg-primary/10 text-primary flex items-center justify-center text-3xl font-bold uppercase">
                   {profile.full_name.charAt(0)}
                 </div>
              </div>
              <h1 className="text-2xl font-bold tracking-tight text-foreground">{profile.full_name}</h1>
              <p className="text-sm font-medium text-muted-foreground mt-1 mb-4 flex items-center gap-1.5 justify-center">
                <GraduationCap className="h-4 w-4 shrink-0" /> Lớp: {profile.current_class_name || "Chưa xếp lớp"}
              </p>
              
              <div className="w-full grid grid-cols-2 gap-2 mt-2">
                 <div className="bg-muted/40 rounded-lg p-3 text-left">
                   <p className="text-xs text-muted-foreground uppercase tracking-wider mb-1">Giới tính</p>
                   <p className="font-medium text-foreground">{genderLabel[profile.gender] || profile.gender}</p>
                 </div>
                 <div className="bg-muted/40 rounded-lg p-3 text-left">
                   <p className="text-xs text-muted-foreground uppercase tracking-wider mb-1">Ngày sinh</p>
                   <p className="font-medium text-foreground">{formatDateVN(profile.dob)} <span className="text-muted-foreground/80 font-normal">{ageString}</span></p>
                 </div>
              </div>
            </CardContent>
          </Card>

          {/* Parent Links Card */}
          <Card className="border-border/60 shadow-sm overflow-hidden">
            <CardHeader className="bg-muted/20 border-b border-border/50 py-4 px-5">
              <CardTitle className="text-sm font-semibold flex items-center gap-2">
                <User className="h-4 w-4 text-primary" /> Thông tin Liên hệ
              </CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              {profile.parents && profile.parents.length > 0 ? (
                <div className="divide-y divide-border/50">
                  {profile.parents.map((parent, idx) => (
                    <div key={parent.parent_id || idx} className="p-5 flex flex-col gap-3 hover:bg-muted/10 transition-colors">
                      <div className="flex items-center gap-3">
                         <div className="h-10 w-10 rounded-full bg-blue-500/10 text-blue-600 flex items-center justify-center text-sm font-semibold uppercase">
                           {parent.full_name.charAt(0)}
                         </div>
                         <div className="flex-1 min-w-0">
                           <p className="text-sm font-semibold text-foreground truncate">{parent.full_name}</p>
                           <p className="text-xs text-muted-foreground">Phụ huynh chính</p>
                         </div>
                      </div>
                      <div className="grid grid-cols-1 gap-2 pl-13">
                         {parent.phone && (
                           <div className="flex items-center gap-2 text-sm text-muted-foreground">
                             <Phone className="h-3.5 w-3.5" /> <span>{parent.phone}</span>
                           </div>
                         )}
                         {parent.email && (
                           <div className="flex items-center gap-2 text-sm text-muted-foreground truncate">
                             <Mail className="h-3.5 w-3.5" /> <span className="truncate">{parent.email}</span>
                           </div>
                         )}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="p-8 text-center text-muted-foreground">
                  <p className="text-sm">Chưa có phụ huynh nào liên kết</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* LỚP BÊN PHẢI: ACTIVITY & LOGS COLUMN */}
        <div className="lg:col-span-8 flex flex-col gap-6">
           {/* Bento Stats Row */}
           <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              <Card className="border-border/60 shadow-sm bg-card hover:bg-muted/10 transition-colors">
                <CardContent className="p-5 flex items-center gap-4">
                  <div className="h-12 w-12 rounded-xl bg-orange-100 dark:bg-orange-500/20 text-orange-500 flex items-center justify-center">
                    <ClipboardList className="h-6 w-6" />
                  </div>
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Điểm danh</h3>
                    <div className="flex items-baseline gap-1 mt-1">
                      <span className="text-2xl font-bold tracking-tight text-foreground">95%</span>
                      <span className="text-xs text-muted-foreground font-medium">/tháng</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border/60 shadow-sm bg-card hover:bg-muted/10 transition-colors">
                <CardContent className="p-5 flex items-center gap-4">
                  <div className="h-12 w-12 rounded-xl bg-green-100 dark:bg-green-500/20 text-green-500 flex items-center justify-center">
                    <HeartPulse className="h-6 w-6" />
                  </div>
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Sức khỏe</h3>
                    <div className="flex items-baseline gap-1 mt-1">
                      <span className="text-2xl font-bold tracking-tight text-foreground">Bình thường</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border/60 shadow-sm bg-card hover:bg-muted/10 transition-colors sm:col-span-2 lg:col-span-1">
                <CardContent className="p-5 flex items-center gap-4">
                  <div className="h-12 w-12 rounded-xl bg-blue-100 dark:bg-blue-500/20 text-blue-500 flex items-center justify-center">
                    <Activity className="h-6 w-6" />
                  </div>
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Hoạt động</h3>
                    <div className="flex items-baseline gap-1 mt-1">
                      <span className="text-2xl font-bold tracking-tight text-foreground">12</span>
                      <span className="text-xs text-muted-foreground font-medium">bài đăng</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
           </div>

           {/* Timeline/Activity Placeholder Blocks */}
           <Card className="flex-1 border-border/60 shadow-sm">
             <CardHeader className="flex flex-row items-center justify-between border-b border-border/50 py-4 px-6 bg-muted/10">
               <CardTitle className="text-base font-semibold">Nhật ký Hoạt động (Timeline)</CardTitle>
               <Button variant="ghost" size="sm" className="h-8 text-primary hover:text-primary hover:bg-primary/10">
                 Xem tất cả
               </Button>
             </CardHeader>
             <CardContent className="p-0">
               {/* Minimal Placeholder timeline */}
               <div className="p-8">
                 <div className="border-l-2 border-border/50 ml-4 pl-8 py-2 relative space-y-10">
                    <div className="relative">
                       <span className="absolute -left-[45px] top-1 h-6 w-6 rounded-full bg-card border-4 border-green-500" />
                       <div className="bg-muted/30 p-4 rounded-xl border border-border/50">
                          <p className="text-xs text-muted-foreground mb-1">Hôm qua, 08:30 <span className="mx-2">•</span> <span className="font-medium text-green-600 dark:text-green-400">Có mặt</span></p>
                          <p className="text-sm font-medium text-foreground">Điểm danh sáng</p>
                          <p className="text-xs text-muted-foreground mt-2">Giáo viên: Nguyễn Văn A ghi nhận</p>
                       </div>
                    </div>
                    <div className="relative">
                       <span className="absolute -left-[45px] top-1 h-6 w-6 rounded-full bg-card border-4 border-orange-500" />
                       <div className="bg-muted/30 p-4 rounded-xl border border-border/50">
                          <p className="text-xs text-muted-foreground mb-1">Tuần trước, 09:15 <span className="mx-2">•</span> <span className="font-medium text-orange-600 dark:text-orange-400">Sức khỏe</span></p>
                          <p className="text-sm font-medium text-foreground">Kiểm tra Y tế định kỳ</p>
                          <p className="text-sm text-foreground/80 mt-1">Chiều cao: 120cm, Cân nặng: 22kg. Sức khỏe bình thường.</p>
                       </div>
                    </div>
                 </div>
               </div>
             </CardContent>
           </Card>
        </div>
      </div>
    </div>
  );
}
