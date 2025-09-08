# 本地开发指南

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/ape7054/purches-backend.git
cd purches-backend
```

### 2. 启动开发环境
```bash
# Linux/Mac
./scripts/dev.sh

# Windows
scripts\dev.bat
```

### 3. 访问服务
- API地址: http://localhost:8080/api
- 健康检查: http://localhost:8080/api/health

## 📝 开发流程

### 日常开发
```bash
# 1. 启动开发服务器
./scripts/dev.sh

# 2. 修改代码后自动重启...
# 3. 测试API功能
# 4. 提交代码
git add .
git commit -m "feat: 新功能描述"
git push origin main
```

### 部署到服务器
```bash
# SSH到服务器执行
ssh root@47.76.90.227
cd /root/purches-backend
./scripts/deploy.sh
```

## 🔧 常用命令

```bash
# 安装依赖
go mod tidy

# 代码检查
go fmt ./...
go vet ./...

# 手动启动
go run main.go

# 编译
go build -o purches-backend main.go
```

## 🌐 环境对比

| 项目 | 本地开发 | 生产环境 |
|------|----------|----------|
| 地址 | http://localhost:8080 | https://www.ency.asia |
| 数据库 | 本地SQLite | 服务器SQLite |
| SSL | 无 | Let's Encrypt |

## 🐛 常见问题

### 端口被占用
```bash
# Linux/Mac
lsof -i :8080
kill -9 进程ID

# Windows
netstat -ano | findstr :8080
taskkill /PID 进程ID /F
```

### 依赖问题
```bash
go mod download
go mod tidy
```

详细的Windows开发说明请查看 [Windows开发指南](windows.md) 