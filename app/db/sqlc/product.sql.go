// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: product.sql

package db

import (
	"context"
)

const getAllProducts = `-- name: GetAllProducts :many
SELECT product_id, name, description, price, stock_quantity, category, data_image, original_image, created_at, is_available, removed_at from Products  ORDER BY name  LIMIT $1 OFFSET $2
`

type GetAllProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllProducts(ctx context.Context, arg GetAllProductsParams) ([]Product, error) {
	rows, err := q.db.Query(ctx, getAllProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ProductID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.StockQuantity,
			&i.Category,
			&i.DataImage,
			&i.OriginalImage,
			&i.CreatedAt,
			&i.IsAvailable,
			&i.RemovedAt,
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

const getProductByID = `-- name: GetProductByID :one
SELECT product_id, name, description, price, stock_quantity, category, data_image, original_image, created_at, is_available, removed_at from Products where product_id = $1
`

func (q *Queries) GetProductByID(ctx context.Context, productID int64) (Product, error) {
	row := q.db.QueryRow(ctx, getProductByID, productID)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StockQuantity,
		&i.Category,
		&i.DataImage,
		&i.OriginalImage,
		&i.CreatedAt,
		&i.IsAvailable,
		&i.RemovedAt,
	)
	return i, err
}
