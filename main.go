package main

import (
	"fmt"
	"math"
	"net/http"
	"purches-backend/database"
	"purches-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 固定用户ID（简化版本，实际项目中应该从token获取）
const DEFAULT_USER_ID = "user_1"

// 生成统一响应
func responseOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func responseError(c *gin.Context, code int, message string, err string) {
	c.JSON(code, models.APIResponse{
		Code:      code,
		Message:   message,
		Data:      err,
		Timestamp: time.Now(),
	})
}

// =================== 商品管理 API ===================

// 获取商品列表
func getProducts(c *gin.Context) {
	var req models.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}

	// 构建查询
	query := database.DB.Model(&models.Product{})

	// 筛选条件
	if req.Supplier != "" {
		query = query.Where("supplier = ?", req.Supplier)
	}
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var products []models.Product
	offset := (req.Page - 1) * req.Limit
	query.Offset(offset).Limit(req.Limit).Find(&products)

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	response := models.ProductListResponse{
		Products: products,
		Pagination: models.PaginationResponse{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	responseOK(c, "获取成功", response)
}

// 获取商品详情
func getProduct(c *gin.Context) {
	productID := c.Param("productId")
	id, err := strconv.Atoi(productID)
	if err != nil {
		responseError(c, 400, "商品ID格式错误", err.Error())
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		responseError(c, 404, "商品不存在", err.Error())
		return
	}

	responseOK(c, "获取成功", product)
}

// =================== 购物车管理 API ===================

// 获取购物车
func getCart(c *gin.Context) {
	userID := DEFAULT_USER_ID

	// 确保用户存在
	var user models.User
	if err := database.DB.FirstOrCreate(&user, models.User{OpenID: userID}).Error; err != nil {
		responseError(c, 500, "用户初始化失败", err.Error())
		return
	}

	// 获取购物车商品
	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", user.ID).Find(&cartItems)

	// 计算统计信息
	var totalPrice float64
	supplierMap := make(map[string]bool)

	for _, item := range cartItems {
		totalPrice += item.TotalPrice
		supplierMap[item.ShopName] = true
	}

	response := models.CartResponse{
		Items: cartItems,
		Summary: models.CartSummary{
			TotalItems:    len(cartItems),
			TotalPrice:    totalPrice,
			SupplierCount: len(supplierMap),
		},
	}

	responseOK(c, "获取成功", response)
}

// 添加商品到购物车
func addToCart(c *gin.Context) {
	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	// 查找商品信息
	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		responseError(c, 404, "商品不存在", err.Error())
		return
	}

	userID := DEFAULT_USER_ID

	// 确保用户存在
	var user models.User
	if err := database.DB.FirstOrCreate(&user, models.User{OpenID: userID}).Error; err != nil {
		responseError(c, 500, "用户初始化失败", err.Error())
		return
	}

	// 检查购物车中是否已存在该商品
	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", user.ID, req.ProductID).First(&existingItem)

	if result.Error == nil {
		// 如果存在，增加数量
		existingItem.Count += req.Count
		existingItem.TotalPrice = float64(existingItem.Count) * existingItem.Price
		database.DB.Save(&existingItem)
		responseOK(c, "添加成功", existingItem)
	} else {
		// 如果不存在，添加新商品
		newItem := models.CartItem{
			UserID:     user.ID,
			ProductID:  req.ProductID,
			Name:       product.Name,
			Count:      req.Count,
			ShopName:   product.Supplier,
			Price:      product.Price,
			TotalPrice: product.Price * float64(req.Count),
			AddedAt:    time.Now(),
		}
		if err := database.DB.Create(&newItem).Error; err != nil {
			responseError(c, 500, "添加失败", err.Error())
			return
		}
		responseOK(c, "添加成功", newItem)
	}
}

// 更新购物车商品数量
func updateCartItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		responseError(c, 400, "商品ID格式错误", err.Error())
		return
	}

	var req models.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		responseError(c, 404, "用户不存在", err.Error())
		return
	}

	// 查找购物车商品
	var cartItem models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", itemID, user.ID).First(&cartItem).Error; err != nil {
		responseError(c, 404, "购物车中没有该商品", err.Error())
		return
	}

	if req.Count > 0 {
		cartItem.Count = req.Count
		cartItem.TotalPrice = float64(cartItem.Count) * cartItem.Price
		database.DB.Save(&cartItem)
	} else {
		// 数量为0时删除商品
		database.DB.Delete(&cartItem)
	}

	responseOK(c, "更新成功", nil)
}

