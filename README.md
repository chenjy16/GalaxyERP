# GalaxyERP

一个基于 Go 语言和 Next.js 构建的现代化企业资源规划 (ERP) 系统，提供全面的业务管理解决方案。

## 🚀 项目概述

GalaxyERP 是一个全功能的企业资源规划系统，采用前后端分离架构，为企业提供完整的业务管理平台。系统涵盖财务、销售、采购、库存、生产、项目、人力资源等核心业务模块，支持多环境部署和灵活配置。

### 🏗️ 架构特点

- **微服务架构**: 采用分层架构设计，清晰的职责分离
- **RESTful API**: 标准化的API设计，支持前后端分离
- **依赖注入**: 使用容器模式管理依赖关系，提高代码可测试性
- **中间件模式**: 统一的认证、日志、错误处理和跨域支持
- **仓储模式**: 数据访问层抽象，支持多种数据库
- **DTO模式**: 数据传输对象，确保API数据结构的一致性
- **响应式设计**: 基于Ant Design的现代化UI界面


## 🛠️ 技术栈

### 后端技术栈
- **语言**: Go 1.24.0
- **Web框架**: Gin v1.9.1 (高性能HTTP Web框架)
- **ORM**: GORM v1.31.0 (Go语言ORM库)
- **数据库**: PostgreSQL / SQLite (开发环境)
- **认证**: JWT v5.3.0 (JSON Web Token)
- **配置管理**: Viper v1.17.0
- **日志**: Zap v1.26.0 (高性能日志库)
- **密码加密**: bcrypt (golang.org/x/crypto v0.42.0)
- **数据验证**: go-playground/validator v10.17.0
- **中间件**: CORS, 认证中间件, 日志中间件
- **数据库驱动**: 
  - PostgreSQL: `gorm.io/driver/postgres v1.6.0`
  - SQLite: `glebarez/sqlite v1.11.0`

### 前端技术栈
- **框架**: Next.js 14.1.0 (React全栈框架)
- **语言**: TypeScript 5.4.0
- **UI组件库**: Ant Design 5.13.0
- **图标库**: Ant Design Icons 5.2.6
- **状态管理**: React Context + Hooks
- **HTTP客户端**: Fetch API
- **日期处理**: Day.js 1.11.18
- **构建工具**: Webpack (Next.js内置)
- **包管理**: npm
- **开发工具**: ESLint 8.57.0, TypeScript

### 开发工具
- **版本控制**: Git
- **构建工具**: Make
- **代码格式化**: gofmt, Prettier
- **数据库迁移**: GORM AutoMigrate
- **环境管理**: 多环境配置 (dev/test/prod)
- **依赖管理**: Go Modules, npm

## 🎯 设计模式与架构模式

### 后端设计模式

#### 1. 依赖注入模式 (Dependency Injection)
```go
// Container 统一管理所有依赖
type Container struct {
    DB *gorm.DB
    UserRepository    repositories.UserRepository
    UserService       services.UserService
    UserController    *controllers.UserController
}
```

#### 2. 仓储模式 (Repository Pattern)
```go
// 接口定义
type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    Search(req *dto.PaginationRequest) (*dto.PaginatedResponse[models.User], error)
}

// 实现
type UserRepositoryImpl struct {
    db *gorm.DB
}
```

#### 3. 服务层模式 (Service Layer Pattern)
```go
// 业务逻辑封装
type UserService interface {
    CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
    GetUser(id uint) (*dto.UserResponse, error)
}
```

#### 4. DTO模式 (Data Transfer Object)
```go
// 统一的数据传输结构
type PaginationRequest struct {
    Page     int    `json:"page" validate:"min=1"`
    PageSize int    `json:"page_size" validate:"min=1,max=100"`
    Search   string `json:"search"`
}
```

#### 5. 中间件模式 (Middleware Pattern)
```go
// 认证中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT验证逻辑
        c.Next()
    }
}
```

#### 6. 工厂模式 (Factory Pattern)
```go
// 统一的响应构建
func (h *APIResponseHelper) Success(data interface{}) *StandardAPIResponse {
    return &StandardAPIResponse{
        Success:   true,
        Data:      data,
        Timestamp: time.Now(),
    }
}
```

### 前端设计模式

#### 1. 组件化模式 (Component Pattern)
```typescript
// 可复用的UI组件
interface TableProps<T> {
    data: T[];
    columns: ColumnType<T>[];
    loading?: boolean;
}
```

