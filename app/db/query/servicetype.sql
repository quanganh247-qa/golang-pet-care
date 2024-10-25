-- name: CreateServiceType :one
INSERT INTO ServiceType (
 ServiceTypeName
) VALUES (
  $1
) RETURNING *;