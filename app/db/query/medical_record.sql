-- name: CreateMedicalRecord :one
INSERT INTO medical_records (pet_id,created_at,updated_at)
VALUES ($1,now(),now())
RETURNING *;

-- name: GetMedicalRecord :one
SELECT * FROM medical_records
WHERE id = $1 LIMIT 1;

-- name: UpdateMedicalRecord :exec
UPDATE medical_records
SET updated_at = NOW()
WHERE id = $1;

-- name: DeleteMedicalRecord :exec
DELETE FROM medical_records
WHERE id = $1;

-- name: CreateMedicalHistory :one
INSERT INTO medical_history(medical_record_id, condition, diagnosis_date, treatment, notes, created_at,updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: GetMedicalHistory :many
SELECT * FROM medical_history
WHERE medical_record_id = $1 LIMIT $2 OFFSET $3;

-- name: GetMedicalHistoryByID :one
SELECT * FROM medical_history
WHERE id = $1 LIMIT 1;

-- name: UpdateMedicalHistory :exec
UPDATE medical_history
SET condition = $2, diagnosis_date = $3, treatment = $4, notes = $5, updated_at = NOW()
WHERE id = $1;

-- name: DeleteMedicalHistory :exec
DELETE FROM medical_history
WHERE id = $1;