#### 2. 服务层模式 (Service Layer Pattern)
```typescript
// API服务封装
class SupplierService {
    async getSuppliers(params: PaginationParams): Promise<PaginatedResponse<Supplier>> {
        return this.apiClient.getPaginated<Supplier>('/suppliers', params);
    }
}
```

#### 3. 状态管理模式 (State Management Pattern)
```typescript
// Context + Reducer模式
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) throw new Error('useAuth must be used within AuthProvider');
    return context;
};
```

#### 4. 适配器模式 (Adapter Pattern)
```typescript
// API响应格式转换
const convertBackendPagination = (response: BackendPaginatedResponse<T>): PaginatedResponse<T> => {
    return {
        data: response.data,
        total: response.pagination.total,
        page: response.pagination.page,
        limit: response.pagination.limit,
        totalPages: response.pagination.total_pages
    };
};
```

### 架构原则

#### SOLID原则
- **单一职责原则**: 每个类/模块只负责一个功能
- **开闭原则**: 对扩展开放，对修改关闭
- **里氏替换原则**: 接口实现可以互相替换
- **接口隔离原则**: 细粒度的接口设计
- **依赖倒置原则**: 依赖抽象而非具体实现


## 📁 项目结构

```
galaxyErp/
├── cmd/                    # 应用程序入口
│   └── server/            # 服务器启动程序
├── configs/               # 配置文件
│   ├── config.yaml       # 基础配置
│   ├── dev.yaml          # 开发环境配置
│   ├── test.yaml         # 测试环境配置
│   └── prod.yaml         # 生产环境配置
├── frontend/              # 前端应用
│   ├── app/              # Next.js 应用页面 (App Router)
│   ├── components/       # React 组件
│   ├── contexts/         # React Context (状态管理)
│   ├── services/         # API 服务层
│   ├── types/            # TypeScript 类型定义
│   └── lib/              # 工具库和API客户端
├── internal/              # 后端核心代码
│   ├── auth/             # 认证模块
│   ├── common/           # 公共组件和工具
│   ├── config/           # 配置管理
│   ├── container/        # 依赖注入容器
│   ├── controllers/      # 控制器层 (HTTP处理)
│   ├── dto/              # 数据传输对象
│   ├── handlers/         # 特殊处理器
│   ├── middleware/       # 中间件 (认证、日志、CORS等)
│   ├── models/           # 数据模型 (GORM实体)
│   ├── repositories/     # 数据访问层 (仓储模式)
│   ├── routes/           # 路由定义
│   ├── services/         # 业务逻辑层
│   └── utils/            # 工具函数
├── sql/                   # SQL 脚本和数据库迁移
├── go.mod                # Go 模块依赖
├── Makefile              # 构建脚本
└── README.md             # 项目文档
```

### 🏛️ 分层架构设计

#### 后端分层架构 (Clean Architecture)

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Controllers │  │ Middleware  │  │      Routes         │  │
│  │   (HTTP)    │  │ (Auth/CORS) │  │   (API Routing)     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Business Logic Layer                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Services   │  │     DTO     │  │    Validation       │  │
│  │ (Business)  │  │ (Transfer)  │  │     (Rules)         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Data Access Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │Repositories │  │   Models    │  │     Database        │  │
│  │ (Interface) │  │  (GORM)     │  │ (PostgreSQL/SQLite) │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

#### 前端分层架构 (Component-Based)

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Pages    │  │ Components  │  │       Layout        │  │
│  │ (App Router)│  │ (Ant Design)│  │    (Navigation)     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    State Management Layer                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Context   │  │    Hooks    │  │      Types          │  │
│  │ (Auth/App)  │  │ (useState)  │  │   (TypeScript)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Data Access Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Services   │  │ API Client  │  │    HTTP Client      │  │
│  │ (Business)  │  │ (Wrapper)   │  │   (Fetch API)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 📋 核心功能模块

### 💰 财务会计模块 (Accounting) - 🚧 待实现
- **科目管理**: 会计科目的创建、编辑和层级管理
- **日记账**: 财务交易记录和凭证管理
- **成本中心**: 成本核算和分析
- **银行账户**: 银行账户管理和对账
- **付款条目**: 付款记录和审批流程
- **预算管理**: 预算制定和执行监控
- **汇率管理**: 多币种汇率管理和转换
- **税务模板**: 税率配置和税务计算
- **会计期间**: 财务年度和会计期间管理

