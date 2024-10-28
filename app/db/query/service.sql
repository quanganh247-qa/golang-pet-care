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

-- name: GetServiceByID :one
SELECT * FROM Service WHERE ServiceID = $1 LIMIT 1;

-- name: DeleteService :exec
DELETE FROM Service WHERE ServiceID = $1;

-- name: GetAllServices :many
SELECT * FROM Service;

-- name: UpdateService :exec
UPDATE Service SET
  TypeID = $2,
  Name = $3,
  Price = $4,
  Duration = $5,
  Description = $6,
  IsAvailable = $7
WHERE ServiceID = $1;