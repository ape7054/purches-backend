package main

import (
	"fmt"
	"net/http"
	"purches-backend/database"
	"purches-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// 固定用户ID（简化版本，实际项目中应该从token获取）
const DEFAULT_USER_ID = "user_1"

// 获取商店列表
func getShops(c *gin.Context) {
	var shops []models.Shop
	result := database.DB.Find(&shops)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "获取商店列表失败: " + result.Error.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    shops,
	})
}

// 获取指定商店的商品列表
func getShopProducts(c *gin.Context) {
	shopID := c.Param("shopId")

	// 先检查商店是否存在
	var shop models.Shop
	if err := database.DB.First(&shop, "id = ?", shopID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "商店不存在",
			Data:    nil,
		})
		return
	}

	// 获取该商店的商品
	var products []models.Product
	database.DB.Where("shop_id = ?", shopID).Find(&products)

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    products,
	})
}

// 获取购物车内容
func getCart(c *gin.Context) {
	// 简化版本：使用固定用户ID，实际项目中应该从token中获取
	userID := DEFAULT_USER_ID

	// 确保用户存在
	var user models.User
	if err := database.DB.FirstOrCreate(&user, models.User{OpenID: userID}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "用户初始化失败",
			Data:    nil,
		})
		return
	}

	// 获取购物车商品
	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", user.ID).Find(&cartItems)

	// 计算总价并转换为响应格式
	var totalPrice float64
	var responseItems []models.CartItemResponse

	for _, item := range cartItems {
		totalPrice += item.Price * float64(item.Quantity)
		responseItems = append(responseItems, models.CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	response := models.CartResponse{
		Items:      responseItems,
		TotalPrice: totalPrice,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    response,
	})
}

// 添加商品到购物车
func addToCart(c *gin.Context) {
	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 查找商品信息
	var product models.Product
	if err := database.DB.First(&product, "id = ?", req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "商品不存在",
			Data:    nil,
		})
		return
	}

	// 简化版本：使用固定用户ID
	userID := DEFAULT_USER_ID

	// 确保用户存在
	var user models.User
	if err := database.DB.FirstOrCreate(&user, models.User{OpenID: userID}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "用户初始化失败",
			Data:    nil,
		})
		return
	}

	// 检查购物车中是否已存在该商品
	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", user.ID, req.ProductID).First(&existingItem)

	if result.Error == nil {
		// 如果存在，增加数量
		existingItem.Quantity += req.Quantity
		database.DB.Save(&existingItem)
	} else {
		// 如果不存在，添加新商品
		newItem := models.CartItem{
			UserID:      user.ID,
			ProductID:   req.ProductID,
			ProductName: product.Name,
			Quantity:    req.Quantity,
			Price:       product.Price,
		}
		database.DB.Create(&newItem)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "添加成功",
		Data:    nil,
	})
}

// 更新购物车商品数量
func updateCartItem(c *gin.Context) {
	productID := c.Param("productId")

	var req models.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 简化版本：使用固定用户ID
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "用户不存在",
			Data:    nil,
		})
		return
	}

	// 查找并更新商品
	var cartItem models.CartItem
	if err := database.DB.Where("user_id = ? AND product_id = ?", user.ID, productID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "购物车中没有该商品",
			Data:    nil,
		})
		return
	}

	if req.Quantity > 0 {
		cartItem.Quantity = req.Quantity
		database.DB.Save(&cartItem)
	} else {
		// 数量为0时删除商品
		database.DB.Delete(&cartItem)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "更新成功",
		Data:    nil,
	})
}

// 从购物车删除商品
func deleteCartItems(c *gin.Context) {
	var req models.DeleteCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 简化版本：使用固定用户ID
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "用户不存在",
			Data:    nil,
		})
		return
	}

	// 删除指定的商品
	database.DB.Where("user_id = ? AND product_id IN ?", user.ID, req.ProductIds).Delete(&models.CartItem{})

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "删除成功",
		Data:    nil,
	})
}

// 提交订单
func createOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	// 允许空请求体
	_ = c.ShouldBindJSON(&req)

	// 简化版本：使用固定用户ID
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "用户不存在",
			Data:    nil,
		})
		return
	}

	// 获取购物车商品
	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", user.ID).Find(&cartItems)

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "购物车为空，无法创建订单",
			Data:    nil,
		})
		return
	}

	// 计算总价
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	// 生成订单ID
	orderID := fmt.Sprintf("order_%d", time.Now().Unix())

	// 创建订单
	order := models.Order{
		ID:         orderID,
		UserID:     user.ID,
		TotalPrice: totalPrice,
		Status:     "pending",
		Remark:     req.Remark,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "创建订单失败",
			Data:    nil,
		})
		return
	}

	// 创建订单商品明细
	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:     orderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
		database.DB.Create(&orderItem)
	}

	// 清空购物车
	database.DB.Where("user_id = ?", user.ID).Delete(&models.CartItem{})

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "订单创建成功",
		Data:    models.CreateOrderResponse{OrderID: orderID},
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
	// 初始化数据库
	database.InitDatabase()

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
			c.JSON(http.StatusOK, models.APIResponse{
				Code:    200,
				Message: "服务正常",
				Data:    "OK",
			})
		})
	}

	fmt.Println("🚀 服务器启动成功!")
	fmt.Println("💾 使用SQLite数据库存储")
	fmt.Println("🔗 后端地址: http://localhost:8080")
	fmt.Println("🔍 健康检查: http://localhost:8080/api/health")

	// 启动服务器
	r.Run(":8080")
}
