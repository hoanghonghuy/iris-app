/**
 * AuthProvider Component
 * Quản lý trạng thái đăng nhập toàn cục của ứng dụng.
 * Cung cấp thông tin user, role và các hàm login/logout cho toàn bộ app.
 */
"use client";

import React, { createContext, useContext, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { authHelpers } from '@/lib/api/client';
import { authApi } from '@/lib/api/auth.api';
import { UserInfo, UserRole } from '@/types';

interface AuthContextType {
  user: UserInfo | null;
  role: UserRole | null;
  isLoading: boolean;
  login: (token: string, role: UserRole) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [role, setRole] = useState<UserRole | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  // Kiểm tra trạng thái đăng nhập khi app khởi chạy
  useEffect(() => {
    const initAuth = async () => {
      const token = authHelpers.getToken();
      const savedRole = authHelpers.getUserRole() as UserRole;

      if (token && savedRole) {
        try {
          // Lấy thông tin user mới nhất từ server
          const userData = await authApi.getMe();
          setUser(userData);
          setRole(savedRole);
        } catch (error) {
          console.error("Auth initialization failed:", error);
          logout();
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = (token: string, userRole: UserRole) => {
    authHelpers.setToken(token);
    authHelpers.setUserRole(userRole);
    setRole(userRole);
    
    // Sau khi login, fetch thông tin user
    authApi.getMe().then(userData => {
      setUser(userData);
      
      // Redirect dựa trên role
      switch (userRole) {
        case 'SUPER_ADMIN':
        case 'SCHOOL_ADMIN':
          router.push('/admin');
          break;
        case 'TEACHER':
          router.push('/teacher');
          break;
        case 'PARENT':
          router.push('/parent');
          break;
        default:
          router.push('/');
      }
    });
  };

  const logout = () => {
    authHelpers.removeToken();
    setUser(null);
    setRole(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, role, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook để sử dụng AuthContext dễ dàng hơn
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};