-- name: CreatePetLog :one
INSERT INTO pet_logs (
    petid,
    datetime,
    title,
    notes
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPetLogsByPetID :many
SELECT * FROM pet_logs
WHERE petid = $1
ORDER BY datetime DESC LIMIT $2 OFFSET $3;

-- name: DeletePetLog :exec
DELETE FROM pet_logs
WHERE  log_id = $1;

-- name: UpdatePetLog :exec
UPDATE pet_logs
SET 
    title = $2,
    notes = $3
WHERE log_id = $1;

-- name: GetPetLogByID :one
SELECT pet_logs.petid, pet_logs.datetime, pet_logs.title, pet_logs.notes
FROM pet_logs
LEFT JOIN pets ON pet_logs.petid = pets.petid
WHERE pet_logs.petid = $1 AND pet_logs.log_id = $2 AND pets.is_active = true 
ORDER BY pet_logs.datetime DESC;