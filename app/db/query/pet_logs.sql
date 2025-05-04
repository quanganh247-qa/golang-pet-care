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
SELECT pet_logs.log_id, pet_logs.petid, pet_logs.datetime, pet_logs.title, pet_logs.notes, pets.name, pets.type, pets.breed
FROM pet_logs
LEFT JOIN pets ON pet_logs.petid = pets.petid
WHERE pet_logs.petid = $1 AND pets.is_active = true
ORDER BY pet_logs.datetime DESC LIMIT $2 OFFSET $3;

-- name: DeletePetLog :exec
DELETE FROM pet_logs
WHERE  log_id = $1;

-- name: UpdatePetLog :exec
UPDATE pet_logs
SET 
    title = $2,
    notes = $3,
    datetime = $4
WHERE log_id = $1;

-- name: GetPetLogByID :one
SELECT pet_logs.log_id, pet_logs.petid, pet_logs.datetime, pet_logs.title, pet_logs.notes, pets.name, pets.type, pets.breed
FROM pet_logs
LEFT JOIN pets ON pet_logs.petid = pets.petid
WHERE pet_logs.log_id = $1 AND pets.is_active = true 
ORDER BY pet_logs.datetime DESC;



-- name: GetAllPetLogsByUsername :many
SELECT 
    pl.log_id,
    pl.petid,
    p.name AS pet_name,
    p.type AS pet_type,
    p.breed AS pet_breed,
    pl.datetime,
    pl.title,
    pl.notes
FROM 
    pet_logs pl
JOIN 
    pets p ON pl.petid = p.petid
WHERE 
    p.username = $1
ORDER BY 
    pl.datetime DESC
LIMIT $2 OFFSET $3;

-- name: CountAllPetLogsByUsername :one
SELECT 
    COUNT(*)
FROM 
    pet_logs pl
JOIN 
    pets p ON pl.petid = p.petid
WHERE 
    p.username = $1;

-- name: GetDetailsPetLogByID :one
SELECT 
    pl.log_id,
    pl.petid,
    p.name AS pet_name,
    p.type AS pet_type,
    p.breed AS pet_breed,
    pl.datetime,
    pl.title,
    pl.notes
FROM
    pet_logs pl
JOIN
    pets p ON pl.petid = p.petid
WHERE 
    pl.log_id = $1 AND p.username = $2;
