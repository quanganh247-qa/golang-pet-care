-- name: GetProductByID :one
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
SELECT * from Products where product_id = $1;

-- name: GetAllProducts :many
SELECT * from Products  ORDER BY name  LIMIT $1 OFFSET $2;

-- name: InsertProduct :one
INSERT INTO Products (name, description, price, category, stock_quantity,data_image,original_image,created_at,is_available) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;
<<<<<<< HEAD
=======
SELECT * from Products where product_id = $1;
>>>>>>> 21608b5 (cart and order api)
=======
SELECT * from Products where product_id = $1;

-- name: GetAllProducts :many
SELECT * from Products  ORDER BY name  LIMIT $1 OFFSET $2;
>>>>>>> bd5945b (get list products)
=======
>>>>>>> 1ec1fee (create product api)
=======
SELECT * from Products where product_id = $1;
>>>>>>> 21608b5 (cart and order api)
=======
SELECT * from Products where product_id = $1;

-- name: GetAllProducts :many
SELECT * from Products  ORDER BY name  LIMIT $1 OFFSET $2;
>>>>>>> bd5945b (get list products)
