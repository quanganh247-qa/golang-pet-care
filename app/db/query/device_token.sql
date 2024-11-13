-- name: InsertDeviceToken :one
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
INSERT INTO device_tokens (
    username, token, device_type, last_used_at, expired_at
) VALUES ($1, $2, $3, $4, $5) RETURNING *;
>>>>>>> 33fcf96 (Big update)

-- name: GetDeviceTokenByUsername :many
SELECT * FROM device_tokens WHERE username = $1;

<<<<<<< HEAD
>>>>>>> 0fb3f30 (user images)
=======
-- name: DeleteDeviceToken :exec
<<<<<<< HEAD
DELETE FROM DeviceTokens WHERE username = $1 AND token = $2;
>>>>>>> 9d28896 (image pet)
=======
DELETE FROM device_tokens WHERE username = $1 AND token = $2;
>>>>>>> 33fcf96 (Big update)
=======
INSERT INTO DeviceTokens (
    username,token,device_type,last_used_at,expired_at
)VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: GetDeviceTokenByUsername :many
SELECT * FROM DeviceTokens WHERE username = $1;

>>>>>>> 0fb3f30 (user images)
