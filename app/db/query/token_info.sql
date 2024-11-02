-- name: InsertTokenInfo :one
INSERT INTO token_info (access_token, refresh_token, token_type ,user_name, expiry)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTokenInfo :one
SELECT * FROM token_info WHERE user_name = $1;
