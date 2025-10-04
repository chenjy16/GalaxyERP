'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { message } from 'antd';
import AuthService from '@/services/auth';
import { User, LoginRequest, RegisterRequest } from '@/types/api';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  // 初始化时检查本地存储的用户信息
  useEffect(() => {
    const initAuth = async () => {
      try {
        let token = localStorage.getItem('auth_token');
        
        // 开发环境自动登录
        if (!token && process.env.NODE_ENV === 'development') {
          try {
            const response = await AuthService.login({
              username: 'admin',
              password: 'password'
            });
            setUser(response.user);
            localStorage.setItem('user', JSON.stringify(response.user));
            return;
          } catch (error) {
            // 自动登录失败，继续正常流程
          }
        }
        
        if (token) {
          // 先设置 token 到 API 客户端
          const { apiClient } = await import('@/lib/api');
          apiClient.setToken(token);
          
          // 验证 token 是否有效
          const currentUser = await AuthService.getCurrentUser();
          setUser(currentUser);
        }
      } catch (error) {
        // Token 无效，清除本地存储
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user');
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, []);

  const login = async (credentials: LoginRequest) => {
    try {
      const response = await AuthService.login(credentials);
      setUser(response.user);
      localStorage.setItem('user', JSON.stringify(response.user));
      message.success('登录成功！');
      router.push('/');
    } catch (error: any) {
      message.error(error.message || '登录失败');
      throw error;
    }
  };

  const register = async (userData: RegisterRequest) => {
    try {
      await AuthService.register(userData);
      message.success('注册成功！请登录');
      router.push('/login');
    } catch (error: any) {
      message.error(error.message || '注册失败');
      throw error;
    }
  };

  const logout = () => {
    AuthService.logout();
    setUser(null);
    localStorage.removeItem('user');
    message.success('已退出登录');
    router.push('/login');
  };

  const refreshUser = async () => {
    try {
      const currentUser = await AuthService.getCurrentUser();
      setUser(currentUser);
      localStorage.setItem('user', JSON.stringify(currentUser));
    } catch (error) {
      // 如果获取用户信息失败，可能是 token 过期
      logout();
    }
  };

  const value: AuthContextType = {
    user,
    loading,
    login,
    register,
    logout,
    refreshUser,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

// 高阶组件：保护需要认证的页面
export function withAuth<P extends object>(Component: React.ComponentType<P>) {
  return function AuthenticatedComponent(props: P) {
    const { user, loading } = useAuth();
    const router = useRouter();

    useEffect(() => {
      if (!loading && !user) {
        router.push('/login');
      }
    }, [user, loading, router]);

    if (loading) {
      return <div>Loading...</div>;
    }

    if (!user) {
      return null;
    }

    // Fix: Ensure we're not passing the raw user object as props
    return <Component {...props} />;
  };
}