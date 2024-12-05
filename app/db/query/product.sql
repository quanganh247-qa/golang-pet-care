-- name: GetProductByID :one
SELECT * from Products where product_id = $1;

-- name: GetAllProducts :many
SELECT * from Products  ORDER BY name  LIMIT $1 OFFSET $2;