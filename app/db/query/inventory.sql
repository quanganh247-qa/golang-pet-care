-- name: CreateInventoryItem :one
INSERT INTO inventory_items (
    name, 
    item_type, 
    description, 
    usage_instructions, 
    dosage, 
    frequency, 
    duration, 
    side_effects, 
    expiration_date, 
    quantity, 
    unit_price, 
    reorder_level, 
    supplier_id, 
    for_species, 
    requires_prescription, 
    storage_condition, 
    batch_number, 
    manufacturer, 
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetInventoryItem :one
SELECT 
    i.*,
    s.name as supplier_name,
    (i.quantity <= i.reorder_level) as needs_reorder,
    (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) as is_expired
FROM 
    inventory_items i
LEFT JOIN 
    medicine_suppliers s ON i.supplier_id = s.id
WHERE 
    i.id = $1;

-- name: ListInventoryItems :many
SELECT 
    i.*,
    s.name as supplier_name,
    (i.quantity <= i.reorder_level) as needs_reorder,
    (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) as is_expired
FROM 
    inventory_items i
LEFT JOIN 
    medicine_suppliers s ON i.supplier_id = s.id
WHERE 
    ($1::inventory_item_type_enum IS NULL OR i.item_type = $1)
    AND ($2::text IS NULL OR i.name ILIKE '%' || $2 || '%')
    AND ($3::boolean IS NULL OR i.is_active = $3)
    AND ($4::boolean IS NULL OR (i.quantity <= i.reorder_level) = $4)
    AND ($5::boolean IS NULL OR (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) = $5)
    AND ($6::bigint IS NULL OR i.supplier_id = $6)
    AND ($7::text IS NULL OR i.for_species = $7)
    AND ($8::boolean IS NULL OR i.requires_prescription = $8)
ORDER BY 
    i.name
LIMIT $9
OFFSET $10;

-- name: CountInventoryItems :one
SELECT COUNT(*)
FROM inventory_items i
WHERE 
    ($1::inventory_item_type_enum IS NULL OR i.item_type = $1)
    AND ($2::text IS NULL OR i.name ILIKE '%' || $2 || '%')
    AND ($3::boolean IS NULL OR i.is_active = $3)
    AND ($4::boolean IS NULL OR (i.quantity <= i.reorder_level) = $4)
    AND ($5::boolean IS NULL OR (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) = $5)
    AND ($6::bigint IS NULL OR i.supplier_id = $6)
    AND ($7::text IS NULL OR i.for_species = $7)
    AND ($8::boolean IS NULL OR i.requires_prescription = $8);

-- name: UpdateInventoryItem :exec
UPDATE inventory_items SET
    name = $1,
    item_type = $2,
    description = $3,
    usage_instructions = $4,
    dosage = $5,
    frequency = $6,
    duration = $7,
    side_effects = $8,
    expiration_date = $9,
    quantity = $10,
    unit_price = $11,
    reorder_level = $12,
    supplier_id = $13,
    for_species = $14,
    requires_prescription = $15,
    storage_condition = $16,
    batch_number = $17,
    manufacturer = $18,
    is_active = $19,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $20;

-- name: DeleteInventoryItem :exec
UPDATE inventory_items SET
    is_active = false,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetInventoryItemsByType :many
SELECT 
    i.*,
    s.name as supplier_name,
    (i.quantity <= i.reorder_level) as needs_reorder,
    (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) as is_expired
FROM 
    inventory_items i
LEFT JOIN 
    medicine_suppliers s ON i.supplier_id = s.id
WHERE 
    i.item_type = $1
    AND i.is_active = true
ORDER BY 
    i.name
LIMIT $2
OFFSET $3;

-- name: GetLowStockItems :many
SELECT 
    i.*,
    s.name as supplier_name,
    (i.quantity <= i.reorder_level) as needs_reorder,
    (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) as is_expired
FROM 
    inventory_items i
LEFT JOIN 
    medicine_suppliers s ON i.supplier_id = s.id
WHERE 
    i.quantity <= i.reorder_level
    AND i.is_active = true
ORDER BY 
    i.quantity ASC
LIMIT $1
OFFSET $2;

-- name: GetExpiringItems :many
SELECT 
    i.*,
    s.name as supplier_name,
    (i.quantity <= i.reorder_level) as needs_reorder,
    (i.expiration_date IS NOT NULL AND i.expiration_date < CURRENT_DATE) as is_expired,
    EXTRACT(DAY FROM (i.expiration_date::timestamp - CURRENT_DATE::timestamp)) as days_until_expiry
FROM 
    inventory_items i
LEFT JOIN 
    medicine_suppliers s ON i.supplier_id = s.id
WHERE 
    i.expiration_date IS NOT NULL
    AND i.expiration_date <= (CURRENT_DATE + $1::integer)
    AND i.is_active = true
    AND i.quantity > 0
ORDER BY 
    i.expiration_date ASC
LIMIT $2
OFFSET $3;

-- Transaction management
-- name: CreateInventoryTransaction :one
INSERT INTO inventory_transactions (
    inventory_item_id,
    transaction_type,
    quantity,
    unit_price,
    total_amount,
    transaction_date,
    supplier_id,
    expiration_date,
    batch_number,
    reference_id,
    reference_type,
    notes,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetInventoryTransaction :one
SELECT 
    t.*,
    i.name,
    i.item_type,
    s.name as supplier_name
FROM 
    inventory_transactions t
JOIN 
    inventory_items i ON t.inventory_item_id = i.id
LEFT JOIN 
    medicine_suppliers s ON t.supplier_id = s.id
WHERE 
    t.id = $1;

-- name: ListInventoryTransactions :many
SELECT 
    t.*,
    i.name,
    i.item_type,
    s.name as supplier_name
FROM 
    inventory_transactions t
JOIN 
    inventory_items i ON t.inventory_item_id = i.id
LEFT JOIN 
    medicine_suppliers s ON t.supplier_id = s.id
WHERE 
    ($1::bigint IS NULL OR t.inventory_item_id = $1)
    AND ($2::inventory_item_type_enum IS NULL OR i.item_type = $2)
    AND ($3::text IS NULL OR t.transaction_type = $3)
    AND ($4::bigint IS NULL OR t.supplier_id = $4)
    AND ($5::timestamptz IS NULL OR t.transaction_date >= $5)
    AND ($6::timestamptz IS NULL OR t.transaction_date <= $6)
    AND ($7::text IS NULL OR t.reference_type = $7)
    AND ($8::bigint IS NULL OR t.reference_id = $8)
    AND ($9::text IS NULL OR t.created_by = $9)
ORDER BY 
    t.transaction_date DESC
LIMIT $10
OFFSET $11;

-- name: CountInventoryTransactions :one
SELECT COUNT(*)
FROM 
    inventory_transactions t
JOIN 
    inventory_items i ON t.inventory_item_id = i.id
WHERE 
    ($1::bigint IS NULL OR t.inventory_item_id = $1)
    AND ($2::inventory_item_type_enum IS NULL OR i.item_type = $2)
    AND ($3::text IS NULL OR t.transaction_type = $3)
    AND ($4::bigint IS NULL OR t.supplier_id = $4)
    AND ($5::timestamptz IS NULL OR t.transaction_date >= $5)
    AND ($6::timestamptz IS NULL OR t.transaction_date <= $6)
    AND ($7::text IS NULL OR t.reference_type = $7)
    AND ($8::bigint IS NULL OR t.reference_id = $8)
    AND ($9::text IS NULL OR t.created_by = $9);

-- name: GetInventoryTransactionsByItem :many
SELECT 
    t.*,
    i.name,
    i.item_type,
    s.name as supplier_name
FROM 
    inventory_transactions t
JOIN 
    inventory_items i ON t.inventory_item_id = i.id
LEFT JOIN 
    medicine_suppliers s ON t.supplier_id = s.id
WHERE 
    t.inventory_item_id = $1
ORDER BY 
    t.transaction_date DESC
LIMIT $2
OFFSET $3;

-- name: GetInventorySummary :many
SELECT 
    i.id, 
    i.name, 
    i.item_type, 
    i.quantity, 
    i.unit_price, 
    (i.quantity * i.unit_price) as total_value,
    i.is_active,
    (i.quantity <= i.reorder_level) as needs_reorder
FROM 
    inventory_items i
WHERE 
    ($1::inventory_item_type_enum IS NULL OR i.item_type = $1)
    AND i.is_active = true
ORDER BY 
    i.name;

-- name: GetInventoryTransactionsSummary :many
SELECT 
    t.transaction_type,
    SUM(t.quantity) as total_quantity,
    SUM(COALESCE(t.total_amount, 0)) as total_value,
    COUNT(*) as transaction_count
FROM 
    inventory_transactions t
JOIN 
    inventory_items i ON t.inventory_item_id = i.id
WHERE 
    t.transaction_date BETWEEN $1 AND $2
    AND ($3::inventory_item_type_enum IS NULL OR i.item_type = $3)
GROUP BY 
    t.transaction_type
ORDER BY 
    t.transaction_type; 