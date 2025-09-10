# 采购订单系统 - API接口文档

## 概述

本文档描述了采购订单系统的后端API接口，用于支持Web前端的所有功能。系统主要包含商品管理、购物车管理、订单管理和供应商管理四大模块。

## 基础信息

- **API基础URL**: `https://api.yourcompany.com/v1`
- **数据格式**: JSON
- **字符编码**: UTF-8
- **时间格式**: ISO 8601 (例如: `2025-09-10T14:30:00.000Z`)

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 具体数据内容
  },
  "timestamp": "2025-09-10T14:30:00.000Z"
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误描述",
  "error": "详细错误信息",
  "timestamp": "2025-09-10T14:30:00.000Z"
}
```

## 数据模型

### 商品 (Product)
```json
{
  "id": 1,
  "name": "牛蛙",
  "price": 31.00,
  "unit": "斤",
  "description": "牛蛙杀好处理干净去掉内脏和眼睛，去掉爪子，50斤",
  "supplier": "F35",
  "status": "available",
  "createdAt": "2025-09-10T14:30:00.000Z",
  "updatedAt": "2025-09-10T14:30:00.000Z"
}
```

### 购物车项 (CartItem)
```json
{
  "id": 1,
  "productId": 1,
  "name": "牛蛙",
  "count": 2,
  "shopName": "F35",
  "price": 31.00,
  "totalPrice": 62.00,
  "addedAt": "2025-09-10T14:30:00.000Z"
}
```

### 订单 (Order)
```json
{
  "id": "ORD202509101430001",
  "supplier": "F35",
  "products": [
    {
      "productId": 1,
      "name": "牛蛙",
      "description": "牛蛙杀好处理干净去掉内脏和眼睛，去掉爪子，50斤",
      "count": 2,
      "unit": "斤",
      "price": 31.00,
      "totalPrice": 62.00
    }
  ],
  "totalPrice": 62.00,
  "finalPrice": 60.00,
  "status": "pending",
  "notes": "明天早上 8 点前送到，谢谢！",
  "createdAt": "2025-09-10T14:30:00.000Z",
  "updatedAt": "2025-09-10T14:30:00.000Z"
}
```

### 供应商 (Supplier)
```json
{
  "name": "F35",
  "contactPerson": "张先生",
  "phone": "13800138000",
  "address": "批发市场F35号",
  "productCount": 25,
  "totalOrders": 156,
  "status": "active",
  "createdAt": "2025-09-10T14:30:00.000Z"
}
```

## API接口列表

## 1. 商品管理 API

### 1.1 获取商品列表
- **URL**: `GET /products`
- **描述**: 获取所有商品列表，支持分页和筛选
- **参数**:
  ```json
  {
    "page": 1,          // 页码，默认1
    "limit": 50,        // 每页数量，默认50
    "supplier": "F35",  // 可选：按供应商筛选
    "search": "牛蛙",   // 可选：搜索关键词
    "status": "available" // 可选：商品状态
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "products": [
        {
          "id": 1,
          "name": "牛蛙",
          "price": 31.00,
          "unit": "斤",
          "description": "牛蛙杀好处理干净去掉内脏和眼睛，去掉爪子，50斤",
          "supplier": "F35",
          "status": "available"
        }
      ],
      "pagination": {
        "total": 82,
        "page": 1,
        "limit": 50,
        "totalPages": 2
      }
    }
  }
  ```

### 1.2 获取商品详情
- **URL**: `GET /products/{productId}`
- **描述**: 获取指定商品的详细信息
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "id": 1,
      "name": "牛蛙",
      "price": 31.00,
      "unit": "斤",
      "description": "牛蛙杀好处理干净去掉内脏和眼睛，去掉爪子，50斤",
      "supplier": "F35",
      "status": "available",
      "createdAt": "2025-09-10T14:30:00.000Z",
      "updatedAt": "2025-09-10T14:30:00.000Z"
    }
  }
  ```

