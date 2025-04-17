-- name: CreateSMTPConfig :one
INSERT INTO smtp_configs (
    name,
    email,
    password,
    smtp_host,
    smtp_port,
    is_default
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSMTPConfig :one
SELECT *
FROM smtp_configs
WHERE id = $1;

-- name: GetDefaultSMTPConfig :one
SELECT *
FROM smtp_configs
WHERE is_default = true
LIMIT 1;

-- name: ListSMTPConfigs :many
SELECT *
FROM smtp_configs
ORDER BY created_at DESC;

-- name: UpdateSMTPConfig :one
UPDATE smtp_configs
SET 
    name = $2,
    email = $3,
    password = $4,
    smtp_host = $5,
    smtp_port = $6,
    is_default = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetAsDefaultSMTPConfig :exec
UPDATE smtp_configs
SET is_default = false
WHERE is_default = true;

-- name: DeleteSMTPConfig :exec
DELETE FROM smtp_configs
WHERE id = $1; 