/**
 * NewConversationPanel Component
 * Form tìm kiếm user để bắt đầu cuộc hội thoại mới (appear khi bấm nút +).
 * Hiển thị input debounce + dropdown danh sách kết quả.
 */

import React from "react";
import { MessageSquare } from "lucide-react";

interface SearchUser {
  user_id: string;
  email: string;
  full_name: string;
}

interface NewConversationPanelProps {
  /** Giá trị tìm kiếm hiện tại */
  value: string;
  onChange: (val: string) => void;
  /** Kết quả tìm kiếm từ API */
  searchResults: SearchUser[];
  /** Callback khi chọn một user để bắt đầu chat */
  onStartConversation: (userId: string) => void;
  /** Hàm tạo initials từ tên */
  getInitials: (name: string) => string;
}

export default function NewConversationPanel({
  value,
  onChange,
  searchResults,
  onStartConversation,
  getInitials,
}: NewConversationPanelProps) {
  return (
    <div className="px-4 pb-4 animate-in fade-in slide-in-from-top-2 duration-200">
      {/* Input tìm kiếm */}
      <div className="relative">
        <input
          type="text"
          placeholder="Tìm kiếm email hoặc tên..."
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="w-full rounded-2xl border-none bg-zinc-100 dark:bg-zinc-900 px-4 py-2.5 pl-10 text-sm text-zinc-900 dark:text-zinc-100 placeholder-zinc-500 focus:ring-2 focus:ring-blue-500 focus:bg-white dark:focus:bg-zinc-950 transition-all outline-none"
        />
        <MessageSquare className="absolute left-4 top-2.5 h-4 w-4 text-zinc-400" />
      </div>

      {/* Dropdown kết quả */}
      {searchResults.length > 0 && (
        <div className="mt-2 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-xl overflow-hidden animate-in fade-in duration-200">
          {searchResults.map((user) => (
            <button
              key={user.user_id}
              onClick={() => onStartConversation(user.user_id)}
              className="w-full text-left px-4 py-3 hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors flex items-center gap-3 border-b border-zinc-100 dark:border-zinc-800 last:border-0"
            >
              <div className="h-10 w-10 flex-shrink-0 rounded-full bg-gradient-to-br from-blue-500 to-indigo-500 flex items-center justify-center text-white font-bold shadow-sm">
                {getInitials(user.full_name || user.email)}
              </div>
              <div className="flex flex-col overflow-hidden">
                <span className="font-semibold text-zinc-900 dark:text-zinc-100 truncate">
                  {user.full_name}
                </span>
                <span className="text-xs text-zinc-500 truncate">{user.email}</span>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
