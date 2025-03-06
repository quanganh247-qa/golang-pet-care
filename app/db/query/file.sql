-- name: CreateFile :one
INSERT INTO files (
    file_name,
    file_path,
    file_size,
    file_type
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetFileByID :one
SELECT * FROM files WHERE id = $1;

-- name: UpdateFile :one
UPDATE files SET
    file_name = $2,
    file_path = $3,
    file_size = $4,
    file_type = $5
WHERE id = $1 RETURNING *;
