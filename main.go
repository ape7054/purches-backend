package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Shop 商店模型（临时简化版）
type Shop struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// Product 商品模型（临时简化版）
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"imageUrl"`
}

// CartItem 购物车商品
type CartItem struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// CartResponse 购物车响应
type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"totalPrice"`
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"` // Allow 0 to effectively delete
}

// DeleteCartRequest 删除购物车请求
type DeleteCartRequest struct {
	ProductIds []string `json:"productIds" binding:"required"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Remark string `json:"remark"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

// 临时数据存储
var shops = []Shop{
	{ID: "shop_1", Name: "快驴", Logo: ""},
	{ID: "shop_2", Name: "华兴街14号", Logo: ""},
	{ID: "shop_3", Name: "肖红梅", Logo: ""},
}

var products = map[string][]Product{
	"shop_1": {
		{ID: "prod_1", Name: "盐", Price: 5.00, ImageURL: ""},
		{ID: "prod_2", Name: "味精", Price: 8.00, ImageURL: ""},
		{ID: "prod_3", Name: "鸡精", Price: 12.00, ImageURL: ""},
	},
	"shop_2": {
		{ID: "prod_6", Name: "大豆油", Price: 25.00, ImageURL: ""},
		{ID: "prod_7", Name: "大米", Price: 30.00, ImageURL: ""},
	},
}

// 内存存储的购物车数据 (userID -> CartItem[])
// 简化版本，实际项目中应该用数据库
var userCarts = make(map[string][]CartItem)

// 简单的产品查找功能
func findProductByID(productID string) *Product {
	for _, shopProducts := range products {
		for _, product := range shopProducts {
			if product.ID == productID {
				return &product
			}
		}
	}
	return nil
}

// 获取商店列表
func getShops(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    shops,
	})
}

// 获取指定商店的商品列表
func getShopProducts(c *gin.Context) {
	shopID := c.Param("shopId")

	shopProducts, exists := products[shopID]
	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "商店不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    shopProducts,
	})
}

// 获取购物车内容
func getCart(c *gin.Context) {
	// 简化版本：使用固定用户ID，实际项目中应该从token中获取
	userID := "user_1"
	
	cartItems := userCarts[userID]
	if cartItems == nil {
		cartItems = []CartItem{}
	}
	
	// 计算总价
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += item.Price * float64(item.Quantity)
	}
	
	response := CartResponse{
		Items:      cartItems,
		TotalPrice: totalPrice,
	}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    response,
	})
}

// 添加商品到购物车
func addToCart(c *gin.Context) {
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	// 查找商品信息
	product := findProductByID(req.ProductID)
	if product == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "商品不存在",
			Data:    nil,
		})
		return
	}
	
	// 简化版本：使用固定用户ID
	userID := "user_1"
	
	if userCarts[userID] == nil {
		userCarts[userID] = []CartItem{}
	}
	
	// 检查购物车中是否已存在该商品
	found := false
	for i, item := range userCarts[userID] {
		if item.ProductID == req.ProductID {
			// 如果存在，增加数量
			userCarts[userID][i].Quantity += req.Quantity
			found = true
			break
		}
	}
	
	// 如果不存在，添加新商品
	if !found {
		newItem := CartItem{
			ProductID:   req.ProductID,
			ProductName: product.Name,
			Quantity:    req.Quantity,
			Price:       product.Price,
		}
		userCarts[userID] = append(userCarts[userID], newItem)
	}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "添加成功",
		Data:    nil,
	})
}

// 更新购物车商品数量
func updateCartItem(c *gin.Context) {
	productID := c.Param("productId")
	
	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	// 简化版本：使用固定用户ID
	userID := "user_1"
	
	if userCarts[userID] == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "购物车中没有该商品",
			Data:    nil,
		})
		return
	}
	
	// 查找并更新商品数量
	found := false
	var newCart []CartItem
	for _, item := range userCarts[userID] {
		if item.ProductID == productID {
			if req.Quantity > 0 {
				item.Quantity = req.Quantity
				newCart = append(newCart, item)
			}
			found = true
		} else {
			newCart = append(newCart, item)
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "购物车中没有该商品",
			Data:    nil,
		})
		return
	}
	
    userCarts[userID] = newCart

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "更新成功",
		Data:    nil,
	})
}

// 从购物车删除商品
func deleteCartItems(c *gin.Context) {
	var req DeleteCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	// 简化版本：使用固定用户ID
	userID := "user_1"
	
	if userCarts[userID] == nil {
		c.JSON(http.StatusOK, APIResponse{
			Code:    200,
			Message: "删除成功",
			Data:    nil,
		})
		return
	}
	
	// 删除指定的商品
	var newCart []CartItem
	for _, item := range userCarts[userID] {
		shouldDelete := false
		for _, deleteID := range req.ProductIds {
			if item.ProductID == deleteID {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			newCart = append(newCart, item)
		}
	}
	
	userCarts[userID] = newCart
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "删除成功",
		Data:    nil,
	})
}

// 提交订单
func createOrder(c *gin.Context) {
	var req CreateOrderRequest
	// 允许空请求体
	_ = c.ShouldBindJSON(&req)
	
	// 简化版本：使用固定用户ID
	userID := "user_1"
	
	cartItems := userCarts[userID]
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "购物车为空，无法创建订单",
			Data:    nil,
		})
		return
	}
	
	// 生成订单ID（简化版本）
	orderID := fmt.Sprintf("order_%d", time.Now().Unix())
	
	// 清空购物车
	userCarts[userID] = []CartItem{}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "订单创建成功",
		Data:    CreateOrderResponse{OrderID: orderID},
	})
}


// 设置CORS中间件
func setupCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	// 创建Gin实例
	r := gin.Default()

	// 设置CORS
	r.Use(setupCORS())

	// API路由组
	api := r.Group("/api")
	{
		// 商店相关路由
		api.GET("/shops", getShops)
		api.GET("/shops/:shopId/products", getShopProducts)

		// 购物车相关路由
		api.GET("/cart", getCart)
		api.POST("/cart/items", addToCart)
		api.PUT("/cart/items/:productId", updateCartItem)
		api.DELETE("/cart/items", deleteCartItems)
		
		// 订单相关路由
		api.POST("/orders", createOrder)

		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, APIResponse{
				Code:    200,
				Message: "服务正常",
				Data:    "OK",
			})
		})
	}

	fmt.Println("🚀 服务器启动成功!")
	fmt.Println("🔗 后端地址: http://YOUR_SERVER_IP:8080")
	fmt.Println("🔍 健康检查: http://YOUR_SERVER_IP:8080/api/health")
	
	// 启动服务器
	r.Run(":8080")
}