## 2. 购物车管理 API

### 2.1 获取购物车
- **URL**: `GET /cart`
- **描述**: 获取当前用户的购物车内容
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "items": [
        {
          "id": 1,
          "productId": 1,
          "name": "牛蛙",
          "count": 2,
          "shopName": "F35",
          "price": 31.00,
          "totalPrice": 62.00
        }
      ],
      "summary": {
        "totalItems": 2,
        "totalPrice": 62.00,
        "supplierCount": 1
      }
    }
  }
  ```

### 2.2 添加商品到购物车
- **URL**: `POST /cart/items`
- **描述**: 添加商品到购物车或更新数量
- **请求体**:
  ```json
  {
    "productId": 1,
    "count": 2
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "添加成功",
    "data": {
      "id": 1,
      "productId": 1,
      "name": "牛蛙",
      "count": 2,
      "shopName": "F35",
      "price": 31.00,
      "totalPrice": 62.00
    }
  }
  ```

### 2.3 更新购物车商品数量
- **URL**: `PUT /cart/items/{itemId}`
- **描述**: 更新购物车中商品的数量
- **请求体**:
  ```json
  {
    "count": 3
  }
  ```

### 2.4 删除购物车商品
- **URL**: `DELETE /cart/items/{itemId}`
- **描述**: 从购物车中删除指定商品

### 2.5 清空购物车
- **URL**: `DELETE /cart`
- **描述**: 清空当前用户的购物车

## 3. 订单管理 API

### 3.1 提交订单
- **URL**: `POST /orders`
- **描述**: 将购物车中的商品提交为订单，自动按供应商分组
- **请求体**:
  ```json
  {
    "items": [
      {
        "productId": 1,
        "count": 2
      }
    ],
    "notes": "明天早上 8 点前送到，谢谢！"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "订单提交成功",
    "data": {
      "orders": [
        {
          "id": "ORD202509101430001",
          "supplier": "F35",
          "products": [
            {
              "productId": 1,
              "name": "牛蛙",
              "count": 2,
              "unit": "斤",
              "price": 31.00,
              "totalPrice": 62.00
            }
          ],
          "totalPrice": 62.00,
          "status": "pending",
          "createdAt": "2025-09-10T14:30:00.000Z"
        }
      ]
    }
  }
  ```

### 3.2 获取订单列表
- **URL**: `GET /orders`
- **描述**: 获取订单列表，按供应商分组展示
- **参数**:
  ```json
  {
    "page": 1,
    "limit": 20,
    "supplier": "F35",    // 可选：按供应商筛选
    "status": "pending"   // 可选：按状态筛选
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "orders": [
        {
          "id": "ORD202509101430001",
          "supplier": "F35",
          "products": [...],
          "totalPrice": 62.00,
          "status": "pending",
          "createdAt": "2025-09-10T14:30:00.000Z"
        }
      ],
      "suppliers": [
        {
          "supplier": "F35",
          "orderCount": 1,
          "totalPrice": 62.00,
          "productCount": 2
        }
      ]
    }
  }
  ```

### 3.3 获取订单详情
- **URL**: `GET /orders/{orderId}`
- **描述**: 获取指定订单的详细信息

### 3.4 更新订单状态
- **URL**: `PUT /orders/{orderId}/status`
- **描述**: 更新订单状态
- **请求体**:
  ```json
  {
    "status": "completed",
    "notes": "已完成配送"
  }
  ```

### 3.5 更新订单最终价格
- **URL**: `PUT /orders/{orderId}/final-price`
- **描述**: 更新订单的最终结算价格
- **请求体**:
  ```json
  {
    "finalPrice": 60.00,
    "reason": "供应商优惠"
  }
  ```

## 4. 供应商管理 API

### 4.1 获取供应商列表
- **URL**: `GET /suppliers`
- **描述**: 获取所有供应商列表及统计信息
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "suppliers": [
        {
          "name": "F35",
          "contactPerson": "张先生",
          "phone": "13800138000",
          "productCount": 25,
          "totalOrders": 156,
          "status": "active"
        }
      ]
    }
  }
  ```

