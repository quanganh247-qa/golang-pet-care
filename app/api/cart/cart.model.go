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
