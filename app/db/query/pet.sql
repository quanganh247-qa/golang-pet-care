-- name: CreatePet :one
INSERT INTO Pet (username, Name, Type, Breed, Age, Weight, Gender, HealthNotes, ProfileImage, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetPetByID :one
SELECT * FROM Pet WHERE PetID = $1;

-- name: UpdatePet :exec
UPDATE Pet
SET Name = $2, Type = $3, Breed = $4, Age = $5, Weight = $6, Gender = $7, HealthNotes = $8, ProfileImage = $9, is_active = $10
WHERE PetID = $1;

-- name: DeletePet :exec
DELETE FROM Pet WHERE PetID = $1;

-- name: ListPets :many
SELECT * FROM Pet ORDER BY PetID LIMIT $1 OFFSET $2;

-- name: ListPetsByUsername :many
SELECT * FROM Pet WHERE username = $1 ORDER BY PetID LIMIT $2 OFFSET $3;

-- name: SetPetInactive :exec
UPDATE Pet SET is_active = $2 WHERE PetID = $1;
