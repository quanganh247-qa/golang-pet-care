-- name: CreatePayment :one
INSERT INTO payments (
    amount,
    payment_method,
    payment_status,
    order_id,
    test_order_id,
    appointment_id,
    transaction_id,
    payment_details
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments
WHERE order_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetPaymentByTestOrderID :one
SELECT * FROM payments
WHERE test_order_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET 
    payment_status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetRevenueLastSevenDays :many
SELECT 
    DATE(created_at) as date,
    SUM(amount) as total_revenue
FROM payments
WHERE 
    payment_status = 'completed' 
    AND created_at >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- name: GetRevenueByPaymentMethod :many
SELECT 
    payment_method,
    SUM(amount) as total_revenue,
    COUNT(*) as transaction_count
FROM payments
WHERE 
    payment_status = 'completed'
    AND created_at >= $1
    AND created_at <= $2
GROUP BY payment_method
ORDER BY total_revenue DESC;