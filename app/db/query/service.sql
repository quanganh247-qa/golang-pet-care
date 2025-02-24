-- name: CreateService :one
INSERT INTO services (
    name, description, duration, cost, category, notes
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetServices :many
SELECT * FROM services where removed_at is NULL ORDER BY name LIMIT $1 OFFSET $2;

-- name: GetServiceByID :one
SELECT * FROM services
WHERE id = $1;

-- name: DeleteService :exec
UPDATE services
SET removed_at = NOW()
WHERE id = $1;

-- name: UpdateService :one
UPDATE services
SET 
    name = $2,
    description = $3,
    duration = $4,
    cost = $5,
    category = $6,
    notes = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
