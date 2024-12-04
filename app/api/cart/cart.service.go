package cart

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type CartServiceInterface interface {
	AddToCartService(c *gin.Context, req CartItem, username string) error
}

func (s *CartService) AddToCartService(c *gin.Context, req CartItem, username string) error {

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

		err = q.AddItemToCart(c, db.AddItemToCartParams{
			CartID:    cartID,
			ProductID: req.ProductID,
			Quantity:  pgtype.Int4{Int32: int32(req.Quantity), Valid: true},
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	return nil
}
