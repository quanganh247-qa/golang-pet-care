-- name: CreateMedicine :one
INSERT INTO medicines (
  name, 
  description, 
  usage, 
  dosage, 
  frequency, 
  duration, 
  side_effects, 
  expiration_date, 
  quantity,
  unit_price,
  reorder_level
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateMedicine :exec
UPDATE medicines 
SET 
  name = $2,
  description = $3,
  usage = $4,
  dosage = $5,
  frequency = $6,
  duration = $7,
  side_effects = $8,
  expiration_date = $9,
  quantity = $10,
  unit_price = $11,
  reorder_level = $12,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ListMedicinesByPet :many
SELECT 
    m.*,
    pm.dosage,
    pm.frequency,
    pm.duration,
    pm.notes AS medicine_notes,
    pt.start_date AS treatment_start_date,
    pt.end_date AS treatment_end_date,
    pt.status AS treatment_status
FROM 
    pet_treatments pt
JOIN 
    treatment_phases tp ON pt.disease_id = tp.disease_id
JOIN 
    phase_medicines pm ON tp.id = pm.phase_id
JOIN 
    medicines m ON pm.medicine_id = m.id
WHERE 
    pt.pet_id = $1 and pt.status = $2 -- Replace with the specific pet_id
ORDER BY 
    tp.start_date, pm.medicine_id LIMIT $3 OFFSET $4;

-- name: GetMedicineByID :one
SELECT * FROM medicines
WHERE id = $1 LIMIT 1;

-- name: GetAllMedicines :many
SELECT 
    medicines.id, 
    medicines.name, 
    medicines.description, 
    medicines.usage, 
    medicines.dosage, 
    medicines.frequency, 
    medicines.duration, 
    medicines.side_effects, 
    medicines.created_at, 
    medicines.updated_at, 
    medicines.expiration_date, 
    medicines.quantity,
    medicines.unit_price,
    medicines.reorder_level,
    ms.name as supplier_name
FROM medicines
JOIN medicine_suppliers ms ON medicines.supplier_id = ms.id
ORDER BY medicines.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllMedicines :one
SELECT COUNT(*) FROM medicines;

-- name: DeleteMedicine :exec
DELETE FROM medicines WHERE id = $1;

-- name: UpdateMedicineQuantity :exec
UPDATE medicines
SET 
  quantity = quantity + $2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: SearchMedicinesByName :many
SELECT * FROM medicines
WHERE name ILIKE $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- Supplier management

-- name: CreateSupplier :one
INSERT INTO medicine_suppliers (
  name,
  email,
  phone,
  address,
  contact_name,
  notes
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSupplierByID :one
SELECT * FROM medicine_suppliers
WHERE id = $1 LIMIT 1;

-- name: GetAllSuppliers :many
SELECT * FROM medicine_suppliers
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateSupplier :exec
UPDATE medicine_suppliers
SET
  name = $2,
  email = $3,
  phone = $4,
  address = $5,
  contact_name = $6,
  notes = $7,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteSupplier :exec
DELETE FROM medicine_suppliers WHERE id = $1;

-- Transaction management

-- name: CreateMedicineTransaction :one
INSERT INTO medicine_transactions (
  medicine_id,
  quantity,
  transaction_type,
  unit_price,
  total_amount,
  supplier_id,
  expiration_date,
  notes,
  prescription_id,
  appointment_id,
  created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetMedicineTransactions :many
SELECT 
  mt.*,
  m.name as medicine_name,
  COALESCE(ms.name, '') as supplier_name
FROM medicine_transactions mt
JOIN medicines m ON mt.medicine_id = m.id
LEFT JOIN medicine_suppliers ms ON mt.supplier_id = ms.id
ORDER BY mt.transaction_date DESC
LIMIT $1 OFFSET $2;

-- name: GetMedicineTransactionsByMedicineID :many
SELECT 
  mt.*,
  m.name as medicine_name,
  COALESCE(ms.name, '') as supplier_name
FROM medicine_transactions mt
JOIN medicines m ON mt.medicine_id = m.id
LEFT JOIN medicine_suppliers ms ON mt.supplier_id = ms.id
WHERE mt.medicine_id = $1
ORDER BY mt.transaction_date DESC
LIMIT $2 OFFSET $3;

-- name: GetMedicineTransactionsByDate :many
SELECT 
  mt.*,
  m.name as medicine_name,
  COALESCE(ms.name, '') as supplier_name
FROM medicine_transactions mt
JOIN medicines m ON mt.medicine_id = m.id
LEFT JOIN medicine_suppliers ms ON mt.supplier_id = ms.id
WHERE mt.transaction_date BETWEEN $1 AND $2
ORDER BY mt.transaction_date DESC;

-- name: GetExpiringMedicines :many
SELECT 
    m.id,
    m.name,
    m.expiration_date,
    (m.expiration_date - CURRENT_DATE) as days_until_expiry,
    m.quantity
FROM medicines m
WHERE 
    m.expiration_date IS NOT NULL
    AND m.expiration_date <= $1::date  -- Pass a date parameter
    AND m.expiration_date >= CURRENT_DATE
    AND m.quantity > 0
ORDER BY m.expiration_date ASC;

-- name: GetLowStockMedicines :many
SELECT 
    m.id,
    m.name,
    m.quantity as current_stock,
    m.reorder_level
FROM 
    medicines m
WHERE 
    m.quantity < m.reorder_level
ORDER BY 
    (m.reorder_level - m.quantity) DESC;

