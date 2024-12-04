// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: cart.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addItemToCart = `-- name: AddItemToCart :one
INSERT INTO CartItem (cart_id, product_id, quantity, unit_price)
VALUES (
    $1, -- cart_id
    $2, -- product_id
    $3, -- quantity
    (SELECT price FROM Products WHERE product_id = $2)
 )
ON CONFLICT (cart_id, product_id)
DO UPDATE SET 
    quantity = CartItem.quantity + EXCLUDED.quantity
RETURNING id, cart_id, product_id, quantity, unit_price, total_price
`

type AddItemToCartParams struct {
	CartID    int64       `json:"cart_id"`
	ProductID int64       `json:"product_id"`
	Quantity  pgtype.Int4 `json:"quantity"`
}

func (q *Queries) AddItemToCart(ctx context.Context, arg AddItemToCartParams) (Cartitem, error) {
	row := q.db.QueryRow(ctx, addItemToCart, arg.CartID, arg.ProductID, arg.Quantity)
	var i Cartitem
	err := row.Scan(
		&i.ID,
		&i.CartID,
		&i.ProductID,
		&i.Quantity,
		&i.UnitPrice,
		&i.TotalPrice,
	)
	return i, err
}

const createCartForUser = `-- name: CreateCartForUser :one
INSERT INTO Cart (user_id)
VALUES ($1)
RETURNING id AS cart_id
`

func (q *Queries) CreateCartForUser(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRow(ctx, createCartForUser, userID)
	var cart_id int64
	err := row.Scan(&cart_id)
	return cart_id, err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO Orders (user_id, total_amount, cart_items, shipping_address, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, order_date, total_amount, payment_status, cart_items, shipping_address, notes
`

type CreateOrderParams struct {
	UserID          int64       `json:"user_id"`
	TotalAmount     float64     `json:"total_amount"`
	CartItems       []byte      `json:"cart_items"`
	ShippingAddress pgtype.Text `json:"shipping_address"`
	Notes           pgtype.Text `json:"notes"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder,
		arg.UserID,
		arg.TotalAmount,
		arg.CartItems,
		arg.ShippingAddress,
		arg.Notes,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderDate,
		&i.TotalAmount,
		&i.PaymentStatus,
		&i.CartItems,
		&i.ShippingAddress,
		&i.Notes,
	)
	return i, err
}

const getCartByUserId = `-- name: GetCartByUserId :many
SELECT id, user_id, created_at, updated_at 
FROM Cart
WHERE user_id = $1
`

func (q *Queries) GetCartByUserId(ctx context.Context, userID int64) ([]Cart, error) {
	rows, err := q.db.Query(ctx, getCartByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Cart{}
	for rows.Next() {
		var i Cart
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCartItems = `-- name: GetCartItems :many
SELECT 
    cartitem.id, cartitem.cart_id, cartitem.product_id, cartitem.quantity, cartitem.unit_price, cartitem.total_price,
    Products.name AS product_name
FROM CartItem
JOIN Products ON CartItem.product_id = Products.product_id
WHERE CartItem.cart_id = $1
`

type GetCartItemsRow struct {
	ID          int64         `json:"id"`
	CartID      int64         `json:"cart_id"`
	ProductID   int64         `json:"product_id"`
	Quantity    pgtype.Int4   `json:"quantity"`
	UnitPrice   float64       `json:"unit_price"`
	TotalPrice  pgtype.Float8 `json:"total_price"`
	ProductName string        `json:"product_name"`
}

func (q *Queries) GetCartItems(ctx context.Context, cartID int64) ([]GetCartItemsRow, error) {
	rows, err := q.db.Query(ctx, getCartItems, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCartItemsRow{}
	for rows.Next() {
		var i GetCartItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.CartID,
			&i.ProductID,
			&i.Quantity,
			&i.UnitPrice,
			&i.TotalPrice,
			&i.ProductName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCartTotal = `-- name: GetCartTotal :one
SELECT SUM(total_price)::FLOAT8
FROM CartItem
WHERE cart_id = $1
`

func (q *Queries) GetCartTotal(ctx context.Context, cartID int64) (float64, error) {
	row := q.db.QueryRow(ctx, getCartTotal, cartID)
	var column_1 float64
	err := row.Scan(&column_1)
	return column_1, err
}
