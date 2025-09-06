# 采购小程序 - 后端API接口文档

## 🔗 基础信息

**API基地址**: `https://www.ency.asia/api`  
**协议**: HTTPS  
**数据格式**: JSON  
**字符编码**: UTF-8  

---

## 📋 统一响应格式

所有接口都返回统一的JSON格式：

```json
{
  "code": 200,           // 状态码
  "message": "操作成功", // 响应消息 
  "data": { ... }        // 具体数据
}
```

**状态码说明**：
- `200` - 成功
- `400` - 请求参数错误
- `404` - 资源不存在
- `500` - 服务器错误

---

## 🏪 商店相关接口

### 1. 获取商店列表

**接口地址**: `GET /shops`

**请求示例**:
```javascript
fetch('https://www.ency.asia/api/shops')
  .then(res => res.json())
  .then(data => console.log(data))
```

**响应数据**:
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
    },
    {
      "id": "shop_3",
      "name": "肖红梅",
      "logo": ""
    }
  ]
}
```

### 2. 获取商店商品列表

**接口地址**: `GET /shops/{shopId}/products`

**请求参数**:
- `shopId` (路径参数): 商店ID，如 `shop_1`

**请求示例**:
```javascript
// 获取 shop_1 的商品列表
fetch('https://www.ency.asia/api/shops/shop_1/products')
  .then(res => res.json())
  .then(data => console.log(data))
```

**响应数据**:
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
    },
    {
      "id": "prod_3",
      "name": "鸡精", 
      "price": 12.00,
      "imageUrl": ""
    }
  ]
}
```

---

## 🛒 购物车相关接口

### 3. 查看购物车

**接口地址**: `GET /cart`

**请求示例**:
```javascript
fetch('https://www.ency.asia/api/cart')
  .then(res => res.json())
  .then(data => console.log(data))
```

**响应数据**:
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

### 4. 添加商品到购物车

**接口地址**: `POST /cart/items`

**请求参数**:
```json
{
  "productId": "prod_1",  // 商品ID（必填）
  "quantity": 2           // 数量（必填，最小值1）
}
```

**请求示例**:
```javascript
fetch('https://www.ency.asia/api/cart/items', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    productId: 'prod_1',
    quantity: 2
  })
})
.then(res => res.json())
.then(data => console.log(data))
```

**响应数据**:
```json
{
  "code": 200,
  "message": "添加成功",
  "data": null
}
```

### 5. 修改购物车商品数量

**接口地址**: `PUT /cart/items/{productId}`

**请求参数**:
- `productId` (路径参数): 商品ID
- 请求体:
```json
{
  "quantity": 5  // 新的数量（必填，最小值0，0表示删除）
}
```

**请求示例**:
```javascript
// 修改 prod_1 的数量为 5
fetch('https://www.ency.asia/api/cart/items/prod_1', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    quantity: 5
  })
})
.then(res => res.json())
.then(data => console.log(data))
```

**响应数据**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

### 6. 删除购物车商品

**接口地址**: `DELETE /cart/items`

**请求参数**:
```json
{
  "productIds": ["prod_1", "prod_2"]  // 要删除的商品ID数组
}
```

**请求示例**:
```javascript
fetch('https://www.ency.asia/api/cart/items', {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    productIds: ['prod_1', 'prod_2']
  })
})
.then(res => res.json())
.then(data => console.log(data))
```

**响应数据**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

---

## 📦 订单相关接口

### 7. 提交订单

**接口地址**: `POST /orders`

**请求参数**:
```json
{
  "remark": "请尽快配送"  // 订单备注（可选）
}
```

**请求示例**:
```javascript
fetch('https://www.ency.asia/api/orders', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    remark: '请尽快配送'
  })
})
.then(res => res.json())
.then(data => console.log(data))
```

**响应数据**:
```json
{
  "code": 200,
  "message": "订单创建成功",
  "data": {
    "orderId": "order_1699123456"
  }
}
```

---

## 🔧 系统接口

### 8. 健康检查

**接口地址**: `GET /health`

**请求示例**:
```javascript
fetch('https://www.ency.asia/api/health')
  .then(res => res.json())
  .then(data => console.log(data))
```

**响应数据**:
```json
{
  "code": 200,
  "message": "服务正常",
  "data": "OK"
}
```

---

## 📱 uni-app 调用示例

### 获取商店列表
```javascript
uni.request({
  url: 'https://www.ency.asia/api/shops',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 200) {
      this.shops = res.data.data
    }
  },
  fail: (err) => {
    console.error('请求失败:', err)
  }
})
```

### 获取商品列表
```javascript
uni.request({
  url: `https://www.ency.asia/api/shops/${shopId}/products`,
  method: 'GET',
  success: (res) => {
    if (res.data.code === 200) {
      this.products = res.data.data
    }
  }
})
```

### 添加到购物车
```javascript
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
    if (res.data.code === 200) {
      uni.showToast({
        title: '添加成功',
        icon: 'success'
      })
    }
  }
})
```

### 查看购物车
```javascript
uni.request({
  url: 'https://www.ency.asia/api/cart',
  method: 'GET',
  success: (res) => {
    if (res.data.code === 200) {
      this.cartItems = res.data.data.items
      this.totalPrice = res.data.data.totalPrice
    }
  }
})
```

### 提交订单
```javascript
uni.request({
  url: 'https://www.ency.asia/api/orders',
  method: 'POST',
  header: {
    'Content-Type': 'application/json'
  },
  data: {
    remark: '请尽快配送'
  },
  success: (res) => {
    if (res.data.code === 200) {
      uni.showToast({
        title: '订单提交成功',
        icon: 'success'
      })
      // 跳转到订单详情或清空购物车
    }
  }
})
```

---

## ⚠️ 重要说明

### 数据存储说明
- 当前版本使用**内存存储**
- **重启服务后数据会重置**
- 购物车数据会丢失，这是正常现象

### 测试建议
1. 先调用健康检查接口确认服务正常
2. 按照业务流程测试：商店 → 商品 → 购物车 → 订单
3. 每次添加商品后查看购物车验证数据正确性
4. 测试各种边界情况（空购物车提交订单等）

### 错误处理
```javascript
// 统一错误处理示例
function handleApiResponse(res) {
  if (res.data.code === 200) {
    return res.data.data
  } else {
    console.error('API错误:', res.data.message)
    uni.showToast({
      title: res.data.message,
      icon: 'error'
    })
    return null
  }
}
```

---

## 📞 联系方式

**API状态**: ✅ 生产环境可用  
**服务器**: 47.76.90.227  
**后端负责人**: [你的联系方式]

如有问题请及时联系！ 