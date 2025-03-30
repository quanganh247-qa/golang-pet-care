-- name: CreateSOAP :one
INSERT INTO consultations (appointment_id, subjective)
VALUES ($1, $2)
RETURNING *;

-- name: GetSOAPByAppointmentID :one
SELECT * FROM consultations WHERE appointment_id = $1;


-- name: UpdateSOAP :one
UPDATE consultations SET subjective = $2, objective = $3, assessment = $4
WHERE appointment_id = $1 RETURNING *;

