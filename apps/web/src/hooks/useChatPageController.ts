import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { chatApi } from "@/lib/api/chat.api";
import { authHelpers } from "@/lib/api/client";
import { useChatWebSocket } from "@/hooks/useChatWebSocket";
import { Conversation, Message } from "@/types";
import { parseJwtPayload } from "@/components/chat/chatHelpers";

type SearchUser = {
  user_id: string;
  email: string;
  full_name: string;
};

export function useChatPageController() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConv, setSelectedConv] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(true);
  const [currentUserId, setCurrentUserId] = useState<string>("");
  const [newConvTarget, setNewConvTarget] = useState("");
  const [searchResults, setSearchResults] = useState<SearchUser[]>([]);

  const [nextCursor, setNextCursor] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const prevMessageCountRef = useRef(0);

  useEffect(() => {
    const token = authHelpers.getToken();
    if (!token) {
      return;
    }

    const payload = parseJwtPayload(token);
    if (payload) {
      setCurrentUserId(payload.user_id);
    }
  }, []);

  const handleNewMessage = useCallback((msg: Message) => {
    setMessages((prev) => {
      if (prev.length > 0 && prev[0]?.conversation_id === msg.conversation_id) {
        if (prev.some((message) => message.message_id === msg.message_id)) {
          return prev;
        }
        return [...prev, msg];
      }
      return prev;
    });
  }, []);

  const { isConnected, sendMessage } = useChatWebSocket(handleNewMessage);

  const fetchConversations = useCallback(async () => {
    try {
      setLoading(true);
      const data = await chatApi.listConversations();
      setConversations(data || []);
    } catch (error) {
      console.error("Failed to load conversations:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchConversations();
  }, [fetchConversations]);

  const loadMessagesForSelectedConversation = useCallback(async (conversation: Conversation) => {
    setMessages([]);
    setNextCursor(null);
    setHasMore(false);

    try {
      const response = await chatApi.listMessages(conversation.conversation_id);
      setMessages([...(response.data || [])].reverse());
      setNextCursor(response.next_cursor);
      setHasMore(response.has_more);
    } catch (error) {
      console.error("Failed to load messages:", error);
    }
  }, []);

  useEffect(() => {
    if (!selectedConv) {
      return;
    }
    void loadMessagesForSelectedConversation(selectedConv);
  }, [selectedConv, loadMessagesForSelectedConversation]);

  useEffect(() => {
    const container = messagesContainerRef.current;
    if (!container) {
      return;
    }

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

  const handleScroll = useCallback(async () => {
    const container = messagesContainerRef.current;
    if (!container || !hasMore || loadingMore || !selectedConv || !nextCursor) {
      return;
    }

    if (container.scrollTop >= 80) {
      return;
    }

    setLoadingMore(true);
    const prevScrollHeight = container.scrollHeight;

    try {
      const response = await chatApi.listMessages(selectedConv.conversation_id, 50, nextCursor);
      const olderMessages = [...(response.data || [])].reverse();
      setMessages((prev) => [...olderMessages, ...prev]);
      setNextCursor(response.next_cursor);
      setHasMore(response.has_more);

      requestAnimationFrame(() => {
        if (container) {
          container.scrollTop = container.scrollHeight - prevScrollHeight;
        }
      });
    } catch (error) {
      console.error("Failed to load more messages:", error);
    } finally {
      setLoadingMore(false);
    }
  }, [hasMore, loadingMore, selectedConv, nextCursor]);

  const handleSend = useCallback(() => {
    if (!input.trim() || !selectedConv) {
      return;
    }

    sendMessage(selectedConv.conversation_id, input.trim());
    setInput("");
  }, [input, selectedConv, sendMessage]);

  const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSend();
    }
  }, [handleSend]);

  const getConvDisplayName = useCallback((conversation: Conversation) => {
    if (conversation.name) {
      return conversation.name;
    }

    if (conversation.participants) {
      const other = conversation.participants.find((participant) => participant.user_id !== currentUserId);
      return other?.full_name || other?.email || "Cuộc hội thoại";
    }

    return "Cuộc hội thoại";
  }, [currentUserId]);

  useEffect(() => {
    if (!newConvTarget.trim()) {
      setSearchResults([]);
      return;
    }

    const timer = setTimeout(async () => {
      try {
        const results = await chatApi.searchUsers(newConvTarget.trim());
        setSearchResults(results || []);
      } catch (error) {
        console.error("Failed to search users:", error);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [newConvTarget]);

  const startConversation = useCallback(async (targetUserId: string) => {
    if (!targetUserId) {
      return;
    }

    try {
      const newConversation = await chatApi.createDirectConversation(targetUserId);
      await fetchConversations();
      setSelectedConv(newConversation);
      setNewConvTarget("");
      setSearchResults([]);
    } catch (error) {
      console.error("Failed to create conversation:", error);
      toast.error("Không thể tạo cuộc trò chuyện. Vui lòng thử lại.");
    }
  }, [fetchConversations]);

  return {
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
  };
}