// 删除购物车商品
func deleteCartItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		responseError(c, 400, "商品ID格式错误", err.Error())
		return
	}

	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		responseError(c, 404, "用户不存在", err.Error())
		return
	}

	// 删除商品
	result := database.DB.Where("id = ? AND user_id = ?", itemID, user.ID).Delete(&models.CartItem{})
	if result.RowsAffected == 0 {
		responseError(c, 404, "购物车中没有该商品", "")
		return
	}

	responseOK(c, "删除成功", nil)
}

// 清空购物车
func clearCart(c *gin.Context) {
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		responseError(c, 404, "用户不存在", err.Error())
		return
	}

	database.DB.Where("user_id = ?", user.ID).Delete(&models.CartItem{})
	responseOK(c, "清空成功", nil)
}

// =================== 订单管理 API ===================

// 提交订单（按供应商分组）
func createOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		responseError(c, 404, "用户不存在", err.Error())
		return
	}

	// 获取商品信息并按供应商分组
	supplierGroups := make(map[string][]models.OrderItemRequest)
	productMap := make(map[int]models.Product)

	for _, item := range req.Items {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			responseError(c, 404, fmt.Sprintf("商品ID %d 不存在", item.ProductID), err.Error())
			return
		}
		productMap[item.ProductID] = product
		supplierGroups[product.Supplier] = append(supplierGroups[product.Supplier], item)
	}

	var createdOrders []models.Order

	// 为每个供应商创建订单
	for supplier, items := range supplierGroups {
		orderID := fmt.Sprintf("ORD%d%03d", time.Now().Unix(), len(createdOrders)+1)

		var totalPrice float64
		var orderItems []models.OrderItem

		for _, item := range items {
			product := productMap[item.ProductID]
			itemTotal := product.Price * float64(item.Count)
			totalPrice += itemTotal

			orderItem := models.OrderItem{
				OrderID:     orderID,
				ProductID:   item.ProductID,
				Name:        product.Name,
				Description: product.Description,
				Count:       item.Count,
				Unit:        product.Unit,
				Price:       product.Price,
				TotalPrice:  itemTotal,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			orderItems = append(orderItems, orderItem)
		}

		// 创建订单
		order := models.Order{
			ID:         orderID,
			UserID:     user.ID,
			Supplier:   supplier,
			TotalPrice: totalPrice,
			Status:     "pending",
			Notes:      req.Notes,
			Products:   orderItems,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := database.DB.Create(&order).Error; err != nil {
			responseError(c, 500, "创建订单失败", err.Error())
			return
		}

		// 创建订单商品明细
		for _, orderItem := range orderItems {
			database.DB.Create(&orderItem)
		}

		createdOrders = append(createdOrders, order)
	}

	// 清空购物车（如果是从购物车提交的）
	database.DB.Where("user_id = ?", user.ID).Delete(&models.CartItem{})

	response := models.CreateOrderResponse{
		Orders: createdOrders,
	}

	responseOK(c, "订单提交成功", response)
}

// 获取订单列表
func getOrders(c *gin.Context) {
	var req models.OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := database.DB.Where("open_id = ?", userID).First(&user).Error; err != nil {
		responseError(c, 404, "用户不存在", err.Error())
		return
	}

	// 构建查询
	query := database.DB.Model(&models.Order{}).Where("user_id = ?", user.ID)

	// 筛选条件
	if req.Supplier != "" {
		query = query.Where("supplier = ?", req.Supplier)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 分页查询
	var orders []models.Order
	offset := (req.Page - 1) * req.Limit
	query.Offset(offset).Limit(req.Limit).Preload("Products").Find(&orders)

	// 统计供应商信息
	var suppliers []models.SupplierSummary
	database.DB.Raw(`
		SELECT supplier, 
			   COUNT(*) as order_count,
			   SUM(total_price) as total_price,
			   COUNT(DISTINCT product_id) as product_count
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = ?
		GROUP BY supplier
	`, user.ID).Scan(&suppliers)

	response := models.OrderListResponse{
		Orders:    orders,
		Suppliers: suppliers,
	}

	responseOK(c, "获取成功", response)
}

// 获取订单详情
func getOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	var order models.Order
	if err := database.DB.Preload("Products").First(&order, "id = ?", orderID).Error; err != nil {
		responseError(c, 404, "订单不存在", err.Error())
		return
	}

	responseOK(c, "获取成功", order)
}

// 更新订单状态
func updateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderId")

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, "id = ?", orderID).Error; err != nil {
		responseError(c, 404, "订单不存在", err.Error())
		return
	}

	order.Status = req.Status
	if req.Notes != "" {
		order.Notes = req.Notes
	}
	order.UpdatedAt = time.Now()

	if err := database.DB.Save(&order).Error; err != nil {
		responseError(c, 500, "更新失败", err.Error())
		return
	}

	responseOK(c, "更新成功", nil)
}

