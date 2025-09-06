package models

import (
	"time"
)

// Shop 商店模型
type Shop struct {
	ID        string    `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product 商品模型
type Product struct {
	ID        string    `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	ImageURL  string    `json:"imageUrl"`
	ShopID    string    `json:"shop_id" gorm:"not null"`
	Shop      Shop      `json:"shop" gorm:"foreignkey:ShopID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	OpenID    string    `json:"open_id" gorm:"unique;not null"` // 微信小程序的openid
	NickName  string    `json:"nick_name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartItem 购物车商品模型
type CartItem struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	ProductID   string    `json:"product_id" gorm:"not null"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	User        User      `json:"user" gorm:"foreignkey:UserID"`
	Product     Product   `json:"product" gorm:"foreignkey:ProductID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Order 订单模型
type Order struct {
	ID         string      `json:"id" gorm:"primary_key"`
	UserID     uint        `json:"user_id" gorm:"not null"`
	TotalPrice float64     `json:"total_price" gorm:"not null"`
	Status     string      `json:"status" gorm:"default:'pending'"` // pending, confirmed, cancelled
	Remark     string      `json:"remark"`
	User       User        `json:"user" gorm:"foreignkey:UserID"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignkey:OrderID"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderItem 订单商品模型
type OrderItem struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	OrderID     string    `json:"order_id" gorm:"not null"`
	ProductID   string    `json:"product_id" gorm:"not null"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Order       Order     `json:"order" gorm:"foreignkey:OrderID"`
	Product     Product   `json:"product" gorm:"foreignkey:ProductID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// UpdateCartRequest 更新购物车请求
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// DeleteCartRequest 删除购物车请求
type DeleteCartRequest struct {
	ProductIds []string `json:"productIds" binding:"required"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Remark string `json:"remark"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

// CartResponse 购物车响应
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"totalPrice"`
}

// CartItemResponse 购物车商品响应
type CartItemResponse struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
