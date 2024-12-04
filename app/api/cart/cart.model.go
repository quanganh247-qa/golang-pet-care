package cart

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

type CartItem struct {
	ID         int64   `json:"id"`
	CartID     int64   `json:"cart_id"`
	ProductID  int64   `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

<<<<<<< HEAD
type CartItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CartItemResponse struct {
	ID          int64   `json:"id"`
	CartID      int64   `json:"cart_id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// PlaceOrderRequest represents the request body for placing an order
type PlaceOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"` // Shipping address
	Notes           string `json:"notes"`                               // Optional notes
}

type Order struct {
	ID              int64              `json:"id"`
	UserID          int64              `json:"user_id"`
	TotalAmount     float64            `json:"total_amount"`
	PaymentStatus   string             `json:"payment_status"`
	OrderDate       string             `json:"order_date"`
	ShippingAddress string             `json:"shipping_address"`
	CartItems       []CartItemResponse `json:"cart_items"`
}

// PlaceOrderResponse represents the response body after placing an order
type OrderResponse struct {
	OrderID       int64   `json:"order_id"`
	OrderDate     string  `json:"order_date"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentStatus string  `json:"payment_status"`
=======
type Order struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	TotalAmount     float64 `json:"total_amount"`
	PaymentStatus   string  `json:"payment_status"`
	ShippingAddress string  `json:"shipping_address"`
}

type OrderItem struct {
	ID         int64   `json:"id"`
	OrderID    int64   `json:"order_id"`
	ProductID  int64   `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
>>>>>>> c449ffc (feat: cart api)
}

type CartApi struct {
	controller CartControllerInterface
}

type CartController struct {
	service CartServiceInterface
}

type CartService struct {
	storeDB db.Store
	redis   *redis.ClientType
}
