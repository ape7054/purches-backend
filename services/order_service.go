package services

import (
	"fmt"
	"purches-backend/config"
	"purches-backend/models"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

// getDefaultUserID 获取默认用户ID
func (os *OrderService) getDefaultUserID() string {
	return config.GetConfig().App.DefaultUser
}

// CreateOrder 创建订单（按供应商分组）
func (os *OrderService) CreateOrder(req models.CreateOrderRequest) ([]models.Order, error) {
	userID := os.getDefaultUserID()

	// 获取用户信息
	var user models.User
	if err := os.db.Where("open_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 获取商品信息并按供应商分组
	supplierGroups := make(map[string][]models.OrderItemRequest)
	productMap := make(map[int]models.Product)

	for _, item := range req.Items {
		var product models.Product
		if err := os.db.First(&product, item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("商品ID %d 不存在", item.ProductID)
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

		if err := os.db.Create(&order).Error; err != nil {
			return nil, err
		}

		// 创建订单商品明细
		for _, orderItem := range orderItems {
			if err := os.db.Create(&orderItem).Error; err != nil {
				return nil, err
			}
		}

		createdOrders = append(createdOrders, order)
	}

	// 清空购物车（如果是从购物车提交的）
	os.db.Where("user_id = ?", user.ID).Delete(&models.CartItem{})

	return createdOrders, nil
}

// GetOrders 获取订单列表
func (os *OrderService) GetOrders(req models.OrderListRequest) (*models.OrderListResponse, error) {
	userID := os.getDefaultUserID()

	// 获取用户信息
	var user models.User
	if err := os.db.Where("open_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 构建查询
	query := os.db.Model(&models.Order{}).Where("user_id = ?", user.ID)

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
	if err := query.Offset(offset).Limit(req.Limit).Preload("Products").Find(&orders).Error; err != nil {
		return nil, err
	}

	// 统计供应商信息
	var suppliers []models.SupplierSummary
	if err := os.db.Raw(`
		SELECT supplier, 
			   COUNT(*) as order_count,
			   SUM(total_price) as total_price,
			   COUNT(DISTINCT product_id) as product_count
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = ?
		GROUP BY supplier
	`, user.ID).Scan(&suppliers).Error; err != nil {
		return nil, err
	}

	response := &models.OrderListResponse{
		Orders:    orders,
		Suppliers: suppliers,
	}

	return response, nil
}

// GetOrderByID 根据ID获取订单
func (os *OrderService) GetOrderByID(orderID string) (*models.Order, error) {
	var order models.Order
	if err := os.db.Preload("Products").First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (os *OrderService) UpdateOrderStatus(orderID string, req models.UpdateOrderStatusRequest) error {
	var order models.Order
	if err := os.db.First(&order, "id = ?", orderID).Error; err != nil {
		return err
	}

	order.Status = req.Status
	if req.Notes != "" {
		order.Notes = req.Notes
	}
	order.UpdatedAt = time.Now()

	return os.db.Save(&order).Error
}

// UpdateOrderFinalPrice 更新订单最终价格
func (os *OrderService) UpdateOrderFinalPrice(orderID string, finalPrice float64) error {
	var order models.Order
	if err := os.db.First(&order, "id = ?", orderID).Error; err != nil {
		return err
	}

	order.FinalPrice = &finalPrice
	order.UpdatedAt = time.Now()

	return os.db.Save(&order).Error
}

// ExportOrders 导出订单数据
func (os *OrderService) ExportOrders(req models.ExportOrdersRequest) ([]models.Order, error) {
	query := os.db.Model(&models.Order{})

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
	if err := query.Preload("Products").Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
