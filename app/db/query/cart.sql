-- name: AddItemToCart :one
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
RETURNING *;


-- name: GetCartByUserId :many
SELECT * 
FROM Cart
WHERE user_id = $1;

-- name: CreateCartForUser :one
INSERT INTO Cart (user_id)
VALUES ($1)
RETURNING id AS cart_id;

-- name: GetCartItems :many
SELECT 
    CartItem.*,
    Products.name AS product_name
FROM CartItem
JOIN Products ON CartItem.product_id = Products.product_id
WHERE CartItem.cart_id = $1;

-- name: GetCartTotal :one
SELECT SUM(total_price)::FLOAT8
FROM CartItem
WHERE cart_id = $1;

-- name: CreateOrder :one
INSERT INTO Orders (user_id, total_amount, cart_items, shipping_address, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *; -- Returning fields you may want to use