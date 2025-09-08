# 采购小程序后端API服务

## 📖 项目简介

这是一个为采购小程序提供后端API服务的Go项目，采用现代化的RESTful API设计，支持商店管理、商品展示、购物车操作和订单处理等核心功能。

- **项目名称**: 采购小程序后端API
- **技术栈**: Go + Gin + Nginx  
- **部署方式**: 宝塔面板 + 域名访问
- **API地址**: https://www.ency.asia/api

## 🚀 快速开始

### 本地开发
```bash
# 1. 克隆项目
git clone https://github.com/ape7054/purches-backend.git
cd purches-backend

# 2. 启动开发环境
# Linux/Mac
./scripts/dev.sh

# Windows  
scripts\dev.bat

# 3. 访问API
# 本地: http://localhost:8080/api
# 线上: https://www.ency.asia/api
```

### 部署到服务器
```bash
# SSH到服务器
ssh root@47.76.90.227
cd /root/purches-backend
./scripts/deploy.sh
```

## 📁 项目结构

```
purches-backend/
├── main.go                    # 主程序入口和API路由
├── models/                    # 数据模型定义
├── database/                  # 数据库相关
├── scripts/                   # 开发和部署脚本
│   ├── dev.sh                # Linux开发脚本
│   ├── dev.bat               # Windows开发脚本
│   ├── deploy.sh             # 部署脚本
│   └── README.md             # 脚本使用说明
├── docs/                      # 文档目录
│   ├── development.md        # 开发指南
│   └── windows.md            # Windows开发说明
├── go.mod                     # Go模块依赖
├── go.sum                     # 依赖版本锁定
├── server.log                 # 服务运行日志
├── API.md                     # API接口文档
└── README.md                  # 项目说明
```

## 🛠️ 技术栈

- **后端框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能Go Web框架
- **数据存储**: SQLite数据库
- **反向代理**: Nginx
- **SSL证书**: Let's Encrypt
- **服务器管理**: 宝塔面板
- **部署平台**: Ubuntu 服务器

## 🔗 API接口文档

### 基地址
```
https://www.ency.asia/api
```

### 统一响应格式
所有API接口都采用统一的JSON响应格式：

```json
{
  "code": 200,           // 状态码：200-成功，4xx-客户端错误，5xx-服务器错误
  "message": "操作成功",  // 响应消息
  "data": { ... }        // 具体数据（可能为null）
}
```

### 主要接口
- `GET /shops` - 获取商店列表
- `GET /shops/{shopId}/products` - 获取商店商品列表
- `GET /cart` - 查看购物车
- `POST /cart/items` - 添加商品到购物车
- `PUT /cart/items/{productId}` - 更新购物车商品数量
- `DELETE /cart/items` - 删除购物车商品
- `POST /orders` - 提交订单
- `GET /health` - 健康检查

详细API文档请查看 [API.md](API.md)

## 🚀 部署架构

```
用户请求 → Nginx (443/80) → 反向代理 → Go应用 (8080端口)
```

### 部署环境
- **服务器IP**: 47.76.90.227
- **域名**: www.ency.asia
- **管理面板**: 宝塔面板

## 📚 开发文档

- [开发指南](docs/development.md) - 本地开发环境搭建
- [Windows开发说明](docs/windows.md) - Windows特殊配置
- [脚本使用说明](scripts/README.md) - 开发和部署脚本
- [API接口文档](API.md) - 详细API文档

## 🔄 版本更新

### v1.0.0 (当前版本)
- ✅ 基础API功能完成
- ✅ 商店和商品管理
- ✅ 购物车操作
- ✅ 订单创建
- ✅ HTTPS部署
- ✅ 域名访问
- ✅ 本地开发环境
- ✅ 自动部署脚本

### 未来计划
- [ ] 用户认证和授权
- [ ] 日志监控系统
- [ ] API限流和缓存
- [ ] 单元测试覆盖
- [ ] Docker容器化

## 📞 联系方式

- **服务器**: 47.76.90.227
- **API地址**: https://www.ency.asia/api
- **GitHub**: https://github.com/ape7054/purches-backend

---

**🎯 项目状态**: ✅ 生产就绪，支持本地开发 