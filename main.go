package main

import (
	"fmt"
	"purches-backend/config"
	"purches-backend/controllers"
	"purches-backend/database"
	"purches-backend/middleware"
	"purches-backend/services"
	"purches-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("配置加载失败: %v", err))
	}

	// 设置Gin模式
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	database.InitDatabase()

	// 创建Gin实例
	r := gin.Default()

	// 设置中间件
	r.Use(middleware.CORS())

	// 初始化服务层
	productService := services.NewProductService(database.DB)
	cartService := services.NewCartService(database.DB)
	orderService := services.NewOrderService(database.DB)
	supplierService := services.NewSupplierService(database.DB)

	// 初始化控制器层
	productController := controllers.NewProductController(productService)
	cartController := controllers.NewCartController(cartService)
	orderController := controllers.NewOrderController(orderService)
	supplierController := controllers.NewSupplierController(supplierService)

	// 设置路由
	setupRoutes(r, productController, cartController, orderController, supplierController, productService)

	// 启动信息
	fmt.Printf("🚀 %s 启动成功!\n", cfg.App.Name)
	fmt.Printf("📦 版本: %s\n", cfg.App.Version)
	fmt.Printf("🔧 环境: %s\n", cfg.App.Environment)
	fmt.Printf("💾 数据库: %s\n", cfg.Database.Type)
	fmt.Printf("🔗 后端地址: http://localhost:%s\n", cfg.Server.Port)
	fmt.Printf("🔍 健康检查: http://localhost:%s/v1/health\n", cfg.Server.Port)
	fmt.Println("📖 API文档: 根据 docs/API_接口文档.md")

	// 启动服务器
	r.Run(":" + cfg.Server.Port)
}

// setupRoutes 设置路由
func setupRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
	cartController *controllers.CartController,
	orderController *controllers.OrderController,
	supplierController *controllers.SupplierController,
	productService *services.ProductService,
) {
	// API路由组
	v1 := r.Group("/v1")
	{
		// 商品管理 API
		v1.GET("/products", productController.GetProducts)
		v1.GET("/products/:productId", productController.GetProduct)
		v1.POST("/products/import", productController.ImportProducts)

		// 购物车管理 API
		v1.GET("/cart", cartController.GetCart)
		v1.POST("/cart/items", cartController.AddToCart)
		v1.PUT("/cart/items/:itemId", cartController.UpdateCartItem)
		v1.DELETE("/cart/items/:itemId", cartController.DeleteCartItem)
		v1.DELETE("/cart", cartController.ClearCart)

		// 订单管理 API
		v1.POST("/orders", orderController.CreateOrder)
		v1.GET("/orders", orderController.GetOrders)
		v1.GET("/orders/:orderId", orderController.GetOrder)
		v1.PUT("/orders/:orderId/status", orderController.UpdateOrderStatus)
		v1.PUT("/orders/:orderId/final-price", orderController.UpdateOrderFinalPrice)
		v1.GET("/orders/export", orderController.ExportOrders)

		// 供应商管理 API
		v1.GET("/suppliers", supplierController.GetSuppliers)
		v1.GET("/suppliers/:supplierName", supplierController.GetSupplier)
		v1.GET("/suppliers/:supplierName/products", supplierController.GetSupplierProducts)
		v1.GET("/suppliers/:supplierName/orders", supplierController.GetSupplierOrders)

		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			utils.ResponseOK(c, "服务正常", "OK")
		})

		// 开发工具接口
		setupDevRoutes(v1, productService)
	}
}

// setupDevRoutes 设置开发工具路由
func setupDevRoutes(v1 *gin.RouterGroup, productService *services.ProductService) {
	// 重置数据（仅用于开发测试）
	v1.POST("/reset-data", func(c *gin.Context) {
		database.ResetData()
		utils.ResponseOK(c, "数据重置完成", "重新初始化了采购订单系统测试数据")
	})

	// 导入JSON数据（仅用于开发测试）
	v1.POST("/import-products-json", func(c *gin.Context) {
		count, err := productService.ImportProductsFromJSON("docs/products.json")
		if err != nil {
			utils.ResponseError(c, 500, "无法读取JSON文件", err.Error())
			return
		}

		utils.ResponseOK(c, "JSON数据导入完成", gin.H{
			"count":   count,
			"message": fmt.Sprintf("成功导入 %d 种商品", count),
		})
	})
}
