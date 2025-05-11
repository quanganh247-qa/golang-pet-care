-- name: AddPetWeightRecord :one
INSERT INTO pet_weight_history (
    pet_id,
    weight_kg,
    notes
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdatePetCurrentWeight :exec
UPDATE pets
SET weight = $2, updated_at = NOW()
WHERE petid = $1;

-- name: GetPetWeightHistory :many
SELECT * FROM pet_weight_history
WHERE pet_id = $1
ORDER BY recorded_at DESC
LIMIT $2 OFFSET $3;

-- name: GetLatestPetWeight :one
SELECT * FROM pet_weight_history
WHERE pet_id = $1
ORDER BY recorded_at DESC
LIMIT 1;

-- name: CountPetWeightRecords :one
SELECT COUNT(*) FROM pet_weight_history
WHERE pet_id = $1;

-- name: DeletePetWeightRecord :exec
DELETE FROM pet_weight_history
WHERE id = $1 AND pet_id = $2; 


