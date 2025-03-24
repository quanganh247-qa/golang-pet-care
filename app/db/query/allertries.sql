-- name: CreatePetAllergy :one
INSERT INTO pet_allergies (
    pet_id,
    type,
    detail
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPetAllergy :one
SELECT * FROM pet_allergies
WHERE id = $1 LIMIT 1;

-- name: ListPetAllergies :many
SELECT * FROM pet_allergies
WHERE pet_id = $1 LIMIT $2 OFFSET $3;

-- name: UpdatePetAllergy :one
UPDATE pet_allergies
SET type = $2,
    detail = $3
WHERE id = $1
RETURNING *;

-- name: DeletePetAllergy :exec
DELETE FROM pet_allergies
<<<<<<< HEAD
WHERE id = $1;
=======
WHERE id = $1;
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
