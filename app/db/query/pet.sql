-- name: CreatePet :one
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
INSERT INTO Pet (username, Name, Type, Breed, Age, Weight, Gender, HealthNotes, data_image, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
=======
INSERT INTO Pet (username, Name, Type, Breed, Age, Weight, Gender, HealthNotes, data_image, original_image, birth_date, microchip_number, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, true)
>>>>>>> 9d28896 (image pet)
RETURNING *;

-- name: GetPetByID :one
SELECT * FROM Pet WHERE PetID = $1 AND is_active is true;

-- name: UpdatePet :exec
UPDATE Pet
SET Name = $2, Type = $3, Breed = $4, Age = $5, Weight = $6, Gender = $7, HealthNotes = $8, data_image = $9, is_active = $10
WHERE PetID = $1;

-- name: DeletePet :exec
DELETE FROM Pet WHERE PetID = $1;
>>>>>>> 0fb3f30 (user images)

-- name: ListPets :many
<<<<<<< HEAD
SELECT * FROM pets
WHERE is_active = true 
ORDER BY name LIMIT $1 OFFSET $2;
=======
SELECT * FROM Pet WHERE is_active is true ORDER BY PetID LIMIT $1 OFFSET $2;
>>>>>>> 3fdf0ad (updated pet status)

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
<<<<<<< HEAD
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
=======
UPDATE Pet SET is_active = $2 WHERE PetID = $1 AND is_active is true;
>>>>>>> 3fdf0ad (updated pet status)
