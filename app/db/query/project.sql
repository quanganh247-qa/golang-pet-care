-- name: CreateProject :one
INSERT INTO projects (
  username,
  name,
  description,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, now(), now()
) RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: GetProjectsByUser :many
SELECT * FROM projects WHERE username = $1;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, 
description = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :one
DELETE FROM projects WHERE id = $1 RETURNING *;

