# 采购小程序后端API服务

## 📖 项目简介

这是一个为采购小程序提供后端API服务的Go项目，采用现代化的RESTful API设计，支持商店管理、商品展示、购物车操作和订单处理等核心功能。

- **项目名称**: 采购小程序后端API
- **技术栈**: Go + Gin + Nginx  
- **部署方式**: 宝塔面板 + 域名访问
- **API地址**: https://www.ency.asia/api

## 🛠️ 技术栈

- **后端框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能Go Web框架
- **数据存储**: 内存存储（生产环境可扩展为数据库）
- **反向代理**: Nginx
- **SSL证书**: Let's Encrypt
- **服务器管理**: 宝塔面板
- **部署平台**: Ubuntu 服务器

## 📁 项目结构

```
purches-backend/
├── main.go              # 主程序入口和API路由定义
├── models/              # 数据模型定义
│   └── models.go        # 所有数据结构和请求响应模型
├── database/            # 数据库相关（预留扩展）
│   └── database.go      # 数据库连接和初始化
├── go.mod               # Go模块依赖
├── go.sum               # 依赖版本锁定
├── server.log           # 服务运行日志
└── README.md            # 项目文档
```

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

### 接口列表

#### 1. 商店相关接口

##### 1.1 获取商店列表
- **URL**: `GET /shops`
- **描述**: 获取所有可用的商店列表
- **响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "shop_1",
      "name": "快驴",
      "logo": ""
    },
    {
      "id": "shop_2", 
      "name": "华兴街14号",
      "logo": ""
    }
  ]
}
```

##### 1.2 获取商店商品列表
- **URL**: `GET /shops/{shopId}/products`
- **描述**: 获取指定商店的所有商品
- **路径参数**:
  - `shopId`: 商店ID（如：shop_1）
- **响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "prod_1",
      "name": "盐",
      "price": 5.00,
      "imageUrl": ""
    },
    {
      "id": "prod_2",
      "name": "味精", 
      "price": 8.00,
      "imageUrl": ""
    }
  ]
}
```

#### 2. 购物车相关接口

##### 2.1 查看购物车
- **URL**: `GET /cart`
- **描述**: 获取当前用户的购物车内容
- **响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "items": [
      {
        "productId": "prod_1",
        "productName": "盐",
        "quantity": 2,
        "price": 5.00
      }
    ],
    "totalPrice": 10.00
  }
}
```

##### 2.2 添加商品到购物车
- **URL**: `POST /cart/items`
- **描述**: 将商品添加到购物车
- **请求体**:
```json
{
  "productId": "prod_1",
  "quantity": 2
}
```
- **响应示例**:
```json
{
  "code": 200,
  "message": "添加成功",
  "data": null
}
```

##### 2.3 更新购物车商品数量
- **URL**: `PUT /cart/items/{productId}`
- **描述**: 修改购物车中指定商品的数量
- **路径参数**:
  - `productId`: 商品ID
- **请求体**:
```json
{
  "quantity": 5
}
```
- **响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

##### 2.4 删除购物车商品
- **URL**: `DELETE /cart/items`
- **描述**: 从购物车中删除指定商品
- **请求体**:
```json
{
  "productIds": ["prod_1", "prod_2"]
}
```
- **响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

#### 3. 订单相关接口

##### 3.1 提交订单
- **URL**: `POST /orders`
- **描述**: 将购物车中的商品生成订单
- **请求体**:
```json
{
  "remark": "请尽快配送"
}
```
- **响应示例**:
```json
{
  "code": 200,
  "message": "订单创建成功",
  "data": {
    "orderId": "order_1699123456"
  }
}
```

#### 4. 系统接口

##### 4.1 健康检查
- **URL**: `GET /health`
- **描述**: 检查服务器运行状态
- **响应示例**:
```json
{
  "code": 200,
  "message": "服务正常",
  "data": "OK"
}
```

## 🚀 部署说明

### 服务器环境
- **操作系统**: Ubuntu
- **服务器IP**: 47.76.90.227
- **域名**: www.ency.asia
- **管理面板**: 宝塔面板

### 部署步骤

1. **代码部署**
   ```bash
   cd /root/purches-backend
   go mod tidy
   go run main.go &
   ```

2. **域名配置**
   - 在宝塔面板添加网站: www.ency.asia
   - 配置反向代理: http://127.0.0.1:8080
   - 代理目录: /api

3. **SSL证书**
   - 申请Let's Encrypt免费证书
   - 开启强制HTTPS访问

### 服务管理

#### 启动服务
```bash
cd /root/purches-backend
nohup go run main.go > server.log 2>&1 &
```

#### 查看服务状态
```bash
ps aux | grep "go run" | grep -v grep
```

#### 查看服务日志
```bash
tail -f /root/purches-backend/server.log
```

#### 停止服务
```bash
pkill -f "go run main.go"
```

## 🔧 开发指南

### 前端联调

#### JavaScript示例
```javascript
// API基地址
const API_BASE = 'https://www.ency.asia/api'

// 获取商店列表
async function getShops() {
  try {
    const response = await fetch(`${API_BASE}/shops`)
    const data = await response.json()
    if (data.code === 200) {
      console.log('商店列表:', data.data)
    }
  } catch (error) {
    console.error('请求失败:', error)
  }
}

// 添加商品到购物车
async function addToCart(productId, quantity) {
  try {
    const response = await fetch(`${API_BASE}/cart/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        productId: productId,
        quantity: quantity
      })
    })
    const data = await response.json()
    console.log('添加结果:', data)
  } catch (error) {
    console.error('添加失败:', error)
  }
}
```

#### uni-app示例
```javascript
// 获取商店列表
uni.request({
  url: 'https://www.ency.asia/api/shops',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 200) {
      console.log('商店列表:', res.data.data)
    }
  },
  fail: (err) => {
    console.error('请求失败:', err)
  }
})

// 添加商品到购物车
uni.request({
  url: 'https://www.ency.asia/api/cart/items',
  method: 'POST',
  header: {
    'Content-Type': 'application/json'
  },
  data: {
    productId: 'prod_1',
    quantity: 2
  },
  success: (res) => {
    console.log('添加结果:', res.data)
  }
})
```

## 🐛 错误码说明

| 错误码 | 说明 | 常见原因 |
|--------|------|----------|
| 200 | 操作成功 | - |
| 400 | 请求参数错误 | 缺少必要参数或参数格式错误 |
| 404 | 资源不存在 | 商店ID或商品ID不存在 |
| 500 | 服务器内部错误 | 服务器异常 |

## 📊 性能说明

- **响应时间**: < 100ms
- **并发支持**: 基于Go协程，支持高并发
- **数据存储**: 当前为内存存储，重启服务数据会重置

## 🔄 版本更新

### v1.0.0 (当前版本)
- ✅ 基础API功能完成
- ✅ 商店和商品管理
- ✅ 购物车操作
- ✅ 订单创建
- ✅ HTTPS部署
- ✅ 域名访问

### 未来计划
- [ ] 数据库持久化存储
- [ ] 用户认证和授权
- [ ] 日志监控系统
- [ ] API限流和缓存
- [ ] 单元测试覆盖

## 📞 联系方式

- **项目负责人**: [你的名字]
- **技术支持**: [你的联系方式]
- **服务器**: 47.76.90.227
- **API地址**: https://www.ency.asia/api

## 📝 更新日志

### 2025-09-06
- 🎉 项目初始化完成
- 🚀 成功部署到生产环境
- 🔒 配置SSL证书
- 📋 完成API文档

---

**🎯 项目状态**: ✅ 生产就绪，可开始前端联调 