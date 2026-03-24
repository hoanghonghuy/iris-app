/**
 * ChatPage Component
 * Trang chat realtime dùng chung cho tất cả roles (Admin, Teacher, Parent).
 * Giao diện split-pane: danh sách hội thoại bên trái, khung chat bên phải.
 *
 * Sử dụng cursor-based pagination: load 50 tin nhắn mới nhất ban đầu,
 * sau đó tải thêm khi user cuộn lên đầu (infinite scroll ngược chiều).
 */
"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";
import { chatApi } from "@/lib/api/chat.api";
import { authHelpers } from "@/lib/api/client";
import { useChatWebSocket } from "@/hooks/useChatWebSocket";
import { Conversation, Message } from "@/types";
import { MessageSquare, Send, ArrowLeft, Wifi, WifiOff, Plus, X, Loader2 } from "lucide-react";

/* ────────────────────────────── helpers ────────────────────────────── */

/** parseJwtPayload giải mã phần payload của JWT token (không xác thực chữ ký) */
function parseJwtPayload(token: string): { user_id: string; email: string } | null {
  try {
    const base64 = token.split(".")[1];
    const json = atob(base64.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(json);
  } catch {
    return null;
  }
}

/** formatTime hiển thị thời gian ngắn gọn (HH:mm) */
function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString("vi-VN", { hour: "2-digit", minute: "2-digit" });
}

/* ────────────────────────────── component ────────────────────────────── */

