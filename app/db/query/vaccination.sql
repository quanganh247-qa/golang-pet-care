-- name: CreateVaccination :one
INSERT INTO Vaccination (petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING vaccinationID, petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes;

-- name: GetVaccinationByID :one
SELECT vaccinationID, petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes
FROM Vaccination
WHERE vaccinationID = $1;

-- name: ListVaccinationsByPetID :many
SELECT vaccinationID, petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes
FROM Vaccination
WHERE petID = $1
ORDER BY dateAdministered DESC LIMIT $2 OFFSET $3;

-- name: UpdateVaccination :exec
UPDATE Vaccination
SET vaccineName = $2,
    dateAdministered = $3,
    nextDueDate = $4,
    vaccineProvider = $5,
    batchNumber = $6,
    notes = $7
WHERE vaccinationID = $1;

-- name: DeleteVaccination :exec
DELETE FROM Vaccination
WHERE vaccinationID = $1;