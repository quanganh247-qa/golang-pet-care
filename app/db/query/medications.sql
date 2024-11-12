-- name: InsertMedicine :one
INSERT INTO Medications (pet_id, medication_name, dosage, frequency, start_date, end_date, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetMedicinesByID :one
SELECT * FROM Medications WHERE medication_id = $1;

-- -- name: UpdateMedicine :exec
-- UPDATE Medications
-- SET medication_name = $2, dosage = $3, frequency = $4, start_date = $5, end_date = $6, notes = $7
-- WHERE medication_id = $1;

-- name: UpdateMedicine :one
UPDATE Medications 
SET 
  medication_name = COALESCE(sqlc.narg(medication_name), medication_name),       
  dosage =  COALESCE(sqlc.narg(dosage),dosage),
  frequency = COALESCE(sqlc.narg(frequency),frequency),
  start_date = COALESCE(sqlc.narg(start_date),start_date),
  end_date = COALESCE(sqlc.narg(end_date),end_date),
  notes = COALESCE(sqlc.narg(notes),notes)
WHERE
  medication_id = sqlc.arg(medication_id)
RETURNING *;



-- name: GetAllMedicinesByPet :many
SELECT * FROM Medications where pet_id = $3  ORDER BY medication_id LIMIT $1 OFFSET $2 ;

