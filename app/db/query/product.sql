-- name: GetProductByID :one
SELECT * from Products where product_id = $1;

-- name: GetAllProducts :many
SELECT * from Products  ORDER BY name  LIMIT $1 OFFSET $2;

-- name: InsertProduct :one
INSERT INTO Products (name, description, price, category, stock_quantity,data_image,original_image,created_at,is_available) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetProductByIDForUpdate :one
SELECT * from Products where product_id = $1 FOR NO KEY UPDATE;

-- name: UpdateProductStock :one
UPDATE Products
SET stock_quantity = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1
RETURNING *;

-- name: UpdateProduct :one
UPDATE Products
SET name = $2,
    description = $3,
    price = $4,
    category = $5,
    stock_quantity = $6,
    data_image = $7,
    original_image = $8,
    is_available = $9
WHERE product_id = $1
RETURNING *;


-- name: DeleteProduct :exec
DELETE FROM Products WHERE product_id = $1;