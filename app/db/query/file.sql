-- name: CreateFile :one
INSERT INTO files (
    file_name,
    file_path,
    file_size,
    file_type,
    pet_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetFileByID :one
SELECT * FROM files WHERE id = $1;

-- name: UpdateFile :one
UPDATE files SET
    file_name = $2,
    file_path = $3,
    file_size = $4,
    file_type = $5,
    pet_id = $6
WHERE id = $1 RETURNING *;

-- name: GetFiles :many
SELECT * FROM files WHERE pet_id = $1;