### 🛒 销售管理模块 (Sales) - ✅ 已实现
- **客户管理**: ✅ 客户信息维护和分类管理 (完整CRUD + 高级搜索 + 分页)
- **报价管理**: ✅ 销售报价单创建和跟踪 (CRUD + 状态管理 + 价格计算)
- **销售订单**: ✅ 订单处理和状态跟踪 (完整生命周期 + 工作流 + 审批)
- **发货单**: 🚧 发货记录和物流跟踪 (待实现)
- **销售分析**: 🚧 销售数据统计和趋势分析 (待实现)
- **技术特点**: RESTful API设计，数据验证，事务管理，状态机模式

### 🛍️ 采购管理模块 (Purchase) - ✅ 已实现
- **供应商管理**: ✅ 供应商信息和评估体系 (CRUD + 资质管理 + 评级系统)
- **采购申请**: ✅ 采购需求申请和审批 (完整工作流 + 多级审批 + 状态跟踪)
- **采购订单**: ✅ 采购订单管理和执行 (CRUD + 确认/取消 + 执行监控)
- **采购收货**: 🚧 收货确认和质量检验 (待实现)
- **采购分析**: ✅ 采购成本分析和供应商绩效 (统计API + 数据可视化)
- **采购合同**: 🚧 合同管理和条款跟踪 (待实现)
- **技术特点**: 工作流引擎，状态机模式，事务管理，并发控制

### 📦 库存管理模块 (Inventory) - ✅ 已实现
- **仓库管理**: ✅ 多仓库管理和库位设置 (CRUD + 层级结构 + 容量管理)
- **物料管理**: ✅ 物料信息和分类管理 (CRUD + 分类体系 + 规格参数)
- **库存跟踪**: ✅ 实时库存监控和预警 (实时查询 + 安全库存 + 自动预警)
- **库存移动**: ✅ 入库、出库、调拨等库存操作 (完整移动API + 批次管理 + 序列号)
- **库存盘点**: 🚧 定期盘点和差异处理 (待实现)
- **库存分析**: 🚧 库存周转率和成本分析 (部分实现)
- **技术特点**: 并发控制，库存锁定，事件驱动，ACID事务

### 🏭 生产管理模块 (Production) - 🔄 部分实现
- **生产计划**: 🚧 生产计划制定和排程 (待实现)
- **物料需求**: 🚧 MRP 物料需求计划 (待实现)
- **工艺路线**: 🚧 生产工艺和操作流程 (模型已定义)
- **工作中心**: 🚧 生产设备和产能管理 (模型已定义)
- **生产订单**: 🚧 生产任务下达和执行 (模型已定义)
- **产品管理**: ✅ 产品信息管理 (CRUD + 搜索 + BOM管理)
- **质量检验**: 🚧 质量控制和不合格品处理 (模型已定义)
- **设备管理**: 🚧 设备维护和故障管理 (模型已定义)
- **技术特点**: 复杂业务逻辑，数据模型完整，扩展性强

### 📊 项目管理模块 (Project) - ✅ 已实现
- **项目管理**: ✅ 项目创建和生命周期管理 (CRUD + 状态管理 + 进度跟踪)
- **里程碑**: ✅ 项目关键节点和进度跟踪 (CRUD + 时间线 + 依赖关系)
- **任务管理**: ✅ 项目任务分解和分配 (CRUD + 层级结构 + 责任分配)
- **时间记录**: ✅ 工时记录和成本核算 (CRUD + 时间统计 + 成本计算)
- **资源管理**: 🚧 项目资源分配和利用率 (待实现)
- **项目报告**: 🚧 项目进度和绩效报告 (待实现)
- **技术特点**: 权限矩阵，时间管理，成本核算，团队协作

### 👥 人力资源模块 (HR) - 🚧 待实现
- **员工管理**: 🚧 员工档案和基本信息维护 (路由已定义)
- **部门管理**: 🚧 组织架构和部门设置 (部分实现)
- **考勤管理**: 🚧 出勤记录和考勤统计 (路由已定义)
- **请假管理**: 🚧 请假申请和审批流程 (模型已定义)
- **加班管理**: 🚧 加班申请和工时统计 (模型已定义)
- **薪资管理**: 🚧 薪资计算和发放记录 (路由已定义)
- **绩效管理**: 🚧 绩效目标设定和评估 (模型已定义)
- **培训管理**: 🚧 培训计划和记录管理 (模型已定义)
- **技能管理**: 🚧 员工技能档案和评估 (待实现)

