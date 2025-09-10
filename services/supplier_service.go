package services

import (
	"purches-backend/models"

	"gorm.io/gorm"
)

type SupplierService struct {
	db *gorm.DB
}

func NewSupplierService(db *gorm.DB) *SupplierService {
	return &SupplierService{
		db: db,
	}
}

// GetSuppliers 获取供应商列表
func (ss *SupplierService) GetSuppliers() (*models.SupplierListResponse, error) {
	var suppliers []models.SupplierInfo

	// 查询供应商及统计信息
	err := ss.db.Raw(`
		SELECT s.name, s.contact_person, s.phone, s.status,
			   COUNT(DISTINCT p.id) as product_count,
			   COUNT(DISTINCT o.id) as total_orders
		FROM suppliers s
		LEFT JOIN products p ON s.name = p.supplier
		LEFT JOIN orders o ON s.name = o.supplier
		GROUP BY s.name, s.contact_person, s.phone, s.status
	`).Scan(&suppliers).Error

	if err != nil {
		return nil, err
	}

	response := &models.SupplierListResponse{
		Suppliers: suppliers,
	}

	return response, nil
}

// GetSupplierDetail 获取供应商详情
func (ss *SupplierService) GetSupplierDetail(supplierName string) (*models.SupplierDetailResponse, error) {
	var supplier models.Supplier
	if err := ss.db.First(&supplier, "name = ?", supplierName).Error; err != nil {
		return nil, err
	}

	// 统计信息
	var stats models.SupplierStatistics
	err := ss.db.Raw(`
		SELECT COUNT(DISTINCT p.id) as product_count,
			   COUNT(DISTINCT o.id) as total_orders,
			   COALESCE(SUM(o.total_price), 0) as total_amount,
			   COALESCE(AVG(o.total_price), 0) as average_order_amount
		FROM suppliers s
		LEFT JOIN products p ON s.name = p.supplier
		LEFT JOIN orders o ON s.name = o.supplier
		WHERE s.name = ?
	`, supplierName).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// 最近订单
	var recentOrders []models.Order
	ss.db.Where("supplier = ?", supplierName).
		Order("created_at DESC").
		Limit(5).
		Preload("Products").
		Find(&recentOrders)

	response := &models.SupplierDetailResponse{
		Supplier:     supplier,
		Statistics:   stats,
		RecentOrders: recentOrders,
	}

	return response, nil
}

// GetSupplierProducts 获取供应商的商品列表
func (ss *SupplierService) GetSupplierProducts(supplierName string) ([]models.Product, error) {
	var products []models.Product
	err := ss.db.Where("supplier = ?", supplierName).Find(&products).Error
	return products, err
}

// GetSupplierOrders 获取供应商的订单列表
func (ss *SupplierService) GetSupplierOrders(supplierName string) ([]models.Order, error) {
	var orders []models.Order
	err := ss.db.Where("supplier = ?", supplierName).
		Preload("Products").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
