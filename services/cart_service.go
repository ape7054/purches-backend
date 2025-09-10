package services

import (
	"purches-backend/models"
	"time"

	"gorm.io/gorm"
)

const DEFAULT_USER_ID = "user_1"

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{
		db: db,
	}
}

// GetCart 获取购物车
func (cs *CartService) GetCart() (*models.CartResponse, error) {
	userID := DEFAULT_USER_ID

	// 确保用户存在
	var user models.User
	if err := cs.db.FirstOrCreate(&user, models.User{OpenID: userID}).Error; err != nil {
		return nil, err
	}

	// 获取购物车商品
	var cartItems []models.CartItem
	cs.db.Where("user_id = ?", user.ID).Find(&cartItems)

	// 计算统计信息
	var totalPrice float64
	supplierMap := make(map[string]bool)

	for _, item := range cartItems {
		totalPrice += item.TotalPrice
		supplierMap[item.ShopName] = true
	}

	response := &models.CartResponse{
		Items: cartItems,
		Summary: models.CartSummary{
			TotalItems:    len(cartItems),
			TotalPrice:    totalPrice,
			SupplierCount: len(supplierMap),
		},
	}

	return response, nil
}

// AddToCart 添加商品到购物车
func (cs *CartService) AddToCart(req models.AddToCartRequest) (*models.CartItem, error) {
	// 查找商品信息
	var product models.Product
	if err := cs.db.First(&product, req.ProductID).Error; err != nil {
		return nil, err
	}

	userID := DEFAULT_USER_ID

	// 确保用户存在
	var user models.User
	if err := cs.db.FirstOrCreate(&user, models.User{OpenID: userID}).Error; err != nil {
		return nil, err
	}

	// 检查购物车中是否已存在该商品
	var existingItem models.CartItem
	result := cs.db.Where("user_id = ? AND product_id = ?", user.ID, req.ProductID).First(&existingItem)

	if result.Error == nil {
		// 如果存在，增加数量
		existingItem.Count += req.Count
		existingItem.TotalPrice = float64(existingItem.Count) * existingItem.Price
		if err := cs.db.Save(&existingItem).Error; err != nil {
			return nil, err
		}
		return &existingItem, nil
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
		if err := cs.db.Create(&newItem).Error; err != nil {
			return nil, err
		}
		return &newItem, nil
	}
}

// UpdateCartItem 更新购物车商品数量
func (cs *CartService) UpdateCartItem(itemID int, count int) error {
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := cs.db.Where("open_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// 查找购物车商品
	var cartItem models.CartItem
	if err := cs.db.Where("id = ? AND user_id = ?", itemID, user.ID).First(&cartItem).Error; err != nil {
		return err
	}

	if count > 0 {
		cartItem.Count = count
		cartItem.TotalPrice = float64(cartItem.Count) * cartItem.Price
		return cs.db.Save(&cartItem).Error
	} else {
		// 数量为0时删除商品
		return cs.db.Delete(&cartItem).Error
	}
}

// DeleteCartItem 删除购物车商品
func (cs *CartService) DeleteCartItem(itemID int) error {
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := cs.db.Where("open_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// 删除商品
	result := cs.db.Where("id = ? AND user_id = ?", itemID, user.ID).Delete(&models.CartItem{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

// ClearCart 清空购物车
func (cs *CartService) ClearCart() error {
	userID := DEFAULT_USER_ID

	// 获取用户信息
	var user models.User
	if err := cs.db.Where("open_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	return cs.db.Where("user_id = ?", user.ID).Delete(&models.CartItem{}).Error
}
