/**
 * ChatPage Component
 * Trang chat realtime dùng chung cho tất cả roles (Admin, Teacher, Parent).
 * Đây là orchestrator layer: quản lý state, hooks, data-fetching.
 * Giao diện được delegate sang ChatSidebar và ChatArea.
 */
"use client";

import React from "react";
import { getInitials } from "@/components/chat/chatHelpers";
import ChatSidebar from "@/components/chat/ChatSidebar";
import ChatArea from "@/components/chat/ChatArea";
import { useChatPageController } from "@/hooks/useChatPageController";

export default function ChatPage() {
  const {
    conversations,
    selectedConv,
    messages,
    input,
    loading,
    currentUserId,
    newConvTarget,
    searchResults,
    hasMore,
    loadingMore,
    isConnected,
    messagesEndRef,
    messagesContainerRef,
    setSelectedConv,
    setInput,
    setNewConvTarget,
    handleSend,
    handleKeyDown,
    handleScroll,
    getConvDisplayName,
    startConversation,
  } = useChatPageController();

  /* ────────────────────────────── render ────────────────────────────── */

  return (
    <div className="relative flex h-full w-full bg-background overflow-hidden rounded-xl md:rounded-none md:border-l md:border-border">
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
