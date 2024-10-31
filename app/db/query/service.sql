-- name: CreateService :one
INSERT INTO Service (
  typeID,
  name,
  price,
  duration,
  description,
  isAvailable
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetServiceByID :one
SELECT * FROM Service WHERE serviceID = $1 LIMIT 1;

-- name: GetAllServices :many
SELECT * FROM Service ORDER BY serviceID LIMIT $1 OFFSET $2;

-- name: DeleteService :exec
DELETE FROM Service WHERE serviceID = $1;

-- name: UpdateService :exec
UPDATE Service Set
  typeID = $2,
  name = $3,
  price = $4,
  duration = $5,
  description = $6,
  isAvailable = $7
WHERE serviceID = $1;