/**
 * Admin Student Detail Page
 * Trang chi tiết hồ sơ học sinh dành cho Admin.
 * Hiển thị thông tin cá nhân, lớp học và liên kết phụ huynh.
 * Thiết kế Modern Mix SaaS (Bento Grid).
 */
"use client";

import React, { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api/admin.api";
import { StudentProfile } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { 
  ArrowLeft, Loader2, User,
  HeartPulse, ClipboardList, Activity, Phone, Mail, GraduationCap
} from "lucide-react";
import { formatDateVN } from "@/lib/utils";
import { extractApiErrorMessage } from "@/lib/api-error";

const genderLabel: Record<string, string> = { male: "Nam", female: "Nữ", other: "Khác" };

type TimelineEvent = {
  id: string;
  timeLabel: string;
  category: "profile" | "classroom" | "contact" | "security";
  title: string;
  description: string;
};

const TIMELINE_BATCH_SIZE = 3;
const PARENT_BATCH_SIZE = 2;

export default function AdminStudentDetailPage() {
  const params = useParams<{ id: string }>();
  const studentId = params?.id;

  const [profile, setProfile] = useState<StudentProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [visibleTimelineCount, setVisibleTimelineCount] = useState(TIMELINE_BATCH_SIZE);
  const [visibleParentCount, setVisibleParentCount] = useState(PARENT_BATCH_SIZE);

  useEffect(() => {
    if (!studentId) return;
    const fetchProfile = async () => {
      try {
        setLoading(true);
        const data = await adminApi.getStudentProfile(studentId);
        setProfile(data);
      } catch (error: unknown) {
        setError(extractApiErrorMessage(error, "Không thể tải hồ sơ học sinh"));
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, [studentId]);

  useEffect(() => {
    setVisibleTimelineCount(TIMELINE_BATCH_SIZE);
    setVisibleParentCount(PARENT_BATCH_SIZE);
  }, [profile?.student_id]);

  const timelineEvents = useMemo<TimelineEvent[]>(() => {
    if (!profile) return [];

    const events: TimelineEvent[] = [
      {
        id: "profile-dob",
        timeLabel: formatDateVN(profile.dob),
        category: "profile",
        title: "Thông tin ngày sinh",
        description: `Hồ sơ ghi nhận ngày sinh ${formatDateVN(profile.dob)}.`,
      },
      {
        id: "classroom-status",
        timeLabel: "Hiện tại",
        category: "classroom",
        title: "Trạng thái lớp học",
        description: profile.current_class_name
          ? `Đang thuộc lớp ${profile.current_class_name}.`
          : "Hiện chưa được xếp lớp.",
      },
      {
        id: "parent-code-status",
        timeLabel: profile.code_expires_at ? formatDateVN(profile.code_expires_at) : "Hiện tại",
        category: "security",
        title: "Mã phụ huynh",
        description: profile.active_parent_code
          ? `Mã đang hoạt động${profile.code_expires_at ? `, hết hạn vào ${formatDateVN(profile.code_expires_at)}.` : "."}`
          : "Hiện chưa có mã phụ huynh đang hoạt động.",
      },
      {
        id: "parent-link-summary",
        timeLabel: "Hiện tại",
        category: "contact",
        title: "Tổng quan liên hệ phụ huynh",
        description: `Đã liên kết ${profile.parents.length} phụ huynh vào hồ sơ học sinh.`,
      },
    ];

    const parentEvents = profile.parents.map((parent) => ({
      id: `parent-${parent.parent_id}`,
      timeLabel: "Hiện tại",
      category: "contact" as const,
      title: `Liên kết phụ huynh: ${parent.full_name}`,
      description: [parent.phone || null, parent.email || null].filter(Boolean).join(" • ") || "Chưa có thông tin liên hệ.",
    }));

    return [...events, ...parentEvents];
  }, [profile]);

  const visibleTimelineEvents = timelineEvents.slice(0, visibleTimelineCount);
  const visibleParents = profile?.parents.slice(0, visibleParentCount) || [];

  const hasMoreTimeline = visibleTimelineCount < timelineEvents.length;
  const hasMoreParents = !!profile && visibleParentCount < profile.parents.length;

  const getTimelineTone = (category: TimelineEvent["category"]) => {
    switch (category) {
      case "profile":
        return { ring: "border-chart-2", text: "text-chart-2", label: "Hồ sơ" };
      case "classroom":
        return { ring: "border-primary", text: "text-primary", label: "Lớp học" };
      case "contact":
        return { ring: "border-success", text: "text-success", label: "Liên hệ" };
      default:
        return { ring: "border-chart-3", text: "text-chart-3", label: "Bảo mật" };
    }
  };

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
    <div className="space-y-4 max-w-7xl mx-auto pb-6 lg:pb-0 lg:h-[calc(100dvh-8rem)] lg:flex lg:flex-col lg:overflow-hidden">
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

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 lg:gap-6 items-start lg:flex-1 lg:min-h-0">
        {/* LỚP BÊN TRÁI: IDENTIFICATION COLUMN */}
        <div className="lg:col-span-4 space-y-4 lg:space-y-5 lg:h-full lg:min-h-0 lg:grid lg:grid-rows-[auto,minmax(0,1fr)]">
          {/* Main ID Card */}
          <Card className="border-border/60 shadow-sm overflow-hidden relative">
            {/* Color banner */}
            <div className="h-20 bg-gradient-to-r from-primary/30 to-primary/20 w-full absolute top-0 left-0" />
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
          <Card className="border-border/60 shadow-sm overflow-hidden lg:h-full lg:flex lg:flex-col">
            <CardHeader className="bg-muted/20 border-b border-border/50 py-4 px-5">
              <CardTitle className="text-sm font-semibold flex items-center gap-2">
                <User className="h-4 w-4 text-primary" /> Thông tin Liên hệ
              </CardTitle>
            </CardHeader>
            <CardContent className="p-0 lg:flex lg:flex-col lg:flex-1 lg:min-h-0">
              {profile.parents && profile.parents.length > 0 ? (
                <>
                  <div className="divide-y divide-border/50 lg:flex-1 lg:min-h-0 lg:overflow-y-auto overscroll-contain">
                  {visibleParents.map((parent, idx) => (
                    <div key={parent.parent_id || idx} className="p-5 flex flex-col gap-3 hover:bg-muted/10 transition-colors">
                      <div className="flex items-center gap-3">
                         <div className="h-10 w-10 rounded-full bg-primary/10 text-primary flex items-center justify-center text-sm font-semibold uppercase">
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
                  {hasMoreParents && (
                    <div className="border-t border-border/60 p-3 shrink-0 bg-card">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="w-full"
                        onClick={() => setVisibleParentCount((prev) => prev + PARENT_BATCH_SIZE)}
                      >
                        Xem thêm liên hệ
                      </Button>
                    </div>
                  )}
                </>
              ) : (
                <div className="p-8 text-center text-muted-foreground">
                  <p className="text-sm">Chưa có phụ huynh nào liên kết</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* LỚP BÊN PHẢI: ACTIVITY & LOGS COLUMN */}
          <div className="lg:col-span-8 flex flex-col gap-4 lg:gap-5 lg:h-full lg:min-h-0">
           {/* Bento Stats Row */}
            <div className="grid grid-cols-3 gap-2 sm:grid-cols-2 sm:gap-3 xl:grid-cols-3">
              <Card className="border-border/60 shadow-sm bg-card hover:bg-muted/10 transition-colors sm:col-span-1">
                <CardContent className="p-3 sm:p-5 flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-4">
                  <div className="h-8 w-8 sm:h-12 sm:w-12 rounded-lg sm:rounded-xl bg-chart-2/10 text-chart-2 flex items-center justify-center">
                    <ClipboardList className="h-4 w-4 sm:h-6 sm:w-6" />
                  </div>
                  <div>
                    <h3 className="text-xs sm:text-sm font-medium text-muted-foreground">Liên hệ</h3>
                    <div className="flex items-baseline gap-1 mt-0.5 sm:mt-1">
                      <span className="text-lg sm:text-2xl font-bold tracking-tight text-foreground">{profile.parents.length}</span>
                      <span className="hidden sm:inline text-xs text-muted-foreground font-medium">phụ huynh</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border/60 shadow-sm bg-card hover:bg-muted/10 transition-colors sm:col-span-1">
                <CardContent className="p-3 sm:p-5 flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-4">
                  <div className="h-8 w-8 sm:h-12 sm:w-12 rounded-lg sm:rounded-xl bg-success/10 text-success flex items-center justify-center">
                    <HeartPulse className="h-4 w-4 sm:h-6 sm:w-6" />
                  </div>
                  <div>
                    <h3 className="text-xs sm:text-sm font-medium text-muted-foreground">Mã phụ huynh</h3>
                    <div className="flex items-baseline gap-1 mt-0.5 sm:mt-1">
                      <span className="text-sm sm:text-2xl font-bold tracking-tight text-foreground">
                        {profile.active_parent_code ? "Đang mở" : "Chưa có"}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border/60 shadow-sm bg-card hover:bg-muted/10 transition-colors sm:col-span-2 xl:col-span-1">
                <CardContent className="p-3 sm:p-5 flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-4">
                  <div className="h-8 w-8 sm:h-12 sm:w-12 rounded-lg sm:rounded-xl bg-chart-3/10 text-chart-3 flex items-center justify-center">
                    <Activity className="h-4 w-4 sm:h-6 sm:w-6" />
                  </div>
                  <div>
                    <h3 className="text-xs sm:text-sm font-medium text-muted-foreground">Lớp hiện tại</h3>
                    <div className="flex items-baseline gap-1 mt-0.5 sm:mt-1">
                      <span className="text-lg sm:text-2xl font-bold tracking-tight text-foreground">
                        {profile.current_class_name ? "Đã xếp" : "Chưa xếp"}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
           </div>

           {/* Timeline */}
           <Card className="flex-1 border-border/60 shadow-sm overflow-hidden lg:flex lg:flex-col lg:min-h-0">
             <CardHeader className="flex flex-row items-center justify-between border-b border-border/50 py-4 px-6 bg-muted/10">
               <CardTitle className="text-base font-semibold">Nhật ký Hoạt động (Timeline)</CardTitle>
               <span className="text-xs text-muted-foreground">Hiển thị {visibleTimelineEvents.length}/{timelineEvents.length}</span>
             </CardHeader>
             <CardContent className="p-0 lg:flex lg:flex-col lg:flex-1 lg:min-h-0">
               <div className="p-4 sm:p-6 lg:flex-1 lg:min-h-0 lg:overflow-y-auto overscroll-contain">
                 <div className="border-l-2 border-border/50 ml-3 sm:ml-4 pl-6 sm:pl-8 py-2 relative space-y-6">
                    {visibleTimelineEvents.map((event) => {
                      const tone = getTimelineTone(event.category);
                      return (
                        <div key={event.id} className="relative">
                          <span className={`absolute -left-[33px] sm:-left-[45px] top-1 h-5 w-5 sm:h-6 sm:w-6 rounded-full bg-card border-4 ${tone.ring}`} />
                          <div className="bg-muted/30 p-3 sm:p-4 rounded-xl border border-border/50">
                            <p className="text-xs text-muted-foreground mb-1">
                              {event.timeLabel}
                              <span className="mx-2">•</span>
                              <span className={`font-medium ${tone.text}`}>{tone.label}</span>
                            </p>
                            <p className="text-sm font-medium text-foreground">{event.title}</p>
                            <p className="text-sm text-foreground/80 mt-1">{event.description}</p>
                          </div>
                        </div>
                      );
                    })}
                 </div>
               </div>
               {hasMoreTimeline && (
                 <div className="p-4 sm:p-5 border-t border-border/60 shrink-0 bg-card">
                   <Button
                     variant="outline"
                     size="sm"
                     className="w-full"
                     onClick={() => setVisibleTimelineCount((prev) => prev + TIMELINE_BATCH_SIZE)}
                   >
                     Xem thêm hoạt động
                   </Button>
                 </div>
               )}
             </CardContent>
           </Card>
        </div>
      </div>
    </div>
  );
}