// 更新订单最终价格
func updateOrderFinalPrice(c *gin.Context) {
	orderID := c.Param("orderId")

	var req models.UpdateOrderPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, "id = ?", orderID).Error; err != nil {
		responseError(c, 404, "订单不存在", err.Error())
		return
	}

	order.FinalPrice = &req.FinalPrice
	order.UpdatedAt = time.Now()

	if err := database.DB.Save(&order).Error; err != nil {
		responseError(c, 500, "更新失败", err.Error())
		return
	}

	responseOK(c, "价格更新成功", nil)
}

// =================== 供应商管理 API ===================

// 获取供应商列表
func getSuppliers(c *gin.Context) {
	var suppliers []models.SupplierInfo

	// 查询供应商及统计信息
	database.DB.Raw(`
		SELECT s.name, s.contact_person, s.phone, s.status,
			   COUNT(DISTINCT p.id) as product_count,
			   COUNT(DISTINCT o.id) as total_orders
		FROM suppliers s
		LEFT JOIN products p ON s.name = p.supplier
		LEFT JOIN orders o ON s.name = o.supplier
		GROUP BY s.name, s.contact_person, s.phone, s.status
	`).Scan(&suppliers)

	response := models.SupplierListResponse{
		Suppliers: suppliers,
	}

	responseOK(c, "获取成功", response)
}

// 获取供应商详情
func getSupplier(c *gin.Context) {
	supplierName := c.Param("supplierName")

	var supplier models.Supplier
	if err := database.DB.First(&supplier, "name = ?", supplierName).Error; err != nil {
		responseError(c, 404, "供应商不存在", err.Error())
		return
	}

	// 统计信息
	var stats models.SupplierStatistics
	database.DB.Raw(`
		SELECT COUNT(DISTINCT p.id) as product_count,
			   COUNT(DISTINCT o.id) as total_orders,
			   COALESCE(SUM(o.total_price), 0) as total_amount,
			   COALESCE(AVG(o.total_price), 0) as average_order_amount
		FROM suppliers s
		LEFT JOIN products p ON s.name = p.supplier
		LEFT JOIN orders o ON s.name = o.supplier
		WHERE s.name = ?
	`, supplierName).Scan(&stats)

	// 最近订单
	var recentOrders []models.Order
	database.DB.Where("supplier = ?", supplierName).
		Order("created_at DESC").
		Limit(5).
		Preload("Products").
		Find(&recentOrders)

	response := models.SupplierDetailResponse{
		Supplier:     supplier,
		Statistics:   stats,
		RecentOrders: recentOrders,
	}

	responseOK(c, "获取成功", response)
}

