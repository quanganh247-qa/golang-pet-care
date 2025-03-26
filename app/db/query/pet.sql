-- name: CreatePet :one
INSERT INTO pets (
    name,
    type,
    breed,
    age,
    gender,
    healthnotes,
    weight,
    birth_date,
    username,
    microchip_number,
    last_checkup_date,
    is_active,
    data_image,
    original_image
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true, $12, $13
) RETURNING *;

-- name: GetPetByID :one
SELECT * FROM pets 
WHERE petid = $1;

-- name: ListPets :many
SELECT * FROM pets
WHERE is_active = true 
ORDER BY name LIMIT $1 OFFSET $2;

-- name: ListPetsByUsername :many
SELECT * FROM pets
WHERE username = $1 AND is_active = true
ORDER BY name LIMIT $2 OFFSET $3;

-- name: UpdatePet :exec
UPDATE pets
SET 
    name = $2,
    type = $3,
    breed = $4,
    age = $5,
    gender = $6,
    healthnotes = $7,
    weight = $8,
    birth_date = $9,
    microchip_number = $10,
    last_checkup_date = $11
WHERE petid = $1;

-- name: DeletePet :exec
DELETE FROM pets WHERE petid = $1;

-- name: SetPetInactive :exec
UPDATE pets
SET is_active = false
WHERE petid = $1;

-- name: UpdatePetAvatar :exec
UPDATE pets
SET 
    data_image = $2,
    original_image = $3
WHERE petid = $1;

-- name: GetAllPets :many
SELECT * FROM pets WHERE is_active is true;


-- name: GetPetProfileSummary :many
SELECT p.*, pt.*, v.* 
FROM pets AS p
LEFT JOIN pet_treatments AS pt ON p.petid = pt.pet_id
LEFT JOIN vaccinations AS v ON p.petid = v.petid
WHERE p.is_active = TRUE AND p.petid = $1;


-- name: GetPetDetailByUserID :one
SELECT p.*, u.full_name
FROM pets AS p
LEFT JOIN users AS u ON p.username = u.username
WHERE p.is_active = TRUE AND p.username = $1 AND p.name = $2;