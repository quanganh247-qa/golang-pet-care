-- name: CreateTreatment :one
INSERT INTO pet_treatments (pet_id, diseases,doctor_id, name, type, start_date, end_date ,status, description, created_at)
VALUES ($1, $2, $3, $4, $5, $6 ,$7 , $8, $9, now()) RETURNING *;

-- name: GetTreatment :one
SELECT * FROM pet_treatments
WHERE id = $1 LIMIT 1;

-- name: UpdateTreatment :exec
UPDATE pet_treatments
SET diseases = $2, start_date = $3, end_date = $4, status = $5, description = $6
WHERE id = $1;

-- name: DeleteTreatment :exec
DELETE FROM pet_treatments
WHERE id = $1;

-- name: ListTreatmentsByPet :many
SELECT * FROM pet_treatments
WHERE pet_id = $1
ORDER BY start_date DESC
LIMIT $2 OFFSET $3;

-- name: CreateTreatmentPhase :one
INSERT INTO treatment_phases (treatment_id, phase_name, description, start_date, status, created_at)
VALUES ($1, $2, $3, $4, $5, now()) RETURNING *;

-- name: GetTreatmentPhase :one
SELECT * FROM treatment_phases
WHERE id = $1 LIMIT 1;

-- name: DeleteTreatmentPhase :exec
DELETE FROM treatment_phases
WHERE id = $1;

-- name: AssignMedicationToTreatmentPhase :one
INSERT INTO phase_medicines (phase_id, medicine_id, dosage, frequency, duration, notes, quantity, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now()) RETURNING *;


-- name: GetTreatmentsByPet :many
SELECT t.*
FROM pet_treatments t
WHERE t.pet_id = $1 LIMIT $2 OFFSET $3;

-- name: GetTreatmentPhasesByTreatment :many
SELECT *  FROM treatment_phases as tp
JOIN pet_treatments t ON t.id = tp.treatment_id
WHERE t.id = $1 LIMIT $2 OFFSET $3;

-- name: GetAllTreatmentPhasesByTreatmentID :many
SELECT * FROM treatment_phases
WHERE treatment_id = $1;

-- name: GetMedicationsByPhase :many
SELECT m.id, m.name, pm.dosage, pm.frequency, pm.duration, pm.notes, pm.quantity, pm.Created_at
FROM medicines m
JOIN phase_medicines pm ON m.id = pm.medicine_id
WHERE pm.phase_id = $1;

-- name: UpdateTreatmentPhaseStatus :exec
UPDATE treatment_phases
SET status = $2, updated_at = now()
WHERE id = $1;

-- name: GetActiveTreatments :many
SELECT t.id, pets.name AS pet_name, d.name AS disease, t.start_date, t.end_date, t.status
FROM pet_treatments t
JOIN pets ON t.pet_id = pets.petid
JOIN diseases d ON t.disease_id = d.id
WHERE t.status = 'ongoing' AND pets.petid = $1 LIMIT $2 OFFSET $3;

-- name: GetTreatmentProgress :many
SELECT tp.phase_name, tp.status, tp.start_date,COUNT(pm.medicine_id) AS num_medicines
FROM treatment_phases tp
LEFT JOIN phase_medicines pm ON tp.id = pm.phase_id
WHERE tp.id = $1
GROUP BY tp.id;

-- Assign Carprofen to the Initial Phase
-- name: AssignCarprofenToInitialPhase :exec
INSERT INTO phase_medicines (phase_id, medicine_id, dosage, frequency, duration, notes)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetMedicineByTreatmentID :many
SELECT pm.medicine_id, pm.dosage, pm.frequency, pm.duration, pm.quantity, pm.notes, pm.is_received
FROM phase_medicines pm
JOIN treatment_phases tp ON pm.phase_id = tp.id
WHERE tp.treatment_id = $1;


-- name: UpdateTreatmentStatus :exec
UPDATE pet_treatments
SET status = $2
WHERE id = $1;