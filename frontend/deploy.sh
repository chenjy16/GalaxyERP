#!/bin/bash

# Galaxy ERP Frontend 部署脚本
# 使用方法: ./deploy.sh [static|docker|standalone]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        log_error "npm 未安装"
        exit 1
    fi
    
    log_info "依赖检查完成"
}

# 清理构建文件
clean_build() {
    log_info "清理构建文件..."
    npm run clean
}

# 安装依赖
install_dependencies() {
    log_info "安装依赖..."
    npm ci --production
}

# 类型检查
type_check() {
    log_info "执行类型检查..."
    npm run type-check
}

# 代码检查
lint_check() {
    log_info "执行代码检查..."
    npm run lint
}

# 静态导出部署
deploy_static() {
    log_info "开始静态导出部署..."
    
    # 使用专用的静态导出构建命令
    npm run build:static
    
    log_info "静态文件已生成到 out/ 目录"
    log_info "可以将 out/ 目录部署到任何静态托管服务"
    log_info ""
    log_info "部署选项:"
    log_info "1. CDN: 上传 out/ 目录到 AWS S3, Cloudflare Pages 等"
    log_info "2. Nginx: 配置 nginx 指向 out/ 目录"
    log_info "3. 静态托管: Vercel, Netlify, GitHub Pages 等"
    log_info ""
    log_info "预览命令: npm run preview:static"
}

# Docker 部署
deploy_docker() {
    log_info "开始 Docker 部署..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装"
        exit 1
    fi
    
    # 修改 next.config.js 启用 standalone
    sed -i.bak 's|// output: '\''standalone'\''|output: '\''standalone'\''|g' next.config.js
    
    # 构建 Docker 镜像
    docker build -t galaxy-erp-frontend:latest .
    
    # 恢复配置文件
    mv next.config.js.bak next.config.js
    
    log_info "Docker 镜像构建完成: galaxy-erp-frontend:latest"
    log_info "使用以下命令启动容器:"
    log_info "docker run -p 3000:3000 galaxy-erp-frontend:latest"
}

# Standalone 部署
deploy_standalone() {
    log_info "开始 Standalone 部署..."
    
    # 修改 next.config.js 启用 standalone
    sed -i.bak 's|// output: '\''standalone'\''|output: '\''standalone'\''|g' next.config.js
    
    # 构建
    npm run build
    
    # 恢复配置文件
    mv next.config.js.bak next.config.js
    
    log_info "Standalone 构建完成"
    log_info "可以使用 'node .next/standalone/server.js' 启动服务"
}

# 主函数
main() {
    local deploy_type=${1:-"static"}
    
    log_info "Galaxy ERP Frontend 部署开始..."
    log_info "部署类型: $deploy_type"
    
    # 检查依赖
    check_dependencies
    
    # 清理构建文件
    clean_build
    
    # 安装依赖
    install_dependencies
    
    # 类型检查
    type_check
    
    # 代码检查
    lint_check
    
    # 根据类型部署
    case $deploy_type in
        "static")
            deploy_static
            ;;
        "docker")
            deploy_docker
            ;;
        "standalone")
            deploy_standalone
            ;;
        *)
            log_error "未知的部署类型: $deploy_type"
            log_info "支持的部署类型: static, docker, standalone"
            exit 1
            ;;
    esac
    
    log_info "部署完成!"
}

# 执行主函数
main "$@"