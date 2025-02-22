-- name: CreateVaccination :one
INSERT INTO vaccinations (
    petID,
    vaccineName,
    dateAdministered,
    nextDueDate,
    vaccineProvider,
    batchNumber,
    notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetVaccinationByID :one
SELECT * FROM vaccinations 
WHERE vaccinationID = $1;

-- name: ListVaccinationsByPetID :many
SELECT * FROM vaccinations
WHERE petID = $1
ORDER BY dateAdministered DESC LIMIT $2 OFFSET $3;

-- name: UpdateVaccination :exec
UPDATE vaccinations
SET 
    vaccineName = $2,
    dateAdministered = $3,
    nextDueDate = $4,
    vaccineProvider = $5,
    batchNumber = $6,
    notes = $7
WHERE vaccinationID = $1;

-- name: DeleteVaccination :exec
DELETE FROM vaccinations
WHERE vaccinationID = $1;