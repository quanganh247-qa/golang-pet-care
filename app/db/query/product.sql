-- name: GetProductByID :one
SELECT * from Products where product_id = $1;

-- name: GetAllProducts :many
SELECT * from Products  ORDER BY name  LIMIT $1 OFFSET $2;

-- name: InsertProduct :one
INSERT INTO Products (name, description, price, category, stock_quantity,data_image,original_image,created_at,is_available) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;