### 4.2 获取供应商详情
- **URL**: `GET /suppliers/{supplierName}`
- **描述**: 获取指定供应商的详细信息和订单统计
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "supplier": {
        "name": "F35",
        "contactPerson": "张先生",
        "phone": "13800138000",
        "address": "批发市场F35号"
      },
      "statistics": {
        "productCount": 25,
        "totalOrders": 156,
        "totalAmount": 15600.00,
        "averageOrderAmount": 100.00
      },
      "recentOrders": [...]
    }
  }
  ```

### 4.3 获取供应商的商品列表
- **URL**: `GET /suppliers/{supplierName}/products`
- **描述**: 获取指定供应商的所有商品

### 4.4 获取供应商的订单列表
- **URL**: `GET /suppliers/{supplierName}/orders`
- **描述**: 获取指定供应商的所有订单

## 5. 数据同步 API

### 5.1 批量导入商品
- **URL**: `POST /products/import`
- **描述**: 批量导入商品数据
- **请求体**:
  ```json
  {
    "products": [
      {
        "name": "商品名称",
        "price": 10.00,
        "unit": "斤",
        "description": "商品描述",
        "supplier": "供应商名称"
      }
    ]
  }
  ```

### 5.2 导出订单数据
- **URL**: `GET /orders/export`
- **描述**: 导出订单数据，支持Excel格式
- **参数**:
  ```json
  {
    "format": "excel",     // excel, csv
    "dateFrom": "2025-09-01",
    "dateTo": "2025-09-30",
    "supplier": "F35"      // 可选
  }
  ```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |

## 状态说明

### 商品状态 (Product Status)
- `available`: 可购买
- `unavailable`: 暂时缺货
- `discontinued`: 已停售

### 订单状态 (Order Status)
- `pending`: 待处理
- `confirmed`: 已确认
- `delivering`: 配送中
- `completed`: 已完成
- `cancelled`: 已取消

## 开发注意事项

1. **时区处理**: 所有时间字段使用UTC时间，前端需要根据用户时区进行转换
2. **价格精度**: 所有价格字段使用2位小数，避免浮点数精度问题
3. **数据验证**: 前端需要对用户输入进行基础验证，后端也要进行完整验证
4. **错误处理**: 前端需要根据错误码进行相应的用户提示
5. **性能优化**: 商品列表等大数据量接口建议使用分页，避免一次性加载过多数据

## 数据库设计建议

### 主要表结构

```sql
-- 商品表
CREATE TABLE products (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  unit VARCHAR(50) NOT NULL,
  description TEXT,
  supplier VARCHAR(255) NOT NULL,
  status ENUM('available','unavailable','discontinued') DEFAULT 'available',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 购物车表
CREATE TABLE cart_items (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(255) NOT NULL,
  product_id INT NOT NULL,
  count INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (product_id) REFERENCES products(id)
);

-- 订单表
CREATE TABLE orders (
  id VARCHAR(50) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  supplier VARCHAR(255) NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,
  final_price DECIMAL(10,2),
  status ENUM('pending','confirmed','delivering','completed','cancelled') DEFAULT 'pending',
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 订单商品表
CREATE TABLE order_items (
  id INT PRIMARY KEY AUTO_INCREMENT,
  order_id VARCHAR(50) NOT NULL,
  product_id INT NOT NULL,
  count INT NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);

-- 供应商表
CREATE TABLE suppliers (
  name VARCHAR(255) PRIMARY KEY,
  contact_person VARCHAR(255),
  phone VARCHAR(20),
  address TEXT,
  status ENUM('active','inactive') DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

*文档版本: v1.0*  
*最后更新: 2025年9月10日*  
*技术支持: 开发团队* 