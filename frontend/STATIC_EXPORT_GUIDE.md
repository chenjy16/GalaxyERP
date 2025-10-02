# Galaxy ERP 静态导出部署指南

## 📋 概述

本指南详细说明如何将 Galaxy ERP 前端项目配置为支持静态导出部署，以及相关的优化策略和部署方法。

## 🔧 静态导出配置

### 1. Layout.tsx 优化

已对 `app/layout.tsx` 进行以下优化：

#### ✅ 静态元数据配置
```typescript
export const metadata: Metadata = {
  title: 'GalaxyERP - 企业资源规划系统',
  description: 'Galaxy ERP 前端 - 现代化企业资源规划系统',
  // ... 完整的 SEO 和 PWA 元数据
};
```

#### ✅ 客户端重定向处理
- 使用 `ClientRedirect` 组件替代服务器端重定向
- 支持认证状态检查和路由保护
- 完全兼容静态导出

#### ✅ PWA 支持
- 添加了 manifest.json 引用
- 配置了图标和主题色
- 支持离线访问

### 2. Next.js 配置

#### 标准配置 (next.config.js)
```javascript
// 动态配置，支持多种部署方式
// output: 'standalone', // Docker 部署
// output: 'export',     // 静态导出
```

#### 静态导出专用配置 (next.config.static.js)
```javascript
const nextConfig = {
  output: 'export',              // 启用静态导出
  trailingSlash: true,           // 添加尾部斜杠
  images: { unoptimized: true }, // 禁用图片优化
  // ... 其他静态导出优化配置
};
```

### 3. 客户端路由处理

#### ClientRedirect 组件特性
- **认证检查**: 自动检查用户登录状态
- **路由保护**: 未登录用户重定向到登录页
- **返回 URL**: 登录后返回原访问页面
- **加载状态**: 认证状态加载期间显示加载界面

#### 支持的路由逻辑
```typescript
// 公开页面（无需认证）
const publicPaths = ['/login', '/register'];

// 受保护页面（需要认证）
// 所有其他页面都需要登录
```

## 🚀 构建和部署

### 1. 构建命令

```bash
# 使用静态导出配置构建
npm run build:static

# 或使用专用导出命令
npm run export:static

# 清理并重新构建
npm run clean:static
```

### 2. 预览静态文件

```bash
# 本地预览静态导出结果
npm run preview:static

# 访问 http://localhost:3000 查看效果
```

### 3. 部署目标

#### CDN 部署
```bash
# 构建静态文件
npm run build:static

# 上传 out/ 目录到 CDN
aws s3 sync out/ s3://your-bucket --delete
```

#### Nginx 部署
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/out;
    index index.html;
    
    # SPA 路由支持
    location / {
        try_files $uri $uri/ $uri.html /index.html;
    }
    
    # 静态资源缓存
    location /_next/static/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

## 📊 静态导出优势

### 1. 性能优势
- **零服务器延迟**: 纯静态文件，无服务器处理时间
- **CDN 加速**: 全球分发，就近访问
- **缓存友好**: 静态资源可长期缓存

### 2. 部署优势
- **简单部署**: 只需上传静态文件
- **低成本**: 无需服务器运行成本
- **高可用**: 静态文件天然高可用

### 3. SEO 优势
- **预渲染**: 所有页面在构建时预渲染
- **完整元数据**: 静态 HTML 包含完整 SEO 信息
- **快速加载**: 提升搜索引擎排名

## ⚠️ 静态导出限制

### 1. 功能限制
- **无 API Routes**: 不支持 Next.js API 路由
- **无服务器函数**: 不支持 `getServerSideProps`
- **无动态重定向**: 重定向需要在客户端处理

### 2. 解决方案
- **外部 API**: 使用独立的后端 API 服务
- **客户端渲染**: 动态内容在客户端获取
- **客户端路由**: 使用 React Router 处理路由

## 🔄 环境切换

### 开发环境
```bash
# 标准开发模式
npm run dev
```

### 静态导出测试
```bash
# 构建静态版本
npm run build:static

# 本地预览
npm run preview:static
```

### 生产部署
```bash
# 根据部署方式选择
./deploy.sh static    # 静态导出
./deploy.sh docker    # Docker 部署
./deploy.sh standalone # 独立部署
```

## 🛠 故障排除

### 1. 构建错误

#### 图片优化错误
```bash
# 错误: Image Optimization using Next.js' default loader is not compatible with `output: 'export'`
# 解决: 在 next.config.static.js 中设置 images.unoptimized: true
```

#### 动态路由错误
```bash
# 错误: Dynamic routes are not supported with `output: 'export'`
# 解决: 使用 exportPathMap 预定义所有路由
```

### 2. 运行时错误

#### 路由问题
```bash
# 问题: 刷新页面 404
# 解决: 配置服务器支持 SPA 路由回退
```

#### API 调用失败
```bash
# 问题: API 调用相对路径错误
# 解决: 使用完整的 API URL (NEXT_PUBLIC_API_URL)
```

## 📈 性能优化

### 1. 构建优化
```javascript
// next.config.static.js
experimental: {
  optimizePackageImports: ['antd', '@ant-design/icons'],
},
```

### 2. 资源优化
```javascript
// 启用压缩
compress: true,
// 禁用不必要的头部
poweredByHeader: false,
generateEtags: false,
```

### 3. 缓存策略
```nginx
# Nginx 缓存配置
location /_next/static/ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}

location /static/ {
    expires 30d;
    add_header Cache-Control "public";
}
```

## 🔍 监控和分析

### 1. 构建分析
```bash
# 分析构建包大小
npm run build:analyze
```

### 2. 性能监控
- 使用 Lighthouse 检查性能
- 监控 Core Web Vitals
- 分析加载时间和交互性

### 3. 用户体验
- 配置 PWA 离线支持
- 实现加载状态指示
- 优化首屏渲染时间

## 📝 最佳实践

1. **定期测试**: 每次发布前测试静态导出
2. **性能监控**: 持续监控页面性能指标
3. **缓存策略**: 合理配置静态资源缓存
4. **错误处理**: 完善的客户端错误处理
5. **用户体验**: 优化加载状态和交互反馈

---

通过以上配置和优化，Galaxy ERP 前端项目已完全支持静态导出部署，可以享受静态网站的所有优势，同时保持现代 SPA 应用的用户体验。