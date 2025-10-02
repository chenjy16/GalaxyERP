/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  
  // 静态导出配置
  output: 'export',
  trailingSlash: true,
  skipTrailingSlashRedirect: true,
  
  // 生产环境优化
  compress: true,
  poweredByHeader: false,
  generateEtags: false,
  
  // 图片优化配置 - 静态导出需要禁用优化
  images: {
    unoptimized: true, // 静态导出必须禁用图片优化
    formats: ['image/webp', 'image/avif'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
  },
  
  // 实验性功能优化
  experimental: {
    optimizePackageImports: ['antd', '@ant-design/icons'],
  },
  
  // 环境变量配置
  env: {
    CUSTOM_KEY: process.env.CUSTOM_KEY,
  },
  
  // 静态导出不支持重定向，需要在客户端处理
  // async redirects() {
  //   return [];
  // },
  
  // 静态导出不支持动态头部，需要在 HTML 中静态定义
  // async headers() {
  //   return [];
  // },
  
  // 静态导出配置
  exportPathMap: async function (
    defaultPathMap,
    { dev, dir, outDir, distDir, buildId }
  ) {
    return {
      '/': { page: '/' },
      '/login': { page: '/login' },
      '/register': { page: '/register' },
      '/accounting': { page: '/accounting' },
      '/hr': { page: '/hr' },
      '/inventory': { page: '/inventory' },
      '/production': { page: '/production' },
      '/project': { page: '/project' },
      '/purchase': { page: '/purchase' },
      '/sales': { page: '/sales' },
      '/system': { page: '/system' },
    };
  },
  
  // 基础路径配置（如果部署到子目录）
  // basePath: '/galaxy-erp',
  
  // 资源前缀配置（如果使用 CDN）
  // assetPrefix: 'https://cdn.example.com',
};

module.exports = nextConfig;