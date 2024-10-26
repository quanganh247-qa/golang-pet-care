-- name: CreateService :one
INSERT INTO Service (
  TypeID,
  Name,
  Price,
  Duration,
  Description,
  IsAvailable
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetService :one
SELECT * FROM Service WHERE ServiceID = $1 LIMIT 1;

-- name: DeleteService :exec
DELETE FROM Service WHERE ServiceID = $1;