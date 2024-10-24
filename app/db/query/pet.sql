-- name: CreatePet :one
INSERT INTO Pet (
    UserID,
    Name,
    Type,
    Breed,
    Age,
    Weight
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPet :one
SELECT * FROM Pet
WHERE PetID = $1 LIMIT 1;