-- name: InsertDeviceToken :one
INSERT INTO DeviceTokens (
    username,token,device_type,last_used_at,expired_at
)VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: GetDeviceTokenByUsername :many
SELECT * FROM DeviceTokens WHERE username = $1;

-- name: DeleteDeviceToken :exec
DELETE FROM DeviceTokens WHERE username = $1 AND token = $2;