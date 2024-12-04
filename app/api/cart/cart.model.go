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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> dc47646 (Optimize SQL query)
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
=======
type CartItemResponse struct {
	ID          int64   `json:"id"`
	CartID      int64   `json:"cart_id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
<<<<<<< HEAD
	Quantity    int     `json:"quantity"`
>>>>>>> 21608b5 (cart and order api)
=======
	Quantity    int32   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
>>>>>>> b0fe977 (place order and make payment)
	TotalPrice  float64 `json:"total_price"`
}

// PlaceOrderRequest represents the request body for placing an order
type PlaceOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"` // Shipping address
	Notes           string `json:"notes"`                               // Optional notes
}

<<<<<<< HEAD
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
=======
>>>>>>> 21608b5 (cart and order api)
type Order struct {
	ID              int64              `json:"id"`
	UserID          int64              `json:"user_id"`
	TotalAmount     float64            `json:"total_amount"`
	PaymentStatus   string             `json:"payment_status"`
	OrderDate       string             `json:"order_date"`
	ShippingAddress string             `json:"shipping_address"`
	CartItems       []CartItemResponse `json:"cart_items"`
}

<<<<<<< HEAD
=======
type Order struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	TotalAmount     float64 `json:"total_amount"`
	PaymentStatus   string  `json:"payment_status"`
	ShippingAddress string  `json:"shipping_address"`
}

>>>>>>> c449ffc (feat: cart api)
type OrderItem struct {
	ID         int64   `json:"id"`
	OrderID    int64   `json:"order_id"`
	ProductID  int64   `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
<<<<<<< HEAD
>>>>>>> c449ffc (feat: cart api)
=======
// PlaceOrderResponse represents the response body after placing an order
<<<<<<< HEAD
type PlaceOrderResponse struct {
	OrderID       int64  `json:"order_id"`       // ID of the created order
	OrderDate     string `json:"order_date"`     // Date the order was placed
	PaymentStatus string `json:"payment_status"` // Status of the payment
>>>>>>> 21608b5 (cart and order api)
=======
type OrderResponse struct {
	OrderID       int64   `json:"order_id"`
	OrderDate     string  `json:"order_date"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentStatus string  `json:"payment_status"`
>>>>>>> b0fe977 (place order and make payment)
=======
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
