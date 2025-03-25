<<<<<<< HEAD
<<<<<<< HEAD
-- name: AddItemToCart :one
INSERT INTO cart_items (
    cart_id,
    product_id,
<<<<<<< HEAD
<<<<<<< HEAD
    unit_price,
    quantity
) VALUES (
    $1, $2, $3, $4
=======
    quantity,
    created_at,
    updated_at
) VALUES (
<<<<<<< HEAD
    $1, $2, $3,  $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
>>>>>>> 33fcf96 (Big update)
=======
    $1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
>>>>>>> ffc9071 (AI suggestion)
=======
    unit_price,
    quantity
) VALUES (
    $1, $2, $3, $4
>>>>>>> dc47646 (Optimize SQL query)
) RETURNING *;


-- name: GetCartByUserId :many
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> dc47646 (Optimize SQL query)
SELECT 
    c.id,
    c.user_id,
    c.created_at,
    c.updated_at
FROM carts c
LEFT JOIN cart_items ci ON ci.cart_id = c.id
WHERE c.user_id = $1
GROUP BY c.id, c.user_id, c.created_at, c.updated_at;
<<<<<<< HEAD

-- name: GetCartItemsByUserId :many
SELECT
    ci.*,
    p.name as product_name,
    p.price as unit_price,
    (p.price * ci.quantity) as total_price
FROM cart_items ci
JOIN products p ON ci.product_id = p.product_id
LEFT JOIN carts c ON ci.cart_id = c.id
WHERE c.user_id = $1;

=======
SELECT * 
FROM carts
WHERE user_id = $1;
>>>>>>> 33fcf96 (Big update)
=======
>>>>>>> dc47646 (Optimize SQL query)

-- name: GetCartItemsByUserId :many
SELECT
    ci.*,
    p.name as product_name,
    p.price as unit_price,
    (p.price * ci.quantity) as total_price
FROM cart_items ci
JOIN products p ON ci.product_id = p.product_id
LEFT JOIN carts c ON ci.cart_id = c.id
WHERE c.user_id = $1;


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
<<<<<<< HEAD
<<<<<<< HEAD
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
<<<<<<< HEAD
LIMIT $2 OFFSET $3;
=======
-- name: AddItemToCart :exec
WITH product_check AS (
    SELECT id FROM CartItem 
    WHERE CartItem.cart_id = $1 AND CartItem.product_id = $2
)
UPDATE CartItem
SET quantity = CartItem.quantity + $3
WHERE CartItem.cart_id = $1 AND CartItem.product_id = $2
=======
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
>>>>>>> 21608b5 (cart and order api)
RETURNING *;


-- name: GetCartByUserId :many
SELECT * 
FROM Cart
WHERE user_id = $1;

-- name: CreateCartForUser :one
INSERT INTO Cart (user_id)
VALUES ($1)
RETURNING id AS cart_id;

<<<<<<< HEAD
-- UPDATE CartItem
-- SET quantity = CartItem.quantity + $3, 
--     total_price = CartItem.quantity * CartItem.unit_price
-- WHERE CartItem.cart_id = $1 AND CartItem.product_id = $2
-- RETURNING *;
>>>>>>> c449ffc (feat: cart api)
=======
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
<<<<<<< HEAD
>>>>>>> 21608b5 (cart and order api)
=======

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
>>>>>>> b0fe977 (place order and make payment)
=======
DELETE FROM CartItem
=======
DELETE FROM cart_items 
>>>>>>> 33fcf96 (Big update)
WHERE cart_id = $1 AND product_id = $2;

-- name: DecreaseItemQuantity :exec
UPDATE cart_items
SET quantity = quantity - $3
WHERE cart_id = $1 AND product_id = $2 AND quantity > $3
RETURNING *;
<<<<<<< HEAD
>>>>>>> 4a16bfc (remove item in cart)
=======

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items 
SET 
    quantity = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE cart_id = $1 AND product_id = $2;
>>>>>>> 33fcf96 (Big update)
=======
LIMIT $2 OFFSET $3;
>>>>>>> dc47646 (Optimize SQL query)
