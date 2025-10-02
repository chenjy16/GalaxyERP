/** @type {import('next').NextConfig} */
const isStatic = process.env.NEXT_CONFIG === 'static';

const nextConfig = {
  reactStrictMode: true,
  
  // 生产环境优化
  compress: true,
  poweredByHeader: false,
  generateEtags: false,
  
  // 图片优化配置
  images: {
    formats: ['image/webp', 'image/avif'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    unoptimized: isStatic, // 静态导出时禁用图片优化
  },
  
  // 实验性功能优化
  experimental: {
    optimizePackageImports: ['antd', '@ant-design/icons'],
  },
  
  // 环境变量配置
  env: {
    CUSTOM_KEY: process.env.CUSTOM_KEY,
  },
  
  // 部署配置
  ...(isStatic && {
    output: 'export',
    trailingSlash: true,
    distDir: 'out',
  }),
  
  // 输出配置 - 根据部署方式选择
  // output: 'standalone', // 用于 Docker 部署
  // output: 'export',     // 用于静态部署
  
  // 重定向配置
  async redirects() {
    return [
      {
        source: '/',
        destination: '/login',
        permanent: false,
        has: [
          {
            type: 'cookie',
            key: 'token',
            value: undefined,
          },
        ],
      },
    ];
  },
  
  // 头部配置
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
        ],
      },
    ];
  },
};

module.exports = nextConfig;