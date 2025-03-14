-- name: AddItemToCart :one
INSERT INTO cart_items (
    cart_id,
    product_id,
    unit_price,
    quantity
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: GetCartByUserId :many
SELECT 
    c.id,
    c.user_id,
    c.created_at,
    c.updated_at
FROM carts c
LEFT JOIN cart_items ci ON ci.cart_id = c.id
WHERE c.user_id = $1
GROUP BY c.id, c.user_id, c.created_at, c.updated_at;

-- name: CreateCartForUser :one
INSERT INTO carts (
    user_id,
    created_at,
    updated_at
) VALUES (
    $1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetCartItems :many
SELECT 
    ci.*,
    p.name as product_name,
    p.price as unit_price,
    (p.price * ci.quantity) as total_price
FROM cart_items ci
JOIN products p ON ci.product_id = p.product_id
WHERE ci.cart_id = $1;

-- name: GetCartTotal :one
SELECT SUM(total_price)::FLOAT8
FROM cart_items
WHERE cart_id = $1;

-- name: CreateOrder :one
INSERT INTO Orders (user_id, total_amount, cart_items, shipping_address, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *; -- Returning fields you may want to use

-- name: GetOrdersByUserId :many
SELECT *
FROM Orders
WHERE user_id = $1 and payment_status = 'pending';

-- name: GetOrderById :one
SELECT *
FROM Orders
WHERE id = $1;

-- name: UpdateOrderPaymentStatus :one
UPDATE Orders
SET payment_status = 'paid'
WHERE id = $1 Returning *;

-- name: RemoveItemFromCart :exec
DELETE FROM cart_items 
WHERE cart_id = $1 AND product_id = $2;

-- name: DecreaseItemQuantity :exec
UPDATE cart_items
SET quantity = quantity - $3
WHERE cart_id = $1 AND product_id = $2 AND quantity > $3
RETURNING *;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items 
SET 
    quantity = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE cart_id = $1 AND product_id = $2;

-- name: GetAllOrders :many
SELECT *
FROM Orders 
WHERE payment_status = $1
ORDER BY order_date DESC
LIMIT $2 OFFSET $3;