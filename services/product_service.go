package services

import (
	"purches-backend/database"
	"purches-backend/models"
	"time"

	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

// GetProducts 获取商品列表
func (ps *ProductService) GetProducts(req models.ProductListRequest) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 构建查询
	query := ps.db.Model(&models.Product{})

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
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductByID 根据ID获取商品
func (ps *ProductService) GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	if err := ps.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// ImportProducts 批量导入商品
func (ps *ProductService) ImportProducts(importProducts []models.ImportProduct) ([]models.Product, error) {
	var createdProducts []models.Product

	for _, importProduct := range importProducts {
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

		if err := ps.db.Create(&product).Error; err != nil {
			return nil, err
		}

		createdProducts = append(createdProducts, product)
	}

	return createdProducts, nil
}

// ImportProductsFromJSON 从JSON文件导入商品
func (ps *ProductService) ImportProductsFromJSON(filename string) (int, error) {
	// 读取JSON文件
	jsonProducts, err := database.LoadProductsFromJSON(filename)
	if err != nil {
		return 0, err
	}

	// 清空现有数据
	ps.db.Exec("DELETE FROM products")
	ps.db.Exec("DELETE FROM suppliers")

	// 创建供应商
	database.CreateSuppliersFromProducts(jsonProducts)

	// 导入商品
	var createdCount int
	for _, jsonProduct := range jsonProducts {
		product := models.Product{
			Name:        jsonProduct.Name,
			Price:       jsonProduct.Price,
			Unit:        jsonProduct.Unit,
			Description: jsonProduct.Description,
			Supplier:    jsonProduct.Supplier,
			Status:      "available",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := ps.db.Create(&product).Error; err == nil {
			createdCount++
		}
	}

	return createdCount, nil
}