### ⚙️ 系统管理模块 (System) - 🔄 部分实现
- **用户管理**: ✅ 用户账户和权限管理 (注册/登录/个人资料 + JWT认证)
- **角色管理**: 🚧 角色定义和权限分配 (待实现)
- **系统配置**: 🚧 系统参数和业务规则配置 (待实现)
- **数据备份**: 🚧 系统数据备份和恢复 (待实现)
- **审计日志**: 🚧 系统操作记录和安全审计 (待实现)
- **系统监控**: 🚧 系统性能和运行状态监控 (待实现)
- **技术特点**: JWT认证，密码加密，会话管理，安全防护

### 📱 前端页面实现状态
- ✅ **主页**: 仪表板和最近活动展示 (响应式设计 + 数据可视化)
- ✅ **销售管理**: 报价单管理页面 (完整CRUD + 表格组件 + 搜索过滤)
- ✅ **采购管理**: 采购订单、供应商、采购请求页面 (工作流界面 + 状态展示)
- ✅ **库存管理**: 库存查询和移动操作页面 (实时数据 + 操作界面)
- ✅ **生产管理**: 工单、物料清单、生产计划页面 (模拟数据 + 界面框架)
- ✅ **项目管理**: 项目、任务、里程碑管理页面 (甘特图 + 时间线)
- ✅ **人力资源**: 基础页面框架 (组件结构 + 路由配置)
- 🚧 **财务管理**: 待实现
- 🚧 **系统管理**: 待实现
- **技术特点**: Ant Design组件，TypeScript类型安全，状态管理，API集成


## 📋 环境要求

### 后端环境
- **Go**: 1.24+ 
- **数据库**: SQLite 3.x (开发环境) / PostgreSQL 12+ (生产环境)

### 前端环境
- **Node.js**: 18.x+
- **npm**: 9.x+ 或 **yarn**: 1.22+

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/galaxyerp/galaxyErp.git
cd galaxyErp
```

### 2. 后端设置
```bash
# 安装 Go 依赖
go mod tidy

# 运行数据库迁移 (使用 SQLite，无需额外配置)
make migrate

# 启动后端服务 (默认端口: 8080)
make run
```

### 3. 前端设置
```bash
# 进入前端目录
cd frontend

# 安装前端依赖
npm install

# 启动前端开发服务器 (默认端口: 3000)
npm run dev
```

### 4. 访问系统
- **前端界面**: http://localhost:3000
- **后端API**: http://localhost:8080

### 5. 默认登录信息
开发环境已配置自动登录功能，如需手动登录请使用：
- **用户名**: `admin`
- **密码**: `password`

> 💡 **开发提示**: 开发环境下前端会自动使用默认账户登录，无需手动输入登录信息。生产环境请务必修改默认密码！

## ⚙️ 配置说明

系统支持三种运行环境，每种环境都有对应的配置文件：

### 🔧 环境配置

| 环境 | 配置文件 | 数据库 | 用途 |
|------|----------|--------|------|
| **开发环境 (dev)** | `configs/dev.yaml` | SQLite | 本地开发和调试 |
| **测试环境 (test)** | `configs/test.yaml` | PostgreSQL | 功能测试和集成测试 |
| **生产环境 (prod)** | `configs/prod.yaml` | PostgreSQL | 生产部署 |

### 📝 环境变量配置

开发环境使用默认配置，无需额外设置。生产环境需要配置环境变量：

```env
# 服务器配置
SERVER_PORT=8080

# JWT 密钥 (生产环境请使用强密钥)
JWT_SECRET=your_super_secret_jwt_key_here

# 生产环境数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=galaxyerp_user
DB_PASSWORD=your_secure_password
DB_NAME=galaxyerp_prod
DB_SSLMODE=require
```

> 📌 **注意**: 开发环境使用 SQLite 数据库，配置文件位于 `configs/dev.yaml`，无需额外配置即可运行。

## 🏃‍♂️ 运行应用

### 开发模式 (推荐)

```bash
# 1. 后端服务 (使用 SQLite，无需额外配置)
make migrate    # 运行数据库迁移
make run        # 启动后端服务 (http://localhost:8080)

