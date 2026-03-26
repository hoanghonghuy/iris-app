import React from "react";
import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface UsersToolbarProps {
  searchQuery: string;
  roleFilter: string;
  onSearchChange: (value: string) => void;
  onRoleFilterChange: (value: string) => void;
}

export function UsersToolbar({
  searchQuery,
  roleFilter,
  onSearchChange,
  onRoleFilterChange,
}: UsersToolbarProps) {
  return (
    <div className="flex items-center gap-3">
      <div className="relative flex-1 max-w-sm">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          type="search"
          placeholder="Tìm theo email..."
          className="pl-8 bg-background min-w-0"
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
        />
      </div>

      <Select value={roleFilter} onValueChange={onRoleFilterChange}>
        <SelectTrigger className="w-[140px] shrink-0">
          <SelectValue placeholder="Tất cả vai trò" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="ALL">Tất cả vai trò</SelectItem>
          <SelectItem value="TEACHER">Giáo viên</SelectItem>
          <SelectItem value="PARENT">Phụ huynh</SelectItem>
          <SelectItem value="SCHOOL_ADMIN">School Admin</SelectItem>
          <SelectItem value="SUPER_ADMIN">Super Admin</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
}
