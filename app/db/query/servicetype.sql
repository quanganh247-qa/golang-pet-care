-- name: CreateServiceType :one
INSERT INTO ServiceType (
 ServiceTypeName,
 Description,
 IconURL
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetServiceType :one
SELECT * FROM ServiceType WHERE TypeID = $1 LIMIT 1;

-- name: DeleteServiceType :exec
DELETE FROM ServiceType WHERE TypeID = $1;