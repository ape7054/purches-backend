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
		&models.Shop{},
		&models.Product{},
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
	var shopCount int64
	DB.Model(&models.Shop{}).Count(&shopCount)
	if shopCount > 0 {
		fmt.Printf("Data already exists (%d shops), skipping seed...\n", shopCount)
		return
	}

	// 插入商店数据
	shops := []models.Shop{
		{ID: "shop_1", Name: "快驴", Logo: ""},
		{ID: "shop_2", Name: "华兴街14号", Logo: ""},
		{ID: "shop_3", Name: "肖红梅", Logo: ""},
		{ID: "shop_4", Name: "文武双全大母指", Logo: ""},
		{ID: "shop_5", Name: "A12号豆腐档", Logo: ""},
		{ID: "shop_6", Name: "F35", Logo: ""},
		{ID: "shop_7", Name: "D30", Logo: ""},
		{ID: "shop_8", Name: "D129", Logo: ""},
		{ID: "shop_9", Name: "淡水锦龙冻品", Logo: ""},
		{ID: "shop_10", Name: "易订货", Logo: ""},
	}

	for _, shop := range shops {
		DB.Create(&shop)
	}

	// 插入部分商品数据（快驴）
	products := []models.Product{
		{ID: "prod_1", Name: "盐", Price: 5.00, ShopID: "shop_1"},
		{ID: "prod_2", Name: "味精", Price: 8.00, ShopID: "shop_1"},
		{ID: "prod_3", Name: "鸡精", Price: 12.00, ShopID: "shop_1"},
		{ID: "prod_4", Name: "生抽", Price: 15.00, ShopID: "shop_1"},
		{ID: "prod_5", Name: "老抽", Price: 16.00, ShopID: "shop_1"},
		{ID: "prod_6", Name: "大豆油", Price: 25.00, ShopID: "shop_2"},
		{ID: "prod_7", Name: "大米", Price: 30.00, ShopID: "shop_2"},
		{ID: "prod_8", Name: "鸡蛋", Price: 18.00, ShopID: "shop_2"},
		{ID: "prod_9", Name: "芹菜", Price: 6.00, ShopID: "shop_3"},
		{ID: "prod_10", Name: "腐竹", Price: 22.00, ShopID: "shop_3"},
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
	DB.Exec("DELETE FROM shops")
	DB.Exec("DELETE FROM products")
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
	// 插入商店数据
	shops := []models.Shop{
		{ID: "shop_1", Name: "快驴", Logo: ""},
		{ID: "shop_2", Name: "华兴街14号", Logo: ""},
		{ID: "shop_3", Name: "肖红梅", Logo: ""},
		{ID: "shop_4", Name: "文武双全大母指", Logo: ""},
		{ID: "shop_5", Name: "A12号豆腐档", Logo: ""},
		{ID: "shop_6", Name: "F35", Logo: ""},
		{ID: "shop_7", Name: "D30", Logo: ""},
		{ID: "shop_8", Name: "D129", Logo: ""},
		{ID: "shop_9", Name: "淡水锦龙冻品", Logo: ""},
		{ID: "shop_10", Name: "易订货", Logo: ""},
	}

	for _, shop := range shops {
		DB.Create(&shop)
	}

	// 插入部分商品数据（快驴）
	products := []models.Product{
		{ID: "prod_1", Name: "盐", Price: 5.00, ShopID: "shop_1"},
		{ID: "prod_2", Name: "味精", Price: 8.00, ShopID: "shop_1"},
		{ID: "prod_3", Name: "鸡精", Price: 12.00, ShopID: "shop_1"},
		{ID: "prod_4", Name: "生抽", Price: 15.00, ShopID: "shop_1"},
		{ID: "prod_5", Name: "老抽", Price: 16.00, ShopID: "shop_1"},
		{ID: "prod_6", Name: "大豆油", Price: 25.00, ShopID: "shop_2"},
		{ID: "prod_7", Name: "大米", Price: 30.00, ShopID: "shop_2"},
		{ID: "prod_8", Name: "鸡蛋", Price: 18.00, ShopID: "shop_2"},
		{ID: "prod_9", Name: "芹菜", Price: 6.00, ShopID: "shop_3"},
		{ID: "prod_10", Name: "腐竹", Price: 22.00, ShopID: "shop_3"},
	}

	for _, product := range products {
		DB.Create(&product)
	}

	fmt.Println("强制插入测试数据完成!")
}
