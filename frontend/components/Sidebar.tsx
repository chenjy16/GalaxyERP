'use client';

import { useState } from 'react';
import { Layout, Menu, Typography, Avatar, Divider, Button } from 'antd';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import {
  DashboardOutlined,
  ShoppingCartOutlined,
  ShoppingOutlined,
  InboxOutlined,
  DollarOutlined,
  TeamOutlined,
  ToolOutlined,
  FolderOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined
} from '@ant-design/icons';

const { Sider } = Layout;
const { Text, Title } = Typography;

export default function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const pathname = usePathname();
  const { user, logout } = useAuth();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link href="/">仪表盘</Link>,
    },
    {
      key: '/sales',
      icon: <ShoppingCartOutlined />,
      label: <Link href="/sales">销售管理</Link>,
    },
    {
      key: '/purchase',
      icon: <ShoppingOutlined />,
      label: <Link href="/purchase">采购管理</Link>,
    },
    {
      key: '/inventory',
      icon: <InboxOutlined />,
      label: <Link href="/inventory">库存管理</Link>,
    },
    {
      key: '/accounting',
      icon: <DollarOutlined />,
      label: <Link href="/accounting">财务管理</Link>,
    },
    {
      key: '/hr',
      icon: <TeamOutlined />,
      label: <Link href="/hr">人力资源</Link>,
    },
    {
      key: '/production',
      icon: <ToolOutlined />,
      label: <Link href="/production">生产管理</Link>,
    },
    {
      key: '/project',
      icon: <FolderOutlined />,
      label: <Link href="/project">项目管理</Link>,
    },
    {
      key: '/system',
      icon: <SettingOutlined />,
      label: <Link href="/system">系统管理</Link>,
    },
  ];

  return (
    <Sider 
      trigger={null} 
      collapsible 
      collapsed={collapsed}
      width={280}
      collapsedWidth={80}
      style={{
        background: 'linear-gradient(180deg, #1f2937 0%, #111827 100%)',
        boxShadow: '2px 0 8px rgba(0,0,0,0.15)',
        position: 'relative',
        zIndex: 100
      }}
    >
      {/* 头部区域 */}
      <div style={{ 
        padding: collapsed ? '16px 8px' : '20px 24px', 
        borderBottom: '1px solid rgba(255,255,255,0.1)',
        textAlign: 'center'
      }}>
        {!collapsed ? (
          <div>
            <div style={{ 
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              width: 48,
              height: 48,
              borderRadius: 12,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 12px',
              fontSize: 24
            }}>
              🏢
            </div>
            <Title level={4} style={{ color: '#fff', margin: 0, fontSize: 18 }}>
              Galaxy ERP
            </Title>
            <Text style={{ color: 'rgba(255,255,255,0.6)', fontSize: 12 }}>
              企业资源规划系统
            </Text>
          </div>
        ) : (
          <div style={{ 
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            width: 40,
            height: 40,
            borderRadius: 10,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            margin: '0 auto',
            fontSize: 20
          }}>
            🏢
          </div>
        )}
      </div>

      {/* 用户信息区域 */}
      {!collapsed && user && (
        <div style={{ 
          padding: '16px 24px',
          borderBottom: '1px solid rgba(255,255,255,0.1)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <Avatar 
              size={40}
              style={{ 
                background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
                border: '2px solid rgba(255,255,255,0.2)'
              }}
              icon={<UserOutlined />}
            />
            <div style={{ flex: 1 }}>
              <Text strong style={{ color: '#fff', display: 'block', fontSize: 14 }}>
                {user.firstName} {user.lastName}
              </Text>
              <Text style={{ color: 'rgba(255,255,255,0.6)', fontSize: 12 }}>
                {user.email}
              </Text>
            </div>
          </div>
        </div>
      )}

      {/* 导航菜单 */}
      <div style={{ flex: 1, padding: '16px 0' }}>
        <Menu
          mode="inline"
          selectedKeys={[pathname]}
          style={{
            background: 'transparent',
            border: 'none',
            color: '#fff'
          }}
          items={menuItems}
          theme="dark"
        />
      </div>

      {/* 底部操作区域 */}
      <div style={{ 
        padding: collapsed ? '16px 8px' : '16px 24px',
        borderTop: '1px solid rgba(255,255,255,0.1)'
      }}>
        {!collapsed && (
          <div style={{ marginBottom: 12 }}>
            <Button 
              type="text" 
              icon={<LogoutOutlined />}
              onClick={logout}
              style={{ 
                color: 'rgba(255,255,255,0.8)',
                width: '100%',
                textAlign: 'left',
                height: 40,
                display: 'flex',
                alignItems: 'center',
                gap: 8
              }}
            >
              退出登录
            </Button>
          </div>
        )}
        
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={() => setCollapsed(!collapsed)}
          style={{
            color: 'rgba(255,255,255,0.8)',
            width: '100%',
            height: 40,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}
        />
      </div>

      {/* 自定义样式 */}
      <style jsx global>{`
        .ant-menu-dark .ant-menu-item {
          margin: 4px 16px !important;
          border-radius: 8px !important;
          height: 44px !important;
          line-height: 44px !important;
          transition: all 0.3s ease !important;
        }
        
        .ant-menu-dark .ant-menu-item:hover {
          background: rgba(255,255,255,0.1) !important;
          transform: translateX(4px) !important;
        }
        
        .ant-menu-dark .ant-menu-item-selected {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4) !important;
        }
        
        .ant-menu-dark .ant-menu-item-selected::after {
          display: none !important;
        }
        
        .ant-menu-dark .ant-menu-item a {
          color: inherit !important;
          text-decoration: none !important;
        }
        
        .ant-layout-sider-collapsed .ant-menu-dark .ant-menu-item {
          margin: 4px 8px !important;
          padding: 0 !important;
          text-align: center !important;
        }
        
        .ant-layout-sider-collapsed .ant-menu-dark .ant-menu-item:hover {
          transform: scale(1.05) !important;
        }
      `}</style>
    </Sider>
  );
}