# 2. 前端服务 (新终端窗口)
cd frontend
npm run dev     # 启动前端开发服务器 (http://localhost:3000)
```

访问 http://localhost:3000 即可使用系统，开发环境会自动登录。

### 测试环境

```bash
# 配置测试数据库后运行
make migrate-test
make run-test
```

### 生产环境

```bash
# 配置生产数据库后运行
make migrate-prod
make run-prod
```

### 手动运行

```bash
# 后端手动运行
GALAXYERP_ENV=dev go run cmd/migrate/main.go  # 迁移数据库
GALAXYERP_ENV=dev go run cmd/server/main.go   # 启动服务器

# 前端手动运行
cd frontend
npm run build   # 构建生产版本
npm start       # 启动生产服务器
```

## 📚 API 文档

### API 基础信息
- **基础URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **内容类型**: `application/json`

### 测试API
可以使用以下方式测试API：
- **curl 命令**: 参考下方API端点示例
- **Postman**: 导入API端点进行测试
- **前端界面**: 通过Web界面直接操作

### 已实现的API端点

#### 🔐 认证模块 (Auth)
```
POST   /api/v1/auth/register     # 用户注册
POST   /api/v1/auth/login        # 用户登录
GET    /api/v1/auth/me           # 获取当前用户信息
POST   /api/v1/auth/logout       # 用户登出
```

#### 👥 用户管理 (Users)
```
GET    /api/v1/users/profile     # 获取用户资料
PUT    /api/v1/users/profile     # 更新用户资料
PUT    /api/v1/users/password    # 修改密码
GET    /api/v1/users/            # 获取用户列表 (管理员)
DELETE /api/v1/users/:id         # 删除用户 (管理员)
POST   /api/v1/users/search      # 搜索用户
```

#### 🛒 销售管理 (Sales)
```
# 客户管理
POST   /api/v1/customers/        # 创建客户
GET    /api/v1/customers/        # 获取客户列表
GET    /api/v1/customers/:id     # 获取客户详情
PUT    /api/v1/customers/:id     # 更新客户
DELETE /api/v1/customers/:id     # 删除客户
POST   /api/v1/customers/search  # 搜索客户

# 销售订单
POST   /api/v1/sales-orders/     # 创建销售订单
GET    /api/v1/sales-orders/     # 获取订单列表
GET    /api/v1/sales-orders/:id  # 获取订单详情
PUT    /api/v1/sales-orders/:id  # 更新订单
DELETE /api/v1/sales-orders/:id  # 删除订单
PUT    /api/v1/sales-orders/:id/status  # 更新订单状态

# 报价单
POST   /api/v1/quotations/       # 创建报价单
GET    /api/v1/quotations/       # 获取报价单列表
GET    /api/v1/quotations/:id    # 获取报价单详情
PUT    /api/v1/quotations/:id    # 更新报价单
DELETE /api/v1/quotations/:id    # 删除报价单
GET    /api/v1/quotations/search # 搜索报价单
```

#### 🛍️ 采购管理 (Purchase)
```
# 供应商管理
POST   /api/v1/suppliers/        # 创建供应商
GET    /api/v1/suppliers/        # 获取供应商列表
GET    /api/v1/suppliers/:id     # 获取供应商详情
PUT    /api/v1/suppliers/:id     # 更新供应商
DELETE /api/v1/suppliers/:id     # 删除供应商

# 采购订单
POST   /api/v1/purchase-orders/  # 创建采购订单
GET    /api/v1/purchase-orders/  # 获取订单列表
GET    /api/v1/purchase-orders/:id  # 获取订单详情
PUT    /api/v1/purchase-orders/:id  # 更新订单
DELETE /api/v1/purchase-orders/:id  # 删除订单
POST   /api/v1/purchase-orders/:id/confirm  # 确认订单
POST   /api/v1/purchase-orders/:id/cancel   # 取消订单

# 采购申请
POST   /api/v1/purchase-requests/  # 创建采购申请
GET    /api/v1/purchase-requests/  # 获取申请列表
GET    /api/v1/purchase-requests/:id  # 获取申请详情
PUT    /api/v1/purchase-requests/:id  # 更新申请
DELETE /api/v1/purchase-requests/:id  # 删除申请
POST   /api/v1/purchase-requests/:id/submit   # 提交申请
POST   /api/v1/purchase-requests/:id/approve  # 审批申请
POST   /api/v1/purchase-requests/:id/reject   # 拒绝申请

