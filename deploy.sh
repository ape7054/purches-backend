#!/bin/bash

# 采购小程序后端自动部署脚本
# 使用方法: ./deploy.sh

set -e  # 遇到错误立即退出

echo "🚀 开始部署采购小程序后端..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_DIR="/root/purches-backend"
SERVICE_PORT="8080"
LOG_FILE="server.log"

# 1. 检查是否在正确目录
if [ ! -f "main.go" ]; then
    echo -e "${RED}❌ 错误: 请在项目根目录执行此脚本${NC}"
    exit 1
fi

# 2. 拉取最新代码
echo -e "${YELLOW}📦 拉取最新代码...${NC}"
git fetch origin
git pull origin main

# 3. 停止旧的服务
echo -e "${YELLOW}⏹️  停止旧服务...${NC}"
pkill -f "go run main.go" || echo "没有运行中的服务"
pkill -f "purches-backend" || echo "没有运行中的二进制服务"

# 等待进程完全停止
sleep 2

# 4. 清理旧的编译文件
echo -e "${YELLOW}🧹 清理旧文件...${NC}"
rm -f purches-backend

# 5. 安装/更新依赖
echo -e "${YELLOW}📚 更新依赖...${NC}"
go mod tidy

# 6. 编译项目
echo -e "${YELLOW}🔨 编译项目...${NC}"
go build -o purches-backend main.go

# 7. 备份数据库（如果存在）
if [ -f "purches.db" ]; then
    echo -e "${YELLOW}💾 备份数据库...${NC}"
    cp purches.db "purches.db.backup.$(date +%Y%m%d_%H%M%S)"
fi

# 8. 启动新服务
echo -e "${YELLOW}🚀 启动新服务...${NC}"
nohup ./purches-backend > $LOG_FILE 2>&1 &

# 获取新进程ID
sleep 2
NEW_PID=$(pgrep -f "purches-backend" | head -1)

if [ -n "$NEW_PID" ]; then
    echo -e "${GREEN}✅ 服务启动成功! PID: $NEW_PID${NC}"
else
    echo -e "${RED}❌ 服务启动失败!${NC}"
    exit 1
fi

# 9. 健康检查
echo -e "${YELLOW}🏥 健康检查...${NC}"
sleep 3

# 检查端口是否监听
if netstat -tlnp | grep ":$SERVICE_PORT" > /dev/null; then
    echo -e "${GREEN}✅ 端口 $SERVICE_PORT 正常监听${NC}"
else
    echo -e "${RED}❌ 端口 $SERVICE_PORT 未监听${NC}"
    exit 1
fi

# 检查API是否正常
if curl -f -s http://localhost:$SERVICE_PORT/api/health > /dev/null; then
    echo -e "${GREEN}✅ API健康检查通过${NC}"
else
    echo -e "${RED}❌ API健康检查失败${NC}"
    exit 1
fi

# 10. 显示部署信息
echo -e "${GREEN}"
echo "========================================"
echo "🎉 部署完成!"
echo "========================================"
echo "服务地址: https://www.ency.asia/api"
echo "进程ID: $NEW_PID"
echo "日志文件: $LOG_FILE"
echo "========================================"
echo -e "${NC}"

# 11. 显示最新日志
echo -e "${YELLOW}📝 最新日志:${NC}"
tail -10 $LOG_FILE 