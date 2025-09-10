package database

import (
	"fmt"
	"log"
	"purches-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error

	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open("purches.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully!")

	// 自动迁移数据表
	err = DB.AutoMigrate(
		&models.Product{},
		&models.Supplier{},
		&models.User{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migration completed!")

	// 初始化测试数据
	SeedData()
}

// SeedData 初始化测试数据
func SeedData() {
	// 检查是否已经有数据，避免重复插入
	var productCount int64
	DB.Model(&models.Product{}).Count(&productCount)
	if productCount > 0 {
		fmt.Printf("Data already exists (%d products), skipping seed...\n", productCount)
		return
	}

	// 插入供应商数据
	suppliers := []models.Supplier{
		{Name: "F35", ContactPerson: "张先生", Phone: "13800138000", Address: "批发市场F35号", Status: "active"},
		{Name: "D30", ContactPerson: "李女士", Phone: "13800138001", Address: "批发市场D30号", Status: "active"},
		{Name: "D129", ContactPerson: "王师傅", Phone: "13800138002", Address: "批发市场D129号", Status: "active"},
		{Name: "快驴", ContactPerson: "客服", Phone: "400-123-4567", Address: "线上平台", Status: "active"},
		{Name: "A12号豆腐档", ContactPerson: "陈师傅", Phone: "13800138003", Address: "批发市场A12号", Status: "active"},
	}

	for _, supplier := range suppliers {
		DB.Create(&supplier)
	}

	// 插入商品数据（符合API文档示例）
	products := []models.Product{
		{Name: "牛蛙", Price: 31.00, Unit: "斤", Description: "牛蛙杀好处理干净去掉内脏和眼睛，去掉爪子，50斤", Supplier: "F35", Status: "available"},
		{Name: "黑鱼片", Price: 28.50, Unit: "斤", Description: "新鲜黑鱼片，无刺", Supplier: "F35", Status: "available"},
		{Name: "基围虾", Price: 45.00, Unit: "斤", Description: "活基围虾，规格40-50只/斤", Supplier: "F35", Status: "available"},
		{Name: "白萝卜", Price: 2.50, Unit: "斤", Description: "新鲜白萝卜", Supplier: "D30", Status: "available"},
		{Name: "胡萝卜", Price: 3.00, Unit: "斤", Description: "新鲜胡萝卜", Supplier: "D30", Status: "available"},
		{Name: "土豆", Price: 2.80, Unit: "斤", Description: "新鲜土豆", Supplier: "D30", Status: "available"},
		{Name: "鸡蛋", Price: 6.50, Unit: "斤", Description: "新鲜鸡蛋，散装", Supplier: "D129", Status: "available"},
		{Name: "鸭蛋", Price: 8.00, Unit: "斤", Description: "新鲜鸭蛋", Supplier: "D129", Status: "available"},
		{Name: "大米", Price: 5.20, Unit: "斤", Description: "优质大米，5斤装", Supplier: "快驴", Status: "available"},
		{Name: "面粉", Price: 4.80, Unit: "斤", Description: "高筋面粉，适合做面条", Supplier: "快驴", Status: "available"},
		{Name: "嫩豆腐", Price: 3.50, Unit: "块", Description: "新鲜嫩豆腐", Supplier: "A12号豆腐档", Status: "available"},
		{Name: "老豆腐", Price: 3.00, Unit: "块", Description: "结实老豆腐，适合炒菜", Supplier: "A12号豆腐档", Status: "available"},
		{Name: "豆腐皮", Price: 8.00, Unit: "斤", Description: "手工豆腐皮", Supplier: "A12号豆腐档", Status: "available"},
		{Name: "青椒", Price: 4.50, Unit: "斤", Description: "新鲜青椒", Supplier: "D30", Status: "available"},
		{Name: "西红柿", Price: 5.00, Unit: "斤", Description: "新鲜西红柿", Supplier: "D30", Status: "available"},
	}

	for _, product := range products {
		DB.Create(&product)
	}

	fmt.Println("Test data inserted successfully!")
}

// ResetData 重新初始化数据（清空并重新插入）
func ResetData() {
	fmt.Println("开始重置数据...")

	// 删除所有现有数据
	DB.Exec("DELETE FROM products")
	DB.Exec("DELETE FROM suppliers")
	DB.Exec("DELETE FROM users")
	DB.Exec("DELETE FROM cart_items")
	DB.Exec("DELETE FROM orders")
	DB.Exec("DELETE FROM order_items")

	fmt.Println("已清空所有数据")

	// 重新插入数据
	SeedDataForce()
}

// SeedDataForce 强制插入测试数据（不检查是否已存在）
func SeedDataForce() {
	// 插入供应商数据
	suppliers := []models.Supplier{
		{Name: "F35", ContactPerson: "张先生", Phone: "13800138000", Address: "批发市场F35号", Status: "active"},
		{Name: "D30", ContactPerson: "李女士", Phone: "13800138001", Address: "批发市场D30号", Status: "active"},
		{Name: "D129", ContactPerson: "王师傅", Phone: "13800138002", Address: "批发市场D129号", Status: "active"},
		{Name: "快驴", ContactPerson: "客服", Phone: "400-123-4567", Address: "线上平台", Status: "active"},
		{Name: "A12号豆腐档", ContactPerson: "陈师傅", Phone: "13800138003", Address: "批发市场A12号", Status: "active"},
	}

	for _, supplier := range suppliers {
		DB.Create(&supplier)
	}

	// 插入商品数据
	products := []models.Product{
		{Name: "牛蛙", Price: 31.00, Unit: "斤", Description: "牛蛙杀好处理干净去掉内脏和眼睛，去掉爪子，50斤", Supplier: "F35", Status: "available"},
		{Name: "黑鱼片", Price: 28.50, Unit: "斤", Description: "新鲜黑鱼片，无刺", Supplier: "F35", Status: "available"},
		{Name: "基围虾", Price: 45.00, Unit: "斤", Description: "活基围虾，规格40-50只/斤", Supplier: "F35", Status: "available"},
		{Name: "白萝卜", Price: 2.50, Unit: "斤", Description: "新鲜白萝卜", Supplier: "D30", Status: "available"},
		{Name: "胡萝卜", Price: 3.00, Unit: "斤", Description: "新鲜胡萝卜", Supplier: "D30", Status: "available"},
		{Name: "土豆", Price: 2.80, Unit: "斤", Description: "新鲜土豆", Supplier: "D30", Status: "available"},
		{Name: "鸡蛋", Price: 6.50, Unit: "斤", Description: "新鲜鸡蛋，散装", Supplier: "D129", Status: "available"},
		{Name: "鸭蛋", Price: 8.00, Unit: "斤", Description: "新鲜鸭蛋", Supplier: "D129", Status: "available"},
		{Name: "大米", Price: 5.20, Unit: "斤", Description: "优质大米，5斤装", Supplier: "快驴", Status: "available"},
		{Name: "面粉", Price: 4.80, Unit: "斤", Description: "高筋面粉，适合做面条", Supplier: "快驴", Status: "available"},
		{Name: "嫩豆腐", Price: 3.50, Unit: "块", Description: "新鲜嫩豆腐", Supplier: "A12号豆腐档", Status: "available"},
		{Name: "老豆腐", Price: 3.00, Unit: "块", Description: "结实老豆腐，适合炒菜", Supplier: "A12号豆腐档", Status: "available"},
		{Name: "豆腐皮", Price: 8.00, Unit: "斤", Description: "手工豆腐皮", Supplier: "A12号豆腐档", Status: "available"},
		{Name: "青椒", Price: 4.50, Unit: "斤", Description: "新鲜青椒", Supplier: "D30", Status: "available"},
		{Name: "西红柿", Price: 5.00, Unit: "斤", Description: "新鲜西红柿", Supplier: "D30", Status: "available"},
	}

	for _, product := range products {
		DB.Create(&product)
	}

	fmt.Println("强制插入测试数据完成!")
}
