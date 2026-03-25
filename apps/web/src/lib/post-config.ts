import { PostScope, PostType } from "@/types";

export const POST_TYPE_META: Record<
  PostType,
  {
    label: string;
    badgeVariant: "default" | "secondary" | "outline" | "destructive";
  }
> = {
  announcement: { label: "Thông báo", badgeVariant: "default" },
  activity: { label: "Hoạt động", badgeVariant: "secondary" },
  daily_note: { label: "Nhận xét ngày", badgeVariant: "outline" },
  health_note: { label: "Sức khỏe", badgeVariant: "destructive" },
};

export const POST_SCOPE_LABELS: Record<PostScope, string> = {
  school: "Toàn trường",
  class: "Cả lớp",
  student: "Từng học sinh",
};

export const POST_TYPE_OPTIONS: Array<{ value: PostType; label: string }> = (
  Object.entries(POST_TYPE_META) as Array<[PostType, { label: string }]>
).map(([value, meta]) => ({
  value,
  label: meta.label,
}));
