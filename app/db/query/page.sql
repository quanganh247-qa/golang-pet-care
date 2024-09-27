-- name: CreatePage :one
INSERT INTO pages (name, content, project_id, slug, created_at, updated_at)
VALUES ($1, $2, $3, $4,now(), now())
RETURNING *;

-- name: GetPages :many
SELECT * FROM pages ORDER BY created_at DESC;

-- name: UpdatePage :one
UPDATE pages
SET name = $1, content = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeletePage :one
DELETE FROM pages WHERE id = $1 RETURNING *;