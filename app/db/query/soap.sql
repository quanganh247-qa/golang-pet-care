-- name: CreateSOAP :one
INSERT INTO consultations (appointment_id, subjective)
VALUES ($1, $2)
RETURNING *;

-- name: GetSOAPByAppointmentID :one
SELECT * FROM consultations WHERE appointment_id = $1;


-- name: UpdateSOAP :one
UPDATE consultations SET subjective = $2, objective = $3, assessment = $4
WHERE appointment_id = $1 RETURNING *;



-- name: GetSOAPNote :one
SELECT 
    c.*,
    d.id as doctor_id,
    u.full_name as doctor_name,
    a.petid as pet_id,
    a.date as consultation_date
FROM consultations c
JOIN appointments a ON c.appointment_id = a.appointment_id
JOIN doctors d ON a.doctor_id = d.id
JOIN users u ON d.user_id = u.id
WHERE c.id = $1;

-- name: ListSOAPNotes :many
SELECT 
    c.*,
    d.id as doctor_id,
    u.full_name as doctor_name,
    a.petid as pet_id,
    a.date as consultation_date
FROM consultations c
JOIN appointments a ON c.appointment_id = a.appointment_id
JOIN doctors d ON a.doctor_id = d.id
JOIN users u ON d.user_id = u.id
WHERE a.petid = $1
ORDER BY c.created_at DESC;