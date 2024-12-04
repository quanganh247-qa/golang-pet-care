package cart

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type CartServiceInterface interface {
	AddToCartService(c *gin.Context, req CartItem, username string) (*CartItem, error)
	GetCartItemsService(c *gin.Context, username string) ([]CartItemResponse, error)
	CreateOrderService(c *gin.Context, username string, arg PlaceOrderRequest) (*PlaceOrderResponse, error)
}

func (s *CartService) AddToCartService(c *gin.Context, req CartItem, username string) (*CartItem, error) {
	var cartItem CartItem
	err := s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
		user, err := s.redis.UserInfoLoadCache(username)
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		cart, err := s.storeDB.GetCartByUserId(c, user.UserID)
		if err != nil {
			return fmt.Errorf("failed to get cart by user id: %w", err)
		}

		var cartID int64

		if len(cart) == 0 {
			cartID, err = s.storeDB.CreateCartForUser(c, user.UserID)
			if err != nil {
				return fmt.Errorf("failed to create cart: %w", err)
			}

		} else {
			cartID = cart[0].ID
		}

		item, err := q.AddItemToCart(c, db.AddItemToCartParams{
			CartID:    cartID,
			ProductID: req.ProductID,
			Quantity:  pgtype.Int4{Int32: int32(req.Quantity), Valid: true},
		})
		if err != nil {
			return err
		}

		cartItem = CartItem{
			ID:         item.ID,
			CartID:     item.CartID,
			ProductID:  item.ProductID,
			Quantity:   int(item.Quantity.Int32), // Convert pgtype.Int4 to int32
			TotalPrice: item.TotalPrice.Float64,
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	return &cartItem, nil
}

func (s *CartService) GetCartItemsService(c *gin.Context, username string) ([]CartItemResponse, error) {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	cart, err := s.storeDB.GetCartByUserId(c, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart by user id: %w", err)
	}

	cartItems, err := s.storeDB.GetCartItems(c, cart[0].ID)
	if err != nil {
		return nil, err
	}

	var items []CartItemResponse

	for _, cart := range cartItems {

		product, _ := s.storeDB.GetProductByID(c, cart.ProductID)

		items = append(items, CartItemResponse{
			ID:          cart.ID,
			CartID:      cart.CartID,
			ProductName: product.Name,
			Quantity:    int(cart.Quantity.Int32),
			TotalPrice:  cart.TotalPrice.Float64,
		})
	}

	return items, nil

}

func (s *CartService) CreateOrderService(c *gin.Context, username string, arg PlaceOrderRequest) (*PlaceOrderResponse, error) {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	cart, err := s.storeDB.GetCartByUserId(c, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart by user id: %w", err)
	}

	cartItems, err := s.storeDB.GetCartItems(c, cart[0].ID)
	if err != nil {
		return nil, err
	}

	// Fetch the total price for a cart
	totalPriceRow, err := s.storeDB.GetCartTotal(c, cart[0].ID)
	if err != nil {
		return nil, err
	}

	// Convert cartItems slice to JSON
	jsonData, err := json.Marshal(cartItems)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cart items to JSON: %w", err)
	}

	var placeOrder PlaceOrderResponse

	err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
		order, err := q.CreateOrder(c, db.CreateOrderParams{
			UserID:          user.UserID,
			TotalAmount:     float64(totalPriceRow),
			CartItems:       []byte(jsonData),
			ShippingAddress: pgtype.Text{String: arg.ShippingAddress, Valid: true},
			Notes:           pgtype.Text{String: arg.Notes, Valid: true},
		})
		if err != nil {
			return err
		}

		placeOrder = PlaceOrderResponse{
			OrderID:       order.ID,
			OrderDate:     order.OrderDate.Time.Format("2006-01-02"),
			PaymentStatus: order.PaymentStatus.String,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &placeOrder, nil

}
