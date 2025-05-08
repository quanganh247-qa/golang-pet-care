-- name: CreateProductStockMovement :one
INSERT INTO product_stock_movements (
    product_id,
    movement_type,
    quantity,
    reason,
    price -- Add price column
)
VALUES (
    $1, $2, $3, $4, $5 -- Add price parameter
)
RETURNING *;

-- name: GetProductStockMovementsByProductID :many
SELECT
    movement_id,
    product_id,
    movement_type,
    quantity,
    reason,
    movement_date,
    price -- Select price column
FROM product_stock_movements
WHERE product_id = $1
ORDER BY movement_date DESC
LIMIT $2
OFFSET $3;

-- name: GetTotalStockMovementsByProductID :one
SELECT
    product_id,
    COALESCE(SUM(CASE WHEN movement_type = 'import' THEN quantity ELSE 0 END), 0) AS total_imported,
    COALESCE(SUM(CASE WHEN movement_type = 'export' THEN quantity ELSE 0 END), 0) AS total_exported
FROM product_stock_movements
WHERE product_id = $1
GROUP BY product_id;

-- name: GetAllProductStockMovements :many
SELECT * FROM product_stock_movements
ORDER BY movement_date DESC
LIMIT $1 OFFSET $2;
