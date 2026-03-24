/**
 * ChatPage Component
 * Trang chat realtime dùng chung cho tất cả roles (Admin, Teacher, Parent).
 * Đây là orchestrator layer: quản lý state, hooks, data-fetching.
 * Giao diện được delegate sang ChatSidebar và ChatArea.
 */
"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";
import { chatApi } from "@/lib/api/chat.api";
import { authHelpers } from "@/lib/api/client";
import { useChatWebSocket } from "@/hooks/useChatWebSocket";
import { Conversation, Message } from "@/types";
import { parseJwtPayload, getInitials } from "@/components/chat/chatHelpers";
import ChatSidebar from "@/components/chat/ChatSidebar";
import ChatArea from "@/components/chat/ChatArea";

export default function ChatPage() {
  /* ── state ── */
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConv, setSelectedConv] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(true);
  const [currentUserId, setCurrentUserId] = useState<string>("");
  const [newConvTarget, setNewConvTarget] = useState("");
  const [searchResults, setSearchResults] = useState<
    { user_id: string; email: string; full_name: string }[]
  >([]);

  /* ── cursor-based pagination ── */
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
      if (prev.length > 0 && prev[0]?.conversation_id === msg.conversation_id) {
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

    if (wasAtBottom || prevMessageCountRef.current === 0) {
      messagesEndRef.current?.scrollIntoView({
        behavior: prevMessageCountRef.current === 0 ? "auto" : "smooth",
      });
    }
    prevMessageCountRef.current = newCount;
  }, [messages]);

  /* ── load more khi user cuộn lên đầu (infinite scroll ngược) ── */
  const handleScroll = useCallback(async () => {
    const container = messagesContainerRef.current;
    if (!container || !hasMore || loadingMore || !selectedConv || !nextCursor) return;

    if (container.scrollTop < 80) {
      setLoadingMore(true);
      const prevScrollHeight = container.scrollHeight;
      try {
        const res = await chatApi.listMessages(selectedConv.conversation_id, 50, nextCursor);
        const older = [...(res.data || [])].reverse();
        setMessages((prev) => [...older, ...prev]);
        setNextCursor(res.next_cursor);
        setHasMore(res.has_more);

        // Giữ nguyên vị trí scroll sau khi prepend tin nhắn cũ
        requestAnimationFrame(() => {
          if (container) container.scrollTop = container.scrollHeight - prevScrollHeight;
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

  /* ── tìm kiếm user (debounce 500ms) ── */
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

  /* ── tạo cuộc hội thoại mới ── */
  const startConversation = async (targetUserId: string) => {
    if (!targetUserId) return;
    try {
      const newConv = await chatApi.createDirectConversation(targetUserId);
      await fetchConversations();
      setSelectedConv(newConv);
      setNewConvTarget("");
      setSearchResults([]);
    } catch (err) {
      console.error(err);
      alert("Không thể tạo cuộc trò chuyện: " + (err as Error).message);
    }
  };

  /* ────────────────────────────── render ────────────────────────────── */

  return (
    <div className="relative flex h-full w-full bg-background overflow-hidden rounded-xl md:rounded-none md:border-l md:border-zinc-200 md:dark:border-zinc-800">
      <ChatSidebar
        loading={loading}
        conversations={conversations}
        selectedConvId={selectedConv?.conversation_id}
        isConnected={isConnected}
        onSelect={setSelectedConv}
        getConvDisplayName={getConvDisplayName}
        getInitials={getInitials}
        searchResults={searchResults}
        newConvTarget={newConvTarget}
        onNewConvTargetChange={setNewConvTarget}
        onStartConversation={startConversation}
      />
      <ChatArea
        selectedConv={selectedConv}
        messages={messages}
        currentUserId={currentUserId}
        loadingMore={loadingMore}
        hasMore={hasMore}
        input={input}
        onInputChange={setInput}
        onSend={handleSend}
        onKeyDown={handleKeyDown}
        onScroll={handleScroll}
        onBack={() => setSelectedConv(null)}
        messagesEndRef={messagesEndRef}
        messagesContainerRef={messagesContainerRef}
      />
    </div>
  );
}
