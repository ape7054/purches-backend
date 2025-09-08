#!/bin/bash

# 本地开发启动脚本
# 使用方法: ./scripts/dev.sh

echo "🚀 启动本地开发环境..."

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go 环境未安装，请先安装 Go"
    exit 1
fi

# 检查项目文件
if [ ! -f "main.go" ]; then
    echo "❌ 找不到main.go文件"
    exit 1
fi

# 安装依赖
echo -e "${YELLOW}📚 安装依赖...${NC}"
go mod tidy

# 启动开发服务器
echo -e "${GREEN}✅ 启动开发服务器 (http://localhost:8080)${NC}"
echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"
echo "=================================="

# 启动并实时显示日志
go run main.go 