# 采购统计
GET    /api/v1/purchase/stats     # 获取采购统计
```

#### 📦 库存管理 (Inventory)
```
# 物料管理
POST   /api/v1/items/            # 创建物料
GET    /api/v1/items/            # 获取物料列表
GET    /api/v1/items/:id         # 获取物料详情
PUT    /api/v1/items/:id         # 更新物料
DELETE /api/v1/items/:id         # 删除物料
POST   /api/v1/items/search      # 搜索物料

# 库存管理
GET    /api/v1/stocks/           # 获取库存列表
POST   /api/v1/stocks/           # 创建库存
GET    /api/v1/stocks/:id        # 获取库存详情
PUT    /api/v1/stocks/:id        # 更新库存
DELETE /api/v1/stocks/:id        # 删除库存
GET    /api/v1/stock/item/:item_id  # 按物料查询库存

# 库存移动
GET    /api/v1/stock-movements/  # 获取移动记录
POST   /api/v1/stock-movements/  # 创建移动记录
POST   /api/v1/stock-movements/in  # 入库操作
POST   /api/v1/stock-movements/out  # 出库操作
POST   /api/v1/stock-movements/adjustment  # 库存调整
POST   /api/v1/stock-movements/transfer    # 库存调拨

# 仓库管理
GET    /api/v1/warehouses/       # 获取仓库列表
POST   /api/v1/warehouses/       # 创建仓库
GET    /api/v1/warehouses/:id    # 获取仓库详情
PUT    /api/v1/warehouses/:id    # 更新仓库
DELETE /api/v1/warehouses/:id    # 删除仓库

# 库存报告
GET    /api/v1/inventory-reports/stats      # 库存统计
GET    /api/v1/inventory-reports/report     # 库存报告
GET    /api/v1/inventory-reports/abc-analysis  # ABC分析
GET    /api/v1/inventory-reports/export     # 导出报告
```

#### 🏭 生产管理 (Production)
```
# 产品管理
POST   /api/v1/products/         # 创建产品
GET    /api/v1/products/         # 获取产品列表
GET    /api/v1/products/:id      # 获取产品详情
PUT    /api/v1/products/:id      # 更新产品
DELETE /api/v1/products/:id      # 删除产品
POST   /api/v1/products/search   # 搜索产品
```

#### 📊 项目管理 (Project)
```
# 项目管理
POST   /api/v1/projects/         # 创建项目
GET    /api/v1/projects/         # 获取项目列表
GET    /api/v1/projects/:id      # 获取项目详情
PUT    /api/v1/projects/:id      # 更新项目
DELETE /api/v1/projects/:id      # 删除项目

# 任务管理
POST   /api/v1/tasks/            # 创建任务
GET    /api/v1/tasks/            # 获取任务列表
GET    /api/v1/tasks/:id         # 获取任务详情
PUT    /api/v1/tasks/:id         # 更新任务
DELETE /api/v1/tasks/:id         # 删除任务

# 里程碑管理
POST   /api/v1/milestones/       # 创建里程碑
GET    /api/v1/milestones/:id    # 获取里程碑详情
PUT    /api/v1/milestones/:id    # 更新里程碑
DELETE /api/v1/milestones/:id    # 删除里程碑
GET    /api/v1/project-milestones/:project_id  # 获取项目里程碑

# 工时记录
POST   /api/v1/time-entries/     # 创建工时记录
GET    /api/v1/time-entries/:id  # 获取工时记录
PUT    /api/v1/time-entries/:id  # 更新工时记录
DELETE /api/v1/time-entries/:id  # 删除工时记录
GET    /api/v1/project-time-entries/:project_id  # 获取项目工时
```

## 🧪 API 测试

### 快速测试

```bash
# 1. 健康检查
curl http://localhost:8080/health

# 2. 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@galaxyerp.com",
    "password": "admin123",
    "first_name": "系统",
    "last_name": "管理员"
  }'

# 3. 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 4. 获取用户信息 (需要token)
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 业务模块测试示例

