/**
 * MessageBubble Component
 * Hiển thị một bong bóng tin nhắn theo phong cách Telegram.
 * Hỗ trợ smart border radius dựa trên vị trí trong nhóm tin liên tiếp.
 */

import React from "react";
import { Message } from "@/types";
import { formatTime } from "./chatHelpers";

interface MessageBubbleProps {
  /** Tin nhắn cần render */
  msg: Message;
  /** true = tin nhắn của người dùng hiện tại (hiển thị bên phải) */
  isMine: boolean;
  /** true = tin nhắn đầu tiên trong chuỗi liên tiếp cùng sender */
  isFirstInGroup: boolean;
  /** true = tin nhắn cuối cùng trong chuỗi liên tiếp cùng sender */
  isLastInGroup: boolean;
  /** true = hiển thị tên sender (dùng cho group chat, không phải direct) */
  showSenderName: boolean;
}

export default function MessageBubble({
  msg,
  isMine,
  isFirstInGroup,
  isLastInGroup,
  showSenderName,
}: MessageBubbleProps) {
  // Smart border radius theo Telegram: góc nhỏ phía sender, góc lớn phía còn lại
  const radiusClass = (() => {
    if (isFirstInGroup && isLastInGroup) return "rounded-2xl";
    if (!isFirstInGroup && !isLastInGroup)
      return isMine ? "rounded-l-2xl rounded-r-sm" : "rounded-r-2xl rounded-l-sm";
    if (isFirstInGroup && !isLastInGroup)
      return isMine
        ? "rounded-t-2xl rounded-bl-2xl rounded-br-sm"
        : "rounded-t-2xl rounded-br-2xl rounded-bl-sm";
    // isLastInGroup && !isFirstInGroup
    return isMine
      ? "rounded-b-2xl rounded-tl-2xl rounded-tr-sm"
      : "rounded-b-2xl rounded-tr-2xl rounded-tl-sm";
  })();

  return (
    <div
      className={`flex w-full ${isMine ? "justify-end pl-12" : "justify-start pr-12"} ${
        isFirstInGroup ? "mt-3" : "mt-0.5"
      }`}
    >
      <div
        className={`max-w-full px-3.5 py-2.5 shadow-sm relative group ${
          isMine
            ? "bg-primary text-primary-foreground"
            : "bg-card border border-border text-foreground"
        } ${radiusClass}`}
      >
        {/* Tên sender cho group chat (chỉ hiện ở tin đầu nhóm, không phải tin của mình) */}
        {showSenderName && (
          <p className="text-[11px] font-bold text-primary mb-1 leading-none">
            {msg.sender_email.split("@")[0]}
          </p>
        )}

        <p className="text-[15px] leading-relaxed whitespace-pre-wrap break-words inline-block">
          {msg.content}
          {/* Spacer nội tuyến để giờ không đè lên nội dung cuối dòng */}
          <span className="inline-block w-12" aria-hidden="true" />
        </p>

        {/* Badge thời gian nhỏ ở góc phải dưới */}
        <span
          className={`absolute bottom-1.5 right-2 text-[10px] font-medium tracking-tight ${
            isMine ? "text-primary-foreground/70" : "text-muted-foreground"
          }`}
        >
          {formatTime(msg.created_at)}
        </span>
      </div>
    </div>
  );
}
