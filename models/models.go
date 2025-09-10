package models

import (
	"time"
)

// Product 商品模型
type Product struct {
	ID          int       `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" gorm:"not null"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Unit        string    `json:"unit" gorm:"not null"`
	Description string    `json:"description"`
	Supplier    string    `json:"supplier" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:available"` // available, unavailable, discontinued
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Supplier 供应商模型
type Supplier struct {
	Name          string    `json:"name" gorm:"primary_key"`
	ContactPerson string    `json:"contactPerson"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	Status        string    `json:"status" gorm:"default:'active'"` // active, inactive
	CreatedAt     time.Time `json:"createdAt"`
}

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	OpenID    string    `json:"openId" gorm:"unique;not null"`
	NickName  string    `json:"nickName"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CartItem 购物车商品模型
type CartItem struct {
	ID         int       `json:"id" gorm:"primary_key"`
	UserID     uint      `json:"userId" gorm:"not null"`
	ProductID  int       `json:"productId" gorm:"not null"`
	Name       string    `json:"name"`
	Count      int       `json:"count" gorm:"not null"`
	ShopName   string    `json:"shopName"` // 实际上是供应商名称
	Price      float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	TotalPrice float64   `json:"totalPrice" gorm:"type:decimal(10,2);not null"`
	User       User      `json:"user" gorm:"foreignkey:UserID"`
	Product    Product   `json:"product" gorm:"foreignkey:ProductID"`
	AddedAt    time.Time `json:"addedAt"`
}

// Order 订单模型
type Order struct {
	ID         string      `json:"id" gorm:"primary_key"`
	UserID     uint        `json:"userId" gorm:"not null"`
	Supplier   string      `json:"supplier" gorm:"not null"`
	TotalPrice float64     `json:"totalPrice" gorm:"type:decimal(10,2);not null"`
	FinalPrice *float64    `json:"finalPrice" gorm:"type:decimal(10,2)"` // 可调整的最终价格
	Status     string      `json:"status" gorm:"default:'pending'"`      // pending, confirmed, delivering, completed, cancelled
	Notes      string      `json:"notes"`
	User       User        `json:"user" gorm:"foreignkey:UserID"`
	Products   []OrderItem `json:"products" gorm:"foreignkey:OrderID"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

// OrderItem 订单商品模型
type OrderItem struct {
	ID          int       `json:"id" gorm:"primary_key"`
	OrderID     string    `json:"orderId" gorm:"not null"`
	ProductID   int       `json:"productId" gorm:"not null"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Count       int       `json:"count" gorm:"not null"`
	Unit        string    `json:"unit"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	TotalPrice  float64   `json:"totalPrice" gorm:"type:decimal(10,2);not null"`
	Order       Order     `json:"order" gorm:"foreignkey:OrderID"`
	Product     Product   `json:"product" gorm:"foreignkey:ProductID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// 请求/响应模型

// ProductListRequest 商品列表请求
type ProductListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Supplier string `form:"supplier"`
	Search   string `form:"search"`
	Status   string `form:"status"`
}

// ProductListResponse 商品列表响应
type ProductListResponse struct {
	Products   []Product          `json:"products"`
	Pagination PaginationResponse `json:"pagination"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
}

// CartResponse 购物车响应
type CartResponse struct {
	Items   []CartItem  `json:"items"`
	Summary CartSummary `json:"summary"`
}

// CartSummary 购物车汇总
type CartSummary struct {
	TotalItems    int     `json:"totalItems"`
	TotalPrice    float64 `json:"totalPrice"`
	SupplierCount int     `json:"supplierCount"`
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	ProductID int `json:"productId" binding:"required"`
	Count     int `json:"count" binding:"required,min=1"`
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	Count int `json:"count" binding:"required,min=0"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
	Notes string             `json:"notes"`
}

// OrderItemRequest 订单商品请求
type OrderItemRequest struct {
	ProductID int `json:"productId" binding:"required"`
	Count     int `json:"count" binding:"required,min=1"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	Orders []Order `json:"orders"`
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Supplier string `form:"supplier"`
	Status   string `form:"status"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Orders    []Order           `json:"orders"`
	Suppliers []SupplierSummary `json:"suppliers"`
}

// SupplierSummary 供应商汇总
type SupplierSummary struct {
	Supplier     string  `json:"supplier"`
	OrderCount   int     `json:"orderCount"`
	TotalPrice   float64 `json:"totalPrice"`
	ProductCount int     `json:"productCount"`
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

// UpdateOrderPriceRequest 更新订单价格请求
type UpdateOrderPriceRequest struct {
	FinalPrice float64 `json:"finalPrice" binding:"required,min=0"`
	Reason     string  `json:"reason"`
}

// SupplierListResponse 供应商列表响应
type SupplierListResponse struct {
	Suppliers []SupplierInfo `json:"suppliers"`
}

// SupplierInfo 供应商信息
type SupplierInfo struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contactPerson"`
	Phone         string `json:"phone"`
	ProductCount  int    `json:"productCount"`
	TotalOrders   int    `json:"totalOrders"`
	Status        string `json:"status"`
}

// SupplierDetailResponse 供应商详情响应
type SupplierDetailResponse struct {
	Supplier     Supplier           `json:"supplier"`
	Statistics   SupplierStatistics `json:"statistics"`
	RecentOrders []Order            `json:"recentOrders"`
}

// SupplierStatistics 供应商统计
type SupplierStatistics struct {
	ProductCount       int     `json:"productCount"`
	TotalOrders        int     `json:"totalOrders"`
	TotalAmount        float64 `json:"totalAmount"`
	AverageOrderAmount float64 `json:"averageOrderAmount"`
}

// ImportProductsRequest 批量导入商品请求
type ImportProductsRequest struct {
	Products []ImportProduct `json:"products" binding:"required"`
}

// ImportProduct 导入商品
type ImportProduct struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Unit        string  `json:"unit" binding:"required"`
	Description string  `json:"description"`
	Supplier    string  `json:"supplier" binding:"required"`
}

// ExportOrdersRequest 导出订单请求
type ExportOrdersRequest struct {
	Format   string `form:"format" binding:"required,oneof=excel csv"`
	DateFrom string `form:"dateFrom"`
	DateTo   string `form:"dateTo"`
	Supplier string `form:"supplier"`
}