```bash
# 获取客户列表 (需要认证)
curl -X GET "http://localhost:8080/api/v1/customers/?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 创建供应商
curl -X POST http://localhost:8080/api/v1/suppliers/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "code": "SUP001",
    "name": "测试供应商",
    "contact_name": "张三",
    "email": "supplier@example.com",
    "phone": "13800138000",
    "address": "北京市朝阳区"
  }'

# 获取库存列表
curl -X GET "http://localhost:8080/api/v1/stocks/?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 创建库存移动记录
curl -X POST http://localhost:8080/api/v1/stock-movements/in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "item_id": 1,
    "warehouse_id": 1,
    "quantity": 100,
    "unit_cost": 10.50,
    "notes": "采购入库"
  }'

# 获取项目列表
curl -X GET "http://localhost:8080/api/v1/projects/?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 认证说明

大部分API端点需要JWT认证，请在请求头中包含：
```
Authorization: Bearer YOUR_JWT_TOKEN
# 使用登录获取的 JWT Token
export TOKEN="your_jwt_token_here"

# 创建客户
curl -X POST http://localhost:8080/api/v1/customers/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试客户",
    "code": "CUST001",
    "email": "customer@test.com",
    "phone": "13800138000"
  }'

# 创建供应商
curl -X POST http://localhost:8080/api/v1/suppliers/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试供应商",
    "code": "SUP001",
    "email": "supplier@test.com",
    "phone": "13900139000"
  }'

# 创建仓库
curl -X POST http://localhost:8080/api/v1/warehouses/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "主仓库",
    "code": "WH001",
    "address": "北京市朝阳区"
  }'
```

### 推荐测试工具
- **Postman**: 导入提供的 Postman 集合文件
- **curl**: 命令行快速测试
- **HTTPie**: 更友好的命令行工具
- **Insomnia**: 现代化的 API 测试工具

## ✨ 项目特性

- **🎨 现代化界面**: 基于 Ant Design 的响应式 UI 设计
- **🔐 安全认证**: JWT 令牌认证和权限管理
- **📱 移动友好**: 支持移动设备访问和操作
- **🌐 多环境支持**: 开发、测试、生产环境配置
- **📊 实时数据**: 实时业务数据统计和分析
- **🔄 RESTful API**: 标准化的 API 接口设计
- **📝 完整文档**: 详细的 API 文档和使用指南
- **🧪 测试支持**: 完整的 API 测试集合
- **🚀 高性能**: Go 语言高并发处理能力
- **🔧 易于扩展**: 模块化架构，便于功能扩展

## 🚀 部署指南

### 环境要求

#### 后端环境
- **Go**: 1.24+ (推荐使用最新版本)
- **数据库**: PostgreSQL 13+ 或 SQLite 3.35+
- **操作系统**: Linux/macOS/Windows

#### 前端环境
- **Node.js**: 18.0+ (推荐 LTS 版本)
- **npm**: 9.0+ 或 **yarn**: 1.22+
- **浏览器**: Chrome 90+, Firefox 88+, Safari 14+

### 开发环境部署

#### 1. 克隆项目
```bash
git clone https://github.com/your-username/galaxy-erp.git
cd galaxy-erp
```

#### 2. 后端配置与启动

```bash
# 安装Go依赖
go mod download

# 复制配置文件
cp config/dev.yaml.example config/dev.yaml

# 编辑配置文件 (可选，默认使用SQLite)
vim config/dev.yaml
```

**配置文件说明** (`config/dev.yaml`):
```yaml
server:
  port: 8080
  mode: debug

database:
  # SQLite配置 (默认)
  driver: sqlite
  dsn: "./data/galaxy_erp.db"
  
  # PostgreSQL配置 (可选)
  # driver: postgres
  # dsn: "host=localhost user=postgres password=password dbname=galaxy_erp port=5432 sslmode=disable"

jwt:
  secret: "your-secret-key-change-in-production"
  expire_hours: 24