export default function ChatPage() {
  /* ── state ── */
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConv, setSelectedConv] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(true);
  const [currentUserId, setCurrentUserId] = useState<string>("");
  const [showNewConv, setShowNewConv] = useState(false);
  const [newConvTarget, setNewConvTarget] = useState("");
  const [searchResults, setSearchResults] = useState<{ user_id: string; email: string; full_name: string }[]>([]);

  // Cursor-based pagination state
  const [nextCursor, setNextCursor] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);

  /* ── lấy user id từ JWT ── */
  useEffect(() => {
    const token = authHelpers.getToken();
    if (token) {
      const payload = parseJwtPayload(token);
      if (payload) setCurrentUserId(payload.user_id);
    }
  }, []);

  /* ── WebSocket ── */
  const handleNewMessage = useCallback((msg: Message) => {
    setMessages((prev) => {
      // Chỉ thêm nếu tin nhắn thuộc cuộc hội thoại đang xem
      if (prev.length > 0 && prev[0]?.conversation_id === msg.conversation_id) {
        // Tránh duplicate
        if (prev.some((m) => m.message_id === msg.message_id)) return prev;
        return [...prev, msg];
      }
      return prev;
    });
  }, []);

  const { isConnected, sendMessage } = useChatWebSocket(handleNewMessage);

  /* ── fetch conversations ── */
  const fetchConversations = useCallback(async () => {
    try {
      setLoading(true);
      const data = await chatApi.listConversations();
      setConversations(data || []);
    } catch (err) {
      console.error("Failed to load conversations:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchConversations();
  }, [fetchConversations]);

  /* ── fetch messages khi chọn conversation (lần đầu) ── */
  useEffect(() => {
    if (!selectedConv) return;
    setMessages([]);
    setNextCursor(null);
    setHasMore(false);

    (async () => {
      try {
        const res = await chatApi.listMessages(selectedConv.conversation_id);
        // API trả về DESC (mới nhất trước) → reverse để hiển thị đúng thứ tự
        setMessages([...(res.data || [])].reverse());
        setNextCursor(res.next_cursor);
        setHasMore(res.has_more);
      } catch (err) {
        console.error("Failed to load messages:", err);
      }
    })();
  }, [selectedConv]);

  /* ── auto scroll xuống cuối khi có tin nhắn mới ── */
  const prevMessageCountRef = useRef(0);
  useEffect(() => {
    const container = messagesContainerRef.current;
    if (!container) return;

    const newCount = messages.length;
    const wasAtBottom =
      container.scrollHeight - container.scrollTop - container.clientHeight < 100;

    // Chỉ scroll xuống nếu user đang ở gần cuối (tin nhắn mới đến)
    // hoặc đây là lần đầu load (prevCount = 0)
    if (wasAtBottom || prevMessageCountRef.current === 0) {
      messagesEndRef.current?.scrollIntoView({ behavior: prevMessageCountRef.current === 0 ? "auto" : "smooth" });
    }
    prevMessageCountRef.current = newCount;
  }, [messages]);

  /* ── load more khi user cuộn lên đầu ── */
  const handleScroll = useCallback(async () => {
    const container = messagesContainerRef.current;
    if (!container || !hasMore || loadingMore || !selectedConv || !nextCursor) return;

    // Khi scroll gần đầu (< 80px)
    if (container.scrollTop < 80) {
      setLoadingMore(true);
      const prevScrollHeight = container.scrollHeight;
      try {
        const res = await chatApi.listMessages(selectedConv.conversation_id, 50, nextCursor);
        const older = [...(res.data || [])].reverse(); // reverse về ASC
        setMessages((prev) => [...older, ...prev]);
        setNextCursor(res.next_cursor);
        setHasMore(res.has_more);

        // Giữ nguyên vị trí scroll sau khi prepend tin nhắn cũ
        requestAnimationFrame(() => {
          if (container) {
            container.scrollTop = container.scrollHeight - prevScrollHeight;
          }
        });
      } catch (err) {
        console.error("Failed to load more messages:", err);
      } finally {
        setLoadingMore(false);
      }
    }
  }, [hasMore, loadingMore, selectedConv, nextCursor]);

  /* ── gửi tin nhắn ── */
  const handleSend = () => {
    if (!input.trim() || !selectedConv) return;
    sendMessage(selectedConv.conversation_id, input.trim());
    setInput("");
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  /* ── tên hiển thị của conversation ── */
  const getConvDisplayName = (conv: Conversation) => {
    if (conv.name) return conv.name;
    if (conv.participants) {
      const other = conv.participants.find((p) => p.user_id !== currentUserId);
      return other?.full_name || other?.email || "Cuộc hội thoại";
    }
    return "Cuộc hội thoại";
  };

  /* ── tìm kiếm user debounce ── */
  useEffect(() => {
    if (!newConvTarget.trim()) {
      setSearchResults([]);
      return;
    }
    const timer = setTimeout(async () => {
      try {
        const results = await chatApi.searchUsers(newConvTarget.trim());
        setSearchResults(results || []);
      } catch (err) {
        console.error("Failed to search users:", err);
      }
    }, 500);
    return () => clearTimeout(timer);
  }, [newConvTarget]);

  /* ── tạo cuộc hội thoại mới bằng user_id cụ thể ── */
  const startConversation = async (targetUserId: string) => {
    if (!targetUserId) return;
    try {
      const newConv = await chatApi.createDirectConversation(targetUserId);
      await fetchConversations();
      setSelectedConv(newConv);
      setShowNewConv(false);
      setNewConvTarget("");
      setSearchResults([]);
    } catch (err) {
      console.error(err);
      alert("Không thể tạo cuộc trò chuyện: " + (err as Error).message);
    }
  };

  /* ── Ảo hóa Avatar từ tên ── */
  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .substring(0, 2)
      .toUpperCase();
  };

  /* ────────────────────────────── render ────────────────────────────── */

  return (
    <div className="relative flex h-[100dvh] md:h-[calc(100vh-3.5rem)] bg-white dark:bg-zinc-950 overflow-hidden">
      {/* ── Sidebar: Danh sách cuộc hội thoại ── */}
      <div className="flex w-full md:w-80 flex-col border-r border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 z-0 h-full">
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-4 md:py-5 border-b border-transparent">
          <h2 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">
            Tin nhắn
          </h2>
          <div className="flex items-center gap-3">
            {isConnected ? (
              <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950/30">
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                </span>
                <span className="text-[10px] font-bold text-emerald-600 dark:text-emerald-500 uppercase tracking-widest hidden md:inline-block">Online</span>
              </div>
            ) : (
              <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-zinc-100 dark:bg-zinc-900">
                <WifiOff className="h-3 w-3 text-zinc-500" />
                <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest hidden md:inline-block">Offline</span>
              </div>
            )}
            <button
              onClick={() => setShowNewConv(!showNewConv)}
              className="rounded-full bg-blue-50 dark:bg-blue-900/30 p-2 text-blue-600 dark:text-blue-400 hover:bg-blue-100 dark:hover:bg-blue-900/50 transition-colors"
              title="Tạo mới"
            >
              {showNewConv ? <X className="h-5 w-5" /> : <Plus className="h-5 w-5" />}
            </button>
          </div>
        </div>

        {/* Form tạo mới */}
        {showNewConv && (
          <div className="px-4 pb-4 animate-in fade-in slide-in-from-top-2 duration-200">
            <div className="relative">
              <input
                type="text"
                placeholder="Tìm kiếm email hoặc tên..."
                value={newConvTarget}
                onChange={(e) => setNewConvTarget(e.target.value)}
                className="w-full rounded-2xl border-none bg-zinc-100 dark:bg-zinc-900 px-4 py-2.5 pl-10 text-sm text-zinc-900 dark:text-zinc-100 placeholder-zinc-500 focus:ring-2 focus:ring-blue-500 focus:bg-white dark:focus:bg-zinc-950 transition-all outline-none"
              />
              <MessageSquare className="absolute left-4 top-2.5 h-4 w-4 text-zinc-400" />
            </div>
            {/* Dropdown tìm kiếm */}
            {searchResults.length > 0 && (
              <div className="mt-2 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-xl overflow-hidden animate-in fade-in duration-200">
                {searchResults.map((user) => (
                  <button
                    key={user.user_id}
                    onClick={() => startConversation(user.user_id)}
                    className="w-full text-left px-4 py-3 hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors flex items-center gap-3 border-b border-zinc-100 dark:border-zinc-800 last:border-0"
                  >
                    <div className="h-10 w-10 flex-shrink-0 rounded-full bg-gradient-to-br from-blue-500 to-indigo-500 flex items-center justify-center text-white font-bold shadow-sm">
                      {getInitials(user.full_name || user.email)}
                    </div>
                    <div className="flex flex-col overflow-hidden">
                      <span className="font-semibold text-zinc-900 dark:text-zinc-100 truncate">{user.full_name}</span>
                      <span className="text-xs text-zinc-500 truncate">{user.email}</span>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        )}

        {/* List Conversations */}
        <div className="flex-1 overflow-y-auto custom-scrollbar px-2">
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <Loader2 className="h-6 w-6 animate-spin text-zinc-300 dark:text-zinc-700" />
            </div>
          ) : conversations.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-zinc-400 dark:text-zinc-600 space-y-4">
              <div className="h-16 w-16 rounded-full bg-zinc-100 dark:bg-zinc-900 flex items-center justify-center">
                <MessageSquare className="h-8 w-8 text-zinc-300 dark:text-zinc-700" />
              </div>
              <p className="text-sm font-medium">Chưa có tin nhắn nào</p>
            </div>
          ) : (
            <div className="space-y-1 pb-4">
              {conversations.map((conv) => {
                const displayName = getConvDisplayName(conv);
                const isSelected = selectedConv?.conversation_id === conv.conversation_id;
                
                return (
                  <button
                    key={conv.conversation_id}
                    onClick={() => setSelectedConv(conv)}
                    className={`w-full text-left p-3 rounded-2xl flex items-center gap-3 transition-all ${
                      isSelected
                        ? "bg-blue-600 text-white shadow-md shadow-blue-500/20"
                        : "hover:bg-zinc-100 dark:hover:bg-zinc-900 text-zinc-900 dark:text-zinc-100"
                    }`}
                  >
                    {/* Avatar */}
                    <div className={`h-12 w-12 flex-shrink-0 rounded-full flex items-center justify-center font-bold text-lg shadow-sm ${
                      isSelected 
                        ? "bg-white/20 text-white" 
                        : "bg-gradient-to-br from-zinc-200 to-zinc-300 dark:from-zinc-800 dark:to-zinc-700 text-zinc-600 dark:text-zinc-300"
                    }`}>
                      {getInitials(displayName)}
                    </div>
                    
                    {/* Info */}
                    <div className="flex-1 overflow-hidden flex flex-col justify-center">
                      <div className="flex items-center justify-between mb-0.5">
                        <p className={`font-semibold truncate text-[15px] ${isSelected ? "text-white" : ""}`}>
                          {displayName}
                        </p>
                      </div>
                      <p className={`text-[13px] truncate ${isSelected ? "text-blue-100" : "text-zinc-500 dark:text-zinc-400"}`}>
                        {conv.type === "direct" ? "Trò chuyện trực tiếp" : `Nhóm ${conv.participants?.length || 0} thành viên`}
                      </p>
                    </div>
                  </button>
                );
              })}
            </div>
          )}
        </div>
      </div>

      {/* ── Chat area (Telegram-like slide over for mobile) ── */}
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
              <button
                onClick={() => setSelectedConv(null)}
                className="md:hidden rounded-full p-2.5 text-blue-600 dark:text-blue-400 hover:bg-zinc-100 dark:hover:bg-zinc-900 transition-colors active:scale-95"
              >
                <ArrowLeft className="h-6 w-6" strokeWidth={2.5} />
              </button>
              
              <div className="h-10 w-10 flex-shrink-0 rounded-full flex items-center justify-center font-bold text-white bg-gradient-to-br from-blue-500 to-indigo-500 shadow-sm ml-1 md:ml-3">
                {getInitials(getConvDisplayName(selectedConv))}
              </div>
              
              <div className="flex-1 overflow-hidden cursor-pointer">
                <p className="text-base font-bold text-zinc-900 dark:text-zinc-100 truncate tracking-tight">
                  {getConvDisplayName(selectedConv)}
                </p>
                <p className="text-[12px] text-zinc-500 dark:text-zinc-400 truncate">
                  {selectedConv.type === "direct" ? "Đang trực tuyến" : `${selectedConv.participants?.length || 0} thành viên`}
                </p>
              </div>
            </div>

            {/* Messages Area */}
            <div
              ref={messagesContainerRef}
              onScroll={handleScroll}
              className="flex-1 overflow-y-auto p-4 space-y-4 bg-[#e5e5ea] dark:bg-[#000000] relative scroll-smooth"
            >
              {/* Load more indicator */}
              {loadingMore && (
                <div className="flex justify-center py-3 sticky top-0 z-10">
                  <div className="bg-white/80 dark:bg-zinc-900/80 backdrop-blur-sm shadow-sm rounded-full px-4 py-1.5 flex items-center gap-2 border border-zinc-200 dark:border-zinc-800">
                    <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
                    <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">Đang tải lịch sử...</span>
                  </div>
                </div>
              )}
              
              {!hasMore && messages.length > 0 && (
                <div className="flex justify-center py-4">
                  <div className="bg-zinc-200/60 dark:bg-zinc-800/60 rounded-xl px-4 py-1.5">
                    <p className="text-[11px] font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                      Khởi đầu cuộc trò chuyện
                    </p>
                  </div>
                </div>
              )}

              {/* Message List */}
              {messages.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-zinc-400 dark:text-zinc-500 opacity-60">
                  <div className="bg-zinc-200/50 dark:bg-zinc-800/50 p-6 rounded-3xl mb-4">
                    <MessageSquare className="h-10 w-10 text-zinc-500 dark:text-zinc-400" />
                  </div>
                  <p className="text-[15px] font-medium text-zinc-600 dark:text-zinc-400">Gửi lời chào đầu tiên!</p>
                </div>
              ) : (
                <div className="flex flex-col gap-1 pb-4">
                  {messages.map((msg, index) => {
                    const isMine = msg.sender_id === currentUserId;
                    
                    // Logic to determine bubble shapes (consecutive messages)
                    const prevMsg = index > 0 ? messages[index - 1] : null;
                    const nextMsg = index < messages.length - 1 ? messages[index + 1] : null;
                    
                    const isFirstInGroup = !prevMsg || prevMsg.sender_id !== msg.sender_id;
                    const isLastInGroup = !nextMsg || nextMsg.sender_id !== msg.sender_id;
                    
                    return (
                      <div
                        key={msg.message_id}
                        className={`flex w-full ${isMine ? "justify-end pl-12" : "justify-start pr-12"} ${isFirstInGroup ? "mt-3" : "mt-0.5"}`}
                      >
                        <div
                          className={`max-w-full px-3.5 py-2.5 shadow-sm relative group
                            ${isMine 
                              ? "bg-blue-600 text-white" 
                              : "bg-white dark:bg-zinc-900 border border-zinc-100 dark:border-zinc-800 text-zinc-900 dark:text-zinc-100"
                            }
                            ${/* Smart Border Radius for Telegram feel */ ""}
                            ${isFirstInGroup && isLastInGroup ? "rounded-2xl" : ""}
                            ${!isFirstInGroup && !isLastInGroup && isMine ? "rounded-l-2xl rounded-r-sm" : ""}
                            ${!isFirstInGroup && !isLastInGroup && !isMine ? "rounded-r-2xl rounded-l-sm" : ""}
                            ${isFirstInGroup && !isLastInGroup && isMine ? "rounded-t-2xl rounded-bl-2xl rounded-br-sm" : ""}
                            ${isFirstInGroup && !isLastInGroup && !isMine ? "rounded-t-2xl rounded-br-2xl rounded-bl-sm" : ""}
                            ${!isFirstInGroup && isLastInGroup && isMine ? "rounded-b-2xl rounded-tl-2xl rounded-tr-sm" : ""}
                            ${!isFirstInGroup && isLastInGroup && !isMine ? "rounded-b-2xl rounded-tr-2xl rounded-tl-sm" : ""}
                          `}
                        >
                          {/* Sender Name for group chats if not mine */}
                          {!isMine && isFirstInGroup && selectedConv.type !== "direct" && (
                            <p className="text-[11px] font-bold text-blue-600 dark:text-blue-400 mb-1 leading-none">
                              {msg.sender_email.split('@')[0]}
                            </p>
                          )}
                          
                          <p className="text-[15px] leading-relaxed whitespace-pre-wrap break-words inline-block">
                            {msg.content}
                            {/* Inline spacer to push time to right */}
                            <span className="inline-block w-12" aria-hidden="true"></span>
                          </p>
                          
                          {/* Subtle Time Badge */}
                          <span
                            className={`absolute bottom-1.5 right-2 text-[10px] font-medium tracking-tight ${
                              isMine ? "text-blue-200" : "text-zinc-400 dark:text-zinc-600"
                            }`}
                          >
                            {formatTime(msg.created_at)}
                          </span>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
              <div ref={messagesEndRef} className="h-2" />
            </div>

            {/* Input Area */}
            <div className="shrink-0 bg-white dark:bg-zinc-950 px-3 pt-3 pb-[max(0.75rem,env(safe-area-inset-bottom))] md:px-4 md:py-4 border-t border-zinc-200 dark:border-zinc-900">
              <div className="flex items-end gap-2 max-w-4xl mx-auto">
                {/* Plus Button (Attachment placeholder) */}
                <button className="flex-shrink-0 p-2.5 rounded-full text-zinc-400 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-zinc-900 transition-colors">
                  <Plus className="h-6 w-6" strokeWidth={2} />
                </button>
                
                <div className="flex-1 relative bg-zinc-100 dark:bg-zinc-900 rounded-3xl border border-transparent focus-within:border-blue-500 focus-within:bg-white dark:focus-within:bg-zinc-950 transition-all flex items-end">
                  <textarea
                    rows={1}
                    value={input}
                    onChange={(e) => {
                      setInput(e.target.value);
                      e.target.style.height = 'auto';
                      e.target.style.height = `${Math.min(e.target.scrollHeight, 120)}px`;
                    }}
                    onKeyDown={handleKeyDown}
                    placeholder="Nhắn tin..."
                    className="w-full bg-transparent border-none px-4 py-3 text-[15px] text-zinc-900 dark:text-zinc-100 placeholder-zinc-500 focus:outline-none focus:ring-0 resize-none max-h-[120px] rounded-3xl"
                    style={{ minHeight: '44px' }}
                  />
                </div>
                
                {/* Send Button */}
                <button
                  onClick={handleSend}
                  disabled={!input.trim()}
                  className={`flex-shrink-0 p-3 rounded-full transition-all duration-200 ${
                    input.trim() 
                      ? "bg-blue-600 text-white shadow-md shadow-blue-600/20 active:scale-90 hover:bg-blue-700" 
                      : "bg-zinc-100 dark:bg-zinc-900 text-zinc-400"
                  }`}
                >
                  <Send className="h-5 w-5" strokeWidth={2} style={{ transform: input.trim() ? "translate(1px, -1px)" : "none" }} />
                </button>
              </div>
            </div>
          </>
        ) : (
          /* Desktop Empty State (Hidden on Mobile due to slide logic) */
          <div className="hidden md:flex flex-col w-full items-center justify-center h-full bg-[#f4f4f5] dark:bg-zinc-950 text-zinc-400 dark:text-zinc-600">
            <div className="w-24 h-24 rounded-full bg-white dark:bg-zinc-900 shadow-sm flex items-center justify-center mb-6">
              <MessageSquare className="h-10 w-10 text-blue-500 opacity-60" />
            </div>
            <p className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-zinc-600 to-zinc-400">Ứng dụng Nhắn tin IRIS</p>
            <p className="text-sm font-medium mt-2">Chọn một cuộc trò chuyện để bắt đầu</p>
          </div>
        )}
      </div>
    </div>
  );
}
