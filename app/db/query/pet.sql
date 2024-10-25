-- name: CreatePet :one
INSERT INTO Pet (username, Name, Type, Breed, Age, Weight, Gender, HealthNotes, ProfileImage)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetPetByID :one
SELECT * FROM Pet WHERE PetID = $1;

-- name: UpdatePet :exec
UPDATE Pet
SET Name = $2, Type = $3, Breed = $4, Age = $5, Weight = $6, Gender = $7, HealthNotes = $8, ProfileImage = $9
WHERE PetID = $1;

-- name: DeletePet :exec
DELETE FROM Pet WHERE PetID = $1;

-- name: ListPets :many
SELECT * FROM Pet ORDER BY PetID LIMIT $1 OFFSET $2;
