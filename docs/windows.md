# Windows开发说明

## 🪟 环境准备

### 必要工具
- [Go 1.18+](https://golang.org/dl/) - 下载 `.msi` 安装包
- [Git for Windows](https://git-scm.com/download/win) - 推荐选择 "Git Bash"
- [VS Code](https://code.visualstudio.com/) + Go扩展

### 验证安装
```cmd
go version
git --version
```

## 🚀 Windows开发

### 启动方式
```cmd
# 方式一: 双击 scripts\dev.bat

# 方式二: 命令行
cd purches-backend
scripts\dev.bat

# 方式三: 使用Git Bash (推荐)
./scripts/dev.sh
```

## 🔄 跨平台兼容性

### ✅ 完全兼容
- Go代码 - 跨平台运行
- API接口 - 完全一致  
- SQLite数据库 - 跨平台
- Git操作 - 完全一致

### 🔧 Windows特殊配置

#### 文件路径
```go
// ✅ 跨平台兼容写法
import "path/filepath"
dbPath := filepath.Join("data", "purches.db")

// ❌ 避免硬编码路径
dbPath := "data\\purches.db"  // 仅Windows
dbPath := "data/purches.db"   // 仅Unix
```

#### 跨平台编译
```cmd
# 在Windows编译Linux版本
set GOOS=linux
set GOARCH=amd64
go build -o purches-backend main.go
```

## 🛠️ 推荐工具

### 终端
- **Git Bash** - 类Unix体验，兼容性最好
- **Windows Terminal** - 微软官方，现代化
- **PowerShell** - 功能强大

### SSH连接
```bash
# 使用Git Bash或PowerShell
ssh root@47.76.90.227
```

## 🐛 Windows常见问题

### 中文乱码
```cmd
chcp 65001  # 设置UTF-8编码
```

### 端口被占用
```cmd
netstat -ano | findstr :8080
taskkill /PID 进程ID /F
```

### PowerShell API测试
```powershell
# 健康检查
Invoke-RestMethod -Uri "http://localhost:8080/api/health"

# 获取商店列表
Invoke-RestMethod -Uri "http://localhost:8080/api/shops"
```

---
**✅ Windows开发Linux部署完全无障碍！** 