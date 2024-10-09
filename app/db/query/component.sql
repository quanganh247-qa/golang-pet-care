-- name: CreateComponents :one
INSERT INTO components (name, description, content , component_code,project_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(),now())
RETURNING *;

-- name: UpdateComponents :one
UPDATE components
SET name = $1, description = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteComponents :one
DELETE FROM components
WHERE id = $1
RETURNING *;

-- name: GetComponentss :many
SELECT * FROM components;

-- name: GetComponentsById :one
SELECT * FROM components
WHERE id = $1 and project_id = $2 and removed_at is null;

-- name: GetComponentsByUser :many
SELECT * FROM components
LEFT JOIN projects on components.project_id = projects.id 
WHERE username = $1 and project_id = $2 and removed_at is null
order by components.updated_at desc  limit $3 offset $4;

-- name: CountingComponentsByUser :one
SELECT count(*) FROM components
LEFT JOIN projects on components.project_id = projects.id 
WHERE username = $1 and project_id = $2 and removed_at is null;

-- name: GetComponentsByName :one
SELECT * FROM components
WHERE name = $1;

-- name: RemoveComponents :one
UPDATE components
SET removed_at = now()
WHERE id = $1 and removed_at is null and project_id = $2
RETURNING *;


