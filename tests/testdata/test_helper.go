package testdata

import (
	"purches-backend/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB 创建测试数据库
func SetupTestDB() (*gorm.DB, error) {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，减少测试输出
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构
	err = db.AutoMigrate(
		&models.Product{},
		&models.Supplier{},
		&models.User{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SeedTestData 创建测试数据
func SeedTestData(db *gorm.DB) error {
	// 创建测试供应商
	suppliers := []models.Supplier{
		{Name: "测试供应商A", ContactPerson: "张三", Phone: "13800000001", Address: "测试地址A", Status: "active", CreatedAt: time.Now()},
		{Name: "测试供应商B", ContactPerson: "李四", Phone: "13800000002", Address: "测试地址B", Status: "active", CreatedAt: time.Now()},
	}

	for _, supplier := range suppliers {
		if err := db.Create(&supplier).Error; err != nil {
			return err
		}
	}

	// 创建测试商品
	products := []models.Product{
		{Name: "测试商品1", Price: 10.50, Unit: "个", Description: "测试商品描述1", Supplier: "测试供应商A", Status: "available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "测试商品2", Price: 25.00, Unit: "斤", Description: "测试商品描述2", Supplier: "测试供应商A", Status: "available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "测试商品3", Price: 8.80, Unit: "包", Description: "测试商品描述3", Supplier: "测试供应商B", Status: "available", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	// 创建测试用户
	user := models.User{
		OpenID:    "test_user_001",
		NickName:  "测试用户",
		AvatarURL: "http://test.com/avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(db *gorm.DB) error {
	// 删除所有表的数据
	tables := []interface{}{
		&models.OrderItem{},
		&models.Order{},
		&models.CartItem{},
		&models.Product{},
		&models.Supplier{},
		&models.User{},
	}

	for _, table := range tables {
		if err := db.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetTestProductData 获取测试商品数据
func GetTestProductData() models.Product {
	return models.Product{
		Name:        "新测试商品",
		Price:       15.99,
		Unit:        "件",
		Description: "新测试商品描述",
		Supplier:    "测试供应商A",
		Status:      "available",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// GetTestUserID 获取测试用户ID
func GetTestUserID() uint {
	return 1 // 通常测试中第一个创建的用户ID是1
}