cors:
  allowed_origins: ["http://localhost:3000"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
```

```bash
# 运行数据库迁移
make migrate

# 启动后端服务
make run
# 或者直接运行
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 3. 前端配置与启动

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install
# 或使用 yarn
yarn install

# 启动开发服务器
npm run dev
# 或使用 yarn
yarn dev
```

前端服务将在 `http://localhost:3000` 启动

#### 4. 验证部署

- **后端健康检查**: http://localhost:8080/health
- **前端页面**: http://localhost:3000
- **API文档**: http://localhost:8080/api/docs (如果配置了Swagger)

### 生产环境部署

#### 方式一：Docker 部署 (推荐)

1. **准备 Docker 环境**
```bash
# 确保安装了 Docker 和 Docker Compose
docker --version
docker-compose --version
```

2. **配置生产环境**
```bash
# 复制生产配置
cp config/prod.yaml.example config/prod.yaml

# 编辑生产配置
vim config/prod.yaml
```

3. **使用 Docker Compose 部署**
```bash
# 构建并启动所有服务
docker-compose -f docker-compose.prod.yml up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 方式二：手动部署

1. **后端构建与部署**
```bash
# 构建后端二进制文件
make build
# 或者
CGO_ENABLED=0 GOOS=linux go build -o bin/galaxy-erp cmd/server/main.go

# 复制文件到服务器
scp bin/galaxy-erp user@server:/opt/galaxy-erp/
scp -r config/ user@server:/opt/galaxy-erp/
scp -r migrations/ user@server:/opt/galaxy-erp/

# 在服务器上运行
./galaxy-erp
```

2. **前端构建与部署**
```bash
# 构建前端静态文件
cd frontend
npm run build
# 或
yarn build

# 部署到 Nginx
sudo cp -r .next/static/* /var/www/galaxy-erp/
sudo cp -r public/* /var/www/galaxy-erp/
```

3. **Nginx 配置示例**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /var/www/galaxy-erp;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 健康检查
    location /health {
        proxy_pass http://localhost:8080;
    }
}
```

### 数据库配置

#### SQLite (开发环境推荐)
```yaml
database:
  driver: sqlite
  dsn: "./data/galaxy_erp.db"
```

#### PostgreSQL (生产环境推荐)
```bash
# 安装 PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# 创建数据库和用户
sudo -u postgres psql
CREATE DATABASE galaxy_erp;
CREATE USER galaxy_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE galaxy_erp TO galaxy_user;
\q
```

```yaml
database:
  driver: postgres
  dsn: "host=localhost user=galaxy_user password=your_password dbname=galaxy_erp port=5432 sslmode=disable"
```

### 环境变量配置

可以使用环境变量覆盖配置文件设置：

```bash
# 服务器配置
export SERVER_PORT=8080
export SERVER_MODE=release

# 数据库配置
export DB_DRIVER=postgres
export DB_DSN="host=localhost user=galaxy_user password=your_password dbname=galaxy_erp port=5432 sslmode=disable"

# JWT配置
export JWT_SECRET="your-production-secret-key"
export JWT_EXPIRE_HOURS=24

# CORS配置
export CORS_ALLOWED_ORIGINS="https://your-domain.com"
```

### 常用命令

```bash
# 后端相关
make run          # 启动开发服务器
make build        # 构建生产版本
make test         # 运行测试
make migrate      # 运行数据库迁移
make clean        # 清理构建文件

# 前端相关
npm run dev       # 启动开发服务器
npm run build     # 构建生产版本
npm run start     # 启动生产服务器
npm run lint      # 代码检查
npm run test      # 运行测试

# Docker相关
docker-compose up -d              # 启动所有服务
docker-compose down               # 停止所有服务
docker-compose logs -f            # 查看日志
docker-compose exec backend sh   # 进入后端容器
```

### 故障排除

#### 常见问题

1. **端口冲突**
   - 检查端口是否被占用：`lsof -i :8080`
   - 修改配置文件中的端口号

2. **数据库连接失败**
   - 检查数据库服务是否启动
   - 验证连接字符串配置
   - 检查防火墙设置

3. **前端无法访问后端API**
   - 检查CORS配置
   - 验证API基础URL配置
   - 检查网络连接

4. **JWT认证失败**
   - 检查JWT密钥配置
   - 验证token是否过期
   - 检查请求头格式

#### 日志查看

```bash
# 后端日志
tail -f logs/app.log

# Docker日志
docker-compose logs -f backend
docker-compose logs -f frontend

# 系统日志
journalctl -u galaxy-erp -f
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！请遵循以下步骤：

1. **Fork** 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 **Pull Request**

### 开发规范

- **Go代码**: 遵循 `gofmt` 格式化标准
- **TypeScript代码**: 遵循 ESLint 配置
- **提交信息**: 使用清晰的提交信息描述更改
- **测试**: 为新功能添加相应的测试用例

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。


## 🙏 致谢

感谢所有为 GalaxyERP 项目做出贡献的开发者和用户！

---

<div align="center">
  <strong>🌟 如果这个项目对您有帮助，请给我们一个 Star！🌟</strong>
</div>