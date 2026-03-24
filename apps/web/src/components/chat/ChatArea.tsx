/**
 * ChatArea Component
 * Vùng bên phải của giao diện chat: header conversation, danh sách tin nhắn,
 * và ô nhập liệu. Nếu chưa chọn conversation, hiển thị empty state (desktop only).
 */

import React, { RefObject } from "react";
import { ArrowLeft, MessageSquare, Loader2, Send, Plus } from "lucide-react";
import { Conversation, Message } from "@/types";
import MessageBubble from "./MessageBubble";

interface ChatAreaProps {
  selectedConv: Conversation | null;
  messages: Message[];
  currentUserId: string;
  loadingMore: boolean;
  hasMore: boolean;
  input: string;
  onInputChange: (val: string) => void;
  onSend: () => void;
  onKeyDown: (e: React.KeyboardEvent) => void;
  onScroll: () => void;
  onBack: () => void;
  messagesEndRef: RefObject<HTMLDivElement | null>;
  messagesContainerRef: RefObject<HTMLDivElement | null>;
}

export default function ChatArea({
  selectedConv,
  messages,
  currentUserId,
  loadingMore,
  hasMore,
  input,
  onInputChange,
  onSend,
  onKeyDown,
  onScroll,
  onBack,
  messagesEndRef,
  messagesContainerRef,
}: ChatAreaProps) {
  return (
    <div
      className={`
        fixed inset-0 z-50 md:static md:flex-1
        flex flex-col bg-[#f4f4f5] dark:bg-zinc-950 md:z-auto
        transition-transform duration-300 ease-[cubic-bezier(0.32,0.72,0,1)]
        ${selectedConv ? "translate-x-0" : "translate-x-full md:translate-x-0"}
      `}
    >
      {selectedConv ? (
        <>
          {/* Chat Header */}
          <div className="shrink-0 flex items-center gap-3 bg-white/80 dark:bg-zinc-950/80 backdrop-blur-md border-b border-zinc-200 dark:border-zinc-800 px-2 py-2.5 z-10 shadow-sm md:shadow-none pt-[max(0.5rem,env(safe-area-inset-top))]">
            {/* Nút back (chỉ hiện trên mobile) */}
            <button
              onClick={onBack}
              className="md:hidden rounded-full p-2.5 text-blue-600 dark:text-blue-400 hover:bg-zinc-100 dark:hover:bg-zinc-900 transition-colors active:scale-95"
            >
              <ArrowLeft className="h-6 w-6" strokeWidth={2.5} />
            </button>

            {/* Avatar */}
            <div className="h-10 w-10 flex-shrink-0 rounded-full flex items-center justify-center font-bold text-white bg-gradient-to-br from-blue-500 to-indigo-500 shadow-sm ml-1 md:ml-3">
              {selectedConv.name
                ? selectedConv.name.substring(0, 2).toUpperCase()
                : (selectedConv.participants?.[0]?.full_name || "?").substring(0, 2).toUpperCase()}
            </div>

            {/* Tên và trạng thái */}
            <div className="flex-1 overflow-hidden cursor-pointer">
              <p className="text-base font-bold text-zinc-900 dark:text-zinc-100 truncate tracking-tight">
                {selectedConv.name ||
                  selectedConv.participants
                    ?.map((p) => p.full_name || p.email)
                    .join(", ") ||
                  "Cuộc hội thoại"}
              </p>
              <p className="text-[12px] text-zinc-500 dark:text-zinc-400 truncate">
                {selectedConv.type === "direct"
                  ? "Đang trực tuyến"
                  : `${selectedConv.participants?.length || 0} thành viên`}
              </p>
            </div>
          </div>

          {/* Vùng tin nhắn */}
          <div
            ref={messagesContainerRef}
            onScroll={onScroll}
            className="flex-1 overflow-y-auto p-4 space-y-4 bg-[#e5e5ea] dark:bg-[#000000] relative scroll-smooth"
          >
            {/* Indicator đang tải thêm */}
            {loadingMore && (
              <div className="flex justify-center py-3 sticky top-0 z-10">
                <div className="bg-white/80 dark:bg-zinc-900/80 backdrop-blur-sm shadow-sm rounded-full px-4 py-1.5 flex items-center gap-2 border border-zinc-200 dark:border-zinc-800">
                  <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
                  <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">
                    Đang tải lịch sử...
                  </span>
                </div>
              </div>
            )}

            {/* Nhãn "Khởi đầu cuộc trò chuyện" khi đã load hết */}
            {!hasMore && messages.length > 0 && (
              <div className="flex justify-center py-4">
                <div className="bg-zinc-200/60 dark:bg-zinc-800/60 rounded-xl px-4 py-1.5">
                  <p className="text-[11px] font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                    Khởi đầu cuộc trò chuyện
                  </p>
                </div>
              </div>
            )}

            {/* Danh sách tin nhắn hoặc empty state */}
            {messages.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-zinc-400 dark:text-zinc-500 opacity-60">
                <div className="bg-zinc-200/50 dark:bg-zinc-800/50 p-6 rounded-3xl mb-4">
                  <MessageSquare className="h-10 w-10 text-zinc-500 dark:text-zinc-400" />
                </div>
                <p className="text-[15px] font-medium text-zinc-600 dark:text-zinc-400">
                  Gửi lời chào đầu tiên!
                </p>
              </div>
            ) : (
              <div className="flex flex-col gap-1 pb-4">
                {messages.map((msg, index) => {
                  const isMine = msg.sender_id === currentUserId;
                  const prevMsg = index > 0 ? messages[index - 1] : null;
                  const nextMsg = index < messages.length - 1 ? messages[index + 1] : null;
                  const isFirstInGroup = !prevMsg || prevMsg.sender_id !== msg.sender_id;
                  const isLastInGroup = !nextMsg || nextMsg.sender_id !== msg.sender_id;
                  const showSenderName =
                    !isMine && isFirstInGroup && selectedConv.type !== "direct";

                  return (
                    <MessageBubble
                      key={msg.message_id}
                      msg={msg}
                      isMine={isMine}
                      isFirstInGroup={isFirstInGroup}
                      isLastInGroup={isLastInGroup}
                      showSenderName={showSenderName}
                    />
                  );
                })}
              </div>
            )}
            <div ref={messagesEndRef} className="h-2" />
          </div>

          {/* Input gửi tin nhắn */}
          <div className="shrink-0 bg-white dark:bg-zinc-950 px-3 pt-3 pb-[max(0.75rem,env(safe-area-inset-bottom))] md:px-4 md:py-4 border-t border-zinc-200 dark:border-zinc-900">
            <div className="flex items-end gap-2 max-w-4xl mx-auto">
              {/* Nút đính kèm (placeholder) */}
              <button className="flex-shrink-0 p-2.5 rounded-full text-zinc-400 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-zinc-900 transition-colors">
                <Plus className="h-6 w-6" strokeWidth={2} />
              </button>

              {/* Textarea tự mở rộng */}
              <div className="flex-1 relative bg-zinc-100 dark:bg-zinc-900 rounded-3xl border border-transparent focus-within:border-blue-500 focus-within:bg-white dark:focus-within:bg-zinc-950 transition-all flex items-end">
                <textarea
                  rows={1}
                  value={input}
                  onChange={(e) => {
                    onInputChange(e.target.value);
                    e.target.style.height = "auto";
                    e.target.style.height = `${Math.min(e.target.scrollHeight, 120)}px`;
                  }}
                  onKeyDown={onKeyDown}
                  placeholder="Nhắn tin..."
                  className="w-full bg-transparent border-none px-4 py-3 text-[15px] text-zinc-900 dark:text-zinc-100 placeholder-zinc-500 focus:outline-none focus:ring-0 resize-none max-h-[120px] rounded-3xl"
                  style={{ minHeight: "44px" }}
                />
              </div>

              {/* Nút gửi */}
              <button
                onClick={onSend}
                disabled={!input.trim()}
                className={`flex-shrink-0 p-3 rounded-full transition-all duration-200 ${
                  input.trim()
                    ? "bg-blue-600 text-white shadow-md shadow-blue-600/20 active:scale-90 hover:bg-blue-700"
                    : "bg-zinc-100 dark:bg-zinc-900 text-zinc-400"
                }`}
              >
                <Send
                  className="h-5 w-5"
                  strokeWidth={2}
                  style={{ transform: input.trim() ? "translate(1px, -1px)" : "none" }}
                />
              </button>
            </div>
          </div>
        </>
      ) : (
        /* Desktop Empty State (ẩn trên mobile vì slide logic đã xử lý) */
        <div className="hidden md:flex flex-col w-full items-center justify-center h-full bg-[#f4f4f5] dark:bg-zinc-950 text-zinc-400 dark:text-zinc-600">
          <div className="w-24 h-24 rounded-full bg-white dark:bg-zinc-900 shadow-sm flex items-center justify-center mb-6">
            <MessageSquare className="h-10 w-10 text-blue-500 opacity-60" />
          </div>
          <p className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-zinc-600 to-zinc-400">
            Ứng dụng Nhắn tin IRIS
          </p>
          <p className="text-sm font-medium mt-2">Chọn một cuộc trò chuyện để bắt đầu</p>
        </div>
      )}
    </div>
  );
}
