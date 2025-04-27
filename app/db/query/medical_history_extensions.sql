-- name: CreateExamination :one
INSERT INTO examinations (
    medical_history_id,
    exam_date,
    exam_type,
    findings,
    vet_notes,
    doctor_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetExaminationByID :one
SELECT * FROM examinations
WHERE id = $1 LIMIT 1;

-- name: ListExaminationsByMedicalHistory :many
SELECT * FROM examinations
WHERE medical_history_id = $1
ORDER BY exam_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdateExamination :one
UPDATE examinations
SET 
    exam_date = $2,
    exam_type = $3,
    findings = $4,
    vet_notes = $5,
    doctor_id = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteExamination :exec
DELETE FROM examinations
WHERE id = $1;

-- name: CreatePrescription :one
INSERT INTO prescriptions (
    medical_history_id,
    examination_id,
    prescription_date,
    doctor_id,
    notes
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPrescriptionByID :one
SELECT * FROM prescriptions
WHERE id = $1 LIMIT 1;

-- name: ListPrescriptionsByMedicalHistory :many
SELECT * FROM prescriptions
WHERE medical_history_id = $1
ORDER BY prescription_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePrescription :one
UPDATE prescriptions
SET 
    examination_id = $2,
    prescription_date = $3,
    doctor_id = $4,
    notes = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeletePrescription :exec
DELETE FROM prescriptions
WHERE id = $1;

-- name: CreatePrescriptionMedication :one
INSERT INTO prescription_medications (
    prescription_id,
    medicine_id,
    dosage,
    frequency,
    duration,
    instructions
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: ListPrescriptionMedications :many
SELECT * FROM prescription_medications
WHERE prescription_id = $1
ORDER BY id;

-- name: UpdatePrescriptionMedication :one
UPDATE prescription_medications
SET 
    medicine_id = $2,
    dosage = $3,
    frequency = $4,
    duration = $5,
    instructions = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeletePrescriptionMedication :exec
DELETE FROM prescription_medications
WHERE id = $1;

-- name: CreateTestResult :one
INSERT INTO test_results (
    medical_history_id,
    examination_id,
    test_type,
    test_date,
    results,
    interpretation,
    file_url,
    doctor_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetTestResultByID :one
SELECT * FROM test_results
WHERE id = $1 LIMIT 1;

-- name: ListTestResultsByMedicalHistory :many
SELECT * FROM test_results
WHERE medical_history_id = $1
ORDER BY test_date DESC
LIMIT $2 OFFSET $3;

-- name: ListTestResultsByExamination :many
SELECT * FROM test_results
WHERE examination_id = $1
ORDER BY test_date DESC;

-- name: UpdateTestResult :one
UPDATE test_results
SET 
    test_type = $2,
    test_date = $3,
    results = $4,
    interpretation = $5,
    file_url = $6,
    doctor_id = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTestResult :exec
DELETE FROM test_results
WHERE id = $1;

-- name: GetMedicalHistoryById :one
SELECT * FROM medical_history
WHERE id = $1 LIMIT 1;

-- name: GetMedicalRecordById :one
SELECT * FROM medical_records
WHERE id = $1 LIMIT 1;

-- name: GetCompleteMedicalHistory :one
SELECT get_pet_medical_history_summary($1) as medical_history;