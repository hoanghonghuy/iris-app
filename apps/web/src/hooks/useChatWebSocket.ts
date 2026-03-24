'use client';

import { useEffect, useRef, useCallback, useState } from 'react';
import { authHelpers } from '@/lib/api/client';
import { WSEvent, Message } from '@/types';

const WS_BASE_URL = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1')
  .replace(/^http/, 'ws');

/**
 * useChatWebSocket - Hook quản lý kết nối WebSocket cho chat realtime.
 * Tự động kết nối/ngắt khi mount/unmount, hỗ trợ auto-reconnect.
 */
export function useChatWebSocket(onNewMessage: (msg: Message) => void) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout>>(undefined);
  const [isConnected, setIsConnected] = useState(false);

  const connect = useCallback(() => {
    const token = authHelpers.getToken();
    if (!token) return;

    // Đóng connection cũ nếu có
    if (wsRef.current) {
      wsRef.current.close();
    }

    // Truyền token qua Sec-WebSocket-Protocol thay vì query string
    // (tránh token bị log trong server access log / browser history)
    // Format: ['Bearer', '<jwt>'] → browser gửi: Sec-WebSocket-Protocol: Bearer, <jwt>
    const ws = new WebSocket(`${WS_BASE_URL}/chat/ws`, ['Bearer', token]);


    ws.onopen = () => {
      console.log('[WS] Connected');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      try {
        const wsEvent: WSEvent = JSON.parse(event.data);
        if (wsEvent.type === 'new_message') {
          onNewMessage(wsEvent.data as Message);
        }
      } catch (err) {
        console.error('[WS] Parse error:', err);
      }
    };

    ws.onclose = (event) => {
      console.log('[WS] Disconnected:', event.code, event.reason);
      setIsConnected(false);

      // Auto-reconnect sau 3 giây (trừ khi đóng chủ ý)
      if (event.code !== 1000) {
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log('[WS] Reconnecting...');
          connect();
        }, 3000);
      }
    };

    ws.onerror = (error) => {
      console.error('[WS] Error:', error);
    };

    wsRef.current = ws;
  }, [onNewMessage]);

  // Gửi tin nhắn qua WebSocket
  const sendMessage = useCallback((conversationId: string, content: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        conversation_id: conversationId,
        content,
      }));
    }
  }, []);

  // Kết nối khi mount, đóng khi unmount
  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close(1000, 'component unmount');
      }
    };
  }, [connect]);

  return { isConnected, sendMessage };
}
