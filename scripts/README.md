# 脚本使用说明

## 📂 脚本列表

### 开发脚本
- `dev.sh` - Linux/Mac本地开发启动脚本
- `dev.bat` - Windows本地开发启动脚本

### 部署脚本  
- `deploy.sh` - 服务器自动部署脚本

## 🚀 使用方法

### 本地开发
```bash
# Linux/Mac
./scripts/dev.sh

# Windows
scripts\dev.bat
```

### 服务器部署
```bash
# 在服务器上执行
./scripts/deploy.sh
```

## ⚡ 脚本功能

### dev.sh/dev.bat
- ✅ 检查Go环境
- ✅ 安装依赖 `go mod tidy`
- ✅ 启动开发服务器
- ✅ 跨平台兼容

### deploy.sh
- ✅ 拉取最新代码
- ✅ 停止旧服务
- ✅ 编译新版本
- ✅ 备份数据库
- ✅ 启动新服务
- ✅ 健康检查

## 🔧 注意事项

1. **执行权限**: Linux脚本需要执行权限
   ```bash
   chmod +x scripts/*.sh
   ```

2. **路径问题**: 脚本会自动切换到项目根目录

3. **依赖检查**: 脚本会检查Go环境和项目文件 