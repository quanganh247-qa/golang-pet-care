-- name: CreateVaccination :one
INSERT INTO Vaccination (petid, vaccine_name, date_administered, next_due_date, vaccine_provider, batch_number, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING vaccinationid, petid, vaccine_name, date_administered, next_due_date, vaccine_provider, batch_number, notes;

-- name: GetVaccinationByID :one
SELECT vaccinationid, petid, vaccine_name, date_administered, next_due_date, vaccine_provider, batch_number, notes
FROM Vaccination
WHERE vaccinationid = $1;

-- name: ListVaccinationsByPetID :many
SELECT vaccinationid, petid, vaccine_name, date_administered, next_due_date, vaccine_provider, batch_number, notes
FROM Vaccination
WHERE petid = $1
ORDER BY date_administered DESC;

-- name: UpdateVaccination :exec
UPDATE Vaccination
SET vaccine_name = $2,
    date_administered = $3,
    next_due_date = $4,
    vaccine_provider = $5,
    batch_number = $6,
    notes = $7
WHERE vaccinationid = $1;

-- name: DeleteVaccination :exec
DELETE FROM Vaccination
WHERE vaccinationid = $1;
