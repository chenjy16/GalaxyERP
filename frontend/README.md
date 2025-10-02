# Galaxy ERP 前端

![Next.js](https://img.shields.io/badge/Next.js-14.1.0-black)
![React](https://img.shields.io/badge/React-18.2.0-blue)
![TypeScript](https://img.shields.io/badge/TypeScript-5.4.0-blue)
![Ant Design](https://img.shields.io/badge/Ant%20Design-5.13.0-blue)

Galaxy ERP 是一个现代化的企业资源规划系统前端应用，基于 Next.js 14 和 React 18 构建，采用 TypeScript 开发，使用 Ant Design 作为 UI 组件库。

## 📋 目录

- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [环境配置](#环境配置)
- [开发指南](#开发指南)
- [部署方式](#部署方式)
- [脚本说明](#脚本说明)
- [环境变量](#环境变量)
- [API 集成](#api-集成)
- [贡献指南](#贡献指南)

## ✨ 功能特性

- 🏢 **多模块管理**: 销售、采购、库存、生产、项目、人力资源、财务、系统管理
- 🔐 **用户认证**: 完整的登录/注册系统，基于 JWT 的身份验证
- 📱 **响应式设计**: 支持桌面端和移动端访问
- 🎨 **现代化 UI**: 基于 Ant Design 的美观界面
- 🚀 **高性能**: Next.js 14 App Router，支持 SSR/SSG/CSR
- 🔧 **TypeScript**: 完整的类型安全支持
- 🐳 **容器化**: 支持 Docker 部署
- 📦 **多种部署方式**: 静态导出、Docker、Standalone 模式

## 🛠 技术栈

### 核心框架
- **Next.js 14.1.0** - React 全栈框架，使用 App Router
- **React 18.2.0** - 用户界面库
- **TypeScript 5.4.0** - 类型安全的 JavaScript

### UI 组件
- **Ant Design 5.13.0** - 企业级 UI 组件库
- **@ant-design/icons 5.2.6** - 图标库
- **Day.js 1.11.18** - 日期处理库

### 开发工具
- **ESLint** - 代码质量检查
- **TypeScript Compiler** - 类型检查

## 📁 项目结构

```
frontend/
├── app/                    # Next.js 14 App Router 页面
│   ├── accounting/         # 财务管理
│   ├── hr/                # 人力资源
│   ├── inventory/         # 库存管理
│   ├── login/             # 登录页面
│   ├── production/        # 生产管理
│   ├── project/           # 项目管理
│   ├── purchase/          # 采购管理
│   ├── register/          # 注册页面
│   ├── sales/             # 销售管理
│   ├── system/            # 系统管理
│   ├── layout.tsx         # 根布局
│   ├── page.tsx           # 首页
│   └── globals.css        # 全局样式
├── components/            # 可复用组件
│   ├── AppLayout.tsx      # 应用布局组件
│   └── Sidebar.tsx        # 侧边栏组件
├── contexts/              # React Context
│   └── AuthContext.tsx    # 认证上下文
├── lib/                   # 工具库
│   └── api.ts             # API 客户端
├── services/              # 业务服务层
│   ├── auth.ts            # 认证服务
│   ├── customer.ts        # 客户服务
│   ├── employee.ts        # 员工服务
│   ├── inventory.ts       # 库存服务
│   └── ...               # 其他业务服务
├── types/                 # TypeScript 类型定义
│   └── api.ts             # API 类型定义
├── deploy.sh              # 部署脚本
├── docker-compose.prod.yml # 生产环境 Docker Compose
├── Dockerfile             # Docker 构建文件
├── nginx.conf             # Nginx 配置
├── next.config.js         # Next.js 配置
└── package.json           # 项目依赖和脚本
```

## 🚀 快速开始

### 环境要求

- Node.js >= 18.0.0
- npm >= 8.0.0 或 yarn >= 1.22.0
- Git

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd galaxyErp/frontend
```

2. **安装依赖**
```bash
npm install
# 或
yarn install
```

3. **配置环境变量**
```bash
# 复制环境变量模板（如果存在）
cp .env.example .env.local

# 编辑环境变量
vim .env.local
```

4. **启动开发服务器**
```bash
npm run dev
# 或
yarn dev
```

5. **访问应用**
打开浏览器访问 [http://localhost:3000](http://localhost:3000)

## ⚙️ 环境配置

### 环境变量

创建 `.env.local` 文件并配置以下变量：

```bash
# API 服务地址
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# 自定义配置
CUSTOM_KEY=your_custom_key

# 运行环境
NODE_ENV=development
```

### 不同环境配置

#### 开发环境 (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NODE_ENV=development
```

#### 测试环境 (.env.test)
```bash
NEXT_PUBLIC_API_URL=https://test-api.galaxy-erp.com/api/v1
NODE_ENV=test
```

#### 生产环境 (.env.production)
```bash
NEXT_PUBLIC_API_URL=https://api.galaxy-erp.com/api/v1
NODE_ENV=production
```

## 👨‍💻 开发指南

### 可用脚本

```bash
# 开发服务器
npm run dev

# 构建应用
npm run build

# 构建分析
npm run build:analyze

# 独立模式构建
npm run build:standalone

# 静态导出构建
npm run build:export

# 启动生产服务器
npm run start

# 生产模式启动
npm run start:prod

# 代码检查
npm run lint

# 自动修复代码问题
npm run lint:fix

# 类型检查
npm run type-check

# 清理构建文件
npm run clean
```

### 开发规范

1. **代码风格**: 使用 ESLint 和 TypeScript 严格模式
2. **组件开发**: 优先使用函数组件和 React Hooks
3. **状态管理**: 使用 React Context 进行全局状态管理
4. **API 调用**: 统一使用 `lib/api.ts` 中的 ApiClient
5. **类型定义**: 所有 API 接口都需要在 `types/api.ts` 中定义类型

### 目录规范

- `app/`: 页面组件，遵循 Next.js 14 App Router 规范
- `components/`: 可复用的 UI 组件
- `services/`: 业务逻辑和 API 调用
- `types/`: TypeScript 类型定义
- `lib/`: 工具函数和配置

## 🚀 部署方式

项目支持三种部署方式，通过 `deploy.sh` 脚本实现：

### 1. 静态导出部署 (推荐用于 CDN)

```bash
./deploy.sh static
```

**特点:**
- 生成纯静态文件到 `out/` 目录
- 适合部署到 CDN、对象存储或静态托管服务
- 最佳性能和 SEO 优化
- 支持 Nginx、Apache 等传统 Web 服务器

**部署目标:**
- AWS S3 + CloudFront
- 阿里云 OSS + CDN
- Vercel、Netlify 等静态托管
- 传统服务器 + Nginx

#### 静态导出详细步骤

```bash
# 构建静态文件
npm run build:static

# 预览静态文件
npm run preview:static

# 清理并重新构建
npm run clean:static

# 使用部署脚本
./deploy.sh static
```

#### 静态导出特性

- ✅ **完全静态化**: 所有页面预渲染为静态 HTML
- ✅ **PWA 支持**: 包含 manifest.json 和离线功能
- ✅ **SEO 优化**: 完整的元数据和结构化数据
- ✅ **客户端路由**: 支持 SPA 路由和认证重定向
- ✅ **性能优化**: 资源压缩和缓存策略

#### 部署选项

1. **CDN 部署**
   ```bash
   # AWS S3 + CloudFront
   aws s3 sync out/ s3://your-bucket --delete
   
   # Cloudflare Pages
   # 直接连接 Git 仓库自动部署
   ```

2. **Nginx 部署**
   ```bash
   # 复制配置文件
   sudo cp nginx.static.conf /etc/nginx/sites-available/galaxy-erp
   sudo ln -s /etc/nginx/sites-available/galaxy-erp /etc/nginx/sites-enabled/
   
   # 重启 Nginx
   sudo nginx -t && sudo systemctl reload nginx
   ```

3. **静态托管平台**
   - Vercel: 自动检测 Next.js 项目
   - Netlify: 构建命令 `npm run build:static`
   - GitHub Pages: 部署 `out/` 目录

### 2. Docker 容器化部署 (推荐用于生产环境)

```bash
./deploy.sh docker
```

**特点:**
- 使用 Docker 多阶段构建
- 包含 Nginx 反向代理配置
- 支持水平扩展和负载均衡
- 完整的生产环境优化

**部署流程:**
```bash
# 构建镜像
docker build -t galaxy-erp-frontend .

# 使用 Docker Compose 部署
docker-compose -f docker-compose.prod.yml up -d
```

### 3. Standalone 独立部署

```bash
./deploy.sh standalone
```

**特点:**
- 生成独立的 Node.js 应用
- 包含所有依赖，无需 node_modules
- 适合传统服务器部署
- 支持 PM2 等进程管理器

## 📜 脚本说明

### deploy.sh 部署脚本

```bash
# 查看帮助
./deploy.sh --help

# 静态导出
./deploy.sh static

# Docker 部署
./deploy.sh docker

# Standalone 部署
./deploy.sh standalone
```

**脚本功能:**
- 自动依赖检查 (Node.js, npm, Docker)
- 清理旧构建文件
- 安装/更新依赖
- TypeScript 类型检查
- ESLint 代码检查
- 根据部署类型修改 `next.config.js`
- 执行构建命令
- 恢复配置文件

### Docker 相关

#### Dockerfile 特性
- 多阶段构建优化镜像大小
- 非 root 用户运行提升安全性
- 健康检查支持
- 生产环境优化

#### docker-compose.prod.yml
- 前端服务 + Nginx 反向代理
- 后端 API 服务集成
- 网络隔离和服务发现
- 健康检查和重启策略

## 🌐 API 集成

### API 客户端配置

项目使用统一的 API 客户端 (`lib/api.ts`)：

```typescript
import { apiClient } from '@/lib/api';

// GET 请求
const data = await apiClient.get('/users');

// POST 请求
const result = await apiClient.post('/users', userData);

// PUT 请求
const updated = await apiClient.put('/users/1', updateData);

// DELETE 请求
await apiClient.delete('/users/1');
```

### 认证集成

- 自动处理 JWT Token
- Token 存储在 localStorage
- 自动添加 Authorization 头
- 支持 Token 刷新机制

## 🔧 故障排除

### 常见问题

1. **端口冲突**
```bash
# 检查端口占用
lsof -i :3000

# 使用其他端口
PORT=3001 npm run dev
```

2. **依赖安装失败**
```bash
# 清理缓存
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

3. **构建失败**
```bash
# 类型检查
npm run type-check

# 代码检查
npm run lint
```

4. **Docker 构建问题**
```bash
# 清理 Docker 缓存
docker system prune -a

# 重新构建
docker build --no-cache -t galaxy-erp-frontend .
```

### 性能优化

1. **构建优化**
- 启用 `compress` 压缩
- 禁用 `poweredByHeader`
- 优化包导入 `optimizePackageImports`

2. **运行时优化**
- 使用 React.memo 优化组件渲染
- 实现虚拟滚动处理大列表
- 使用 useMemo 和 useCallback 优化计算

3. **部署优化**
- 启用 Gzip 压缩
- 配置静态资源缓存
- 使用 CDN 加速

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- 遵循 ESLint 配置
- 使用 TypeScript 严格模式
- 编写单元测试
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 支持

如有问题或建议，请：

1. 查看 [Issues](../../issues) 页面
2. 创建新的 Issue
3. 联系开发团队

---

**Galaxy ERP Frontend** - 现代化企业资源规划系统前端