package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// ImportProductFromJSON 从JSON导入商品的结构
type ImportProductFromJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
	Supplier    string  `json:"supplier"`
}

// LoadProductsFromJSON 从JSON文件加载商品数据
func LoadProductsFromJSON(filename string) ([]ImportProductFromJSON, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var products []ImportProductFromJSON
	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
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

	// 尝试从JSON文件导入真实数据
	fmt.Println("尝试从 docs/products.json 导入真实商品数据...")
	jsonProducts, err := LoadProductsFromJSON("docs/products.json")
	if err != nil {
		fmt.Printf("无法读取JSON文件，使用默认测试数据: %v\n", err)
		SeedDefaultData()
		return
	}

	// 先创建供应商
	fmt.Println("创建供应商数据...")
	CreateSuppliersFromProducts(jsonProducts)

	// 创建商品
	fmt.Println("导入商品数据...")
	for _, jsonProduct := range jsonProducts {
		product := models.Product{
			Name:        jsonProduct.Name,
			Price:       jsonProduct.Price,
			Unit:        jsonProduct.Unit,
			Description: jsonProduct.Description,
			Supplier:    jsonProduct.Supplier,
			Status:      "available",
		}

		if err := DB.Create(&product).Error; err != nil {
			fmt.Printf("导入商品失败 %s: %v\n", jsonProduct.Name, err)
		}
	}

	fmt.Printf("成功导入 %d 种商品数据!\n", len(jsonProducts))
}

// CreateSuppliersFromProducts 从商品数据中提取并创建供应商
func CreateSuppliersFromProducts(products []ImportProductFromJSON) {
	// 提取唯一供应商名称
	supplierMap := make(map[string]bool)
	for _, product := range products {
		supplierMap[product.Supplier] = true
	}

	// 供应商联系信息映射（根据实际情况填写）
	supplierInfo := map[string]models.Supplier{
		"F35": {
			Name:          "F35",
			ContactPerson: "张先生",
			Phone:         "13800138000",
			Address:       "批发市场F35号",
			Status:        "active",
		},
		"快驴": {
			Name:          "快驴",
			ContactPerson: "客服",
			Phone:         "400-123-4567",
			Address:       "线上平台",
			Status:        "active",
		},
		"华兴街14号": {
			Name:          "华兴街14号",
			ContactPerson: "李女士",
			Phone:         "13800138001",
			Address:       "华兴街14号",
			Status:        "active",
		},
		"肖红梅": {
			Name:          "肖红梅",
			ContactPerson: "肖红梅",
			Phone:         "13800138002",
			Address:       "蔬菜批发区",
			Status:        "active",
		},
		"艾武双全大母指": {
			Name:          "艾武双全大母指",
			ContactPerson: "王师傅",
			Phone:         "13800138003",
			Address:       "批发市场大母指档口",
			Status:        "active",
		},
		"A12号豆腐档": {
			Name:          "A12号豆腐档",
			ContactPerson: "陈师傅",
			Phone:         "13800138004",
			Address:       "批发市场A12号",
			Status:        "active",
		},
		"D30": {
			Name:          "D30",
			ContactPerson: "赵女士",
			Phone:         "13800138005",
			Address:       "批发市场D30号",
			Status:        "active",
		},
		"D129": {
			Name:          "D129",
			ContactPerson: "钱师傅",
			Phone:         "13800138006",
			Address:       "批发市场D129号",
			Status:        "active",
		},
	}

	// 创建供应商记录
	for supplierName := range supplierMap {
		supplier, exists := supplierInfo[supplierName]
		if !exists {
			// 如果没有预设信息，创建默认信息
			supplier = models.Supplier{
				Name:          supplierName,
				ContactPerson: "待完善",
				Phone:         "待完善",
				Address:       "待完善",
				Status:        "active",
			}
		}

		if err := DB.Create(&supplier).Error; err != nil {
			fmt.Printf("创建供应商失败 %s: %v\n", supplierName, err)
		}
	}

	fmt.Printf("成功创建 %d 个供应商\n", len(supplierMap))
}

// SeedDefaultData 使用默认测试数据（当JSON文件不可用时）
func SeedDefaultData() {
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
	}

	for _, product := range products {
		DB.Create(&product)
	}

	fmt.Println("Default test data inserted successfully!")
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

// SeedDataForce 强制插入数据（不检查是否已存在）
func SeedDataForce() {
	// 尝试从JSON文件导入
	jsonProducts, err := LoadProductsFromJSON("docs/products.json")
	if err != nil {
		fmt.Printf("无法读取JSON文件，使用默认测试数据: %v\n", err)
		SeedDefaultData()
		return
	}

	// 创建供应商和商品
	CreateSuppliersFromProducts(jsonProducts)

	for _, jsonProduct := range jsonProducts {
		product := models.Product{
			Name:        jsonProduct.Name,
			Price:       jsonProduct.Price,
			Unit:        jsonProduct.Unit,
			Description: jsonProduct.Description,
			Supplier:    jsonProduct.Supplier,
			Status:      "available",
		}
		DB.Create(&product)
	}

	fmt.Printf("强制插入 %d 种商品数据完成!\n", len(jsonProducts))
}
