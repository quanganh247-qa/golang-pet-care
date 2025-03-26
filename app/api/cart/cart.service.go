package cart

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type CartServiceInterface interface {
	AddToCartService(c *gin.Context, req CartItemRequest, username string) (*CartItem, error)
	GetCartItemsService(c *gin.Context, username string) ([]CartItemResponse, error)
	CreateOrderService(c *gin.Context, username string, arg PlaceOrderRequest) (*OrderResponse, error)
	GetOrdersService(c *gin.Context, username string) ([]OrderResponse, error)
	GetOrderByIdService(c *gin.Context, username string, orderID int64) (*Order, error)
	DeleteItemFromCartService(c *gin.Context, username string, itemID int64) error
	GetAllOrdersService(c *gin.Context, pagination *util.Pagination, status string) ([]OrderResponse, error)
}

func (s *CartService) AddToCartService(c *gin.Context, req CartItemRequest, username string) (*CartItem, error) {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	var product db.Product
	var carts []db.Cart
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		product, err = s.storeDB.GetProductByID(c, req.ProductID)
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		res, err := s.storeDB.GetCartByUserId(c, user.UserID)
		if err != nil {
			return
		}

		for _, cart := range res {
			carts = append(carts, db.Cart{
				ID:        cart.ID,
				UserID:    cart.UserID,
				CreatedAt: cart.CreatedAt,
				UpdatedAt: cart.UpdatedAt,
			})
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("error fetching data: %w", err)
	}

	// Check if cart exists, if not create one
	var cartID int64
	if len(carts) == 0 {
		newCart, err := s.storeDB.CreateCartForUser(c, user.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to create cart: %w", err)
		}
		cartID = newCart.ID
	} else {
		cartID = carts[0].ID
	}

	// Now use the product information
	cartItem, err := s.storeDB.AddItemToCart(c, db.AddItemToCartParams{
		CartID:    cartID,
		ProductID: req.ProductID,
		UnitPrice: product.Price,
		Quantity:  pgtype.Int4{Int32: int32(req.Quantity), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	return &CartItem{
		ID:         cartItem.ID,
		CartID:     cartItem.CartID,
		ProductID:  cartItem.ProductID,
		Quantity:   int(cartItem.Quantity.Int32),
		UnitPrice:  cartItem.UnitPrice,
		TotalPrice: cartItem.TotalPrice.Float64,
	}, nil
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

	if len(cart) == 0 {
		return nil, fmt.Errorf("cart not found")
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
			UnitPrice:   product.Price,
			ProductID:   cart.ProductID,
			Quantity:    cart.Quantity.Int32,
			TotalPrice:  cart.TotalPrice.Float64,
		})
	}

	return items, nil

}

func (s *CartService) CreateOrderService(c *gin.Context, username string, arg PlaceOrderRequest) (*OrderResponse, error) {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	carts, err := s.storeDB.GetCartByUserId(c, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart by user id: %w", err)
	}

	// Fetch the total price for a cart
	totalPriceRow, err := s.storeDB.GetCartTotal(c, carts[0].ID)
	if err != nil {
		return nil, err
	}

	// Convert cartItems slice to JSON
	jsonData, err := json.Marshal(carts)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cart items to JSON: %w", err)
	}

	var placeOrder OrderResponse

	err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
		order, err := q.CreateOrder(c, db.CreateOrderParams{
			UserID:          user.UserID,
			TotalAmount:     float64(totalPriceRow),
			CartItems:       jsonData,
			ShippingAddress: pgtype.Text{String: arg.ShippingAddress, Valid: true},
			Notes:           pgtype.Text{String: arg.Notes, Valid: true},
		})
		if err != nil {
			return err
		}

		placeOrder = OrderResponse{
			OrderID:       order.ID,
			OrderDate:     order.OrderDate.Time.Format("2006-01-02"),
			TotalAmount:   order.TotalAmount,
			PaymentStatus: order.PaymentStatus.String,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &placeOrder, nil

}

// get list of orders
func (s *CartService) GetOrdersService(c *gin.Context, username string) ([]OrderResponse, error) {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	orders, err := s.storeDB.GetOrdersByUserId(c, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user id: %w", err)
	}

	var orderResponses []OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, OrderResponse{
			OrderID:       order.ID,
			OrderDate:     order.OrderDate.Time.Format("2006-01-02"),
			TotalAmount:   order.TotalAmount,
			PaymentStatus: order.PaymentStatus.String,
		})
	}

	return orderResponses, nil
}

// get oder by id

func (s *CartService) GetOrderByIdService(c *gin.Context, username string, orderID int64) (*Order, error) {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	order, err := s.storeDB.GetOrderById(c, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	if order.UserID != user.UserID {
		return nil, fmt.Errorf("order not found")
	}

	// Convert cartItems slice to JSON
	var cartItems []CartItemResponse
	err = json.Unmarshal([]byte(order.CartItems), &cartItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart items: %w", err)
	}

	orderResponse := Order{
		ID:              order.ID,
		UserID:          order.UserID,
		OrderDate:       order.OrderDate.Time.Format("2006-01-02"),
		TotalAmount:     order.TotalAmount,
		ShippingAddress: order.ShippingAddress.String,
		PaymentStatus:   order.PaymentStatus.String,
		CartItems:       cartItems,
	}

	return &orderResponse, nil
}

func (s *CartService) DeleteItemFromCartService(c *gin.Context, username string, itemID int64) error {
	user, err := s.redis.UserInfoLoadCache(username)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	cart, err := s.storeDB.GetCartByUserId(c, user.UserID)
	if err != nil {
		return fmt.Errorf("failed to get cart by user id: %w", err)
	}

	err = s.storeDB.ExecWithTransaction(c, func(q *db.Queries) error {
		err := q.RemoveItemFromCart(c, db.RemoveItemFromCartParams{
			CartID:    cart[0].ID,
			ProductID: itemID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil

}

func (s *CartService) GetAllOrdersService(c *gin.Context, pagination *util.Pagination, status string) ([]OrderResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize
	orders, err := s.storeDB.GetAllOrders(c, db.GetAllOrdersParams{
		PaymentStatus: pgtype.Text{String: status, Valid: true},
		Limit:         int32(pagination.PageSize),
		Offset:        int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var orderResponses []OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, OrderResponse{
			OrderID:       order.ID,
			OrderDate:     order.OrderDate.Time.Format("2006-01-02"),
			TotalAmount:   order.TotalAmount,
			PaymentStatus: order.PaymentStatus.String,
		})
	}

	return orderResponses, nil

}
