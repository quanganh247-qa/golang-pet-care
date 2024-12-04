-- name: AddItemToCart :exec
WITH product_check AS (
    SELECT id FROM CartItem 
    WHERE CartItem.cart_id = $1 AND CartItem.product_id = $2
)
UPDATE CartItem
SET quantity = CartItem.quantity + $3
WHERE CartItem.cart_id = $1 AND CartItem.product_id = $2
RETURNING *;

-- If the product does not exist, insert a new CartItem
INSERT INTO CartItem (cart_id, product_id, quantity, unit_price, total_price)
SELECT $1, $2, $3, Products.price, $3 * Products.price
FROM Products
WHERE Products.product_id = $2
ON CONFLICT (cart_id, product_id) DO NOTHING
RETURNING *;

-- name: GetCartByUserId :many
SELECT * 
FROM Cart
WHERE user_id = $1;

-- name: CreateCartForUser :one
INSERT INTO Cart (user_id)
VALUES ($1)
RETURNING id AS cart_id;

-- UPDATE CartItem
-- SET quantity = CartItem.quantity + $3, 
--     total_price = CartItem.quantity * CartItem.unit_price
-- WHERE CartItem.cart_id = $1 AND CartItem.product_id = $2
-- RETURNING *;