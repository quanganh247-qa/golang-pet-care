-- name: CreateServiceType :one
INSERT INTO ServiceType (
 serviceTypeName,
 description,
 iconURL
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetServiceType :one
SELECT * FROM ServiceType WHERE typeID = $1 LIMIT 1;

-- name: DeleteServiceType :exec
DELETE FROM ServiceType WHERE typeID = $1;