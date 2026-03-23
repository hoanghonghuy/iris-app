/**
 * ChatPage Component
 * Trang chat realtime dùng chung cho tất cả roles (Admin, Teacher, Parent).
 * Giao diện split-pane: danh sách hội thoại bên trái, khung chat bên phải.
 */
"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";
import { chatApi } from "@/lib/api/chat.api";
import { authHelpers } from "@/lib/api/client";
import { useChatWebSocket } from "@/hooks/useChatWebSocket";
import { Conversation, Message } from "@/types";
import { MessageSquare, Send, ArrowLeft, Wifi, WifiOff, Plus, X } from "lucide-react";

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
  const messagesEndRef = useRef<HTMLDivElement>(null);

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
      // Nếu đang không xem cuộc hội thoại → vẫn cập nhật nếu match
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

  /* ── fetch messages khi chọn conversation ── */
  useEffect(() => {
    if (!selectedConv) return;
    (async () => {
      try {
        const data = await chatApi.listMessages(selectedConv.conversation_id);
        setMessages(data || []);
      } catch (err) {
        console.error("Failed to load messages:", err);
      }
    })();
  }, [selectedConv]);

  /* ── auto scroll ── */
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

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
    }, 500); // 500ms debounce
    return () => clearTimeout(timer);
  }, [newConvTarget]);

  /* ── tạo cuộc hội thoại mới bằng user_id cụ thể ── */
  const startConversation = async (targetUserId: string) => {
    if (!targetUserId) return;
    try {
      const newConv = await chatApi.createDirectConversation(targetUserId);
      await fetchConversations(); // refresh list
      setSelectedConv(newConv);
      setShowNewConv(false);
      setNewConvTarget("");
      setSearchResults([]);
    } catch (err) {
      console.error(err);
      alert("Không thể tạo cuộc trò chuyện: " + (err as Error).message);
    }
  };

  /* ────────────────────────────── render ────────────────────────────── */

  return (
    <div className="flex h-[calc(100vh-3.5rem)] bg-white dark:bg-zinc-950">
      {/* ── Sidebar: Danh sách cuộc hội thoại ── */}
      <div
        className={`${
          selectedConv ? "hidden md:flex" : "flex"
        } w-full md:w-80 flex-col border-r border-zinc-200 dark:border-zinc-800`}
      >
        {/* Header */}
        <div className="flex items-center justify-between border-b border-zinc-200 dark:border-zinc-800 px-4 py-3">
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            Tin nhắn
          </h2>
          <div className="flex items-center gap-2">
            {isConnected ? (
              <Wifi className="h-4 w-4 text-emerald-500" />
            ) : (
              <WifiOff className="h-4 w-4 text-red-500" />
            )}
            <button
              onClick={() => setShowNewConv(!showNewConv)}
              className="rounded-md p-1.5 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
              title="Tạo mới"
            >
              {showNewConv ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4 text-blue-600 dark:text-blue-500" />}
            </button>
          </div>
        </div>

        {/* Form tạo mới */}
        {showNewConv && (
          <div className="p-3 border-b border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/50 relative">
            <input
              type="text"
              placeholder="Nhập tên hoặc email..."
              value={newConvTarget}
              onChange={(e) => setNewConvTarget(e.target.value)}
              className="w-full rounded-md border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 px-3 py-1.5 text-sm outline-none focus:border-blue-500"
            />
            {/* Dropdown tìm kiếm */}
            {searchResults.length > 0 && (
              <div className="absolute z-10 left-3 right-3 top-14 bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-md shadow-lg max-h-48 overflow-y-auto">
                {searchResults.map((user) => (
                  <button
                    key={user.user_id}
                    onClick={() => startConversation(user.user_id)}
                    className="w-full text-left px-3 py-2 text-sm hover:bg-zinc-100 dark:hover:bg-zinc-700 transition-colors border-b border-zinc-100 dark:border-zinc-700/50 last:border-0 flex flex-col gap-0.5"
                  >
                    <span className="font-medium text-zinc-900 dark:text-zinc-100">{user.full_name}</span>
                    <span className="text-xs text-zinc-500">{user.email}</span>
                  </button>
                ))}
              </div>
            )}
            {newConvTarget && searchResults.length === 0 && (
              <p className="text-xs text-zinc-500 mt-1">Đang tìm kiếm...</p>
            )}
          </div>
        )}

        {/* List */}
        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <div className="p-4 text-center text-zinc-500 dark:text-zinc-400 text-sm">
              Đang tải...
            </div>
          ) : conversations.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-zinc-400 dark:text-zinc-500 px-4">
              <MessageSquare className="h-12 w-12 mb-3" />
              <p className="text-sm text-center">Chưa có cuộc hội thoại nào</p>
            </div>
          ) : (
            conversations.map((conv) => (
              <button
                key={conv.conversation_id}
                onClick={() => setSelectedConv(conv)}
                className={`w-full text-left px-4 py-3 border-b border-zinc-100 dark:border-zinc-800/50 hover:bg-zinc-50 dark:hover:bg-zinc-900 transition-colors ${
                  selectedConv?.conversation_id === conv.conversation_id
                    ? "bg-zinc-100 dark:bg-zinc-800"
                    : ""
                }`}
              >
                <p className="text-sm font-medium text-zinc-900 dark:text-zinc-100 truncate">
                  {getConvDisplayName(conv)}
                </p>
                <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">
                  {conv.type === "direct" ? "Trực tiếp" : "Nhóm"} •{" "}
                  {conv.participants?.length || 0} thành viên
                </p>
              </button>
            ))
          )}
        </div>
      </div>

      {/* ── Chat area ── */}
      <div
        className={`${
          selectedConv ? "flex" : "hidden md:flex"
        } flex-1 flex-col`}
      >
        {selectedConv ? (
          <>
            {/* Chat header */}
            <div className="flex items-center gap-3 border-b border-zinc-200 dark:border-zinc-800 px-4 py-3">
              <button
                onClick={() => setSelectedConv(null)}
                className="md:hidden rounded-md p-1 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
              >
                <ArrowLeft className="h-5 w-5 text-zinc-900 dark:text-zinc-100" />
              </button>
              <div>
                <p className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                  {getConvDisplayName(selectedConv)}
                </p>
                <p className="text-xs text-zinc-500 dark:text-zinc-400">
                  {selectedConv.participants?.map((p) => p.email).join(", ")}
                </p>
              </div>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {messages.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-zinc-400 dark:text-zinc-500">
                  <MessageSquare className="h-10 w-10 mb-2" />
                  <p className="text-sm">Bắt đầu cuộc trò chuyện</p>
                </div>
              ) : (
                messages.map((msg) => {
                  const isMine = msg.sender_id === currentUserId;
                  return (
                    <div
                      key={msg.message_id}
                      className={`flex ${isMine ? "justify-end" : "justify-start"}`}
                    >
                      <div
                        className={`max-w-[70%] rounded-2xl px-4 py-2 ${
                          isMine
                            ? "bg-blue-600 text-white"
                            : "bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100"
                        }`}
                      >
                        {!isMine && (
                          <p className="text-xs font-medium text-zinc-500 dark:text-zinc-400 mb-0.5">
                            {msg.sender_email}
                          </p>
                        )}
                        <p className="text-sm whitespace-pre-wrap break-words">{msg.content}</p>
                        <p
                          className={`text-[10px] mt-1 ${
                            isMine ? "text-blue-200" : "text-zinc-400 dark:text-zinc-500"
                          }`}
                        >
                          {formatTime(msg.created_at)}
                        </p>
                      </div>
                    </div>
                  );
                })
              )}
              <div ref={messagesEndRef} />
            </div>

            {/* Input */}
            <div className="border-t border-zinc-200 dark:border-zinc-800 px-4 py-3">
              <div className="flex items-center gap-2">
                <input
                  type="text"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Nhập tin nhắn..."
                  className="flex-1 rounded-full border border-zinc-300 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-900 px-4 py-2 text-sm text-zinc-900 dark:text-zinc-100 placeholder-zinc-400 dark:placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button
                  onClick={handleSend}
                  disabled={!input.trim()}
                  className="rounded-full bg-blue-600 p-2 text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <Send className="h-4 w-4" />
                </button>
              </div>
            </div>
          </>
        ) : (
          /* Empty state khi chưa chọn conversation */
          <div className="flex flex-col items-center justify-center h-full text-zinc-400 dark:text-zinc-500">
            <MessageSquare className="h-16 w-16 mb-4" />
            <p className="text-lg font-medium">Chọn một cuộc hội thoại</p>
            <p className="text-sm mt-1">để bắt đầu nhắn tin</p>
          </div>
        )}
      </div>
    </div>
  );
}
