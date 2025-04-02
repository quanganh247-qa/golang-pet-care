-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_number,
    amount,
    date,
    due_date,
    status,
    description,
    customer_name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: CreateInvoiceItem :one
INSERT INTO invoice_items (
    invoice_id,
    name,
    price,
    quantity
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoices
WHERE id = $1 LIMIT 1;

-- name: GetInvoiceByNumber :one
SELECT * FROM invoices
WHERE invoice_number = $1 LIMIT 1;

-- name: GetInvoiceItems :many
SELECT * FROM invoice_items
WHERE invoice_id = $1;

-- name: ListInvoices :many
SELECT * FROM invoices
ORDER BY date DESC
LIMIT $1
OFFSET $2;

-- name: UpdateInvoiceStatus :exec
UPDATE invoices
SET status = $2
WHERE id = $1;

-- name: DeleteInvoice :exec
DELETE FROM invoices
WHERE id = $1;

-- name: GetInvoiceWithItems :many
SELECT 
    i.id,
    i.invoice_number,
    i.amount,
    i.date,
    i.due_date,
    i.status,
    i.description,
    i.customer_name,
    i.created_at,
    ii.id as item_id,
    ii.name as item_name,
    ii.price as item_price,
    ii.quantity as item_quantity
FROM invoices i
LEFT JOIN invoice_items ii ON i.id = ii.invoice_id
WHERE i.id = $1;

-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices;

-- name: UpdateInvoiceItem :one
UPDATE invoice_items
SET 
    name = COALESCE($2, name),
    price = COALESCE($3, price),
    quantity = COALESCE($4, quantity)
WHERE id = $1
RETURNING *;

-- name: DeleteInvoiceItem :exec
DELETE FROM invoice_items
WHERE id = $1;

-- name: GetInvoiceItemByID :one
SELECT * FROM invoice_items
WHERE id = $1;

-- name: UpdateInvoiceAmount :exec
UPDATE invoices
SET 
    amount = (
        SELECT COALESCE(SUM(price * quantity), 0) 
        FROM invoice_items 
        WHERE invoice_id = $1
    )
WHERE id = $1;