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
  const reconnectAttemptsRef = useRef(0);
  const connectRef = useRef<() => void>(() => undefined);
  const [isConnected, setIsConnected] = useState(false);

  const connect = useCallback(() => {
    const token = authHelpers.getToken();
    if (!token) return;

    // Đóng connection cũ nếu có
    if (wsRef.current) {
      wsRef.current.close();
    }

    // Ưu tiên token qua Sec-WebSocket-Protocol.
    // Đồng thời gửi thêm query param token để tương thích proxy không forward sub-protocol header.
    const wsUrl = `${WS_BASE_URL}/chat/ws?token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(wsUrl, ['Bearer', token]);


    ws.onopen = () => {
      console.log('[WS] Connected');
      setIsConnected(true);
      reconnectAttemptsRef.current = 0;
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
        if (reconnectAttemptsRef.current >= 5) {
          console.warn('[WS] Stop reconnect after 5 failed attempts');
          return;
        }

        reconnectAttemptsRef.current += 1;
        const delayMs = Math.min(3000 * reconnectAttemptsRef.current, 15000);
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log('[WS] Reconnecting... attempt', reconnectAttemptsRef.current);
          connectRef.current();
        }, delayMs);
      }
    };

    ws.onerror = () => {
      // Browser không expose chi tiết WebSocket error event (thường chỉ thấy {}).
      // Xem thêm thông tin ở onclose(code/reason) và backend logs.
      console.warn('[WS] Transport error detected (details hidden by browser)');
    };

    wsRef.current = ws;
  }, [onNewMessage]);

  useEffect(() => {
    connectRef.current = connect;
  }, [connect]);

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
