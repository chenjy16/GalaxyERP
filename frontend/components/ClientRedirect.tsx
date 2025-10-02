'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

interface ClientRedirectProps {
  children: React.ReactNode;
}

/**
 * 客户端重定向组件
 * 用于静态导出部署中处理认证重定向逻辑
 * 替代 next.config.js 中的服务器端重定向
 */
export default function ClientRedirect({ children }: ClientRedirectProps) {
  const router = useRouter();
  const { user, loading } = useAuth();

  useEffect(() => {
    // 等待认证状态加载完成
    if (loading) return;

    // 获取当前路径
    const currentPath = window.location.pathname;
    
    // 公开页面列表（不需要认证的页面）
    const publicPaths = ['/login', '/register'];
    
    // 如果是公开页面，不需要重定向
    if (publicPaths.includes(currentPath)) {
      // 如果已登录用户访问登录页，重定向到首页
      if (user && currentPath === '/login') {
        router.replace('/');
      }
      return;
    }

    // 如果未登录且不在公开页面，重定向到登录页
    if (!user && !publicPaths.includes(currentPath)) {
      // 保存当前路径，登录后可以返回
      const returnUrl = encodeURIComponent(currentPath);
      router.replace(`/login?returnUrl=${returnUrl}`);
      return;
    }

    // 如果是根路径且已登录，可以重定向到默认页面
    if (currentPath === '/' && user) {
      // 可以根据用户角色重定向到不同的默认页面
      // 这里暂时不做重定向，显示首页
    }
  }, [user, loading, router]);

  // 在认证状态加载期间显示加载状态
  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        fontSize: '16px',
        color: '#666'
      }}>
        加载中...
      </div>
    );
  }

  return <>{children}</>;
}