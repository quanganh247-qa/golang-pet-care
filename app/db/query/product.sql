-- name: GetProductByID :one
SELECT * from Products where product_id = $1;