'use client';

import React from 'react';
import { usePathname } from 'next/navigation';
import { Layout, ConfigProvider, Spin } from 'antd';
import Sidebar from '@/components/Sidebar';
import { useAuth } from '@/contexts/AuthContext';
import 'antd/dist/reset.css';

const { Content } = Layout;

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { user, loading } = useAuth();
  
  // 不需要认证的页面
  const publicPages = ['/login', '/register'];
  const isPublicPage = publicPages.includes(pathname);

  // 如果正在加载认证状态，显示加载页面
  if (loading) {
    return (
      <ConfigProvider
        theme={{
          token: {
            colorPrimary: '#667eea',
          },
        }}
      >
        <div style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100vh',
          background: '#f5f5f5'
        }}>
          <Spin size="large" />
        </div>
      </ConfigProvider>
    );
  }

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#667eea',
        },
      }}
    >
      {isPublicPage ? (
        // 公共页面（登录/注册）不显示侧边栏
        <div style={{ minHeight: '100vh' }}>
          {children}
        </div>
      ) : (
        // 需要认证的页面显示完整布局
        <Layout style={{ minHeight: '100vh' }}>
          <Sidebar />
          <Layout>
            <Content style={{ padding: 24 }}>
              {children}
            </Content>
          </Layout>
        </Layout>
      )}
    </ConfigProvider>
  );
}