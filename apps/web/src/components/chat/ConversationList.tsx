/**
 * ConversationList Component
 * Danh sách các cuộc hội thoại ở cột bên trái.
 * Hiển thị loading spinner, empty state, hoặc danh sách conversation buttons.
 */

import React from "react";
import { Loader2, MessageSquare } from "lucide-react";
import { Conversation } from "@/types";

interface ConversationListProps {
  loading: boolean;
  conversations: Conversation[];
  /** ID cuộc hội thoại đang được chọn (để highlight) */
  selectedConvId: string | undefined;
  /** Callback khi user click vào một conversation */
  onSelect: (conv: Conversation) => void;
  /** Hàm tính tên hiển thị của conversation */
  getConvDisplayName: (conv: Conversation) => string;
  /** Hàm tạo initials từ tên */
  getInitials: (name: string) => string;
}

export default function ConversationList({
  loading,
  conversations,
  selectedConvId,
  onSelect,
  getConvDisplayName,
  getInitials,
}: ConversationListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (conversations.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-muted-foreground space-y-4">
        <div className="h-16 w-16 rounded-full bg-muted flex items-center justify-center">
          <MessageSquare className="h-8 w-8 text-muted-foreground" />
        </div>
        <p className="text-sm font-medium">Chưa có tin nhắn nào</p>
      </div>
    );
  }

  return (
    <div className="space-y-1 pb-4">
      {conversations.map((conv) => {
        const displayName = getConvDisplayName(conv);
        const isSelected = selectedConvId === conv.conversation_id;

        return (
          <button
            key={conv.conversation_id}
            onClick={() => onSelect(conv)}
            className={`w-full text-left p-3 rounded-2xl flex items-center gap-3 transition-all ${
              isSelected
                ? "bg-primary text-primary-foreground shadow-md shadow-primary/20"
                : "hover:bg-secondary text-foreground"
            }`}
          >
            {/* Avatar */}
            <div
              className={`h-12 w-12 flex-shrink-0 rounded-full flex items-center justify-center font-bold text-lg shadow-sm ${
                isSelected
                  ? "bg-primary-foreground/20 text-primary-foreground"
                  : "bg-muted text-muted-foreground"
              }`}
            >
              {getInitials(displayName)}
            </div>

            {/* Info */}
            <div className="flex-1 overflow-hidden flex flex-col justify-center">
              <div className="flex items-center justify-between mb-0.5">
                <p className={`font-semibold truncate text-[15px] ${isSelected ? "text-primary-foreground" : ""}`}>
                  {displayName}
                </p>
              </div>
              <p
                className={`text-[13px] truncate ${
                  isSelected ? "text-primary-foreground/80" : "text-muted-foreground"
                }`}
              >
                {conv.type === "direct"
                  ? "Trò chuyện trực tiếp"
                  : `Nhóm ${conv.participants?.length || 0} thành viên`}
              </p>
            </div>
          </button>
        );
      })}
    </div>
  );
}
