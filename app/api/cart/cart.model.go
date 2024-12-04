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

type CartItemResponse struct {
	ID          int64   `json:"id"`
	CartID      int64   `json:"cart_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

// PlaceOrderRequest represents the request body for placing an order
type PlaceOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"` // Shipping address
	Notes           string `json:"notes"`                               // Optional notes
}

type Order struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	TotalAmount     float64 `json:"total_amount"`
	PaymentStatus   string  `json:"payment_status"`
	ShippingAddress string  `json:"shipping_address"`
}

// PlaceOrderResponse represents the response body after placing an order
type PlaceOrderResponse struct {
	OrderID       int64  `json:"order_id"`       // ID of the created order
	OrderDate     string `json:"order_date"`     // Date the order was placed
	PaymentStatus string `json:"payment_status"` // Status of the payment
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
