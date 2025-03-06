-- name: GetState :one
SELECT * FROM states WHERE id = $1;