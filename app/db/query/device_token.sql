-- name: InsertDeviceToken :one
<<<<<<< HEAD
INSERT INTO device_tokens (
    username, token, device_type, last_used_at, expired_at
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetDeviceTokenByUsername :many
SELECT * FROM device_tokens WHERE username = $1;

-- name: DeleteDeviceToken :exec
DELETE FROM device_tokens WHERE username = $1 AND token = $2;
=======
INSERT INTO DeviceTokens (
    username,token,device_type,last_used_at,expired_at
)VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: GetDeviceTokenByUsername :many
SELECT * FROM DeviceTokens WHERE username = $1;

>>>>>>> 0fb3f30 (user images)
