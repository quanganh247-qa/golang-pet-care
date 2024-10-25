-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: UpdateUser :one
UPDATE users 
SET 
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  full_name =  COALESCE(sqlc.narg(full_name),full_name),
  email = COALESCE(sqlc.narg(email),email),
  is_verified_email = COALESCE(sqlc.narg(is_verified_email),is_verified_email)
WHERE
  username = sqlc.arg(username)
RETURNING *;

