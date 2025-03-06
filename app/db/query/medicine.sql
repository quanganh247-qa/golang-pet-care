
-- name: CreateMedicine :one
INSERT INTO medicines (name, description, usage, dosage, frequency, duration, side_effects, expiration_date, quantity)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListMedicinesByPet :many
SELECT 
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
    pm.dosage,
    pm.frequency,
    pm.duration,
    pm.notes AS medicine_notes,
    pt.start_date AS treatment_start_date,
    pt.end_date AS treatment_end_date,
    pt.status AS treatment_status
FROM 
    pet_treatments pt
JOIN 
    treatment_phases tp ON pt.disease_id = tp.disease_id
JOIN 
    phase_medicines pm ON tp.id = pm.phase_id
JOIN 
    medicines m ON pm.medicine_id = m.id
WHERE 
    pt.pet_id = $1 and pt.status = $2 -- Replace with the specific pet_id
ORDER BY 
    tp.start_date, pm.medicine_id LIMIT $3 OFFSET $4;

-- name: CreateAllergy :one
INSERT INTO allergies (medical_record_id, allergen, severity, reaction, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllergies :many
SELECT * FROM allergies
WHERE medical_record_id = $1;

-- name: UpdateAllergy :exec
UPDATE allergies
SET allergen = $2, severity = $3, reaction = $4, notes = $5, updated_at = NOW()
WHERE id = $1;

-- name: DeleteAllergy :exec
DELETE FROM Allergies
WHERE id = $1;


-- name: GetMedicineByID :one
SELECT * FROM medicines
WHERE id = $1 LIMIT 1;