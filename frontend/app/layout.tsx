import type { Metadata } from 'next';
import './globals.css';
import AppLayout from '@/components/AppLayout';
import ClientRedirect from '@/components/ClientRedirect';
import { AuthProvider } from '@/contexts/AuthContext';

// 静态元数据配置 - 完全支持静态导出
export const metadata: Metadata = {
  title: 'GalaxyERP - 企业资源规划系统',
  description: 'Galaxy ERP 前端 - 现代化企业资源规划系统，支持销售、采购、库存、生产、项目、人力资源、财务管理',
  keywords: 'ERP, 企业资源规划, 销售管理, 采购管理, 库存管理, 生产管理',
  authors: [{ name: 'Galaxy ERP Team' }],
  creator: 'Galaxy ERP',
  publisher: 'Galaxy ERP',
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  openGraph: {
    type: 'website',
    locale: 'zh_CN',
    title: 'GalaxyERP - 企业资源规划系统',
    description: 'Galaxy ERP 前端 - 现代化企业资源规划系统',
    siteName: 'GalaxyERP',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'GalaxyERP - 企业资源规划系统',
    description: 'Galaxy ERP 前端 - 现代化企业资源规划系统',
  },
  viewport: {
    width: 'device-width',
    initialScale: 1,
    maximumScale: 1,
  },
  themeColor: '#1890ff',
  manifest: '/manifest.json',
};

// 根布局组件 - 完全支持静态导出
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <head>
        {/* 静态导出优化的 meta 标签 */}
        <meta charSet="utf-8" />
        <meta name="format-detection" content="telephone=no" />
        <meta name="msapplication-tap-highlight" content="no" />
        {/* PWA 支持 */}
        <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
        <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
        <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
      </head>
      <body>
        {/* 客户端组件包装 - 支持静态导出 */}
        <AuthProvider>
          <ClientRedirect>
            <AppLayout>{children}</AppLayout>
          </ClientRedirect>
        </AuthProvider>
      </body>
    </html>
  );
}