// 获取供应商的商品列表
func getSupplierProducts(c *gin.Context) {
	supplierName := c.Param("supplierName")

	var products []models.Product
	database.DB.Where("supplier = ?", supplierName).Find(&products)

	responseOK(c, "获取成功", products)
}

// 获取供应商的订单列表
func getSupplierOrders(c *gin.Context) {
	supplierName := c.Param("supplierName")

	var orders []models.Order
	database.DB.Where("supplier = ?", supplierName).
		Preload("Products").
		Order("created_at DESC").
		Find(&orders)

	responseOK(c, "获取成功", orders)
}

// =================== 数据同步 API ===================

// 批量导入商品
func importProducts(c *gin.Context) {
	var req models.ImportProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	var createdProducts []models.Product

	for _, importProduct := range req.Products {
		product := models.Product{
			Name:        importProduct.Name,
			Price:       importProduct.Price,
			Unit:        importProduct.Unit,
			Description: importProduct.Description,
			Supplier:    importProduct.Supplier,
			Status:      "available",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := database.DB.Create(&product).Error; err != nil {
			responseError(c, 500, "导入失败", err.Error())
			return
		}

		createdProducts = append(createdProducts, product)
	}

	responseOK(c, fmt.Sprintf("成功导入 %d 个商品", len(createdProducts)), createdProducts)
}

// 导出订单数据（简化版本，实际应该生成Excel文件）
func exportOrders(c *gin.Context) {
	var req models.ExportOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, 400, "请求参数错误", err.Error())
		return
	}

	query := database.DB.Model(&models.Order{})

	// 筛选条件
	if req.DateFrom != "" {
		query = query.Where("created_at >= ?", req.DateFrom)
	}
	if req.DateTo != "" {
		query = query.Where("created_at <= ?", req.DateTo)
	}
	if req.Supplier != "" {
		query = query.Where("supplier = ?", req.Supplier)
	}

	var orders []models.Order
	query.Preload("Products").Find(&orders)

	responseOK(c, "导出成功", orders)
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
	v1 := r.Group("/v1")
	{
		// 商品管理 API
		v1.GET("/products", getProducts)
		v1.GET("/products/:productId", getProduct)
		v1.POST("/products/import", importProducts)

		// 购物车管理 API
		v1.GET("/cart", getCart)
		v1.POST("/cart/items", addToCart)
		v1.PUT("/cart/items/:itemId", updateCartItem)
		v1.DELETE("/cart/items/:itemId", deleteCartItem)
		v1.DELETE("/cart", clearCart)

		// 订单管理 API
		v1.POST("/orders", createOrder)
		v1.GET("/orders", getOrders)
		v1.GET("/orders/:orderId", getOrder)
		v1.PUT("/orders/:orderId/status", updateOrderStatus)
		v1.PUT("/orders/:orderId/final-price", updateOrderFinalPrice)
		v1.GET("/orders/export", exportOrders)

		// 供应商管理 API
		v1.GET("/suppliers", getSuppliers)
		v1.GET("/suppliers/:supplierName", getSupplier)
		v1.GET("/suppliers/:supplierName/products", getSupplierProducts)
		v1.GET("/suppliers/:supplierName/orders", getSupplierOrders)

		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			responseOK(c, "服务正常", "OK")
		})

		// 重置数据（仅用于开发测试）
		v1.POST("/reset-data", func(c *gin.Context) {
			database.ResetData()
			responseOK(c, "数据重置完成", "重新初始化了采购订单系统测试数据")
		})
	}

	fmt.Println("🚀 采购订单系统后端启动成功!")
	fmt.Println("💾 使用SQLite数据库存储")
	fmt.Println("🔗 后端地址: http://localhost:8080")
	fmt.Println("🔍 健康检查: http://localhost:8080/v1/health")
	fmt.Println("📖 API文档: 根据 docs/API_接口文档.md")

	// 启动服务器
	r.Run(":8080")
}
