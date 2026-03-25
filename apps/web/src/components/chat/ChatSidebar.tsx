/**
 * ChatSidebar Component
 * Cột trái của giao diện chat: header (tên + trạng thái kết nối + nút tạo mới),
 * panel tìm kiếm user, và danh sách hội thoại.
 */

import React, { useState } from "react";
import { Plus, X } from "lucide-react";
import { Conversation } from "@/types";
import ConversationList from "./ConversationList";
import NewConversationPanel from "./NewConversationPanel";

interface SearchUser {
  user_id: string;
  email: string;
  full_name: string;
}

interface ChatSidebarProps {
  loading: boolean;
  conversations: Conversation[];
  selectedConvId: string | undefined;
  isConnected: boolean;
  onSelect: (conv: Conversation) => void;
  getConvDisplayName: (conv: Conversation) => string;
  getInitials: (name: string) => string;
  /** Kết quả tìm kiếm user (debounce từ parent) */
  searchResults: SearchUser[];
  newConvTarget: string;
  onNewConvTargetChange: (val: string) => void;
  onStartConversation: (userId: string) => void;
}

export default function ChatSidebar({
  loading,
  conversations,
  selectedConvId,
  isConnected,
  onSelect,
  getConvDisplayName,
  getInitials,
  searchResults,
  newConvTarget,
  onNewConvTargetChange,
  onStartConversation,
}: ChatSidebarProps) {
  const [showNewConv, setShowNewConv] = useState(false);

  const handleToggleNewConv = () => {
    setShowNewConv((prev) => !prev);
    // Reset input khi đóng panel
    if (showNewConv) onNewConvTargetChange("");
  };

  return (
    <div className="flex w-full md:w-80 flex-col border-r border-border bg-background z-0 h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-4 md:py-5 border-b border-transparent">
        <h2 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">
          Tin nhắn
        </h2>
        <div className="flex items-center gap-3">
          {/* Trạng thái kết nối WebSocket */}
          {isConnected ? (
            <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950/30">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
                <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500" />
              </span>
              <span className="text-[10px] font-bold text-emerald-600 dark:text-emerald-500 uppercase tracking-widest hidden md:inline-block">
                Online
              </span>
            </div>
          ) : (
            <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-zinc-100 dark:bg-zinc-900">
              {/* WifiOff icon inline để tránh import thêm prop */}
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-3 w-3 text-zinc-500"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth={2}
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <line x1="1" y1="1" x2="23" y2="23" />
                <path d="M16.72 11.06A10.94 10.94 0 0 1 19 12.55" />
                <path d="M5 12.55a10.94 10.94 0 0 1 5.17-2.39" />
                <path d="M10.71 5.05A16 16 0 0 1 22.56 9" />
                <path d="M1.42 9a15.91 15.91 0 0 1 4.7-2.88" />
                <path d="M8.53 16.11a6 6 0 0 1 6.95 0" />
                <line x1="12" y1="20" x2="12.01" y2="20" />
              </svg>
              <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest hidden md:inline-block">
                Offline
              </span>
            </div>
          )}

          {/* Nút tạo cuộc hội thoại mới */}
          <button
            onClick={handleToggleNewConv}
            className="rounded-full bg-primary/10 p-2 text-primary hover:bg-primary/20 transition-colors"
            title="Tạo mới"
          >
            {showNewConv ? <X className="h-5 w-5" /> : <Plus className="h-5 w-5" />}
          </button>
        </div>
      </div>

      {/* Panel tạo hội thoại mới */}
      {showNewConv && (
        <NewConversationPanel
          value={newConvTarget}
          onChange={onNewConvTargetChange}
          searchResults={searchResults}
          onStartConversation={(userId) => {
            onStartConversation(userId);
            setShowNewConv(false);
          }}
          getInitials={getInitials}
        />
      )}

      {/* Danh sách hội thoại */}
      <div className="flex-1 overflow-y-auto custom-scrollbar px-2">
        <ConversationList
          loading={loading}
          conversations={conversations}
          selectedConvId={selectedConvId}
          onSelect={onSelect}
          getConvDisplayName={getConvDisplayName}
          getInitials={getInitials}
        />
      </div>
    </div>
  );
}
